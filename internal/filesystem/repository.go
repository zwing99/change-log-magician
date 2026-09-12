// Package filesystem provides the small filesystem boundary used by CLM.
package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/zwing99/change-log-magician/internal/changelog"
)

const (
	ChangelogFile = "CHANGELOG.md"
	FragmentsDir  = "changelog.d"
	ArchiveDir    = "changelog.d/archive"
	ConfigFile    = "clm.toml"
)

// Repository maps logical CLM paths onto one working directory.
type Repository struct{ Root string }

func (r Repository) path(parts ...string) string {
	return filepath.Join(append([]string{r.Root}, parts...)...)
}

// Initialize creates the default CLM files without replacing existing files.
func (r Repository) Initialize() error {
	if err := os.MkdirAll(r.path(ArchiveDir), 0o755); err != nil {
		return fmt.Errorf("create archive directory: %w", err)
	}
	if err := writeNew(r.path(ChangelogFile), changelog.RenderRoot(changelog.NewRoot())); err != nil {
		return err
	}
	const config = "[paths]\nchangelog = \"CHANGELOG.md\"\nfragments = \"changelog.d\"\narchive = \"changelog.d/archive\"\n\n[tags]\nrelease_prefix = \"v\"\nclm_prefix = \"clm/v\"\n"
	return writeNew(r.path(ConfigFile), config)
}

func writeNew(path, content string) error {
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("%s already exists; refusing to overwrite it", filepath.Base(path))
	} else if !os.IsNotExist(err) {
		return err
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

// ActiveFragments returns parsed direct fragment files and ignores archives.
func (r Repository) ActiveFragments() ([]changelog.Fragment, error) {
	entries, err := os.ReadDir(r.path(FragmentsDir))
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read fragment directory: %w", err)
	}
	fragments := make([]changelog.Fragment, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		path := r.path(FragmentsDir, entry.Name())
		contents, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read fragment %s: %w", path, err)
		}
		fragment, err := changelog.ParseFragment(filepath.ToSlash(filepath.Join(FragmentsDir, entry.Name())), string(contents))
		if err != nil {
			return nil, err
		}
		fragments = append(fragments, fragment)
	}
	sort.Slice(fragments, func(i, j int) bool { return fragments[i].Path < fragments[j].Path })
	return fragments, nil
}

// CreateFragment creates a level-specific source fragment.
func (r Repository) CreateFragment(level changelog.Level, name string, now time.Time) (string, error) {
	if _, err := changelog.ParseLevel(string(level)); err != nil {
		return "", err
	}
	if err := os.MkdirAll(r.path(FragmentsDir), 0o755); err != nil {
		return "", err
	}
	slug := slugify(name)
	if slug == "" {
		slug = "change"
	}
	base := now.UTC().Format("20060102-150405") + "-" + slug
	for sequence := 0; ; sequence++ {
		filename := base + ".md"
		if sequence > 0 {
			filename = fmt.Sprintf("%s-%d.md", base, sequence+1)
		}
		path := r.path(FragmentsDir, filename)
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			continue
		}
		content := fmt.Sprintf("# [%s]\n", level)
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			return "", err
		}
		return filepath.ToSlash(filepath.Join(FragmentsDir, filename)), nil
	}
}

// SelectFragment resolves an explicit target or exactly one active fragment.
func (r Repository) SelectFragment(target string) (string, error) {
	if strings.TrimSpace(target) != "" {
		path := target
		if !filepath.IsAbs(path) {
			path = r.path(target)
			if !strings.HasPrefix(filepath.Clean(path), filepath.Clean(r.path(FragmentsDir))+string(os.PathSeparator)) {
				path = r.path(FragmentsDir, target)
			}
		}
		if _, err := os.Stat(path); err != nil {
			return "", fmt.Errorf("selected fragment %s: %w", target, err)
		}
		return path, nil
	}
	paths, err := r.activeFragmentPaths()
	if err != nil {
		return "", err
	}
	if len(paths) == 0 {
		return "", fmt.Errorf("no active changelog fragment exists; create one with clm new-patch, new-minor, new-major, or new-unreleased")
	}
	if len(paths) > 1 {
		return "", fmt.Errorf("multiple active fragments exist; select one with --change <path>")
	}
	return r.path(paths[0]), nil
}

func (r Repository) activeFragmentPaths() ([]string, error) {
	entries, err := os.ReadDir(r.path(FragmentsDir))
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read fragment directory: %w", err)
	}
	paths := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
			paths = append(paths, filepath.ToSlash(filepath.Join(FragmentsDir, entry.Name())))
		}
	}
	sort.Strings(paths)
	return paths, nil
}

// AppendEntry updates a valid fragment through the canonical renderer.
func (r Repository) AppendEntry(target string, category changelog.Category, entry string) error {
	if strings.TrimSpace(entry) == "" {
		return fmt.Errorf("changelog entry cannot be empty")
	}
	path, err := r.SelectFragment(target)
	if err != nil {
		return err
	}
	contents, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	fragment, err := changelog.ParseFragment(path, string(contents))
	if err != nil {
		// New fragments intentionally start empty; permit their first authored entry.
		level, levelErr := parseEmptyFragment(path, string(contents))
		if levelErr != nil {
			return err
		}
		fragment = changelog.Fragment{Path: path, Level: level, Entries: map[changelog.Category][]string{}}
	}
	fragment.Entries[category] = append(fragment.Entries[category], strings.TrimSpace(entry))
	return os.WriteFile(path, []byte(changelog.RenderFragment(fragment)), 0o644)
}

func parseEmptyFragment(path, source string) (changelog.Level, error) {
	line := strings.TrimSpace(source)
	if !strings.HasPrefix(line, "# [") || !strings.HasSuffix(line, "]") {
		return "", fmt.Errorf("%s is not a valid empty fragment", path)
	}
	return changelog.ParseLevel(strings.TrimSuffix(strings.TrimPrefix(line, "# ["), "]"))
}

func slugify(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	var builder strings.Builder
	dash := false
	for _, runeValue := range value {
		if (runeValue >= 'a' && runeValue <= 'z') || (runeValue >= '0' && runeValue <= '9') {
			builder.WriteRune(runeValue)
			dash = false
		} else if !dash && builder.Len() > 0 {
			builder.WriteByte('-')
			dash = true
		}
	}
	return strings.Trim(builder.String(), "-")
}
