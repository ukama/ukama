#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2025-present, Ukama Inc.

source "${LIB_DIR}/extract.sh"
source "${LIB_DIR}/serial.sh"

HTTP_PID=""
SERIAL_PID=""
ORIGINAL_SSH_STATE=""
SSH_STARTED=0
HOST_ETH=""
DNSMASQ_CONF="/etc/dnsmasq.d/smarc.conf"

_detect_ssh_state() {
    local status=0
    systemctl status sshd.service &>/dev/null || status=$?
    case $status in
        0) echo "active" ;;
        3) echo "inactive" ;;
        4) echo "not-installed" ;;
        *) echo "unknown" ;;
    esac
}

_network_cleanup() {
    [ -n "$HTTP_PID" ] && kill "$HTTP_PID" 2>/dev/null || true
    [ -n "$SERIAL_PID" ] && kill "$SERIAL_PID" 2>/dev/null || true

    if [ "$ORIGINAL_SSH_STATE" = "inactive" ] && [ "$SSH_STARTED" -eq 1 ]; then
        sudo systemctl stop sshd 2>/dev/null || true
    fi

    if [ -n "$HOST_ETH" ]; then
        nmcli device set "$HOST_ETH" managed yes 2>/dev/null || true
    fi

    sudo rm -f "$DNSMASQ_CONF" 2>/dev/null || true
    sudo systemctl reload dnsmasq 2>/dev/null || true

    rm -f alpine.iso "flash-${BOARD}.sh" 2>/dev/null || true
}

method_validate() {
    local image_path host_dev
    image_path=$(yq_read "$BOARD_CONFIG" image.path)
    host_dev=$(yq_read "$BOARD_CONFIG" target_usb)

    if [ ! -f "$image_path" ]; then
        echo "  [FAIL] image not found at $image_path"
        return 1
    fi

    if declare -F detect_image_format >/dev/null; then
        local fmt
        fmt=$(detect_image_format "$image_path" 2>/dev/null || true)
        [ -z "$fmt" ] && fmt="unknown"
        echo "  [ OK ] image: $image_path ($fmt)"
    fi

    if [ ! -b "$host_dev" ]; then
        echo "  [FAIL] target USB '$host_dev' is not a block device"
        return 1
    fi
    echo "  [ OK ] target USB: $host_dev"

    return 0
}

method_confirm() {
    local host_dev image_path
    host_dev=$(yq_read "$BOARD_CONFIG" target_usb)
    image_path=$(yq_read "$BOARD_CONFIG" image.path)

    echo ""
    echo "Plan:"
    echo "  - configure host ethernet and DHCP for the target"
    echo "  - build a bootable Alpine USB on $host_dev"
    echo "  - serve $image_path over HTTP for the target to download"
    echo "  - target boots Alpine, flashes its eMMC, reboots"
    echo ""
    echo "This will ERASE $host_dev and modify host networking."
    read -rp "Type 'yes' to continue: " confirm
    [ "$confirm" = "yes" ]
}

method_apply() {
    trap _network_cleanup EXIT

    local host_eth host_ip target_ip image_path img_name host_dev iso_url http_port
    host_eth=$(yq_read "$BOARD_CONFIG" network.host_eth)
    host_ip=$(yq_read "$BOARD_CONFIG" network.host_ip)
    target_ip=$(yq_read "$BOARD_CONFIG" network.target_ip)
    image_path=$(yq_read "$BOARD_CONFIG" image.path)
    img_name=$(yq_read "$BOARD_CONFIG" image.name)
    host_dev=$(yq_read "$BOARD_CONFIG" target_usb)
    iso_url=$(yq_read "$BOARD_CONFIG" alpine.iso_url)
    http_port=$(yq_read "$BOARD_CONFIG" http.port)
    [ "$http_port" = "null" ] && http_port=8000

    HOST_ETH="$host_eth"
    ORIGINAL_SSH_STATE=$(_detect_ssh_state)

    echo "Configuring host ethernet ($host_eth)..."
    nmcli device set "$host_eth" managed no
    sudo ip link set "$host_eth" down || true
    sudo ip addr flush dev "$host_eth" || true
    sudo ip addr add "$host_ip/24" dev "$host_eth"
    sudo ip link set "$host_eth" up

    echo "Installing and configuring dnsmasq..."
    sudo apt-get update -qq
    sudo apt-get install -y dnsmasq

    sudo tee "$DNSMASQ_CONF" >/dev/null <<EOF
interface=${host_eth}
bind-interfaces
dhcp-range=${host_ip%.*}.100,${host_ip%.*}.200,12h
dhcp-option=option:router,${host_ip}
dhcp-option=option:dns-server,8.8.8.8,8.8.4.4
EOF
    sudo systemctl restart dnsmasq

    if [ "$ORIGINAL_SSH_STATE" = "inactive" ]; then
        sudo systemctl start sshd
        SSH_STARTED=1
    fi

    echo "Downloading Alpine ISO..."
    curl -L "$iso_url" -o alpine.iso

    local flash_script="flash-${BOARD}.sh"
    render_template "${BOARD_DIR}/payloads/flash-network.sh.tpl" "$flash_script" \
        BOARD_NAME "$BOARD" \
        HOST_IP "$host_ip" \
        HTTP_PORT "$http_port" \
        IMG_NAME "$img_name"
    chmod +x "$flash_script"

    echo "Building bootable Alpine USB on $host_dev..."
    DEV="$host_dev" \
    FLASH_SCRIPT="$flash_script" \
    BOARD_NAME="$BOARD" \
    ISO_FILE="alpine.iso" \
        "${LIB_DIR}/alpine_iso.sh"

    echo "Starting HTTP server on port ${http_port}..."
    local img_dir prev_dir
    img_dir=$(dirname "$image_path")
    prev_dir=$(pwd)
    cd "$img_dir"
    python3 -m http.server "$http_port" >/dev/null 2>&1 &
    HTTP_PID=$!
    cd "$prev_dir"

    sudo eject "$host_dev" 2>/dev/null || true

    echo ""
    echo "USB ready. Insert it into the $BOARD board and power on."
    echo "The board will boot Alpine, download the image, and flash its eMMC."
}

method_verify() {
    echo "Network-method verification happens on the target via serial markers (see monitor)."
    return 0
}

method_monitor() {
    local serial_dev success_marker boot_marker
    serial_dev=$(yq_read "$BOARD_CONFIG" serial.device)
    success_marker=$(yq_read "$BOARD_CONFIG" serial.success_marker)
    boot_marker=$(yq_read "$BOARD_CONFIG" serial.boot_marker)

    if [ ! -e "$serial_dev" ]; then
        echo "Serial $serial_dev not available — skipping monitor."
        return 0
    fi

    echo ""
    read -rp "Press ENTER once the target is booting, or 's' to skip serial monitor: " resp
    [ "$resp" = "s" ] && return 0

    echo "Watching $serial_dev for '$success_marker'..."
    serial_wait_for_marker "$serial_dev" "${LOG_DIR}/serial.log" "$success_marker" 300 || {
        echo "Did not see success marker within timeout."
        return 1
    }
    echo "Flash completed on target."

    if [ -n "$boot_marker" ] && [ "$boot_marker" != "null" ]; then
        echo "Watching for boot marker '$boot_marker'..."
        serial_wait_for_marker "$serial_dev" "${LOG_DIR}/serial.log" "$boot_marker" 120 || {
            echo "Did not see boot marker within timeout."
            return 1
        }
        echo "Target booted."
    fi
}
