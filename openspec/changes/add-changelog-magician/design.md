## Context

This is a greenfield CLI. The proposal defines a Keep a Changelog-compatible fragment workflow, stable release cuts, and GitHub/GitLab pipeline support. The implementation must be idiomatic Go, deliberately small, and thoroughly unit-tested; it must also tolerate user-authored Markdown rather than treating generated text as the only valid input.

## Goals / Non-Goals

**Goals:**

- Keep release-note rules in a pure, testable domain layer.
- Make a release cut predictable and safe: validate first, then render, archive, commit, and tag.
- Keep CI integration limited to local Git state and explicit command output so both GitHub and GitLab can invoke it.
- Make the default interaction polished while keeping noninteractive output stable for automation.

**Non-Goals:**

- Hosting releases, artifacts, or changelogs in a GitHub/GitLab API.
- Managing prerelease publication, pushing to remotes, or credentials inside CLM.
- Supporting prerelease promotion, package-manager publishing, or arbitrary changelog taxonomies in the first release.
- Reformatting or interpreting manually maintained root-changelog content outside CLM-managed release sections.

## Decisions

### Go executable with domain-led package boundaries

Use a single Go module and a thin `cmd/clm` entrypoint. Keep commands, terminal interaction, and serialization at the edge; keep fragment parsing, changelog rendering, version planning, and release planning free of Cobra, Git, filesystem, and terminal concerns. Organize internals around `cli`, `changelog`, `release`, `gitx`, and `filesystem` responsibilities.

This makes table-driven unit tests the default for behavior and limits interfaces to real side-effect boundaries. Cobra is appropriate for commands and completion generation; a Go-native banner implementation avoids requiring an external `figlet` binary.

### Stable, constrained Markdown format

Fragments use one level heading (`[major]`, `[minor]`, `[patch]`, or `[unreleased]`) and the standard six category headings. The parser accepts category ordering differences and omitted empty sections, but rejects ambiguous structure. Rendering uses an internal document model and preserves authored entry text. The root parser edits only CLM-owned `Unreleased` and numbered release sections, preserving all unrelated content.

This balances manual editability with the structural guarantees required for automated release assembly. Free-form Markdown parsing was rejected because it cannot provide reliable validation or deterministic rendering.

### Explicit fragment selection with safe convenience

Entry commands accept an explicit fragment target. If exactly one active fragment exists, CLM selects it automatically; otherwise it fails with a selection prompt. `-m` supplies one-line text, while omitted text follows an interactive prompt/editor path.

Implicitly tracking a current fragment by branch or hidden state was rejected because it introduces surprising state and makes CI/reproducibility harder.

### Release plan before release mutation

`cut`, `cut --check`, and `next-version` share one release-plan builder. It discovers only active fragments, validates all of them, finds the latest stable `vX.Y.Z` tag, applies the highest requested bump against `0.0.0` when absent, and builds a proposed root-changelog document and archive plan in memory.

`cut --check` and `next-version` serialize that plan without writes. `cut` performs all preflight checks before it mutates files. This avoids drift between CI validation and release behavior.

### Unreleased promotion semantics

An unreleased-only cut adds pending entries to root `Unreleased`, archives those source fragments under a unique unreleased archive run, and creates a release commit without tags. Once any versioned fragment is present, a stable cut combines existing `Unreleased` entries and all active fragment entries into the numbered release; it then leaves a fresh empty `Unreleased` section and archives every processed fragment under `changelog.d/archive/vX.Y.Z/`.

This makes unversioned changes visible immediately while ensuring they are included in the next stable release. Leaving current unreleased fragments behind during a versioned cut was rejected because it would make the future release boundary unclear.

### Git-owned cut boundary

Before a cut, CLM requires a suitable repository state, checks that target tags and archive destinations do not exist, and validates its complete release plan. It writes the changelog, uses Git-aware moves for source fragments, creates one release commit, then creates annotated `vX.Y.Z` and `clm/vX.Y.Z` tags for stable cuts. CLM never pushes; GitHub Actions or GitLab CI owns credentialed pushes.

Creating only working-tree changes was rejected because the cut must reliably bind the rendered changelog and archive to the tag's commit. The second namespaced tag provides CLM-specific release provenance while retaining conventional SemVer tags.

### Forge-neutral prerelease inputs

Version planning accepts a base revision for MR/PR fragment enforcement and produces JSON suitable for CI. A caller-supplied identifier, such as `pr.123.abcdef0`, is appended to a planned stable version to form a SemVer prerelease. CLM does not invoke forge APIs or publish prereleases. An unreleased-only set has no planned stable version and therefore no prerelease.

This avoids coupling core behavior to GitHub Actions or GitLab CI variables while allowing either workflow to name and publish its own artifacts.

### Reproducible toolchain and output behavior

`mise.toml` pins the Go toolchain and offers formatting, test, lint, and aggregate check tasks. Human terminal output uses semantic color with `auto|always|never` behavior, honors `NO_COLOR`, and never relies on color alone. Banners are limited to interactive discovery; scripts, JSON, pipelines, and piped output remain decoration-free.

## Risks / Trade-offs

- [Manual Markdown can exceed the supported structure] → Validate with file-specific errors, document the supported form, and preserve content outside CLM-owned structures.
- [A Git failure after a release commit can leave a partial cut] → Perform all possible checks before mutation, make tags only after commit success, and print exact recovery steps for post-commit failures.
- [Shallow CI clones can lack historical tags] → Document full-history/tag fetch requirements in GitHub and GitLab examples and report a clear Git discovery error.
- [An archive move can conflict with an existing file] → Preflight every destination and never overwrite archive content.
- [Color or banners can corrupt machine-readable output] → Centralize output mode selection and cover TTY/non-TTY/JSON cases in tests.

## Migration Plan

1. Publish the initial CLM binary from GitHub with checksums and installation instructions.
2. Users run `clm init` in an existing Git repository to add CLM configuration and changelog structure without replacing existing changelog content.
3. Teams add the documented GitHub Actions or GitLab CI validation command to MR/PR pipelines.
4. Release automation grants push permission only to the post-merge stable-cut job; validation and version-preview jobs remain read-only.

Rollback consists of reverting the release commit and deleting the two local/remote stable tags before any downstream artifact publication. CLM will not push automatically.
