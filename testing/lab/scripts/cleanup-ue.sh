#!/bin/sh
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

set -u

if [ "$#" -ne 2 ]; then
    echo "usage: $0 <ue-id-or-ref> <run-dir>" >&2
    exit 2
fi

UE_KEY="$1"
RUN_DIR="$2"
STATE_FILE="$RUN_DIR/runtime-ues/$(printf "%s" "$UE_KEY" | tr -c 'A-Za-z0-9_.-' '-').env"

if [ ! -f "$STATE_FILE" ]; then
    echo "cleanup-ue: state not found $STATE_FILE"
    exit 0
fi

# shellcheck disable=SC1090
. "$STATE_FILE"

if [ -n "${UE_CONTAINER:-}" ]; then
    mkdir -p "$RUN_DIR/runtime-ues"
    podman logs "$UE_CONTAINER" > "$RUN_DIR/runtime-ues/$UE_KEY.container.log" 2>&1 || true
    echo "cleanup-ue: rm $UE_CONTAINER"
    podman rm -f "$UE_CONTAINER" >/dev/null 2>&1 || true
fi

exit 0
