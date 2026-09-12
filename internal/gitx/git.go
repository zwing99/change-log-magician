// Package gitx is CLM's narrow Git command boundary.
package gitx

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/zwing99/change-log-magician/internal/release"
)

// Client runs Git against one repository directory.
type Client struct{ Dir string }

func (g Client) run(args ...string) (string, error) {
	command := exec.Command("git", args...)
	command.Dir = g.Dir
	var stderr bytes.Buffer
	command.Stderr = &stderr
	output, err := command.Output()
	if err != nil {
		return "", fmt.Errorf("git %s: %s", strings.Join(args, " "), strings.TrimSpace(stderr.String()))
	}
	return strings.TrimSpace(string(output)), nil
}

// EnsureRepository verifies that CLM is operating inside a Git working tree.
func (g Client) EnsureRepository() error {
	output, err := g.run("rev-parse", "--is-inside-work-tree")
	if err != nil || output != "true" {
		return fmt.Errorf("CLM must run inside a Git working tree")
	}
	return nil
}

// EnsureClean rejects unrelated or uncommitted changes before a release commit.
func (g Client) EnsureClean() error {
	output, err := g.run("status", "--porcelain")
	if err != nil {
		return err
	}
	if output != "" {
		return fmt.Errorf("working tree is not clean; commit or stash changes before cutting a release")
	}
	return nil
}

// LatestStableVersion returns the highest stable vX.Y.Z tag reachable from HEAD.
func (g Client) LatestStableVersion() (release.Version, error) {
	// An unborn repository has no reachable tags and starts from the implicit 0.0.0 baseline.
	if _, err := g.run("rev-parse", "--verify", "HEAD"); err != nil {
		return release.Version{}, nil
	}
	output, err := g.run("tag", "--merged", "HEAD", "--list", "v*")
	if err != nil {
		return release.Version{}, err
	}
	versions := make([]release.Version, 0)
	for _, tag := range strings.Fields(output) {
		version, err := release.ParseVersion(tag)
		if err == nil {
			versions = append(versions, version)
		}
	}
	if len(versions) == 0 {
		return release.Version{}, nil
	}
	sort.Slice(versions, func(i, j int) bool {
		if versions[i].Major != versions[j].Major {
			return versions[i].Major < versions[j].Major
		}
		if versions[i].Minor != versions[j].Minor {
			return versions[i].Minor < versions[j].Minor
		}
		return versions[i].Patch < versions[j].Patch
	})
	return versions[len(versions)-1], nil
}

// TagExists reports whether a tag is already present.
func (g Client) TagExists(tag string) bool {
	_, err := g.run("rev-parse", "--verify", "--quiet", "refs/tags/"+tag)
	return err == nil
}

// ChangedFiles returns files changed between a base revision and HEAD.
func (g Client) ChangedFiles(base string) ([]string, error) {
	output, err := g.run("diff", "--name-only", base+"...HEAD")
	if err != nil {
		return nil, err
	}
	if output == "" {
		return nil, nil
	}
	return strings.Fields(output), nil
}

// Move stages a Git-aware rename.
func (g Client) Move(source, destination string) error {
	if _, err := g.run("mv", "--", source, destination); err != nil {
		return err
	}
	return nil
}

// Add stages a path.
func (g Client) Add(path string) error { _, err := g.run("add", "--", path); return err }

// Commit creates one release commit.
func (g Client) Commit(message string) error { _, err := g.run("commit", "-m", message); return err }

// AnnotatedTag creates an annotated tag on HEAD.
func (g Client) AnnotatedTag(tag, message string) error {
	_, err := g.run("tag", "-a", tag, "-m", message)
	return err
}

// Relative returns an OS-independent repository-relative path.
func (g Client) Relative(path string) (string, error) {
	relative, err := filepath.Rel(g.Dir, path)
	if err != nil {
		return "", err
	}
	return filepath.ToSlash(relative), nil
}
