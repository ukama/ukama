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

podman rm -f "$MEDIA_NAME" >/dev/null 2>&1 || true
podman run -d --name "$MEDIA_NAME" --network host "$MEDIA_IMAGE"
echo "$MEDIA_NAME"
