#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2025-present, Ukama Inc.
#
# TRX boot soak test, run on the bench box after a board is flashed.
# Reboots the TRX N times, waits for each boot to reach the login prompt and
# records PASS/FAIL plus any DDR read-leveling errors from that boot's serial log.
# sshd comes up before the known bringup hang point, so the soak survives hung boots.
#
# Usage: export TRX_ROOT_PASSWORD='cavium.lte'; sudo -E ./tools/trx-boot-soak.sh [cycles]
# Env: TRX_IP (10.102.81.61), SERIAL (/dev/ttyUSB0), BOOT_TIMEOUT (900)

set -u

TRX_IP="${TRX_IP:-10.102.81.61}"
TRX_PW="${TRX_ROOT_PASSWORD:-cavium.lte}"
SERIAL="${SERIAL:-/dev/ttyUSB0}"
CYCLES="${1:-5}"
BOOT_TIMEOUT="${BOOT_TIMEOUT:-900}"

OUT="soak-$(date +%Y%m%d_%H%M%S)"
mkdir -p "$OUT"

ssh_cmd() {
    sshpass -p "$TRX_PW" ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null \
        -o LogLevel=ERROR -o ConnectTimeout=5 "root@${TRX_IP}" "$1"
}

if ! ssh_cmd 'true' 2>/dev/null; then
    echo "ERROR: TRX not reachable at ${TRX_IP} - board must be booted with SSH up to start."
    echo "  (Check: bench box has the post-flash IP: sudo ip addr add 10.102.81.100/24 dev enp0s31f6)"
    exit 1
fi

echo "TRX boot soak: ${CYCLES} cycles, board ${TRX_IP}, serial ${SERIAL}"
echo "PASS = boot completes to login prompt; FAIL = bringup hangs before login."
echo "Logs in ${OUT}/"

pass=0
fail=0
for i in $(seq 1 "$CYCLES"); do
    log="${OUT}/boot-${i}.serial.log"
    echo ""
    echo "=== Cycle ${i}/${CYCLES} ==="

    cat_pid=""
    if [ -e "$SERIAL" ]; then
        sudo stty -F "$SERIAL" 115200 cs8 -cstopb -parenb -ixon -echo raw 2>/dev/null || true
        sudo bash -c "cat '$SERIAL' > '$log'" &
        cat_pid=$!
    else
        echo "  (serial device ${SERIAL} not found - continuing without serial capture)"
    fi

    ssh_cmd reboot >/dev/null 2>&1 || true
    echo "  reboot sent; waiting up to ${BOOT_TIMEOUT}s for full boot (login/getty up)..."
    sleep 30
    elapsed=30
    booted=0
    while [ "$elapsed" -lt "$BOOT_TIMEOUT" ]; do
        if ssh_cmd 'ps w 2>/dev/null | grep -v grep | grep -q getty' 2>/dev/null; then
            booted=1
            break
        fi
        sleep 10
        elapsed=$((elapsed + 10))
    done

    if [ "$booted" -eq 1 ]; then
        pass=$((pass + 1))
        echo "  PASS - boot completed to login in ~${elapsed}s"
    else
        fail=$((fail + 1))
        echo "  FAIL - bringup did not reach login within ${BOOT_TIMEOUT}s"
        ssh_cmd 'cat /proc/loadavg; echo ---; ps w' > "${OUT}/boot-${i}.hang-state.txt" 2>/dev/null || true
        echo "  hang state saved to ${OUT}/boot-${i}.hang-state.txt"
    fi

    if [ -n "$cat_pid" ]; then
        sudo kill "$cat_pid" 2>/dev/null || true
    fi

    # clean read-leveling rows end in "(0)"; anything else = bit errors during DDR training
    if [ -s "$log" ]; then
        if grep -a "Rlevel Rank" "$log" | grep -qvE '\(0\)$'; then
            grep -a "Rlevel Rank" "$log" | grep -vE '\(0\)$' | sed 's/^/  DDR OUTLIER: /'
        else
            echo "  DDR leveling clean"
        fi
    fi
done

echo ""
echo "=== Soak result: ${pass} PASS / ${fail} FAIL out of ${CYCLES} cycles ==="
echo "Correlate FAILs with 'DDR OUTLIER' lines above - a match points at marginal"
echo "DDR on this board rather than software. Logs: ${OUT}/"
