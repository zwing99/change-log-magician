## 1. Project foundation

- [x] 1.1 Initialize the Go module, `cmd/clm` entrypoint, domain-led internal package layout, `.gitignore`, and a pinned `mise.toml`; verify `mise run check` can run formatting, static analysis, and tests.
- [x] 1.2 Add an MIT `LICENSE`, `README.md`, `AGENTS.md`, and `CONTRIBUTING.md` covering installation, quality expectations, fragment workflow, and GitHub/GitLab CI examples; verify all documented commands and paths match the executable.
- [x] 1.3 Build the Cobra command tree, global `--color` mode, structured-output convention, and completion subcommands; verify help and generated Bash, Zsh, Fish, and PowerShell completion fixtures.
- [x] 1.4 Implement terminal output mode detection, semantic color, and interactive-only block-letter banner behavior; verify `NO_COLOR`, `--color=never`, forced color, and piped output in unit tests.

## 2. Changelog domain and fragment authoring

- [x] 2.1 Define the pure changelog document, release-level, category, entry, and validation-error model; verify table-driven tests cover all four levels and six categories.
- [x] 2.2 Implement parsing and rendering of manually editable fragment Markdown, including reordered categories and actionable malformed-file diagnostics; verify golden fixtures for valid and invalid fragments.
- [x] 2.3 Implement root `CHANGELOG.md` parsing/rendering that retains an `Unreleased` section and preserves content outside CLM-owned sections; verify golden tests retain manual preamble, links, and unrelated content.
- [x] 2.4 Implement `clm init` with safe non-overwriting creation of `CHANGELOG.md`, `changelog.d/`, archive structure, and default `clm.toml`; verify temporary-repository tests for new and existing files.
- [x] 2.5 Implement `new-major`, `new-minor`, `new-patch`, and `new-unreleased` with collision-safe filenames and starter fragment content; verify each command creates the expected valid fragment.
- [x] 2.6 Implement category entry commands with `-m`, interactive prompt/editor handling, explicit target selection, and single-fragment inference; verify append behavior, ambiguous-target failures, and no empty entries.

## 3. Release planning and stable cuts

- [x] 3.1 Implement SemVer discovery from reachable stable tags and next-version calculation using implicit `0.0.0`; verify first patch/minor/major and mixed-level precedence cases.
- [x] 3.2 Implement a shared release-plan builder that discovers active fragments, validates the full pending set, computes archive destinations, and produces a non-mutating rendered changelog plan; verify no-fragment, malformed, and collision cases.
- [x] 3.3 Implement unreleased-only cut planning and rendering, including archived source fragments and a retained `Unreleased` section; verify it creates no stable version in integration tests.
- [x] 3.4 Implement versioned-cut rendering that promotes existing root unreleased entries plus all pending fragments into a dated numbered release and recreates empty `Unreleased`; verify golden output and archive layout.
- [x] 3.5 Implement Git preflight, Git-aware archive moves, release commit creation, and annotated `vX.Y.Z` plus `clm/vX.Y.Z` tags without pushing; verify a temporary Git-repository integration test observes the commit, moved files, and tags.
- [x] 3.6 Add failure handling and recovery diagnostics for dirty state, archive/tag collisions, and post-commit Git failures; verify preflight failures leave repository files and tags unchanged.

## 4. Pipeline checks and prerelease planning

- [x] 4.1 Implement `clm cut --check` as a non-mutating release-plan validation path and `--require-fragment --base` branch-diff enforcement; verify success, missing-MR-fragment, and unrelated-invalid-fragment Git fixtures.
- [x] 4.2 Implement machine-readable `clm next-version` output for stable plans and explicit unreleased-only no-version plans; verify JSON schema and no-write behavior.
- [x] 4.3 Implement caller-identified SemVer prerelease derivation from a stable plan and reject unreleased-only prerelease requests; verify identifiers such as `pr.123.abcdef0` and invalid identifier errors.
- [x] 4.4 Document and test GitHub Actions and GitLab CI examples for full tag history, MR/PR validation, version-preview artifact naming, prerelease publication handoff, and protected post-merge push credentials; verify examples use only CLM's documented public commands.

## 5. Quality and release readiness

- [x] 5.1 Add table-driven unit coverage for domain, version, parser, renderer, and error paths plus golden fixtures for full changelog output; verify `go test ./...` passes with race detection where applicable.
- [x] 5.2 Add lint, formatting, static-analysis, and test tasks to the aggregate Mise check and CI workflow; verify the initial repository's CI invokes `mise run check`.
- [x] 5.3 Configure GitHub-based CLM binary release automation with reproducible build artifacts and checksums; verify the workflow is triggered from the documented release-tag path without coupling CLM's runtime to GitHub APIs.
