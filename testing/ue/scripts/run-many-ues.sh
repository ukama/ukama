#!/usr/bin/env bash
#
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/csv-lib.sh"

CSV=""
COUNT=""

while [[ $# -gt 0 ]]; do
    case "$1" in
        --csv) CSV="$2"; shift 2 ;;
        --count) COUNT="$2"; shift 2 ;;
        *) echo "unknown arg $1" >&2; exit 1 ;;
    esac
done

: "${CSV:?--csv required}"
: "${COUNT:?--count required}"

started=0
while IFS= read -r imsi; do
    [[ -z "$imsi" ]] && continue
    "$SCRIPT_DIR/run-ue.sh" --csv "$CSV" --imsi "$imsi"
    started=$((started + 1))
    [[ "$started" -ge "$COUNT" ]] && break
done < <(csv_enabled_imsi_list "$CSV")

echo "started $started UE(s)"
