## Context

The initial stable cut proved the artifact pipeline, but the cut workflow's
fragment path filter matched an archived fragment in its own release commit.
GitHub also evaluates the Copilot setup workflow for tag pushes despite its
path-only push filter.

## Goals / Non-Goals

**Goals:**

- Prevent archive-only commits from reaching `clm cut`.
- Keep Copilot setup available for workflow-file changes on branches and pull
  requests without running it for release tags.
- Protect both boundaries with repository tests that inspect workflow YAML.

**Non-Goals:**

- Change CLM release planning, version calculation, tag ownership, or artifact
  publication.
- Suppress ordinary CI runs on the release commit.

## Decisions

### Use ordered negative path matching for archived fragments

The cut workflow will match the changelog fragment tree and explicitly exclude
`changelog.d/archive/**`. This is clearer than relying on a single wildcard
that can match nested paths, and leaves active fragments eligible to trigger a
cut.

### Limit Copilot setup pushes to branches

The setup workflow will retain its existing workflow-file path filter and add a
branch filter. Pull-request and manual-dispatch triggers remain unchanged.
This preserves setup validation for authored workflow changes while excluding
all tag pushes.

### Test trigger contracts as workflow source

The existing Go workflow-boundary tests will assert the archive exclusion and
branch restriction. GitHub trigger semantics cannot be executed locally, so
source-level contract tests prevent accidental removal of the guardrails.

## Risks / Trade-offs

- [An incorrectly ordered negative pattern re-includes archives] → Keep the
  exclusion after the inclusive fragment pattern and test both strings.
- [A future desired tag-triggered setup task is blocked] → Add a separate,
  intentionally tag-triggered workflow rather than broadening setup automation.

## Migration Plan

1. Update the two workflow trigger definitions and their tests.
2. Verify with `mise run check` and a workflow syntax review.
3. Merge normally; no rerun is needed for the completed `v0.0.2` release.
