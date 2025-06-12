#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2025-present, Ukama Inc.

#!/bin/bash
set -euo pipefail

CONFIG="config.yaml"
YQ_BIN="./.bin/yq"
FLASH_SCRIPT="flash-smarc.sh"
ISO_BUILDER="./create_auto_iso.sh"

# === Logging Setup ===
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
TMP_LOG_DIR="logs/${TIMESTAMP}_UNKNOWN"
mkdir -p "$TMP_LOG_DIR"

ORCHESTRATOR_LOG="${TMP_LOG_DIR}/orchestrator.log"
SERIAL_LOG="${TMP_LOG_DIR}/serial_console.log"
RAW_SERIAL="${TMP_LOG_DIR}/serial_raw.log"
MAC_FILE="${TMP_LOG_DIR}/mac.txt"
SN_FILE="${TMP_LOG_DIR}/serial.txt"
STATUS_FILE="${TMP_LOG_DIR}/status.txt"
ORIGINAL_SSH_STATE=$(systemctl is-active sshd || echo "unknown")

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

ensure_yq
validate_config

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

{
    echo "=== [1] Configure dev Ethernet (${DEV_ETH}) ==="
    sudo ip addr flush dev "$DEV_ETH" || true
    sudo ip addr add "${STATIC_IP}/24" dev "$DEV_ETH"
    sudo ip link set dev "$DEV_ETH" up

    echo "=== [2] Start SSH (as needed) ==="
    if [ "$ORIGINAL_SSH_STATE" != "active" ]; then
        echo "🔐 Starting SSH temporarily for image transfer"
        sudo systemctl start sshd
    fi

    echo "=== [3] Download Alpine ISO ==="
    wget -O alpine.iso "$ISO_URL"

    echo "=== [4] Generate flash script ==="
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

    echo "=== [5] Build auto-run Alpine ISO ==="
    bash "$ISO_BUILDER"

    echo "=== [6] Flash ISO to USB ${USB_DEV} ==="
    sudo dd if=alpine-auto.iso of="${USB_DEV}" bs=4M status=progress && sync

    echo "=== [7] Insert USB into SMARC board and power it up ==="
    echo "⚠️  No user input is needed — SMARC will auto-run flash script."
    echo "📡 Monitoring serial port at ${SERIAL_DEV}..."

    cat "${SERIAL_DEV}" | tee "${RAW_SERIAL}" | head -n 100 > /dev/null &

    sleep 10  # Allow time for serial to emit logs

    MAC=$(grep -oE '([a-f0-9]{2}:){5}[a-f0-9]{2}' "$RAW_SERIAL" | head -n1 || true)
    SN=$(grep -E '^.*-[0-9A-Fa-f]{4,}$' "$RAW_SERIAL" | head -n1 || true)

    [ -n "$MAC" ] && echo "$MAC" > "$MAC_FILE"
    [ -n "$SN" ] && echo "$SN" > "$SN_FILE"

    MAC_CLEAN=$(echo "$MAC" | tr -d ':' | tr '[:lower:]' '[:upper:]')
    SN_CLEAN=$(echo "$SN" | tr -d ' ' | tr '[:lower:]' '[:upper:]' | tr -cd '[:alnum:]-')

    NEW_LOG_DIR="logs/${TIMESTAMP}_${MAC_CLEAN}_${SN_CLEAN}"
    mv "$TMP_LOG_DIR" "$NEW_LOG_DIR"

    ORCHESTRATOR_LOG="${NEW_LOG_DIR}/orchestrator.log"
    SERIAL_LOG="${NEW_LOG_DIR}/serial_console.log"
    STATUS_FILE="${NEW_LOG_DIR}/status.txt"

    echo "=== [8] Waiting for '${SUCCESS_MARKER}' ==="
    timeout 300 grep -q "$SUCCESS_MARKER" < <(tee "$SERIAL_LOG" < "$SERIAL_DEV")
    echo "✅ Flash completed."

    echo "=== [9] Waiting for '${BOOT_MARKER}' ==="
    timeout 120 grep -q "$BOOT_MARKER" < <(tee -a "$SERIAL_LOG" < "$SERIAL_DEV")
    echo "✅ System booted."

    echo "PASS" > "$STATUS_FILE"

    if [ "$ORIGINAL_SSH_STATE" != "active" ]; then
        echo "🧹 Restoring SSH state — stopping SSHD"
        sudo systemctl stop sshd
    fi
} 2>&1 | tee -a "$ORCHESTRATOR_LOG" || {
    echo "❌ Flashing failed — logs in: $TMP_LOG_DIR"
    echo "FAIL" > "$STATUS_FILE"

    if [ "$ORIGINAL_SSH_STATE" != "active" ]; then
        echo "🧹 Cleaning up: stopping SSHD"
        sudo systemctl stop sshd
    fi
    exit 1
}
