#!/bin/bash
set -euo pipefail

CONFIG="smarc_config.yaml"
YQ_BIN="./.bin/yq"
FLASH_SCRIPT="flash-smarc.sh"
ISO_BUILDER="./create_auto_iso.sh"
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
    if [ "$ORIGINAL_SSH_STATE" != "active" ]; then
        echo "🧹 Restoring SSH state — stopping SSHD" | tee -a "$ORCHESTRATOR_LOG"
        sudo systemctl stop sshd || true
    fi
    rm -f "$YQ_BIN"
    rm -f alpine.iso
}
trap cleanup EXIT

REQUIRED_KEYS=(
    ".network.dev_eth" ".network.static_ip" ".network.target_ip"
    ".image.name" ".image.path"
    ".usb.device" ".usb.iso_url"
    ".serial.device" ".serial.baud"
    ".flash.target_device" ".flash.success_marker" ".flash.boot_marker"
    ".system.target_hostname"
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

retry() {
    local n=1 max=$RETRIES delay=5
    until "$@"; do
        if (( n == max )); then
            echo "❌ Command failed after $n attempts." | tee -a "$ORCHESTRATOR_LOG"
            exit 1
        else
            echo "🔁 Retry $n/$max: $*" | tee -a "$ORCHESTRATOR_LOG"
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
STATIC_IP=$(yq_read       '.network.static_ip')
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
TARGET_HOSTNAME=$(yq_read '.system.target_hostname')

ORIGINAL_SSH_STATE=$(detect_ssh_state)

{
    echo "=== [1] Configure dev Ethernet (${DEV_ETH}) ==="
    sudo ip addr flush dev "$DEV_ETH" || true
    sudo ip addr add "${STATIC_IP}/24" dev "$DEV_ETH"
    sudo ip link set dev "$DEV_ETH" up

    echo "=== [2] Start SSH (as needed) ==="
    if [ "$ORIGINAL_SSH_STATE" = "inactive" ]; then
        echo "🔐 Starting SSH temporarily for image transfer"
        sudo systemctl start sshd
    elif [ "$ORIGINAL_SSH_STATE" = "not-installed" ]; then
        echo "⚠️ SSHD is not installed — skipping SSH-related steps."
    fi

    echo "=== [3] Download Alpine ISO ==="
    curl -L "$ISO_URL" -o alpine.iso

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

    echo "=== [5] Create bootable USB with custom autorun ==="
    USB_DEV="$USB_DEV" HOSTNAME="$TARGET_HOSTNAME" "$ISO_BUILDER"

    echo "🔁 Do you want to test this image in QEMU before inserting into SMARC? (y/N): "
    read -r qemu_choice

    if [[ "$qemu_choice" == "y" || "$qemu_choice" == "Y" ]]; then
        if [ ! -e /dev/kvm ]; then
            echo "⚠️  /dev/kvm not found. Running without KVM acceleration."
            KVM=""
        else
            KVM="-enable-kvm"
        fi

        echo "🚀 Booting actual USB device ($USB_DEV) in QEMU..."
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
        echo "⚠️  No user interaction required — SMARC will auto-run the flash script from USB."
        echo "🔌 Please ensure the board is powered on and connected via serial (${SERIAL_DEV})."
        echo "⏳ Press ENTER to begin monitoring the serial port..."
        read -r

        echo "📡 Monitoring serial output from SMARC via ${SERIAL_DEV}..."
        touch "$SERIAL_LOG"
        cat "$SERIAL_DEV" | tee "$RAW_SERIAL" "$SERIAL_LOG" &
        SERIAL_PID=$!

        echo "=== [7] Waiting for '${SUCCESS_MARKER}' ==="
        retry timeout 300 grep -q "$SUCCESS_MARKER" "$SERIAL_LOG"
        echo "✅ Flash completed."

        echo "=== [8] Waiting for '${BOOT_MARKER}' ==="
        retry timeout 120 grep -q "$BOOT_MARKER" "$SERIAL_LOG"
        echo "✅ System booted."

        echo "PASS" > "$STATUS_FILE"
    fi
} 2>&1 | tee -a "$ORCHESTRATOR_LOG" || {
    echo "❌ Flashing failed — logs in: $TMP_LOG_DIR" | tee -a "$ORCHESTRATOR_LOG"
    echo "FAIL" > "$STATUS_FILE"
    exit 1
}
