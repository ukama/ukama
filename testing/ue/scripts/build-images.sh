#!/usr/bin/env bash
#
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
UE_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

UE_IMAGE="${UE_IMAGE:-ukama/ue:dev}"
MEDIA_IMAGE="${MEDIA_IMAGE:-ukama/media:dev}"
BUILD_MEDIA="${BUILD_MEDIA:-false}"

usage() {
    cat <<USAGE
Usage: $0 [--media]

Builds the UE image by default.

Options:
  --media     Also build the local media image for lab-only use.

Environment:
  UE_IMAGE      UE image tag. Default: ukama/ue:dev
  MEDIA_IMAGE   Media image tag. Default: ukama/media:dev
  BUILD_MEDIA   true/false. Default: false
USAGE
}

while [[ $# -gt 0 ]]; do
    case "$1" in
        --media)
            BUILD_MEDIA="true"
            shift
            ;;
        -h|--help)
            usage
            exit 0
            ;;
        *)
            echo "unknown arg $1" >&2
            usage >&2
            exit 1
            ;;
    esac
done

if ! command -v podman >/dev/null 2>&1; then
    echo "podman is required" >&2
    exit 1
fi

podman build \
    -t "$UE_IMAGE" \
    -f "$UE_ROOT/ue/Containerfile" \
    "$UE_ROOT"

if [[ "$BUILD_MEDIA" == "true" ]]; then
    echo "building media image for lab-only use"
    podman build \
        -t "$MEDIA_IMAGE" \
        -f "$UE_ROOT/media/Containerfile" \
        "$UE_ROOT"
else
    echo "skipping media image; final E2E expects MEDIA_IP to be external"
fi
