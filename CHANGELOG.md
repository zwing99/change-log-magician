# Changelog

All notable changes to this project will be documented in this file.

This changelog is maintained with [Change Log Magician](https://github.com/zwing99/change-log-magician) (`clm`). Add one fragment for every change, then let CLM assemble releases.

Helpful commands:

- `clm new-patch <name>` creates a release-note fragment.
- `clm added -m "Description"` adds a categorized entry.
- `clm cut --check` validates the next release without changing files.

## [Unreleased]

## [0.0.1] - 2026-09-12

### Added

- Initial release of Change Log Magician.
- Create Keep a Changelog fragments with explicit major, minor, patch, or unreleased intent.
- Append Added, Changed, Deprecated, Removed, Fixed, and Security notes from the CLI.
- Cut validated changelog releases with archived fragments, release commits, and annotated Git tags.
- Validate pull requests and merge requests, preview versions, and derive prerelease artifact versions.
- Generate shell completions and provide accessible color-aware terminal output.
- Include CLM maintenance guidance and command hints in initialized changelogs.
