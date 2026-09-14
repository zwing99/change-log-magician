package ci

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReleaseWorkflowsKeepReleaseBoundaries(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		file     string
		contains []string
		omits    []string
	}{
		{
			name: "main cut only consumes active fragments and pushes CLM state",
			file: "cut-release.yml",
			contains: []string{
				"branches: [main]", "changelog.d/*.md", "fetch-depth: 0", "mise run check", "go run ./cmd/clm cut", "git push origin HEAD:main --follow-tags", "RELEASE_PUSH_TOKEN",
			},
			omits: []string{"git tag", "clm next-version"},
		},
		{
			name: "stable publisher only builds from stable tag and changelog body",
			file: "release.yml",
			contains: []string{
				"v[0-9]+.[0-9]+.[0-9]+", "workflow_dispatch", "mise run release-build -- \"$TAG\"", "go run ./cmd/release-notes \"$TAG\" > release-notes.md", "body_path: release-notes.md",
			},
			omits: []string{"./cmd/clm cut", "git push", "git tag"},
		},
		{
			name: "pull requests validate before same repository beta publication",
			file: "ci.yml",
			contains: []string{
				"clm cut --check --require-fragment", "beta.pr.${{ github.event.number }}.${{ github.event.pull_request.head.sha }}", "head.repo.full_name == github.repository", "git tag -a \"$TAG\"", "prerelease: true",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			source, err := os.ReadFile(filepath.Join("..", "..", ".github", "workflows", test.file))
			if err != nil {
				t.Fatal(err)
			}
			for _, expected := range test.contains {
				if !strings.Contains(string(source), expected) {
					t.Errorf("workflow is missing %q", expected)
				}
			}
			for _, forbidden := range test.omits {
				if strings.Contains(string(source), forbidden) {
					t.Errorf("workflow must not contain %q", forbidden)
				}
			}
		})
	}
}
