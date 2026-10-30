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
LTE_NODE_REPO=""
MODE=""
BAND_CFG="${BAND_CFG:-}"
STACK_BUILD_DIR=""
STAGE_DIR=""

UBOOT_PART_SIZE=$((0x220000))
UBOOT_COPY_OFFSET=$((0xC0000))
OS_PART_SIZE=$((0x500000))
RD_PART_SIZE=$((0x2B00000))
RD_MAX_SIZE=$((0x2800000))
APP_PART_SIZE=$((0xEE0000))
ENV_PART_SIZE=$((0x20000))
JFFS2_ERASE_SIZE=0x20000

STACK_FILES=(lsmD.gz LFMSOFT_OCT_D.tgz dsp.tgz Dimark_Client.tgz)

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

check_status() {
    local code=$1 msg=$2 stage=$3
    if [[ $code -ne 0 ]]; then
        log ERROR "Stage '${stage}' failed"
        exit 1
    fi
    log SUCCESS "$msg"
}

cleanup() {
    if [[ -n "$STAGE_DIR" ]]; then
        rm -rf "$STAGE_DIR"
    fi
}

trap cleanup EXIT

check_pre_requisite() {
    for cmd in make gzip mkfs.jffs2 mkenvimage sha256sum; do
        command -v "$cmd" >/dev/null || { log ERROR "'$cmd' not found in PATH"; exit 1; }
    done

    for f in "${LTE_NODE_REPO}/Makefile" "${LTE_NODE_REPO}/distro/env-setup" \
             "${MISC_FILES_DIR}/uboot.env" "${MISC_FILES_DIR}/app/rc_pre.local" \
             "${MISC_FILES_DIR}/app/${BAND_CFG}"; do
        [[ -f "$f" ]] || { log ERROR "Missing ${f}"; exit 1; }
    done
}

build_os() {
    log INFO "Building u-boot, kernel and rootfs from ${LTE_NODE_REPO}"
    make -C "${LTE_NODE_REPO}" os fs
    check_status $? "u-boot, kernel and rootfs built" "make os fs"
}

build_stack() {
    log INFO "Building ${MODE} LTE stack"
    make -C "${LTE_NODE_REPO}" stack TYPE="${MODE}"
    check_status $? "${MODE} stack built" "make stack"

    STACK_BUILD_DIR=$(ls -d "${LTE_NODE_REPO}/stack/${MODE,,}/source/LTE_Stack_${MODE}_Bin_startup_NMMdisable_IPV6disable_"* 2>/dev/null | tail -1)
    [[ -n "$STACK_BUILD_DIR" ]] || { log ERROR "No ${MODE} stack build directory found"; exit 1; }
}

validate_build_artifacts() {
    local distro_out="${LTE_NODE_REPO}/distro/output"
    UBOOT_BIN="${distro_out}/u-boot-octeon_trx.bin"
    OS_GZ="${distro_out}/lsm_os.gz"
    RD_GZ="${distro_out}/lsm_rd.gz"

    local missing=()
    for f in "$UBOOT_BIN" "$OS_GZ" "$RD_GZ"; do
        [[ -f "$f" ]] || missing+=("$f")
    done
    for f in "${STACK_FILES[@]}"; do
        [[ -f "${STACK_BUILD_DIR}/${f}" ]] || missing+=("${STACK_BUILD_DIR}/${f}")
    done
    if [[ ${#missing[@]} -ne 0 ]]; then
        log ERROR "Missing build artifacts: ${missing[*]}"
        exit 1
    fi

    (( $(stat -c %s "$UBOOT_BIN") <= UBOOT_COPY_OFFSET )) || { log ERROR "u-boot larger than $UBOOT_COPY_OFFSET bytes"; exit 1; }
    (( $(stat -c %s "$OS_GZ") <= OS_PART_SIZE )) || { log ERROR "lsm_os.gz larger than $OS_PART_SIZE bytes"; exit 1; }
    (( $(stat -c %s "$RD_GZ") <= RD_MAX_SIZE )) || { log ERROR "lsm_rd.gz larger than $RD_MAX_SIZE bytes"; exit 1; }

    log SUCCESS "Build artifacts present and within partition sizes"
}

stage_app() {
    STAGE_DIR=$(mktemp -d)
    cp -a "${MISC_FILES_DIR}/app/." "${STAGE_DIR}/"
    for f in "${STACK_FILES[@]}"; do
        cp "${STACK_BUILD_DIR}/${f}" "${STAGE_DIR}/"
    done
    cp "${MISC_FILES_DIR}/app/${BAND_CFG}" "${STAGE_DIR}/startup.cfg"
    touch "${STAGE_DIR}/.sign1985"
    log SUCCESS "App partition staged with ${BAND_CFG} as startup.cfg"
}

blank_image() {
    head -c "$2" /dev/zero | tr '\000' '\377' > "$1"
}

pad_image() {
    blank_image "$2" "$3"
    dd if="$1" of="$2" conv=notrunc status=none
}

make_flash_images() {
    mkdir -p "$OUT_DIR"
    rm -f "${OUT_DIR}"/flash_*.img

    blank_image "${OUT_DIR}/flash_uboot.img" "$UBOOT_PART_SIZE"
    dd if="$UBOOT_BIN" of="${OUT_DIR}/flash_uboot.img" conv=notrunc status=none
    dd if="$UBOOT_BIN" of="${OUT_DIR}/flash_uboot.img" bs="$UBOOT_COPY_OFFSET" seek=1 conv=notrunc status=none

    pad_image "$OS_GZ" "${OUT_DIR}/flash_os0.img" "$OS_PART_SIZE"
    pad_image "$RD_GZ" "${OUT_DIR}/flash_rd0.img" "$RD_PART_SIZE"
    cp "${OUT_DIR}/flash_os0.img" "${OUT_DIR}/flash_os1.img"
    cp "${OUT_DIR}/flash_rd0.img" "${OUT_DIR}/flash_rd1.img"

    mkfs.jffs2 -b -U -s 0x1000 -e "$JFFS2_ERASE_SIZE" --pad="$APP_PART_SIZE" -d "$STAGE_DIR" -o "${OUT_DIR}/flash_app0.img"
    blank_image "${OUT_DIR}/flash_app1.img" "$APP_PART_SIZE"

    mkenvimage -b -p 0x00 -s "$ENV_PART_SIZE" -o "${OUT_DIR}/flash_env.img" "${MISC_FILES_DIR}/uboot.env"
    chmod 644 "${OUT_DIR}/flash_env.img"

    cp "$UBOOT_BIN" "${OUT_DIR}/u-boot-octeon_trx.bin"
    cp "$OS_GZ" "${OUT_DIR}/lsm_os_trx.gz"
    cp "$RD_GZ" "${OUT_DIR}/lsm_rd_trx.gz"

    log SUCCESS "Flash images written to ${OUT_DIR}"
}

validate_flash_images() {
    declare -A sizes=(
        [uboot]=$UBOOT_PART_SIZE
        [os0]=$OS_PART_SIZE [rd0]=$RD_PART_SIZE [app0]=$APP_PART_SIZE
        [os1]=$OS_PART_SIZE [rd1]=$RD_PART_SIZE [app1]=$APP_PART_SIZE
        [env]=$ENV_PART_SIZE
    )
    local total=0
    for p in "${!sizes[@]}"; do
        local f="${OUT_DIR}/flash_${p}.img"
        if [[ ! -f "$f" || $(stat -c %s "$f") -ne ${sizes[$p]} ]]; then
            log ERROR "flash_${p}.img missing or not ${sizes[$p]} bytes"
            exit 1
        fi
        total=$((total + sizes[$p]))
    done
    (( total == 0x8000000 )) || { log ERROR "Partitions do not add up to 128 MiB NOR"; exit 1; }

    (cd "$OUT_DIR" && sha256sum flash_*.img u-boot-octeon_trx.bin lsm_os_trx.gz lsm_rd_trx.gz > SHA256SUMS)
    log SUCCESS "All flash images match the NOR layout"
}

if [[ $# -lt 1 ]]; then
    log ERROR "Usage: $0 <lte-node-repo> [FDD|TDD]"
    exit 1
fi

LTE_NODE_REPO=$(realpath "$1")
MODE="${2:-FDD}"
if [[ "$MODE" != "FDD" && "$MODE" != "TDD" ]]; then
    log ERROR "Invalid mode '$MODE'. Expected 'FDD' or 'TDD'."
    exit 1
fi

if [[ -z "$BAND_CFG" ]]; then
    if [[ "$MODE" == "FDD" ]]; then
        BAND_CFG="startup_band3_20Mhz.cfg"
    else
        log ERROR "No default TDD band config, set BAND_CFG"
        exit 1
    fi
fi

check_pre_requisite
build_os
build_stack
validate_build_artifacts
stage_app
make_flash_images
validate_flash_images

log SUCCESS "TRX ${MODE} image creation complete (${BAND_CFG})"
exit 0
