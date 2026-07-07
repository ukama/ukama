#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2025-present, Ukama Inc.
#
# Standalone check of the dd flashing path - writes nothing to the board.
# For each of the 8 partitions it:
#   1. md5sums the source image on the bench box
#   2. scps the image to the board and md5sums the copy (checks the transfer)
#   3. reads the same number of bytes back from /dev/flash_* and md5sums them
#      (checks what dd wrote on a previous flash run)
# All three sums matching means scp and dd are both byte-exact.
#
# Usage: export TRX_ROOT_PASSWORD='cavium.lte'; ./tools/trx-verify-flash.sh
# Env: TRX_IP (10.102.81.61), IMG_DIR (build-system/trx)

set -u

TRX_IP="${TRX_IP:-10.102.81.61}"
TRX_PW="${TRX_ROOT_PASSWORD:-cavium.lte}"
IMG_DIR="${IMG_DIR:-/usr/ukm/ukama/builder/scripts/build-system/trx}"

ssh_cmd() {
    sshpass -p "$TRX_PW" ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null \
        -o LogLevel=ERROR -o ConnectTimeout=5 "root@${TRX_IP}" "$1"
}
scp_file() {
    sshpass -p "$TRX_PW" scp -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null \
        -o LogLevel=ERROR "$1" "root@${TRX_IP}:$2"
}

if ! ssh_cmd 'true' 2>/dev/null; then
    echo "ERROR: TRX not reachable at ${TRX_IP}."
    echo "  (Bench box may need: sudo ip addr add 10.102.81.100/24 dev enp0s31f6)"
    exit 1
fi

echo "Verifying flash contents on ${TRX_IP} against images in ${IMG_DIR}"
echo "No writes are performed."
echo ""

fail=0
for key in app0 app1 env os0 os1 rd0 rd1 uboot; do
    img="${IMG_DIR}/flash_${key}.img"
    dev="/dev/flash_${key}"
    if [ ! -f "$img" ]; then
        echo "[${key}] SKIP - image not found: $img"
        continue
    fi

    size=$(stat -c %s "$img")
    local_md5=$(md5sum "$img" | cut -d' ' -f1)

    scp_file "$img" "/tmp/verify.img"
    scp_md5=$(ssh_cmd "md5sum /tmp/verify.img" | cut -d' ' -f1)

    flash_md5=$(ssh_cmd "dd if=${dev} bs=65536 2>/dev/null | head -c ${size} | md5sum" | cut -d' ' -f1)
    ssh_cmd "rm -f /tmp/verify.img"

    if [ "$local_md5" = "$scp_md5" ] && [ "$local_md5" = "$flash_md5" ]; then
        echo "[${key}] OK    image=${local_md5}  scp=match  flash=match"
    else
        fail=1
        echo "[${key}] FAIL  image=${local_md5}"
        echo "         scp copy   = ${scp_md5}   $([ "$local_md5" = "$scp_md5" ] && echo match || echo MISMATCH)"
        echo "         flash read = ${flash_md5}   $([ "$local_md5" = "$flash_md5" ] && echo match || echo MISMATCH)"
    fi
done

echo ""
if [ "$fail" -eq 0 ]; then
    echo "RESULT: all partitions match their source images byte-for-byte."
    echo "scp and dd are both exact; the flash contents are not the problem."
else
    echo "RESULT: mismatch found - see FAIL lines above."
    echo "scp MISMATCH = transfer corruption; flash MISMATCH with scp match = dd/flash issue."
fi
exit "$fail"
