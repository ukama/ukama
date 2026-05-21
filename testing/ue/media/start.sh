#!/usr/bin/env bash
#
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

set -euo pipefail

mkdir -p /srv/media
echo "ukama media test target" >/srv/media/index.html
iperf3 -s -D
cd /srv/media
exec python3 -m http.server 8080
