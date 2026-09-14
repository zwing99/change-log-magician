package release

import (
	"fmt"
	"strings"
)

// ExtractReleaseNotes returns a stable release body without its version heading.
func ExtractReleaseNotes(changelog, tag string) (string, error) {
	version, err := ParseVersion(tag)
	if err != nil {
		return "", fmt.Errorf("release notes require a stable tag: %w", err)
	}
	heading := "## [" + version.String() + "] - "
	lines := strings.Split(strings.ReplaceAll(changelog, "\r\n", "\n"), "\n")
	start := -1
	for index, line := range lines {
		if strings.HasPrefix(line, heading) {
			start = index + 1
			break
		}
	}
	if start == -1 {
		return "", fmt.Errorf("release %q was not found in CHANGELOG.md", tag)
	}
	end := len(lines)
	for index := start; index < len(lines); index++ {
		if strings.HasPrefix(lines[index], "## ") {
			end = index
			break
		}
	}
	body := strings.Trim(strings.Join(lines[start:end], "\n"), "\n")
	if body == "" {
		return "", fmt.Errorf("release %q has no changelog body", tag)
	}
	return body + "\n", nil
}
