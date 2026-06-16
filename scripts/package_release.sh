#!/bin/sh
set -eu

if [ "$#" -lt 1 ]; then
  echo "usage: scripts/package_release.sh <version> [out_dir]" >&2
  exit 2
fi

version="$1"
out_dir="${2:-dist/release}"
root="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
cache_dir="${GOCACHE:-$root/.gocache}"

mkdir -p "$root/$out_dir"

build_one() {
  goos="$1"
  goarch="$2"
  work_dir="$(mktemp -d "${TMPDIR:-/tmp}/uncut-${version}-${goos}-${goarch}.XXXXXX")"

  env \
    GOOS="$goos" \
    GOARCH="$goarch" \
    CGO_ENABLED=0 \
    GOCACHE="$cache_dir" \
    go build -buildvcs=false -ldflags="-s -w" -o "$work_dir/uncut" "$root"

  cp "$root/README.md" "$work_dir/README.md"
  cp "$root/METHODS.md" "$work_dir/METHODS.md"
  cp -R "$root/docs" "$work_dir/docs"
  mkdir -p "$work_dir/man"
  cp "$root/man/uncut.1" "$work_dir/man/uncut.1"

  tar -C "$work_dir" -czf "$root/$out_dir/uncut-${goos}-${goarch}.tar.gz" \
    uncut README.md METHODS.md docs man

  echo "$root/$out_dir/uncut-${goos}-${goarch}.tar.gz"
}

build_one darwin arm64
build_one darwin amd64
build_one linux amd64
build_one linux arm64

