## Why

Teams need release notes to be easy to create during normal development and reliable to publish after merge. Hand-editing a shared changelog creates merge conflicts and makes semantic versioning, release commits, tags, and preview artifacts inconsistent across GitHub and GitLab repositories.

## What Changes

- Introduce Change Log Magician (`clm`), an idiomatic Go CLI for Keep a Changelog-style, fragment-based release notes.
- Add commands that create level-specific (`major`, `minor`, `patch`, or `unreleased`) Markdown fragments and append entries to the six Keep a Changelog categories through prompts, `$EDITOR`, or `-m` text.
- Maintain a root `CHANGELOG.md` with a permanent `Unreleased` section; process fragments into either that section or a dated SemVer release, promoting existing unreleased content when a versioned release is cut.
- Add deterministic release planning, validation, archive, Git commit, and annotated-tag behavior. A stable cut creates both `vX.Y.Z` and `clm/vX.Y.Z` tags.
- Add non-mutating MR/PR checks and version-preview/prerelease support for GitLab and GitHub workflows. Unreleased-only changes pass validation but do not create prereleases.
- Establish a well-tested Go project with Mise-managed tooling, shell completions, accessible colored terminal output, friendly banner output, and contributor/agent documentation under the MIT license.

## Capabilities

### New Capabilities

- `changelog-fragments`: Create, edit, parse, validate, and append to human-editable changelog fragments.
- `release-cutting`: Render and archive pending fragments into the root changelog, calculate stable versions, commit releases, and create tags.
- `pipeline-release-planning`: Validate merge-request/pull-request fragments and compute stable or prerelease version plans without mutating the repository.
- `cli-developer-experience`: Provide installation-oriented CLI behavior, completion output, terminal presentation, and project bootstrap conventions.

### Modified Capabilities

- None.

## Impact

- Adds a new Go module, `cmd/clm` executable, internal domain packages, tests, fixtures, and release/build configuration.
- Adds `CHANGELOG.md`, `changelog.d/`, `clm.toml`, `mise.toml`, `README.md`, `AGENTS.md`, `CONTRIBUTING.md`, and an MIT `LICENSE` to initialized repositories or the initial CLM repository as applicable.
- Requires Git for repository inspection, release commits, archive moves, and annotated tags; GitHub/GitLab pipeline examples require repository credentials only when publishing tags or artifacts.
