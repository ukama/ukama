# Quick Start - Flash Ukama Boards

## TL;DR - Flash in 3 Steps

### Controller Board (RPi CM4)
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
├── ukama-access-node.img      ← Controller
├── ukama-com-image.img        ← COM
└── ukama-amplifier-image.img  ← FEM
```

Config files already point to these locations (relative path `../build-system/`).

## What You Need

**Controller Board:**
- ✅ USB cable
- ✅ BOOT jumper set
- ❌ No USB stick needed
- ❌ No network needed

**COM/FEM Boards:**
- ✅ USB stick (2GB+, will be erased)
- ✅ Ethernet cable
- ✅ Network interface on laptop
- ❌ No BOOT jumper needed

## Full Documentation

See [README.md](README.md) for detailed instructions and troubleshooting.
