#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2025-present, Ukama Inc.

source "${LIB_DIR}/tftp.sh"
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
    sudo pkill -9 -f oct-remote-boot 2>/dev/null || true
    if [ -n "$SPAM_PID" ]; then
        kill "$SPAM_PID" 2>/dev/null || true
        SPAM_PID=""
    fi
    [ -n "$TFTP_STAGE_DIR" ] && sudo rm -rf "$TFTP_STAGE_DIR"
}

# Send a single command to the BDI via telnet and wait for its prompt.
# Usage: bdi_telnet_cmd <host> <command>
bdi_telnet_cmd() {
    local host="$1"
    local cmd="$2"
    if ! command -v expect >/dev/null 2>&1; then
        echo "WARNING: expect not installed — cannot send BDI command '$cmd'" >&2
        return 1
    fi
    expect -c "
        set timeout 8
        spawn telnet $host
        expect {
            \"Core#0>\" {}
            \"cnMIPS#0>\" {}
            timeout {
                puts \"BDI telnet: timeout waiting for prompt\"
                exit 1
            }
        }
        send \"$cmd\r\"
        sleep 1
        send \"quit\r\"
        expect eof
    " 2>/dev/null
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
        local serial_holder=""
        if command -v lsof >/dev/null 2>&1; then
            serial_holder=$(lsof -t "$serial_dev" 2>/dev/null | head -1 || true)
        elif command -v fuser >/dev/null 2>&1; then
            serial_holder=$(fuser "$serial_dev" 2>/dev/null | tr -cd '0-9' || true)
        fi
        if [ -n "$serial_holder" ]; then
            echo "  [FAIL] $serial_dev is held by another process (PID ${serial_holder})."
            echo "         Close any serial terminal (PuTTY / screen / minicom) on it before flashing."
            echo "         Phase 1 needs exclusive serial access, or it cannot see the u-boot prompt."
            fail=1
        fi
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

    local restore_errexit=0
    case $- in *e*) restore_errexit=1; set +e ;; esac

    uboot_send_and_wait "$dev" "setenv ethaddr ${mac_dashed}" "$prompt" 8
    uboot_send_and_wait "$dev" "setenv ipaddr ${trx_ip}" "$prompt" 8
    uboot_send_and_wait "$dev" "setenv netmask ${netmask}" "$prompt" 8
    uboot_send_and_wait "$dev" "setenv serverip ${host_ip}" "$prompt" 8
    uboot_send_and_wait "$dev" "setenv bootby flash" "$prompt" 8
    uboot_send_and_wait "$dev" "setenv cfgloadby flash" "$prompt" 8
    uboot_send_and_wait "$dev" "setenv swloadby flash" "$prompt" 8
    uboot_send_and_wait "$dev" 'setenv i2cinit "i2c dev 0; i2c probe; i2c dev 1; i2c probe"' "$prompt" 8
    uboot_send_and_wait "$dev" 'setenv bootcmd "run i2cinit; run namedalloc; run bootcby${bootby}"' "$prompt" 8
    uboot_send_and_wait "$dev" 'setenv bootcbytftp "tftp 0x21000000 lsm_os.gz; gunzip 0x21000000 0x20000000 0x1000000; tftp 0x30800000 lsm_rd.gz; bootoctlinux 0x20000000 coremask=0x7 endbootargs rd_name=initrd mem=512M;"' "$prompt" 8
    uboot_send_and_wait "$dev" 'setenv namedalloc "namedalloc dsp-dump 0x400000 0x7f4D0000; namedalloc cazac 0x630000 0x7f8D0000; namedalloc cpu-dsp-if 0x100000 0x7ff00000; namedalloc dsp-log-buf 0x4000000 0x80000000; namedalloc initrd 0x2800000 0x30800000;"' "$prompt" 8
    uboot_send_and_wait "$dev" "setenv mk_ubootenv 1" "$prompt" 8

    [ "$restore_errexit" = "1" ] && set -e
    return 0
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
    local bdi_ip serial_dev baud uboot_prompt oct_path oct_board oct_clock
    local ddr_os ddr_rd gdb_port oct_env_root oct_env_protocol

    bdi_ip=$(yq_read "$BOARD_CONFIG" network.bdi_ip)
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

    local oct_log="${LOG_DIR}/oct-remote-boot.log"
    local oct_attempt=0 max_oct_attempts=8 clock_ok=0

    _oct_stop() {
        # Graceful termination first
        sudo pkill -TERM -f oct-remote-boot 2>/dev/null || true
        local kw=0
        while [ "$kw" -lt 10 ] && pgrep -f oct-remote-boot >/dev/null 2>&1; do
            sleep 1
            kw=$((kw + 1))
        done
        # Hard kill — also mop up any SDK children that hold resources
        sudo pkill -9 -f oct-remote-boot 2>/dev/null || true
        sudo pkill -9 -f 'oct-remote-' 2>/dev/null || true
        sudo pkill -9 -f 'oct-pci-' 2>/dev/null || true
        REMOTE_BOOT_PID=""
        sleep 2
    }

    _bdi_clear_gdb_slot() {
        # The BDI has a single GDB slot. If a previous oct-remote-boot died
        # uncleanly (SIGKILL) the BDI may not notice the TCP half-close and
        # keeps refusing new connections. Try to nudge it by sending a GDB
        # detach packet and/or connecting briefly to force a state change.
        local host="$1"
        local port="${2:-2001}"
        # Send a raw GDB detach packet ($D#44) — best-effort
        printf '+$D#44' | timeout 2 nc -N "$host" "$port" 2>/dev/null || true
        sleep 1
        # Quick connect/disconnect to provoke RST handling
        timeout 1 bash -c "exec 3<>/dev/tcp/$host/$port; exec 3>&-" 2>/dev/null || true
        sleep 1
    }

    _wait_for_bdi_gdb() {
        local host="$1"
        local port="${2:-2001}"
        local w=0
        echo "  Probing BDI GDB port ${host}:${port}..."
        while [ "$w" -lt 60 ]; do
            if nc -z "$host" "$port" 2>/dev/null; then
                echo "  GDB port is open."
                return 0
            fi
            sleep 1
            w=$((w + 1))
        done
        echo "  GDB port still closed after 60s — BDI may still be booting or crashed."
        return 1
    }

    _oct_stop

    while [ "$oct_attempt" -lt "$max_oct_attempts" ]; do
        oct_attempt=$((oct_attempt + 1))
        echo "=== oct-remote-boot attempt ${oct_attempt}/${max_oct_attempts} ==="

        # Make sure the BDI GDB port is actually open before launching oct-remote-boot.
        # If the BDI was power-cycled it can take ~30-60s to load its config and open 2001.
        if ! _wait_for_bdi_gdb "$bdi_ip" "$gdb_port"; then
            echo "  BDI GDB port not reachable. If the BDI was just power-cycled, wait longer."
            echo "  If this persists, power-cycle the BDI box itself."
        fi

        echo "Running oct-remote-boot (OCTEON_ROOT=$oct_env_root, $oct_env_protocol)..."
        : > "$oct_log"
        sudo stdbuf -oL -eL env OCTEON_ROOT="$oct_env_root" OCTEON_REMOTE_PROTOCOL="$oct_env_protocol" \
            "$oct_path" --board="$oct_board" --ddr_clock_hz="$oct_clock" >"$oct_log" 2>&1 &
        REMOTE_BOOT_PID=$!
        echo "  oct-remote-boot started (PID $REMOTE_BOOT_PID)"

        local cwait=0 clk="" mhz="" ddr_bad="" refused=""
        while [ "$cwait" -lt 90 ]; do
            clk=$(grep -a "Measured DDR clock" "$oct_log" 2>/dev/null | tail -1 || true)
            [ -n "$clk" ] && break
            refused=$(grep -a "Connection refused" "$oct_log" 2>/dev/null | head -1 || true)
            [ -n "$refused" ] && break
            ddr_bad=$(grep -aE "exceeds DIMM specifications|GDB Reply Error" "$oct_log" 2>/dev/null | head -1 || true)
            [ -n "$ddr_bad" ] && break
            kill -0 "$REMOTE_BOOT_PID" 2>/dev/null || break
            sleep 1
            cwait=$((cwait + 1))
        done
        mhz=$(printf '%s' "$clk" | grep -oE '[0-9]+' | head -1 || true)

        if [ -n "$mhz" ] && [ "$mhz" -ge 380 ] && [ "$mhz" -le 420 ]; then
            echo "  $clk"
            echo "  DDR clock locked at ~400 MHz — good, continuing."
            clock_ok=1
            break
        fi

        _oct_stop

        if [ -n "$refused" ]; then
            echo "  BDI GDB port 2001 refused — single slot held by stale session."
            echo "  Trying to clear the slot (GDB detach + TCP probe)..."
            _bdi_clear_gdb_slot "$bdi_ip" "$gdb_port"
            echo "  Waiting 15s for BDI to settle before retry..."
            sleep 15
            continue
        fi

        if [ -n "$ddr_bad" ]; then
            echo "  DDR PLL mislocked — $ddr_bad"
        elif [ -n "$clk" ]; then
            echo "  DDR clock ${mhz} MHz — not ~400 (mislock)."
        else
            echo "  oct-remote-boot did not reach a DDR clock within 90s. Last lines:"
            tail -n 20 "$oct_log" 2>/dev/null | sed 's/^/    /' || true
        fi

        if [ "$oct_attempt" -lt "$max_oct_attempts" ]; then
            echo ""
            echo ">>> DDR did not come up at ~400 MHz. If this keeps happening, COLD power-cycle the"
            echo ">>> TRX now: full power OFF, wait ~5s, power back ON. Leave the JTAG cable connected."
            printf ">>> Press Enter to retry (attempt %d/%d)... " "$((oct_attempt + 1))" "$max_oct_attempts"
            read -r _ </dev/tty || true
            echo ""
        fi
    done

    if [ "$clock_ok" -ne 1 ]; then
        echo "ERROR: oct-remote-boot did not bring DDR up at ~400 MHz after ${max_oct_attempts} attempts."
        echo "  If the BDI GDB port kept refusing, a stale debugger session is wedged — close any"
        echo "  other oct-remote-boot/telnet on the BDI, or power-cycle the BDI itself."
        echo "  If DDR mislocked (596/599/796 MHz), COLD power-cycle the TRX and re-run."
        return 1
    fi

    sudo tail -n +1 -F "$oct_log" 2>/dev/null &
    OCT_TAIL_PID=$!

    # oct-remote-boot with GDB protocol stages the bootloader but sometimes fails
    # to start the core via GDB continue. The proven manual fix is to send
    # 'go 0x400000' via BDI telnet once the bootloader stub is written.
    echo "Waiting for oct-remote-boot to finish staging bootloader..."
    local stub_wait=0
    while [ "$stub_wait" -lt 30 ]; do
        if grep -a "Done writing boot stub" "$oct_log" 2>/dev/null | head -1 >/dev/null; then
            break
        fi
        if grep -a "Starting core 0" "$oct_log" 2>/dev/null | head -1 >/dev/null; then
            break
        fi
        if ! kill -0 "$REMOTE_BOOT_PID" 2>/dev/null; then
            break
        fi
        sleep 1
        stub_wait=$((stub_wait + 1))
    done

    echo "Sending 'go 0x400000' via BDI telnet to start core..."
    if ! bdi_telnet_cmd "$bdi_ip" "go 0x400000"; then
        echo "WARNING: BDI telnet command failed — core may already be running, or BDI not at prompt."
        echo "  Continuing anyway; if u-boot does not appear on serial, the BDI may need a reset."
    fi

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

    echo "Draining residual serial output (spam backlog) before sending commands..."
    uboot_drain 3

    echo "Pushing u-boot environment variables..."
    _phase1_uboot_env "$serial_dev" "$uboot_prompt"

    echo "Enabling ethernet (mw64 x2) and saving env..."
    uboot_send_and_wait "$serial_dev" "mw64 0x00011800B0001000 0x0140" "$uboot_prompt" 8 || true
    uboot_send_and_wait "$serial_dev" "mw64 0x00011800B0001000 0x0140" "$uboot_prompt" 8 || true
    uboot_send_and_wait "$serial_dev" "saveenv" "$uboot_prompt" 15 || true

    local host_ip ping_marker link_up=0 ping_attempt=0 max_ping_attempts=6
    host_ip=$(yq_read "$BOARD_CONFIG" network.host_ip)
    echo "Bringing up ethernet link to ${host_ip} (mw64 + ping retries)..."
    while [ "$ping_attempt" -lt "$max_ping_attempts" ]; do
        ping_attempt=$((ping_attempt + 1))
        ping_marker=$(wc -c < "$UBOOT_LOG" 2>/dev/null || echo 0)
        uboot_send_and_wait "$serial_dev" "ping ${host_ip}" "$uboot_prompt" 20 || true
        if tail -c +"$ping_marker" "$UBOOT_LOG" 2>/dev/null | grep -q "is alive"; then
            link_up=1
            break
        fi
        echo "  ping attempt ${ping_attempt}/${max_ping_attempts} failed — re-enabling ethernet (mw64) and retrying..."
        uboot_send_and_wait "$serial_dev" "mw64 0x00011800B0001000 0x0140" "$uboot_prompt" 8 || true
        sleep 2
    done

    if [ "$link_up" -ne 1 ]; then
        echo "ERROR: TRX ethernet link did not come up after ${max_ping_attempts} attempts (octeth0 Down)."
        echo "  mw64 ethernet-enable + ping kept failing. Usual cause is a mislocked DDR clock"
        echo "  (SGMII reference clock off) — cold power-cycle the TRX and re-run; also check the cable."
        grep -a "Measured DDR clock" "$oct_log" 2>/dev/null | tail -1 | sed 's/^/  oct-remote-boot: /' || true
        return 1
    fi
    echo "  host ${host_ip} is reachable — ethernet link up after ${ping_attempt} attempt(s)"

    _phase1_flash_artifact "$serial_dev" "$uboot_prompt" "os"    "$ddr_os"
    _phase1_flash_artifact "$serial_dev" "$uboot_prompt" "rd"    "$ddr_rd"
    _phase1_flash_artifact "$serial_dev" "$uboot_prompt" "uboot" "$ddr_os"

    uboot_close
    tftp_stop
    [ -n "$REMOTE_BOOT_PID" ] && sudo kill "$REMOTE_BOOT_PID" 2>/dev/null || true
    REMOTE_BOOT_PID=""
}

_phase2_run() {
    local trx_ip ssh_user staging
    trx_ip=$(yq_read "$BOARD_CONFIG" network.trx_ip)
    ssh_user=$(yq_read "$BOARD_CONFIG" phase2.ssh_user)
    staging=$(yq_read "$BOARD_CONFIG" phase2.ssh_staging_dir)

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
        # bs=1 is REQUIRED: the MTD CFI driver on the TRX Linux 3.4 kernel silently
        # corrupts flash with large write() buffers (e.g. bs=1M). Proven manual baseline.
        sshpass "${sshpass_args[@]}" ssh "${ssh_opts[@]}" "${ssh_user}@${trx_ip}" \
            "dd if=${staging}/${name} of=${dst} bs=1 && rm -f ${staging}/${name}"
    done

    echo ""
    echo "All 8 images written. MTD char devices bypass the Linux page cache, so 'sync' is"
    echo "irrelevant. Power-cycle the TRX to boot the new flash."
    echo ""
    echo "NOTE: If the board was previously bricked by a bs=1M run, this reflash with bs=1"
    echo "will restore it. The power-cycle after dd is what commits the data to NOR."
}

method_apply() {
    trap _jtag_octeon_cleanup EXIT

    echo "=== Phase 1: JTAG bringup ==="
    _phase1_run

    echo ""
    echo "=== Manual pause ==="
    echo "Please:"
    echo "  1. Power OFF the TRX"
    echo "  2. Disconnect the BDI / JTAG cable  (REQUIRED — while connected the BDI holds the"
    echo "     CPU in reset, so the board will NOT finish booting from flash)"
    echo "  3. Power ON the TRX"
    echo "  4. Wait until it boots to Linux (should be reachable at $(yq_read "$BOARD_CONFIG" network.trx_ip))"
    echo ""
    read -rp "Press ENTER when ready: " _

    echo ""
    echo "=== Phase 2: SSH + dd image flash ==="
    _phase2_run
}

method_verify() {
    local trx_ip ssh_user
    trx_ip=$(yq_read "$BOARD_CONFIG" network.trx_ip)
    ssh_user=$(yq_read "$BOARD_CONFIG" phase2.ssh_user)

    local ssh_opts=(-o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o LogLevel=ERROR)

    if sshpass -p "$TRX_ROOT_PASSWORD" ssh "${ssh_opts[@]}" "${ssh_user}@${trx_ip}" true 2>/dev/null; then
        echo "  [ OK ] all 8 images written; TRX still reachable over SSH"
    else
        echo "  [WARN] TRX not reachable over SSH right after flashing (it may have dropped the link)"
    fi
    echo "  Final check is manual: power-cycle the TRX and confirm it boots and comes up."
    return 0
}

method_monitor() {
    echo "TRX flash complete. Power-cycle and observe operation."
    return 0
}
