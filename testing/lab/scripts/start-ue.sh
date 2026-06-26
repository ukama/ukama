#!/bin/sh
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

set -eu

if [ "$#" -ne 8 ]; then
    echo "usage: $0 <repo> <ue-ref> <ue-id> <imsi> <iccid> <ue-ip> <site-ref> <run-dir>" >&2
    exit 2
fi

REPO="$1"
UE_REF="$2"
UE_ID="$3"
IMSI="$4"
ICCID="$5"
UE_IP="$6"
SITE_REF="$7"
RUN_DIR="$8"

UE_DIR="$REPO/testing/ue"
NET_STATE="$RUN_DIR/runtime-net/net.env"
SITE_STATE="$RUN_DIR/runtime-sites/$(printf "%s" "$SITE_REF" | tr -c 'A-Za-z0-9_.-' '-').env"
MEDIA_STATE="$RUN_DIR/runtime-media/media.env"
UE_STATE_DIR="$RUN_DIR/runtime-ues"

UE_IMAGE="ukama/ue:dev"
UE_CONTAINER="ue-$IMSI"

mkdir -p "$UE_STATE_DIR"

need_cmd() {
    if ! command -v "$1" >/dev/null 2>&1; then
        echo "missing required command: $1" >&2
        exit 1
    fi
}

need_file() {
    if [ ! -f "$1" ]; then
        echo "missing $1" >&2
        exit 1
    fi
}

need_char_device() {
    if [ ! -c "$1" ]; then
        echo "missing character device $1" >&2
        exit 1
    fi
}

container_ip_on_network() {
    container="$1"
    network="$2"

    podman inspect -f '{{range $name, $net := .NetworkSettings.Networks}}{{if eq $name "'"$network"'"}}{{$net.IPAddress}}{{end}}{{end}}' \
        "$container" 2>/dev/null
}

ue_data_port() {
    last3="$(printf "%s" "$IMSI" | sed 's/.*\(...\)$/\1/')"
    awk -v x="$last3" 'BEGIN { printf "%d", 41000 + x }'
}

container_running() {
    podman inspect -f '{{.State.Running}}' "$1" 2>/dev/null | grep -q '^true$'
}

wait_container_http() {
    container="$1"
    name="$2"
    url="$3"
    i=0

    echo "ue: wait $name from $container url=$url"

    while [ "$i" -lt 60 ]; do
        if podman exec "$container" curl -fsS --max-time 2 "$url" >/dev/null 2>&1; then
            echo "ue: $name ready"
            return 0
        fi

        i=$((i + 1))
        sleep 2
    done

    echo "$name not ready from $container: $url" >&2
    return 1
}


wait_pcrf_service_on() {
    container="$1"
    url="$2"
    i=0
    body=""

    echo "ue: wait pcrf service from $container url=$url"

    while [ "$i" -lt 60 ]; do
        body="$(podman exec "$container" curl -fsS --max-time 2 "$url" 2>/dev/null || true)"
        if printf "%s" "$body" | grep -Eq '"state"[[:space:]]*:[[:space:]]*"on"' && \
           printf "%s" "$body" | grep -Eq '"admission"[[:space:]]*:[[:space:]]*"enabled"'; then
            echo "ue: pcrf service ready"
            return 0
        fi

        i=$((i + 1))
        sleep 2
    done

    echo "pcrf service not enabled from $container: $url" >&2
    echo "last service body: $body" >&2
    return 1
}

print_attach_debug() {
    echo "---- UE container ----" >&2
    podman ps -a --filter "name=$UE_CONTAINER" >&2 || true

    echo "---- UE logs ----" >&2
    podman logs --tail 120 "$UE_CONTAINER" >&2 || true
    echo >&2

    echo "---- EPCEMU status from tower ----" >&2
    podman exec "$TNODE_CONTAINER" sh -lc \
        'curl -sS --max-time 5 http://127.0.0.1:18028/v1/status || true' >&2 || true
    echo >&2

    echo "---- PCRF service from tower ----" >&2
    podman exec "$TNODE_CONTAINER" sh -lc \
        'curl -sS --max-time 5 http://127.0.0.1:18030/v1/service || true' >&2 || true
    echo >&2

    echo "---- PCRF status from tower ----" >&2
    podman exec "$TNODE_CONTAINER" sh -lc \
        'curl -sS --max-time 5 http://127.0.0.1:18030/v1/status || true' >&2 || true
    echo >&2
}

need_cmd podman
need_cmd awk
need_cmd curl
need_cmd grep
need_char_device /dev/net/tun

need_file "$UE_DIR/ue/Containerfile"
need_file "$NET_STATE"
need_file "$SITE_STATE"
need_file "$MEDIA_STATE"

# shellcheck disable=SC1090
. "$NET_STATE"
# shellcheck disable=SC1090
. "$SITE_STATE"
# shellcheck disable=SC1090
. "$MEDIA_STATE"

if [ -z "${LAB_NET:-}" ]; then
    echo "LAB_NET missing in $NET_STATE" >&2
    exit 1
fi

if [ -z "${TNODE_CONTAINER:-}" ]; then
    echo "TNODE_CONTAINER missing in $SITE_STATE" >&2
    exit 1
fi

if [ -z "${MEDIA_CONTAINER:-}" ]; then
    echo "MEDIA_CONTAINER missing in $MEDIA_STATE" >&2
    exit 1
fi

if [ -z "${MEDIA_IP:-}" ]; then
    echo "MEDIA_IP missing in $MEDIA_STATE" >&2
    exit 1
fi

HTTP_PORT="${HTTP_PORT:-8080}"
IPERF_PORT="${IPERF_PORT:-5201}"

UE_IP_ADDR="$(printf "%s" "$UE_IP" | cut -d/ -f1)"

if [ -z "$UE_IP_ADDR" ]; then
    echo "UE IP missing for $UE_ID" >&2
    exit 1
fi

if [ -z "$IMSI" ]; then
    echo "IMSI missing for $UE_ID" >&2
    exit 1
fi

if [ -z "$ICCID" ]; then
    echo "ICCID missing for $UE_ID" >&2
    exit 1
fi

if ! podman network exists "$LAB_NET" >/dev/null 2>&1; then
    echo "podman network does not exist: $LAB_NET" >&2
    exit 1
fi

if ! container_running "$TNODE_CONTAINER"; then
    echo "tower container is not running: $TNODE_CONTAINER" >&2
    podman ps -a --filter "name=$TNODE_CONTAINER" >&2 || true
    exit 1
fi

if ! container_running "$MEDIA_CONTAINER"; then
    echo "media container is not running: $MEDIA_CONTAINER" >&2
    podman ps -a --filter "name=$MEDIA_CONTAINER" >&2 || true
    exit 1
fi

TOWER_IP="$(container_ip_on_network "$TNODE_CONTAINER" "$LAB_NET")"
if [ -z "$TOWER_IP" ]; then
    echo "tower container has no IP on $LAB_NET: $TNODE_CONTAINER" >&2
    podman inspect "$TNODE_CONTAINER" >&2 || true
    exit 1
fi

# Do not use host-published ports for the E2E path.  Validate the same network
# path the UE will use by curling from another container already on LAB_NET.
# /v1/ping returns 503 until the service is really ready; /v1/status may return
# 200 while ready=false.
wait_container_http "$MEDIA_CONTAINER" epcemu "http://$TOWER_IP:18028/v1/ping"
wait_container_http "$MEDIA_CONTAINER" pcrf "http://$TOWER_IP:18030/v1/ping"
wait_pcrf_service_on "$MEDIA_CONTAINER" "http://$TOWER_IP:18030/v1/service"
wait_container_http "$MEDIA_CONTAINER" media "http://$MEDIA_IP:$HTTP_PORT/"

if [ ! -f "$UE_STATE_DIR/.ue-image-built" ]; then
    echo "ue: build $UE_IMAGE"
    podman build -t "$UE_IMAGE" -f "$UE_DIR/ue/Containerfile" "$UE_DIR"
    touch "$UE_STATE_DIR/.ue-image-built"
fi

CSV="$UE_STATE_DIR/$UE_ID.csv"
cat > "$CSV" <<CSVEOF
IMSI,ICCID,MSISDN,SmDpAddress,ActivationCode,IsPhysical,QrCode,UE_IP,APN,Enabled
$IMSI,$ICCID,10000000000,lab,lab,TRUE,lab,$UE_IP_ADDR,internet,TRUE
CSVEOF

UE_DATA_PORT="$(ue_data_port)"

echo "ue: start $UE_CONTAINER network=$LAB_NET tower=$TOWER_IP media=$MEDIA_IP ip=$UE_IP_ADDR"
podman rm -f "$UE_CONTAINER" >/dev/null 2>&1 || true

if ! podman run -d \
    --name "$UE_CONTAINER" \
    --privileged \
    --device /dev/net/tun \
    --network "$LAB_NET" \
    -e UE_IMSI="$IMSI" \
    -e UE_ICCID="$ICCID" \
    -e UE_IP="$UE_IP_ADDR/22" \
    -e UE_APN="internet" \
    -e UE_TUN="tun0" \
    -e EPCEMU_URL="http://$TOWER_IP:18028" \
    -e EPCEMU_DATA_HOST="$TOWER_IP" \
    -e EPCEMU_DATA_PORT="18029" \
    -e UE_DATA_HOST="0.0.0.0" \
    -e UE_DATA_PORT="$UE_DATA_PORT" \
    -e PCRF_URL="http://$TOWER_IP:18030" \
    -e MEDIA_IP="$MEDIA_IP" \
    -e HTTP_PORT="$HTTP_PORT" \
    -e IPERF_PORT="$IPERF_PORT" \
    -e UE_DETACH_ON_EXIT="1" \
    "$UE_IMAGE" \
    /bin/sh -c 'exec /opt/ukama/ue-agent/ue-agent' \
    > "$UE_STATE_DIR/$UE_ID.start.out" 2>&1; then
    echo "failed to start UE container: $UE_CONTAINER" >&2
    cat "$UE_STATE_DIR/$UE_ID.start.out" >&2 || true
    exit 1
fi

UE_CONTAINER_IP="$(container_ip_on_network "$UE_CONTAINER" "$LAB_NET")"
if [ -z "$UE_CONTAINER_IP" ]; then
    echo "UE container has no IP on $LAB_NET: $UE_CONTAINER" >&2
    print_attach_debug
    podman inspect "$UE_CONTAINER" >&2 || true
    exit 1
fi

echo "ue: container ip=$UE_CONTAINER_IP network=$LAB_NET"

STATE_FILE="$UE_STATE_DIR/$UE_ID.env"
cat > "$STATE_FILE" <<STATE
UE_REF=$UE_REF
UE_ID=$UE_ID
UE_CONTAINER=$UE_CONTAINER
IMSI=$IMSI
ICCID=$ICCID
UE_IP=$UE_IP_ADDR
UE_DATA_PORT=$UE_DATA_PORT
UE_CONTAINER_IP=$UE_CONTAINER_IP
SITE_REF=$SITE_REF
TNODE_CONTAINER=$TNODE_CONTAINER
TOWER_IP=$TOWER_IP
MEDIA_CONTAINER=$MEDIA_CONTAINER
MEDIA_IP=$MEDIA_IP
HTTP_PORT=$HTTP_PORT
IPERF_PORT=$IPERF_PORT
LAB_NET=$LAB_NET
CSV=$CSV
UE_DIR=$UE_DIR
STATE

cp "$STATE_FILE" "$UE_STATE_DIR/$UE_REF.env"

sleep 2

if ! container_running "$UE_CONTAINER"; then
    echo "UE container exited early: $UE_CONTAINER" >&2
    print_attach_debug
    exit 1
fi

echo "ue-started ref=$UE_REF imsi=$IMSI ip=$UE_IP_ADDR tower=$TOWER_IP media=$MEDIA_IP network=$LAB_NET"
