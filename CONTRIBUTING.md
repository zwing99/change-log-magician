# Contributing

Install the repository toolchain with `mise install`, then run `mise run check` before opening a pull request. Keep changes focused, formatted, and accompanied by tests for observable behavior.

Keep Go comments short and useful: explain intent or a non-obvious constraint, not a line that readers can already understand from the code.

Every pull request also needs a `changelog.d/*.md` fragment. Use one of `clm new-major`, `clm new-minor`, `clm new-patch`, or `clm new-unreleased`, then add categorized entries with commands such as `clm added -m "..."`.

CI validates this with:

```bash
clm cut --check --require-fragment --base "$BASE_REVISION"
```

Stable releases are cut only after merge by an authorized workflow; CLM commits and tags locally, while that workflow pushes the result.
