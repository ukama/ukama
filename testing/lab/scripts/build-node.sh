#!/bin/sh
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

set -eu

if [ "$#" -lt 2 ]; then
    echo "usage: $0 <repo> <node-id> [runtime]" >&2
    exit 2
fi

REPO="$1"
NODE_ID="$2"
NODE_RUNTIME="${3:-${NODE_RUNTIME:-starter}}"

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

NODE_DIR="$(find_node_dir)"

echo "build-node: $NODE_ID runtime=$NODE_RUNTIME"
(
    cd "$NODE_DIR"
    ./mk_local_vnode.sh --node-id "$NODE_ID" --runtime "$NODE_RUNTIME"
)
