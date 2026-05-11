#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2025-present, Ukama Inc.

validate_tarball() {
    local path="$1"
    local label="$2"
    local must_contain_csv="$3"
    local min_entries="${4:-100}"

    if [ ! -f "$path" ]; then
        echo "  [FAIL] $label: file not found — $path" >&2
        return 1
    fi

    local listing
    if ! listing=$(tar tzf "$path" 2>/dev/null); then
        echo "  [FAIL] $label: not a readable gzip tarball — $path" >&2
        return 1
    fi

    local entry_count
    entry_count=$(printf '%s\n' "$listing" | grep -c '.')
    if [ "$entry_count" -lt "$min_entries" ]; then
        echo "  [FAIL] $label: only $entry_count entries (need >= $min_entries) — likely empty or truncated" >&2
        return 1
    fi

    if printf '%s\n' "$listing" | grep -qE '^/'; then
        local sample
        sample=$(printf '%s\n' "$listing" | grep -E '^/' | head -3 | tr '\n' ' ')
        echo "  [FAIL] $label: contains absolute paths (e.g. $sample)" >&2
        echo "         re-create with: cd <mounted_dir> && tar czf out.tgz ." >&2
        return 1
    fi

    local top_levels
    top_levels=$(printf '%s\n' "$listing" \
        | sed 's|^\./||' \
        | awk -F/ '{print $1}' \
        | sort -u \
        | grep -v '^$' || true)

    local missing=()
    local entry
    local old_ifs="$IFS"
    IFS=','
    local required=( $must_contain_csv )
    IFS="$old_ifs"
    for entry in "${required[@]}"; do
        entry="${entry// /}"
        [ -z "$entry" ] && continue
        if ! printf '%s\n' "$top_levels" | grep -qx "$entry"; then
            missing+=("$entry")
        fi
    done

    if [ "${#missing[@]}" -gt 0 ]; then
        local found
        found=$(printf '%s\n' "$top_levels" | head -8 | tr '\n' ' ')
        echo "  [FAIL] $label: missing required top-level entries: ${missing[*]}" >&2
        echo "         found at root: $found" >&2
        local non_mount
        non_mount=$(printf '%s\n' "$top_levels" | grep -vE '^(media|mnt|home|run)$' || true)
        if [ -z "$non_mount" ]; then
            echo "         tarball looks like it was created from a mount path, not from inside it" >&2
            echo "         re-create with: cd <mounted_dir> && sudo tar czf out.tgz ." >&2
        fi
        return 1
    fi

    local size
    size=$(du -h "$path" 2>/dev/null | cut -f1)
    echo "  [ OK ] $label: $entry_count entries, ${size:-unknown}, has [$must_contain_csv]"
    return 0
}
