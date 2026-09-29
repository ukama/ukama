#!/bin/sh
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

set -eu
[ "$#" -eq 2 ] || { echo 'usage: workload-prepare-ue.sh REPO RUN_DIR' >&2; exit 2; }
mkdir -p "$2/runtime-ues"
podman build -t ukama/ue:dev -f "$1/testing/ue/ue/Containerfile" "$1/testing/ue"
touch "$2/runtime-ues/.ue-image-built"
