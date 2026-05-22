#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2025-present, Ukama Inc.

source "${LIB_DIR}/tftp.sh"
source "${LIB_DIR}/bdi.sh"
source "${LIB_DIR}/uboot_serial.sh"

REMOTE_BOOT_PID=""
OCT_TAIL_PID=""
SPAM_PID=""
TFTP_STAGE_DIR=""

_jtag_octeon_cleanup() {
    uboot_close
    tftp_stop
    if [ -n "$OCT_TAIL_PID" ]; then
        kill "$OCT_TAIL_PID" 2>/dev/null || true
        OCT_TAIL_PID=""
    fi
    if [ -n "$REMOTE_BOOT_PID" ]; then
        sudo kill "$REMOTE_BOOT_PID" 2>/dev/null || true
        REMOTE_BOOT_PID=""
    fi
    if [ -n "$SPAM_PID" ]; then
        kill "$SPAM_PID" 2>/dev/null || true
        SPAM_PID=""
    fi
    [ -n "$TFTP_STAGE_DIR" ] && sudo rm -rf "$TFTP_STAGE_DIR"
}

method_validate() {
    local fail=0
    local bdi_ip oct_path serial_dev band

    bdi_ip=$(yq_read "$BOARD_CONFIG" network.bdi_ip)
    oct_path=$(yq_read "$BOARD_CONFIG" oct_remote_boot.path)
    serial_dev=$(yq_read "$BOARD_CONFIG" serial.device)
    band="${BAND:-$(yq_read "$BOARD_CONFIG" band.default)}"

    if ! ping -c1 -W2 "$bdi_ip" >/dev/null 2>&1; then
        echo "  [FAIL] BDI not reachable at $bdi_ip"
        fail=1
    else
        echo "  [ OK ] BDI reachable at $bdi_ip"
    fi

    if [ ! -x "$oct_path" ]; then
        echo "  [FAIL] oct-remote-boot not found at $oct_path"
        fail=1
    else
        echo "  [ OK ] oct-remote-boot present"
    fi

    if [ ! -e "$serial_dev" ]; then
        echo "  [FAIL] serial device not found: $serial_dev"
        fail=1
    else
        echo "  [ OK ] serial device: $serial_dev"
    fi

    local key
    for key in os rd uboot; do
        local p
        p=$(yq_read "$BOARD_CONFIG" "phase1.artifacts.${key}.path")
        if [ ! -f "$p" ]; then
            echo "  [FAIL] phase1 artifact missing ($key): $p"
            fail=1
        else
            echo "  [ OK ] phase1 $key: $p"
        fi
    done

    local img_keys
    img_keys=$(yq_keys "$BOARD_CONFIG" phase2.images)
    for key in $img_keys; do
        local p
        p=$(yq_read "$BOARD_CONFIG" "phase2.images.${key}.src")
        if [ ! -f "$p" ]; then
            echo "  [FAIL] phase2 image missing ($key): $p"
            fail=1
        else
            echo "  [ OK ] phase2 $key"
        fi
    done

    local band_cfg
    band_cfg="$(yq_read "$BOARD_CONFIG" band.configs_dir)/${band}.cfg"
    if [ ! -f "$band_cfg" ]; then
        echo "  [FAIL] band config not found: $band_cfg"
        fail=1
    else
        echo "  [ OK ] band: ${band} ($band_cfg)"
    fi

    if [ -z "${TRX_ROOT_PASSWORD:-}" ]; then
        echo "  [FAIL] TRX_ROOT_PASSWORD environment variable is not set"
        fail=1
    else
        echo "  [ OK ] TRX_ROOT_PASSWORD is set"
    fi

    if ! command -v sshpass >/dev/null 2>&1; then
        echo "  [..] installing sshpass..."
        sudo apt-get update -qq
        sudo apt-get install -y sshpass
        if command -v sshpass >/dev/null 2>&1; then
            echo "  [ OK ] sshpass installed"
        else
            echo "  [FAIL] could not install sshpass"
            fail=1
        fi
    else
        echo "  [ OK ] sshpass available"
    fi

    return $fail
}

method_confirm() {
    local band trx_ip
    band="${BAND:-$(yq_read "$BOARD_CONFIG" band.default)}"
    trx_ip=$(yq_read "$BOARD_CONFIG" network.trx_ip)

    echo ""
    echo "Plan:"
    echo "  Phase 1 (JTAG)  : TFTP+JTAG bringup, flash OS/RD/uboot to ${trx_ip}"
    echo "  Manual pause    : you power-cycle TRX and remove BDI cable"
    echo "  Phase 2 (SSH)   : scp+dd 8 .img files to /dev/flash_*"
    echo "  Band            : ${band}"
    echo ""
    echo "This will overwrite the TRX flash."
    read -rp "Type 'yes' to continue: " confirm
    [ "$confirm" = "yes" ]
}

_phase1_uboot_env() {
    local dev="$1"
    local prompt="$2"
    local trx_mac trx_ip netmask host_ip

    trx_mac=$(yq_read "$BOARD_CONFIG" network.trx_mac)
    trx_ip=$(yq_read "$BOARD_CONFIG" network.trx_ip)
    netmask=$(yq_read "$BOARD_CONFIG" network.netmask)
    host_ip=$(yq_read "$BOARD_CONFIG" network.host_ip)

    local mac_dashed="${trx_mac//:/-}"

    uboot_send_and_wait "$dev" "setenv ethaddr ${mac_dashed}" "$prompt" 5
    uboot_send_and_wait "$dev" "setenv ipaddr ${trx_ip}" "$prompt" 5
    uboot_send_and_wait "$dev" "setenv netmask ${netmask}" "$prompt" 5
    uboot_send_and_wait "$dev" "setenv serverip ${host_ip}" "$prompt" 5
    uboot_send_and_wait "$dev" "setenv bootby flash" "$prompt" 5
    uboot_send_and_wait "$dev" "setenv cfgloadby flash" "$prompt" 5
    uboot_send_and_wait "$dev" "setenv swloadby flash" "$prompt" 5
    uboot_send_and_wait "$dev" 'setenv i2cinit "i2c dev 0; i2c probe; i2c dev 1; i2c probe"' "$prompt" 5
    uboot_send_and_wait "$dev" 'setenv bootcmd "run i2cinit; run namedalloc; run bootcby${bootby}"' "$prompt" 5
    uboot_send_and_wait "$dev" 'setenv bootcbytftp "tftp 0x21000000 lsm_os.gz; gunzip 0x21000000 0x20000000 0x1000000; tftp 0x30800000 lsm_rd.gz; bootoctlinux 0x20000000 coremask=0x7 endbootargs rd_name=initrd mem=512M;"' "$prompt" 5
    uboot_send_and_wait "$dev" 'setenv namedalloc "namedalloc dsp-dump 0x400000 0x7f4D0000; namedalloc cazac 0x630000 0x7f8D0000; namedalloc cpu-dsp-if 0x100000 0x7ff00000; namedalloc dsp-log-buf 0x4000000 0x80000000; namedalloc initrd 0x2800000 0x30800000;"' "$prompt" 5
    uboot_send_and_wait "$dev" "setenv mk_ubootenv 1" "$prompt" 5
}

_phase1_flash_artifact() {
    local dev="$1"
    local prompt="$2"
    local key="$3"
    local ddr_addr="$4"

    local src flash_addr name
    src=$(yq_read "$BOARD_CONFIG" "phase1.artifacts.${key}.path")
    flash_addr=$(yq_read "$BOARD_CONFIG" "phase1.artifacts.${key}.flash_addr")
    name=$(tftp_stage_file "$src")

    echo "  Flashing $key ($name) to ${flash_addr}..."

    local marker_before
    marker_before=$(wc -c < "$UBOOT_LOG" 2>/dev/null || echo 0)

    uboot_send_and_wait "$dev" "tftp ${ddr_addr} ${name}" "$prompt" 120

    if ! tail -c +"$marker_before" "$UBOOT_LOG" 2>/dev/null | grep -q "Bytes transferred = "; then
        echo "ERROR: tftp failed for $key — no 'Bytes transferred' seen in log"
        return 1
    fi

    uboot_send_and_wait "$dev" "protect off all" "$prompt" 10
    uboot_send_and_wait "$dev" "erase ${flash_addr} +\${filesize}" "$prompt" 60
    uboot_send_and_wait "$dev" "cp.b ${ddr_addr} ${flash_addr} \${filesize}" "$prompt" 120
}

_phase1_run() {
    local bdi_ip bdi_prompt serial_dev baud uboot_prompt oct_path oct_board oct_clock
    local ddr_os ddr_rd gdb_port oct_env_root oct_env_protocol

    bdi_ip=$(yq_read "$BOARD_CONFIG" network.bdi_ip)
    bdi_prompt=$(yq_read "$BOARD_CONFIG" bdi.prompt)
    gdb_port=$(yq_read "$BOARD_CONFIG" bdi.gdb_port)
    serial_dev=$(yq_read "$BOARD_CONFIG" serial.device)
    baud=$(yq_read "$BOARD_CONFIG" serial.baud)
    uboot_prompt=$(yq_read "$BOARD_CONFIG" serial.uboot_prompt)
    oct_path=$(yq_read "$BOARD_CONFIG" oct_remote_boot.path)
    oct_board=$(yq_read "$BOARD_CONFIG" oct_remote_boot.board)
    oct_clock=$(yq_read "$BOARD_CONFIG" oct_remote_boot.ddr_clock_hz)
    ddr_os=$(yq_read "$BOARD_CONFIG" phase1.ddr_os_load_addr)
    ddr_rd=$(yq_read "$BOARD_CONFIG" phase1.ddr_rd_load_addr)

    oct_env_root=$(dirname "$(dirname "$(dirname "$oct_path")")")
    oct_env_protocol="GDB:${bdi_ip},${gdb_port}"

    TFTP_STAGE_DIR=$(mktemp -d /tmp/ukama-trx-tftp.XXXXXX)
    echo "TFTP staging at $TFTP_STAGE_DIR"
    sudo pkill -x in.tftpd 2>/dev/null || true
    sleep 1
    tftp_serve "$TFTP_STAGE_DIR"
    sleep 1
    if ! ss -lnup | grep -q ':69 '; then
        echo "ERROR: TFTP server failed to start on port 69"
        return 1
    fi

    local bdi_config_src
    bdi_config_src=$(yq_read "$BOARD_CONFIG" bdi.config_file)
    if [ -f "$bdi_config_src" ]; then
        sudo cp "$bdi_config_src" "${TFTP_STAGE_DIR}/cnf71xx.cfg"
        sudo chmod 644 "${TFTP_STAGE_DIR}/cnf71xx.cfg"
        echo "Staged cnf71xx.cfg in TFTP root for BDI auto-load"

        local bdi_config_dir
        bdi_config_dir=$(dirname "$bdi_config_src")
        local def_file
        for def_file in "$bdi_config_dir"/*.def; do
            [ -f "$def_file" ] || continue
            sudo cp "$def_file" "${TFTP_STAGE_DIR}/$(basename "$def_file")"
            sudo chmod 644 "${TFTP_STAGE_DIR}/$(basename "$def_file")"
            echo "Staged $(basename "$def_file") in TFTP root for BDI auto-load"
        done
    else
        echo "WARNING: bdi.config_file not found at $bdi_config_src"
    fi

    local host_ip
    host_ip=$(yq_read "$BOARD_CONFIG" network.host_ip)

    local host_ip
    host_ip=$(yq_read "$BOARD_CONFIG" network.host_ip)

    echo "Checking BDI state..."
    if ! bdi_send_sequence "$bdi_ip" "$bdi_prompt" 10 "HALT" >/dev/null 2>&1; then
        echo "  BDI unconfigured (prompt is not '${bdi_prompt}') — loading cnf71xx.cfg via TFTP..."
        local bdi_load_log="${LOG_DIR}/bdi-config-load.log"

        # Set TFTP host explicitly, then load config
        bdi_send_command "$bdi_ip" "Core#0>" "HOST ${host_ip}" 15 >"$bdi_load_log" 2>&1 || true
        bdi_send_command "$bdi_ip" "Core#0>" "CONFIG cnf71xx.cfg" 60 >>"$bdi_load_log" 2>&1 || true
        sleep 5

        if ! bdi_send_sequence "$bdi_ip" "$bdi_prompt" 10 "HALT" >/dev/null 2>&1; then
            echo "ERROR: BDI config still not loaded after CONFIG attempt"
            echo "--- BDI config load output ---"
            cat "$bdi_load_log" 2>/dev/null || true
            return 1
        fi
        echo "  BDI config loaded successfully"
    fi

    echo "Telneting BDI at ${bdi_ip}: halting core..."
    if ! bdi_send_sequence "$bdi_ip" "$bdi_prompt" 90 "HALT"; then
        echo "ERROR: BDI did not respond with '${bdi_prompt}' after HALT"
        return 1
    fi

    local oct_log="${LOG_DIR}/oct-remote-boot.log"
    echo "Running oct-remote-boot:"
    echo "  OCTEON_ROOT=$oct_env_root"
    echo "  OCTEON_REMOTE_PROTOCOL=$oct_env_protocol"
    echo "  sudo $oct_path --board=$oct_board --ddr_clock_hz=$oct_clock"
    echo "  log: $oct_log"
    sudo env OCTEON_ROOT="$oct_env_root" OCTEON_REMOTE_PROTOCOL="$oct_env_protocol" \
        "$oct_path" --board="$oct_board" --ddr_clock_hz="$oct_clock" >"$oct_log" 2>&1 &
    REMOTE_BOOT_PID=$!
    echo "  oct-remote-boot started (PID $REMOTE_BOOT_PID)"

    sudo tail -n +1 -F "$oct_log" 2>/dev/null &
    OCT_TAIL_PID=$!

    echo "Opening serial console at $serial_dev ($baud)..."
    uboot_open "$serial_dev" "$baud" "${LOG_DIR}/uboot.log"

    echo "Spamming keys to interrupt zero-second autoboot (until prompt appears)..."
    (
        exec 3>"$serial_dev"
        while true; do
            printf ' \r\n' >&3
            sleep 0.03
        done
    ) &
    SPAM_PID=$!

    echo "Waiting for u-boot prompt '${uboot_prompt}' (up to 120s)..."
    if ! uboot_wait_for "$uboot_prompt" 120; then
        # Stop tail first so its output doesn't race with our diagnostics
        kill "$OCT_TAIL_PID" 2>/dev/null || true
        sleep 1

        echo "ERROR: u-boot prompt did not appear within 120s"
        echo "--- last 40 lines of oct-remote-boot output ---"
        tail -n 40 "$oct_log" 2>/dev/null || true
        echo "--- oct-remote-boot exit status ---"
        if kill -0 "$REMOTE_BOOT_PID" 2>/dev/null; then
            echo "  still running (PID $REMOTE_BOOT_PID) — did not error out, but no u-boot on serial"
            sudo kill "$REMOTE_BOOT_PID" 2>/dev/null || true
        else
            wait "$REMOTE_BOOT_PID" 2>/dev/null
            echo "  exited with status $?"
        fi
        echo "--- last 40 lines of serial (uboot.log) ---"
        tail -n 40 "${LOG_DIR}/uboot.log" 2>/dev/null || true
        return 1
    fi

    kill "$OCT_TAIL_PID" 2>/dev/null || true
    if [ -n "$SPAM_PID" ]; then
        kill "$SPAM_PID" 2>/dev/null || true
        SPAM_PID=""
    fi

    echo "Pushing u-boot environment variables..."
    _phase1_uboot_env "$serial_dev" "$uboot_prompt"

    echo "Enabling ethernet (mw64 x2) and saving env..."
    uboot_send_and_wait "$serial_dev" "mw64 0x00011800B0001000 0x0140" "$uboot_prompt" 5
    uboot_send_and_wait "$serial_dev" "mw64 0x00011800B0001000 0x0140" "$uboot_prompt" 5
    uboot_send_and_wait "$serial_dev" "saveenv" "$uboot_prompt" 15

    local host_ip
    host_ip=$(yq_read "$BOARD_CONFIG" network.host_ip)
    echo "Pinging host ${host_ip} from u-boot..."
    uboot_send_and_wait "$serial_dev" "ping ${host_ip}" "$uboot_prompt" 15 || {
        echo "WARNING: ping failed — retrying mw64 x2..."
        uboot_send_and_wait "$serial_dev" "mw64 0x00011800B0001000 0x0140" "$uboot_prompt" 5
        uboot_send_and_wait "$serial_dev" "mw64 0x00011800B0001000 0x0140" "$uboot_prompt" 5
        uboot_send_and_wait "$serial_dev" "ping ${host_ip}" "$uboot_prompt" 15 || {
            echo "ERROR: TRX cannot reach host. Check cables."
            return 1
        }
    }

    _phase1_flash_artifact "$serial_dev" "$uboot_prompt" "os"    "$ddr_os"
    _phase1_flash_artifact "$serial_dev" "$uboot_prompt" "rd"    "$ddr_rd"
    _phase1_flash_artifact "$serial_dev" "$uboot_prompt" "uboot" "$ddr_os"

    uboot_close
    tftp_stop
    [ -n "$REMOTE_BOOT_PID" ] && sudo kill "$REMOTE_BOOT_PID" 2>/dev/null || true
    REMOTE_BOOT_PID=""
}

_phase2_run() {
    local trx_ip ssh_user staging band band_cfg target_path
    trx_ip=$(yq_read "$BOARD_CONFIG" network.trx_ip)
    ssh_user=$(yq_read "$BOARD_CONFIG" phase2.ssh_user)
    staging=$(yq_read "$BOARD_CONFIG" phase2.ssh_staging_dir)
    band="${BAND:-$(yq_read "$BOARD_CONFIG" band.default)}"
    band_cfg="$(yq_read "$BOARD_CONFIG" band.configs_dir)/${band}.cfg"
    target_path=$(yq_read "$BOARD_CONFIG" band.target_path)

    local sshpass_args=(-p "$TRX_ROOT_PASSWORD")
    local ssh_opts=(-o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o LogLevel=ERROR)

    echo "Waiting for TRX SSH at ${trx_ip}..."
    local elapsed=0
    while [ "$elapsed" -lt 120 ]; do
        if sshpass "${sshpass_args[@]}" ssh "${ssh_opts[@]}" "${ssh_user}@${trx_ip}" true 2>/dev/null; then
            break
        fi
        sleep 2
        elapsed=$((elapsed + 2))
    done
    if [ "$elapsed" -ge 120 ]; then
        echo "ERROR: TRX not reachable via SSH within 120s"
        return 1
    fi
    echo "SSH up."

    sshpass "${sshpass_args[@]}" ssh "${ssh_opts[@]}" "${ssh_user}@${trx_ip}" "mkdir -p ${staging}"

    local img_keys
    img_keys=$(yq_keys "$BOARD_CONFIG" phase2.images)
    local key
    for key in $img_keys; do
        local src dst name
        src=$(yq_read "$BOARD_CONFIG" "phase2.images.${key}.src")
        dst=$(yq_read "$BOARD_CONFIG" "phase2.images.${key}.dst")
        name=$(basename "$src")

        echo "  [${key}] scp ${name} -> ${trx_ip}:${staging}/"
        sshpass "${sshpass_args[@]}" scp "${ssh_opts[@]}" "$src" "${ssh_user}@${trx_ip}:${staging}/${name}"

        echo "  [${key}] dd to ${dst}"
        sshpass "${sshpass_args[@]}" ssh "${ssh_opts[@]}" "${ssh_user}@${trx_ip}" \
            "dd if=${staging}/${name} of=${dst} bs=1 && rm -f ${staging}/${name}"
    done

    echo "Copying band config (${band}) to ${target_path}..."
    sshpass "${sshpass_args[@]}" scp "${ssh_opts[@]}" "$band_cfg" "${ssh_user}@${trx_ip}:${target_path}"
}

method_apply() {
    trap _jtag_octeon_cleanup EXIT

    echo "=== Phase 1: JTAG bringup ==="
    _phase1_run

    echo ""
    echo "=== Manual pause ==="
    echo "Please:"
    echo "  1. Power OFF the TRX"
    echo "  2. Disconnect the BDI / JTAG cable"
    echo "  3. Power ON the TRX"
    echo "  4. Wait until it boots to Linux (should be reachable at $(yq_read "$BOARD_CONFIG" network.trx_ip))"
    echo ""
    read -rp "Press ENTER when ready: " _

    echo ""
    echo "=== Phase 2: SSH + dd image flash ==="
    _phase2_run
}

method_verify() {
    local trx_ip ssh_user target_path
    trx_ip=$(yq_read "$BOARD_CONFIG" network.trx_ip)
    ssh_user=$(yq_read "$BOARD_CONFIG" phase2.ssh_user)
    target_path=$(yq_read "$BOARD_CONFIG" band.target_path)

    local ssh_opts=(-o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o LogLevel=ERROR)

    if sshpass -p "$TRX_ROOT_PASSWORD" ssh "${ssh_opts[@]}" "${ssh_user}@${trx_ip}" "test -f ${target_path}"; then
        echo "  [ OK ] band config present at ${target_path}"
    else
        echo "  [FAIL] band config missing at ${target_path}"
        return 1
    fi
}

method_monitor() {
    echo "TRX flash complete. Power-cycle and observe operation."
    return 0
}
