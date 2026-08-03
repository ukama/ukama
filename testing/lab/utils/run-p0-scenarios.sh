#!/usr/bin/env bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

# Run P0 scenarios sequentially and aggregate their reports.
#
# With no category arguments, every top-level directory under scenarios/p0
# is discovered and run recursively. Scenario status is still respected by
# ukama-lab: active/xfail scenarios run, while wip/skip scenarios are reported
# as skipped without provisioning a world.
#
# Usage:
#   ./utils/run-p0-scenarios.sh
#   ./utils/run-p0-scenarios.sh software-update console site
#   ./utils/run-p0-scenarios.sh --factory-nodes auto
#   ./utils/run-p0-scenarios.sh --list

set -uo pipefail

# Credentials must be supplied by the caller. Non-secret service URLs retain
# useful udev defaults.
export UKAMA_IDENTIFIER="${UKAMA_IDENTIFIER:-}"
export UKAMA_PASSWORD="${UKAMA_PASSWORD:-}"
export PAUTH_URL="${PAUTH_URL:-https://pauth.udev.ukama.com}"
export BFF_BASE_URL="${BFF_BASE_URL:-https://bff.udev.ukama.com}"
export UKAMA_LAB_DUMP_BFF_CURL="${UKAMA_LAB_DUMP_BFF_CURL:-1}"
export ULAB_CDR_WAIT_SEC="${ULAB_CDR_WAIT_SEC:-60}"
export ULAB_UKAMA_AGENT_NODE_GW_URL="${ULAB_UKAMA_AGENT_NODE_GW_URL:-https://ukamaagent-nodegateway-ukama.udev.ukama.com:8080/}"
export ULAB_UKAMA_AGENT_API_GW_URL="${ULAB_UKAMA_AGENT_API_GW_URL:-https://ukamaagent-ukama.udev.ukama.com:8080/}"
export ULAB_CDR_DIAG_STRICT="${ULAB_CDR_DIAG_STRICT:-0}"
export ULAB_OVS_ALLOW_UNMETERED_FALLBACK="${ULAB_OVS_ALLOW_UNMETERED_FALLBACK:-1}"
export ULAB_CDR_DIAG_DISABLE="${ULAB_CDR_DIAG_DISABLE:-1}"

usage() {
    cat <<EOF_USAGE
usage: $0 [options] [category ...]

With no categories, every P0 category below SCENARIO_ROOT is run.
Categories are top-level directories such as billing, console,
console-kpi, data-package, node, sim, site, software-update, ue and usage.

Options:
  --factory-nodes auto       Ensure enough complete bundles for runnable scenarios
  --factory-nodes N          Ensure at least N complete unprovisioned bundles
  --scenario-list FILE       Run exact scenario paths listed in FILE
  --batch-id ID              Use ID instead of a generated batch timestamp
  --prepare-only             Prepare factory nodes, then exit
  --list                     List selected scenarios and exit
  --fail-fast                Stop after the first failed scenario
  -h, --help                 Show this help

Environment overrides:
  UKAMA_REPO                 Ukama repository root
  UKAMA_LAB_BFF              BFF GraphQL URL
  UKAMA_LAB_WAREHOUSE_URL    Warehouse API URL
  UKAMA_LAB_FACTORY_URL      Factory API URL used by ukama-lab
  UKAMA_LAB_SIM_TYPE         SIM type (default: ukama_data)
  ULAB_FACTORY_NODE_COUNT    Default factory target: auto or N (default: 0/off)
  ULAB_FACTORY_HEADROOM_PERCENT
                             Auto-mode safety margin (default: 10)
  ULAB_FACTORY_SEED_URL      Factory URL used to generate node sets
  P0_RUNS_DIR                Parent directory for batch results
  SCENARIO_ROOT              P0 scenario root (default: scenarios/p0)
  LAB_BIN                    ukama-lab executable (default: ./bin/ukama-lab)
  P0_STATUS_FILE             Optional live worker status TSV
EOF_USAGE
}

require_env() {
    local name="$1"

    if [[ -z "${!name:-}" ]]; then
        printf 'error: required environment variable %s is not set\n' \
            "$name" >&2
        return 1
    fi
}

LAB_BIN="${LAB_BIN:-./bin/ukama-lab}"
SCENARIO_ROOT="${SCENARIO_ROOT:-scenarios/p0}"
UKAMA_REPO="${UKAMA_REPO:-}"
BFF_GRAPHQL_URL="${UKAMA_LAB_BFF:-${BFF_BASE_URL%/}/gateway/graphql}"
WAREHOUSE_URL="${UKAMA_LAB_WAREHOUSE_URL:-http://warehouse-ukama.udev.ukama.com}"
FACTORY_URL_FOR_LAB="${UKAMA_LAB_FACTORY_URL:-http://factory-ukama.udev.ukama.com}"
FACTORY_SEED_URL="${ULAB_FACTORY_SEED_URL:-https://factory-ukama.udev.ukama.com}"
SIM_TYPE="${UKAMA_LAB_SIM_TYPE:-ukama_data}"
P0_RUNS_DIR="${P0_RUNS_DIR:-runs/p0-batches}"
FACTORY_NODE_TARGET="${ULAB_FACTORY_NODE_COUNT:-0}"
FACTORY_HEADROOM_PERCENT="${ULAB_FACTORY_HEADROOM_PERCENT:-10}"
LIST_ONLY=0
FAIL_FAST=0
PREPARE_ONLY=0
SCENARIO_LIST_FILE=""
BATCH_ID_OVERRIDE=""
STATUS_FILE="${P0_STATUS_FILE:-}"

requested_categories=()
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
        --scenario-list)
            if (($# < 2)); then
                printf 'error: --scenario-list requires a file\n' >&2
                exit 2
            fi
            SCENARIO_LIST_FILE="$2"
            shift 2
            ;;
        --scenario-list=*)
            SCENARIO_LIST_FILE="${1#*=}"
            shift
            ;;
        --batch-id)
            if (($# < 2)); then
                printf 'error: --batch-id requires a value\n' >&2
                exit 2
            fi
            BATCH_ID_OVERRIDE="$2"
            shift 2
            ;;
        --batch-id=*)
            BATCH_ID_OVERRIDE="${1#*=}"
            shift
            ;;
        --prepare-only)
            PREPARE_ONLY=1
            shift
            ;;
        --list)
            LIST_ONLY=1
            shift
            ;;
        --fail-fast)
            FAIL_FAST=1
            shift
            ;;
        -h|--help)
            usage
            exit 0
            ;;
        --)
            shift
            requested_categories+=("$@")
            break
            ;;
        -*)
            printf 'error: unknown option: %s\n' "$1" >&2
            usage >&2
            exit 2
            ;;
        *)
            requested_categories+=("$1")
            shift
            ;;
    esac
done

if [[ "$FACTORY_NODE_TARGET" != "auto" &&
      ! "$FACTORY_NODE_TARGET" =~ ^[0-9]+$ ]]; then
    printf 'error: --factory-nodes must be auto or a non-negative integer\n' >&2
    exit 2
fi
if [[ ! "$FACTORY_HEADROOM_PERCENT" =~ ^[0-9]+$ ]]; then
    printf 'error: ULAB_FACTORY_HEADROOM_PERCENT must be a non-negative integer\n' >&2
    exit 2
fi
if [[ ! -d "$SCENARIO_ROOT" ]]; then
    printf 'error: scenario root does not exist: %s\n' "$SCENARIO_ROOT" >&2
    exit 2
fi
SCENARIO_ROOT_ABS="$(CDPATH= cd -- "$SCENARIO_ROOT" && pwd)"

# Discover every top-level P0 category dynamically.
mapfile -d '' available_categories < <(
    find "$SCENARIO_ROOT" -mindepth 1 -maxdepth 1 -type d \
        -printf '%f\0' | sort -z
)

if ((${#available_categories[@]} == 0)); then
    printf 'error: no P0 category directories found under %s\n' \
        "$SCENARIO_ROOT" >&2
    exit 2
fi

categories=()
selected_scenarios=()

if [[ -n "$SCENARIO_LIST_FILE" ]]; then
    if ((${#requested_categories[@]} > 0)); then
        printf 'error: categories cannot be combined with --scenario-list\n' >&2
        exit 2
    fi
    if [[ ! -r "$SCENARIO_LIST_FILE" ]]; then
        printf 'error: scenario list is not readable: %s\n' \
            "$SCENARIO_LIST_FILE" >&2
        exit 2
    fi

    declare -A seen_scenarios=()
    declare -A seen_categories=()
    while IFS= read -r scenario || [[ -n "$scenario" ]]; do
        scenario="${scenario%%#*}"
        scenario="$(sed 's/^[[:space:]]*//;s/[[:space:]]*$//' <<<"$scenario")"
        [[ -n "$scenario" ]] || continue
        if [[ "$scenario" != /* ]]; then
            scenario="$(pwd)/$scenario"
        fi
        if [[ ! -f "$scenario" || ! "$scenario" =~ \.ya?ml$ ]]; then
            printf 'error: invalid scenario in %s: %s\n' \
                "$SCENARIO_LIST_FILE" "$scenario" >&2
            exit 2
        fi
        if [[ -z "${seen_scenarios[$scenario]+x}" ]]; then
            selected_scenarios+=("$scenario")
            seen_scenarios[$scenario]=1
        fi

        scenario_abs="$(readlink -f -- "$scenario")"
        relative="${scenario_abs#${SCENARIO_ROOT_ABS%/}/}"
        category="${relative%%/*}"
        if [[ -n "$category" && -z "${seen_categories[$category]+x}" ]]; then
            categories+=("$category")
            seen_categories[$category]=1
        fi
    done <"$SCENARIO_LIST_FILE"
else
    if ((${#requested_categories[@]} == 0)); then
        categories=("${available_categories[@]}")
    else
        declare -A seen_categories=()
        for category in "${requested_categories[@]}"; do
            if [[ "$category" == "all" ]]; then
                categories=("${available_categories[@]}")
                seen_categories=()
                break
            fi
            if [[ ! "$category" =~ ^[a-z0-9-]+$ ||
                  ! -d "${SCENARIO_ROOT%/}/${category}" ]]; then
                printf 'error: unknown or missing P0 category: %s\n' \
                    "$category" >&2
                exit 2
            fi
            if [[ -z "${seen_categories[$category]+x}" ]]; then
                categories+=("$category")
                seen_categories[$category]=1
            fi
        done
    fi

    for category in "${categories[@]}"; do
        category_dir="${SCENARIO_ROOT%/}/${category}"
        while IFS= read -r -d '' scenario; do
            selected_scenarios+=("$scenario")
        done < <(
            find "$category_dir" -type f \
                \( -name '*.yaml' -o -name '*.yml' \) \
                -print0 | sort -z
        )
    done
fi

if ((${#selected_scenarios[@]} == 0)); then
    printf 'error: no scenario YAML files selected\n' >&2
    exit 2
fi

list_scenarios() {
    python3 - "$SCENARIO_ROOT" "${selected_scenarios[@]}" <<'PY'
import collections
import pathlib
import sys

root = pathlib.Path(sys.argv[1]).resolve()
counts = collections.Counter()

for value in sys.argv[2:]:
    path = pathlib.Path(value)
    status = "active"
    name = ""
    try:
        with path.open(encoding="utf-8") as stream:
            for raw in stream:
                line = raw.split("#", 1)[0].rstrip()
                if not line or line[0].isspace() or ":" not in line:
                    continue
                key, value = line.split(":", 1)
                key = key.strip()
                value = value.strip().strip("\"'")
                if key == "name":
                    name = value
                elif key == "status":
                    status = value or status
                if name and status != "active":
                    # We have both values needed for the common non-active case.
                    pass
    except OSError as exc:
        print(f"error\t{path}\t{exc}")
        counts["error"] += 1
        continue

    try:
        relative = path.resolve().relative_to(root)
    except ValueError:
        relative = path
    category = relative.parts[0] if len(relative.parts) > 1 else "_root"
    counts[status] += 1
    print(f"{status:<7}  {category:<16}  {relative}  {name}")

print("-" * 96)
print(
    "total={} active={} xfail={} wip={} skip={} error={}".format(
        sum(counts.values()),
        counts["active"],
        counts["xfail"],
        counts["wip"],
        counts["skip"],
        counts["error"],
    )
)
PY
}

if ((LIST_ONLY)); then
    printf 'P0 categories:'
    printf ' %s' "${categories[@]}"
    printf '\n\n'
    list_scenarios
    exit 0
fi

require_env UKAMA_IDENTIFIER || exit 2
require_env UKAMA_PASSWORD || exit 2
require_env PAUTH_URL || exit 2
require_env BFF_BASE_URL || exit 2

if [[ ! -x "$LAB_BIN" ]]; then
    printf 'error: ukama-lab is not executable: %s\n' "$LAB_BIN" >&2
    exit 2
fi
if [[ ! -d "$UKAMA_REPO" ]]; then
    printf 'error: Ukama repository does not exist: %s\n' "$UKAMA_REPO" >&2
    exit 2
fi

if [[ -n "$BATCH_ID_OVERRIDE" ]]; then
    if [[ ! "$BATCH_ID_OVERRIDE" =~ ^[A-Za-z0-9._-]+$ ]]; then
        printf 'error: invalid --batch-id: %s\n' "$BATCH_ID_OVERRIDE" >&2
        exit 2
    fi
    BATCH_STAMP="$BATCH_ID_OVERRIDE"
else
    BATCH_STAMP="$(date -u +%Y%m%dt%H%M%Sz)"
fi
BATCH_DIR="${P0_RUNS_DIR%/}/${BATCH_STAMP}"
RUNS_DIR="${BATCH_DIR}/runs"
LOGS_DIR="${BATCH_DIR}/logs"
SUMMARY_TSV="${BATCH_DIR}/scenarios.tsv"

mkdir -p "$RUNS_DIR" "$LOGS_DIR"
printf 'category\tscenario\trun_id\texit_code\treport\tlog\n' \
    >"$SUMMARY_TSV"

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
            value = value.strip().strip("\"'")

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
        print(
            f"error: cannot calculate factory nodes for {scenario}: {exc}",
            file=sys.stderr,
        )
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
        for key in (
            "nodes",
            "Nodes",
            "data",
            "Data",
            "items",
            "Items",
            "results",
            "Results",
        ):
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

if ((PREPARE_ONLY)); then
    printf 'factory preparation complete; no scenarios executed\n'
    exit 0
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

write_status() {
    local state="$1"
    local current="${2:--}"
    local message="${3:-}"

    [[ -n "$STATUS_FILE" ]] || return 0
    mkdir -p "$(dirname -- "$STATUS_FILE")"
    {
        printf 'state\t%s\n' "$state"
        printf 'total\t%s\n' "${#selected_scenarios[@]}"
        printf 'completed\t%s\n' "$completed_count"
        printf 'passed\t%s\n' "$passed_count"
        printf 'failed\t%s\n' "$failed_count"
        printf 'skipped\t%s\n' "$skipped_count"
        printf 'current\t%s\n' "$current"
        printf 'message\t%s\n' "$message"
        printf 'started_at\t%s\n' "$STATUS_STARTED_AT"
        printf 'updated_at\t%s\n' "$(date -u +%Y-%m-%dT%H:%M:%SZ)"
    } >"$STATUS_FILE.tmp"
    mv "$STATUS_FILE.tmp" "$STATUS_FILE"
}

scenario_outcome() {
    local report_path="$1"
    local rc="$2"

    python3 - "$report_path" "$rc" <<'PY_STATUS'
import json
import pathlib
import sys

path = pathlib.Path(sys.argv[1])
rc = int(sys.argv[2])
report = {}
if path.is_file():
    try:
        report = json.loads(path.read_text(encoding="utf-8"))
    except (OSError, json.JSONDecodeError):
        pass
status = report.get("status", "")
if rc == 0 and status in {"wip", "skip"}:
    print("SKIP")
elif rc == 0 and report.get("passed") is True:
    print("PASS")
else:
    print("FAIL")
PY_STATUS
}

printf '\n============================================================\n'
printf 'P0 batch: %s\n' "$BATCH_STAMP"
printf 'Categories (%s):' "${#categories[@]}"
printf ' %s' "${categories[@]}"
printf '\nScenarios selected: %s\n' "${#selected_scenarios[@]}"
printf '============================================================\n'

batch_failed=0
completed_count=0
passed_count=0
failed_count=0
skipped_count=0
STATUS_STARTED_AT="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
write_status RUNNING - 'starting batch'

for scenario in "${selected_scenarios[@]}"; do
    scenario_abs="$(readlink -f -- "$scenario")"
    relative="${scenario_abs#${SCENARIO_ROOT_ABS%/}/}"
    category="${relative%%/*}"
    slug="${relative%.yaml}"
    slug="${slug%.yml}"
    slug="${slug//\//-}"
    slug="${slug,,}"
    slug="${slug//[^a-z0-9-]/-}"
    run_id="p0-${BATCH_STAMP}-${slug}"
    report_path="${RUNS_DIR}/${run_id}/report.json"
    log_path="${LOGS_DIR}/${slug}.log"

    write_status RUNNING "$relative" 'scenario started'
    printf '\n-- Running: %s --\n' "$relative"
    "$LAB_BIN" validate "$scenario" \
        "${common_args[@]}" \
        --run-id "$run_id" 2>&1 | tee "$log_path"
    rc=${PIPESTATUS[0]}

    printf '%s\t%s\t%s\t%s\t%s\t%s\n' \
        "$category" "$relative" "$run_id" "$rc" \
        "$report_path" "$log_path" >>"$SUMMARY_TSV"

    outcome="$(scenario_outcome "$report_path" "$rc")"
    completed_count=$((completed_count + 1))
    case "$outcome" in
        PASS) passed_count=$((passed_count + 1)) ;;
        SKIP) skipped_count=$((skipped_count + 1)) ;;
        *)
            failed_count=$((failed_count + 1))
            batch_failed=1
            ;;
    esac
    write_status RUNNING - "last=$outcome $relative"

    if ((rc != 0 && FAIL_FAST)); then
        printf 'stopping after failed scenario (--fail-fast): %s\n' \
            "$relative" >&2
        break
    fi
done

write_status FINALIZING - 'building batch report'

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

        results.append(
            {
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
            }
        )

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
        f"{item['outcome']:<4}  {item['category']:<18}  "
        f"{item['scenario']}  ({item['duration_sec']}s)"
    )
text_path.write_text("\n".join(lines) + "\n", encoding="utf-8")

print("\n============================================================")
print("P0 batch summary")
print("============================================================")
print(f"{'RESULT':<7} {'CATEGORY':<19} SCENARIO")
for item in results:
    print(f"{item['outcome']:<7} {item['category']:<19} {item['scenario']}")
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
report_rc=$?
write_status DONE - "batch_report_exit_code=$report_rc"
exit "$report_rc"
