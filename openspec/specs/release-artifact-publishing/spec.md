# release-artifact-publishing Specification

## Purpose

Build stable and beta CLM release artifacts consistently and publish them with
the tags and changelog content that define each release.

## Requirements

### Requirement: Reproducible versioned artifact build
The repository SHALL expose a Mise task that accepts one explicit stable or
prerelease SemVer tag and produces version-stamped `clm` binaries for Linux
amd64 and arm64, macOS amd64 and arm64, and Windows amd64. The task MUST emit
SHA-256 checksums for every published binary and MUST fail for a malformed
release version.

#### Scenario: Build artifacts for a versioned release
- **WHEN** a maintainer runs the release-build task for `v1.2.3` or `v1.2.3-beta.pr.42.abcdef`
- **THEN** it produces the five version-stamped platform binaries and a checksum file without creating Git commits or tags

#### Scenario: Reject a malformed build version
- **WHEN** a maintainer runs the release-build task with a malformed version
- **THEN** the task fails before publishing any release artifact

### Requirement: CLM-owned main-branch release cut
The authorized GitHub Actions workflow SHALL run on qualifying pushes to `main`,
run `clm cut`, and push only the release commit and its `vX.Y.Z` and
`clm/vX.Y.Z` tags created by CLM. The workflow MUST use an authorized
credential that permits the stable-tag workflow to run after the push, and it
MUST NOT calculate, create, or rename release tags itself.

#### Scenario: Cut a versioned main-branch release
- **WHEN** a qualifying push to `main` has pending versioned fragments
- **THEN** the workflow pushes CLM's release commit and both CLM-created stable tags without building or publishing release artifacts itself

#### Scenario: Cut an unreleased-only main-branch change
- **WHEN** a qualifying push to `main` has only unreleased fragments
- **THEN** the workflow pushes CLM's release commit, creates no stable tag, and starts no artifact publication

### Requirement: Stable-tag artifact publication
The GitHub Actions artifact workflow SHALL run only for a stable `vX.Y.Z` tag
created by CLM, use that tag as the explicit input to the Mise release-build
task, and publish the generated assets to the matching GitHub Release. It MUST
ignore the companion `clm/vX.Y.Z` tag and every prerelease tag, including PR
beta tags, and MUST NOT invoke `clm cut` or mutate Git tags.

#### Scenario: Build from a CLM stable tag
- **WHEN** CI receives a pushed `v1.2.3` tag created by CLM
- **THEN** it runs the Mise build task for `v1.2.3` and publishes its artifacts to the `v1.2.3` GitHub Release

#### Scenario: Ignore a companion CLM tag
- **WHEN** CI receives a pushed `clm/v1.2.3` tag
- **THEN** it does not start the artifact workflow

#### Scenario: Ignore a beta prerelease tag
- **WHEN** CI receives a pushed `v1.2.3-beta.pr.42.abcdef` tag
- **THEN** it does not start the stable artifact workflow

### Requirement: Existing stable tag publication retry
The artifact workflow SHALL support an authorized manual retry for an existing
stable `vX.Y.Z` tag. A retry MUST rebuild the tag's artifacts and create or
update that tag's GitHub Release without invoking CLM or pushing commits or
tags.

#### Scenario: Republish an already-pushed tag
- **WHEN** a maintainer retries publication for an existing `v1.2.3` tag
- **THEN** the workflow updates the `v1.2.3` GitHub Release assets and description without creating another release cut or tag

### Requirement: Pull-request beta prerelease publication
The pull-request workflow SHALL validate changelog fragments and derive a
prerelease version with `clm next-version --prerelease-id` for a versioned
change. For a pull request from the repository, it SHALL create and push a
unique `vX.Y.Z-beta.pr.<number>.<sha>` prerelease tag, build assets through
Mise, and publish them as a downloadable GitHub prerelease from that
pull-request workflow. CI MUST NOT create a beta tag or prerelease for an
unreleased-only change or a pull request from a fork without repository write
authority. Pushing a beta tag MUST NOT start the stable artifact workflow.

#### Scenario: Publish a versioned pull-request beta
- **WHEN** repository pull request 42 has a pending patch release and commit SHA `abcdef`
- **THEN** CI validates it and its pull-request workflow publishes a downloadable prerelease tagged `v0.0.2-beta.pr.42.abcdef` without starting the stable artifact workflow

#### Scenario: Do not publish an unreleased-only pull request
- **WHEN** a pull request has only unreleased fragments
- **THEN** CI validates the pull request but creates no beta tag, artifacts, or GitHub prerelease

#### Scenario: Safely validate a fork pull request
- **WHEN** a pull request originates from a fork
- **THEN** CI validates its changelog fragments without using write credentials or publishing a beta prerelease

### Requirement: Changelog-derived GitHub Release description
The GitHub Release description SHALL use the Markdown body of the stable
version's entry in `CHANGELOG.md` after `clm cut`. The workflow MUST preserve
the entry's category headings and text without generating substitute release
notes.

#### Scenario: Publish release notes from the cut entry
- **WHEN** CLM cuts `v1.2.3` with Added and Fixed entries
- **THEN** the `v1.2.3` GitHub Release description contains the same Added and Fixed Markdown content as the `1.2.3` changelog entry
