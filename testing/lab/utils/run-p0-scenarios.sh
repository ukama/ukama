#!/usr/bin/env bash
#!/bin/sh
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

# Run P0 scenario categories sequentially and aggregate their reports.
#
# The caller must export the Ukama authentication/environment variables.
# Credentials are intentionally not stored in this script.
#
# Usage:
#   ./scripts/run-p0-scenarios.sh
#   ./scripts/run-p0-scenarios.sh data-package usage
#
# Optional factory replenishment before the batch:
#   ULAB_FACTORY_NODE_COUNT=100 ./scripts/run-p0-scenarios.sh data-package

set -uo pipefail

usage() {
    printf '%s\n' \
        "usage: $0 [data-package|usage|billing|console-kpi ...]" \
        "" \
        "Environment overrides:" \
        "  UKAMA_REPO                 Ukama repository root" \
        "  UKAMA_LAB_BFF              BFF GraphQL URL" \
        "  UKAMA_LAB_WAREHOUSE_URL    Warehouse API URL" \
        "  UKAMA_LAB_FACTORY_URL      Factory API URL used by ukama-lab" \
        "  ULAB_FACTORY_NODE_COUNT    Node sets to generate before running (default: 0)" \
        "  ULAB_FACTORY_SEED_URL      Factory URL used to generate node sets" \
        "  P0_RUNS_DIR                Parent directory for batch results" \
        "  SCENARIO_ROOT              P0 scenario root" \
        "  LAB_BIN                    ukama-lab executable"
}

if [[ "${1:-}" == "-h" || "${1:-}" == "--help" ]]; then
    usage
    exit 0
fi

require_env() {
    local name="$1"

    if [[ -z "${!name:-}" ]]; then
        printf 'error: required environment variable %s is not set\n' "$name" >&2
        return 1
    fi
}

require_env UKAMA_IDENTIFIER || exit 2
require_env UKAMA_PASSWORD   || exit 2
require_env PAUTH_URL        || exit 2
require_env BFF_BASE_URL     || exit 2

LAB_BIN="${LAB_BIN:-./bin/ukama-lab}"
SCENARIO_ROOT="${SCENARIO_ROOT:-scenarios/p0}"
UKAMA_REPO="${UKAMA_REPO:-/home/kashif/work/ukama/repos/ukama/}"
BFF_GRAPHQL_URL="${UKAMA_LAB_BFF:-${BFF_BASE_URL%/}/gateway/graphql}"
WAREHOUSE_URL="${UKAMA_LAB_WAREHOUSE_URL:-http://warehouse-ukama.udev.ukama.com}"
FACTORY_URL_FOR_LAB="${UKAMA_LAB_FACTORY_URL:-http://factory-ukama.udev.ukama.com}"
FACTORY_SEED_URL="${ULAB_FACTORY_SEED_URL:-https://factory-ukama.udev.ukama.com}"
SIM_TYPE="${UKAMA_LAB_SIM_TYPE:-ukama_data}"
P0_RUNS_DIR="${P0_RUNS_DIR:-runs/p0-batches}"
FACTORY_NODE_COUNT="${ULAB_FACTORY_NODE_COUNT:-0}"
BATCH_STAMP="$(date -u +%Y%m%dT%H%M%SZ)"
BATCH_DIR="${P0_RUNS_DIR%/}/${BATCH_STAMP}"
RUNS_DIR="${BATCH_DIR}/runs"
LOGS_DIR="${BATCH_DIR}/logs"
SUMMARY_TSV="${BATCH_DIR}/scenarios.tsv"

if [[ ! -x "$LAB_BIN" ]]; then
    printf 'error: ukama-lab is not executable: %s\n' "$LAB_BIN" >&2
    exit 2
fi
if [[ ! -d "$SCENARIO_ROOT" ]]; then
    printf 'error: scenario root does not exist: %s\n' "$SCENARIO_ROOT" >&2
    exit 2
fi
if [[ ! -d "$UKAMA_REPO" ]]; then
    printf 'error: Ukama repository does not exist: %s\n' "$UKAMA_REPO" >&2
    exit 2
fi
if [[ ! "$FACTORY_NODE_COUNT" =~ ^[0-9]+$ ]]; then
    printf 'error: ULAB_FACTORY_NODE_COUNT must be a non-negative integer\n' >&2
    exit 2
fi

if (($#)); then
    categories=("$@")
else
    categories=(data-package usage billing console-kpi)
fi

for category in "${categories[@]}"; do
    if [[ ! "$category" =~ ^[a-z0-9-]+$ ||
          ! -d "${SCENARIO_ROOT%/}/${category}" ]]; then
        printf 'error: unknown or missing P0 category: %s\n' "$category" >&2
        exit 2
    fi
done

mkdir -p "$RUNS_DIR" "$LOGS_DIR"
printf 'category\tscenario\trun_id\texit_code\treport\tlog\n' >"$SUMMARY_TSV"

if ((FACTORY_NODE_COUNT > 0)); then
    FACTORY_NODE_SCRIPT="${FACTORY_NODE_SCRIPT:-./generate-factory-nodes.sh}"
    if [[ ! -x "$FACTORY_NODE_SCRIPT" &&
          -x ./scripts/generate-factory-nodes.sh ]]; then
        FACTORY_NODE_SCRIPT=./scripts/generate-factory-nodes.sh
    fi
    if [[ ! -x "$FACTORY_NODE_SCRIPT" ]]; then
        printf 'error: factory node generator is not executable\n' >&2
        exit 2
    fi

    printf '\n== Factory: generate %s node sets ==\n' "$FACTORY_NODE_COUNT"
    if ! PROVISION=false \
         ALLOCATE=false \
         FACTORY_URL="$FACTORY_SEED_URL" \
         OUT="${BATCH_DIR}/generated-factory-nodes.csv" \
         "$FACTORY_NODE_SCRIPT" "$FACTORY_NODE_COUNT"; then
        printf 'error: factory node generation failed\n' >&2
        exit 1
    fi
fi

common_args=(
    --sim-type "$SIM_TYPE"
    --warehouse-url "$WAREHOUSE_URL"
    --factory-url "$FACTORY_URL_FOR_LAB"
    --bff "$BFF_GRAPHQL_URL"
    --out "$RUNS_DIR"
    --verbose
    --repo "$UKAMA_REPO"
)

if [[ -n "${UKAMA_LAB_ASR_URL:-}" ]]; then
    common_args+=(--asr-url "$UKAMA_LAB_ASR_URL")
fi

for category in "${categories[@]}"; do
    category_dir="${SCENARIO_ROOT%/}/${category}"
    mapfile -d '' scenarios < <(
        find "$category_dir" -type f \
            \( -name '*.yaml' -o -name '*.yml' \) -print0 | sort -z
    )

    printf '\n============================================================\n'
    printf 'Category: %s (%s scenarios)\n' "$category" "${#scenarios[@]}"
    printf '============================================================\n'

    if ((${#scenarios[@]} == 0)); then
        printf 'warning: no scenarios found in %s\n' "$category_dir" >&2
        continue
    fi

    for scenario in "${scenarios[@]}"; do
        relative="${scenario#${SCENARIO_ROOT%/}/}"
        slug="${relative%.yaml}"
        slug="${slug%.yml}"
        slug="${slug//\//-}"
        slug="${slug//[^a-zA-Z0-9_.-]/-}"
        run_id="p0-${BATCH_STAMP}-${slug}"
        report_path="${RUNS_DIR}/${run_id}/report.json"
        log_path="${LOGS_DIR}/${slug}.log"

        printf '\n-- Running: %s --\n' "$relative"
        "$LAB_BIN" validate "$scenario" \
            "${common_args[@]}" \
            --run-id "$run_id" 2>&1 | tee "$log_path"
        rc=${PIPESTATUS[0]}

        printf '%s\t%s\t%s\t%s\t%s\t%s\n' \
            "$category" "$relative" "$run_id" "$rc" \
            "$report_path" "$log_path" >>"$SUMMARY_TSV"
    done
done

python3 - "$SUMMARY_TSV" "$BATCH_DIR" <<'PY'
import csv
import json
import sys
from collections import Counter
from pathlib import Path

summary_path = Path(sys.argv[1])
batch_dir = Path(sys.argv[2])
results = []

with summary_path.open(newline="", encoding="utf-8") as stream:
    for row in csv.DictReader(stream, delimiter="\t"):
        exit_code = int(row["exit_code"])
        report_path = Path(row["report"])
        report = None
        report_error = ""

        if report_path.is_file():
            try:
                report = json.loads(report_path.read_text(encoding="utf-8"))
            except (OSError, json.JSONDecodeError) as exc:
                report_error = str(exc)

        scenario_status = report.get("status", "") if report else ""
        if scenario_status in {"wip", "skip"} and exit_code == 0:
            outcome = "SKIP"
        elif exit_code == 0 and report and report.get("passed") is True:
            outcome = "PASS"
        else:
            outcome = "FAIL"

        results.append({
            "category": row["category"],
            "scenario": row["scenario"],
            "run_id": row["run_id"],
            "outcome": outcome,
            "scenario_status": scenario_status,
            "exit_code": exit_code,
            "duration_sec": report.get("duration_sec", 0) if report else 0,
            "events": report.get("events", {}) if report else {},
            "checks": report.get("checks", {}) if report else {},
            "cleanup": report.get("cleanup", "unknown") if report else "unknown",
            "report": row["report"],
            "log": row["log"],
            "report_error": report_error,
        })

counts = Counter(item["outcome"] for item in results)
payload = {
    "total": len(results),
    "passed": counts["PASS"],
    "failed": counts["FAIL"],
    "skipped": counts["SKIP"],
    "duration_sec": sum(item["duration_sec"] for item in results),
    "results": results,
}

json_path = batch_dir / "batch-report.json"
text_path = batch_dir / "batch-report.txt"
json_path.write_text(json.dumps(payload, indent=2) + "\n", encoding="utf-8")

lines = [
    "P0 scenario batch report",
    "",
    f"total={payload['total']} pass={payload['passed']} "
    f"fail={payload['failed']} skip={payload['skipped']} "
    f"duration_sec={payload['duration_sec']}",
    "",
]
for item in results:
    lines.append(
        f"{item['outcome']:<4}  {item['category']:<14}  "
        f"{item['scenario']}  ({item['duration_sec']}s)"
    )
text_path.write_text("\n".join(lines) + "\n", encoding="utf-8")

print("\n============================================================")
print("P0 batch summary")
print("============================================================")
print(f"{'RESULT':<7} {'CATEGORY':<15} SCENARIO")
for item in results:
    print(f"{item['outcome']:<7} {item['category']:<15} {item['scenario']}")
print("------------------------------------------------------------")
print(
    f"total={payload['total']} pass={payload['passed']} "
    f"fail={payload['failed']} skip={payload['skipped']} "
    f"duration_sec={payload['duration_sec']}"
)
print(f"text report: {text_path}")
print(f"json report: {json_path}")

sys.exit(1 if payload["failed"] else 0)
PY
