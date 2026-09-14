## ADDED Requirements

### Requirement: Authorized stable release handoff
The repository SHALL document the two GitHub Actions workflows that consume
local CLM release behavior: a main-branch workflow runs the quality gate,
invokes `clm cut`, and pushes the resulting release commit and CLM-created
tags; a stable-tag workflow builds and publishes artifacts. The CLM executable
MUST remain forge-neutral and MUST NOT make network calls or push Git state.

#### Scenario: Document CI ownership boundaries
- **WHEN** a maintainer follows the repository's release documentation
- **THEN** it identifies CLM as the creator of release commits and tags, the main-branch workflow as the authorized Git publisher, and the stable-tag workflow as the artifact publisher

### Requirement: Pull-request prerelease publication boundary
The repository SHALL document that the pull-request workflow uses CLM only to
validate fragments and derive a prerelease version. It SHALL identify CI as the
creator and publisher of temporary beta tags and SHALL state that fork pull
requests are validated without publication credentials.

#### Scenario: Document beta ownership
- **WHEN** a contributor follows the pull-request pipeline documentation
- **THEN** it distinguishes CLM-derived prerelease versions from CI-created beta tags and GitHub prereleases

### Requirement: Stable pipeline trigger boundary
The repository SHALL document that the stable-tag workflow accepts only
CLM-created stable `vX.Y.Z` tags. It MUST state that a beta tag is published by
the pull-request workflow itself and cannot trigger the stable artifact
workflow.

#### Scenario: Document tag-trigger isolation
- **WHEN** a maintainer follows the release pipeline documentation
- **THEN** it distinguishes the PR beta publication path from the stable tag-triggered release path
