// Command release-notes prints a stable changelog entry body for GitHub Releases.
package main

import (
	"fmt"
	"os"

	"github.com/zwing99/change-log-magician/internal/release"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: release-notes vX.Y.Z")
		os.Exit(2)
	}
	changelog, err := os.ReadFile("CHANGELOG.md")
	if err == nil {
		var notes string
		notes, err = release.ExtractReleaseNotes(string(changelog), os.Args[1])
		if err == nil {
			fmt.Print(notes)
			return
		}
	}
	fmt.Fprintln(os.Stderr, "error:", err)
	os.Exit(1)
}
