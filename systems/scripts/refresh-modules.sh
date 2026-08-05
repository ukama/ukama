#!/usr/bin/env bash
#
# Refresh Go modules for every ukama service (gRPC services and REST
# gateways): delete go.sum, strip all require sections from go.mod so the
# dependency graph is rebuilt from scratch, then run the full local
# pipeline: go mod tidy && make lint && make gen && make test && make.
#
# Usage:
#   ./refresh-modules.sh [options] [MODULE_DIR...]
#
#   MODULE_DIR  paths relative to systems/ (e.g. registry/network). When
#               omitted, every directory under systems/ containing both a
#               go.mod and a Makefile is processed (systems/common is
#               excluded — it is a library, not a service; it still gets
#               a plain 'go mod tidy' first so services resolve against a
#               clean common).
#
# Options:
#   --skip-lint   skip 'make lint'
#   --skip-gen    skip 'make gen'  (needs protoc/mockery installed)
#   --skip-test   skip 'make test'
#   --tidy-only   only strip + 'go mod tidy' (no make targets at all)
#
# Failures don't stop the run; a summary is printed at the end.
# go.mod edits are destructive by design — git is your undo.

set -uo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SYSTEMS_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"

DO_LINT=1 DO_GEN=1 DO_TEST=1 TIDY_ONLY=0
MODULES=()

while [[ $# -gt 0 ]]; do
  case "$1" in
    --skip-lint) DO_LINT=0; shift ;;
    --skip-gen)  DO_GEN=0; shift ;;
    --skip-test) DO_TEST=0; shift ;;
    --tidy-only) TIDY_ONLY=1; shift ;;
    -h|--help) grep '^#' "$0" | sed 's/^# \{0,1\}//'; exit 0 ;;
    *) MODULES+=("$1"); shift ;;
  esac
done

log()  { printf '\n\033[1;34m==> %s\033[0m\n' "$*"; }

# Discover modules when none given: any dir (2 levels deep) with go.mod + Makefile.
if [[ ${#MODULES[@]} -eq 0 ]]; then
  while IFS= read -r modfile; do
    dir="$(dirname "$modfile")"
    rel="${dir#"$SYSTEMS_DIR"/}"
    [[ "$rel" == "common" || "$rel" == common/* ]] && continue
    [[ -f "$dir/Makefile" ]] || continue
    MODULES+=("$rel")
  done < <(find "$SYSTEMS_DIR" -mindepth 2 -maxdepth 3 -name go.mod -not -path '*/node_modules/*' | sort)
fi

command -v go >/dev/null || { echo "go not found in PATH" >&2; exit 1; }
if [[ $TIDY_ONLY -eq 0 && $DO_LINT -eq 1 ]] && ! command -v golangci-lint >/dev/null; then
  echo "WARNING: golangci-lint not found - 'make lint' will fail. Use --skip-lint to skip." >&2
fi
if [[ $TIDY_ONLY -eq 0 && $DO_GEN -eq 1 ]] && ! command -v mockery >/dev/null; then
  echo "WARNING: mockery not found - 'make gen' may fail. Use --skip-gen to skip." >&2
fi

# strip_requires FILE - remove every 'require ( ... )' block and every
# single-line 'require x vN' directive from go.mod, keep everything else.
strip_requires() {
  awk '
    /^require[ \t]*\(/ { inblock=1; next }
    inblock && /^\)/   { inblock=0; next }
    inblock            { next }
    /^require[ \t]+[^ \t(]/ { next }
    { print }
  ' "$1" > "$1.tmp" && mv "$1.tmp" "$1"
}

# Clean common first so every service resolves against a tidy library.
log "systems/common: go mod tidy"
(cd "$SYSTEMS_DIR/common" && GOFLAGS=-mod=mod go mod tidy) || echo "WARNING: tidy failed in common"

FAILED=() SKIPPED=()
for rel in "${MODULES[@]}"; do
  dir="$SYSTEMS_DIR/$rel"

  # Skip anything that isn't a Go service: needs a go.mod and at least
  # one .go source file (filters out node/BFF services, ory-based auth,
  # chart-only or doc directories passed explicitly).
  if [[ ! -f "$dir/go.mod" ]]; then
    echo "skipping $rel: not a Go service (no go.mod)"
    SKIPPED+=("$rel")
    continue
  fi
  if ! find "$dir" -name '*.go' -not -path '*/node_modules/*' -print -quit | grep -q .; then
    echo "skipping $rel: not a Go service (no .go files)"
    SKIPPED+=("$rel")
    continue
  fi

  log "$rel"
  (
    cd "$dir"
    rm -f go.sum
    strip_requires go.mod
    GOFLAGS=-mod=mod go mod tidy || exit 1
    if [[ $TIDY_ONLY -eq 0 ]]; then
      [[ $DO_LINT -eq 1 ]] && { make lint || exit 1; }
      [[ $DO_GEN  -eq 1 ]] && { make gen  || exit 1; }
      [[ $DO_TEST -eq 1 ]] && { make test || exit 1; }
      make || exit 1
    fi
  ) || FAILED+=("$rel")
done

echo
if [[ ${#SKIPPED[@]} -gt 0 ]]; then
  printf 'Skipped (not Go services): %s\n' "${SKIPPED[*]}"
fi
if [[ ${#FAILED[@]} -gt 0 ]]; then
  printf '\033[1;31mFAILED (%d):\033[0m\n' "${#FAILED[@]}"
  printf '  %s\n' "${FAILED[@]}"
  exit 1
fi
log "All ${#MODULES[@]} modules refreshed successfully"
