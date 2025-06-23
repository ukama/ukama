#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
# Copyright (c) 2025-present, Ukama Inc.

set -euo pipefail

CONFIG="smarc_config.yaml"
YQ_BIN="./.bin/yq"
FLASH_SCRIPT="flash-smarc.sh"
ISO_BUILDER="./create_dual_partition_usb.sh"
RETRIES=3

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

    if [ "$ORIGINAL_SSH_STATE" != "active" ]; then
        echo "Restoring SSH state — stopping SSHD" | tee -a "$ORCHESTRATOR_LOG"
        sudo systemctl stop sshd || true
    fi
    echo "Restoring NetworkManager control of $DEV_ETH"
    nmcli device set "$DEV_ETH" managed yes 2>/dev/null || true
    rm -f "$YQ_BIN"
    rm -f alpine.iso
    rm -f ${FLASH_SCRIPT}
}
trap cleanup EXIT

REQUIRED_KEYS=(
    ".network.dev_eth" ".network.host_ip" ".network.target_ip"
    ".image.name" ".image.path"
    ".usb.device" ".usb.iso_url"
    ".serial.device" ".serial.baud"
    ".flash.target_device" ".flash.success_marker" ".flash.boot_marker"
)

ensure_yq() {
    if [ ! -x "$YQ_BIN" ]; then
        echo "Downloading yq..." | tee -a "$ORCHESTRATOR_LOG"
        mkdir -p .bin
        curl -L https://github.com/mikefarah/yq/releases/latest/download/yq_linux_amd64 -o "$YQ_BIN"
        chmod +x "$YQ_BIN"
    fi
}

yq_read() {
    "$YQ_BIN" eval "$1" "$CONFIG"
}

validate_config() {
    echo "Validating config..." | tee -a "$ORCHESTRATOR_LOG"
    for key in "${REQUIRED_KEYS[@]}"; do
        if ! "$YQ_BIN" eval "$key" "$CONFIG" &>/dev/null; then
            echo "Missing config: $key" | tee -a "$ORCHESTRATOR_LOG"
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
    systemctl status sshd.service &>/dev/null
    case $? in
        0) echo "active" ;;
        3) echo "inactive" ;;
        4) echo "not-installed" ;;
        *) echo "unknown" ;;
    esac
}

ensure_yq
validate_config

DEV_ETH=$(yq_read         '.network.dev_eth')
HOST_IP=$(yq_read         '.network.host_ip')
TARGET_IP=$(yq_read       '.network.target_ip')
IMG_NAME=$(yq_read        '.image.name')
IMG_PATH=$(yq_read        '.image.path')
USB_DEV=$(yq_read         '.usb.device')
ISO_URL=$(yq_read         '.usb.iso_url')
SERIAL_DEV=$(yq_read      '.serial.device')
SERIAL_BAUD=$(yq_read     '.serial.baud')
TARGET_DEV=$(yq_read      '.flash.target_device')
SUCCESS_MARKER=$(yq_read  '.flash.success_marker')
BOOT_MARKER=$(yq_read     '.flash.boot_marker')

ORIGINAL_SSH_STATE=$(detect_ssh_state)

{
    # === [0] Verify USB device exists ===
    if [ ! -b "$USB_DEV" ]; then
        echo "USB device '$USB_DEV' not found or is not a block device."
        echo "Make sure it's plugged in and use the full device path (e.g., /dev/sdb)"
        exit 1
    fi

    echo "=== [1] Configure dev Ethernet (${DEV_ETH}) ==="
    nmcli device set enp0s25 managed no 
    sudo ip link set "$DEV_ETH" down || true
    sudo ip addr flush dev "$DEV_ETH" || true
    sudo ip addr add "$HOST_IP/24" dev "$DEV_ETH"
    sudo ip link set "$DEV_ETH" up
    ip addr show "$DEV_ETH"

    sleep 2

    echo "=== [2] Start SSH (as needed) ==="
    if [ "$ORIGINAL_SSH_STATE" = "inactive" ]; then
        echo "Starting SSH temporarily for image transfer"
        sudo systemctl start sshd
    elif [ "$ORIGINAL_SSH_STATE" = "not-installed" ]; then
        echo "SSHD is not installed — skipping SSH-related steps."
    fi

    echo "=== [3] Download Alpine ISO ==="
    curl -L "$ISO_URL" -o alpine.iso

    echo "=== [4] Generate flash script ==="
    cat > "$FLASH_SCRIPT" <<EOF
#!/bin/sh
set -eux

echo "[SMARC] Setting static IP ${TARGET_IP} on eth0"
ip addr flush dev eth0 || true
ip addr add "${TARGET_IP}/24" dev eth0
ip link set eth0 up

echo "[SMARC] Detecting eMMC device..."
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
echo "[SMARC] Detected eMMC device: \$EMMC_DEV"

echo "[SMARC] Downloading image from ${HOST_IP}"
wget "http://${HOST_IP}:8000/${IMG_NAME}" -O "/mnt/${IMG_NAME}"

ls -lh "/mnt/${IMG_NAME}"

if [ ! -f "/mnt/${IMG_NAME}" ]; then
  echo "[SMARC] Image not found after wget!"
  exit 1
fi

echo "[SMARC] Zeroing first 64MB of \$EMMC_DEV"
dd if=/dev/zero of="\$EMMC_DEV" bs=1M count=64

echo "[SMARC] Flashing image to \$EMMC_DEV"
echo "[WARN] pv not found — flashing without progress bar"
dd if="/mnt/${IMG_NAME}" of="\$EMMC_DEV" bs=4M
sync

echo "[SMARC] Flash complete"
reboot
EOF
    chmod +x "$FLASH_SCRIPT"

    echo "=== [5] Create bootable USB with custom autorun ==="
    USB_DEV="$USB_DEV" FLASH_SCRIPT="$FLASH_SCRIPT" "$ISO_BUILDER"

    # === [5.5] Start HTTP server to serve image ===
    echo "=== [5.5] Starting temporary HTTP server to serve image"
    IMG_DIR=$(dirname "$IMG_PATH")
    cd "$IMG_DIR"
    python3 -m http.server 8000 > /dev/null 2>&1 &
    HTTP_PID=$!
    echo "HTTP server started with PID $HTTP_PID"
    cd - > /dev/null

    echo "Do you want to test this image in QEMU before inserting into SMARC? (y/N): "
    read -r qemu_choice

    if [[ "$qemu_choice" == "y" || "$qemu_choice" == "Y" ]]; then
        if [ ! -e /dev/kvm ]; then
            echo "/dev/kvm not found. Running without KVM acceleration."
            KVM=""
        else
            KVM="-enable-kvm"
        fi

        echo "Booting actual USB device ($USB_DEV) in QEMU..."
        sudo qemu-system-x86_64 \
             $KVM \
             -m 1024 \
             -machine type=pc,accel=kvm \
             -boot order=d \
             -drive file="$USB_DEV",format=raw,if=virtio,media=disk \
             -serial mon:stdio \
             -display none \
             -name "AlpineUSBTest"
    else
        # Eject USB
        sudo eject "$USB_DEV"

        echo "=== [6] Insert USB into SMARC board and power it up ==="
        echo "USB is ready. Insert into target and run flash-smarc.sh manually."
        echo "(Boot to Alpine, mount /dev/sda2, run /mnt/flash-smarc.sh)"
        echo "Please ensure the board is powered on and connected via serial (${SERIAL_DEV})."
        echo "Once ready, press ENTER to begin monitoring the serial port..."
        read -r

        echo "Monitoring serial output from SMARC via ${SERIAL_DEV}..."
        touch "$SERIAL_LOG"
        cat "$SERIAL_DEV" | tee "$RAW_SERIAL" "$SERIAL_LOG" &
        SERIAL_PID=$!

        echo "=== [7] Waiting for '${SUCCESS_MARKER}' ==="
        retry timeout 300 grep -q "$SUCCESS_MARKER" "$SERIAL_LOG"
        echo "Flash completed."

        echo "=== [8] Waiting for '${BOOT_MARKER}' ==="
        retry timeout 120 grep -q "$BOOT_MARKER" "$SERIAL_LOG"
        echo "System booted."

        echo "PASS" > "$STATUS_FILE"
    fi
} 2>&1 | tee -a "$ORCHESTRATOR_LOG" || {
    echo "❌ Flashing failed — logs in: $TMP_LOG_DIR" | tee -a "$ORCHESTRATOR_LOG"
    echo "FAIL" > "$STATUS_FILE"
    exit 1
}
