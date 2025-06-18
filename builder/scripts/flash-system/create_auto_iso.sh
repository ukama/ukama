#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
# Copyright (c) 2025-present, Ukama Inc.

set -euo pipefail
set -x

ORIG_ISO="alpine.iso"
MNT_ISO="iso-mount"
MNT_USB="usb-mount"
APKOVL_DIR="apkovl-root"

: "${USB_DEV:?Must set USB_DEV}"
: "${HOSTNAME:=localhost}"

USB_PART="${USB_DEV}1"
APKOVL_FILE="${HOSTNAME}.apkovl.tar.gz"
FLASH_SCRIPT="flash-smarc.sh"

TEMP_DIRS=($MNT_ISO $MNT_USB $APKOVL_DIR)
TEMP_FILES=($APKOVL_FILE)

cleanup() {
    echo "🧹 Cleaning up..."
    for d in "${TEMP_DIRS[@]}"; do [ -d "$d" ] && sudo umount "$d" || true; rm -rf "$d"; done
    for f in "${TEMP_FILES[@]}"; do [ -f "$f" ] && rm -f "$f"; done
}
trap cleanup EXIT

verify_apkovl() {
    local apkovl_path="$1"
    local tmp_list="/tmp/apkovl-contents.txt"

    echo "🔍 Verifying overlay tarball at $apkovl_path..."

    if [ ! -f "$apkovl_path" ]; then
        echo "❌ ERROR: apkovl file missing: $apkovl_path"
        exit 1
    fi

    tar -tzf "$apkovl_path" > "$tmp_list"

    local required_files=(
        etc/hostname
        etc/local.d/autorun.start
        flash-smarc.sh
    )

    for f in "${required_files[@]}"; do
        if ! grep -q "$f" "$tmp_list"; then
            echo "❌ MISSING: $f in apkovl"
            exit 1
        fi
    done

    echo "✅ apkovl overlay verified: all required files present."
}

# 🔥 Wipe and format USB
echo "💣 Wiping USB device $USB_DEV..."
sudo wipefs -a "$USB_DEV"
sudo dd if=/dev/zero of="$USB_DEV" bs=1M count=10
sudo sync

echo "🔧 Partitioning USB device $USB_DEV"
sudo parted --script "$USB_DEV" \
  mklabel msdos \
  mkpart primary fat32 1MiB 100% \
  set 1 boot on

sudo mkfs.vfat -F 32 -n ALPINE_DATA "$USB_PART"

# 📁 Mount ISO and USB
mkdir -p "$MNT_ISO" "$MNT_USB"
sudo mount -o loop "$ORIG_ISO" "$MNT_ISO"
sudo mount "$USB_PART" "$MNT_USB"

# 📦 Copy ISO contents to USB
sudo rsync -a "$MNT_ISO"/ "$MNT_USB"/

read

echo "⚙️ Detect and patch boot config (GRUB, syslinux, or extlinux)"
GRUB_CFG="$MNT_USB/boot/grub/grub.cfg"
SYS_CFG="$MNT_USB/syslinux.cfg"
UEFI_CFG="$MNT_USB/boot/extlinux/extlinux.conf"

echo "🔍 Patching bootloader to enable 'data' mode for apkovl loading..."

if [ -f "$SYS_CFG" ]; then
    echo "🛠️  Patching syslinux.cfg..."
    sudo sed -i '/^APPEND / s|$| modules=loop,squashfs,sd-mod,usb-storage quiet data|' "$SYS_CFG"
elif [ -f "$UEFI_CFG" ]; then
    echo "🛠️  Patching extlinux.conf..."
    sudo sed -i '/^  APPEND / s|$| modules=loop,squashfs,sd-mod,usb-storage quiet data|' "$UEFI_CFG"
elif [ -f "$GRUB_CFG" ]; then
    echo "🛠️  Patching grub.cfg..."
    # ✅ Replace any existing line with a known-good one
    sudo sed -i '/^[[:space:]]*linux /c\    linux    /boot/vmlinuz-lts modules=loop,squashfs,sd-mod,usb-storage quiet data' "$GRUB_CFG"
else
    echo "⚠️  No known boot config found (syslinux, extlinux, or grub). apkovl may not be loaded."
fi

# 🧬 Add EFI boot if present
if [ -d "$MNT_ISO/EFI/BOOT" ]; then
    sudo mkdir -p "$MNT_USB/EFI/BOOT"
    sudo cp -r "$MNT_ISO/EFI/BOOT/"* "$MNT_USB/EFI/BOOT/"
fi

# 📦 Build apkovl overlay
mkdir -p "$APKOVL_DIR/etc/local.d"
mkdir -p "$APKOVL_DIR/etc/runlevels/default"
mkdir -p "$APKOVL_DIR/etc"

# Set hostname
echo "$HOSTNAME" > "$APKOVL_DIR/etc/hostname"

# Auto-run flash script
cat > "$APKOVL_DIR/etc/local.d/autorun.start" <<EOF
#!/bin/sh
echo "[AutoUSB] Starting autorun" > /flash.log
env >> /flash.log
echo "[AutoUSB] Running flash script..." >> /flash.log
/flash-smarc.sh >> /flash.log 2>&1
echo "[AutoUSB] Flash script completed" >> /flash.log
EOF
chmod +x "$APKOVL_DIR/etc/local.d/autorun.start"
ln -sf ../../init.d/local "$APKOVL_DIR/etc/runlevels/default/local"

# Copy flash script into overlay
cp "$FLASH_SCRIPT" "$APKOVL_DIR/flash-smarc.sh"
chmod +x "$APKOVL_DIR/flash-smarc.sh"

# Create and copy apkovl
tar -C "$APKOVL_DIR" -czf "$APKOVL_FILE" .
sudo mkdir -p "$MNT_USB/apkovl"
sudo cp "$APKOVL_FILE" "$MNT_USB/apkovl/"

# Verify apkovl contents before unmounting
APKOVL_USB_PATH="$MNT_USB/apkovl/$APKOVL_FILE"
verify_apkovl "$APKOVL_USB_PATH"

# Finish
sync
sudo umount "$MNT_USB"
sudo umount "$MNT_ISO"

echo "✅ Bootable USB created with hostname=${HOSTNAME} and autorun enabled."
