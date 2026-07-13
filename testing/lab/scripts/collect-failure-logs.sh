#!/bin/sh
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

set -u

if [ "$#" -lt 1 ] || [ "$#" -gt 2 ]; then
    echo "usage: $0 <run-dir> [out-dir-name]" >&2
    exit 2
fi

RUN_DIR="$1"
OUT_NAME="${2:-failure-diagnostics}"

case "$OUT_NAME" in
    ""|*/*|*..*|*[!A-Za-z0-9_.-]*)
        echo "invalid diagnostics out-dir-name: $OUT_NAME" >&2
        exit 2
        ;;
esac

OUT_DIR="$RUN_DIR/$OUT_NAME"
SUMMARY="$OUT_DIR/summary.txt"
LEVELS="${ULAB_FAILURE_LOG_LEVELS:-error,critical}"
ATTEMPTED=0
COLLECTED=0
FAILURES=0

rm -rf "$OUT_DIR"
mkdir -p "$OUT_DIR"
: > "$SUMMARY"

log() {
    printf "%s\n" "$*" | tee -a "$SUMMARY"
}

warn() {
    FAILURES=$((FAILURES + 1))
    printf "WARN %s\n" "$*" | tee -a "$SUMMARY" >&2
}

safe_name() {
    printf "%s" "$1" | tr -c 'A-Za-z0-9_.-' '-'
}

json_escape() {
    printf "%s" "$1" | sed 's/\\/\\\\/g; s/"/\\"/g'
}

write_manifest() {
    manifest="$OUT_DIR/manifest.json"
    tmp="$manifest.tmp"

    {
        printf '{\n'
        printf '  "schema": "ukama.lab.failure-logs.v1",\n'
        printf '  "run_dir": "%s",\n' "$(json_escape "$RUN_DIR")"
        printf '  "diagnostics_dir": "%s",\n' \
            "$(json_escape "$OUT_NAME")"
        printf '  "command": "/sbin/ukama-log",\n'
        printf '  "boot": "current",\n'
        printf '  "levels": "%s",\n' "$(json_escape "$LEVELS")"
        printf '  "format": "jsonl",\n'
        printf '  "collected_at": %s,\n' "$(date +%s)"
        printf '  "nodes_attempted": %s,\n' "$ATTEMPTED"
        printf '  "nodes_collected": %s,\n' "$COLLECTED"
        printf '  "failures": %s\n' "$FAILURES"
        printf '}\n'
    } > "$tmp"
    mv "$tmp" "$manifest"
}

write_node_env() {
    out="$1"

    {
        echo "NODE_FILE=$NODE_FILE"
        echo "LOGICAL_NODE_ID=${LOGICAL_NODE_ID:-}"
        echo "FACTORY_NODE_ID=${FACTORY_NODE_ID:-}"
        echo "NODE_KIND=${NODE_KIND:-${NODE_TYPE:-}}"
        echo "SITE_REF=${SITE_REF:-}"
        echo "CONTAINER_NAME=${CONTAINER_NAME:-}"
    } > "$out"
}

collect_node() {
    NODE_FILE="$1"
    LOGICAL_NODE_ID=""
    FACTORY_NODE_ID=""
    NODE_KIND=""
    NODE_TYPE=""
    SITE_REF=""
    CONTAINER_NAME=""
    # shellcheck disable=SC1090
    . "$NODE_FILE"

    node_key="${LOGICAL_NODE_ID:-${FACTORY_NODE_ID:-${CONTAINER_NAME:-node}}}"
    node_safe="$(safe_name "$node_key")"
    node_dir="$OUT_DIR/$node_safe"
    output="$node_dir/errors.jsonl"
    stderr_file="$node_dir/ukama-log.stderr"
    rc_file="$node_dir/ukama-log.rc"

    ATTEMPTED=$((ATTEMPTED + 1))
    mkdir -p "$node_dir"
    write_node_env "$node_dir/node.env"
    : > "$output"
    : > "$stderr_file"

    if [ -z "${CONTAINER_NAME:-}" ]; then
        printf '1\n' > "$rc_file"
        warn "node=$node_key missing CONTAINER_NAME"
        return 0
    fi

    if ! podman inspect -f '{{.State.Running}}' "$CONTAINER_NAME" \
        2>/dev/null | grep -q '^true$'; then
        printf '1\n' > "$rc_file"
        warn "node=$node_key container=$CONTAINER_NAME not running"
        return 0
    fi

    if ! podman exec "$CONTAINER_NAME" test -x /sbin/ukama-log \
        >/dev/null 2>"$stderr_file"; then
        printf '1\n' > "$rc_file"
        warn "node=$node_key /sbin/ukama-log unavailable"
        return 0
    fi

    if podman exec "$CONTAINER_NAME" /sbin/ukama-log \
        --boot current \
        --level "$LEVELS" \
        --format jsonl \
        > "$output" 2> "$stderr_file"; then
        rc=0
        COLLECTED=$((COLLECTED + 1))
    else
        rc=$?
        warn "node=$node_key ukama-log failed rc=$rc"
    fi

    printf '%s\n' "$rc" > "$rc_file"
    records="$(wc -l < "$output" | tr -d ' ')"
    log "node=$node_key container=$CONTAINER_NAME records=$records rc=$rc"
}

log "failure log collection begin run_dir=$RUN_DIR out=$OUT_DIR"
log "command=/sbin/ukama-log boot=current levels=$LEVELS format=jsonl"

if [ ! -d "$RUN_DIR/runtime-nodes" ]; then
    warn "runtime-nodes directory not found: $RUN_DIR/runtime-nodes"
else
    for node_file in "$RUN_DIR"/runtime-nodes/*.env; do
        [ -f "$node_file" ] || continue
        collect_node "$node_file"
    done
fi

write_manifest
log "failure log collection complete attempted=$ATTEMPTED "\
    "collected=$COLLECTED failures=$FAILURES"

exit 0
