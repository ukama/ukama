#!/usr/bin/env bash
#
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

set -euo pipefail

MEDIA_IMAGE="${MEDIA_IMAGE:-ukama/media:dev}"
MEDIA_NAME="${MEDIA_NAME:-ukama-media}"
MEDIA_NETWORK_MODE="${MEDIA_NETWORK_MODE:-container}"
VNODE_NAME="${VNODE_NAME:-ukama-vnode}"
ALLOW_LOCAL_MEDIA="${ALLOW_LOCAL_MEDIA:-true}"

usage() {
    cat <<USAGE
Usage: $0

Environment:
  MEDIA_IMAGE          Default: ukama/media:dev
  MEDIA_NAME           Default: ukama-media
  MEDIA_NETWORK_MODE   container|host|podman. Default: container
  VNODE_NAME           Default: ukama-vnode
  ALLOW_LOCAL_MEDIA    Default: true for lab mode

Examples:
  MEDIA_IMAGE=localhost/ukama/media:dev $0
  MEDIA_NETWORK_MODE=host $0
USAGE
}

if [[ "${1:-}" == "-h" || "${1:-}" == "--help" ]]; then
    usage
    exit 0
fi

if [[ "$ALLOW_LOCAL_MEDIA" != "true" ]]; then
    cat >&2 <<MSG
media must be external for the real E2E path.
Set MEDIA_IP to an external media/sink server and do not run local media.

For temporary lab-only testing only:
  ALLOW_LOCAL_MEDIA=true $0
MSG
    exit 1
fi

podman rm -f "$MEDIA_NAME" >/dev/null 2>&1 || true

case "$MEDIA_NETWORK_MODE" in
    container)
        podman run -d \
            --name "$MEDIA_NAME" \
            --network "container:${VNODE_NAME}" \
            "$MEDIA_IMAGE"
        ;;
    host)
        podman run -d \
            --name "$MEDIA_NAME" \
            --network host \
            "$MEDIA_IMAGE"
        ;;
    podman)
        podman run -d \
            --name "$MEDIA_NAME" \
            -p 8080:8080 \
            -p 5201:5201 \
            "$MEDIA_IMAGE"
        ;;
    *)
        echo "unknown MEDIA_NETWORK_MODE=$MEDIA_NETWORK_MODE" >&2
        exit 1
        ;;
esac

echo "$MEDIA_NAME"
