package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

var workingDirectoryLock sync.Mutex

func TestCutCreatesCommitAndTags(t *testing.T) {
	workingDirectoryLock.Lock()
	defer workingDirectoryLock.Unlock()
	directory := newGitRepository(t)
	oldDirectory, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldDirectory)
	if err := os.Chdir(directory); err != nil {
		t.Fatal(err)
	}

	runCLM(t, "init")
	gitRun(t, directory, "add", ".")
	gitRun(t, directory, "commit", "-m", "initialize clm")
	runCLM(t, "new-patch", "fix")
	runCLM(t, "fixed", "-m", "Correct release rendering")
	gitRun(t, directory, "add", "changelog.d")
	gitRun(t, directory, "commit", "-m", "add fragment")

	runCLM(t, "cut", "--check")
	runCLM(t, "cut")

	if got := gitRun(t, directory, "tag", "--list", "v0.0.1"); got != "v0.0.1" {
		t.Fatalf("primary tag: %q", got)
	}
	if got := gitRun(t, directory, "tag", "--list", "clm/v0.0.1"); got != "clm/v0.0.1" {
		t.Fatalf("CLM tag: %q", got)
	}
	if _, err := os.Stat(filepath.Join(directory, "changelog.d", "archive", "v0.0.1")); err != nil {
		t.Fatal(err)
	}
	contents, err := os.ReadFile(filepath.Join(directory, "CHANGELOG.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(contents), "## [0.0.1]") || !strings.Contains(string(contents), "Correct release rendering") {
		t.Fatalf("unexpected changelog: %s", contents)
	}
}

func TestCheckSupportsAnUnbornRepository(t *testing.T) {
	workingDirectoryLock.Lock()
	defer workingDirectoryLock.Unlock()
	directory := newGitRepository(t)
	oldDirectory, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldDirectory)
	if err := os.Chdir(directory); err != nil {
		t.Fatal(err)
	}
	runCLM(t, "init")
	runCLM(t, "new-patch", "initial-release")
	runCLM(t, "added", "-m", "Initial release")
	runCLM(t, "cut", "--check")
}

func TestUnreleasedCutDoesNotTagAndVersionPreviewIsEmpty(t *testing.T) {
	workingDirectoryLock.Lock()
	defer workingDirectoryLock.Unlock()
	directory := newGitRepository(t)
	oldDirectory, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldDirectory)
	if err := os.Chdir(directory); err != nil {
		t.Fatal(err)
	}
	runCLM(t, "init")
	gitRun(t, directory, "add", ".")
	gitRun(t, directory, "commit", "-m", "initialize")
	runCLM(t, "new-unreleased", "docs")
	runCLM(t, "changed", "-m", "Clarify documentation")
	gitRun(t, directory, "add", "changelog.d")
	gitRun(t, directory, "commit", "-m", "unreleased fragment")
	output := runCLM(t, "next-version", "--json")
	var plan struct {
		Version *string `json:"version"`
	}
	if err := json.Unmarshal([]byte(output), &plan); err != nil {
		t.Fatal(err)
	}
	if plan.Version != nil {
		t.Fatalf("expected no version, got %s", *plan.Version)
	}
	runCLM(t, "cut")
	if got := gitRun(t, directory, "tag", "--list"); got != "" {
		t.Fatalf("unexpected tags: %q", got)
	}
	contents, err := os.ReadFile(filepath.Join(directory, "CHANGELOG.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(contents), "## [Unreleased]") || !strings.Contains(string(contents), "Clarify documentation") {
		t.Fatalf("unexpected changelog: %s", contents)
	}
}

func TestCheckRequiresChangedFragment(t *testing.T) {
	workingDirectoryLock.Lock()
	defer workingDirectoryLock.Unlock()
	directory := newGitRepository(t)
	oldDirectory, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldDirectory)
	if err := os.Chdir(directory); err != nil {
		t.Fatal(err)
	}
	runCLM(t, "init")
	gitRun(t, directory, "add", ".")
	gitRun(t, directory, "commit", "-m", "initialize")
	base := gitRun(t, directory, "rev-parse", "HEAD")
	runCLM(t, "new-patch", "fix")
	runCLM(t, "fixed", "-m", "Correct issue")
	gitRun(t, directory, "add", "changelog.d")
	gitRun(t, directory, "commit", "-m", "fragment")
	runCLM(t, "cut", "--check", "--require-fragment", "--base", base)
	if output := runCLMError(t, "cut", "--check", "--require-fragment", "--base", "HEAD"); !strings.Contains(output, "does not add or modify") {
		t.Fatalf("unexpected required-fragment error: %s", output)
	}
}

func TestCutPreflightLeavesFilesUnchangedForTagCollision(t *testing.T) {
	workingDirectoryLock.Lock()
	defer workingDirectoryLock.Unlock()
	directory := newGitRepository(t)
	oldDirectory, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldDirectory)
	if err := os.Chdir(directory); err != nil {
		t.Fatal(err)
	}
	runCLM(t, "init")
	gitRun(t, directory, "add", ".")
	gitRun(t, directory, "commit", "-m", "initialize")
	gitRun(t, directory, "tag", "v0.0.1")
	gitRun(t, directory, "tag", "clm/v0.0.2")
	runCLM(t, "new-patch", "fix")
	runCLM(t, "fixed", "-m", "Correct issue")
	gitRun(t, directory, "add", "changelog.d")
	gitRun(t, directory, "commit", "-m", "fragment")
	before, err := os.ReadFile(filepath.Join(directory, "CHANGELOG.md"))
	if err != nil {
		t.Fatal(err)
	}
	output := runCLMError(t, "cut")
	if !strings.Contains(output, "release tag already exists") {
		t.Fatalf("unexpected error: %s", output)
	}
	after, err := os.ReadFile(filepath.Join(directory, "CHANGELOG.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(before) != string(after) {
		t.Fatal("preflight failure changed CHANGELOG.md")
	}
}

func newGitRepository(t *testing.T) string {
	t.Helper()
	directory := t.TempDir()
	gitRun(t, directory, "init")
	gitRun(t, directory, "config", "user.email", "clm@example.test")
	gitRun(t, directory, "config", "user.name", "CLM Test")
	return directory
}

func gitRun(t *testing.T, directory string, args ...string) string {
	t.Helper()
	command := exec.Command("git", args...)
	command.Dir = directory
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, output)
	}
	return strings.TrimSpace(string(output))
}

func runCLM(t *testing.T, arguments ...string) string {
	t.Helper()
	var output bytes.Buffer
	command := NewRootCommand()
	command.SetArgs(arguments)
	command.SetOut(&output)
	command.SetErr(&output)
	if err := command.Execute(); err != nil {
		t.Fatalf("clm %s: %v\n%s", strings.Join(arguments, " "), err, output.String())
	}
	return output.String()
}

func runCLMError(t *testing.T, arguments ...string) string {
	t.Helper()
	var output bytes.Buffer
	command := NewRootCommand()
	command.SetArgs(arguments)
	command.SetOut(&output)
	command.SetErr(&output)
	if err := command.Execute(); err == nil {
		t.Fatalf("clm %s unexpectedly succeeded", strings.Join(arguments, " "))
	}
	return output.String()
}
