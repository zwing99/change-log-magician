package filesystem

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/zwing99/change-log-magician/internal/changelog"
)

func TestRepositoryInitializesAndAuthorsFragment(t *testing.T) {
	t.Parallel()
	directory := t.TempDir()
	repo := Repository{Root: directory}
	if err := repo.Initialize(); err != nil {
		t.Fatal(err)
	}
	if err := repo.Initialize(); err == nil {
		t.Fatal("expected safe overwrite failure")
	}
	path, err := repo.CreateFragment(changelog.Minor, "Shell completions", time.Date(2026, 9, 12, 1, 2, 3, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	if err := repo.AppendEntry(path, changelog.Added, "Generate completion files"); err != nil {
		t.Fatal(err)
	}
	fragments, err := repo.ActiveFragments()
	if err != nil || len(fragments) != 1 || fragments[0].Level != changelog.Minor {
		t.Fatalf("fragments: %#v, %v", fragments, err)
	}
	if _, err := os.Stat(filepath.Join(directory, ChangelogFile)); err != nil {
		t.Fatal(err)
	}
}

func TestSelectFragmentRejectsAmbiguity(t *testing.T) {
	t.Parallel()
	directory := t.TempDir()
	repo := Repository{Root: directory}
	if err := os.Mkdir(filepath.Join(directory, FragmentsDir), 0o755); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"a.md", "b.md"} {
		if err := os.WriteFile(filepath.Join(directory, FragmentsDir, name), []byte("# [patch]\n\n## Fixed\n\n- x\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := repo.SelectFragment(""); err == nil {
		t.Fatal("expected ambiguity error")
	}
}
