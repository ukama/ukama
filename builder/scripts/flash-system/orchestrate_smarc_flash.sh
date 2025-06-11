#!/bin/bash
#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2025-present, Ukama Inc.

set -euo pipefail

CONFIG="config.yaml"
YQ_BIN="./.bin/yq"
FLASH_SCRIPT="flash-smarc.sh"

# === Setup Logging ===
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
TMP_LOG_DIR="logs/${TIMESTAMP}_UNKNOWN"
mkdir -p "$TMP_LOG_DIR"

ORCHESTRATOR_LOG="${TMP_LOG_DIR}/orchestrator.log"
SERIAL_LOG="${TMP_LOG_DIR}/serial_console.log"
RAW_SERIAL="${TMP_LOG_DIR}/serial_raw.log"
MAC_FILE="${TMP_LOG_DIR}/mac.txt"
SN_FILE="${TMP_LOG_DIR}/serial.txt"
STATUS_FILE="${TMP_LOG_DIR}/status.txt"

# Required config keys
REQUIRED_KEYS=(
    ".network.dev_eth"
    ".network.static_ip"
    ".network.target_ip"
    ".image.name"
    ".image.path"
    ".usb.device"
    ".usb.iso_url"
    ".serial.device"
    ".serial.baud"
    ".flash.target_device"
    ".flash.success_marker"
    ".flash.boot_marker"
)

# === Helpers ===
ensure_yq() {
    if [ ! -x "$YQ_BIN" ]; then
        echo "📦 Downloading yq..." | tee -a "$ORCHESTRATOR_LOG"
        mkdir -p .bin
        curl -L https://github.com/mikefarah/yq/releases/latest/download/yq_linux_amd64 -o "$YQ_BIN"
        chmod +x "$YQ_BIN"
    fi
}

yq_read() {
    "$YQ_BIN" eval "$1" "$CONFIG"
}

validate_config() {
    echo "🔍 Validating config..." | tee -a "$ORCHESTRATOR_LOG"
    for key in "${REQUIRED_KEYS[@]}"; do
        if ! "$YQ_BIN" eval "$key" "$CONFIG" &>/dev/null; then
            echo "❌ Missing config: $key" | tee -a "$ORCHESTRATOR_LOG"
            exit 1
        fi
    done
}

# === Init ===
ensure_yq
validate_config

# === Load Config ===
DEV_ETH=$(yq_read        '.network.dev_eth')
STATIC_IP=$(yq_read      '.network.static_ip')
TARGET_IP=$(yq_read      '.network.target_ip')
IMG_NAME=$(yq_read       '.image.name')
IMG_PATH=$(yq_read       '.image.path')
USB_DEV=$(yq_read        '.usb.device')
ISO_URL=$(yq_read        '.usb.iso_url')
SERIAL_DEV=$(yq_read     '.serial.device')
SERIAL_BAUD=$(yq_read    '.serial.baud')
TARGET_DEV=$(yq_read     '.flash.target_device')
SUCCESS_MARKER=$(yq_read '.flash.success_marker')
BOOT_MARKER=$(yq_read    '.flash.boot_marker')

# === Main Process ===
{
    echo "=== [1] Configuring ${DEV_ETH} with static IP ==="
    sudo ip addr flush dev "$DEV_ETH" || true
    sudo ip addr add "${STATIC_IP}/24" dev "$DEV_ETH"
    sudo ip link set dev "$DEV_ETH" up

    echo "=== [2] Starting SSH server ==="
    sudo systemctl start sshd

    echo "=== [3] Writing Alpine ISO to USB ${USB_DEV} ==="
    wget -O alpine.iso "$ISO_URL"
    sudo dd if=alpine.iso of="$USB_DEV" bs=4M status=progress && sync

    echo "=== [4] Writing SMARC flash script ==="
    cat > "$FLASH_SCRIPT" <<EOF
#!/bin/sh
set -e
udhcpc -i eth0
scp root@${STATIC_IP}:${IMG_PATH} /tmp/${IMG_NAME}
dd if=/tmp/${IMG_NAME} of=${TARGET_DEV} bs=4M status=progress && sync
echo "[SMARC] ${SUCCESS_MARKER}"
reboot
EOF
    chmod +x "$FLASH_SCRIPT"

    echo "=== [5] Boot SMARC board and connect UART ==="
    echo "➡️  Press ENTER when you see Alpine prompt"
    read -r

    # Start serial log capture
    cat "$SERIAL_DEV" | tee "$RAW_SERIAL" | head -n 100 > /dev/null &

    echo "ℹ️  On the SMARC serial console, run:"
    echo "   ip link show eth0"
    echo "   dmidecode -s system-serial-number"
    echo "➡️  Press ENTER when done"
    read -r

    # === Extract MAC and Serial ===
    MAC=$(grep -oE '([a-f0-9]{2}:){5}[a-f0-9]{2}' "$RAW_SERIAL" | head -n1 || true)
    SN=$(grep -E '^.*-[0-9A-Fa-f]{4,}$' "$RAW_SERIAL" | head -n1 || true)

    [ -n "$MAC" ] && echo "$MAC" > "$MAC_FILE"
    [ -n "$SN" ]  && echo "$SN" > "$SN_FILE"

    echo "MAC: ${MAC:-Not found}"   | tee -a "$ORCHESTRATOR_LOG"
    echo "Serial: ${SN:-Not found}" | tee -a "$ORCHESTRATOR_LOG"

    # Rename log folder
    MAC_CLEAN=$(echo "$MAC" | tr -d ':' | tr '[:lower:]' '[:upper:]')
    SN_CLEAN=$(echo "$SN"   | tr -d ' ' | tr '[:lower:]' '[:upper:]' | tr -cd '[:alnum:]-')
    NEW_LOG_DIR="logs/${TIMESTAMP}_${MAC_CLEAN}_${SN_CLEAN}"

    mv "$TMP_LOG_DIR" "$NEW_LOG_DIR"
    ORCHESTRATOR_LOG="${NEW_LOG_DIR}/orchestrator.log"
    SERIAL_LOG="${NEW_LOG_DIR}/serial_console.log"
    STATUS_FILE="${NEW_LOG_DIR}/status.txt"

    # === Flash Success Marker ===
    echo "=== [6] Waiting for flash success marker: $SUCCESS_MARKER ==="
    timeout 300 grep -q "$SUCCESS_MARKER" < <(tee "$SERIAL_LOG" < "$SERIAL_DEV")
    echo "✅ Flash succeeded."

    echo "=== [7] Waiting for OS boot marker: $BOOT_MARKER ==="
    timeout 120 grep -q "$BOOT_MARKER" < <(tee -a "$SERIAL_LOG" < "$SERIAL_DEV")
    echo "✅ System booted successfully."

    echo "PASS" > "$STATUS_FILE"

} 2>&1 | tee -a "$ORCHESTRATOR_LOG" || {
    echo "❌ Flashing failed. See logs in $TMP_LOG_DIR or $NEW_LOG_DIR"
    echo "FAIL" > "$STATUS_FILE"
    exit 1
}
