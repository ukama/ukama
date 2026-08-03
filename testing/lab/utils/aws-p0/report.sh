#!/usr/bin/env bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

# Merge worker reports already downloaded to the local batch directory.

set -Eeuo pipefail

SCRIPT_DIR="$(CDPATH= cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=lib.sh
. "$SCRIPT_DIR/lib.sh"

BATCH_DIR="${1:-}"
[[ -n "$BATCH_DIR" ]] || p0_die "usage: $0 LOCAL_BATCH_DIRECTORY"
[[ -d "$BATCH_DIR" ]] || p0_die "batch directory not found: $BATCH_DIR"

p0_require_cmd jq

ROWS_JSONL="$BATCH_DIR/.results.jsonl"
INFRA_TSV="$BATCH_DIR/infrastructure-failures.tsv"
COMBINED_JSON="$BATCH_DIR/batch-report.json"
COMBINED_TSV="$BATCH_DIR/combined.tsv"
SUMMARY_TXT="$BATCH_DIR/summary.txt"
FAILED_TXT="$BATCH_DIR/failed.txt"

: >"$ROWS_JSONL"
printf 'worker\treason\n' >"$INFRA_TSV"

worker_dirs=()
if [[ -d "$BATCH_DIR/workers" ]]; then
    mapfile -t worker_dirs < <(find "$BATCH_DIR/workers" \
        -mindepth 1 -maxdepth 1 -type d -name 'worker-*' | sort)
fi

for worker_dir in "${worker_dirs[@]}"; do
    worker="$(basename "$worker_dir")"
    report="$(find "$worker_dir" -type f -name batch-report.json | head -1)"

    if [[ -n "$report" && -f "$report" ]]; then
        jq -c --arg worker "$worker" \
            '.results[] | . + {worker:$worker}' "$report" \
            >>"$ROWS_JSONL"
    else
        reason='missing batch-report.json'
        if [[ -f "$worker_dir/$worker.failed" ]]; then
            reason="$(tr '\n' ' ' <"$worker_dir/$worker.failed" | sed 's/[[:space:]]*$//')"
        fi
        printf '%s\t%s\n' "$worker" "$reason" >>"$INFRA_TSV"
    fi
done

if [[ -s "$ROWS_JSONL" ]]; then
    jq -s '
      def count_result($name): map(select(.outcome == $name)) | length;
      {
        total:length,
        passed:count_result("PASS"),
        failed:count_result("FAIL"),
        skipped:count_result("SKIP"),
        duration_sec:(map(.duration_sec // 0) | add // 0),
        results:.
      }
    ' "$ROWS_JSONL" >"$COMBINED_JSON"
else
    printf '{"total":0,"passed":0,"failed":0,"skipped":0,"duration_sec":0,"results":[]}\n' \
        >"$COMBINED_JSON"
fi

{
    printf 'result\tworker\tcategory\tscenario\tduration_sec\treport\tlog\n'
    jq -r '.results[] | [
        .outcome,
        .worker,
        .category,
        .scenario,
        (.duration_sec // 0),
        (.report // ""),
        (.log // "")
      ] | @tsv' "$COMBINED_JSON"
} >"$COMBINED_TSV"

jq -r '.results[] | select(.outcome == "FAIL") |
    [.worker, .scenario, (.log // "")] | @tsv' \
    "$COMBINED_JSON" >"$FAILED_TXT"

infra_failed=$(( $(wc -l <"$INFRA_TSV") - 1 ))
total="$(jq -r '.total' "$COMBINED_JSON")"
passed="$(jq -r '.passed' "$COMBINED_JSON")"
failed="$(jq -r '.failed' "$COMBINED_JSON")"
skipped="$(jq -r '.skipped' "$COMBINED_JSON")"
duration="$(jq -r '.duration_sec' "$COMBINED_JSON")"

{
    printf 'Ukama distributed P0 report\n\n'
    printf 'total=%s pass=%s fail=%s skip=%s infra_fail=%s duration_sec=%s\n' \
        "$total" "$passed" "$failed" "$skipped" "$infra_failed" "$duration"
    printf '\nFailed scenarios:\n'
    if [[ -s "$FAILED_TXT" ]]; then
        awk -F '\t' '{printf "  %-10s %s\n", $1, $2}' "$FAILED_TXT"
    else
        printf '  none\n'
    fi
    printf '\nInfrastructure failures:\n'
    if ((infra_failed > 0)); then
        tail -n +2 "$INFRA_TSV" | awk -F '\t' '{printf "  %-10s %s\n", $1, $2}'
    else
        printf '  none\n'
    fi
} >"$SUMMARY_TXT"

printf '\n============================================================\n'
printf 'Distributed P0 summary\n'
printf '============================================================\n'
printf '%-7s %-10s %-18s %s\n' RESULT WORKER CATEGORY SCENARIO
jq -r '.results[] | [
    .outcome,
    .worker,
    .category,
    .scenario
  ] | @tsv' "$COMBINED_JSON" |
while IFS=$'\t' read -r outcome worker category scenario; do
    printf '%-7s %-10s %-18s %s\n' \
        "$outcome" "$worker" "$category" "$scenario"
done
printf '%s\n' '------------------------------------------------------------'
printf 'total=%s pass=%s fail=%s skip=%s infra_fail=%s duration_sec=%s\n' \
    "$total" "$passed" "$failed" "$skipped" "$infra_failed" "$duration"
printf 'summary: %s\n' "$SUMMARY_TXT"
printf 'json:    %s\n' "$COMBINED_JSON"
printf 'tsv:     %s\n' "$COMBINED_TSV"

if ((failed > 0 || infra_failed > 0)); then
    exit 1
fi
