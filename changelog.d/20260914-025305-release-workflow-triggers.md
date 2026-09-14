# [patch]

## Added

- Prevent the main release-cut workflow from rerunning on CLM's archive-only release commits.
- Prevent Copilot setup jobs from running for stable, companion, and prerelease tag pushes.
- Require pull-request branches to rebase on `origin/main` and review their changelog fragment before updates.
- Require release-affecting changes to verify the release credential before merge and the complete cut-to-published-release path afterward.
