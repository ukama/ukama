#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2025-present, Ukama Inc.

#!/bin/bash
set -euo pipefail

ORIG_ISO="alpine.iso"
AUTO_ISO="alpine-auto.iso"
FLASH_SCRIPT="flash-smarc.sh"

echo "📦 Creating auto-bootable Alpine ISO..."

# Step 1: Mount original ISO
mkdir -p iso-mount iso-root
sudo mount -o loop "$ORIG_ISO" iso-mount
rsync -a iso-mount/ iso-root/
sudo umount iso-mount
rmdir iso-mount

# Step 2: Inject autorun script
mkdir -p iso-root/etc/local.d
cat > iso-root/etc/local.d/autorun.start <<EOF
#!/bin/sh
echo "[AutoISO] Running flash script..."
/flash-smarc.sh > /flash.log 2>&1
EOF
chmod +x iso-root/etc/local.d/autorun.start

# Step 3: Inject flash script
cp "$FLASH_SCRIPT" iso-root/flash-smarc.sh
chmod +x iso-root/flash-smarc.sh

# Step 4: Enable autorun in Alpine
echo "local" >> iso-root/etc/runlevels/default

# Step 5: Rebuild ISO
mkisofs -quiet -l -R -V "AlpineAuto" \
        -b boot/syslinux/isolinux.bin \
        -c boot/syslinux/boot.cat \
        -no-emul-boot -boot-load-size 4 -boot-info-table \
        -o "$AUTO_ISO" iso-root

echo "✅ Custom ISO created: $AUTO_ISO"
