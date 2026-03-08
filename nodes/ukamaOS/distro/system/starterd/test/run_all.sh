#!/usr/bin/env bash
#
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.
#
set -euo pipefail

SCRIPT_DIR="$(CDPATH='' cd -- "$(dirname -- "$0")" && pwd)"

: "${STARTERD_BIN:=starter.d}"
: "${INIT_STARTER_BIN:=init-starter}"

"$SCRIPT_DIR/smoke_example_app.sh" "$STARTERD_BIN"
"$SCRIPT_DIR/smoke_init_switch.sh" "$INIT_STARTER_BIN"

echo "run_all: PASS"
