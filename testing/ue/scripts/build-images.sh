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

if ! command -v podman >/dev/null 2>&1; then
    echo "podman is required" >&2
    exit 1
fi

podman build \
    -t "$UE_IMAGE" \
    -f "$UE_ROOT/ue/Containerfile" \
    "$UE_ROOT"

podman build \
    -t "$MEDIA_IMAGE" \
    -f "$UE_ROOT/media/Containerfile" \
    "$UE_ROOT"
