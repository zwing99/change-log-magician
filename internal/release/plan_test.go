package release

import (
	"strings"
	"testing"
	"time"

	"github.com/zwing99/change-log-magician/internal/changelog"
)

func TestBuildPlanPromotesUnreleased(t *testing.T) {
	t.Parallel()
	root := changelog.NewRoot()
	root.Unreleased().Entries[changelog.Added] = []string{"Already visible"}
	fragments := []changelog.Fragment{{Path: "changelog.d/f.md", Level: changelog.Minor, Entries: map[changelog.Category][]string{changelog.Fixed: {"Bug fix"}}}}
	plan, err := BuildPlan(root, fragments, Version{}, time.Date(2026, 9, 12, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	if plan.Version == nil || plan.Version.String() != "0.1.0" {
		t.Fatalf("unexpected version: %#v", plan.Version)
	}
	for _, expected := range []string{"## [Unreleased]", "## [0.1.0] - 2026-09-12", "- Already visible", "- Bug fix"} {
		if !strings.Contains(plan.Rendered, expected) {
			t.Fatalf("missing %q", expected)
		}
	}
}

func TestBuildPlanUnreleasedOnly(t *testing.T) {
	t.Parallel()
	root := changelog.NewRoot()
	fragments := []changelog.Fragment{{Path: "changelog.d/f.md", Level: changelog.Unreleased, Entries: map[changelog.Category][]string{changelog.Changed: {"Internal change"}}}}
	plan, err := BuildPlan(root, fragments, Version{}, time.Date(2026, 9, 12, 3, 4, 5, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	if plan.Version != nil || !strings.Contains(plan.ArchiveSegment, "unreleased/") || !strings.Contains(plan.Rendered, "- Internal change") {
		t.Fatalf("unexpected plan: %#v", plan)
	}
}
