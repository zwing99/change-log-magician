#!/usr/bin/env bash

set -euo pipefail

version="${1:-}"
if [[ ! "$version" =~ ^v(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(-([0-9A-Za-z-]+)(\.[0-9A-Za-z-]+)*)?$ ]]; then
  echo "release version must be a stable or prerelease vX.Y.Z tag" >&2
  exit 1
fi
if [[ "$version" == *-* ]]; then
  prerelease="${version#*-}"
  IFS='.' read -r -a identifiers <<< "$prerelease"
  for identifier in "${identifiers[@]}"; do
    if [[ "$identifier" =~ ^[0-9]+$ && "$identifier" != "0" && "$identifier" == 0* ]]; then
      echo "release version must use SemVer prerelease identifiers" >&2
      exit 1
    fi
  done
fi

# Validate before replacing dist so malformed input never publishes artifacts.
rm -rf dist
mkdir -p dist

for target in linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64; do
  os="${target%/*}"
  arch="${target#*/}"
  extension=""
  if [[ "$os" == "windows" ]]; then
    extension=".exe"
  fi
  CGO_ENABLED=0 GOOS="$os" GOARCH="$arch" go build -trimpath -ldflags "-s -w -X main.version=$version" -o "dist/clm-${version}-${os}-${arch}${extension}" ./cmd/clm
done

(
  cd dist
  shasum -a 256 clm-* > checksums.txt
)
