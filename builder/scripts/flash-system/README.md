# Ukama Board Flashing Guide

This guide explains how to flash UkamaOS images to different board types using the unified orchestrator script.

## Overview

The flashing system supports multiple board types with different flashing methods:
- **Controller Board** (RPi CM4): Direct USB flashing via rpiboot
- **COM Board** (x86 SMARC): Network boot via USB stick
- **FEM Boards** (x86 SMARC): Network boot via USB stick

## Prerequisites

### System Requirements
- Ubuntu Linux (recommended) or any Linux distribution
- USB port(s)
- Network interface (for COM/FEM boards)
- Serial console cable (optional, for monitoring)

### Install Dependencies
```bash
sudo apt update
sudo apt install -y git build-essential libusb-1.0-0-dev curl
```

## Building Images

Before flashing, you need to build the images. Navigate to the build-system directory:

```bash
cd builder/scripts/build-system
```

### Build Controller Image (RPi CM4)
```bash
./build-access-image.sh
```
Output: `ukama-access-node.img` (created in current directory)

### Build COM Board Image
```bash
./build-com-image.sh
```
Output: `ukama-com-image.img` (created in current directory)

### Build FEM Board Image
```bash
./build-amplifier-image.sh
```
Output: `ukama-amplifier-image.img` (created in current directory)

## Flashing Instructions

Navigate to the flash-system directory:
```bash
cd builder/scripts/flash-system
```

---

## 1. Flashing Controller Board (RPi CM4)

### What is eMMC?
eMMC (embedded MultiMediaCard) is permanent storage soldered on the board. You don't need to remove or insert any memory card - it's already inside the board.

### Hardware Setup
1. Connect CM4 module to carrier board
2. **Enable USB boot mode**: Set BOOT jumper or switch to ON
3. Connect USB cable from carrier board to your Ubuntu laptop
4. Power on the board

### Configuration
Edit `controller_config.yaml`:
```bash
nano controller_config.yaml
```

Update the image path:
```yaml
image:
  name: "ukama-access-node.img"
  path: "/full/path/to/ukama-access-node.img"  # Update this!
```

### Flash the Board
```bash
./orchestrate_board_flash.sh -c controller_config.yaml -b controller
```

The script will:
- Wait for CM4 in USB boot mode
- Automatically detect the eMMC device
- Flash the image directly to eMMC
- Show progress

### Post-Flash Steps
1. Power off the board
2. **Disable USB boot mode**: Remove BOOT jumper or switch to OFF
3. Power on the board
4. Board will boot UkamaOS from eMMC

### Verify Boot (Optional)
```bash
./flash-access-node.sh --verify
```

---

## 2. Flashing COM Board (x86 SMARC)

### Hardware Setup
1. Prepare a USB stick (minimum 2GB, will be erased!)
2. Find the USB device name:
```bash
lsblk
# Example output: /dev/sdc (your USB stick)
```

### Configuration
Edit `smarc_config.yaml`:
```bash
nano smarc_config.yaml
```

Update these paths:
```yaml
image:
  name: "ukama-com-image.img"
  path: "/full/path/to/ukama-com-image.img"  # Update this!

host_device:
  device: "/dev/sdc"  # Update with your USB device from lsblk!

network:
  host_eth: "eth0"  # Update with your network interface (check with 'ip a')
```

### Flash the Board
```bash
./orchestrate_board_flash.sh -c smarc_config.yaml -b SMARC
```

The script will:
1. Setup network on your laptop (DHCP + HTTP server)
2. Create bootable USB stick with Alpine Linux
3. Eject USB when ready

### Hardware Steps
1. Remove USB stick from laptop
2. Insert USB stick into COM board
3. Power on COM board
4. Board will:
   - Boot Alpine Linux from USB
   - Connect to your laptop via network
   - Download UkamaOS image
   - Flash to internal eMMC
   - Reboot automatically

### Post-Flash Steps
1. Remove USB stick from COM board
2. Board will boot UkamaOS from eMMC
3. USB stick can be reused for other boards

---

## 3. Flashing FEM Boards (x86 SMARC)

Same process as COM board, but use `fem_config.yaml`:

```bash
# Edit config
nano fem_config.yaml

# Update image path and USB device
# Then flash
./orchestrate_board_flash.sh -c fem_config.yaml -b FEM-Control
```

---

## Troubleshooting

### Controller Board Issues

**CM4 not detected in USB boot mode:**
- Verify BOOT jumper is set correctly
- Check USB cable connection
- Try a different USB port
- Run `lsusb` and look for "Broadcom BCM2711 Boot"

**Permission denied errors:**
- Run with sudo: `sudo ./orchestrate_board_flash.sh ...`

**Missing tools:**
- Script will prompt to install automatically
- Or manually: `sudo apt install -y lsblk timeout lsusb libusb-1.0-0-dev git`

### COM/FEM Board Issues

**USB device not found:**
- Check USB stick is connected
- Verify device name with `lsblk`
- Make sure USB is not mounted: `sudo umount /dev/sdc*`

**Network issues:**
- Verify network interface name: `ip a`
- Update `host_eth` in config file
- Check Ethernet cable connection

**Board not booting from USB:**
- Verify USB boot is enabled in BIOS
- Try recreating the USB stick
- Check serial console for boot messages

**Image download fails:**
- Verify HTTP server is running (script starts it automatically)
- Check firewall settings: `sudo ufw status`
- Verify network connectivity between laptop and board

### Serial Console Monitoring

To monitor boot process:
```bash
# Find serial device
ls /dev/ttyUSB*

# Connect to serial console
screen /dev/ttyUSB0 115200
# Press Ctrl+A then K to exit
```

---

## Image Locations

After building, images are created in the build-system directory:

```
builder/scripts/build-system/
├── ukama-access-node.img      # Controller board image
├── ukama-com-image.img        # COM board image
└── ukama-amplifier-image.img  # FEM board image
```

Update the `path` field in config files to point to these images.

---

## Quick Reference

| Board | Config File | Flash Method | USB Stick Needed? |
|-------|-------------|--------------|-------------------|
| Controller | controller_config.yaml | rpiboot (direct USB) | ❌ No |
| COM | smarc_config.yaml | network boot | ✅ Yes |
| FEM | fem_config.yaml | network boot | ✅ Yes |

---

## Support

For issues or questions, contact the Ukama development team.
