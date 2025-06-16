#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
# Copyright (c) 2025-present, Ukama Inc.

set -euo pipefail
set -x

ORIG_ISO="alpine.iso"
AUTO_ISO="alpine-auto-usb"
FLASH_SCRIPT="flash-smarc.sh"
MNT_ISO="iso-mount"
MNT_USB="usb-mount"
APKOVL_DIR="apkovl-root"
APKOVL_FILE="custom.apkovl.tar.gz"
USB_PART="${USB_DEV}1"

TEMP_DIRS=($MNT_ISO $MNT_USB $APKOVL_DIR)
TEMP_FILES=($APKOVL_FILE)

cleanup() {
    echo "🧹 Cleaning up..."
    for d in "${TEMP_DIRS[@]}"; do [ -d "$d" ] && sudo umount "$d" || true; rm -rf "$d"; done
    for f in "${TEMP_FILES[@]}"; do [ -f "$f" ] && rm -f "$f"; done
}
trap cleanup EXIT

# Step 1: Format USB device
echo "🔧 Partitioning USB device $USB_DEV"
[[ "$USB_DEV" =~ [0-9]+$ ]] && {
    echo "❌ ERROR: usb.device must be a full block device (e.g. /dev/sdb), not a partition like /dev/sdb1"
    exit 1
}
sudo parted --script "$USB_DEV" \
  mklabel msdos \
  mkpart primary fat32 1MiB 100% \
  set 1 boot on

sudo mkfs.vfat -F 32 -n UKAMA_ALPINE "$USB_PART"

# Step 2: Mount ISO and USB
mkdir -p "$MNT_ISO" "$MNT_USB"
sudo mount -o loop "$ORIG_ISO" "$MNT_ISO"
sudo mount "$USB_PART" "$MNT_USB"

# Step 3: Copy Alpine ISO files to USB
sudo rsync -a "$MNT_ISO"/ "$MNT_USB"/

# Step 3.5: Enable serial logging based on boot mode
echo "🔍 Detecting boot mode configuration (BIOS vs UEFI)..."

SYS_CFG="$MNT_USB/syslinux.cfg"
UEFI_CFG="$MNT_USB/boot/extlinux/extlinux.conf"  # Alpine often places this for UEFI

if [ -f "$SYS_CFG" ]; then
    echo "🛠️  Detected BIOS boot (syslinux), patching serial console"
    if ! grep -q '^SERIAL' "$SYS_CFG"; then
        echo 'SERIAL 0 115200' | sudo tee -a "$SYS_CFG" > /dev/null
    fi
    sudo sed -i '/^APPEND / s|$| console=ttyS0,115200|' "$SYS_CFG"
elif [ -f "$UEFI_CFG" ]; then
    echo "🛠️  Detected UEFI boot (extlinux), patching serial console"
    sudo sed -i '/^  APPEND / s|$| console=ttyS0,115200|' "$UEFI_CFG"
else
    echo "⚠️  No known boot config found (syslinux or extlinux), skipping serial patch"
fi

# Step 4: Install syslinux bootloader
sudo syslinux --install "$USB_PART"

# Optional: Add UEFI boot support
if [ -d "$MNT_ISO/EFI/BOOT" ]; then
    sudo mkdir -p "$MNT_USB/EFI/BOOT"
    sudo cp -r "$MNT_ISO/EFI/BOOT/"* "$MNT_USB/EFI/BOOT/"
fi

# Step 5: Create custom apkovl with autorun
mkdir -p "$APKOVL_DIR/etc/local.d"
mkdir -p "$APKOVL_DIR/etc/runlevels/default"

cat > "$APKOVL_DIR/etc/local.d/autorun.start" <<EOF
#!/bin/sh
echo "[AutoUSB] Running flash script..."
/flash-smarc.sh > /flash.log 2>&1
EOF
chmod +x "$APKOVL_DIR/etc/local.d/autorun.start"
ln -sf /etc/init.d/local "$APKOVL_DIR/etc/runlevels/default/local"

cp "$FLASH_SCRIPT" "$APKOVL_DIR/flash-smarc.sh"
chmod +x "$APKOVL_DIR/flash-smarc.sh"

tar -C "$APKOVL_DIR" -czf "$APKOVL_FILE" .
mkdir -p "$MNT_USB/apkovl"
sudo cp "$APKOVL_FILE" "$MNT_USB/apkovl/"

# Step 6: Sync and unmount
sync
sudo umount "$MNT_USB"
sudo umount "$MNT_ISO"

echo "✅ Bootable USB created successfully with autorun and serial logging."
