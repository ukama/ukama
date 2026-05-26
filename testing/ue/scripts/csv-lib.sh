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
        function clean(s) {
            gsub(/\r/, "", s)
            return s
        }

        NR == 1 {
            for (i = 1; i <= NF; i++) {
                h[clean($i)] = i
            }
            next
        }

        {
            line = clean($0)
            if (clean($h["IMSI"]) == imsi) {
                print line
                exit
            }
        }
    ' "$csv"
}

csv_row_by_index() {
    local csv="$1"
    local idx="$2"

    awk -F',' -v row="$idx" '
        function clean(s) {
            gsub(/\r/, "", s)
            return s
        }

        NR == 1 {
            next
        }

        NR == row + 1 {
            print clean($0)
            exit
        }
    ' "$csv"
}

csv_enabled_imsi_list() {
    local csv="$1"

    awk -F',' '
        function clean(s) {
            gsub(/\r/, "", s)
            return s
        }

        NR == 1 {
            for (i = 1; i <= NF; i++) {
                h[clean($i)] = i
            }
            next
        }

        {
            enabled = clean($h["Enabled"])
            if (enabled == "TRUE" || enabled == "true") {
                print clean($h["IMSI"])
            }
        }
    ' "$csv"
}

csv_field() {
    local csv="$1"
    local row="$2"
    local field="$3"

    printf '%s\n' "$row" | awk -F',' -v field="$field" '
        function clean(s) {
            gsub(/\r/, "", s)
            return s
        }

        BEGIN {
            n = split(field, f, "\034")
        }

        {
            for (i = 1; i <= NF; i++) {
                value[i] = clean($i)
            }
        }

        END {
            print value[1]
        }
    '
}

csv_field_by_name() {
    local csv="$1"
    local row="$2"
    local field="$3"

    awk -F',' -v row="$row" -v field="$field" '
        function clean(s) {
            gsub(/\r/, "", s)
            return s
        }

        NR == 1 {
            for (i = 1; i <= NF; i++) {
                h[clean($i)] = i
            }
            next
        }

        clean($0) == row {
            print clean($h[field])
            exit
        }
    ' "$csv"
}

csv_count_enabled() {
    local csv="$1"

    awk -F',' '
        function clean(s) {
            gsub(/\r/, "", s)
            return s
        }

        NR == 1 {
            for (i = 1; i <= NF; i++) {
                h[clean($i)] = i
            }
            next
        }

        {
            enabled = clean($h["Enabled"])
            if (enabled == "TRUE" || enabled == "true") {
                c++
            }
        }

        END {
            print c + 0
        }
    ' "$csv"
}
