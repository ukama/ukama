#!/bin/sh
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

set -eu

if [ "$#" -ne 2 ]; then
    echo "usage: $0 <repo> <run-dir>" >&2
    exit 2
fi

REPO="$1"
RUN_DIR="$2"
UE_DIR="$REPO/testing/ue"
STATE_DIR="$RUN_DIR/runtime-media"
STATE_FILE="$STATE_DIR/media.env"
MEDIA_IMAGE="ukama/media:dev"
MEDIA_CONTAINER="ukama-media-$(basename "$RUN_DIR" | tr -c 'A-Za-z0-9_.-' '-')"

mkdir -p "$STATE_DIR"

if ! command -v podman >/dev/null 2>&1; then
    echo "podman is required" >&2
    exit 1
fi

if [ ! -f "$UE_DIR/media/Containerfile" ]; then
    echo "missing $UE_DIR/media/Containerfile" >&2
    exit 1
fi

echo "media: build $MEDIA_IMAGE"
podman build -t "$MEDIA_IMAGE" -f "$UE_DIR/media/Containerfile" "$UE_DIR"

echo "media: start $MEDIA_CONTAINER"
podman rm -f "$MEDIA_CONTAINER" >/dev/null 2>&1 || true
podman run -d --name "$MEDIA_CONTAINER" "$MEDIA_IMAGE" >/dev/null

MEDIA_IP="$(podman inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' "$MEDIA_CONTAINER")"
if [ -z "$MEDIA_IP" ]; then
    echo "media container has no IP: $MEDIA_CONTAINER" >&2
    exit 1
fi

cat > "$STATE_FILE" <<STATE
MEDIA_CONTAINER=$MEDIA_CONTAINER
MEDIA_IP=$MEDIA_IP
HTTP_PORT=8080
IPERF_PORT=5201
STATE

echo "media-ready container=$MEDIA_CONTAINER ip=$MEDIA_IP"
