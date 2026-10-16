# cnode image

Builds the image for the controller node (Raspberry Pi CM5 or CM4).

## Requirements

- An x86 Ubuntu machine with about 40 GB free
- These packages:

      sudo apt install -y fdisk kpartx dosfstools rsync wget unzip jq binutils

- The repo with its submodules:

      git clone --recurse-submodules https://github.com/ukama/ukama.git

  If you already have the repo, run `git submodule update --init --recursive`.

## Build

From the root of the repo:

    sudo builder/scripts/build-system/build-cnode-image.sh

It needs sudo because it builds in a chroot and mounts the image. It takes a while. The image is saved as `builder/scripts/build-system/cnode-image.img`.

The apps are listed in `builder/boards/common.config` and `builder/boards/cnode.config`. The node uses `dev.bootstrap.ukama.com` by default. To use another server:

    sudo BOOTSTRAP_SERVER=<server> builder/scripts/build-system/build-cnode-image.sh

## Check

The build checks the image at the end with `validate-cnode-image.sh`. To check it again:

    sudo builder/scripts/build-system/validate-cnode-image.sh builder/scripts/build-system/cnode-image.img

## Flash

In the factory, use the flash tool in the manufacturing repo (`tools/flash/cnode`). It also gives the board its node ID.

To flash by hand:

1. Put the boot jumper on and connect the board to your machine with USB.
2. Build [usbboot](https://github.com/raspberrypi/usbboot) and run:

       sudo ./rpiboot -d mass-storage-gadget64

3. Find the new disk with `lsblk` and write the image to it:

       sudo dd if=builder/scripts/build-system/cnode-image.img of=/dev/sdX bs=4M conv=fsync status=progress

4. Remove the jumper and power cycle the board.

Login is `ukama` / `ukama`.

## Node ID

The image has no node ID. In the image, `/ukama/mocksysfs/sys` is a link to the board's real `/sys`. The flash tool replaces it with files that hold the board's node ID. Without them the node can't come online.
