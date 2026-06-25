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

if [ -n "${LAB_NET:-}" ]; then
    echo "cleanup-network: rm $LAB_NET"
    podman network rm "$LAB_NET" >/dev/null 2>&1 || true
fi

if [ -n "${CNI_CONFIG:-}" ]; then
    rm -f "$CNI_CONFIG" 2>/dev/null || true
fi

exit 0
