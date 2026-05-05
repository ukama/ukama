#!/usr/bin/env bash

# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

set -euo pipefail

ROOT="/tmp/sys"
VERBOSE=1

usage() {
    cat <<EOF
Usage: $0 [--root <path>] [--clean]

Options:
  --root <path>        Root directory for mock tree (default: /tmp/sys)
  --clean              Remove root directory and exit
EOF
}

log() {
    [ "$VERBOSE" -eq 1 ] && echo "$@"
}

mkdirp() {
    mkdir -p "$1"
}

mkfile() {
    local path="$1"
    local content="${2:-0}"

    mkdir -p "$(dirname "$path")"
    printf "%s\n" "$content" > "$path"
}

mklink() {
    local target="$1"
    local link="$2"

    mkdir -p "$(dirname "$link")"
    ln -snf "$target" "$link"
}

mk_i2c_bus() {
    local bus="$1"

    mkdirp "$ROOT/bus/i2c/devices/i2c-${bus}"
}

mk_i2c_dev() {
    local bus="$1"
    local addrHex="$2"
    local addrDec=$((addrHex))
    local addr4
    local devdir

    addr4="$(printf "%04x" "$addrDec")"
    devdir="$ROOT/bus/i2c/devices/i2c-${bus}/${bus}-${addr4}"

    mkdirp "$devdir"
    mkfile "$devdir/name" ""
    mkfile "$devdir/uevent" ""

    echo "$devdir"
}

mk_eeprom_file() {
    local devdir="$1"

    mkdirp "$devdir"

    # Needs to be large enough for schema payload offsets.
    dd if=/dev/zero of="$devdir/eeprom" bs=1 count=65536 status=none
}

mk_hwmon() {
    local base="$ROOT/class/hwmon/hwmon0"

    mkdirp "$base"
    mkfile "$base/name" "cnode-cm4-mock"
    mkfile "$base/temp1_input" "42000"
}

mk_dev_i2c_endpoints() {
    mkdirp "$ROOT/dev"
    : > "$ROOT/dev/i2c-0"
}

init_cnode_tree() {
    local inv

    mkdirp "$ROOT/bus/i2c/devices"
    mkdirp "$ROOT/class/hwmon"
    mkdirp "$ROOT/dev"

    mk_i2c_bus 0

    inv="$(mk_i2c_dev 0 0x50)"
    mkfile "$inv/name" "cm4-inventory-eeprom"
    mk_eeprom_file "$inv"

    mklink "$ROOT/bus/i2c/devices/i2c-0/0-0050/eeprom" \
           "$ROOT/inventory_db"

    mklink "$ROOT/bus/i2c/devices/i2c-0/0-0050/eeprom" \
           "$ROOT/cnode_inventory_db"

    #
    # Keep this for compatibility with current node.d default:
    #   #define INVENTORY_DB "/tmp/sys/tnode_inventory_db"
    #
    mklink "$ROOT/bus/i2c/devices/i2c-0/0-0050/eeprom" \
           "$ROOT/tnode_inventory_db"

    mk_hwmon
    mk_dev_i2c_endpoints

    log "CNode mock created at: $ROOT"
}

ACTION="init"

while [ $# -gt 0 ]; do
    case "$1" in
        --root)
            ROOT="$2"
            shift 2
            ;;
        --clean)
            ACTION="clean"
            shift
            ;;
        -h|--help)
            usage
            exit 0
            ;;
        *)
            echo "Unknown arg: $1" >&2
            usage
            exit 1
            ;;
    esac
done

case "$ACTION" in
    clean)
        rm -rf "$ROOT"
        echo "Cleaned: $ROOT"
        ;;
    init)
        rm -rf "$ROOT"
        init_cnode_tree
        ;;
    *)
        echo "Unknown action: $ACTION" >&2
        exit 1
        ;;
esac
