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

pick_tnode() {
    json="$1"

    jq -r '
      [
        .. | objects
        | select((.nodeId? // .nodeID? // .node_id? // .id?) != null)
        | {
            id:   ((.nodeId // .nodeID // .node_id // .id) | tostring),
            type: ((.nodeType // .node_type // .type // .kind // "") |
                   tostring | ascii_downcase)
          }
        | select((.id | ascii_downcase | test("-tnode-")) or
                 (.type == "tnode"))
        | .id
      ]
      | unique
      | first // empty
    ' "$json"
}

derive_node_id() {
    node_id="$1"
    kind="$2"

    case "$node_id" in
        *-tnode-*)
            echo "$node_id" | sed "s/-tnode-/-$kind-/"
            ;;
        *)
            echo "cannot derive $kind from non-tnode id: $node_id" >&2
            exit 1
            ;;
    esac
}

mark_provisioned() {
    node_id="$1"

    echo "factory: marking node provisioned $node_id"
    curl -fsS -X PATCH \
        "$FACTORY_URL/v1/nodefactory/node/$node_id" \
        -H "accept: application/json" \
        >/dev/null
}

assign_org() {
    node_id="$1"

    echo "factory: assigning node $node_id to org $FACTORY_ORG"
    curl -fsS -X PATCH \
        "$FACTORY_URL/v1/nodefactory/node/$node_id/org/$FACTORY_ORG" \
        -H "accept: application/json" \
        >/dev/null
}

container_name() {
    echo "ukama-vnode-$(safe_name "$1")"
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
        echo "SLOT=$SLOT"
        echo "LAB_NET=${LAB_NET:-}"
    } > "$state_file"

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
cleanup() {
    rm -f "$FACTORY_JSON"
}
trap cleanup EXIT INT TERM

echo "factory: fetching unprovisioned tnodes for site=$SITE_REF"

curl -fsS -X GET \
    "$FACTORY_URL/v1/nodefactory/nodes?isProvisioned=false" \
    -H "accept: application/json" \
    > "$FACTORY_JSON"
TNODE_ID="$(pick_tnode "$FACTORY_JSON")"

if [ -z "$TNODE_ID" ]; then
    echo "no available tnode factory node for site=$SITE_REF" >&2
    exit 1
fi

CNODE_ID="$(derive_node_id "$TNODE_ID" cnode)"
ANODE_ID="$(derive_node_id "$TNODE_ID" anode)"

TNODE_CONTAINER="$(container_name "$TNODE_ID")"
CNODE_CONTAINER="$(container_name "$CNODE_ID")"
ANODE_CONTAINER="$(container_name "$ANODE_ID")"

echo "factory: selected site bundle site=$SITE_REF tnode=$TNODE_ID cnode=$CNODE_ID anode=$ANODE_ID"

# Build first. Do not mutate factory state until all three images build.
"$SCRIPT_DIR/build-node.sh" "$REPO" "$TNODE_ID" "$NODE_RUNTIME"
"$SCRIPT_DIR/build-node.sh" "$REPO" "$CNODE_ID" "$NODE_RUNTIME"
"$SCRIPT_DIR/build-node.sh" "$REPO" "$ANODE_ID" "$NODE_RUNTIME"

# Only after all builds succeed, mark all three provisioned.
mark_provisioned "$TNODE_ID"
mark_provisioned "$CNODE_ID"
mark_provisioned "$ANODE_ID"

# Only after all three are provisioned, assign all three to the org.
assign_org "$TNODE_ID"
assign_org "$CNODE_ID"
assign_org "$ANODE_ID"

LAB_NET="$LAB_NET" "$SCRIPT_DIR/start-node.sh" "$TNODE_ID" "$TNODE_CONTAINER" "$RUN_DIR"
LAB_NET="$LAB_NET" "$SCRIPT_DIR/start-node.sh" "$CNODE_ID" "$CNODE_CONTAINER" "$RUN_DIR"
LAB_NET="$LAB_NET" "$SCRIPT_DIR/start-node.sh" "$ANODE_ID" "$ANODE_CONTAINER" "$RUN_DIR"

write_site_state
write_node_state "$SITE_REF-tnode" "$TNODE_ID" "tnode" "$TNODE_CONTAINER"
write_node_state "$SITE_REF-cnode" "$CNODE_ID" "cnode" "$CNODE_CONTAINER"
write_node_state "$SITE_REF-anode" "$ANODE_ID" "anode" "$ANODE_CONTAINER"

echo "site-started site=$SITE_REF tnode=$TNODE_ID cnode=$CNODE_ID anode=$ANODE_ID"
