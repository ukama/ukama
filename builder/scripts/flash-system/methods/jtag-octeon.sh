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

# TRX_VERBOSE=0 silences the serial/oct-remote-boot tails
TRX_VERBOSE="${TRX_VERBOSE:-1}"

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
    sudo pkill -9 -f '[o]ct-remote-boot' 2>/dev/null || true
    if [ -n "$SPAM_PID" ]; then
        kill "$SPAM_PID" 2>/dev/null || true
        SPAM_PID=""
    fi
    [ -n "$TFTP_STAGE_DIR" ] && sudo rm -rf "$TFTP_STAGE_DIR"
}

bdi_telnet_cmd() {
    local host="$1"
    local cmd="$2"
    if ! command -v expect >/dev/null 2>&1; then
        echo "WARNING: expect not installed - cannot send BDI command '$cmd'" >&2
        return 1
    fi
    expect -c "
        set timeout 300
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
        # wait for the prompt before quit, otherwise they merge and the command never runs
        expect {
            \"cnMIPS#0>\" {}
            \"Core#0>\"   {}
            timeout      {}
        }
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

    local rc_post_src
    rc_post_src="$(dirname "$BOARD_CONFIG")/payloads/rc_post.local"
    if [ ! -f "$rc_post_src" ]; then
        echo "  [FAIL] rc_post.local payload not found: $rc_post_src"
        fail=1
    else
        echo "  [ OK ] rc_post.local payload: $rc_post_src"
    fi

    local post_trx_ip post_host_ip
    post_trx_ip=$(yq_read "$BOARD_CONFIG" network.post_flash.trx_ip)
    post_host_ip=$(yq_read "$BOARD_CONFIG" network.post_flash.host_ip)
    if [ -z "$post_trx_ip" ] || [ "$post_trx_ip" = "null" ] || [ -z "$post_host_ip" ] || [ "$post_host_ip" = "null" ]; then
        echo "  [FAIL] network.post_flash.trx_ip/host_ip not configured"
        fail=1
    else
        echo "  [ OK ] post-flash subnet: host ${post_host_ip}, TRX ${post_trx_ip}"
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
    echo "  Phase 2 (SSH)   : scp+dd 8 .img files to /dev/flash_*, install band.cfg + rc_post.local"
    echo "  Post-flash      : TRX moves to $(yq_read "$BOARD_CONFIG" network.post_flash.trx_ip); bench box -> $(yq_read "$BOARD_CONFIG" network.post_flash.host_ip)"
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

    # exact file sizes for the bootcbyflash cp.b commands
    local os_path rd_path os_size rd_size os_size_hex rd_size_hex
    os_path=$(yq_read "$BOARD_CONFIG" phase1.artifacts.os.path)
    rd_path=$(yq_read "$BOARD_CONFIG" phase1.artifacts.rd.path)
    os_size=$(stat -c %s "$os_path" 2>/dev/null || yq_read "$BOARD_CONFIG" phase1.artifacts.os.size)
    rd_size=$(stat -c %s "$rd_path" 2>/dev/null || yq_read "$BOARD_CONFIG" phase1.artifacts.rd.size)
    os_size_hex=$(printf '0x%x' "$os_size")
    rd_size_hex=$(printf '0x%x' "$rd_size")

    local restore_errexit=0
    case $- in *e*) restore_errexit=1; set +e ;; esac

    uboot_send_and_wait "$dev" "setenv ethaddr ${mac_dashed}" "$prompt" 300
    uboot_send_and_wait "$dev" "setenv ipaddr ${trx_ip}" "$prompt" 300
    uboot_send_and_wait "$dev" "setenv netmask ${netmask}" "$prompt" 300
    uboot_send_and_wait "$dev" "setenv serverip ${host_ip}" "$prompt" 300
    uboot_send_and_wait "$dev" "setenv bootby flash" "$prompt" 300
    uboot_send_and_wait "$dev" "setenv cfgloadby flash" "$prompt" 300
    uboot_send_and_wait "$dev" "setenv swloadby flash" "$prompt" 300
    uboot_send_and_wait "$dev" 'setenv i2cinit "i2c dev 0; i2c probe; i2c dev 1; i2c probe"' "$prompt" 300
    uboot_send_and_wait "$dev" 'setenv bootcmd "mw64 0x00011800B0001000 0x0140; run i2cinit; run namedalloc; run bootcby${bootby}"' "$prompt" 300
    uboot_send_and_wait "$dev" 'setenv bootcbytftp "tftp 0x21000000 lsm_os_trx.gz; gunzip 0x21000000 0x20000000 0x1000000; tftp 0x30800000 lsm_rd_trx.gz; bootoctlinux 0x20000000 coremask=0x7 endbootargs rd_name=initrd mem=512M;"' "$prompt" 300
    local bootcbyflash_cmd
    bootcbyflash_cmd="setenv bootcbyflash \"cp.b 0x17E20000 0x21000000 ${os_size_hex}; gunzip 0x21000000 0x20000000 0x1000000; cp.b 0x18320000 0x30800000 ${rd_size_hex}; bootoctlinux 0x20000000 coremask=0x7 endbootargs rd_name=initrd mem=512M;\""
    uboot_send_and_wait "$dev" "$bootcbyflash_cmd" "$prompt" 300
    uboot_send_and_wait "$dev" 'setenv namedalloc "namedalloc dsp-dump 0x400000 0x7f4D0000; namedalloc cazac 0x630000 0x7f8D0000; namedalloc cpu-dsp-if 0x100000 0x7ff00000; namedalloc dsp-log-buf 0x4000000 0x80000000; namedalloc initrd 0x2800000 0x30800000;"' "$prompt" 300
    uboot_send_and_wait "$dev" "setenv mk_ubootenv 1" "$prompt" 300
    # SGMII autoneg must be on before Linux boots or octeth0 stays down
    uboot_send_and_wait "$dev" "setenv preboot 'mw64 0x00011800B0001000 0x0140'" "$prompt" 300

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

    uboot_send_and_wait "$dev" "tftp ${ddr_addr} ${name}" "$prompt" 300

    if ! tail -c +"$marker_before" "$UBOOT_LOG" 2>/dev/null | grep -q "Bytes transferred = "; then
        echo "ERROR: tftp failed for $key - no 'Bytes transferred' seen in log"
        return 1
    fi

    uboot_send_and_wait "$dev" "protect off all" "$prompt" 300
    uboot_send_and_wait "$dev" "erase ${flash_addr} +\${filesize}" "$prompt" 300
    uboot_send_and_wait "$dev" "cp.b ${ddr_addr} ${flash_addr} \${filesize}" "$prompt" 300
}

_phase1_run() {
    local bdi_ip serial_dev baud uboot_prompt oct_path oct_board oct_clock
    local ddr_os ddr_rd gdb_port oct_env_root oct_env_protocol host_ip

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
    host_ip=$(yq_read "$BOARD_CONFIG" network.host_ip)

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

    # only send CONFIG to a bare BDI (Core#0>); reloading a configured one hangs the next DDR init
    echo "Checking BDI state at ${bdi_ip} (cnMIPS#0> = configured, Core#0> = bare)..."
    local bdi_state
    bdi_state=$(expect -c "
        set timeout 300
        log_user 0
        spawn telnet $bdi_ip
        expect {
            \"cnMIPS#0>\" { puts BDI_READY }
            \"Core#0>\"   { puts BDI_BARE }
            timeout       { puts BDI_DOWN }
        }
        catch { send \"quit\r\"; expect eof }
    " 2>/dev/null | grep -oE 'BDI_READY|BDI_BARE|BDI_DOWN' | tail -1)

    if [ "$bdi_state" = "BDI_READY" ]; then
        echo "  BDI already configured - skipping CONFIG reload (matches the manual flow)."
    elif [ "$bdi_state" = "BDI_BARE" ]; then
        echo "  BDI is unconfigured - sending HOST + CONFIG (the BDI will reboot and auto-load)..."
        if ! expect -c "
            set timeout 300
            spawn telnet $bdi_ip
            expect {
                \"Core#0>\" {}
                timeout { puts \"BDI telnet timeout\"; exit 1 }
            }
            send \"HOST $host_ip\r\"
            expect \"Core#0>\"
            send \"CONFIG cnf71xx.cfg\r\"
            expect {
                \"configuration passed\" { puts \"Config load succeeded.\"; exit 0 }
                \"cannot open\" { puts \"Config load failed (file not found on TFTP).\"; exit 1 }
                timeout { puts \"Config load timeout (TFTP slow or unreachable).\"; exit 1 }
            }
        " 2>/dev/null; then
            echo "ERROR: BDI config load failed."
            echo "  Please manually telnet to the BDI and run: HOST $host_ip then CONFIG cnf71xx.cfg"
            return 1
        fi
        echo "  Config sent - BDI is rebooting and auto-loading cnf71xx.cfg over TFTP."
        echo "  Keeping the TFTP server up and polling until the BDI comes up configured..."
        # keep TFTP up while the BDI reboots and auto-loads the config
        local cfg_wait=0 bdi_now=""
        while [ "$cfg_wait" -lt 300 ]; do
            sleep 10
            cfg_wait=$((cfg_wait + 10))
            bdi_now=$(expect -c "
                set timeout 300
                log_user 0
                spawn telnet $bdi_ip
                expect {
                    \"cnMIPS#0>\" { puts BDI_READY }
                    \"Core#0>\"   { puts BDI_BARE }
                    timeout       { puts BDI_DOWN }
                }
                catch { send \"quit\r\"; expect eof }
            " 2>/dev/null | grep -oE 'BDI_READY|BDI_BARE|BDI_DOWN' | tail -1)
            if [ "$bdi_now" = "BDI_READY" ]; then
                echo "  BDI came up configured (cnMIPS#0>) after ${cfg_wait}s."
                break
            fi
            echo "  ...BDI not ready yet (${bdi_now:-no response}, ${cfg_wait}s)"
        done
        if [ "$bdi_now" != "BDI_READY" ]; then
            echo "ERROR: BDI did not auto-load its config (still ${bdi_now:-unknown} after ${cfg_wait}s)."
            echo "  Almost always means CNF71XX.cfg (+ its .def) is missing so TFTP auto-load fails."
            echo "  Confirm the file exists at: $bdi_config_src"
            return 1
        fi
    else
        echo "ERROR: could not reach the BDI telnet prompt at ${bdi_ip}."
        echo "  Check BDI power and network, then re-run."
        return 1
    fi

    echo "  Probing BDI GDB port ${bdi_ip}:2001..."
    local gdb_wait=0
    while ! nc -z "$bdi_ip" 2001 2>/dev/null; do
        gdb_wait=$((gdb_wait + 5))
        if [ "$gdb_wait" -ge 300 ]; then
            echo "ERROR: BDI GDB port still closed after 60s - the BDI may not have auto-loaded"
            echo "  its config (TFTP must be running when the BDI boots)."
            echo "  Power-cycle the BDI now (TFTP is still being served) and re-run this script."
            return 1
        fi
        sleep 5
    done
    echo "  GDB port is open."

    # bring-up: go 0x400000 once, then oct-remote-boot; on failure re-run only oct-remote-boot
    local oct_log="${LOG_DIR}/oct-remote-boot.log"
    local prompt_seen=0

    echo "Sending 'go 0x400000' via BDI telnet (once)..."
    if ! bdi_telnet_cmd "$bdi_ip" "go 0x400000"; then
        echo "ERROR: BDI 'go 0x400000' failed (is the BDI at cnMIPS#0>?)."
        return 1
    fi
    sleep 5

    echo "Opening serial console at $serial_dev ($baud)..."
    uboot_open "$serial_dev" "$baud" "${LOG_DIR}/uboot.log"

    local serial_tail_pid=""
    if [ "$TRX_VERBOSE" = "1" ]; then
        echo "(verbose) Tailing serial log to console..."
        tail -f "${LOG_DIR}/uboot.log" &
        serial_tail_pid=$!
    fi

    echo "Spamming serial with key presses to stop zero-second autoboot..."
    (
        # autoboot delay is zero, keep a key held down
        exec 3>"$serial_dev"
        while true; do
            printf ' ' >&3
            sleep 0.03
        done
    ) &
    SPAM_PID=$!

    # the DDR clock lock varies per attempt; only ~400 MHz gives a working SGMII clock
    local oct_try=0 max_oct_tries=8
    local elapsed clk mhz oct_exit ddr_mislock prompt_this
    while [ "$oct_try" -lt "$max_oct_tries" ]; do
        oct_try=$((oct_try + 1))
        echo ""
        echo "=== oct-remote-boot attempt ${oct_try}/${max_oct_tries} (no re-go, no power-cycle) ==="

        sudo pkill -9 -f '[o]ct-remote-boot' 2>/dev/null || true
        sleep 1

        echo "Starting oct-remote-boot (OCTEON_ROOT=$oct_env_root, $oct_env_protocol)..."
        : > "$oct_log"
        sudo stdbuf -oL -eL env OCTEON_ROOT="$oct_env_root" OCTEON_REMOTE_PROTOCOL="$oct_env_protocol" \
            "$oct_path" --board="$oct_board" --ddr_clock_hz="$oct_clock" >"$oct_log" 2>&1 &
        REMOTE_BOOT_PID=$!
        disown "$REMOTE_BOOT_PID" 2>/dev/null || true
        echo "  oct-remote-boot started (PID $REMOTE_BOOT_PID)"

        if [ "$TRX_VERBOSE" = "1" ]; then
            echo "(verbose) Tailing oct-remote-boot log to console..."
            tail -f "$oct_log" &
            OCT_TAIL_PID=$!
        fi

        # prompt can be "Octeon zen(ram)=>" or "Octeon zen(Failsafe)=>"
        echo "Waiting for u-boot prompt 'Octeon zen…=>' with DDR ~400 MHz (up to 300s)..."
        elapsed=0; clk=""; mhz=""; ddr_mislock=0; prompt_this=0
        while [ "$elapsed" -lt 300 ]; do
            clk=$(grep -a "Measured DDR clock" "$oct_log" 2>/dev/null | tail -1 || true)
            mhz=$(printf '%s' "$clk" | grep -oE '[0-9]+' | head -1 || true)
            if [ -n "$mhz" ] && { [ "$mhz" -lt 380 ] || [ "$mhz" -gt 420 ]; }; then
                ddr_mislock=1
                break
            fi
            if grep -qE "Octeon zen.*=>" "${LOG_DIR}/uboot.log" 2>/dev/null; then
                prompt_this=1
                break
            fi
            # oct-remote-boot exited; u-boot may still be printing
            if ! kill -0 "$REMOTE_BOOT_PID" 2>/dev/null; then
                local grace=0
                while [ "$grace" -lt 300 ]; do
                    if grep -qE "Octeon zen.*=>" "${LOG_DIR}/uboot.log" 2>/dev/null; then
                        prompt_this=1
                        break
                    fi
                    sleep 1
                    grace=$((grace + 1))
                done
                break
            fi
            sleep 1
            elapsed=$((elapsed + 1))
        done

        if [ "$prompt_this" -eq 1 ] && [ -n "$mhz" ] && [ "$mhz" -ge 380 ] && [ "$mhz" -le 420 ]; then
            # oct-remote-boot stays running; it hosts u-boot over GDB while we flash
            prompt_seen=1
            echo "  u-boot prompt + good DDR clock (${mhz} MHz) on attempt ${oct_try}."
            if [ -n "$OCT_TAIL_PID" ]; then
                kill "$OCT_TAIL_PID" 2>/dev/null || true
                OCT_TAIL_PID=""
            fi
            break
        fi

        if kill -0 "$REMOTE_BOOT_PID" 2>/dev/null; then
            sudo kill "$REMOTE_BOOT_PID" 2>/dev/null || true
        fi
        wait "$REMOTE_BOOT_PID" 2>/dev/null
        oct_exit=$?
        REMOTE_BOOT_PID=""

        if [ "$ddr_mislock" -eq 1 ]; then
            echo "  DDR mislocked at ${mhz} MHz (need ~400) - re-running to re-roll the DDR clock."
        elif [ "$prompt_this" -eq 1 ]; then
            echo "  u-boot came up but DDR clock unknown/!~400 - re-running."
        elif [ "$oct_exit" -eq 139 ] || grep -qaE "Segmentation fault|GDB Reply Error|in reset, told to continue" "$oct_log" 2>/dev/null; then
            echo "  oct-remote-boot exited (GDB error / segfault) - re-running it (normal recovery)."
        else
            echo "  no u-boot prompt within 300s - re-running."
        fi
        echo "--- last 15 lines of oct-remote-boot output ---"
        tail -n 15 "$oct_log" 2>/dev/null | sed 's/^/    /' || true
        sleep 2
    done

    if [ -n "$SPAM_PID" ]; then
        kill "$SPAM_PID" 2>/dev/null || true
        SPAM_PID=""
    fi

    if [ "$prompt_seen" -ne 1 ]; then
        if [ -n "$serial_tail_pid" ]; then
            kill "$serial_tail_pid" 2>/dev/null || true
            serial_tail_pid=""
        fi
        uboot_close
        echo "ERROR: u-boot prompt did not appear after ${max_oct_tries} oct-remote-boot attempts."
        echo "  If oct-remote-boot kept segfaulting, the BDI GDB stub may be wedged: cold"
        echo "  power-cycle BOTH the TRX and the BDI, let the BDI settle, then re-run."
        echo "--- last 20 lines of serial (uboot.log) ---"
        tail -n 20 "${LOG_DIR}/uboot.log" 2>/dev/null | sed 's/^/    /' || true
        return 1
    fi

    echo "  u-boot prompt reached."
    echo "Draining residual serial output (spam backlog) before sending commands..."
    uboot_drain 3

    echo "Pushing u-boot environment variables..."
    _phase1_uboot_env "$serial_dev" "$uboot_prompt"

    echo "Enabling ethernet (mw64 x2) and saving env..."
    uboot_send_and_wait "$serial_dev" "mw64 0x00011800B0001000 0x0140" "$uboot_prompt" 300 || true
    uboot_send_and_wait "$serial_dev" "mw64 0x00011800B0001000 0x0140" "$uboot_prompt" 300 || true
    uboot_send_and_wait "$serial_dev" "saveenv" "$uboot_prompt" 300 || true

    local ping_marker link_up=0 ping_attempt=0 max_ping_attempts=6
    echo "Bringing up ethernet link to ${host_ip} (mw64 + ping retries)..."
    while [ "$ping_attempt" -lt "$max_ping_attempts" ]; do
        ping_attempt=$((ping_attempt + 1))
        ping_marker=$(wc -c < "$UBOOT_LOG" 2>/dev/null || echo 0)
        uboot_send_and_wait "$serial_dev" "ping ${host_ip}" "$uboot_prompt" 300 || true
        if tail -c +"$ping_marker" "$UBOOT_LOG" 2>/dev/null | grep -q "is alive"; then
            link_up=1
            break
        fi
        echo "  ping attempt ${ping_attempt}/${max_ping_attempts} failed - re-enabling ethernet (mw64) and retrying..."
        uboot_send_and_wait "$serial_dev" "mw64 0x00011800B0001000 0x0140" "$uboot_prompt" 300 || true
        sleep 2
    done

    if [ "$link_up" -ne 1 ]; then
        echo "ERROR: TRX ethernet link did not come up after ${max_ping_attempts} attempts (octeth0 Down)."
        echo "  mw64 ethernet-enable + ping kept failing. Usual cause is a mislocked DDR clock"
        echo "  (SGMII reference clock off) - cold power-cycle the TRX and re-run; also check the cable."
        grep -a "Measured DDR clock" "$oct_log" 2>/dev/null | tail -1 | sed 's/^/  oct-remote-boot: /' || true
        return 1
    fi
    echo "  host ${host_ip} is reachable - ethernet link up after ${ping_attempt} attempt(s)"

    _phase1_flash_artifact "$serial_dev" "$uboot_prompt" "os"    "$ddr_os"
    _phase1_flash_artifact "$serial_dev" "$uboot_prompt" "rd"    "$ddr_rd"
    _phase1_flash_artifact "$serial_dev" "$uboot_prompt" "uboot" "$ddr_os"

    if [ -n "$serial_tail_pid" ]; then
        kill "$serial_tail_pid" 2>/dev/null || true
        serial_tail_pid=""
    fi
    uboot_close
    tftp_stop
    [ -n "$REMOTE_BOOT_PID" ] && sudo kill "$REMOTE_BOOT_PID" 2>/dev/null || true
    REMOTE_BOOT_PID=""
}

# A freshly-erased board boots with ethernet down (no rc_post.local yet); log in
# over serial and enable it with devmem. Best-effort, falls back to the SSH wait.
_phase2_enable_ethernet_over_serial() {
    local serial_dev baud
    serial_dev=$(yq_read "$BOARD_CONFIG" serial.device)
    baud=$(yq_read "$BOARD_CONFIG" serial.baud)

    local holder=""
    if command -v lsof >/dev/null 2>&1; then
        holder=$(lsof -t "$serial_dev" 2>/dev/null | head -1 || true)
    fi
    if [ -n "$holder" ]; then
        echo "  NOTE: $serial_dev is held by another process (PID $holder) - skipping auto"
        echo "        ethernet-enable. Close PuTTY/screen, or on the console run twice:"
        echo "          devmem 0x00011800B0001000 64 0x0140"
        return 0
    fi

    echo "Enabling TRX ethernet over serial (login as root, devmem x2)..."
    uboot_open "$serial_dev" "$baud" "${LOG_DIR}/phase2-serial.log" || return 0

    local phase2_tail_pid=""
    if [ "$TRX_VERBOSE" = "1" ]; then
        echo "(verbose) Tailing phase2 serial log to console..."
        tail -f "${LOG_DIR}/phase2-serial.log" &
        phase2_tail_pid=$!
    fi

    # only match output newer than each send; stale prompts in the log must not count
    local login_attempt logged_in=0 marker
    for login_attempt in 1 2 3; do
        marker=$(uboot_log_mark)
        uboot_send "$serial_dev" ""
        if uboot_wait_for_from "$marker" "~ #" 8; then
            logged_in=1
            break
        fi
        if ! uboot_wait_for_from "$marker" "login:" 30; then
            echo "  no fresh login prompt on attempt ${login_attempt}/3 (still booting?); retrying..."
            continue
        fi
        marker=$(uboot_log_mark)
        uboot_send "$serial_dev" "root"
        if ! uboot_wait_for_from "$marker" "assword" 20; then
            echo "  no password prompt on attempt ${login_attempt}/3; retrying..."
            continue
        fi
        marker=$(uboot_log_mark)
        uboot_send "$serial_dev" "${TRX_ROOT_PASSWORD:-cavium.lte}"
        if uboot_wait_for_from "$marker" "~ #" 20; then
            logged_in=1
            break
        fi
        echo "  login attempt ${login_attempt}/3 didn't reach a shell prompt; retrying..."
    done

    if [ "$logged_in" -ne 1 ]; then
        echo "  WARNING: serial login didn't reach a shell prompt; relying on SSH wait."
        if [ -n "$phase2_tail_pid" ]; then
            kill "$phase2_tail_pid" 2>/dev/null || true
            phase2_tail_pid=""
        fi
        uboot_close
        return 0
    fi

    uboot_send "$serial_dev" "devmem 0x00011800B0001000 64 0x0140"
    sleep 1
    uboot_send "$serial_dev" "devmem 0x00011800B0001000 64 0x0140"
    sleep 1
    echo "  ethernet-enable (devmem x2) sent over serial."
    if [ -n "$phase2_tail_pid" ]; then
        kill "$phase2_tail_pid" 2>/dev/null || true
        phase2_tail_pid=""
    fi
    uboot_close
    return 0
}

_phase2_remount_app() {
    local trx_ip="$1" ssh_user="$2"
    local sshpass_args=(-p "$TRX_ROOT_PASSWORD")
    local ssh_opts=(-o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o LogLevel=ERROR)

    echo "Remounting /mnt/app to pick up the freshly-flashed app image..."
    local app_dev
    app_dev=$(sshpass "${sshpass_args[@]}" ssh "${ssh_opts[@]}" "${ssh_user}@${trx_ip}" "mount | grep ' /mnt/app ' | awk '{print \$1}'") || app_dev=""
    if [ -z "$app_dev" ]; then
        app_dev="/dev/flash_app0"
        echo "  WARNING: could not parse /mnt/app device; falling back to ${app_dev}."
    fi

    # sshd keeps its host keys under /mnt/app, so a plain umount is busy; lazy-unmount first
    local out
    out=$(sshpass "${sshpass_args[@]}" ssh "${ssh_opts[@]}" "${ssh_user}@${trx_ip}" "
        sync
        umount /mnt/app 2>/dev/null || true
        umount -l /mnt/app 2>/dev/null || true
        if mount -t jffs2 ${app_dev} /mnt/app >/dev/null 2>&1; then
            if touch /mnt/app/.probe 2>/dev/null && rm -f /mnt/app/.probe 2>/dev/null; then
                echo MOUNT_OK
            else
                echo MOUNT_RO
            fi
        else
            echo MOUNT_FAIL
        fi
    ") || out=""

    if printf '%s\n' "$out" | grep -q '^MOUNT_OK$'; then
        echo "  remounted /mnt/app from ${app_dev} and verified writable."
        return 0
    fi

    echo "  WARNING: /mnt/app remount did not verify writable."
    echo "           Diagnostics from TRX:"
    sshpass "${sshpass_args[@]}" ssh "${ssh_opts[@]}" "${ssh_user}@${trx_ip}" "
        echo '--- /proc/mounts ---'; grep /mnt/app /proc/mounts || true
        echo '--- dmesg tail ---'; dmesg | tail -n 30 || true
        echo '--- open files on /mnt/app ---'
        if command -v lsof >/dev/null 2>&1; then lsof /mnt/app 2>/dev/null || true; else fuser -mv /mnt/app 2>/dev/null || true; fi
    " || true
    return 1
}

# Inject rc_post.local into flash_app0.img on the bench box before it is dd'd.
_phase2_inject_rc_post_into_app0_image() {
    local app0_image="$1" rc_post_src="$2" rc_post_target="$3"
    local internal_path="${rc_post_target#/mnt/app}"
    internal_path="${internal_path#/}"

    echo "Injecting rc_post.local into ${app0_image} on the bench box (pre-dd)..."

    # jffs2 only mounts on MTD devices, so stage the image into mtdram (128 KiB erase size)
    local erase_kib=128 img_bytes total_kib
    img_bytes=$(stat -c %s "$app0_image" 2>/dev/null || echo 0)
    if [ "$img_bytes" -eq 0 ] || [ $((img_bytes % (erase_kib * 1024))) -ne 0 ]; then
        echo "  WARNING: app0 image size ${img_bytes} is not a multiple of ${erase_kib}KiB; will try post-flash copy instead."
        return 1
    fi
    total_kib=$((img_bytes / 1024))

    modprobe jffs2 >/dev/null 2>&1 || true
    modprobe -r mtdram >/dev/null 2>&1 || true
    if ! modprobe mtdram total_size="$total_kib" erase_size="$erase_kib" >/dev/null 2>&1; then
        echo "  WARNING: mtdram kernel module unavailable; will try post-flash copy instead."
        return 1
    fi
    modprobe mtdblock >/dev/null 2>&1 || true

    local mtd_idx mtd_dev
    mtd_idx=$(grep -i mtdram /proc/mtd 2>/dev/null | head -1 | sed 's/^mtd\([0-9]*\):.*/\1/')
    mtd_dev="/dev/mtdblock${mtd_idx}"
    if [ -z "$mtd_idx" ] || [ ! -e "$mtd_dev" ]; then
        echo "  WARNING: no mtdblock device for mtdram; will try post-flash copy instead."
        modprobe -r mtdram >/dev/null 2>&1 || true
        return 1
    fi

    local tmpdir mntdir
    tmpdir=$(mktemp -d) || { modprobe -r mtdram >/dev/null 2>&1 || true; return 1; }
    mntdir="${tmpdir}/mnt"
    mkdir -p "$mntdir"

    _inject_cleanup() {
        umount "$mntdir" >/dev/null 2>&1 || true
        rm -rf "$tmpdir"
        modprobe -r mtdram >/dev/null 2>&1 || true
    }

    if ! dd if="$app0_image" of="$mtd_dev" bs=1M >/dev/null 2>&1; then
        echo "  WARNING: could not load app0 image into mtdram; will try post-flash copy instead."
        _inject_cleanup
        return 1
    fi

    if ! mount -t jffs2 "$mtd_dev" "$mntdir" >/dev/null 2>&1; then
        echo "  WARNING: could not mount app0 image as jffs2 (via mtdram); will try post-flash copy instead."
        _inject_cleanup
        return 1
    fi

    mkdir -p "${mntdir}/$(dirname "$internal_path")"
    if ! cp "$rc_post_src" "${mntdir}/${internal_path}"; then
        echo "  WARNING: could not copy rc_post.local into app0 image; will try post-flash copy instead."
        _inject_cleanup
        return 1
    fi
    chmod +x "${mntdir}/${internal_path}"
    sync
    umount "$mntdir" >/dev/null 2>&1 || true

    [ -f "${app0_image}.orig" ] || cp "$app0_image" "${app0_image}.orig"
    if ! dd if="$mtd_dev" of="${tmpdir}/app0.new" bs=1M >/dev/null 2>&1; then
        echo "  WARNING: could not dump edited app0 image from mtdram; will try post-flash copy instead."
        _inject_cleanup
        return 1
    fi
    if [ "$(stat -c %s "${tmpdir}/app0.new")" != "$img_bytes" ]; then
        echo "  WARNING: edited app0 image size mismatch; keeping original, will try post-flash copy instead."
        _inject_cleanup
        return 1
    fi
    mv "${tmpdir}/app0.new" "$app0_image"
    _inject_cleanup

    echo "  injected rc_post.local into app0 image; it will be flashed by dd along with the other images."
    echo "  (original image saved as ${app0_image}.orig)"
    return 0
}

_phase2_run() {
    local trx_ip ssh_user staging
    trx_ip=$(yq_read "$BOARD_CONFIG" network.trx_ip)
    ssh_user=$(yq_read "$BOARD_CONFIG" phase2.ssh_user)
    staging=$(yq_read "$BOARD_CONFIG" phase2.ssh_staging_dir)

    local sshpass_args=(-p "$TRX_ROOT_PASSWORD")
    local ssh_opts=(-o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o LogLevel=ERROR)

    local rc_post_src rc_post_target app0_image
    rc_post_src="$(dirname "$BOARD_CONFIG")/payloads/rc_post.local"
    rc_post_target=$(yq_read "$BOARD_CONFIG" phase2.rc_post_local)
    app0_image=$(yq_read "$BOARD_CONFIG" phase2.images.app0.src)

    local rc_post_injected=0
    if [ -f "$rc_post_src" ] && [ -f "$app0_image" ]; then
        if _phase2_inject_rc_post_into_app0_image "$app0_image" "$rc_post_src" "$rc_post_target"; then
            rc_post_injected=1
        fi
    fi

    _phase2_enable_ethernet_over_serial

    echo "Waiting for TRX SSH at ${trx_ip}..."
    local elapsed=0
    while [ "$elapsed" -lt 300 ]; do
        if sshpass "${sshpass_args[@]}" ssh "${ssh_opts[@]}" "${ssh_user}@${trx_ip}" true 2>/dev/null; then
            break
        fi
        sleep 2
        elapsed=$((elapsed + 2))
    done
    if [ "$elapsed" -ge 300 ]; then
        echo "ERROR: TRX not reachable via SSH within 300s."
        echo "  Ethernet likely still down. On the TRX serial console (root / cavium.lte) run"
        echo "  twice:  devmem 0x00011800B0001000 64 0x0140   then re-run the flash."
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
            "dd if=${staging}/${name} of=${dst} bs=1M && rm -f ${staging}/${name}"
    done

    echo ""
    echo "Syncing TRX filesystems..."
    sshpass "${sshpass_args[@]}" ssh "${ssh_opts[@]}" "${ssh_user}@${trx_ip}" "sync"
    echo "  sync done."

    echo "Inspecting /mnt/app mount source:"
    sshpass "${sshpass_args[@]}" ssh "${ssh_opts[@]}" "${ssh_user}@${trx_ip}" \
        "mount | grep ' /mnt/app ' || true; df -h /mnt/app || true"

    # Never write to /mnt/app here: the flash under the live mount was just replaced
    # by dd and writes through it corrupt the new image. Config installs happen after
    # the power-cycle in _phase2_post_flash_config.
    if [ "$rc_post_injected" -eq 1 ]; then
        echo "rc_post.local was pre-injected into flash_app0.img."
    else
        echo "rc_post.local ships inside flash_app0.img; nothing to copy before the reboot."
    fi

    echo ""
    echo "All 8 images written."
    echo "Next: power-cycle the TRX to boot from the newly flashed images."
    echo "Band config is installed after the reboot, once /mnt/app is cleanly mounted."
}

# Install band config (and rc_post.local if missing) after the final power-cycle,
# once /mnt/app is cleanly mounted. Reboots the TRX if anything changed.
_phase2_post_flash_config() {
    local trx_ip ssh_user
    if yq_exists "$BOARD_CONFIG" network.post_flash.trx_ip; then
        trx_ip=$(yq_read "$BOARD_CONFIG" network.post_flash.trx_ip)
    else
        trx_ip=$(yq_read "$BOARD_CONFIG" network.trx_ip)
    fi
    ssh_user=$(yq_read "$BOARD_CONFIG" phase2.ssh_user)

    local sshpass_args=(-p "$TRX_ROOT_PASSWORD")
    local ssh_opts=(-o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o LogLevel=ERROR)

    local band_default band_configs_dir band_cfg_src band_cfg_target rc_post_src rc_post_target
    band_default=$(yq_read "$BOARD_CONFIG" band.default)
    band_configs_dir=$(yq_read "$BOARD_CONFIG" band.configs_dir)
    band_cfg_src="${band_configs_dir}/${band_default}.cfg"
    band_cfg_target=$(yq_read "$BOARD_CONFIG" band.target_path)
    rc_post_src="$(dirname "$BOARD_CONFIG")/payloads/rc_post.local"
    rc_post_target=$(yq_read "$BOARD_CONFIG" phase2.rc_post_local)

    echo "Waiting for TRX SSH at ${trx_ip} (booted from flash)..."
    local elapsed=0
    while [ "$elapsed" -lt 300 ]; do
        if sshpass "${sshpass_args[@]}" ssh "${ssh_opts[@]}" "${ssh_user}@${trx_ip}" true 2>/dev/null; then
            break
        fi
        sleep 2
        elapsed=$((elapsed + 2))
    done
    if [ "$elapsed" -ge 300 ]; then
        echo "ERROR: TRX not reachable at ${trx_ip} after the final power-cycle."
        echo "  Wait for it to finish booting, then re-check; band config not installed."
        return 1
    fi
    echo "  TRX is up at ${trx_ip}."

    local reboot_needed=0

    # rc_post.local normally ships inside the app image; only copy if missing
    if sshpass "${sshpass_args[@]}" ssh "${ssh_opts[@]}" "${ssh_user}@${trx_ip}" "test -f ${rc_post_target}"; then
        echo "  rc_post.local already present at ${rc_post_target} (from the app image)."
    else
        echo "  rc_post.local missing - installing ${rc_post_src} -> ${rc_post_target}..."
        cat "$rc_post_src" | sshpass "${sshpass_args[@]}" ssh "${ssh_opts[@]}" "${ssh_user}@${trx_ip}" \
            "cat > ${rc_post_target} && chmod +x ${rc_post_target}" || {
            echo "ERROR: rc_post.local install failed."
            return 1
        }
        reboot_needed=1
    fi

    # skip the write (and the reboot) if the target already matches
    if [ -f "$band_cfg_src" ]; then
        local local_md5 remote_md5
        local_md5=$(md5sum "$band_cfg_src" | awk '{print $1}')
        remote_md5=$(sshpass "${sshpass_args[@]}" ssh "${ssh_opts[@]}" "${ssh_user}@${trx_ip}" \
            "md5sum ${band_cfg_target} 2>/dev/null | cut -d' ' -f1" || true)
        if [ "$local_md5" = "$remote_md5" ]; then
            echo "  Band config at ${band_cfg_target} already matches ${band_default}."
        else
            echo "  Installing band config (${band_cfg_src} -> ${band_cfg_target})..."
            if ! cat "$band_cfg_src" | sshpass "${sshpass_args[@]}" ssh "${ssh_opts[@]}" "${ssh_user}@${trx_ip}" "cat > ${band_cfg_target}"; then
                echo "ERROR: band config install failed."
                return 1
            fi
            reboot_needed=1
        fi
    else
        echo "WARNING: band config source not found at ${band_cfg_src}"
    fi

    sshpass "${sshpass_args[@]}" ssh "${ssh_opts[@]}" "${ssh_user}@${trx_ip}" "sync"

    if [ "$reboot_needed" -eq 1 ]; then
        echo "  Rebooting the TRX so the new config takes effect..."
        sshpass "${sshpass_args[@]}" ssh "${ssh_opts[@]}" "${ssh_user}@${trx_ip}" "reboot" 2>/dev/null || true
        echo "  (If it doesn't come back within ~5 minutes, power-cycle it manually.)"
        sleep 20
    fi
    return 0
}

# The flashed TRX comes up on its production subnet; add a matching IP on the bench NIC.
_phase2_rehost_for_post_flash() {
    local post_trx_ip post_host_ip post_netmask post_iface pre_trx_ip
    post_trx_ip=$(yq_read "$BOARD_CONFIG" network.post_flash.trx_ip)
    post_host_ip=$(yq_read "$BOARD_CONFIG" network.post_flash.host_ip)
    post_netmask=$(yq_read "$BOARD_CONFIG" network.post_flash.netmask)
    post_iface=$(yq_read "$BOARD_CONFIG" network.post_flash.interface)

    if [ -z "$post_trx_ip" ] || [ "$post_trx_ip" = "null" ]; then
        echo "  No network.post_flash.trx_ip configured; skipping host IP reconfiguration."
        return 0
    fi
    if [ -z "$post_host_ip" ] || [ "$post_host_ip" = "null" ]; then
        return 0
    fi

    if [ -z "$post_iface" ] || [ "$post_iface" = "null" ]; then
        pre_trx_ip=$(yq_read "$BOARD_CONFIG" network.trx_ip)
        post_iface=$(ip route get "$pre_trx_ip" 2>/dev/null | awk '{for(i=1;i<=NF;i++) if($i=="dev") {print $(i+1); exit}}')
    fi
    if [ -z "$post_iface" ] || [ "$post_iface" = "null" ]; then
        echo "  WARNING: could not detect host interface for post-flash subnet ${post_host_ip}; skipping rehost."
        return 0
    fi

    local cidr
    cidr=$(python3 -c "import ipaddress; print(ipaddress.IPv4Network('0.0.0.0/${post_netmask}').prefixlen)" 2>/dev/null || echo "24")

    echo "Rehosting bench box to post-flash subnet: ${post_iface} -> ${post_host_ip}/${cidr}"
    ip addr add "${post_host_ip}/${cidr}" dev "$post_iface" 2>/dev/null || true
    echo "  Bench box now has ${post_host_ip}/${cidr} on ${post_iface}."
    echo "  After power-cycle, the TRX will be at ${post_trx_ip}."
}

method_apply() {
    trap _jtag_octeon_cleanup EXIT

    local trx_ip serial_dev baud
    trx_ip=$(yq_read "$BOARD_CONFIG" network.trx_ip)
    serial_dev=$(yq_read "$BOARD_CONFIG" serial.device)
    baud=$(yq_read "$BOARD_CONFIG" serial.baud)

    echo "=== Phase 1: JTAG bringup ==="
    _phase1_run

    echo ""
    echo "=== Manual pause 1 ==="
    echo "Please:"
    echo "  1. Power OFF the TRX"
    echo "  2. Disconnect the BDI / JTAG cable  (REQUIRED - while connected the BDI holds the"
    echo "     CPU in reset, so the board will NOT finish booting from flash)"
    echo "  3. CLOSE PuTTY/screen on /dev/ttyUSB0 (the script needs the serial port)"
    echo "  4. Power ON the TRX"
    echo ""
    echo "The script will wait for the 'LSM login:' prompt on ${serial_dev} and continue"
    echo "automatically. Press ENTER at any time to skip waiting and proceed immediately."
    echo ""

    local pause1_log="${LOG_DIR}/pause1-serial.log"
    local prompt_seen=0 elapsed=0 max_pause1_wait=300
    local pause1_tail_pid=""
    if uboot_open "$serial_dev" "$baud" "$pause1_log" 2>/dev/null; then
        if [ "$TRX_VERBOSE" = "1" ]; then
            echo "(verbose) Tailing pause1 serial log to console..."
            tail -f "$pause1_log" &
            pause1_tail_pid=$!
        fi
        while [ "$elapsed" -lt "$max_pause1_wait" ]; do
            local key=""
            if IFS= read -rs -t 1 -n 1 key 2>/dev/null; then
                echo "  skipped by user"
                prompt_seen=1
                break
            fi
            if grep -qF "LSM login:" "$pause1_log" 2>/dev/null; then
                echo "  'LSM login:' seen on serial - continuing automatically."
                prompt_seen=1
                break
            fi
            elapsed=$((elapsed + 1))
        done
        if [ -n "$pause1_tail_pid" ]; then
            kill "$pause1_tail_pid" 2>/dev/null || true
            pause1_tail_pid=""
        fi
        uboot_close
    fi

    if [ "$prompt_seen" -ne 1 ]; then
        echo "  WARNING: 'LSM login:' not seen within ${max_pause1_wait}s; continuing anyway."
    fi

    echo ""
    echo "=== Phase 2: SSH + dd image flash ==="
    _phase2_run

    echo ""
    echo "=== Manual pause 2 ==="
    echo "Please power-cycle the TRX now to boot from the newly flashed images."
    echo "rc_post.local (inside the app image) will bring ethernet up automatically."
    echo ""
    read -rp "Press ENTER once the TRX is powered back on (no need to wait for full boot): " _

    _phase2_rehost_for_post_flash
    _phase2_post_flash_config
}

method_verify() {
    local trx_ip ssh_user band_cfg_target rc_post_target
    if yq_exists "$BOARD_CONFIG" network.post_flash.trx_ip; then
        trx_ip=$(yq_read "$BOARD_CONFIG" network.post_flash.trx_ip)
    else
        trx_ip=$(yq_read "$BOARD_CONFIG" network.trx_ip)
    fi
    ssh_user=$(yq_read "$BOARD_CONFIG" phase2.ssh_user)
    band_cfg_target=$(yq_read "$BOARD_CONFIG" band.target_path)
    rc_post_target=$(yq_read "$BOARD_CONFIG" phase2.rc_post_local)

    local sshpass_args=(-p "$TRX_ROOT_PASSWORD")
    local ssh_opts=(-o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o LogLevel=ERROR)

    echo "  Waiting for TRX SSH after final power-cycle..."
    local elapsed=0
    while [ "$elapsed" -lt 300 ]; do
        if sshpass "${sshpass_args[@]}" ssh "${ssh_opts[@]}" "${ssh_user}@${trx_ip}" true 2>/dev/null; then
            break
        fi
        sleep 2
        elapsed=$((elapsed + 2))
    done

    if [ "$elapsed" -ge 300 ]; then
        echo "  [WARN] TRX not reachable over SSH after final power-cycle (it may still be booting)"
        echo "  Final check is manual: confirm the TRX boots and ${rc_post_target} enables ethernet."
        return 0
    fi

    echo "  [ OK ] TRX reachable over SSH after final power-cycle"
    echo "  Verifying installed config files..."
    sshpass "${sshpass_args[@]}" ssh "${ssh_opts[@]}" "${ssh_user}@${trx_ip}" \
        "ls -l '${band_cfg_target}' '${rc_post_target}'" 2>/dev/null || true
    return 0
}

method_monitor() {
    echo "TRX flash complete. Power-cycle and observe operation."
    return 0
}
