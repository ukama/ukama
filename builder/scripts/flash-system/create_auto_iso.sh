#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
# Copyright (c) 2025-present, Ukama Inc.

set -euo pipefail

ORIG_ISO="alpine.iso"
AUTO_ISO="alpine-auto.iso"
FLASH_SCRIPT="flash-smarc.sh"
ISO_MOUNT="iso-mount"
ISO_ROOT="iso-root"
APKOVL_DIR="apkovl-root"
APKOVL_FILE="custom.apkovl.tar.gz"

# Temporary work files to clean on exit
TEMP_DIRS=("$ISO_MOUNT" "$ISO_ROOT" "$APKOVL_DIR")
TEMP_FILES=("$APKOVL_FILE" "ORIG_ISO")

cleanup() {
    echo "🧹 Cleaning up temporary files and directories..."
    for d in "${TEMP_DIRS[@]}"; do
        [ -d "$d" ] && rm -rf "$d"
    done
    for f in "${TEMP_FILES[@]}"; do
        [ -f "$f" ] && rm -f "$f"
    done
}
trap cleanup EXIT

echo "📦 Creating auto-bootable Alpine ISO with apkovl..."

# Step 1: Mount original ISO read-only
mkdir -p "$ISO_MOUNT" "$ISO_ROOT"
sudo mount -o loop "$ORIG_ISO" "$ISO_MOUNT"
sudo rsync -a "$ISO_MOUNT"/ "$ISO_ROOT"/
sudo umount "$ISO_MOUNT"
rmdir "$ISO_MOUNT"

# Fix ownership and permissions
sudo chown -R "$USER:$USER" "$ISO_ROOT"
chmod -R u+w "$ISO_ROOT"

# Step 2: Create apkovl with autorun and flash script
mkdir -p "$APKOVL_DIR/etc/local.d"
mkdir -p "$APKOVL_DIR/etc/runlevels/default"

cat > "$APKOVL_DIR/etc/local.d/autorun.start" <<EOF
#!/bin/sh
echo "[AutoISO] Running flash script..."
/flash-smarc.sh > /flash.log 2>&1
EOF
chmod +x "$APKOVL_DIR/etc/local.d/autorun.start"
ln -sf /etc/init.d/local "$APKOVL_DIR/etc/runlevels/default/local"

cp "$FLASH_SCRIPT" "$APKOVL_DIR/flash-smarc.sh"
chmod +x "$APKOVL_DIR/flash-smarc.sh"

tar -C "$APKOVL_DIR" -czf "$APKOVL_FILE" .

# Step 3: Inject apkovl into ISO root
mkdir -p "$ISO_ROOT/apkovl"
cp "$APKOVL_FILE" "$ISO_ROOT/apkovl/"

# Step 4: Rebuild ISO
mkisofs -quiet -l -R -V "Ukama-AlpineAuto" \
    -b boot/syslinux/isolinux.bin \
    -c boot/syslinux/boot.cat \
    -no-emul-boot -boot-load-size 4 -boot-info-table \
    -o "$AUTO_ISO" "$ISO_ROOT"

echo "✅ Custom ISO created: $AUTO_ISO"
exit 1
