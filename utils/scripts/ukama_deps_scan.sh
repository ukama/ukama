#!/usr/bin/env bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

set -euo pipefail

# ------------------------------------------------------------
# Ukama deps scanner:
#  - scans ELF binaries + shared objects and runs ldd
#  - reads .gitmodules to learn vendored submodules
#  - classifies libs: submodule-built vs apt-installed (best-effort mapping)
#  - infers build tools needed by inspecting submodule folders
#
# Output:
#  1) tools needed
#  2) libraries needed:
#     2.1) will be installed via build (submodules)
#     2.2) needs to be installed via apt-get
#
# Usage:
#   ./ukama_deps_scan.sh --repo <repo_root> --gitmodules <path_to_.gitmodules> --scan <path_to_tree> [--exclude <dir>]...
# ------------------------------------------------------------

die() { echo "ERROR: $*" >&2; exit 1; }
have() { command -v "$1" >/dev/null 2>&1; }
log() { echo "[$(date +'%H:%M:%S')] $*"; }

REPO_ROOT=""
GITMODULES=""
SCAN_ROOT=""
EXCLUDE_DIRS=()

usage() {
    cat <<EOF
Usage:
  $0 --repo <repo_root> --gitmodules <path_to_.gitmodules> --scan <path_to_tree> [--exclude <dir>]...

Args:
  --repo <path>        Repo root to locate submodule directories.
  --gitmodules <path>  Path to .gitmodules file.
  --scan <path>        Root to scan for ELF executables/shared objects.
  --exclude <path>     Directory to skip (can repeat).

Notes:
  - "submodule-built" classification is best-effort via SUBMOD_LIB_PATTERNS mapping.
  - Extend SUBMOD_LIB_PATTERNS if you add new vendored libs/sonames.
EOF
}

while [[ $# -gt 0 ]]; do
    case "$1" in
        --repo)       REPO_ROOT="${2:-}"; shift 2 ;;
        --gitmodules) GITMODULES="${2:-}"; shift 2 ;;
        --scan)       SCAN_ROOT="${2:-}"; shift 2 ;;
        --exclude)
            [[ $# -ge 2 ]] || die "--exclude requires a directory path"
            EXCLUDE_DIRS+=("$2")
            shift 2
            ;;
        -h|--help) usage; exit 0 ;;
        *) die "Unknown arg: $1 (try --help)" ;;
    esac
done

[[ -n "$REPO_ROOT" ]]  || die "--repo is required"
[[ -n "$GITMODULES" ]] || die "--gitmodules is required"
[[ -n "$SCAN_ROOT" ]]  || die "--scan is required"
[[ -d "$REPO_ROOT" ]]  || die "repo root not a dir: $REPO_ROOT"
[[ -f "$GITMODULES" ]] || die "gitmodules not a file: $GITMODULES"
[[ -d "$SCAN_ROOT" ]]  || die "scan root not a dir: $SCAN_ROOT"

have file || die "'file' not found. Install: sudo apt-get install -y file"
have ldd  || die "'ldd' not found (Ubuntu: sudo apt-get install -y libc-bin)"
have awk  || die "'awk' not found?"
have sort || die "'sort' not found?"

TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT

BIN_LIST="$TMP/binaries.txt"
LDD_RAW="$TMP/ldd_raw.txt"
ALL_LIBS="$TMP/all_libs.txt"
MISSING_LIBS="$TMP/missing_libs.txt"

SUBMODULES_TXT="$TMP/submodules.txt"
SUBMODULE_DIRS="$TMP/submodule_dirs.txt"
TOOLS_TXT="$TMP/tools.txt"

SUBMOD_LIBS="$TMP/submodule_libs.txt"
APT_LIBS="$TMP/apt_libs.txt"

# ------------------------------------------------------------
# 1) Parse .gitmodules -> list of submodule paths + short names
# ------------------------------------------------------------
awk '$1=="path" && $2=="=" { print $3 }' "$GITMODULES" | sort -u > "$SUBMODULE_DIRS"
awk -F/ '{ print $NF }' "$SUBMODULE_DIRS" | sort -u > "$SUBMODULES_TXT"

SUBMOD_COUNT=$(wc -l < "$SUBMODULE_DIRS" | tr -d ' ')
log "Found $SUBMOD_COUNT submodules in .gitmodules"

# ------------------------------------------------------------
# 2) Find ELF executables/shared objects under SCAN_ROOT (with excludes)
# ------------------------------------------------------------
log "Scanning ELF binaries under: $SCAN_ROOT"
log "Excluding paths: ${EXCLUDE_DIRS[*]:-(none)}"

FIND_ARGS=()
for ex in "${EXCLUDE_DIRS[@]}"; do
    if [[ "$ex" = /* ]]; then
        FIND_ARGS+=( -path "$ex" -prune -o )
    else
        FIND_ARGS+=( -path "$SCAN_ROOT/$ex" -prune -o )
    fi
done

# shellcheck disable=SC2068
find "$SCAN_ROOT" \
     ${FIND_ARGS[@]} \
     -type f -print0 \
    | xargs -0 file \
    | grep -E 'ELF .* (executable|shared object)' \
    | cut -d: -f1 \
    | sort -u > "$BIN_LIST"

BIN_COUNT=$(wc -l < "$BIN_LIST" | tr -d ' ')
log "Found $BIN_COUNT ELF files"
[[ "$BIN_COUNT" -gt 0 ]] || die "No ELF binaries found under --scan path"

# ------------------------------------------------------------
# 3) Run ldd for each binary
# ------------------------------------------------------------
log "Running ldd..."
: > "$LDD_RAW"
while read -r bin; do
    echo "### $bin" >> "$LDD_RAW"
    # ldd can exit non-zero for static or unusual binaries; we still want output
    ldd "$bin" 2>&1 >> "$LDD_RAW" || true
    echo >> "$LDD_RAW"
done < "$BIN_LIST"

# ------------------------------------------------------------
# 4) Extract ALL libs + missing libs (make greps non-fatal)
# ------------------------------------------------------------
: > "$ALL_LIBS"
: > "$MISSING_LIBS"

# ALL libs (sonames) mentioned by ldd.
# Make grep non-fatal if no matches (static binaries, etc.)
grep -E '=>|ld-linux|linux-vdso' "$LDD_RAW" 2>/dev/null \
    | awk '{print $1}' \
    | sed 's/://g' \
    | sort -u > "$ALL_LIBS" || true

# Missing libs (grep returns 1 when no matches; do not treat as failure)
grep "not found" "$LDD_RAW" 2>/dev/null \
    | awk '{print $1}' \
    | sort -u > "$MISSING_LIBS" || true

# ------------------------------------------------------------
# 5) Classify libs: submodule build vs apt-get install
# ------------------------------------------------------------
# Best-effort mapping of submodule -> expected soname patterns.
declare -A SUBMOD_LIB_PATTERNS=(
    # babelouest stack
    ["orcania"]="^liborcania\\.so"
    ["yder"]="^libyder\\.so"
    ["ulfius"]="^libulfius\\.so"

    # data formats
    ["tomlc"]="^libtoml(c99)?\\.so|^libtoml\\.so"

    # json
    ["jansson"]="^libjansson\\.so"

    # rabbitmq-c
    ["amqp"]="^librabbitmq(\\-c)?\\.so"

    # protobuf-c
    ["protobuf-c"]="^libprotobuf-c\\.so"

    # prometheus client (ukama fork - adjust if soname differs)
    ["prometheus-client"]="^libprom\\.so|^libprometheus.*\\.so"

    # forks that might also exist system-wide
    ["libcap"]="^libcap\\.so"
    ["libuuid"]="^libuuid\\.so"

    # Unity is test framework; typically no runtime soname
    ["Unity"]="$^"
    ["firmware"]="$^"
    ["kernel"]="$^"
    ["busybox"]="$^"
)

declare -A HAVE_SUBMOD=()
while read -r sm; do
    [[ -n "$sm" ]] && HAVE_SUBMOD["$sm"]=1
done < "$SUBMODULES_TXT"

is_submodule_lib() {
    local lib="$1"
    local sm
    for sm in "${!SUBMOD_LIB_PATTERNS[@]}"; do
        [[ -n "${HAVE_SUBMOD[$sm]:-}" ]] || continue
        if [[ "$lib" =~ ${SUBMOD_LIB_PATTERNS[$sm]} ]]; then
            return 0
        fi
    done
  return 1
}

: > "$SUBMOD_LIBS"
: > "$APT_LIBS"

if [[ -s "$ALL_LIBS" ]]; then
    while read -r lib; do
        [[ -n "$lib" ]] || continue
        [[ "$lib" == "linux-vdso.so.1" ]] && continue

        if is_submodule_lib "$lib"; then
            echo "$lib" >> "$SUBMOD_LIBS"
        else
            echo "$lib" >> "$APT_LIBS"
        fi
    done < "$ALL_LIBS"
fi

sort -u -o "$SUBMOD_LIBS" "$SUBMOD_LIBS" 2>/dev/null || true
sort -u -o "$APT_LIBS" "$APT_LIBS" 2>/dev/null || true

# ------------------------------------------------------------
# 6) Infer tools needed for proper build of these submodules
# ------------------------------------------------------------
add_tool() { echo "$1" >> "$TOOLS_TXT"; }
: > "$TOOLS_TXT"

# Baseline tools
add_tool "git (git submodule update --init --recursive)"
add_tool "build-essential (gcc/g++/make)"
add_tool "pkg-config"

while read -r relpath; do
    [[ -n "$relpath" ]] || continue
    smdir="$REPO_ROOT/$relpath"
    [[ -d "$smdir" ]] || continue

  # Autotools
  if [[ -f "$smdir/configure.ac" || -f "$smdir/configure.in" || -f "$smdir/autogen.sh" || -f "$smdir/bootstrap" ]]; then
      add_tool "autoconf"
      add_tool "automake"
      add_tool "libtool"
      add_tool "m4"
  fi

  # CMake
  if [[ -f "$smdir/CMakeLists.txt" ]]; then
      add_tool "cmake"
      add_tool "ninja-build (optional)"
  fi

  # Meson
  if [[ -f "$smdir/meson.build" ]]; then
      add_tool "meson"
      add_tool "ninja-build"
  fi

  # Go
  if [[ -f "$smdir/go.mod" ]]; then
      add_tool "golang"
  fi

  # Node
  if [[ -f "$smdir/package.json" ]]; then
      add_tool "nodejs"
      add_tool "npm (or yarn/pnpm depending on repo)"
  fi

  # Python
  if [[ -f "$smdir/pyproject.toml" || -f "$smdir/setup.py" || -f "$smdir/requirements.txt" ]]; then
      add_tool "python3"
      add_tool "python3-pip"
      add_tool "python3-venv"
  fi

  # Rust
  if [[ -f "$smdir/Cargo.toml" ]]; then
      add_tool "rustc + cargo"
  fi
done < "$SUBMODULE_DIRS"

sort -u -o "$TOOLS_TXT" "$TOOLS_TXT"

# ------------------------------------------------------------
# 7) Print final output
# ------------------------------------------------------------
echo
echo "=================================================="
echo "1. tools needed"
echo "=================================================="
cat "$TOOLS_TXT"

echo
echo "=================================================="
echo "2. libraries needed"
echo "=================================================="
echo "  2.1) will be installed via build (those in the submodule)"
if [[ -s "$SUBMOD_LIBS" ]]; then
    sed 's/^/    - /' "$SUBMOD_LIBS"
else
    echo "    - (none detected via mapping)"
fi

echo
echo "  2.2) those which needs to be installed via apt-get"
if [[ -s "$APT_LIBS" ]]; then
    sed 's/^/    - /' "$APT_LIBS"
else
    echo "    - (none)"
fi

echo
echo "--------------------------------------------------"
echo "Missing libs (if any):"
if [[ -s "$MISSING_LIBS" ]]; then
    sed 's/^/  - /' "$MISSING_LIBS"
else
    echo "  - None"
fi

echo
echo "Raw ldd output stored at: $LDD_RAW"
