# Change Log Magician

![Change Log Magician banner](assets/clm-banner.png)

Change Log Magician (`clm`) keeps release notes tidy while your team ships. Write small, merge-friendly Markdown fragments during development; CLM turns them into an audited Keep a Changelog release when it is time to cut.

No shared changelog fights. No guessing the next version. Just a little release magic.

## Install and develop

```bash
mise install
mise run check
go run ./cmd/clm --help
```

The committed `mise.toml` is the source of truth for Go and quality-tool versions.

Build the five downloadable release artifacts locally with an explicit stable or beta tag:

```bash
mise run release-build -- v1.2.3
mise run release-build -- v1.2.3-beta.pr.42.abcdef
mise run release-build-check
```

Artifacts and SHA-256 checksums are written to `dist/`.

## Quick start

```bash
clm init
clm new-minor shell-completions
clm added -m "Add shell completion support"
clm cut --check --require-fragment --base origin/main
clm next-version --json
clm cut
```

## Changelog fragments for an MR

Each merge request (MR) or pull request adds one small Markdown fragment directly under `changelog.d/`, alongside its code change. The fragment is the release-note source for that MR: it tells CLM both whether the change should create a stable release and what users should read in its changelog entry. Do not edit `CHANGELOG.md` as part of the MR; CLM assembles it when a release is cut.

Use the release intent that best describes the MR:

| Command | Use when the MR contains | Result at the next cut |
| --- | --- | --- |
| `clm new-major` | a breaking change | next major version |
| `clm new-minor` | a backwards-compatible feature | next minor version |
| `clm new-patch` | a backwards-compatible fix or small change | next patch version |
| `clm new-unreleased` | a note that belongs in `Unreleased` but must not create a stable tag | no stable version by itself |

For example, a feature MR can create its fragment and add user-visible notes as it is developed:

```bash
clm new-minor shell-completions
clm added -m "Add shell completion support"
clm fixed -m "Preserve quoted completion arguments"
```

`new-minor` creates a timestamped file, such as `changelog.d/20260915-120000-shell-completions.md`, with the required level heading. Each entry command adds a bullet to the correct category and writes the category heading when needed. The completed file looks like this:

```md
# [minor]

## Added

- Add shell completion support

## Fixed

- Preserve quoted completion arguments
```

The first heading must be exactly one of `[major]`, `[minor]`, `[patch]`, or `[unreleased]`. Entries are bullets under the standard Keep a Changelog categories: `Added`, `Changed`, `Deprecated`, `Removed`, `Fixed`, or `Security`. Empty categories are omitted. Keep entries concise and user-facing: describe the behavior someone receives, rather than an implementation detail or MR number.

CLM normally selects the only active fragment automatically, so a patch MR can be authored without opening the file:

```bash
clm new-patch escaped-json
clm fixed -m "Escape control characters in JSON output"
clm changed -m "Use the configured indentation for JSON output"
```

If your checkout has more than one active fragment, choose the MR's file explicitly with `--change` (or `-c`):

```bash
clm security --change changelog.d/20260915-120000-escaped-json.md \
  --message "Reject malformed signature headers"
```

Omit `--message` to enter a short entry interactively, or use `--editor` to compose one in `$EDITOR`. You can also edit a fragment by hand, provided it retains the required level heading, recognized category headings, and non-empty `- ` bullet entries.

Before pushing or opening the MR, validate the exact change against its target branch:

```bash
clm cut --check --require-fragment --base origin/main
```

Commit the generated `changelog.d/*.md` file with the MR. After the MR merges, the authorized release workflow runs `clm cut`; it combines active fragments into `CHANGELOG.md`, archives the consumed files, and creates release tags when a major, minor, or patch fragment is present.

`clm cut` creates the changelog/archive commit. A stable cut creates annotated `vX.Y.Z` and `clm/vX.Y.Z` tags but never pushes them. CI must provide credentials only to the post-merge job that pushes those results.

## Pipelines

GitHub Actions and GitLab CI both need a checkout with reachable tags and a target-base revision:

```bash
clm cut --check --require-fragment --base "$BASE_REVISION"
clm next-version --json --prerelease-id "pr.123.$SHORT_SHA"
```

For an unreleased-only change, `next-version` reports no stable version and the pipeline must skip prerelease publication. CLM does not call forge APIs or publish artifacts; the workflow uses its output to name/publish them.

Copy the ready-to-adapt examples from [`docs/ci/github-actions.yml`](docs/ci/github-actions.yml) and [`docs/ci/gitlab-ci.yml`](docs/ci/gitlab-ci.yml). On a repository pull request, GitHub Actions validates fragments, derives a `vX.Y.Z-beta.pr.<number>.<sha>` version with CLM, and publishes its CI-owned beta tag and prerelease assets. Fork pull requests validate only: they receive no write credential, tag, or prerelease.

On a qualifying `main` push, the authorized cut workflow runs `mise run check` and `clm cut`. CLM creates the release commit plus `vX.Y.Z` and `clm/vX.Y.Z`; the workflow, authenticated with a GitHub App token or scoped `RELEASE_PUSH_TOKEN` PAT, pushes that exact state. CLM never pushes or calls GitHub APIs.

Only the CLM-created stable `vX.Y.Z` tag starts the stable publisher. It builds with `mise run release-build`, extracts the matching changelog entry as the GitHub Release body, and creates or updates the release. Beta and `clm/vX.Y.Z` tags cannot start that workflow. Maintainers can use its manual dispatch with an existing stable tag to rebuild and republish without cutting another release or pushing Git state.

Generate completion scripts with `clm completion bash|zsh|fish|powershell`.
