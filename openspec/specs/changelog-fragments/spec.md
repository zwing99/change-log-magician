# changelog-fragments Specification

## Purpose

Provide a human-editable changelog fragment workflow that makes every merge request or pull request contribute structured release-note content.

## Requirements

### Requirement: Repository initialization
The system SHALL initialize a repository with a root `CHANGELOG.md` containing a Keep a Changelog `Unreleased` section, a `changelog.d/` fragment directory, an archive directory, and default CLM configuration when they do not already exist. The system MUST NOT overwrite existing changelog content or configuration without explicit user confirmation.

#### Scenario: Initialize an empty repository
- **WHEN** a user runs `clm init` in a Git repository without CLM files
- **THEN** the repository contains the initialized changelog structure and configuration

#### Scenario: Protect existing changelog
- **WHEN** a user runs `clm init` where `CHANGELOG.md` already exists
- **THEN** CLM reports the existing file and does not replace it without confirmation

### Requirement: Level-specific fragment creation
The system SHALL create a collision-safe Markdown fragment through `clm new-major`, `clm new-minor`, `clm new-patch`, or `clm new-unreleased`. Each created fragment MUST begin with exactly one corresponding level header and MUST be placed outside the archive directory.

#### Scenario: Create a minor fragment
- **WHEN** a user runs `clm new-minor shell-completions`
- **THEN** CLM creates a new fragment headed `[minor]` with the supplied name represented in its filename

#### Scenario: Create an unreleased fragment
- **WHEN** a user runs `clm new-unreleased documentation`
- **THEN** CLM creates a fragment headed `[unreleased]` that is eligible for validation but does not independently request a stable release

### Requirement: Categorized entry authoring
The system SHALL support `added`, `changed`, `deprecated`, `removed`, `fixed`, and `security` entry commands. Each command MUST append an entry to its matching Keep a Changelog section in a selected fragment. The `-m` option SHALL accept entry text; when it is absent, CLM SHALL interactively collect the entry and offer the user's configured editor for multiline content.

#### Scenario: Append supplied text
- **WHEN** a user runs `clm added -m "Support shell completions"` against a selected fragment
- **THEN** the fragment contains the text under `## Added`

#### Scenario: Prompt for omitted text
- **WHEN** a user runs `clm fixed` without `-m`
- **THEN** CLM prompts for the entry instead of silently adding an empty item

#### Scenario: Resolve an unambiguous target
- **WHEN** exactly one active fragment exists and the user omits a target
- **THEN** CLM appends to that fragment

#### Scenario: Reject an ambiguous target
- **WHEN** multiple active fragments exist and the user omits a target
- **THEN** CLM reports that a fragment selection is required and leaves all fragments unchanged

### Requirement: Manual-edit-compatible validation
The system SHALL accept manually edited fragments that retain one valid level header and valid Keep a Changelog category headings. It MUST reject fragments with missing or conflicting level headers, unsupported categories, malformed category structure, or no changelog entries, with actionable file-specific diagnostics.

#### Scenario: Validate a manually reordered fragment
- **WHEN** a valid fragment's category sections are manually reordered
- **THEN** CLM accepts and processes the fragment without changing the authored entry text

#### Scenario: Report an invalid level header
- **WHEN** a fragment has no valid level header
- **THEN** CLM reports the fragment path and the allowed headers
