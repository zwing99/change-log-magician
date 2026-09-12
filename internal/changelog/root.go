package changelog

import (
	"fmt"
	"strings"
)

// Release is a rendered root-changelog section. Version is empty for Unreleased.
type Release struct {
	Version string
	Date    string
	Entries map[Category][]string
}

// Root is the CLM-managed portion of a root changelog plus preserved outer text.
type Root struct {
	Preamble string
	Releases []Release
	Trailer  string
}

// NewRoot returns the default Keep a Changelog root document.
func NewRoot() Root {
	return Root{Preamble: "# Changelog\n\nAll notable changes to this project will be documented in this file.\n\nThis changelog is maintained with [Change Log Magician](https://github.com/zwing99/change-log-magician) (`clm`). Add one fragment for every change, then let CLM assemble releases.\n\nHelpful commands:\n\n- `clm new-patch <name>` creates a release-note fragment.\n- `clm added -m \"Description\"` adds a categorized entry.\n- `clm cut --check` validates the next release without changing files.\n", Releases: []Release{{Entries: map[Category][]string{}}}}
}

// ParseRoot reads CLM release sections while retaining other root content.
func ParseRoot(source string) (Root, error) {
	lines := strings.Split(strings.ReplaceAll(source, "\r\n", "\n"), "\n")
	root := Root{}
	start := -1
	for i, line := range lines {
		if line == "## [Unreleased]" || isVersionHeader(line) {
			start = i
			break
		}
	}
	if start == -1 {
		return Root{}, fmt.Errorf("CHANGELOG.md has no CLM release section")
	}
	root.Preamble = strings.TrimRight(strings.Join(lines[:start], "\n"), "\n") + "\n"
	i := start
	for i < len(lines) {
		if !(lines[i] == "## [Unreleased]" || isVersionHeader(lines[i])) {
			root.Trailer = strings.TrimLeft(strings.Join(lines[i:], "\n"), "\n")
			break
		}
		release, next, err := parseRelease(lines, i)
		if err != nil {
			return Root{}, err
		}
		root.Releases = append(root.Releases, release)
		i = next
		for i < len(lines) && strings.TrimSpace(lines[i]) == "" {
			i++
		}
	}
	if len(root.Releases) == 0 || root.Releases[0].Version != "" {
		return Root{}, fmt.Errorf("CHANGELOG.md must begin its managed releases with ## [Unreleased]")
	}
	return root, nil
}

func isVersionHeader(line string) bool {
	return strings.HasPrefix(line, "## [") && strings.Contains(line, "] - ")
}

func parseRelease(lines []string, start int) (Release, int, error) {
	header := lines[start]
	release := Release{Entries: map[Category][]string{}}
	if header != "## [Unreleased]" {
		parts := strings.SplitN(strings.TrimSuffix(strings.TrimPrefix(header, "## ["), ""), "] - ", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return Release{}, 0, fmt.Errorf("invalid release heading %q", header)
		}
		release.Version, release.Date = parts[0], parts[1]
	}
	current := Category("")
	i := start + 1
	for i < len(lines) {
		line := strings.TrimSpace(lines[i])
		if line == "## [Unreleased]" || isVersionHeader(lines[i]) {
			break
		}
		if strings.HasPrefix(line, "## ") {
			break // start of manually maintained trailing content
		}
		if strings.HasPrefix(line, "### ") {
			category, err := ParseCategory(strings.TrimPrefix(line, "### "))
			if err != nil {
				return Release{}, 0, fmt.Errorf("%s: unsupported category in root changelog: %w", header, err)
			}
			current = category
		} else if strings.HasPrefix(line, "- ") {
			if current == "" {
				return Release{}, 0, fmt.Errorf("%s: entry appears before a category", header)
			}
			release.Entries[current] = append(release.Entries[current], strings.TrimSpace(strings.TrimPrefix(line, "- ")))
		}
		i++
	}
	return release, i, nil
}

// RenderRoot writes CLM-managed release sections and preserves other text.
func RenderRoot(root Root) string {
	var builder strings.Builder
	builder.WriteString(strings.TrimRight(root.Preamble, "\n"))
	builder.WriteString("\n\n")
	for index, release := range root.Releases {
		if index > 0 {
			builder.WriteString("\n")
		}
		if release.Version == "" {
			builder.WriteString("## [Unreleased]\n")
		} else {
			fmt.Fprintf(&builder, "## [%s] - %s\n", release.Version, release.Date)
		}
		for _, category := range Categories() {
			entries := release.Entries[category]
			if len(entries) == 0 {
				continue
			}
			fmt.Fprintf(&builder, "\n### %s\n\n", category)
			for _, entry := range entries {
				fmt.Fprintf(&builder, "- %s\n", entry)
			}
		}
	}
	if strings.TrimSpace(root.Trailer) != "" {
		builder.WriteString("\n")
		builder.WriteString(strings.TrimLeft(root.Trailer, "\n"))
	}
	return strings.TrimRight(builder.String(), "\n") + "\n"
}

// Unreleased returns the active root release, creating it where necessary.
func (r *Root) Unreleased() *Release {
	if len(r.Releases) == 0 || r.Releases[0].Version != "" {
		r.Releases = append([]Release{{Entries: map[Category][]string{}}}, r.Releases...)
	}
	return &r.Releases[0]
}
