#!/usr/bin/env bash
#
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

set -euo pipefail

csv_row_by_imsi() {
    local csv="$1"
    local imsi="$2"

    awk -F',' -v imsi="$imsi" '
        NR==1 {
            for (i = 1; i <= NF; i++) {
                gsub(//, "", $i)
                h[$i] = i
            }
            next
        }
        {
            for (i = 1; i <= NF; i++) gsub(//, "", $i)
        }
        $h["IMSI"] == imsi {print; exit}
    ' "$csv"
}

csv_row_by_index() {
    local csv="$1"
    local idx="$2"

    awk -F',' -v row="$idx" '
        NR==1 {next}
        {
            for (i = 1; i <= NF; i++) gsub(//, "", $i)
        }
        NR == row + 1 {print; exit}
    ' "$csv"
}

csv_enabled_imsi_list() {
    local csv="$1"

    awk -F',' '
        NR==1 {
            for (i = 1; i <= NF; i++) {
                gsub(//, "", $i)
                h[$i] = i
            }
            next
        }
        {
            for (i = 1; i <= NF; i++) gsub(//, "", $i)
        }
        $h["Enabled"] == "TRUE" || $h["Enabled"] == "true" {print $h["IMSI"]}
    ' "$csv"
}

csv_field() {
    local csv="$1"
    local row="$2"
    local field="$3"

    awk -F',' -v row="$row" -v field="$field" '
        NR==1 {
            for (i = 1; i <= NF; i++) {
                gsub(//, "", $i)
                h[$i] = i
            }
            next
        }
        {
            line = $0
            gsub(//, "", line)
            for (i = 1; i <= NF; i++) gsub(//, "", $i)
        }
        line == row {print $h[field]; exit}
    ' "$csv"
}

csv_count_enabled() {
    local csv="$1"

    awk -F',' '
        NR==1 {
            for (i = 1; i <= NF; i++) {
                gsub(//, "", $i)
                h[$i] = i
            }
            next
        }
        {
            for (i = 1; i <= NF; i++) gsub(//, "", $i)
        }
        $h["Enabled"] == "TRUE" || $h["Enabled"] == "true" {c++}
        END {print c + 0}
    ' "$csv"
}
