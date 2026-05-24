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
ALLOW_LOCAL_MEDIA="${ALLOW_LOCAL_MEDIA:-false}"

usage() {
    cat <<USAGE
Usage: $0

This script is disabled by default because the E2E target expects media to be
an external Internet service, not a local container on the UE/virtualnode host.

For temporary lab-only testing:
  ALLOW_LOCAL_MEDIA=true $0
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
podman run -d --name "$MEDIA_NAME" --network host "$MEDIA_IMAGE"
echo "$MEDIA_NAME"
