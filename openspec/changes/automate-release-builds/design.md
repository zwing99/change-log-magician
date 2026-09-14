## Context

The current tag-triggered workflow contains its own cross-compilation loop and
uploads raw binaries. CLM already has the necessary stable-cut boundary: it
creates the release commit, both annotated tags, and the final changelog entry,
but it deliberately does not push or call forge APIs. See proposal.md and the
release-artifact-publishing spec for the intended behavior.

## Goals / Non-Goals

**Goals:**

- Make the repository's Mise configuration the single source of build commands
  used both locally and by GitHub Actions.
- Publish assets only after CLM has created the version and tags.
- Keep the GitHub Release description mechanically derived from the cut
  changelog entry.
- Give versioned repository pull requests a downloadable beta artifact before
  merge.
- Exercise the flow by releasing this change through CLM.

**Non-Goals:**

- Add GoReleaser, package-manager publishers, installers, signing, SBOMs, or
  macOS notarization.
- Add network or push behavior to the CLM executable.
- Support targets beyond the existing five desktop/server targets.
- Publish prereleases for fork pull requests or create stable tags outside CLM.

## Decisions

### Main-branch cutting and stable-tag publishing are separate workflows

A qualifying push to `main` starts the cut workflow. It checks out full
history, runs `mise run check`, then runs `clm cut`. For a versioned cut, CLM
leaves the release commit and both tags in the checkout; the workflow pushes
that exact Git state. Its path filters exclude CLM's archived fragments so the
release commit does not start another cut.

The push credential must be a GitHub App installation token or scoped PAT,
because pushes made with the default `GITHUB_TOKEN` do not start another
workflow. The resulting `vX.Y.Z` tag starts the separate stable-tag workflow;
the `clm/vX.Y.Z` tag and every prerelease tag do not match that workflow's
effective job condition. This keeps CLM as tag authority while allowing the
release build to run exclusively from an immutable stable tag.

The stable-tag workflow also accepts a manually selected existing stable tag
for recovery. It checks out that tag, performs no cut or Git push, and
idempotently creates or updates the corresponding GitHub Release and assets.

### Mise owns artifact assembly

`mise.toml` gains a release-build task that receives a stable or prerelease tag
through an explicit task environment value. It validates that value, uses
`CGO_ENABLED=0` and Go cross-compilation for the existing target matrix, applies
the version linker value, and writes a deterministic `dist/` layout with
checksums.

An inline workflow shell loop was rejected because it cannot be run and checked
locally through the same project interface. A Justfile was rejected because
Mise already pins the toolchain and supplies every existing quality task.

### The release body is extracted from the post-cut changelog

The stable-tag workflow extracts the Markdown body between that
version's `##` heading and the next release heading (or end of file), retaining
category headings and entry text. The tag supplies the GitHub Release title,
so the extracted body omits the duplicate version/date heading. The publisher
uses this generated file as its release body rather than deriving notes from
Git commits.

### Pull requests receive CI-owned beta tags only after validation

The pull-request workflow retains `clm cut --check --require-fragment` as its
gate. When `clm next-version --json --prerelease-id
beta.pr.<number>.<sha>` reports a planned stable version, CI creates the
matching unique beta tag, pushes it, invokes the same Mise build task, and
publishes a GitHub prerelease within the pull-request workflow. These temporary
tags are CI-owned and deliberately do not start the stable tag workflow; CLM's
stable `vX.Y.Z` and `clm/vX.Y.Z` tags remain CLM-owned.

The beta workflow is restricted to same-repository pull requests. Fork pull
requests receive validation only because exposing repository write credentials
to untrusted fork code would be unsafe. Unreleased-only changes do not have a
stable version plan, so they receive no beta release, consistent with CLM's
existing prerelease behavior.

Beta release descriptions identify the candidate stable version and source pull
request. The exact changelog-entry body remains a stable-release guarantee,
because it is available only after `clm cut` renders the final entry.

### Dogfood through a real patch release

Implementation adds a patch fragment using CLM. After the change is merged,
the authorized workflow cuts and publishes that fragment as the first release
using the new path; this verifies the build task, Git handoff, release assets,
and changelog-derived description together.

## Risks / Trade-offs

- [The release workflow can push protected `main`] -> Restrict its triggering
  paths and write credential to authorized maintainers; grant `contents: write`
  only to the release job.
- [A failed artifact publication occurs after CLM's Git state is pushed] ->
  make publication idempotent for the existing stable tag and provide a
  workflow-dispatch retry path that rebuilds and republishes without another
  cut.
- [Shell parsing extracts the wrong release entry] -> Use the exact stable tag
  and a fixture-backed extraction check covering multiple changelog sections.
- [Cross-platform output regresses] -> Validate the artifact names, executable
  extensions, embedded version, and checksum coverage in task-level tests or
  CI verification.
- [Beta releases accumulate after each pull request update] -> Use the PR
  number and commit SHA for unique immutable previews and document that they
  are temporary prereleases rather than stable releases.

## Migration Plan

1. Add the Mise task and its local verification.
2. Replace the inline tag-triggered workflow with a main-branch cut workflow
   and a stable-tag artifact workflow, then document the credential and retry
   procedures.
3. Add the CLM patch fragment for this change.
4. Open a repository pull request to verify its beta prerelease, then merge it
   so the main workflow cuts it and the stable tag workflow publishes the
   dogfooded release.

Rollback removes the new workflow and task; already pushed CLM release commits
and tags remain auditable and a failed GitHub Release can be retried from the
same stable tag.
