#!/bin/sh
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

set -eu

if [ "$#" -lt 4 ]; then
    echo "usage: $0 <repo> <site-ref> <network-ref> <run-dir> [slot]" >&2
    exit 2
fi

REPO="$1"
SITE_REF="$2"
NETWORK_REF="$3"
RUN_DIR="$4"
SLOT="${5:-0}"

FACTORY_URL="${FACTORY_URL:-https://factory-ukama.udev.ukama.com}"
FACTORY_ORG="${FACTORY_ORG:-ukama}"
NODE_RUNTIME="${NODE_RUNTIME:-starter}"

SCRIPT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)"
STATE_DIR="$RUN_DIR/runtime-sites"
NODE_STATE_DIR="$RUN_DIR/runtime-nodes"
NET_STATE="$RUN_DIR/runtime-net/net.env"
mkdir -p "$STATE_DIR" "$NODE_STATE_DIR"

need_cmd() {
    if ! command -v "$1" >/dev/null 2>&1; then
        echo "missing required command: $1" >&2
        exit 1
    fi
}

safe_name() {
    printf "%s" "$1" | tr -c 'A-Za-z0-9_.-' '-'
}

list_site_bundles() {
    json="$1"

    jq -r --arg mode "${2:-complete}" '
      def nid:
        ((.nodeId // .nodeID // .node_id // .id // "") | tostring);

      def ntype:
        ((.nodeType // .node_type // .type // .kind // "") |
         tostring | ascii_downcase);

      def assoc_tnode:
        ((.associatedTNodeId // .associated_tnode_id //
          .associatedTnodeId // .associated_tnode //
          .tnodeId // .tnode_id //
          .parentTNodeId // .parent_tnode_id // "") | tostring);

      def kind:
        if (.type == "tnode" or (.id | ascii_downcase | test("-tnode-"))) then
          "tnode"
        elif (.type == "cnode" or (.id | ascii_downcase | test("-cnode-"))) then
          "cnode"
        elif (.type == "anode" or (.id | ascii_downcase | test("-anode-"))) then
          "anode"
        else
          .type
        end;

      def derived($id; $kind):
        ($id | gsub("-tnode-"; "-" + $kind + "-"));

      [
        .. | objects
        | select((nid | length) > 0)
        | select(((.isProvisioned // .is_provisioned //
                   .provisioned // false) | tostring |
                   ascii_downcase) != "true")
        | {id: nid, type: ntype, associated_tnode_id: assoc_tnode}
        | .kind = kind
      ]
      | unique_by(.id)
      as $nodes
      | (reduce $nodes[] as $node (
          {by_id: {}, by_parent: {}};
          .by_id[$node.id] = $node
          | if $node.associated_tnode_id != "" then
              .by_parent[$node.associated_tnode_id][$node.kind]
                //= $node.id
            else . end
        )) as $index

      | def child($tnode; $kind):
          [
            ($index.by_id[derived($tnode; $kind)]
             | select(.kind == $kind) | .id),
            ($index.by_parent[$tnode][$kind] // empty)
          ] | min // "";

        $nodes[]
        | select(.kind == "tnode")
        | .id as $tnode
        | child($tnode; "cnode") as $cnode
        | child($tnode; "anode") as $anode
        | if $mode == "incomplete" then
            select($cnode == "" or $anode == "")
            | "skip incomplete bundle tnode=\($tnode)" +
              " cnode=\($cnode != "") anode=\($anode != "")"
          else
            select($cnode != "" and $anode != "")
            | [$tnode, $cnode, $anode] | @tsv
          end
    ' "$json"
}

print_incomplete_bundles() {
    list_site_bundles "$1" incomplete | head -n 20 >&2 || true
}

reserve_node() {
    node_id="$1"
    body="$(mktemp "${TMPDIR:-/tmp}/ukama-reserve-node.XXXXXX")"

    echo "factory: reserving node $node_id"
    if ! code="$(curl -sS --connect-timeout 10 --max-time 30 \
        -o "$body" -w "%{http_code}" -X PATCH \
        "$FACTORY_URL/v1/nodefactory/node/$node_id" \
        -H "accept: application/json")"; then
        echo "factory: reserve request failed node=$node_id;" \
            "reservation state unknown" >&2
        rm -f "$body"
        exit 1
    fi

    case "$code" in
        2??)
            if [ "${ULAB_WEBAPP_RUN:-0}" = 1 ]; then
                printf "%s\t%s\n" "$SITE_REF" "$node_id" >> "$STATE_DIR/factory-claims.tsv"
            fi
            rm -f "$body"
            return 0
            ;;
        409)
            echo "factory: reservation conflict node=$node_id status=409" >&2
            cat "$body" >&2 || true
            printf '\n' >&2
            rm -f "$body"
            return 1
            ;;
        *)
            echo "factory: reserve failed node=$node_id status=$code" >&2
            cat "$body" >&2 || true
            rm -f "$body"
            exit 1
            ;;
    esac
}

assign_org() {
    node_id="$1"
    body="$(mktemp "${TMPDIR:-/tmp}/ukama-assign-org.XXXXXX")"

    echo "factory: assigning node $node_id to org $FACTORY_ORG"
    if ! code="$(curl -sS --connect-timeout 10 --max-time 30 \
        -o "$body" -w "%{http_code}" -X PATCH \
        "$FACTORY_URL/v1/nodefactory/node/$node_id/org/$FACTORY_ORG" \
        -H "accept: application/json")"; then
        echo "factory: assign request failed node=$node_id" >&2
        rm -f "$body"
        exit 1
    fi

    case "$code" in
        2??)
            rm -f "$body"
            ;;
        *)
            echo "factory: assign failed node=$node_id status=$code" >&2
            cat "$body" >&2 || true
            rm -f "$body"
            exit 1
            ;;
    esac
}

container_name() {
    echo "ukama-vnode-$(safe_name "$1")"
}

container_ip_on_network() {
    container="$1"
    network="$2"

    podman inspect -f '{{range $name, $net := .NetworkSettings.Networks}}{{if eq $name "'"$network"'"}}{{$net.IPAddress}}{{end}}{{end}}' \
        "$container" 2>/dev/null
}

write_node_state() {
    logical_node_id="$1"
    node_id="$2"
    node_kind="$3"
    container="$4"
    state_file="$NODE_STATE_DIR/$(safe_name "$logical_node_id").env"

    {
        echo "LOGICAL_NODE_ID=$logical_node_id"
        echo "FACTORY_NODE_ID=$node_id"
        echo "NODE_TYPE=$node_kind"
        echo "NODE_KIND=$node_kind"
        echo "SITE_REF=$SITE_REF"
        echo "NETWORK_REF=$NETWORK_REF"
        echo "CONTAINER_NAME=$container"
        echo "IMAGE=testing/virtualnode:$node_id"
        echo "SLOT=$SLOT"
        echo "LAB_NET=${LAB_NET:-}"
    } > "$state_file"
}

write_site_state() {
    state_file="$STATE_DIR/$(safe_name "$SITE_REF").env"

    {
        echo "SITE_REF=$SITE_REF"
        echo "NETWORK_REF=$NETWORK_REF"
        echo "TNODE_ID=$TNODE_ID"
        echo "CNODE_ID=$CNODE_ID"
        echo "ANODE_ID=$ANODE_ID"
        echo "TNODE_CONTAINER=$TNODE_CONTAINER"
        echo "CNODE_CONTAINER=$CNODE_CONTAINER"
        echo "ANODE_CONTAINER=$ANODE_CONTAINER"
        echo "TNODE_IP=${TNODE_IP:-}"
        echo "CNODE_IP=${CNODE_IP:-}"
        echo "ANODE_IP=${ANODE_IP:-}"
        echo "SLOT=$SLOT"
        echo "LAB_NET=${LAB_NET:-}"
    } > "$state_file"

    [ "${1:-complete}" != early ] || return 0
    {
        printf "%s\t%s\t%s\t%s\t%s\t%s\n" \
            "$SITE_REF" "$NETWORK_REF" "$TNODE_ID" "$CNODE_ID" \
            "$ANODE_ID" "$SLOT"
    } >> "$STATE_DIR/sites.tsv"
}

need_cmd curl
need_cmd jq
need_cmd podman

if [ ! -f "$NET_STATE" ]; then
    echo "lab network state not found: $NET_STATE" >&2
    echo "runtime must create the lab network before starting site nodes" >&2
    exit 1
fi

# shellcheck disable=SC1090
. "$NET_STATE"

if [ -z "${LAB_NET:-}" ]; then
    echo "LAB_NET missing in $NET_STATE" >&2
    exit 1
fi

FACTORY_JSON="$(mktemp "${TMPDIR:-/tmp}/ukama-factory-sites.XXXXXX")"
FACTORY_BUNDLES=""
cleanup() {
    rm -f "$FACTORY_JSON" "$FACTORY_BUNDLES"
}
trap cleanup EXIT
trap 'exit 130' INT
trap 'exit 143' TERM
FACTORY_BUNDLES="$(mktemp "${TMPDIR:-/tmp}/ukama-factory-bundles.XXXXXX")"

echo "factory: fetching unprovisioned factory nodes for site=$SITE_REF"

curl -fsS --connect-timeout 10 --max-time 30 -X GET \
    "$FACTORY_URL/v1/nodefactory/nodes?isProvisioned=false" \
    -H "accept: application/json" > "$FACTORY_JSON"

echo "factory: finding complete unprovisioned bundles for site=$SITE_REF"
list_site_bundles "$FACTORY_JSON" > "$FACTORY_BUNDLES"

if [ ! -s "$FACTORY_BUNDLES" ]; then
    echo "no complete unprovisioned factory bundle for site=$SITE_REF" >&2
    print_incomplete_bundles "$FACTORY_JSON"
    exit 1
fi

BUNDLE_COUNT="$(wc -l < "$FACTORY_BUNDLES" | tr -d '[:space:]')"
echo "factory: found $BUNDLE_COUNT candidate bundles for site=$SITE_REF"
BUNDLE_RESERVED=0

while IFS="$(printf '\t')" read -r TNODE_ID CNODE_ID ANODE_ID; do
    echo "factory: selected bundle site=$SITE_REF tnode=$TNODE_ID" \
        "cnode=$CNODE_ID anode=$ANODE_ID"

    # The tower node is the bundle claim. Factory updates this flag
    # atomically, so only one lab runner can receive a successful response.
    if ! reserve_node "$TNODE_ID"; then
        echo "factory: bundle already reserved tnode=$TNODE_ID;" \
            "trying next candidate"
        continue
    fi

    if ! reserve_node "$CNODE_ID"; then
        echo "factory: incomplete reservation cnode=$CNODE_ID" >&2
        exit 1
    fi
    if ! reserve_node "$ANODE_ID"; then
        echo "factory: incomplete reservation anode=$ANODE_ID" >&2
        exit 1
    fi

    assign_org "$TNODE_ID"
    assign_org "$CNODE_ID"
    assign_org "$ANODE_ID"

    echo "factory: reserved bundle site=$SITE_REF tnode=$TNODE_ID" \
        "cnode=$CNODE_ID anode=$ANODE_ID"
    BUNDLE_RESERVED=1
    break
done < "$FACTORY_BUNDLES"

if [ "$BUNDLE_RESERVED" != 1 ]; then
    echo "factory: all $BUNDLE_COUNT candidate bundles conflicted;" \
        "no bundle reserved for site=$SITE_REF" >&2
    exit 1
fi

TNODE_CONTAINER="$(container_name "$TNODE_ID")"
CNODE_CONTAINER="$(container_name "$CNODE_ID")"
ANODE_CONTAINER="$(container_name "$ANODE_ID")"

# Persist the exclusively reserved, assigned bundle before any build/start can
# fail. The webapp parent reloads this state for bounded teardown on failure.
if [ "${ULAB_WEBAPP_RUN:-0}" = 1 ]; then
    write_site_state early
    write_node_state "$SITE_REF-tnode" "$TNODE_ID" "tnode" "$TNODE_CONTAINER"
    write_node_state "$SITE_REF-cnode" "$CNODE_ID" "cnode" "$CNODE_CONTAINER"
    write_node_state "$SITE_REF-anode" "$ANODE_ID" "anode" "$ANODE_CONTAINER"
fi

# Build only after the complete Factory bundle has been reserved.
"$SCRIPT_DIR/build-node.sh" "$REPO" "$TNODE_ID" "$NODE_RUNTIME"
"$SCRIPT_DIR/build-node.sh" "$REPO" "$CNODE_ID" "$NODE_RUNTIME"
"$SCRIPT_DIR/build-node.sh" "$REPO" "$ANODE_ID" "$NODE_RUNTIME"

"$SCRIPT_DIR/start-node.sh" "$TNODE_ID" "$TNODE_CONTAINER" "$RUN_DIR"
"$SCRIPT_DIR/start-node.sh" "$CNODE_ID" "$CNODE_CONTAINER" "$RUN_DIR"
"$SCRIPT_DIR/start-node.sh" "$ANODE_ID" "$ANODE_CONTAINER" "$RUN_DIR"

TNODE_IP="$(container_ip_on_network "$TNODE_CONTAINER" "$LAB_NET")"
CNODE_IP="$(container_ip_on_network "$CNODE_CONTAINER" "$LAB_NET")"
ANODE_IP="$(container_ip_on_network "$ANODE_CONTAINER" "$LAB_NET")"

if [ -z "$TNODE_IP" ] || [ -z "$CNODE_IP" ] || [ -z "$ANODE_IP" ]; then
    echo "one or more site containers have no IP on $LAB_NET" >&2
    podman ps -a --filter "name=ukama-vnode-" >&2 || true
    exit 1
fi

write_site_state
write_node_state "$SITE_REF-tnode" "$TNODE_ID" "tnode" "$TNODE_CONTAINER"
write_node_state "$SITE_REF-cnode" "$CNODE_ID" "cnode" "$CNODE_CONTAINER"
write_node_state "$SITE_REF-anode" "$ANODE_ID" "anode" "$ANODE_CONTAINER"

echo "site-started site=$SITE_REF tnode=$TNODE_ID cnode=$CNODE_ID anode=$ANODE_ID"
