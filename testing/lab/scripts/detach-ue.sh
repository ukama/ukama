#!/bin/sh
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

set -eu

if [ "$#" -ne 2 ]; then
    echo "usage: $0 <ue-id-or-ref> <run-dir>" >&2
    exit 2
fi

SCRIPT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)"
STATE_FILE="$2/runtime-ues/$(printf '%s' "$1" | \
    tr -c 'A-Za-z0-9_.-' '-').env"
# A UE which was never started has no session to close.
[ -f "$STATE_FILE" ] || exit 0
"$SCRIPT_DIR/ue-state.py" "$1" "$2" detach "${ULAB_UE_DETACH_TIMEOUT:-60}"
printf 'DETACHED_AT=%s\n' "$(date +%s)" >> "$STATE_FILE"
