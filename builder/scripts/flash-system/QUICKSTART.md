# Quick Start - Flash Ukama Boards

## TL;DR - Flash in 3 Steps

### Controller Board (RPi CM4 / Access Node)
```bash
# 1. Build image
cd builder/scripts/build-system
./build-access-image.sh

# 2. Connect hardware: Set BOOT jumper, connect USB, power on

# 3. Flash
cd ../flash-system
./orchestrate_board_flash.sh -c controller_config.yaml -b controller

# 4. Remove BOOT jumper, reboot
```

### Controller Board (Microchip SOM)
```bash
# 1. Download image from Google Drive (or build if available)

# 2. Prepare SD card (find device with 'lsblk')
cd builder/scripts/flash-system
nano microchip-controller_config.yaml  # Update paths

# 3. Flash (creates auto-flash SD card)
./orchestrate_board_flash.sh -c microchip-controller_config.yaml -b microchip-controller

# 4. Insert SD into board, power on, wait for auto-flash, remove SD after reboot
```

### COM Board (x86 SMARC)
```bash
# 1. Build image
cd builder/scripts/build-system
./build-com-image.sh

# 2. Prepare USB stick (find device with 'lsblk')
cd ../flash-system
nano smarc_config.yaml  # Update host_device.device to your USB stick

# 3. Flash (creates bootable USB)
./orchestrate_board_flash.sh -c smarc_config.yaml -b SMARC

# 4. Insert USB into COM board, power on, wait for auto-flash
```

### FEM Board (x86 SMARC)
```bash
# Same as COM board, but use:
./build-amplifier-image.sh  # for building
./orchestrate_board_flash.sh -c fem_config.yaml -b FEM-Control  # for flashing
```

## Image Locations

After building, images are here:
```
builder/scripts/build-system/
├── ukama-access-node.img      ← Controller (RPi CM4)
├── ukama-com-image.img        ← COM
└── ukama-amplifier-image.img  ← FEM
```

**Controller (Microchip SOM):** Download from Google Drive:
https://drive.google.com/file/d/1JuZ1EDS4p4mB_rid_gxyWzO-xecvMHeP/view

Config files already point to these locations (relative path `../build-system/`).

## What You Need

**Controller (RPi CM4):**
- ✅ USB cable
- ✅ BOOT jumper set
- ❌ No USB stick needed
- ❌ No network needed

**Controller (Microchip SOM):**
- ✅ SD card (8GB+, will be erased)
- ❌ No USB cable needed
- ❌ No network needed

**COM/FEM Boards:**
- ✅ USB stick (2GB+, will be erased)
- ✅ Ethernet cable
- ✅ Network interface on laptop
- ❌ No BOOT jumper needed

## Full Documentation

See [README.md](README.md) for detailed instructions and troubleshooting.
