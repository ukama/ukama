#!/bin/sh
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

set -eu

if [ "$#" -ne 2 ]; then
    echo "usage: $0 <payment|software> <on|off>" >&2
    exit 2
fi

TARGET="$1"
STATE="$2"

case "$TARGET:$STATE" in
    payment:on)
        VAR=ULAB_PAYMENT_FAILURE_ON_CMD
        ;;
    payment:off)
        VAR=ULAB_PAYMENT_FAILURE_OFF_CMD
        ;;
    software:on)
        VAR=ULAB_SOFTWARE_FAILURE_ON_CMD
        ;;
    software:off)
        VAR=ULAB_SOFTWARE_FAILURE_OFF_CMD
        ;;
    *)
        echo "unsupported test control: target=$TARGET state=$STATE" >&2
        exit 2
        ;;
esac

# Resolve the allow-listed environment variable without eval. The command
# itself is deliberately run through a shell because test environments may
# need kubectl, podman, compose, or another deployment-specific control.
CMD="$(printenv "$VAR" 2>/dev/null || true)"
if [ -z "$CMD" ]; then
    echo "test control command is not configured: $VAR" >&2
    exit 1
fi

printf 'test-control target=%s state=%s command_var=%s\n' \
    "$TARGET" "$STATE" "$VAR"
/bin/sh -c "$CMD"
printf 'test-control-complete target=%s state=%s\n' "$TARGET" "$STATE"
