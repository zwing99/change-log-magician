# Agent Guide

- Use the Go version and tasks in `mise.toml`.
- Keep the module path `github.com/zwing99/change-log-magician`; the executable name is `clm`.
- Keep command handlers thin; keep changelog and release rules pure and unit-tested.
- Keep `cut`, `cut --check`, and `next-version` on the same release-plan logic.
- Comment intent, invariants, and surprising constraints in short plain-language sentences. Do not narrate obvious code or add verbose block comments.
- Add table-driven tests and fixtures for behavior changes. Use temporary Git repositories for Git integration tests.
- When fixing a bug, add a regression test that demonstrates the failure and protects the intended behavior.
- Run `mise run check` before handing off changes.
- Before opening or updating a pull request, fetch and rebase its branch onto `origin/main`.
- When ignore rules need to change, generate the base template with gitignore.io and retain its source URL in `.gitignore`.
- Do not add hidden state, silently rewrite user-authored changelog content, or make network/forge API calls from CLM.
- CLM never pushes. CI owns credentials and pushes the release commit and tags after `clm cut` succeeds.
- Before merging a release-affecting change, verify required repository secrets (currently `RELEASE_PUSH_TOKEN`) are configured and usable.
- After merging a release-affecting change, verify the cut workflow, stable tags, stable publisher, release assets, and changelog-derived GitHub Release description.
- Stable cuts require a clean Git worktree; cover release behavior with temporary Git repositories in tests.
- Keep terminal output automation-safe: JSON and completion output must not contain banners or ANSI decoration.
