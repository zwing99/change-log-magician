// Package changelog defines the Markdown data model used by CLM.
package changelog

import (
	"fmt"
	"strings"
)

// Level is the semantic versioning intent of a fragment.
type Level string

const (
	Major      Level = "major"
	Minor      Level = "minor"
	Patch      Level = "patch"
	Unreleased Level = "unreleased"
)

// ParseLevel validates a fragment level.
func ParseLevel(value string) (Level, error) {
	level := Level(strings.ToLower(strings.TrimSpace(value)))
	switch level {
	case Major, Minor, Patch, Unreleased:
		return level, nil
	default:
		return "", fmt.Errorf("invalid release level %q (expected major, minor, patch, or unreleased)", value)
	}
}

// IsVersioned reports whether a level produces a stable version.
func (l Level) IsVersioned() bool {
	return l == Major || l == Minor || l == Patch
}

// Category is a Keep a Changelog section.
type Category string

const (
	Added      Category = "Added"
	Changed    Category = "Changed"
	Deprecated Category = "Deprecated"
	Removed    Category = "Removed"
	Fixed      Category = "Fixed"
	Security   Category = "Security"
)

var categoryOrder = []Category{Added, Changed, Deprecated, Removed, Fixed, Security}

// Categories returns categories in Keep a Changelog display order.
func Categories() []Category {
	return append([]Category(nil), categoryOrder...)
}

// ParseCategory validates a category name case-insensitively.
func ParseCategory(value string) (Category, error) {
	for _, category := range categoryOrder {
		if strings.EqualFold(string(category), strings.TrimSpace(value)) {
			return category, nil
		}
	}
	return "", fmt.Errorf("unsupported changelog category %q", value)
}

// Fragment is one active or archived changelog source file.
type Fragment struct {
	Path     string
	Level    Level
	Entries  map[Category][]string
	RawOrder []Category
}

// HasEntries reports whether the fragment contains any release-note entries.
func (f Fragment) HasEntries() bool {
	for _, entries := range f.Entries {
		if len(entries) > 0 {
			return true
		}
	}
	return false
}

// CloneEntries returns a deep copy suitable for rendering.
func CloneEntries(entries map[Category][]string) map[Category][]string {
	copyOf := make(map[Category][]string, len(entries))
	for category, values := range entries {
		copyOf[category] = append([]string(nil), values...)
	}
	return copyOf
}
