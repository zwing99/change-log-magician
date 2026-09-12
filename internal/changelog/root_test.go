package changelog

import (
	"strings"
	"testing"
)

func TestRootRenderPreservesPreambleAndTrailer(t *testing.T) {
	t.Parallel()
	source := "# Changelog\n\nIntro text.\n\n## [Unreleased]\n\n### Added\n\n- Existing\n\n## Links\n\n[Unreleased]: example\n"
	root, err := ParseRoot(source)
	if err != nil {
		t.Fatal(err)
	}
	root.Unreleased().Entries[Fixed] = []string{"Correction"}
	rendered := RenderRoot(root)
	for _, expected := range []string{"Intro text.", "- Existing", "- Correction", "## Links", "[Unreleased]: example"} {
		if !strings.Contains(rendered, expected) {
			t.Fatalf("missing %q in %s", expected, rendered)
		}
	}
}

func TestNewRootIncludesCLMMaintenanceHints(t *testing.T) {
	t.Parallel()
	rendered := RenderRoot(NewRoot())
	for _, expected := range []string{"maintained with", "clm new-patch <name>", "clm added -m", "clm cut --check"} {
		if !strings.Contains(rendered, expected) {
			t.Fatalf("default changelog is missing %q", expected)
		}
	}
}
