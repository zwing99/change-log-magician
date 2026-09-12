package release

import (
	"testing"

	"github.com/zwing99/change-log-magician/internal/changelog"
)

func TestNextVersion(t *testing.T) {
	t.Parallel()
	base := Version{}
	cases := []struct {
		level changelog.Level
		want  string
	}{{changelog.Patch, "0.0.1"}, {changelog.Minor, "0.1.0"}, {changelog.Major, "1.0.0"}}
	for _, test := range cases {
		got, err := NextVersion(base, test.level)
		if err != nil {
			t.Fatal(err)
		}
		if got.String() != test.want {
			t.Fatalf("%s: got %s, want %s", test.level, got, test.want)
		}
	}
}

func TestHighestLevelAndPrerelease(t *testing.T) {
	t.Parallel()
	fragments := []changelog.Fragment{{Level: changelog.Patch}, {Level: changelog.Major}, {Level: changelog.Unreleased}}
	level, ok := HighestLevel(fragments)
	if !ok || level != changelog.Major {
		t.Fatalf("got %s %v", level, ok)
	}
	got, err := Prerelease(Version{Major: 1, Minor: 4}, "pr.123.abcdef0")
	if err != nil || got != "1.4.0-pr.123.abcdef0" {
		t.Fatalf("got %q, %v", got, err)
	}
	if _, err := Prerelease(Version{}, "bad id"); err == nil {
		t.Fatal("expected invalid prerelease id")
	}
}
