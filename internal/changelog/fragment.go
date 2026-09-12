package changelog

import (
	"bufio"
	"fmt"
	"strings"
)

// ParseFragment parses CLM's intentionally constrained, human-editable format.
func ParseFragment(path, source string) (Fragment, error) {
	fragment := Fragment{Path: path, Entries: make(map[Category][]string)}
	scanner := bufio.NewScanner(strings.NewReader(source))
	lineNumber := 0
	var current *Category
	seenLevel := false
	seenCategories := make(map[Category]bool)

	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "# ") {
			if seenLevel {
				return Fragment{}, fmt.Errorf("%s:%d: fragment has more than one level header", path, lineNumber)
			}
			header := strings.TrimSpace(strings.TrimPrefix(line, "# "))
			if !strings.HasPrefix(header, "[") || !strings.HasSuffix(header, "]") {
				return Fragment{}, fmt.Errorf("%s:%d: expected level header such as # [patch]", path, lineNumber)
			}
			level, err := ParseLevel(strings.TrimSuffix(strings.TrimPrefix(header, "["), "]"))
			if err != nil {
				return Fragment{}, fmt.Errorf("%s:%d: %w", path, lineNumber, err)
			}
			fragment.Level, seenLevel = level, true
			continue
		}
		if strings.HasPrefix(line, "## ") {
			category, err := ParseCategory(strings.TrimSpace(strings.TrimPrefix(line, "## ")))
			if err != nil {
				return Fragment{}, fmt.Errorf("%s:%d: %w", path, lineNumber, err)
			}
			if seenCategories[category] {
				return Fragment{}, fmt.Errorf("%s:%d: duplicate %s category", path, lineNumber, category)
			}
			seenCategories[category] = true
			fragment.RawOrder = append(fragment.RawOrder, category)
			current = &category
			continue
		}
		if strings.HasPrefix(line, "- ") {
			if current == nil {
				return Fragment{}, fmt.Errorf("%s:%d: entry must follow a category heading", path, lineNumber)
			}
			entry := strings.TrimSpace(strings.TrimPrefix(line, "- "))
			if entry == "" {
				return Fragment{}, fmt.Errorf("%s:%d: changelog entry cannot be empty", path, lineNumber)
			}
			fragment.Entries[*current] = append(fragment.Entries[*current], entry)
			continue
		}
		return Fragment{}, fmt.Errorf("%s:%d: unsupported content; use headings and bullet entries", path, lineNumber)
	}
	if err := scanner.Err(); err != nil {
		return Fragment{}, fmt.Errorf("read %s: %w", path, err)
	}
	if !seenLevel {
		return Fragment{}, fmt.Errorf("%s: missing level header (expected # [major], # [minor], # [patch], or # [unreleased])", path)
	}
	if !fragment.HasEntries() {
		return Fragment{}, fmt.Errorf("%s: fragment has no changelog entries", path)
	}
	return fragment, nil
}

// RenderFragment renders a stable, readable fragment. Empty sections are omitted.
func RenderFragment(fragment Fragment) string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "# [%s]\n", fragment.Level)
	for _, category := range Categories() {
		entries := fragment.Entries[category]
		if len(entries) == 0 {
			continue
		}
		fmt.Fprintf(&builder, "\n## %s\n\n", category)
		for _, entry := range entries {
			fmt.Fprintf(&builder, "- %s\n", strings.TrimSpace(entry))
		}
	}
	return builder.String()
}
