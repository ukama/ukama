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

UE_KEY="$1"
RUN_DIR="$2"
SCRIPT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)"
STATE_FILE="$RUN_DIR/runtime-ues/$(printf "%s" "$UE_KEY" | tr -c 'A-Za-z0-9_.-' '-').env"

if [ ! -f "$STATE_FILE" ]; then
    echo "UE state not found: $STATE_FILE" >&2
    exit 1
fi

# shellcheck disable=SC1090
. "$STATE_FILE"

require_value() {
    name="$1"
    eval "value=\${$name:-}"
    if [ -z "$value" ]; then
        echo "$name missing in $STATE_FILE" >&2
        exit 1
    fi
}

require_value UE_CONTAINER
require_value TNODE_CONTAINER
require_value TOWER_IP
require_value MEDIA_CONTAINER
require_value MEDIA_IP
require_value IMSI
require_value UE_IP

"$SCRIPT_DIR/verify-ue-session.sh" "$UE_KEY" "$RUN_DIR"

if ! podman exec "$UE_CONTAINER" ip route get "$MEDIA_IP" 2>/dev/null | \
    grep -q "dev tun0"; then
    echo "UE media route is not using tun0: ue=$UE_KEY media=$MEDIA_IP" >&2
    exit 1
fi

if ! podman exec "$MEDIA_CONTAINER" ip route get "$UE_IP" 2>/dev/null | \
    grep -q "via $TOWER_IP"; then
    echo "media return route is not using site tower: ue=$UE_KEY via=$TOWER_IP" >&2
    exit 1
fi

FLOW_TMP="/tmp/ulab-flow-$IMSI.$$"
FLOW_RESPONSE="$(podman exec "$TNODE_CONTAINER" sh -lc \
    "code=\$(curl -sS --max-time 5 -o '$FLOW_TMP' -w '%{http_code}' \\
      'http://127.0.0.1:18030/v1/subscriber/imsi/$IMSI/flow' || printf 000); \\
     printf '%s\\n' \"\$code\"; cat '$FLOW_TMP' 2>/dev/null; rm -f '$FLOW_TMP'" \
    2>/dev/null || true)"
FLOW_CODE="$(printf '%s\n' "$FLOW_RESPONSE" | sed -n '1p')"
FLOW_BODY="$(printf '%s\n' "$FLOW_RESPONSE" | sed '1d')"

case "$FLOW_CODE" in
    204|404)
        echo "ue-policy-blocked ue=$UE_KEY imsi=$IMSI reason=pcrf-flow-withdrawn status=$FLOW_CODE"
        exit 0
        ;;
    200)
        ;;
    *)
        echo "PCRF flow query failed: imsi=$IMSI status=${FLOW_CODE:-missing}" >&2
        exit 1
        ;;
esac

COMPACT_BODY="$(printf '%s' "$FLOW_BODY" | tr -d '[:space:]')"
case "$COMPACT_BODY" in
    ''|'[]'|'{}'|'null')
        echo "ue-policy-blocked ue=$UE_KEY imsi=$IMSI reason=pcrf-flow-empty"
        exit 0
        ;;
esac

OVS_FLOWS="$(podman exec "$TNODE_CONTAINER" \
    ovs-ofctl -O OpenFlow15 dump-flows br0 2>/dev/null || true)"

if ! printf '%s\n' "$OVS_FLOWS" | \
    grep -q "priority=100.*nw_src=$UE_IP.*NORMAL"; then
    echo "PCRF declares a flow but TX OVS flow is missing: imsi=$IMSI ip=$UE_IP" >&2
    exit 1
fi

if ! printf '%s\n' "$OVS_FLOWS" | \
    grep -q "priority=100.*nw_dst=$UE_IP.*NORMAL"; then
    echo "PCRF declares a flow but RX OVS flow is missing: imsi=$IMSI ip=$UE_IP" >&2
    exit 1
fi

# With an attached UE, valid routes, reachable media, and both PCRF-backed OVS
# flows present, a failed transfer is attributable to the active policy/meter.
echo "ue-policy-blocked ue=$UE_KEY imsi=$IMSI reason=policy-meter-enforced"
