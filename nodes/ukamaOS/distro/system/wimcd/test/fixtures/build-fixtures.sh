#!/bin/sh
set -eu

ROOT="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)"
OUT="$ROOT/tarballs"

rm -rf "$OUT"
mkdir -p "$OUT"

tar -czf "$OUT/example_v1-abc.tar.gz" \
    -C "$ROOT/pkgs/example_v1-abc" .

tar -czf "$OUT/example_v1-xyz.tar.gz" \
    -C "$ROOT/pkgs/example_v1-xyz" .

tar -czf "$OUT/example_bad_missing_version.tar.gz" \
    -C "$ROOT/pkgs/example_bad_missing_version" .

tar -czf "$OUT/example_bad_wrong_version.tar.gz" \
    -C "$ROOT/pkgs/example_bad_wrong_version" .
