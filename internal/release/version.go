// Package release contains pure version and cut planning rules.
package release

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/zwing99/change-log-magician/internal/changelog"
)

var versionPattern = regexp.MustCompile(`^v?(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)$`)
var prereleasePattern = regexp.MustCompile(`^[0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*$`)

// Version is a stable SemVer release without a leading v.
type Version struct{ Major, Minor, Patch int }

func (v Version) String() string { return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch) }

// ParseVersion accepts stable vX.Y.Z or X.Y.Z tags.
func ParseVersion(value string) (Version, error) {
	match := versionPattern.FindStringSubmatch(strings.TrimSpace(value))
	if match == nil {
		return Version{}, fmt.Errorf("invalid stable version %q", value)
	}
	values := [3]int{}
	for i := range values {
		parsed, err := strconv.Atoi(match[i+1])
		if err != nil {
			return Version{}, fmt.Errorf("parse version %q: %w", value, err)
		}
		values[i] = parsed
	}
	return Version{values[0], values[1], values[2]}, nil
}

// NextVersion returns the SemVer increment requested by level.
func NextVersion(base Version, level changelog.Level) (Version, error) {
	switch level {
	case changelog.Major:
		return Version{Major: base.Major + 1}, nil
	case changelog.Minor:
		return Version{Major: base.Major, Minor: base.Minor + 1}, nil
	case changelog.Patch:
		return Version{Major: base.Major, Minor: base.Minor, Patch: base.Patch + 1}, nil
	default:
		return Version{}, fmt.Errorf("%s does not produce a stable version", level)
	}
}

// HighestLevel returns the highest versioned level, if any.
func HighestLevel(fragments []changelog.Fragment) (changelog.Level, bool) {
	rank := map[changelog.Level]int{changelog.Patch: 1, changelog.Minor: 2, changelog.Major: 3}
	var highest changelog.Level
	for _, fragment := range fragments {
		if rank[fragment.Level] > rank[highest] {
			highest = fragment.Level
		}
	}
	return highest, highest != ""
}

// Prerelease derives a valid SemVer prerelease version from a stable plan.
func Prerelease(version Version, identifier string) (string, error) {
	identifier = strings.TrimSpace(identifier)
	if !prereleasePattern.MatchString(identifier) {
		return "", fmt.Errorf("invalid prerelease identifier %q", identifier)
	}
	return version.String() + "-" + identifier, nil
}
