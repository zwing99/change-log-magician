## 1. Local release assembly

- [x] 1.1 Add a `mise run release-build` task and a repository-owned build script that require an explicit stable or prerelease version tag, cross-compile the five supported targets with the embedded version, and write named artifacts plus SHA-256 checksums to `dist/`; verify valid stable, valid beta, and invalid version invocations locally.
- [x] 1.2 Add automated coverage for the build script's version validation, target artifact names/extensions, embedded version, and checksum coverage; verify it runs through `mise run check` or a documented dedicated verification task.
- [x] 1.3 Add a fixture-backed release-notes extraction helper that returns exactly one stable changelog entry body while preserving category Markdown; verify multiple-version and missing-version cases with table-driven tests.

## 2. Main-branch cut workflow

- [x] 2.1 Add a qualifying `main` push workflow that checks out full history, runs `mise run check`, and invokes `clm cut`; verify its path filters process active fragments but not CLM's archive-only release commit.
- [x] 2.2 Make the cut workflow push only the resulting release commit and CLM-created `vX.Y.Z` and `clm/vX.Y.Z` tags for versioned cuts, and push only the release commit for unreleased-only cuts; use a GitHub App token or scoped PAT so the pushed stable tag starts the tag workflow, and verify no workflow step calculates or creates tags.

## 3. Pull-request beta workflow

- [x] 3.1 Extend the pull-request workflow to run `clm cut --check --require-fragment` and derive `vX.Y.Z-beta.pr.<number>.<sha>` with `clm next-version --json`; verify malformed fragments and missing fragments fail before any beta publication.
- [x] 3.2 For same-repository versioned pull requests, create and push the CI-owned beta tag, build it through `mise run release-build`, and publish a downloadable GitHub prerelease within the pull-request workflow; verify the beta's embedded version and assets match its tag and its push does not start the stable workflow.
- [x] 3.3 Restrict beta publication to repository pull requests and skip it for forks and unreleased-only plans; verify those cases retain validation without write credentials, tags, artifacts, or prereleases.

## 4. Stable-tag artifact workflow

- [x] 4.1 Replace the inline-build release workflow with a workflow whose effective trigger accepts only stable `vX.Y.Z` tags and an authorized manual existing-tag retry; verify beta and `clm/vX.Y.Z` tag pushes do not start its release job, it checks out the requested stable tag, and it never invokes `clm cut` or pushes Git state.
- [x] 4.2 Invoke `mise run release-build` with the stable tag and create or update the matching GitHub Release assets idempotently; verify a retry for an already-pushed tag does not cut another release or create another tag.
- [x] 4.3 Generate the GitHub Release body from the checked-out tag's changelog entry and publish it unchanged as the release description; verify the fixture-backed extraction output is the workflow input.

## 5. Documentation and dogfood release

- [x] 5.1 Update README and GitHub Actions CI documentation to describe local artifact builds, beta prereleases and their fork boundary, the main-branch cut workflow, stable-tag publication and retry, the authorized push credential, CLM tag ownership, and the GitHub Actions publishing boundary; verify all documented commands match task and workflow names.
- [ ] 5.2 Add this change's patch fragment through `clm new-patch` and `clm added`; open a repository pull request to verify its beta prerelease, then merge it and confirm the main workflow cuts it, the stable-tag workflow publishes it, and the stable GitHub release description matches the CLM-generated changelog entry.
- [x] 5.3 Run `mise run check`, `openspec validate automate-release-builds --strict`, and the release-build verification before handoff.
