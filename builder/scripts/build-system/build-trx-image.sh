#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2025-present, Ukama Inc.

set -euo pipefail

DIR="$(pwd)"
UKAMA_ROOT=$(realpath ../../../)
MISC_FILES_DIR="${UKAMA_ROOT}/builder/scripts/build-system/trx"
OUT_DIR="${OUT_DIR:-${DIR}/trx-image}"
MODE=""
BAND=""
EARFCN=""
BAND_CFG=""
TMP_CFG=""

FDD_UL_OFFSET=18000
UL_BW_20MHZ=5
declare -A EARFCN_MIN=([1]=0 [41]=39650 [48]=55240)
declare -A EARFCN_MAX=([1]=599 [41]=41589 [48]=56739)

log() {
    local type="$1" msg="$2"
    local color
    case "$type" in
        INFO)    color="\033[1;34m";;
        SUCCESS) color="\033[1;32m";;
        ERROR)   color="\033[1;31m";;
        *)       color="\033[1;37m";;
    esac
    echo -e "${color}${type}: ${msg}\033[0m"
}

cleanup() {
    if [[ -n "$TMP_CFG" ]]; then
        rm -f "$TMP_CFG"
    fi
}

trap cleanup EXIT

check_pre_requisite() {
    for cmd in awk sed sha256sum; do
        command -v "$cmd" >/dev/null || { log ERROR "'$cmd' not found in PATH"; exit 1; }
    done

    [[ -f "$BAND_CFG" ]] || { log ERROR "No ${MODE} config for band ${BAND}: ${BAND_CFG}"; exit 1; }

    if [[ -n "$EARFCN" ]]; then
        [[ "$EARFCN" =~ ^[0-9]+$ ]] || { log ERROR "EARFCN '${EARFCN}' is not a number"; exit 1; }
    fi
}

cfg_value() {
    awk -v k="$1" '{for (i = 1; i < NF; i++) if ($i == k) {print $NF; exit}}' "$2"
}

ul_earfcn() {
    if [[ "$MODE" == "FDD" ]]; then
        echo $(($1 + FDD_UL_OFFSET))
    else
        echo "$1"
    fi
}

make_cfg() {
    mkdir -p "$OUT_DIR"
    rm -f "${OUT_DIR}/startup.cfg" "${OUT_DIR}/SHA256SUMS"
    TMP_CFG=$(mktemp)

    if [[ -z "$EARFCN" ]]; then
        cp "$BAND_CFG" "$TMP_CFG"
        return
    fi

    local old_dl ul
    old_dl=$(cfg_value u16dlcarrierfreq "$BAND_CFG")
    ul=$(ul_earfcn "$EARFCN")
    sed -E \
        -e "s/^(configure cell u16dlcarrierfreq 0 )[0-9]+/\1${EARFCN}/" \
        -e "s/^(configure cell freq u16ulcarrierfreq 0 )[0-9]+/\1${ul}/" \
        -e "s/^(configure cell eutranfreqrelation u16arfcnvalueeutrandl 0 [0-9]+ )${old_dl}([[:space:]]*)$/\1${EARFCN}\2/" \
        "$BAND_CFG" > "$TMP_CFG"
}

validate_cfg() {
    local band dl ul ul_bw
    band=$(cfg_value u8freqbandindicator "$TMP_CFG")
    dl=$(cfg_value u16dlcarrierfreq "$TMP_CFG")
    ul=$(cfg_value u16ulcarrierfreq "$TMP_CFG")
    ul_bw=$(cfg_value eulbandwidth "$TMP_CFG")

    [[ "$band" == "$BAND" ]] || { log ERROR "Config is band ${band}, expected ${BAND}"; exit 1; }
    (( dl >= EARFCN_MIN[$BAND] && dl <= EARFCN_MAX[$BAND] )) || {
        log ERROR "EARFCN ${dl} is outside band ${BAND} (${EARFCN_MIN[$BAND]}-${EARFCN_MAX[$BAND]})"
        exit 1
    }
    [[ "$ul" == "$(ul_earfcn "$dl")" ]] || { log ERROR "UL EARFCN ${ul} does not match DL ${dl}"; exit 1; }
    [[ "$ul_bw" == "$UL_BW_20MHZ" ]] || { log ERROR "Config is not 20 MHz"; exit 1; }

    mv "$TMP_CFG" "${OUT_DIR}/startup.cfg"
    chmod 644 "${OUT_DIR}/startup.cfg"
    (cd "$OUT_DIR" && sha256sum startup.cfg > SHA256SUMS)
    log SUCCESS "Band ${BAND}, EARFCN DL ${dl} UL ${ul}, 20 MHz"
}

MODE="${1:-FDD}"
if [[ "$MODE" != "FDD" && "$MODE" != "TDD" ]]; then
    log ERROR "Usage: $0 [FDD|TDD] [band] [earfcn]"
    exit 1
fi

if [[ "$MODE" == "FDD" ]]; then
    BAND="${2:-1}"
else
    BAND="${2:-41}"
fi
EARFCN="${3:-}"
BAND_CFG="${MISC_FILES_DIR}/${MODE,,}/startup_B${BAND}.cfg"

check_pre_requisite
make_cfg
validate_cfg

log SUCCESS "TRX ${MODE} config written to ${OUT_DIR}/startup.cfg"
exit 0
