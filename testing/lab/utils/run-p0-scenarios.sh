#!/usr/bin/env bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

# Run P0 scenario categories sequentially and aggregate their reports.
#
# Udev authentication and service values have defaults below. The caller may
# override any of them by exporting a different value before starting.
#
# Usage:
#   ./scripts/run-p0-scenarios.sh
#   ./scripts/run-p0-scenarios.sh data-package usage
#   ./scripts/run-p0-scenarios.sh --factory-nodes auto data-package
#   ./scripts/run-p0-scenarios.sh --factory-nodes 100 usage billing

set -uo pipefail

# Udev defaults. A value already exported by the caller takes precedence.
export UKAMA_IDENTIFIER="${UKAMA_IDENTIFIER:-ukama@ukama.com}"
export UKAMA_PASSWORD="${UKAMA_PASSWORD:-@Pass2026.}"
export PAUTH_URL="${PAUTH_URL:-https://pauth.udev.ukama.com}"
export BFF_BASE_URL="${BFF_BASE_URL:-https://bff.udev.ukama.com}"
export UKAMA_LAB_DUMP_BFF_CURL="${UKAMA_LAB_DUMP_BFF_CURL:-1}"
export ULAB_CDR_WAIT_SEC="${ULAB_CDR_WAIT_SEC:-5}"
export ULAB_UKAMA_AGENT_NODE_GW_URL="${ULAB_UKAMA_AGENT_NODE_GW_URL:-https://ukamaagent-nodegateway-ukama.udev.ukama.com:8080/}"
export ULAB_UKAMA_AGENT_API_GW_URL="${ULAB_UKAMA_AGENT_API_GW_URL:-https://ukamaagent-ukama.udev.ukama.com:8080/}"
export ULAB_CDR_DIAG_STRICT="${ULAB_CDR_DIAG_STRICT:-0}"
export ULAB_OVS_ALLOW_UNMETERED_FALLBACK="${ULAB_OVS_ALLOW_UNMETERED_FALLBACK:-1}"
export ULAB_CDR_DIAG_DISABLE="${ULAB_CDR_DIAG_DISABLE:-1}"

usage() {
    printf '%s\n' \
        "usage: $0 [--factory-nodes auto|N] [data-package|usage|billing|console-kpi ...]" \
        "" \
        "Options:" \
        "  --factory-nodes auto       Ensure enough complete bundles for runnable scenarios" \
        "  --factory-nodes N          Ensure at least N complete unprovisioned bundles" \
        "  -h, --help                 Show this help" \
        "" \
        "Environment overrides:" \
        "  UKAMA_REPO                 Ukama repository root" \
        "  UKAMA_LAB_BFF              BFF GraphQL URL" \
        "  UKAMA_LAB_WAREHOUSE_URL    Warehouse API URL" \
        "  UKAMA_LAB_FACTORY_URL      Factory API URL used by ukama-lab" \
        "  ULAB_FACTORY_NODE_COUNT    Default factory target: auto or N (default: 0/off)" \
        "  ULAB_FACTORY_HEADROOM_PERCENT  Auto-mode safety margin (default: 10)" \
        "  ULAB_FACTORY_SEED_URL      Factory URL used to generate node sets" \
        "  P0_RUNS_DIR                Parent directory for batch results" \
        "  SCENARIO_ROOT              P0 scenario root" \
        "  LAB_BIN                    ukama-lab executable"
}

require_env() {
    local name="$1"

    if [[ -z "${!name:-}" ]]; then
        printf 'error: required environment variable %s is not set\n' "$name" >&2
        return 1
    fi
}

require_env UKAMA_IDENTIFIER || exit 2
require_env UKAMA_PASSWORD || exit 2
require_env PAUTH_URL || exit 2
require_env BFF_BASE_URL || exit 2

LAB_BIN="${LAB_BIN:-./bin/ukama-lab}"
SCENARIO_ROOT="${SCENARIO_ROOT:-scenarios/p0}"
UKAMA_REPO="${UKAMA_REPO:-/home/kashif/work/ukama/repos/ukama/}"
BFF_GRAPHQL_URL="${UKAMA_LAB_BFF:-${BFF_BASE_URL%/}/gateway/graphql}"
WAREHOUSE_URL="${UKAMA_LAB_WAREHOUSE_URL:-http://warehouse-ukama.udev.ukama.com}"
FACTORY_URL_FOR_LAB="${UKAMA_LAB_FACTORY_URL:-http://factory-ukama.udev.ukama.com}"
FACTORY_SEED_URL="${ULAB_FACTORY_SEED_URL:-https://factory-ukama.udev.ukama.com}"
SIM_TYPE="${UKAMA_LAB_SIM_TYPE:-ukama_data}"
P0_RUNS_DIR="${P0_RUNS_DIR:-runs/p0-batches}"
FACTORY_NODE_TARGET="${ULAB_FACTORY_NODE_COUNT:-0}"
FACTORY_HEADROOM_PERCENT="${ULAB_FACTORY_HEADROOM_PERCENT:-10}"
BATCH_STAMP="$(date -u +%Y%m%dT%H%M%SZ)"
BATCH_DIR="${P0_RUNS_DIR%/}/${BATCH_STAMP}"
RUNS_DIR="${BATCH_DIR}/runs"
LOGS_DIR="${BATCH_DIR}/logs"
SUMMARY_TSV="${BATCH_DIR}/scenarios.tsv"

categories=()
while (($#)); do
    case "$1" in
        --factory-nodes)
            if (($# < 2)); then
                printf 'error: --factory-nodes requires auto or a number\n' >&2
                exit 2
            fi
            FACTORY_NODE_TARGET="$2"
            shift 2
            ;;
        --factory-nodes=*)
            FACTORY_NODE_TARGET="${1#*=}"
            shift
            ;;
        -h|--help)
            usage
            exit 0
            ;;
        --)
            shift
            categories+=("$@")
            break
            ;;
        -*)
            printf 'error: unknown option: %s\n' "$1" >&2
            usage >&2
            exit 2
            ;;
        *)
            categories+=("$1")
            shift
            ;;
    esac
done

if ((${#categories[@]} == 0)); then
    categories=(data-package usage billing console-kpi)
fi

if [[ "$FACTORY_NODE_TARGET" != "auto" &&
      ! "$FACTORY_NODE_TARGET" =~ ^[0-9]+$ ]]; then
    printf 'error: --factory-nodes must be auto or a non-negative integer\n' >&2
    exit 2
fi
if [[ ! "$FACTORY_HEADROOM_PERCENT" =~ ^[0-9]+$ ]]; then
    printf 'error: ULAB_FACTORY_HEADROOM_PERCENT must be a non-negative integer\n' >&2
    exit 2
fi

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

for category in "${categories[@]}"; do
    if [[ ! "$category" =~ ^[a-z0-9-]+$ ||
          ! -d "${SCENARIO_ROOT%/}/${category}" ]]; then
        printf 'error: unknown or missing P0 category: %s\n' "$category" >&2
        exit 2
    fi
done

mkdir -p "$RUNS_DIR" "$LOGS_DIR"
printf 'category\tscenario\trun_id\texit_code\treport\tlog\n' >"$SUMMARY_TSV"

selected_scenarios=()
for category in "${categories[@]}"; do
    category_dir="${SCENARIO_ROOT%/}/${category}"
    while IFS= read -r -d '' scenario; do
        selected_scenarios+=("$scenario")
    done < <(
        find "$category_dir" -type f \
            \( -name '*.yaml' -o -name '*.yml' \) -print0 | sort -z
    )
done

factory_scenario_plan() {
    python3 - "$@" <<'PY'
import sys

def scenario_values(path):
    status = "active"
    networks = 0
    sites_per_network = 0
    nodes = {"tower": 0, "amplifier": 0, "controller": 0}
    section = ""
    in_nodes = False

    with open(path, encoding="utf-8") as stream:
        for raw in stream:
            line = raw.split("#", 1)[0].rstrip()
            if not line.strip() or ":" not in line:
                continue

            indent = len(line) - len(line.lstrip(" "))
            key, value = line.strip().split(":", 1)
            key = key.strip()
            value = value.strip().strip('"\'')

            if indent == 0:
                section = key if not value else ""
                in_nodes = False
                if key == "status" and value:
                    status = value
                continue

            if section != "world":
                continue

            if indent == 2:
                in_nodes = key == "nodes_per_site"
                if key == "networks" and value:
                    networks = int(value)
                elif key == "sites_per_network" and value:
                    sites_per_network = int(value)
                continue

            if in_nodes and indent == 4 and key in nodes and value:
                nodes[key] = int(value)

    return status, networks, sites_per_network, nodes


required = 0
runnable = 0
ignored = 0
for scenario in sys.argv[1:]:
    try:
        status, networks, sites, nodes = scenario_values(scenario)
    except (OSError, ValueError) as exc:
        print(f"error: cannot calculate factory nodes for {scenario}: {exc}",
              file=sys.stderr)
        sys.exit(2)

    if status not in {"active", "xfail"}:
        ignored += 1
        continue

    runnable += 1
    bundles_per_site = max(nodes.values())
    required += networks * sites * bundles_per_site

print(required, runnable, ignored)
PY
}

factory_available_bundles() {
    local response
    local rc

    response="$(mktemp)" || return 1
    if ! curl -fsS \
        -H 'accept: application/json' \
        "${FACTORY_SEED_URL%/}/v1/nodefactory/nodes?isProvisioned=false" \
        -o "$response"; then
        rm -f "$response"
        return 1
    fi

    python3 - "$response" <<'PY'
import json
import sys

def find_list(value):
    if isinstance(value, list):
        return value
    if isinstance(value, dict):
        for key in ("nodes", "Nodes", "data", "Data", "items", "Items",
                    "results", "Results"):
            if key in value:
                result = find_list(value[key])
                if result is not None:
                    return result
    return None


def get(record, *keys):
    for key in keys:
        if isinstance(record, dict) and record.get(key) not in (None, ""):
            return record[key]
    return ""


def node_type(value, node_id):
    value = str(value).lower()
    node_id = str(node_id).lower()
    if value in {"tnode", "tower", "tower_node"} or "-tnode-" in node_id:
        return "tnode"
    if value in {"cnode", "controller", "controller_node"} or "-cnode-" in node_id:
        return "cnode"
    if value in {"anode", "amplifier", "access", "amplifier_node"} or "-anode-" in node_id:
        return "anode"
    return value


with open(sys.argv[1], encoding="utf-8") as stream:
    data = json.load(stream)

nodes = find_list(data)
if nodes is None:
    print("error: factory response does not contain a node list", file=sys.stderr)
    sys.exit(2)

by_id = {}
tnodes = []
for record in nodes:
    node_id = str(get(record, "id", "Id", "nodeId", "NodeId", "node_id"))
    if not node_id:
        continue
    kind = node_type(
        get(record, "type", "Type", "nodeType", "NodeType", "node_type"),
        node_id,
    )
    by_id[node_id] = kind
    if kind == "tnode":
        tnodes.append(node_id)

complete = 0
for tower_id in tnodes:
    controller_id = tower_id.replace("-tnode-", "-cnode-")
    amplifier_id = tower_id.replace("-tnode-", "-anode-")
    if by_id.get(controller_id) == "cnode" and by_id.get(amplifier_id) == "anode":
        complete += 1

print(complete)
PY
    rc=$?
    rm -f "$response"
    return "$rc"
}

if [[ "$FACTORY_NODE_TARGET" != "0" ]]; then
    if ! command -v curl >/dev/null 2>&1; then
        printf 'error: curl is required for factory preparation\n' >&2
        exit 2
    fi

    if ! read -r required_bundles runnable_scenarios ignored_scenarios < <(
        factory_scenario_plan "${selected_scenarios[@]}"
    ); then
        printf 'error: unable to calculate the factory requirement\n' >&2
        exit 2
    fi

    if [[ "$FACTORY_NODE_TARGET" == "auto" ]]; then
        if ((required_bundles > 0)); then
            headroom=$(((required_bundles * FACTORY_HEADROOM_PERCENT + 99) / 100))
            if ((headroom < 2)); then
                headroom=2
            fi
            target_bundles=$((required_bundles + headroom))
        else
            headroom=0
            target_bundles=0
        fi
    else
        headroom=0
        target_bundles=$FACTORY_NODE_TARGET
        if ((target_bundles < required_bundles)); then
            target_bundles=$required_bundles
        fi
    fi

    printf '\n== Factory preparation ==\n'
    printf 'selected=%s runnable=%s ignored=%s required=%s headroom=%s target=%s\n' \
        "${#selected_scenarios[@]}" "$runnable_scenarios" \
        "$ignored_scenarios" "$required_bundles" "$headroom" \
        "$target_bundles"

    if ((target_bundles > 0)); then
        if ! available_bundles="$(factory_available_bundles)"; then
            printf 'error: failed to query available factory bundles\n' >&2
            exit 1
        fi
        printf 'available=%s\n' "$available_bundles"

        if ((available_bundles < target_bundles)); then
            deficit=$((target_bundles - available_bundles))
            FACTORY_NODE_SCRIPT="${FACTORY_NODE_SCRIPT:-./generate-factory-nodes.sh}"
            if [[ ! -x "$FACTORY_NODE_SCRIPT" &&
                  -x ./scripts/generate-factory-nodes.sh ]]; then
                FACTORY_NODE_SCRIPT=./scripts/generate-factory-nodes.sh
            fi
            if [[ ! -x "$FACTORY_NODE_SCRIPT" ]]; then
                printf 'error: factory node generator is not executable\n' >&2
                exit 2
            fi

            printf 'generating=%s complete node bundles\n' "$deficit"
            if ! PROVISION=false \
                 ALLOCATE=false \
                 FACTORY_URL="$FACTORY_SEED_URL" \
                 OUT="${BATCH_DIR}/generated-factory-nodes.csv" \
                 "$FACTORY_NODE_SCRIPT" "$deficit"; then
                printf 'error: factory node generation failed\n' >&2
                exit 1
            fi

            if ! available_bundles="$(factory_available_bundles)"; then
                printf 'error: failed to verify factory bundles\n' >&2
                exit 1
            fi
            if ((available_bundles < target_bundles)); then
                printf 'error: factory has %s bundles after generation; need %s\n' \
                    "$available_bundles" "$target_bundles" >&2
                exit 1
            fi
            printf 'available_after_generation=%s\n' "$available_bundles"
        else
            printf 'factory pool already satisfies the target\n'
        fi
    else
        printf 'no runnable scenarios require factory bundles\n'
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
