## Why

The first stable cut correctly produced and published `v0.0.2`, but its
archive-only release commit started another cut workflow that failed because no
active fragments remained. Tag pushes also started Copilot setup jobs despite
their path filter, adding unnecessary work to every release.

## What Changes

- Exclude archived changelog fragments from the main-branch release-cut trigger.
- Restrict Copilot setup's push trigger to branches while retaining pull-request
  validation and manual dispatch.
- Add workflow-level regression checks for both trigger boundaries.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `release-artifact-publishing`: Ensure the main cut workflow ignores its own
  archive-only commit and ancillary setup automation ignores release tags.

## Impact

- Affects the release-cut and Copilot setup GitHub Actions workflows and their
  workflow-boundary tests.
- Does not change CLM commands, tag ownership, artifact contents, or release
  publishing behavior.
