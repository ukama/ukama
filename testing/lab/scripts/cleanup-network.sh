#!/bin/sh
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

set -eu

if [ "$#" -lt 1 ]; then
    echo "usage: $0 <run-dir>" >&2
    exit 2
fi

RUN_DIR="$1"
ENV_FILE="$RUN_DIR/runtime.env"

network_name=""

if [ -f "$ENV_FILE" ]; then
    # shellcheck disable=SC1090
    . "$ENV_FILE"
    network_name="${LAB_NET:-${ULAB_NETWORK:-${NETWORK_NAME:-}}}"
fi

if [ -z "$network_name" ]; then
    run_name="$(basename "$RUN_DIR")"
    network_name="ukama-lab-$run_name"
fi

echo "cleanup-network: $network_name"

# Remove any leftover containers still attached to this run network.
containers="$(podman ps -a --filter "network=$network_name" --format '{{.Names}}' 2>/dev/null || true)"
for c in $containers; do
    echo "cleanup-network-container: $c"
    podman rm -f "$c" >/dev/null 2>&1 || true
done

podman network rm "$network_name" >/dev/null 2>&1 || true

exit 0
