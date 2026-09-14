package release

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExtractReleaseNotes(t *testing.T) {
	t.Parallel()
	changelog, err := os.ReadFile(filepath.Join("testdata", "release-notes.md"))
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name    string
		tag     string
		want    string
		wantErr bool
	}{
		{name: "first stable release", tag: "v1.2.3", want: "### Added\n\n- Keep this **Markdown** intact.\n\n### Fixed\n\n- Correct the important thing.\n"},
		{name: "later stable release", tag: "v1.2.2", want: "### Changed\n\n- Prior release body.\n"},
		{name: "missing release", tag: "v9.9.9", wantErr: true},
		{name: "prerelease tag", tag: "v1.2.3-beta.1", wantErr: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := ExtractReleaseNotes(string(changelog), test.tag)
			if test.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if got != test.want {
				t.Fatalf("ExtractReleaseNotes() = %q, want %q", got, test.want)
			}
		})
	}
}
