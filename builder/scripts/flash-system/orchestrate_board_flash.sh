#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2025-present, Ukama Inc.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

BOARD_NAME=""
CONFIG="${SCRIPT_DIR}/boards.yaml"
YQ_BIN="${SCRIPT_DIR}/.bin/yq"
FLASH_SCRIPT=""
ISO_BUILDER="${SCRIPT_DIR}/create_dual_partition.sh"
RETRIES=3
HOST_ETH=""
SERIAL_PID=""
HTTP_PID=""
SSH_STARTED=0

TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
TMP_LOG_DIR="logs/${TIMESTAMP}_UNKNOWN"
mkdir -p "$TMP_LOG_DIR"

ORCHESTRATOR_LOG="${TMP_LOG_DIR}/orchestrator.log"
SERIAL_LOG="${TMP_LOG_DIR}/serial_console.log"
RAW_SERIAL="${TMP_LOG_DIR}/serial_raw.log"
MAC_FILE="${TMP_LOG_DIR}/mac.txt"
SN_FILE="${TMP_LOG_DIR}/serial.txt"
STATUS_FILE="${TMP_LOG_DIR}/status.txt"

cleanup() {
    if [[ -n "${SERIAL_PID:-}" ]]; then
        kill "$SERIAL_PID" 2>/dev/null || true
    fi

    if [[ -n "${HTTP_PID:-}" ]]; then
        echo "Stopping temporary HTTP server..."
        kill "$HTTP_PID" 2>/dev/null || true
    fi

    if [ "${ORIGINAL_SSH_STATE:-}" = "inactive" ] && [ "${SSH_STARTED:-0}" -eq 1 ]; then
        echo "Restoring SSH state — stopping SSHD" | tee -a "$ORCHESTRATOR_LOG"
        sudo systemctl stop sshd || true
    fi
    if [[ -n "${HOST_ETH:-}" ]]; then
        echo "Restoring NetworkManager control of $HOST_ETH"
        nmcli device set "$HOST_ETH" managed yes 2>/dev/null || true
    fi
    rm -f "$YQ_BIN"
    rm -f alpine.iso
    rm -f "${FLASH_SCRIPT:-}"
}
trap cleanup EXIT

NETWORK_REQUIRED_KEYS=(
    "network.host_eth" "network.host_ip" "network.target_ip"
    "image.name" "image.path"
    "host_device.device" "host_device.iso_url"
    "serial.device" "serial.baud"
    "flash.target_device" "flash.success_marker" "flash.boot_marker"
)

SD_CARD_REQUIRED_KEYS=(
    "image.name" "image.path"
    "sd_card.device"
    "network.host_eth" "network.host_ip"
    "serial.device" "serial.baud"
    "flash.target_device" "flash.success_marker" "flash.boot_marker"
)

RPIBOOT_REQUIRED_KEYS=(
    "image.name" "image.path"
    "serial.device" "serial.baud"
)

ensure_yq() {
    if [ ! -x "$YQ_BIN" ]; then
        echo "Downloading yq..." | tee -a "$ORCHESTRATOR_LOG"
        mkdir -p "$(dirname "$YQ_BIN")"
        curl -L https://github.com/mikefarah/yq/releases/latest/download/yq_linux_amd64 -o "$YQ_BIN"
        chmod +x "$YQ_BIN"
    fi
}

yq_read() {
    "$YQ_BIN" eval ".boards.$BOARD_NAME.$1" "$CONFIG"
}

config_value_exists() {
    local key="$1"
    local value

    value=$("$YQ_BIN" eval ".boards.${BOARD_NAME}.${key} // \"__MISSING__\"" "$CONFIG")
    [[ "$value" != "__MISSING__" && "$value" != "null" ]]
}

validate_config() {
    local flash_method
    local -a required_keys=()

    echo "Validating config..." | tee -a "$ORCHESTRATOR_LOG"

    if ! config_value_exists "method"; then
        echo "Board '$BOARD_NAME' is missing or has no flash method in $CONFIG" | tee -a "$ORCHESTRATOR_LOG"
        exit 1
    fi

    flash_method=$(yq_read 'method')
    case "$flash_method" in
        network)
            required_keys=("${NETWORK_REQUIRED_KEYS[@]}")
            ;;
        sd-card)
            required_keys=("${SD_CARD_REQUIRED_KEYS[@]}")
            ;;
        rpiboot)
            required_keys=("${RPIBOOT_REQUIRED_KEYS[@]}")
            ;;
        *)
            echo "Unsupported flash method '$flash_method' for board '$BOARD_NAME'" | tee -a "$ORCHESTRATOR_LOG"
            exit 1
            ;;
    esac

    for key in "${required_keys[@]}"; do
        if ! config_value_exists "$key"; then
            echo "Missing config for board '$BOARD_NAME': $key" | tee -a "$ORCHESTRATOR_LOG"
            exit 1
        fi
    done
}

retry() {
    local n=1 max=$RETRIES delay=5
    until "$@"; do
        if (( n == max )); then
            echo "Command failed after $n attempts." | tee -a "$ORCHESTRATOR_LOG"
            exit 1
        else
            echo "Retry $n/$max: $*" | tee -a "$ORCHESTRATOR_LOG"
            sleep $delay
            ((n++))
        fi
    done
}

detect_ssh_state() {
    local status=0

    systemctl status sshd.service &>/dev/null || status=$?

    case $status in
        0) echo "active" ;;
        3) echo "inactive" ;;
        4) echo "not-installed" ;;
        *) echo "unknown" ;;
    esac
}

# Main
ORIGINAL_SSH_STATE=$(detect_ssh_state)

# Parse options: -c for config, -b for board name (SMARC, FEM-Control)
while getopts ":c:b:" opt; do
  case "${opt}" in
    c) CONFIG="${OPTARG}" ;;
    b) BOARD_NAME="${OPTARG}" ;;
    *) echo "Usage: $0 -b <board_name> [-c <config_file>]" >&2; exit 1 ;;
  esac
done
shift $((OPTIND-1))

if [[ -z "$BOARD_NAME" ]]; then
    echo "Error: -b <board_name> is required."
    echo "Usage: $0 -b <board_name> [-c <config_file>]"
    echo "Available boards: tnode, anode, cnode"
    exit 1
fi

if [ ! -f "$CONFIG" ]; then
    echo "Config file '$CONFIG' not found."
    exit 1
fi

FLASH_SCRIPT="flash-${BOARD_NAME}.sh"

ensure_yq
validate_config

FLASH_METHOD=$(yq_read 'method')

if [ "$FLASH_METHOD" = "rpiboot" ]; then
    IMG_NAME=$(yq_read 'image.name')
    IMG_PATH=$(yq_read 'image.path')

    if [ ! -f "$IMG_PATH" ]; then
        echo "Image not found: $IMG_PATH"
        exit 1
    fi

    cp "$IMG_PATH" "$IMG_NAME"

    if [ ! -f "${SCRIPT_DIR}/flash-cnode.sh" ]; then
        echo "flash-cnode.sh not found"
        exit 1
    fi

    bash "${SCRIPT_DIR}/flash-cnode.sh"
    exit 0
fi

if [ "$FLASH_METHOD" = "sd-card" ]; then
    IMG_PATH=$(yq_read 'image.path')
    SD_DEV=$(yq_read 'sd_card.device')
    HOST_IP=$(yq_read 'network.host_ip')
    SERIAL_DEV=$(yq_read 'serial.device')
    SUCCESS_MARKER=$(yq_read 'flash.success_marker')
    BOOT_MARKER=$(yq_read 'flash.boot_marker')

    if [ ! -f "$IMG_PATH" ]; then
        echo "Image not found: $IMG_PATH"
        exit 1
    fi

    if [ ! -b "$SD_DEV" ]; then
        echo "SD card device '$SD_DEV' not found."
        echo "Make sure SD card is inserted and use the correct device path."
        exit 1
    fi

    echo "Creating SD card for $BOARD_NAME..."
    DEV="$SD_DEV" \
    IMAGE_PATH="$IMG_PATH" \
    BOARD_NAME="$BOARD_NAME" \
    HOST_IP="$HOST_IP" \
        bash "${SCRIPT_DIR}/create_sdcard.sh"

    echo ""
    echo "========================================"
    echo "SD card created successfully!"
    echo "========================================"
    echo ""
    echo "Next steps:"
    echo "1. Remove SD card from Ubuntu"
    echo "2. Insert SD card into $BOARD_NAME board"
    echo "3. Connect serial cable: $SERIAL_DEV"
    echo "4. Power on the board"
    echo "5. Board will auto-flash to eMMC"
    echo ""
    echo "Press ENTER when ready to monitor serial..."
    read -r

    touch "$SERIAL_LOG"
    cat "$SERIAL_DEV" | tee "$RAW_SERIAL" "$SERIAL_LOG" &
    SERIAL_PID=$!

    echo "Monitoring serial output from $BOARD_NAME via ${SERIAL_DEV}..."
    echo "Waiting for '${SUCCESS_MARKER}'"
    retry timeout 300 grep -q "$SUCCESS_MARKER" "$SERIAL_LOG"
    echo "Flash completed."

    echo "Waiting for '${BOOT_MARKER}'"
    retry timeout 120 grep -q "$BOOT_MARKER" "$SERIAL_LOG"
    echo "System booted."

    echo "PASS" > "$STATUS_FILE"
    exit 0
fi

HOST_ETH=$(yq_read        'network.host_eth')
HOST_IP=$(yq_read         'network.host_ip')
TARGET_IP=$(yq_read       'network.target_ip')
IMG_NAME=$(yq_read        'image.name')
IMG_PATH=$(yq_read        'image.path')
HOST_DEV=$(yq_read        'host_device.device')
ISO_URL=$(yq_read         'host_device.iso_url')
SERIAL_DEV=$(yq_read      'serial.device')
SERIAL_BAUD=$(yq_read     'serial.baud')
TARGET_DEV=$(yq_read      'flash.target_device')
SUCCESS_MARKER=$(yq_read  'flash.success_marker')
BOOT_MARKER=$(yq_read     'flash.boot_marker')

{
    # Verify device exists
    if [ ! -b "$HOST_DEV" ]; then
        echo "Host device '$HOST_DEV' not found or is not a block device."
        echo "Make sure it's plugged in and use the full device path (e.g., /dev/sdb)"
        exit 1
    fi

    # Configure host Ethernet
    echo "Configure dev Ethernet (${HOST_ETH})" | tee -a "$ORCHESTRATOR_LOG"
    nmcli device set  "$HOST_ETH" managed no
    sudo ip link set  "$HOST_ETH" down  || true
    sudo ip addr flush dev "$HOST_ETH"  || true
    sudo ip addr add   "$HOST_IP/24" dev "$HOST_ETH"
    sudo ip link set   "$HOST_ETH"   up
    ip addr show       "$HOST_ETH"

    # install and configure dnsmasq
    echo "Installing & configuring dnsmasq" | tee -a "$ORCHESTRATOR_LOG"
    sudo apt-get update -qq
    sudo apt-get install -y dnsmasq

    sudo tee /etc/dnsmasq.d/smarc.conf >/dev/null <<EOF
interface=${HOST_ETH}
bind-interfaces
dhcp-range=${HOST_IP%.*}.100,${HOST_IP%.*}.200,12h
dhcp-option=option:router,${HOST_IP}
dhcp-option=option:dns-server,8.8.8.8,8.8.4.4
EOF

    sudo systemctl restart dnsmasq
    echo "dnsmasq DHCP server is up on $HOST_ETH" | tee -a "$ORCHESTRATOR_LOG"

    # start sshd
    echo "Start SSH (as needed)"
    if [ "$ORIGINAL_SSH_STATE" = "inactive" ]; then
        echo "Starting SSH temporarily for image transfer"
        sudo systemctl start sshd
        SSH_STARTED=1
    elif [ "$ORIGINAL_SSH_STATE" = "not-installed" ]; then
        echo "SSHD is not installed — skipping SSH-related steps."
    fi

    echo "Download Alpine ISO"
    curl -L "$ISO_URL" -o alpine.iso

    echo "Generate flash script"
    cat > "$FLASH_SCRIPT" <<EOF
#!/bin/sh
set -eux

echo "[$BOARD_NAME] Enabling eth0"
ip link set eth0 up

echo "[$BOARD_NAME] Bringing up eth0 via DHCP (udhcpc)"
udhcpc -i eth0 -q

echo "[$BOARD_NAME] Waiting a couple seconds for lease…"
sleep 2

echo "[$BOARD_NAME] Detecting eMMC device..."
for dev in /dev/mmcblk*; do
    if [ -e "\${dev}boot0" ] && [ -e "\${dev}boot1" ]; then
        EMMC_DEV="\$dev"
        break
    fi
done

if [ -z "\${EMMC_DEV:-}" ]; then
    echo "[ERROR] No eMMC device found with boot0/boot1"
    exit 1
fi
echo "[$BOARD_NAME] Detected eMMC device: \$EMMC_DEV"

echo "[$BOARD_NAME] Downloading image from ${HOST_IP}:8000/${IMG_NAME}"
wget "http://${HOST_IP}:8000/${IMG_NAME}" -O "/mnt/${IMG_NAME}"

ls -lh "/mnt/${IMG_NAME}"

if [ ! -f "/mnt/${IMG_NAME}" ]; then
    echo "[$BOARD_NAME] Image not found after wget!"
    exit 1
fi

echo "[$BOARD_NAME] Zeroing first 64MB of \$EMMC_DEV"
dd if=/dev/zero of="\$EMMC_DEV" bs=1M count=64

echo "[$BOARD_NAME] Flashing image to \$EMMC_DEV"
echo "[WARN] pv not found — flashing without progress bar"
dd if="/mnt/${IMG_NAME}" of="\$EMMC_DEV" bs=4M
sync

echo "[$BOARD_NAME] Flash complete"
reboot
EOF
    chmod +x "$FLASH_SCRIPT"

    echo "Create bootable $HOST_DEV"
    DEV="$HOST_DEV" \
    FLASH_SCRIPT="$FLASH_SCRIPT" \
    BOARD_NAME="$BOARD_NAME" \
    ISO_FILE="alpine.iso" \
       "$ISO_BUILDER"

    # Start HTTP server to serve image
    echo "Starting temporary HTTP server to serve image"
    IMG_DIR=$(dirname "$IMG_PATH")
    cd "$IMG_DIR"
    python3 -m http.server 8000 > /dev/null 2>&1 &
    HTTP_PID=$!
    echo "HTTP server started with PID $HTTP_PID"
    cd - > /dev/null

    echo "Do you want to test this image in QEMU before inserting into target? (y/N): "
    read -r qemu_choice

    if [[ "$qemu_choice" == "y" || "$qemu_choice" == "Y" ]]; then
        if [ ! -e /dev/kvm ]; then
            echo "/dev/kvm not found. Running without KVM acceleration."
            KVM=""
        else
            KVM="-enable-kvm"
        fi

        echo "Booting actual device ($HOST_DEV) in QEMU..."
        sudo qemu-system-x86_64 \
             $KVM \
             -m 1024 \
             -machine type=pc,accel=kvm \
             -boot order=d \
             -drive file="$HOST_DEV",format=raw,if=virtio,media=disk \
             -serial mon:stdio \
             -display none \
             -name "AlpineUSBTest"
    else
        # Eject USB
        sudo eject "$HOST_DEV"

        echo "Insert USB into $BOARD_NAME and power it up"
        echo "USB is ready. Insert into target and run flash-smarc.sh manually."
        echo "(Boot to Alpine, mount /dev/sda2, run /mnt/flash-smarc.sh)"
        echo "Please ensure the board is powered on and connected via serial (${SERIAL_DEV})."
        echo "Once ready, press ENTER to begin monitoring the serial port..."
        read -r

        echo "Monitoring serial output from $BOARD_NAME via ${SERIAL_DEV}..."
        touch "$SERIAL_LOG"
        cat "$SERIAL_DEV" | tee "$RAW_SERIAL" "$SERIAL_LOG" &
        SERIAL_PID=$!

        echo "Waiting for '${SUCCESS_MARKER}'"
        retry timeout 300 grep -q "$SUCCESS_MARKER" "$SERIAL_LOG"
        echo "Flash completed."

        echo "Waiting for '${BOOT_MARKER}'"
        retry timeout 120 grep -q "$BOOT_MARKER" "$SERIAL_LOG"
        echo "System booted."

        echo "PASS" > "$STATUS_FILE"
    fi
} 2>&1 | tee -a "$ORCHESTRATOR_LOG" || {
    echo "Flashing failed — logs in: $TMP_LOG_DIR" | tee -a "$ORCHESTRATOR_LOG"
    echo "FAIL" > "$STATUS_FILE"
    exit 1
}
