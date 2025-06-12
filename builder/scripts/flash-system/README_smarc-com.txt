Ukama SMARC Module Auto-Flashing Orchestrator
============================================

This tool automates the process of flashing a Ukama x86_64 image to SMARC modules using
an I-Pi Elkhart Lake dev board. It builds a custom Alpine ISO that auto-executes a flash
script and logs each flashing session including MAC address, serial number, and PASS/FAIL status.

What This Tool Does:
--------------------
    
    o Builds a custom Alpine Linux ISO that auto-runs a flash script on boot
    o Flashes the ISO to a USB stick
    o Boots the SMARC board and automatically flashes the Ukama OS image to eMMC
    o Temporarily enables SSH if not already active
    o Monitors the process over serial connection
    o Logs MAC address, serial number, and flash result

Requirements:

    o Linux host machine
    o USB port to flash the ISO
    o Ethernet connection between host and SMARC board
    o Serial USB connection to the SMARC board (e.g. /dev/ttyUSB0)

Packages to install:
    > sudo apt install curl wget rsync xorriso syslinux-utils coreutils

Configuration (config.yaml):

network:
dev_eth: "eth1"
static_ip: "192.168.53.100"
target_ip: "192.168.53.151"

image:
name: "ukama-os.img"
path: "/home/factory/images/ukama-os.img"

usb:
device: "/dev/sdX"
iso_url: "https://dl-cdn.alpinelinux.org/alpine/v3.20/releases/x86_64/alpine-standard-3.20.0-x86_64.iso"

serial:
device: "/dev/ttyUSB0"
baud: 115200

flash:
target_device: "/dev/mmcblk0"
success_marker: "Flash complete. Rebooting."
boot_marker: "login:"

Operator Workflow:
----------------
    
Step 1: Prepare Host

    o Connect Ethernet from host to SMARC board
    o Connect USB-to-serial from host to SMARC board (/dev/ttyUSB0)
    o Insert USB stick into host for flashing ISO

Step 2: Run the Orchestrator

    > ./orchestrate_smarc_flash.sh

This will:
    o Validate config
    o Generate flash-smarc.sh
    o Build a custom Alpine ISO that runs the flash script
    o Flash the ISO to the USB stick
    o Start SSH temporarily if needed
    o Prompt to insert USB into SMARC board

Step 3: Boot SMARC Board

    o Insert USB into SMARC board
    o Power on
    o Do not log in or touch anything
    o SMARC will boot and automatically:
    o Run flash-smarc.sh
    o Pull Ukama OS image from host
    o Flash to eMMC
    o Reboot into new system

Step 4: Monitor and Review Logs
    Logs will be saved to:

    logs/
    └── YYYYMMDD_HHMMSS_MAC_SERIAL/
    ├── orchestrator.log
    ├── serial_console.log
    ├── serial_raw.log
    ├── mac.txt
    ├── serial.txt
    └── status.txt (contains PASS or FAIL)

Step 5: On Success
    o Messages will appear in terminal:
        "Flash completed."
        "System booted."
    o Remove USB
    o Proceed to next unit
