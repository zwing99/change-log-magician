## 1. Release-cut trigger isolation

- [x] 1.1 Update the main-branch cut workflow to include active changelog fragments and explicitly exclude `changelog.d/archive/**`; verify the workflow-boundary test covers both paths and an archive-only release commit cannot invoke `clm cut`.

## 2. Setup workflow tag isolation

- [x] 2.1 Limit Copilot setup push events to branches while retaining its workflow-file path filter, pull-request trigger, and manual dispatch; verify the workflow-boundary test confirms release tags cannot start it.

## 3. Verification

- [x] 3.1 Run `mise run check` and `openspec validate fix-release-workflow-triggers --strict`; verify the workflow YAML and OpenSpec artifacts remain valid before handoff.
