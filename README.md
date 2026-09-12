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

## Quick start

```bash
clm init
clm new-minor shell-completions
clm added -m "Add shell completion support"
clm cut --check --require-fragment --base origin/main
clm next-version --json
clm cut
```

Fragments live directly under `changelog.d/`. They use one level heading (`[major]`, `[minor]`, `[patch]`, or `[unreleased]`) and the standard Keep a Changelog categories. Every merge request or pull request must contribute a valid fragment. Use `new-unreleased` for changes that should appear under the persistent root `Unreleased` section without producing a tag by themselves.

`clm cut` creates the changelog/archive commit. A stable cut creates annotated `vX.Y.Z` and `clm/vX.Y.Z` tags but never pushes them. CI must provide credentials only to the post-merge job that pushes those results.

## Pipelines

GitHub Actions and GitLab CI both need a checkout with reachable tags and a target-base revision:

```bash
clm cut --check --require-fragment --base "$BASE_REVISION"
clm next-version --json --prerelease-id "pr.123.$SHORT_SHA"
```

For an unreleased-only change, `next-version` reports no stable version and the pipeline must skip prerelease publication. CLM does not call forge APIs or publish artifacts; the workflow uses its output to name/publish them.

Copy the ready-to-adapt examples from [`docs/ci/github-actions.yml`](docs/ci/github-actions.yml) and [`docs/ci/gitlab-ci.yml`](docs/ci/gitlab-ci.yml). The authorized post-merge job runs `clm cut`, then pushes the resulting commit and tags.

Generate completion scripts with `clm completion bash|zsh|fish|powershell`.
