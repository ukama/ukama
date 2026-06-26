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

# Temporary verbose mode for debugging Phase 1. Set TRX_VERBOSE=0 to disable.
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
        # Wait for the BDI prompt to come back before quitting. Without this the
        # quit can merge onto the command line (\"go 0x400000quit\" -> syntax error)
        # and the command never executes.
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

    # Compute actual file sizes for bootcbyflash so cp.b copies exactly the right amount.
    local os_path rd_path os_size rd_size os_size_hex rd_size_hex
    os_path=$(yq_read "$BOARD_CONFIG" phase1.artifacts.os.path)
    rd_path=$(yq_read "$BOARD_CONFIG" phase1.artifacts.rd.path)
    os_size=$(stat -c %s "$os_path" 2>/dev/null || yq_read "$BOARD_CONFIG" phase1.artifacts.os.size)
    rd_size=$(stat -c %s "$rd_path" 2>/dev/null || yq_read "$BOARD_CONFIG" phase1.artifacts.rd.size)
    os_size_hex=$(printf '0x%x' "$os_size")
    rd_size_hex=$(printf '0x%x' "$rd_size")

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
    # Add mw64 to bootcmd as a fallback in case preboot doesn't run (e.g. bootdelay=0).
    uboot_send_and_wait "$dev" 'setenv bootcmd "mw64 0x00011800B0001000 0x0140; run i2cinit; run namedalloc; run bootcby${bootby}"' "$prompt" 8
    # bootcbytftp uses the actual staged filenames (lsm_os_trx.gz / lsm_rd_trx.gz).
    uboot_send_and_wait "$dev" 'setenv bootcbytftp "tftp 0x21000000 lsm_os_trx.gz; gunzip 0x21000000 0x20000000 0x1000000; tftp 0x30800000 lsm_rd_trx.gz; bootoctlinux 0x20000000 coremask=0x7 endbootargs rd_name=initrd mem=512M;"' "$prompt" 8
    # bootcbyflash: copy OS + RD from flash (where Phase 1 wrote them) and boot Linux.
    local bootcbyflash_cmd
    bootcbyflash_cmd="setenv bootcbyflash \"cp.b 0x17E20000 0x21000000 ${os_size_hex}; gunzip 0x21000000 0x20000000 0x1000000; cp.b 0x18320000 0x30800000 ${rd_size_hex}; bootoctlinux 0x20000000 coremask=0x7 endbootargs rd_name=initrd mem=512M;\""
    uboot_send_and_wait "$dev" "$bootcbyflash_cmd" "$prompt" 8
    uboot_send_and_wait "$dev" 'setenv namedalloc "namedalloc dsp-dump 0x400000 0x7f4D0000; namedalloc cazac 0x630000 0x7f8D0000; namedalloc cpu-dsp-if 0x100000 0x7ff00000; namedalloc dsp-log-buf 0x4000000 0x80000000; namedalloc initrd 0x2800000 0x30800000;"' "$prompt" 8
    uboot_send_and_wait "$dev" "setenv mk_ubootenv 1" "$prompt" 8
    # SGMII autoneg must be enabled BEFORE Linux boots, otherwise octeth0 stays down.
    # u-boot's 'preboot' runs automatically before bootcmd.
    uboot_send_and_wait "$dev" "setenv preboot 'mw64 0x00011800B0001000 0x0140'" "$prompt" 8

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

    # --- BDI state check ---
    # The bench procedure never reloads CONFIG on an already-configured BDI: at
    # cnMIPS#0> the operator goes straight to go + oct-remote-boot. Every bring-up
    # attempted right after a CONFIG reload (= BDI reboot) has hung in DDR init,
    # so CONFIG is sent only when the BDI is actually bare (Core#0>), and after
    # that the TRX must be cold power-cycled before bring-up.
    echo "Checking BDI state at ${bdi_ip} (cnMIPS#0> = configured, Core#0> = bare)..."
    local bdi_state
    bdi_state=$(expect -c "
        set timeout 20
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
        echo "  BDI already configured — skipping CONFIG reload (matches the manual flow)."
    elif [ "$bdi_state" = "BDI_BARE" ]; then
        echo "  BDI is unconfigured — sending HOST + CONFIG (the BDI will reboot and auto-load)..."
        if ! expect -c "
            set timeout 25
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
        echo "  Config sent — BDI is rebooting and auto-loading cnf71xx.cfg over TFTP."
        echo "  Keeping the TFTP server up and polling until the BDI comes up configured..."
        # The BDI's post-reboot auto-load needs our TFTP server (still running here) to be
        # serving cnf71xx.cfg. Poll the prompt until it returns cnMIPS#0>, then continue in
        # this same run — do NOT exit (exiting tears down TFTP and the BDI stays bare).
        local cfg_wait=0 bdi_now=""
        while [ "$cfg_wait" -lt 120 ]; do
            sleep 10
            cfg_wait=$((cfg_wait + 10))
            bdi_now=$(expect -c "
                set timeout 15
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
        if [ "$gdb_wait" -ge 60 ]; then
            echo "ERROR: BDI GDB port still closed after 60s — the BDI may not have auto-loaded"
            echo "  its config (TFTP must be running when the BDI boots)."
            echo "  Power-cycle the BDI now (TFTP is still being served) and re-run this script."
            return 1
        fi
        sleep 5
    done
    echo "  GDB port is open."

    # --- Phase 1 core bring-up ---
    # The bench procedure (confirmed with Supreeth):
    #   1. go 0x400000   ONCE, via BDI telnet.
    #   2. oct-remote-boot. If it segfaults / GDB-errors ("Core 0, in reset, told to
    #      continue / Segmentation fault"), just RUN OCT-REMOTE-BOOT AGAIN — no
    #      power-cycle, no re-sending go. It succeeds on the 1st run, rarely the 2nd.
    #   3. watch serial, interrupt autoboot, flash via TFTP.
    # We automate exactly that: send go once, keep the serial open, and retry only
    # oct-remote-boot until the u-boot prompt appears.
    local oct_log="${LOG_DIR}/oct-remote-boot.log"
    local prompt_seen=0

    echo "Sending 'go 0x400000' via BDI telnet (once)..."
    if ! bdi_telnet_cmd "$bdi_ip" "go 0x400000"; then
        echo "ERROR: BDI 'go 0x400000' failed (is the BDI at cnMIPS#0>?)."
        return 1
    fi
    sleep 5

    # Open the serial console and start interrupting autoboot. Kept open across all
    # oct-remote-boot retries so the u-boot prompt is caught whenever it appears.
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
        # The TRX u-boot has a zero-second autoboot delay; keep a key pressed
        # (space) so it stops immediately when it checks for input.
        exec 3>"$serial_dev"
        while true; do
            printf ' ' >&3
            sleep 0.03
        done
    ) &
    SPAM_PID=$!

    # The Octeon DDR PLL is a cold-boot lottery: each oct-remote-boot run re-rolls the
    # measured DDR clock (seen 267 / 400 / 201 MHz across re-runs). Only ~400 MHz works —
    # at the wrong clock the SGMII ref clock is off and ethernet never links. So we
    # accept an attempt ONLY if u-boot comes up AND the clock is ~400; otherwise we
    # re-run oct-remote-boot (no re-go, no power-cycle) to re-roll the DDR clock.
    local oct_try=0 max_oct_tries=8
    local elapsed clk mhz oct_exit ddr_mislock prompt_this
    while [ "$oct_try" -lt "$max_oct_tries" ]; do
        oct_try=$((oct_try + 1))
        echo ""
        echo "=== oct-remote-boot attempt ${oct_try}/${max_oct_tries} (no re-go, no power-cycle) ==="

        # Clear any stale oct-remote-boot still holding the BDI GDB port.
        sudo pkill -9 -f oct-remote-boot 2>/dev/null || true
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

        # Match the u-boot prompt broadly — it may be "Octeon zen(ram)=>" or
        # "Octeon zen(Failsafe)=>" depending on board/flash state (per Supreeth).
        echo "Waiting for u-boot prompt 'Octeon zen…=>' with DDR ~400 MHz (up to 180s)..."
        elapsed=0; clk=""; mhz=""; ddr_mislock=0; prompt_this=0
        while [ "$elapsed" -lt 180 ]; do
            # Once a DDR clock is measured, reject a bad lock immediately (no point
            # waiting for u-boot on a clock that won't link ethernet).
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
            # If oct-remote-boot exited, u-boot may still be printing — short grace window.
            if ! kill -0 "$REMOTE_BOOT_PID" 2>/dev/null; then
                local grace=0
                while [ "$grace" -lt 15 ]; do
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
            # Leave oct-remote-boot running — it hosts u-boot in DDR over GDB while we
            # flash. REMOTE_BOOT_PID stays set and is reaped by the cleanup after flashing.
            prompt_seen=1
            echo "  u-boot prompt + good DDR clock (${mhz} MHz) on attempt ${oct_try}."
            if [ -n "$OCT_TAIL_PID" ]; then
                kill "$OCT_TAIL_PID" 2>/dev/null || true
                OCT_TAIL_PID=""
            fi
            break
        fi

        # Failed attempt: stop oct-remote-boot if still alive, report why, then re-run.
        if kill -0 "$REMOTE_BOOT_PID" 2>/dev/null; then
            sudo kill "$REMOTE_BOOT_PID" 2>/dev/null || true
        fi
        wait "$REMOTE_BOOT_PID" 2>/dev/null
        oct_exit=$?
        REMOTE_BOOT_PID=""

        if [ "$ddr_mislock" -eq 1 ]; then
            echo "  DDR mislocked at ${mhz} MHz (need ~400) — re-running to re-roll the DDR clock."
        elif [ "$prompt_this" -eq 1 ]; then
            echo "  u-boot came up but DDR clock unknown/!~400 — re-running."
        elif [ "$oct_exit" -eq 139 ] || grep -qaE "Segmentation fault|GDB Reply Error|in reset, told to continue" "$oct_log" 2>/dev/null; then
            echo "  oct-remote-boot exited (GDB error / segfault) — re-running it (normal recovery)."
        else
            echo "  no u-boot prompt within 180s — re-running."
        fi
        echo "--- last 15 lines of oct-remote-boot output ---"
        tail -n 15 "$oct_log" 2>/dev/null | sed 's/^/    /' || true
        sleep 2
    done

    # Stop the autoboot spammer regardless of outcome.
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
    uboot_send_and_wait "$serial_dev" "mw64 0x00011800B0001000 0x0140" "$uboot_prompt" 8 || true
    uboot_send_and_wait "$serial_dev" "mw64 0x00011800B0001000 0x0140" "$uboot_prompt" 8 || true
    uboot_send_and_wait "$serial_dev" "saveenv" "$uboot_prompt" 15 || true

    local ping_marker link_up=0 ping_attempt=0 max_ping_attempts=6
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

    if [ -n "$serial_tail_pid" ]; then
        kill "$serial_tail_pid" 2>/dev/null || true
        serial_tail_pid=""
    fi
    uboot_close
    tftp_stop
    [ -n "$REMOTE_BOOT_PID" ] && sudo kill "$REMOTE_BOOT_PID" 2>/dev/null || true
    REMOTE_BOOT_PID=""
}

# On a freshly-erased board the post_config script (rc_post.local) that enables SGMII
# autoneg isn't on the board yet (it arrives with the Phase 2 app images), so Linux
# boots with ethernet down and SSH can't connect. Bootstrap it the way it's done by
# hand: log in over serial and run "devmem 0x00011800B0001000 64 0x0140" twice.
# Best-effort — needs exclusive /dev/ttyUSB0; if it can't, we fall back to SSH wait.
_phase2_enable_ethernet_over_serial() {
    local serial_dev baud
    serial_dev=$(yq_read "$BOARD_CONFIG" serial.device)
    baud=$(yq_read "$BOARD_CONFIG" serial.baud)

    local holder=""
    if command -v lsof >/dev/null 2>&1; then
        holder=$(lsof -t "$serial_dev" 2>/dev/null | head -1 || true)
    fi
    if [ -n "$holder" ]; then
        echo "  NOTE: $serial_dev is held by another process (PID $holder) — skipping auto"
        echo "        ethernet-enable. Close PuTTY/screen, or on the console run twice:"
        echo "          devmem 0x00011800B0001000 64 0x0140"
        return 0
    fi

    echo "Enabling TRX ethernet over serial (login as root, devmem x2)..."
    uboot_open "$serial_dev" "$baud" "${LOG_DIR}/phase2-serial.log" || return 0

    uboot_send "$serial_dev" ""
    sleep 2
    if ! uboot_wait_for "~#" 5; then
        uboot_send "$serial_dev" ""
        if uboot_wait_for "login:" 90; then
            uboot_send "$serial_dev" "root"
            uboot_wait_for "assword" 15 || true
            uboot_send "$serial_dev" "cavium.lte"
            if ! uboot_wait_for "~#" 30; then
                echo "  WARNING: serial login didn't reach a shell prompt; relying on SSH wait."
                uboot_close
                return 0
            fi
        else
            echo "  WARNING: no serial login prompt (still booting?); relying on SSH wait."
            uboot_close
            return 0
        fi
    fi

    uboot_send "$serial_dev" "devmem 0x00011800B0001000 64 0x0140"
    sleep 1
    uboot_send "$serial_dev" "devmem 0x00011800B0001000 64 0x0140"
    sleep 1
    echo "  ethernet-enable (devmem x2) sent over serial."
    uboot_close
    return 0
}

_phase2_run() {
    local trx_ip ssh_user staging
    trx_ip=$(yq_read "$BOARD_CONFIG" network.trx_ip)
    ssh_user=$(yq_read "$BOARD_CONFIG" phase2.ssh_user)
    staging=$(yq_read "$BOARD_CONFIG" phase2.ssh_staging_dir)

    local sshpass_args=(-p "$TRX_ROOT_PASSWORD")
    local ssh_opts=(-o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o LogLevel=ERROR)

    _phase2_enable_ethernet_over_serial

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
        echo "ERROR: TRX not reachable via SSH within 120s."
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
        # Supreeth confirmed the bench uses bs=1M for the 8 .img dd's on TRX.
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

    # Copy band config and rc_post.local while we still have SSH/ethernet from boot 1.
    local band_default band_configs_dir band_cfg_src band_cfg_target rc_post_src rc_post_target
    band_default=$(yq_read "$BOARD_CONFIG" band.default)
    band_configs_dir=$(yq_read "$BOARD_CONFIG" band.configs_dir)
    band_cfg_src="${band_configs_dir}/${band_default}.cfg"
    band_cfg_target=$(yq_read "$BOARD_CONFIG" band.target_path)

    if [ -f "$band_cfg_src" ]; then
        echo "Installing band config (${band_cfg_src} -> ${band_cfg_target})..."
        sshpass "${sshpass_args[@]}" ssh "${ssh_opts[@]}" "${ssh_user}@${trx_ip}" "mkdir -p $(dirname "$band_cfg_target")"
        sshpass "${sshpass_args[@]}" scp "${ssh_opts[@]}" "$band_cfg_src" "${ssh_user}@${trx_ip}:${band_cfg_target}"
    else
        echo "WARNING: band config source not found at ${band_cfg_src}"
    fi

    rc_post_src="$(dirname "$BOARD_CONFIG")/payloads/rc_post.local"
    rc_post_target=$(yq_read "$BOARD_CONFIG" phase2.rc_post_local)
    echo "Installing rc_post.local (${rc_post_src} -> ${rc_post_target})..."
    sshpass "${sshpass_args[@]}" ssh "${ssh_opts[@]}" "${ssh_user}@${trx_ip}" "mkdir -p /mnt/app"
    sshpass "${sshpass_args[@]}" scp "${ssh_opts[@]}" "$rc_post_src" "${ssh_user}@${trx_ip}:${rc_post_target}"
    sshpass "${sshpass_args[@]}" ssh "${ssh_opts[@]}" "${ssh_user}@${trx_ip}" "chmod +x ${rc_post_target}"

    echo ""
    echo "All 8 images written and config files installed."
    echo "  Band config : ${band_cfg_target}"
    echo "  Ethernet fix: ${rc_post_target}"
    echo ""
    echo "Next: power-cycle the TRX to boot from the newly flashed images."
    echo "After power-cycle, rc_post.local will enable ethernet automatically."
}

# After the full app images are flashed, the TRX reboots into its production
# network (e.g. 10.102.81.61). Move the bench box NIC to the matching subnet so
# method_verify can reach it.
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

    # Auto-detect interface if not specified.
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
    echo "  2. Disconnect the BDI / JTAG cable  (REQUIRED — while connected the BDI holds the"
    echo "     CPU in reset, so the board will NOT finish booting from flash)"
    echo "  3. CLOSE PuTTY/screen on /dev/ttyUSB0 (the script needs the serial port)"
    echo "  4. Power ON the TRX"
    echo ""
    echo "The script will wait for the 'LSM login:' prompt on ${serial_dev} and continue"
    echo "automatically. Press ENTER at any time to skip waiting and proceed immediately."
    echo ""

    local pause1_log="${LOG_DIR}/pause1-serial.log"
    local prompt_seen=0 elapsed=0 max_pause1_wait=240
    if uboot_open "$serial_dev" "$baud" "$pause1_log" 2>/dev/null; then
        while [ "$elapsed" -lt "$max_pause1_wait" ]; do
            # Allow the operator to skip the wait by pressing ENTER.
            local key=""
            if IFS= read -rs -t 1 -n 1 key 2>/dev/null; then
                echo "  skipped by user"
                prompt_seen=1
                break
            fi
            if grep -qF "LSM login:" "$pause1_log" 2>/dev/null; then
                echo "  'LSM login:' seen on serial — continuing automatically."
                prompt_seen=1
                break
            fi
            elapsed=$((elapsed + 1))
        done
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
    echo "rc_post.local is already installed, so ethernet will come up automatically."
    echo ""
    read -rp "Press ENTER once the TRX has booted: " _

    # The flashed app images reconfigure the TRX to its production IP. Move the
    # bench box to the matching subnet so verification can reach it.
    _phase2_rehost_for_post_flash
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
    while [ "$elapsed" -lt 120 ]; do
        if sshpass "${sshpass_args[@]}" ssh "${ssh_opts[@]}" "${ssh_user}@${trx_ip}" true 2>/dev/null; then
            break
        fi
        sleep 2
        elapsed=$((elapsed + 2))
    done

    if [ "$elapsed" -ge 120 ]; then
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
