## MODIFIED Requirements

### Requirement: CLM-owned main-branch release cut
The authorized GitHub Actions workflow SHALL run on qualifying pushes to `main`
that change active changelog fragments, run `clm cut`, and push only the release
commit and its `vX.Y.Z` and `clm/vX.Y.Z` tags created by CLM. The workflow MUST
ignore CLM archive-only release commits, use an authorized credential that
permits the stable-tag workflow to run after the push, and MUST NOT calculate,
create, or rename release tags itself.

#### Scenario: Cut a versioned main-branch release
- **WHEN** a qualifying push to `main` has pending versioned fragments
- **THEN** the workflow pushes CLM's release commit and both CLM-created stable tags without building or publishing release artifacts itself

#### Scenario: Cut an unreleased-only main-branch change
- **WHEN** a qualifying push to `main` has only unreleased fragments
- **THEN** the workflow pushes CLM's release commit, creates no stable tag, and starts no artifact publication

#### Scenario: Ignore an archive-only release commit
- **WHEN** CLM pushes a release commit that changes only `CHANGELOG.md` and archived fragments
- **THEN** the main-branch cut workflow does not start another cut job

## ADDED Requirements

### Requirement: Ancillary workflow tag isolation
Repository setup and validation workflows SHALL not start from stable, companion,
or prerelease tag pushes unless they are artifact-publishing workflows that
explicitly consume those tags.

#### Scenario: Ignore release tags in setup automation
- **WHEN** CLM or CI pushes `vX.Y.Z`, `clm/vX.Y.Z`, or a beta prerelease tag
- **THEN** the Copilot setup workflow does not start while the stable artifact publisher retains its stable-tag trigger
