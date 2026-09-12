package cli

import (
	"bytes"
	"os"
	"testing"
)

func TestPrinterColorModes(t *testing.T) {
	t.Parallel()
	var output bytes.Buffer
	newPrinter(&output, "never").success("ready")
	if output.String() != "ready\n" {
		t.Fatalf("unexpected uncolored output %q", output.String())
	}
	output.Reset()
	newPrinter(&output, "always").success("ready")
	if output.String() != "\x1b[32mready\x1b[0m\n" {
		t.Fatalf("unexpected colored output %q", output.String())
	}
}

func TestNoColorWinsInAutoMode(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	printer := newPrinter(os.Stdout, "auto")
	if printer.color {
		t.Fatal("NO_COLOR must disable automatic color")
	}
}

func TestColorModeValidation(t *testing.T) {
	t.Parallel()
	if err := validateColorMode("always"); err != nil {
		t.Fatal(err)
	}
	if err := validateColorMode("rainbow"); err == nil {
		t.Fatal("expected invalid color mode")
	}
}

func TestCompletionScriptsAreNotEmpty(t *testing.T) {
	t.Parallel()
	for _, shell := range []string{"bash", "zsh", "fish", "powershell"} {
		var output bytes.Buffer
		command := NewRootCommand()
		command.SetArgs([]string{"completion", shell})
		command.SetOut(&output)
		command.SetErr(&output)
		if err := command.Execute(); err != nil {
			t.Fatalf("%s completion: %v", shell, err)
		}
		if output.Len() == 0 {
			t.Fatalf("%s completion was empty", shell)
		}
	}
}
