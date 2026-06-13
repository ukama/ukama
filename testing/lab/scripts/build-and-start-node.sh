#!/bin/sh
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

set -eu

if [ "$#" -lt 6 ]; then
    echo "usage: $0 <repo> <logical-node-id> <node-type> <site-ref> <network-ref> <run-dir> [slot]" >&2
    exit 2
fi

REPO="$1"
LOGICAL_NODE_ID="$2"
NODE_TYPE="$3"
SITE_REF="$4"
NETWORK_REF="$5"
RUN_DIR="$6"
SLOT="${7:-0}"

FACTORY_URL="${FACTORY_URL:-https://factory-ukama.udev.ukama.com}"
FACTORY_ORG="${FACTORY_ORG:-ukama}"
NODE_RUNTIME="${NODE_RUNTIME:-starter}"
IMAGE_REPO="${IMAGE_REPO:-testing/virtualnode}"

STATE_DIR="$RUN_DIR/runtime-nodes"
mkdir -p "$STATE_DIR"

need_cmd() {
    if ! command -v "$1" >/dev/null 2>&1; then
        echo "missing required command: $1" >&2
        exit 1
    fi
}

safe_name() {
    echo "$1" | tr -c 'A-Za-z0-9_.-' '-'
}

node_kind_for_type() {
    case "$1" in
        tower|tnode)
            echo "tnode"
            ;;
        controller|cnode)
            echo "cnode"
            ;;
        amplifier|amp|anode)
            echo "anode"
            ;;
        *)
            echo "$1"
            ;;
    esac
}

find_node_dir() {
    if [ -x "$REPO/testing/node/mk_local_vnode.sh" ]; then
        echo "$REPO/testing/node"
        return
    fi

    if [ -x "$REPO/node/mk_local_vnode.sh" ]; then
        echo "$REPO/node"
        return
    fi

    if [ -x "$REPO/mk_local_vnode.sh" ]; then
        echo "$REPO"
        return
    fi

    echo "mk_local_vnode.sh not found under repo: $REPO" >&2
    exit 1
}

pick_factory_node() {
    json="$1"
    kind="$2"

    jq -r --arg kind "$kind" '
      [
        .. | objects
        | select((.nodeId? // .nodeID? // .node_id? // .id?) != null)
        | {
            id:   ((.nodeId // .nodeID // .node_id // .id) | tostring),
            type: ((.nodeType // .node_type // .type // .kind // "") | tostring | ascii_downcase)
          }
        | select(
            (.id | ascii_downcase | test("-" + $kind + "-"))
            or
            (.type == $kind)
          )
        | .id
      ]
      | unique
      | first // empty
    ' "$json"
}

write_state() {
    state_file="$STATE_DIR/$(safe_name "$LOGICAL_NODE_ID").env"

    {
        echo "LOGICAL_NODE_ID=$LOGICAL_NODE_ID"
        echo "FACTORY_NODE_ID=$FACTORY_NODE_ID"
        echo "NODE_TYPE=$NODE_TYPE"
        echo "NODE_KIND=$NODE_KIND"
        echo "SITE_REF=$SITE_REF"
        echo "NETWORK_REF=$NETWORK_REF"
        echo "CONTAINER_NAME=$CONTAINER_NAME"
        echo "IMAGE=$IMAGE"
        echo "SLOT=$SLOT"
    } > "$state_file"

    {
        printf "%s\t%s\t%s\t%s\t%s\t%s\t%s\n" \
            "$LOGICAL_NODE_ID" \
            "$FACTORY_NODE_ID" \
            "$NODE_TYPE" \
            "$NODE_KIND" \
            "$SITE_REF" \
            "$NETWORK_REF" \
            "$CONTAINER_NAME"
    } >> "$STATE_DIR/nodes.tsv"
}

need_cmd curl
need_cmd jq
need_cmd podman

NODE_KIND="$(node_kind_for_type "$NODE_TYPE")"
NODE_DIR="$(find_node_dir)"
CONTAINER_NAME="ukama-vnode-$(safe_name "$LOGICAL_NODE_ID")"
FACTORY_JSON="$(mktemp "${TMPDIR:-/tmp}/ukama-factory-nodes.XXXXXX")"

cleanup() {
    rm -f "$FACTORY_JSON"
}
trap cleanup EXIT INT TERM

echo "factory: fetching unprovisioned nodes for type=$NODE_TYPE kind=$NODE_KIND"

curl -fsS -X GET \
    "$FACTORY_URL/v1/nodefactory/nodes?isProvisioned=false" \
    -H "accept: application/json" \
    > "$FACTORY_JSON"

FACTORY_NODE_ID="$(pick_factory_node "$FACTORY_JSON" "$NODE_KIND")"

if [ -z "$FACTORY_NODE_ID" ]; then
    echo "no available factory node for logical node=$LOGICAL_NODE_ID type=$NODE_TYPE kind=$NODE_KIND" >&2
    exit 1
fi

IMAGE="$IMAGE_REPO:$FACTORY_NODE_ID"

echo "factory: selected $FACTORY_NODE_ID for logical node $LOGICAL_NODE_ID"
echo "build: $NODE_DIR/mk_local_vnode.sh --node-id $FACTORY_NODE_ID --runtime $NODE_RUNTIME"

(
    cd "$NODE_DIR"
    ./mk_local_vnode.sh --node-id "$FACTORY_NODE_ID" --runtime "$NODE_RUNTIME"
)

echo "factory: marking node provisioned $FACTORY_NODE_ID"

curl -fsS -X PATCH \
    "$FACTORY_URL/v1/nodefactory/node/$FACTORY_NODE_ID" \
    -H "accept: application/json" \
    >/dev/null

echo "factory: assigning node $FACTORY_NODE_ID to org $FACTORY_ORG"

curl -fsS -X PATCH \
    "$FACTORY_URL/v1/nodefactory/node/$FACTORY_NODE_ID/org/$FACTORY_ORG" \
    -H "accept: application/json" \
    >/dev/null

echo "podman: removing existing container if present: $CONTAINER_NAME"
podman rm -f "$CONTAINER_NAME" >/dev/null 2>&1 || true

PUBLISH_ARGS=""

#
# Do not publish host ports by default because scale scenarios start many nodes.
# For one-node local debugging, run with:
#
#   ULAB_PUBLISH_NODE_PORTS=1 ./ulab ...
#
if [ "${ULAB_PUBLISH_NODE_PORTS:-0}" = "1" ]; then
    PUBLISH_ARGS="-p 18001:18001 -p 18026:18026 -p 18028:18028 -p 18029:18029/udp -p 18030:18030"
fi

echo "podman: starting $CONTAINER_NAME from $IMAGE"

if [ -n "${ULAB_NODE_ENTRYPOINT:-}" ]; then
    # Debug/override mode, for example:
    #   ULAB_NODE_ENTRYPOINT=/bin/bash ./ulab ...
    # or:
    #   ULAB_NODE_ENTRYPOINT=/bin/bash ULAB_NODE_CMD="-lc 'env && sleep infinity'" ./ulab ...
    # shellcheck disable=SC2086
    podman run -d \
        --name "$CONTAINER_NAME" \
        --privileged \
        --device /dev/net/tun \
        --entrypoint "$ULAB_NODE_ENTRYPOINT" \
        $PUBLISH_ARGS \
        "$IMAGE" \
        ${ULAB_NODE_CMD:-}
else
    # Normal mode: image ENTRYPOINT starts /sbin/starter.d automatically.
    # shellcheck disable=SC2086
    podman run -d \
        --name "$CONTAINER_NAME" \
        --privileged \
        --device /dev/net/tun \
        $PUBLISH_ARGS \
        "$IMAGE"
fi

write_state

echo "node-started logical=$LOGICAL_NODE_ID factory=$FACTORY_NODE_ID container=$CONTAINER_NAME image=$IMAGE"
