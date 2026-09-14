#!/usr/bin/env bash

set -euo pipefail

assert_artifacts() {
  local version="$1"
  local linux_binary="dist/clm-${version}-linux-amd64"
  local artifact
  for artifact in \
    "$linux_binary" \
    "dist/clm-${version}-linux-arm64" \
    "dist/clm-${version}-darwin-amd64" \
    "dist/clm-${version}-darwin-arm64" \
    "dist/clm-${version}-windows-amd64.exe"; do
    test -f "$artifact"
    grep -F "$(basename "$artifact")" dist/checksums.txt >/dev/null
  done
  # Cross-compiled binaries are not necessarily executable on the host.
  grep -aF "$version" "$linux_binary" >/dev/null
}

mise run release-build -- v1.2.3
assert_artifacts v1.2.3

mise run release-build -- v1.2.3-beta.pr.42.abcdef
assert_artifacts v1.2.3-beta.pr.42.abcdef

for invalid_version in invalid-version v1.2.3-beta.01; do
  if mise run release-build -- "$invalid_version"; then
    echo "malformed version unexpectedly built artifacts" >&2
    exit 1
  fi
done
