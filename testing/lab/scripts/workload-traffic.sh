#!/bin/sh
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

# One measured payload transfer through the existing UE TUN and node datapath.
set -eu
[ "$#" -eq 5 ] || { echo 'usage: workload-traffic.sh UE_ID MB RUN_DIR SLOT RESULT_JSON' >&2; exit 2; }
ue_key=$1
amount_mb=$2
run_dir=$3
slot=$4
result=$5
state="$run_dir/runtime-ues/$ue_key.env"
[ -f "$state" ] || { echo "missing UE state: $state" >&2; exit 1; }
. "$state"
: "${UE_CONTAINER:?}" "${TNODE_CONTAINER:?}" "${MEDIA_CONTAINER:?}" "${MEDIA_IP:?}"
podman exec "$UE_CONTAINER" test -d /sys/class/net/tun0
podman exec "$UE_CONTAINER" ip route replace "$MEDIA_IP/32" dev tun0
# Each worker slot gets a server port. The usual single iperf server accepts
# one test at a time and would otherwise become the load generator bottleneck.
port=$((5202 + slot))
podman exec "$MEDIA_CONTAINER" sh -c '
  port=$1
  pidfile=/tmp/ulab-iperf-$port.pid
  if [ -f "$pidfile" ] && kill -0 "$(cat "$pidfile")" 2>/dev/null; then exit 0; fi
  iperf3 -s -D -p "$port" --pidfile "$pidfile"
' sh "$port"
bytes=$((amount_mb * 1024 * 1024))
if ! podman exec "$UE_CONTAINER" iperf3 -c "$MEDIA_IP" -p "$port" -n "$bytes" -J > "$result.tmp"; then
  mv "$result.tmp" "$result"
  exit 1
fi
mv "$result.tmp" "$result"
# The C collector reads the receiver's actual bytes; no requested-byte credit.
jq -e '.error == null and (.end.sum_received.bytes | type == "number") and .end.sum_received.bytes > 0' "$result" >/dev/null
