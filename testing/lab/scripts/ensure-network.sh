#!/bin/sh
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

set -eu

if [ "$#" -ne 1 ]; then
    echo "usage: $0 <run-dir>" >&2
    exit 2
fi

RUN_DIR="$1"
RUN_ID="$(basename "$RUN_DIR")"
SAFE_RUN_ID="$(printf "%s" "$RUN_ID" | tr -c 'A-Za-z0-9-' '-' | sed 's/^-*//;s/-*$//')"
LAB_NET="ukama-lab-$SAFE_RUN_ID"
STATE_DIR="$RUN_DIR/runtime-net"
STATE_FILE="$STATE_DIR/net.env"

need_cmd() {
    if ! command -v "$1" >/dev/null 2>&1; then
        echo "missing required command: $1" >&2
        exit 1
    fi
}

need_cmd podman
mkdir -p "$STATE_DIR"

if podman network exists "$LAB_NET" >/dev/null 2>&1; then
    echo "network: exists $LAB_NET"
else
    echo "network: create $LAB_NET"
    podman network create "$LAB_NET" >/dev/null
fi

cat > "$STATE_FILE" <<STATE
LAB_NET=$LAB_NET
STATE

echo "network-ready name=$LAB_NET"
