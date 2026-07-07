#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2025-present, Ukama Inc.
#
# Standalone phase-2 dry run: does everything the flash script's phase 2 does
# except the dd writes. Waits for SSH, scps all 8 images to the board's staging
# dir, md5-checks each transfer, cleans up, then syncs. The flash is untouched,
# so the board's boot behavior afterwards reflects whatever dd wrote previously.
#
# Usage: export TRX_ROOT_PASSWORD='cavium.lte'; ./tools/trx-phase2-nodd.sh
# Env: TRX_IP (192.168.53.151), IMG_DIR (build-system/trx), STAGING (/mnt/tmp)

set -u

TRX_IP="${TRX_IP:-192.168.53.151}"
TRX_PW="${TRX_ROOT_PASSWORD:-cavium.lte}"
IMG_DIR="${IMG_DIR:-/usr/ukm/ukama/builder/scripts/build-system/trx}"
STAGING="${STAGING:-/mnt/tmp}"

ssh_cmd() {
    sshpass -p "$TRX_PW" ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null \
        -o LogLevel=ERROR -o ConnectTimeout=5 "root@${TRX_IP}" "$1"
}
scp_file() {
    sshpass -p "$TRX_PW" scp -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null \
        -o LogLevel=ERROR "$1" "root@${TRX_IP}:$2"
}

echo "Phase-2 dry run against ${TRX_IP} - no dd, flash is not written."
echo "Waiting for TRX SSH..."
elapsed=0
while [ "$elapsed" -lt 300 ]; do
    if ssh_cmd 'true' 2>/dev/null; then
        break
    fi
    sleep 2
    elapsed=$((elapsed + 2))
done
if [ "$elapsed" -ge 300 ]; then
    echo "ERROR: TRX not reachable at ${TRX_IP} within 300s."
    exit 1
fi
echo "SSH up."

ssh_cmd "mkdir -p ${STAGING}"

fail=0
for key in app0 app1 env os0 os1 rd0 rd1 uboot; do
    img="${IMG_DIR}/flash_${key}.img"
    name="flash_${key}.img"
    if [ ! -f "$img" ]; then
        echo "  [${key}] SKIP - image not found: $img"
        continue
    fi

    local_md5=$(md5sum "$img" | cut -d' ' -f1)

    echo "  [${key}] scp ${name} -> ${TRX_IP}:${STAGING}/"
    scp_file "$img" "${STAGING}/${name}"

    remote_md5=$(ssh_cmd "md5sum ${STAGING}/${name}" | cut -d' ' -f1)
    if [ "$local_md5" = "$remote_md5" ]; then
        echo "  [${key}] transfer ok (${local_md5}); skipping dd"
    else
        fail=1
        echo "  [${key}] TRANSFER MISMATCH local=${local_md5} remote=${remote_md5}"
    fi
    ssh_cmd "rm -f ${STAGING}/${name}"
done

echo ""
echo "Syncing TRX filesystems..."
ssh_cmd "sync"

if [ "$fail" -eq 0 ]; then
    echo "Done: all 8 transfers byte-exact, no dd performed, flash unchanged."
else
    echo "Done with TRANSFER MISMATCH above - scp corruption on this link."
fi
exit "$fail"
