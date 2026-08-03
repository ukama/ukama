#!/bin/sh
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

set -eu

if [ "$#" -ne 2 ]; then
    echo "usage: $0 <target> <on|off>" >&2
    exit 2
fi

TARGET="$1"
STATE="$2"

case "$TARGET" in
    payment|software|site_restart|node_restart|software_timeout)
        ;;
    *)
        echo "unsupported test control target: $TARGET" >&2
        exit 2
        ;;
esac

case "$STATE" in
    on|off)
        ;;
    *)
        echo "unsupported test control state: $STATE" >&2
        exit 2
        ;;
esac

TARGET_UPPER="$(printf '%s' "$TARGET" | tr '[:lower:]' '[:upper:]')"
STATE_UPPER="$(printf '%s' "$STATE" | tr '[:lower:]' '[:upper:]')"
VAR="ULAB_${TARGET_UPPER}_FAILURE_${STATE_UPPER}_CMD"

# Resolve only an allow-listed variable name. The command is intentionally
# run through a shell because test deployments may use kubectl, podman,
# compose, or another environment-specific control.
CMD="$(printenv "$VAR" 2>/dev/null || true)"
if [ -z "$CMD" ]; then
    echo "test control command is not configured: $VAR" >&2
    exit 1
fi

printf 'test-control target=%s state=%s command_var=%s\n' \
    "$TARGET" "$STATE" "$VAR"
/bin/sh -c "$CMD"
printf 'test-control-complete target=%s state=%s\n' "$TARGET" "$STATE"
