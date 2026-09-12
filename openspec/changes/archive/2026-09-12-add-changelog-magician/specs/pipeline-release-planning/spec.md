## Purpose

Enable GitHub and GitLab merge-request pipelines to enforce fragment policy and obtain deterministic release-version inputs without modifying repository state.

## ADDED Requirements

### Requirement: Merge-request fragment enforcement
The system SHALL provide a non-mutating cut check that validates the complete active fragment set and, when given a base revision with the required-fragment option, verifies that the proposed branch adds or modifies at least one valid active fragment.

#### Scenario: Accept a merge request with a valid fragment
- **WHEN** CI runs `clm cut --check --require-fragment --base <base-revision>` and the branch contributes a valid fragment
- **THEN** CLM exits successfully without editing files, creating commits, or creating tags

#### Scenario: Reject a merge request without a fragment
- **WHEN** CI runs the required-fragment check and the branch contributes no active fragment
- **THEN** CLM exits unsuccessfully with a message explaining the missing fragment requirement

#### Scenario: Reject a globally uncuttable pending set
- **WHEN** the branch adds a valid fragment but another active fragment is malformed
- **THEN** CLM fails the check and identifies the malformed fragment

### Requirement: Version preview
The system SHALL provide a non-mutating machine-readable version plan for the active fragment set. The plan MUST report the calculated stable version and included fragments when any versioned fragment is pending, and MUST explicitly report that no stable version is planned for an unreleased-only set.

#### Scenario: Preview a minor release
- **WHEN** CI requests a version plan for a pending minor release
- **THEN** CLM returns the calculated next stable version without changing repository state

#### Scenario: Preview an unreleased-only change
- **WHEN** CI requests a version plan and only unreleased fragments are pending
- **THEN** CLM reports no planned stable version and exits successfully

### Requirement: Prerelease version derivation
The system SHALL derive a valid SemVer prerelease version from a planned stable version and a caller-provided prerelease identifier. The CLI MUST NOT derive or publish a prerelease when the active set is unreleased-only.

#### Scenario: Derive a pull-request prerelease
- **WHEN** CI supplies a planned version `1.4.0` and identifier `pr.123.abcdef0`
- **THEN** CLM returns `1.4.0-pr.123.abcdef0` for artifact or tag publication

#### Scenario: Skip an unreleased-only prerelease
- **WHEN** CI requests a prerelease for an unreleased-only active set
- **THEN** CLM reports that no prerelease is available and does not create a tag

### Requirement: Forge-neutral operation
The system SHALL work against local Git state without requiring GitHub or GitLab APIs. It MUST document equivalent GitHub Actions and GitLab CI invocation patterns, including that credentials are required only by the workflow step that pushes commits, tags, or published artifacts.

#### Scenario: Run from a GitLab merge-request pipeline
- **WHEN** GitLab CI provides a checked-out repository and merge-request base revision
- **THEN** the validation and version-preview commands operate without GitLab API credentials

#### Scenario: Run from a GitHub pull-request workflow
- **WHEN** GitHub Actions provides a checked-out repository and pull-request base revision
- **THEN** the validation and version-preview commands operate without GitHub API credentials
