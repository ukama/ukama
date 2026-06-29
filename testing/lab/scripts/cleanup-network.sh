#!/bin/sh
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

set -u

if [ "$#" -ne 1 ]; then
    echo "usage: $0 <run-dir>" >&2
    exit 2
fi

RUN_DIR="$1"
STATE_FILE="$RUN_DIR/runtime-net/net.env"

if [ ! -f "$STATE_FILE" ]; then
    echo "cleanup-network: state not found $STATE_FILE"
    exit 0
fi

# shellcheck disable=SC1090
. "$STATE_FILE"

if [ -z "${LAB_NET:-}" ]; then
    echo "cleanup-network: LAB_NET missing in $STATE_FILE"
    exit 0
fi

echo "cleanup-network: rm containers on $LAB_NET"
containers="$(podman ps -a --filter "network=$LAB_NET" --format '{{.Names}}' 2>/dev/null || true)"
for c in $containers; do
    echo "cleanup-network: rm container $c"
    podman rm -f "$c" >/dev/null 2>&1 || true
done

echo "cleanup-network: rm $LAB_NET"
podman network rm "$LAB_NET" >/dev/null 2>&1 || true

if [ -n "${CNI_CONFIG:-}" ]; then
    rm -f "$CNI_CONFIG" 2>/dev/null || true
fi

rm -f "$HOME/.config/cni/net.d/$LAB_NET.conflist" 2>/dev/null || true

exit 0
