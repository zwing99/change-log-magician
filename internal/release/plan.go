package release

import (
	"fmt"
	"time"

	"github.com/zwing99/change-log-magician/internal/changelog"
)

// Plan is the pure result shared by check, preview, prerelease, and cut.
type Plan struct {
	Fragments      []changelog.Fragment
	Level          changelog.Level
	Version        *Version
	Date           string
	Rendered       string
	ArchiveSegment string
}

// BuildPlan merges fragments into a parsed root changelog without side effects.
func BuildPlan(root changelog.Root, fragments []changelog.Fragment, base Version, now time.Time) (Plan, error) {
	if len(fragments) == 0 {
		return Plan{}, fmt.Errorf("no active changelog fragments to cut")
	}
	level, versioned := HighestLevel(fragments)
	plan := Plan{Fragments: fragments, Level: level, Date: now.UTC().Format("2006-01-02")}
	entries := make(map[changelog.Category][]string)
	for _, fragment := range fragments {
		for category, values := range fragment.Entries {
			entries[category] = append(entries[category], values...)
		}
	}
	if !versioned {
		unreleased := root.Unreleased()
		for category, values := range entries {
			unreleased.Entries[category] = append(unreleased.Entries[category], values...)
		}
		plan.ArchiveSegment = "unreleased/" + now.UTC().Format("20060102-150405")
		plan.Rendered = changelog.RenderRoot(root)
		return plan, nil
	}
	next, err := NextVersion(base, level)
	if err != nil {
		return Plan{}, err
	}
	plan.Version = &next
	plan.ArchiveSegment = "v" + next.String()
	unreleased := root.Unreleased()
	releaseEntries := changelog.CloneEntries(unreleased.Entries)
	for category, values := range entries {
		releaseEntries[category] = append(releaseEntries[category], values...)
	}
	root.Releases = append([]changelog.Release{{Entries: map[changelog.Category][]string{}}}, root.Releases...)
	root.Releases[1] = changelog.Release{Version: next.String(), Date: plan.Date, Entries: releaseEntries}
	plan.Rendered = changelog.RenderRoot(root)
	return plan, nil
}
