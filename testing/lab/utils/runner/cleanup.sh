#!/usr/bin/env bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

# Terminate every non-terminated EC2 worker tagged with one batch ID.

set -Eeuo pipefail

SCRIPT_DIR="$(CDPATH= cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=lib.sh
. "$SCRIPT_DIR/lib.sh"
p0_load_config "$SCRIPT_DIR"

BATCH_ID="${1:-}"
[[ -n "$BATCH_ID" ]] || p0_die "usage: $0 BATCH_ID"
p0_validate_simple_id BATCH_ID "$BATCH_ID"
p0_require_cmd aws
p0_cleanup_batch "$BATCH_ID"
