#!/usr/bin/env bash
#
# Generate N Ukama node sets in factory via API only.
#
# For each generated tower node, the script creates:
#   - one tnode using factory random ID generation
#   - one cnode derived from that tnode using associatedTNodeId
#   - one anode derived from that tnode using associatedTNodeId
#
# No direct DB access is used.
#
# Usage:
#   FACTORY_URL=http://localhost:8070 ./generate_factory_nodes.sh <count> [org_name]
#
# Examples:
#   ./generate_factory_nodes.sh 100
#   ./generate_factory_nodes.sh 100 ukama
#
# Env:
#   FACTORY_URL          Base URL for factory API gateway. Default: http://localhost:8070
#   PROVISION            true/false. Default: false
#   ALLOCATE             true/false. Default: false
#   FACTORY_AUTH_HEADER  Optional extra auth header, e.g. 'X-Session-Token: xxx'
#   OUT                  Output CSV path. Default: generated_factory_nodes_<timestamp>.csv
#

set -eu

usage() {
    cat >&2 <<USAGE
Usage:
  FACTORY_URL=http://localhost:8070 $0 <count> [org_name]

Examples:
  $0 100
  $0 100 ukama

Env:
  FACTORY_URL          Factory API gateway URL. Default: http://localhost:8070
  PROVISION            true/false. Default: false
  ALLOCATE             true/false. Default: false
  FACTORY_AUTH_HEADER  Optional extra header, e.g. 'X-Session-Token: xxx'
  OUT                  Output CSV path
USAGE
    exit 1
}

if [ "$#" -lt 1 ] || [ "$#" -gt 2 ]; then
    usage
fi

COUNT="$1"
ORG_NAME="${2:-}"
FACTORY_URL="${FACTORY_URL:-http://localhost:8070}"
PROVISION="${PROVISION:-false}"
ALLOCATE="${ALLOCATE:-false}"
OUT="${OUT:-generated_factory_nodes_$(date +%Y%m%d_%H%M%S).csv}"

case "$COUNT" in
    ''|*[!0-9]*)
        echo "error: count must be a positive integer" >&2
        exit 1
        ;;
esac

if [ "$COUNT" -le 0 ]; then
    echo "error: count must be > 0" >&2
    exit 1
fi

case "$PROVISION" in
    true|false) ;;
    *) echo "error: PROVISION must be true or false" >&2; exit 1 ;;
esac

case "$ALLOCATE" in
    true|false) ;;
    *) echo "error: ALLOCATE must be true or false" >&2; exit 1 ;;
esac

if [ "$ALLOCATE" = "true" ] && [ -z "$ORG_NAME" ]; then
    echo "error: org_name is required when ALLOCATE=true" >&2
    exit 1
fi

need_cmd() {
    if ! command -v "$1" >/dev/null 2>&1; then
        echo "error: missing required command: $1" >&2
        exit 1
    fi
}

need_cmd curl
need_cmd python3
need_cmd date

add_curl_headers() {
    CURL_OPTS+=("-H" "accept: application/json")

    if [ -n "${FACTORY_AUTH_HEADER:-}" ]; then
        CURL_OPTS+=("-H" "$FACTORY_AUTH_HEADER")
    fi
}

api_call() {
    method="$1"
    path="$2"

    tmp_body="$(mktemp)"

    CURL_OPTS=(-sS)
    add_curl_headers
    CURL_OPTS+=(
        -X "$method"
        -o "$tmp_body"
        -w '%{http_code}'
        "${FACTORY_URL}${path}"
    )

    code="$(curl "${CURL_OPTS[@]}" || true)"

    case "$code" in
        ''|*[!0-9]*)
            echo "error: ${method} ${path} returned invalid HTTP code: ${code}" >&2
            cat "$tmp_body" >&2 || true
            echo >&2
            rm -f "$tmp_body"
            exit 1
            ;;
    esac

    if [ "$code" -lt 200 ] || [ "$code" -ge 300 ]; then
        echo "error: ${method} ${path} failed with HTTP ${code}" >&2
        cat "$tmp_body" >&2
        echo >&2
        rm -f "$tmp_body"
        exit 1
    fi

    cat "$tmp_body"
    rm -f "$tmp_body"
}

api_call_allow_conflict() {
    method="$1"
    path="$2"

    tmp_body="$(mktemp)"

    CURL_OPTS=(-sS)
    add_curl_headers
    CURL_OPTS+=(
        -X "$method"
        -o "$tmp_body"
        -w '%{http_code}'
        "${FACTORY_URL}${path}"
    )

    code="$(curl "${CURL_OPTS[@]}" || true)"

    case "$code" in
        ''|*[!0-9]*)
            echo "error: ${method} ${path} returned invalid HTTP code: ${code}" >&2
            cat "$tmp_body" >&2 || true
            echo >&2
            rm -f "$tmp_body"
            exit 1
            ;;
    esac

    if [ "$code" -ge 200 ] && [ "$code" -lt 300 ]; then
        rm -f "$tmp_body"
        return 0
    fi

    if [ "$code" = "409" ]; then
        echo "warn: ${method} ${path} returned 409, treating as already done" >&2
        rm -f "$tmp_body"
        return 0
    fi

    echo "error: ${method} ${path} failed with HTTP ${code}" >&2
    cat "$tmp_body" >&2
    echo >&2
    rm -f "$tmp_body"
    exit 1
}

json_id() {
    python3 -c '
import json, sys
try:
    data = json.load(sys.stdin)
except Exception as e:
    print(f"error: failed to parse JSON: {e}", file=sys.stderr)
    sys.exit(1)
node_id = data.get("id") or data.get("Id")
if not node_id:
    print("error: response does not contain id", file=sys.stderr)
    print(data, file=sys.stderr)
    sys.exit(1)
print(node_id)
'
}

urlenc() {
    python3 -c 'import sys, urllib.parse; print(urllib.parse.quote(sys.argv[1], safe=""))' "$1"
}

generate_node() {
    node_type="$1"
    tnode_id="${2:-}"

    path="/v1/nodefactory/node/${node_type}/generate"
    if [ -n "$tnode_id" ]; then
        path="${path}?associatedTNodeId=$(urlenc "$tnode_id")"
    fi

    api_call POST "$path" | json_id
}

provision_node() {
    node_id="$1"
    api_call_allow_conflict PATCH "/v1/nodefactory/node/${node_id}"
}

allocate_node() {
    node_id="$1"
    api_call_allow_conflict PATCH "/v1/nodefactory/node/${node_id}/org/${ORG_NAME}"
}

printf 'set,tnode,cnode,anode,org,provisioned,allocated\n' >"$OUT"

i=1
while [ "$i" -le "$COUNT" ]; do
    echo "factory: generating node set ${i}/${COUNT}" >&2

    tnode_id="$(generate_node tnode)"
    cnode_id="$(generate_node cnode "$tnode_id")"
    anode_id="$(generate_node anode "$tnode_id")"

    if [ "$PROVISION" = "true" ]; then
        echo "factory: provisioning ${tnode_id} ${cnode_id} ${anode_id}" >&2
        provision_node "$tnode_id"
        provision_node "$cnode_id"
        provision_node "$anode_id"
    fi

    if [ "$ALLOCATE" = "true" ]; then
        echo "factory: allocating ${tnode_id} ${cnode_id} ${anode_id} to ${ORG_NAME}" >&2
        allocate_node "$tnode_id"
        allocate_node "$cnode_id"
        allocate_node "$anode_id"
    fi

    printf '%s,%s,%s,%s,%s,%s,%s\n' \
        "$i" \
        "$tnode_id" \
        "$cnode_id" \
        "$anode_id" \
        "$ORG_NAME" \
        "$PROVISION" \
        "$ALLOCATE" >>"$OUT"

    i=$((i + 1))
done

echo "factory: done"
echo "factory: wrote ${OUT}"
