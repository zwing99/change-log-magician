// Package cli exposes the CLM command line.
package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/zwing99/change-log-magician/internal/changelog"
	"github.com/zwing99/change-log-magician/internal/filesystem"
	"github.com/zwing99/change-log-magician/internal/gitx"
	"github.com/zwing99/change-log-magician/internal/release"
)

const banner = "   ____ _     __  __\n  / ___| |   |  \\/  |\n | |   | |   | |\\/| |\n | |___| |___| |  | |\n  \\____|_____|_|  |_|\n"

// NewRootCommand constructs the complete CLI and is testable without process globals.
func NewRootCommand() *cobra.Command {
	var color string
	root := &cobra.Command{
		Use:     "clm",
		Short:   "Change Log Magician",
		Long:    "Change Log Magician manages Keep a Changelog fragments and Git release cuts.",
		Version: "devel",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := validateColorMode(color); err != nil {
				return err
			}
			printer := newPrinter(cmd.OutOrStdout(), color)
			if printer.color || printer.interactive {
				fmt.Fprint(cmd.OutOrStdout(), banner)
			}
			return cmd.Help()
		},
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error { return validateColorMode(color) },
	}
	root.PersistentFlags().StringVar(&color, "color", "auto", "color mode: auto, always, or never")
	root.AddCommand(newInitCommand(&color))
	for _, level := range []changelog.Level{changelog.Major, changelog.Minor, changelog.Patch, changelog.Unreleased} {
		root.AddCommand(newNewCommand(level, &color))
	}
	for _, category := range changelog.Categories() {
		root.AddCommand(newEntryCommand(category, &color))
	}
	root.AddCommand(newCutCommand(&color), newNextVersionCommand(&color), newCompletionCommand())
	return root
}

func validateColorMode(mode string) error {
	switch mode {
	case "auto", "always", "never":
		return nil
	default:
		return fmt.Errorf("invalid --color value %q (expected auto, always, or never)", mode)
	}
}

func repository() (filesystem.Repository, gitx.Client, error) {
	workingDirectory, err := os.Getwd()
	if err != nil {
		return filesystem.Repository{}, gitx.Client{}, err
	}
	return filesystem.Repository{Root: workingDirectory}, gitx.Client{Dir: workingDirectory}, nil
}

func newInitCommand(color *string) *cobra.Command {
	return &cobra.Command{Use: "init", Short: "Initialize CLM files in the current Git repository", RunE: func(cmd *cobra.Command, _ []string) error {
		repo, git, err := repository()
		if err != nil {
			return err
		}
		if err := git.EnsureRepository(); err != nil {
			return err
		}
		if err := repo.Initialize(); err != nil {
			return err
		}
		newPrinter(cmd.OutOrStdout(), *color).success("Initialized CHANGELOG.md and changelog.d/")
		return nil
	}}
}

func newNewCommand(level changelog.Level, color *string) *cobra.Command {
	return &cobra.Command{Use: "new-" + string(level) + " [name]", Short: "Create a " + string(level) + " changelog fragment", Args: cobra.MaximumNArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		repo, _, err := repository()
		if err != nil {
			return err
		}
		name := "change"
		if len(args) == 1 {
			name = args[0]
		}
		path, err := repo.CreateFragment(level, name, time.Now())
		if err != nil {
			return err
		}
		newPrinter(cmd.OutOrStdout(), *color).success("Created " + path)
		return nil
	}}
}

func newEntryCommand(category changelog.Category, color *string) *cobra.Command {
	var message, target string
	var useEditor bool
	command := &cobra.Command{Use: strings.ToLower(string(category)), Short: "Add an entry to the " + string(category) + " section", RunE: func(cmd *cobra.Command, _ []string) error {
		repo, _, err := repository()
		if err != nil {
			return err
		}
		if message == "" {
			message, err = collectEntry(cmd, useEditor)
			if err != nil {
				return err
			}
		}
		if err := repo.AppendEntry(target, category, message); err != nil {
			return err
		}
		newPrinter(cmd.OutOrStdout(), *color).success("Added " + string(category) + " entry")
		return nil
	}}
	command.Flags().StringVarP(&message, "message", "m", "", "entry text")
	command.Flags().StringVarP(&target, "change", "c", "", "fragment path")
	command.Flags().BoolVar(&useEditor, "editor", false, "collect entry in $EDITOR")
	return command
}

func collectEntry(cmd *cobra.Command, useEditor bool) (string, error) {
	if useEditor {
		return editEntry()
	}
	fmt.Fprint(cmd.OutOrStdout(), "Entry text (use --editor for multiline editing): ")
	entry, err := bufio.NewReader(cmd.InOrStdin()).ReadString('\n')
	if err != nil && err != io.EOF {
		return "", err
	}
	entry = strings.TrimSpace(entry)
	if entry == "" {
		return "", fmt.Errorf("changelog entry cannot be empty")
	}
	return entry, nil
}

func editEntry() (string, error) {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		return "", fmt.Errorf("$EDITOR is not set; use -m or configure an editor")
	}
	file, err := os.CreateTemp("", "clm-entry-*.md")
	if err != nil {
		return "", err
	}
	path := file.Name()
	defer os.Remove(path)
	if err := file.Close(); err != nil {
		return "", err
	}
	command := exec.Command("sh", "-c", editor+" \""+path+"\"")
	command.Stdin, command.Stdout, command.Stderr = os.Stdin, os.Stdout, os.Stderr
	if err := command.Run(); err != nil {
		return "", fmt.Errorf("run editor: %w", err)
	}
	contents, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	entry := strings.TrimSpace(string(contents))
	if entry == "" {
		return "", fmt.Errorf("changelog entry cannot be empty")
	}
	return entry, nil
}

func newCutCommand(color *string) *cobra.Command {
	var check, requireFragment bool
	var base string
	command := &cobra.Command{Use: "cut", Short: "Validate or cut pending changelog fragments", RunE: func(cmd *cobra.Command, _ []string) error {
		repo, git, err := repository()
		if err != nil {
			return err
		}
		if err := git.EnsureRepository(); err != nil {
			return err
		}
		plan, root, err := buildPlan(repo, git)
		if err != nil {
			return err
		}
		if requireFragment {
			if base == "" {
				return fmt.Errorf("--require-fragment requires --base <revision>")
			}
			if err := requireChangedFragment(git, base); err != nil {
				return err
			}
		}
		printer := newPrinter(cmd.OutOrStdout(), *color)
		if check {
			printer.success(planMessage(plan))
			return nil
		}
		if err := git.EnsureClean(); err != nil {
			return err
		}
		if err := executeCut(repo, git, plan, root); err != nil {
			return err
		}
		printer.success("Cut " + cutName(plan))
		return nil
	}}
	command.Flags().BoolVar(&check, "check", false, "validate without writing files or Git state")
	command.Flags().BoolVar(&requireFragment, "require-fragment", false, "require this branch to change an active fragment")
	command.Flags().StringVar(&base, "base", "", "base revision for --require-fragment")
	return command
}

func newNextVersionCommand(color *string) *cobra.Command {
	var prereleaseID string
	var jsonOutput bool
	command := &cobra.Command{Use: "next-version", Short: "Print the next cut version without changing the repository", RunE: func(cmd *cobra.Command, _ []string) error {
		repo, git, err := repository()
		if err != nil {
			return err
		}
		if err := git.EnsureRepository(); err != nil {
			return err
		}
		plan, _, err := buildPlan(repo, git)
		if err != nil {
			return err
		}
		if prereleaseID != "" && plan.Version == nil {
			return fmt.Errorf("no prerelease is available for unreleased-only fragments")
		}
		if jsonOutput {
			payload := struct {
				Cuttable           bool    `json:"cuttable"`
				Level              string  `json:"level,omitempty"`
				Version            *string `json:"version"`
				Prerelease         string  `json:"prerelease,omitempty"`
				Fragments          int     `json:"fragments"`
				IncludesUnreleased bool    `json:"includesUnreleased"`
			}{Cuttable: true, Fragments: len(plan.Fragments)}
			for _, fragment := range plan.Fragments {
				if fragment.Level == changelog.Unreleased {
					payload.IncludesUnreleased = true
				}
			}
			if plan.Version != nil {
				version := plan.Version.String()
				payload.Version = &version
				payload.Level = string(plan.Level)
				if prereleaseID != "" {
					payload.Prerelease, err = release.Prerelease(*plan.Version, prereleaseID)
					if err != nil {
						return err
					}
				}
			}
			return json.NewEncoder(cmd.OutOrStdout()).Encode(payload)
		}
		if plan.Version == nil {
			fmt.Fprintln(cmd.OutOrStdout(), "No stable version planned (unreleased-only fragments). ")
			return nil
		}
		value := plan.Version.String()
		if prereleaseID != "" {
			value, err = release.Prerelease(*plan.Version, prereleaseID)
			if err != nil {
				return err
			}
		}
		fmt.Fprintln(cmd.OutOrStdout(), value)
		return nil
	}}
	command.Flags().BoolVar(&jsonOutput, "json", false, "emit a JSON release plan")
	command.Flags().StringVar(&prereleaseID, "prerelease-id", "", "derive a prerelease identifier")
	return command
}

func newCompletionCommand() *cobra.Command {
	command := &cobra.Command{Use: "completion [bash|zsh|fish|powershell]", Short: "Generate shell completion scripts", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		root := cmd.Root()
		switch args[0] {
		case "bash":
			return root.GenBashCompletion(cmd.OutOrStdout())
		case "zsh":
			return root.GenZshCompletion(cmd.OutOrStdout())
		case "fish":
			return root.GenFishCompletion(cmd.OutOrStdout(), true)
		case "powershell":
			return root.GenPowerShellCompletionWithDesc(cmd.OutOrStdout())
		default:
			return fmt.Errorf("unsupported shell %q", args[0])
		}
	}}
	return command
}

func buildPlan(repo filesystem.Repository, git gitx.Client) (release.Plan, changelog.Root, error) {
	fragments, err := repo.ActiveFragments()
	if err != nil {
		return release.Plan{}, changelog.Root{}, err
	}
	contents, err := os.ReadFile(filepath.Join(repo.Root, filesystem.ChangelogFile))
	if err != nil {
		return release.Plan{}, changelog.Root{}, fmt.Errorf("read CHANGELOG.md: %w", err)
	}
	root, err := changelog.ParseRoot(string(contents))
	if err != nil {
		return release.Plan{}, changelog.Root{}, err
	}
	base, err := git.LatestStableVersion()
	if err != nil {
		return release.Plan{}, changelog.Root{}, err
	}
	plan, err := release.BuildPlan(root, fragments, base, time.Now())
	return plan, root, err
}

func requireChangedFragment(git gitx.Client, base string) error {
	files, err := git.ChangedFiles(base)
	if err != nil {
		return err
	}
	for _, file := range files {
		if strings.HasPrefix(filepath.ToSlash(file), filesystem.FragmentsDir+"/") && strings.HasSuffix(file, ".md") && !strings.HasPrefix(filepath.ToSlash(file), filesystem.ArchiveDir+"/") {
			return nil
		}
	}
	return fmt.Errorf("this branch does not add or modify an active changelog fragment")
}

func executeCut(repo filesystem.Repository, git gitx.Client, plan release.Plan, _ changelog.Root) error {
	if plan.Version != nil && (git.TagExists("v"+plan.Version.String()) || git.TagExists("clm/v"+plan.Version.String())) {
		return fmt.Errorf("release tag already exists; aborting before mutation")
	}
	archive := filepath.Join(repo.Root, filesystem.ArchiveDir, plan.ArchiveSegment)
	if _, err := os.Stat(archive); err == nil {
		return fmt.Errorf("archive destination %s already exists", archive)
	} else if !os.IsNotExist(err) {
		return err
	}
	if err := os.MkdirAll(archive, 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(repo.Root, filesystem.ChangelogFile), []byte(plan.Rendered), 0o644); err != nil {
		return err
	}
	if err := git.Add(filesystem.ChangelogFile); err != nil {
		return err
	}
	for _, fragment := range plan.Fragments {
		source := filepath.FromSlash(fragment.Path)
		destination := filepath.Join(filesystem.ArchiveDir, plan.ArchiveSegment, filepath.Base(source))
		if err := git.Move(source, destination); err != nil {
			return err
		}
	}
	if err := git.Add(filepath.ToSlash(filepath.Join(filesystem.ArchiveDir, plan.ArchiveSegment))); err != nil {
		return err
	}
	message := "chore(release): " + cutName(plan)
	if err := git.Commit(message); err != nil {
		return err
	}
	if plan.Version != nil {
		version := plan.Version.String()
		if err := git.AnnotatedTag("v"+version, "Release v"+version); err != nil {
			return fmt.Errorf("release commit created but primary tag failed: %w", err)
		}
		if err := git.AnnotatedTag("clm/v"+version, "CLM release v"+version); err != nil {
			return fmt.Errorf("release commit and v%s tag created but CLM tag failed: %w", version, err)
		}
	}
	return nil
}

func cutName(plan release.Plan) string {
	if plan.Version == nil {
		return "unreleased changelog entries"
	}
	return "v" + plan.Version.String()
}
func planMessage(plan release.Plan) string {
	if plan.Version == nil {
		return fmt.Sprintf("Cuttable: unreleased-only (%d fragments, no tag)", len(plan.Fragments))
	}
	return fmt.Sprintf("Release ready: v%s (%d fragments, %s)", plan.Version.String(), len(plan.Fragments), plan.Level)
}

type printer struct {
	writer             io.Writer
	color, interactive bool
}

func newPrinter(writer io.Writer, mode string) printer {
	interactive := false
	if file, ok := writer.(*os.File); ok {
		if info, err := file.Stat(); err == nil {
			interactive = info.Mode()&os.ModeCharDevice != 0
		}
	}
	color := mode == "always" || (mode == "auto" && interactive && os.Getenv("NO_COLOR") == "")
	return printer{writer: writer, color: color, interactive: interactive}
}
func (p printer) success(message string) {
	if p.color {
		fmt.Fprintf(p.writer, "\x1b[32m%s\x1b[0m\n", message)
		return
	}
	fmt.Fprintln(p.writer, message)
}
