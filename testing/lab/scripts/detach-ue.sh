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
    echo "detach-ue: state not found $STATE_FILE"
    exit 0
fi

# shellcheck disable=SC1090
. "$STATE_FILE"

if [ -z "${TNODE_CONTAINER:-}" ] || [ -z "${IMSI:-}" ]; then
    echo "detach-ue: TNODE_CONTAINER or IMSI missing in $STATE_FILE"
    exit 0
fi

if ! podman inspect -f '{{.State.Running}}' "$TNODE_CONTAINER" 2>/dev/null | grep -q '^true$'; then
    echo "detach-ue: tower container is not running: $TNODE_CONTAINER" >&2
    exit 1
fi

echo "detach-ue: detach imsi=$IMSI tower=$TNODE_CONTAINER"

BODY="$RUN_DIR/runtime-ues/$UE_KEY.detach.body"
CODE="$RUN_DIR/runtime-ues/$UE_KEY.detach.code"
TMP="$BODY.tmp"
mkdir -p "$RUN_DIR/runtime-ues"

HTTP_CODE="$(podman exec "$TNODE_CONTAINER" sh -lc \
    "curl -sS --max-time 8 -o /tmp/ulab-detach-$IMSI.out -w '%{http_code}' -X DELETE 'http://127.0.0.1:18028/v1/ue/$IMSI'; rc=\$?; cat /tmp/ulab-detach-$IMSI.out 2>/dev/null > /tmp/ulab-detach-$IMSI.body; rm -f /tmp/ulab-detach-$IMSI.out; if [ \$rc -ne 0 ]; then echo CURLERR; fi" \
    2>"$TMP.err" || true)"

# The command above returns the HTTP code on stdout. Keep the body separately
# by doing a second best-effort read when the first command did not cleanly
# expose it to the host.  The body is diagnostic only; the code controls flow.
podman exec "$TNODE_CONTAINER" sh -lc \
    "cat /tmp/ulab-detach-$IMSI.body 2>/dev/null; rm -f /tmp/ulab-detach-$IMSI.body" \
    > "$TMP" 2>/dev/null || true
mv "$TMP" "$BODY"
printf "%s\n" "$HTTP_CODE" > "$CODE"

case "$HTTP_CODE" in
    200|202|204|404)
        echo "detach-ue: imsi=$IMSI status=$HTTP_CODE"
        ;;
    *)
        echo "detach-ue: imsi=$IMSI unexpected status=$HTTP_CODE" >&2
        cat "$BODY" >&2 || true
        cat "$TMP.err" >&2 || true
        exit 1
        ;;
esac

# Mark the intended session-close point.  cleanup-ue.sh will remove the
# container after CDR diagnostics and usage validation have completed.
printf "DETACHED_AT=%s\n" "$(date +%s)" >> "$STATE_FILE"

exit 0
