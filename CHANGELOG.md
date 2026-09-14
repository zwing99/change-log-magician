# Changelog

All notable changes to this project will be documented in this file.

This changelog is maintained with [Change Log Magician](https://github.com/zwing99/change-log-magician) (`clm`). Add one fragment for every change, then let CLM assemble releases.

Helpful commands:

- `clm new-patch <name>` creates a release-note fragment.
- `clm added -m "Description"` adds a categorized entry.
- `clm cut --check` validates the next release without changing files.

## [Unreleased]

### Added

- Archive the completed release workflow trigger OpenSpec change.

### Fixed

- Avoid failed release-cut jobs when a push contains only archived changelog fragments.

## [0.0.3] - 2026-09-14

### Added

- Prevent the main release-cut workflow from rerunning on CLM's archive-only release commits.
- Prevent Copilot setup jobs from running for stable, companion, and prerelease tag pushes.
- Require pull-request branches to rebase on `origin/main` and review their changelog fragment before updates.
- Require release-affecting changes to verify the release credential before merge and the complete cut-to-published-release path afterward.

## [0.0.2] - 2026-09-14

### Added

- Publish durable OpenSpec specifications for CLM's changelog, release, pipeline, and CLI capabilities.
- Automate cross-platform release builds, pull-request beta prereleases, and stable GitHub Releases with changelog-derived descriptions.

## [0.0.1] - 2026-09-12

### Added

- Initial release of Change Log Magician.
- Create Keep a Changelog fragments with explicit major, minor, patch, or unreleased intent.
- Append Added, Changed, Deprecated, Removed, Fixed, and Security notes from the CLI.
- Cut validated changelog releases with archived fragments, release commits, and annotated Git tags.
- Validate pull requests and merge requests, preview versions, and derive prerelease artifact versions.
- Generate shell completions and provide accessible color-aware terminal output.
- Include CLM maintenance guidance and command hints in initialized changelogs.
