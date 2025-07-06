#!/bin/bash
# Script to create and test a bootable Alpine USB from config.yaml
# Combines flashing + QEMU test with YAML-config-driven values

## Create vanilla alpine version (for testing only)

set -euo pipefail

CONFIG="smarc_config.yaml"
YQ_BIN=".bin/yq"

# Step 1: Ensure yq exists
if [ ! -x "$YQ_BIN" ]; then
  echo "📦 Downloading yq..."
  mkdir -p .bin
  curl -L https://github.com/mikefarah/yq/releases/latest/download/yq_linux_amd64 -o "$YQ_BIN"
  chmod +x "$YQ_BIN"
fi

# Step 2: Extract vars from YAML
ALPINE_ISO_URL=$($YQ_BIN e '.usb.iso_url' "$CONFIG")
USB_DEV=$($YQ_BIN e '.usb.device' "$CONFIG")
ISO_NAME="alpine.iso"

# Step 3: Confirm device
if [[ "$USB_DEV" == "/dev/sdX" || -z "$USB_DEV" ]]; then
  echo "❌ USB device not set properly in config.yaml"
  exit 1
fi

# Step 4: Prompt before erase
read -rp "⚠️  This will erase all data on ${USB_DEV}1. Proceed? (y/N): " confirm
if [[ "$confirm" != "y" && "$confirm" != "Y" ]]; then
  echo "Aborted."
  exit 0
fi

# Step 5: Download ISO
if [ ! -f "$ISO_NAME" ]; then
  echo "⬇️  Downloading Alpine ISO from $ALPINE_ISO_URL..."
  curl -L "$ALPINE_ISO_URL" -o "$ISO_NAME"
else
  echo "✅ Alpine ISO already exists."
fi

# Step 6: Unmount USB and Flash
for part in $(lsblk -ln -o NAME "${USB_DEV}1" | tail -n +2); do
  sudo umount "/dev/$part" || true
done

echo "📝 Flashing $ISO_NAME to ${USB_DEV}1..."
ISO_SIZE=$(stat --format=%s "$ISO_NAME")
pv -s "$ISO_SIZE" "$ISO_NAME" | sudo dd of="${USB_DEV}1" bs=4M conv=fsync
sync

echo "✅ USB flash completed: ${USB_DEV}1"

# Step 7: Offer QEMU test
read -rp "🔁 Do you want to test this image in QEMU? (y/N): " qconfirm
if [[ "$qconfirm" == "y" || "$qconfirm" == "Y" ]]; then

  if [ ! -e /dev/kvm ]; then
    echo "⚠️  /dev/kvm not found. Running without KVM acceleration."
    KVM=""
  else
    KVM="-enable-kvm"
  fi

  echo "🚀 Booting USB in QEMU..."
  sudo qemu-system-x86_64 \
    $KVM \
    -m 1024 \
    -machine type=pc,accel=kvm \
    -boot order=d \
    -drive format=raw,file="${USB_DEV}1",if=virtio \
    -serial mon:stdio \
    -display none \
    -name "AlpineUSBTest"
else
  echo "💡 Skipping QEMU test. Insert USB into your target machine to boot."
fi
