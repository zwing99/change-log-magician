## Why

CLM can cut an auditable release and create its local tags, but this repository has no complete, dogfooded path from that cut to downloadable GitHub release assets. The existing tag-triggered workflow duplicates build logic outside the project task runner and publishes assets without using CLM's release commit as the authoritative release description.

## What Changes

- Add a Mise release-build task that creates version-stamped, cross-platform CLM binaries and their checksums from an explicit stable or prerelease version.
- Add an authorized GitHub Actions main-branch workflow that runs CLM to cut releases and pushes only the commit and tags CLM created.
- Add a GitHub Actions stable-tag workflow that invokes the Mise build task and publishes or updates the matching GitHub Release.
- Extend pull-request validation to publish downloadable beta prereleases for versioned changes, using a CI-created prerelease tag derived from CLM's version plan.
- Restrict the stable release pipeline to CLM-created stable tags; beta tags remain inputs to no tag-triggered release pipeline.
- Publish the GitHub Release for the `vX.Y.Z` tag created by CLM and use that version's rendered changelog section as the release description, including when the stable tag already exists.
- Dogfood CLM by adding and cutting the changelog fragment for this change through the release workflow it introduces.
- Retire the standalone tag-triggered release workflow's inline build loop.

## Capabilities

### New Capabilities

- `release-artifact-publishing`: Build stable and beta CLM artifacts through Mise and publish GitHub releases that correspond to their version tags.

### Modified Capabilities

- `pipeline-release-planning`: Document the authorized GitHub release workflow as the consumer of CLM's local release plan and tags while retaining CLM's forge-neutral behavior.

## Impact

- Affects `mise.toml`, GitHub Actions main and tag release automation, README and GitHub Actions documentation, and the existing release workflow.
- Adds a release-build task and release workflow tests or task-level verification as appropriate.
- Does not add a package manager, GoReleaser, a network call from the CLM executable, or any CLM tag-push behavior. Stable release tags remain exclusively CLM-created; CI creates only temporary PR prerelease tags.
