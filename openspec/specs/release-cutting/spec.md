# release-cutting Specification

## Purpose

Turn validated pending fragments into an auditable Keep a Changelog release, Git commit, archive, and annotated tags.

## Requirements

### Requirement: Release-plan calculation
The system SHALL calculate the next stable version from the highest pending versioned fragment level using the latest reachable stable `vX.Y.Z` tag as its base. When no stable tag exists, it MUST use `0.0.0` as the base. `major` takes precedence over `minor`, which takes precedence over `patch`.

#### Scenario: Calculate the first minor release
- **WHEN** no stable release tag exists and the highest pending fragment level is `minor`
- **THEN** the calculated stable version is `0.1.0`

#### Scenario: Calculate a release from mixed fragments
- **WHEN** pending fragments include both `patch` and `major` levels
- **THEN** the calculated release uses a major version increment

### Requirement: Unreleased changelog lifecycle
The system SHALL keep an `Unreleased` section in the root changelog at all times. A cut containing only `unreleased` fragments MUST merge their entries into that section, archive their fragments, and create no release tag. A cut containing any versioned fragment MUST include both existing root `Unreleased` entries and all pending fragment entries in the dated numbered release, then recreate an empty `Unreleased` section.

#### Scenario: Cut unreleased-only fragments
- **WHEN** the pending directory contains only `unreleased` fragments
- **THEN** CLM updates `CHANGELOG.md`, archives the fragments, commits the change, and creates no tags

#### Scenario: Promote existing unreleased entries
- **WHEN** the root changelog contains unreleased entries and a pending minor fragment is cut
- **THEN** the numbered minor release contains those existing entries and the root changelog retains an empty `Unreleased` section

### Requirement: Stable cut execution
The system SHALL make `clm cut` fail before mutation when there are no active fragments or validation fails. For a valid cut, it MUST render the root changelog in Keep a Changelog order with a `YYYY-MM-DD` release date, use Git to move processed fragments into the committed archive, create one release commit, and create annotated tags `vX.Y.Z` and `clm/vX.Y.Z` for versioned cuts.

#### Scenario: Cut a patch release
- **WHEN** a validated patch fragment is pending
- **THEN** CLM writes the dated patch release, moves its fragment to the versioned archive, commits the release changes, and creates both annotated stable tags on that commit

#### Scenario: Reject an empty cut
- **WHEN** a user runs `clm cut` with no active fragments
- **THEN** CLM exits with an actionable error and creates no commit or tag

### Requirement: Safe release boundaries
The system SHALL preflight repository state, changelog rendering, archive collisions, and tag collisions before changing tracked release files. It MUST preserve changelog content outside the sections it owns and MUST report recovery guidance if a Git operation fails after a release commit is created.

#### Scenario: Detect a tag collision
- **WHEN** the calculated stable tag already exists
- **THEN** CLM aborts before changing the changelog or archive

#### Scenario: Preserve manually maintained changelog content
- **WHEN** the root changelog contains content outside CLM-managed release sections
- **THEN** a successful cut retains that content unchanged
