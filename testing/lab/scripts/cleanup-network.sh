#!/bin/sh
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

set -u

# EC2 systemd workers do not necessarily receive HOME.
: "${HOME:=/root}"
export HOME

if [ "$#" -ne 1 ]; then
    echo "usage: $0 <run-dir>" >&2
    exit 2
fi

RUN_DIR="$1"
SCRIPT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)"
HOST_CONTROL="$SCRIPT_DIR/node-host-control.sh"
STATE_FILE="$RUN_DIR/runtime-net/net.env"
FAILED=0

# Always run the mesh sweep, including partial-start and missing-network cases.
cleanup_mesh() {
    result=$?
    if ! python3 "$SCRIPT_DIR/cleanup-mesh.py" "$RUN_DIR"; then
        result=1
    fi
    trap - EXIT
    exit "$result"
}
trap cleanup_mesh EXIT

if [ ! -f "$STATE_FILE" ]; then
    echo "cleanup-network: state not found $STATE_FILE"
    exit 0
fi

# shellcheck disable=SC1090
. "$STATE_FILE"

if [ -z "${LAB_NET:-}" ]; then
    echo "cleanup-network: LAB_NET missing in $STATE_FILE" >&2
    exit 1
fi

echo "cleanup-network: rm containers on $LAB_NET"
if ! containers="$(podman ps -a --filter "network=$LAB_NET" \
    --format '{{.Names}}')"; then
    echo "cleanup-network: failed to list containers" >&2
    FAILED=1
fi
for c in $containers; do
    "$HOST_CONTROL" remove "$c" || FAILED=1
    echo "cleanup-network: rm container $c"
    podman rm -f "$c" >/dev/null 2>&1 || FAILED=1
done

echo "cleanup-network: rm $LAB_NET"
if podman network exists "$LAB_NET"; then
    podman network rm "$LAB_NET" >/dev/null 2>&1 || FAILED=1
fi

if [ -n "${CNI_CONFIG:-}" ]; then
    rm -f "$CNI_CONFIG" 2>/dev/null || true
fi

rm -f "$HOME/.config/cni/net.d/$LAB_NET.conflist" 2>/dev/null || true

exit "$FAILED"
