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

dump_policy_debug() {
    echo "---- Policy block failure: ue=$UE_KEY imsi=$IMSI ----" >&2
    echo "flow response used by check: status=${FLOW_CODE:-not-queried}" >&2
    printf '%s\n' "${FLOW_BODY:-}" >&2
    echo "OVS flows used by check:" >&2
    printf '%s\n' "${OVS_FLOWS:-not-queried}" >&2

    for path in /v1/status /v1/service \
        "/v1/subscriber/imsi/$IMSI" \
        "/v1/policy/imsi/$IMSI" \
        "/v1/cdr/imsi/$IMSI"; do
        echo "---- PCRF $path ----" >&2
        timeout -k 2s 10s podman exec "$TNODE_CONTAINER" \
            curl -sS --max-time 5 -w '\nHTTP %{http_code}\n' \
            "http://127.0.0.1:18030$path" >&2 || true
    done

    echo "---- PCRF persisted usage, sessions, and meters ----" >&2
    case "$IMSI" in
        ''|*[!0-9]*)
            echo "Skipping database query: IMSI must contain only digits" >&2
            ;;
        *)
            if timeout -k 2s 10s podman exec "$TNODE_CONTAINER" \
                sh -c 'command -v sqlite3 >/dev/null'; then
                timeout -k 2s 10s podman exec "$TNODE_CONTAINER" \
                    sqlite3 -readonly -header -column /ukama/apps/db/pcrf.db "
                    BEGIN;
                    SELECT s.imsi, s.service_on, hex(p.id) AS policy_id,
                           p.data, p.consumed, p.ulbr, p.dlbr,
                           p.starttime, p.endtime, u.data AS local_usage,
                           u.updatedat AS usage_updatedat
                    FROM subscribers s JOIN policies p ON p.id = s.policy_id
                    LEFT JOIN usages u ON u.subscriber_id = s.id
                    WHERE s.imsi = '$IMSI';
                    SELECT id, hex(policy_id) AS policy_id, starttime, endtime,
                           txbytes, rxbytes, totalbytes, state, flowstate, sync,
                           txmeter_id, rxmeter_id, updatedat
                    FROM sessions WHERE subscriber_id =
                        (SELECT id FROM subscribers WHERE imsi = '$IMSI')
                    ORDER BY id;
                    SELECT f.id, f.cookie, f.ueipaddr, f.meter_id,
                           m.rate, m.burst, m.type
                    FROM flows f JOIN meters m ON m.id = f.meter_id
                    WHERE m.id IN (
                        SELECT rxmeter_id FROM sessions WHERE subscriber_id =
                            (SELECT id FROM subscribers WHERE imsi = '$IMSI')
                        UNION
                        SELECT txmeter_id FROM sessions WHERE subscriber_id =
                            (SELECT id FROM subscribers WHERE imsi = '$IMSI'));
                    COMMIT;" >&2 || true
            else
                echo "sqlite3 unavailable; see PCRF policy and CDR responses above" >&2
            fi
            ;;
    esac

    echo "---- OVS meters at failure ----" >&2
    timeout -k 2s 10s podman exec "$TNODE_CONTAINER" \
        ovs-ofctl -O OpenFlow15 dump-meters br0 >&2 || true
    echo "---- PCRF logs from current boot (all levels) ----" >&2
    timeout -k 2s 15s podman exec "$TNODE_CONTAINER" \
        /sbin/ukama-log --boot current --format jsonl | \
        grep -E '"app"[[:space:]]*:[[:space:]]*"pcrf"' >&2 || true
}

policy_check_exit() {
    check_rc=$?
    trap - 0
    if [ "$check_rc" -ne 0 ]; then
        dump_policy_debug || true
    fi
    exit "$check_rc"
}

trap policy_check_exit 0

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
