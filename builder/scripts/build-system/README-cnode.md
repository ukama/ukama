# cnode image

`build-cnode-image.sh` builds the image for the controller node (cnode), a Raspberry Pi CM5 or CM4.

## What you need

- An x86 Ubuntu machine with about 40 GB free
- These tools:

      sudo apt install -y fdisk kpartx dosfstools rsync wget unzip

## Build the image

From anywhere in the repo:

    sudo builder/scripts/build-system/build-cnode-image.sh

It builds the apps the cnode uses, then the OS, then the image. The first run takes a while. You get `cnode-image.img` in `builder/scripts/build-system/`.

To build only the image again after a change, add `--image-only`.

The apps in the image come from `builder/boards/common.config` and `builder/boards/cnode.config`.

The node finds its org through `dev.bootstrap.ukama.com`. To use another bootstrap server:

    sudo BOOTSTRAP_SERVER=<server> builder/scripts/build-system/build-cnode-image.sh

## Flash it

In the factory, the flash tool in the manufacturing repo (`tools/flash/cnode`) writes the image and gives the board its node ID.

To flash a board by hand:

1. Put the boot jumper on and connect the board's USB port to your machine.
2. Build [usbboot](https://github.com/raspberrypi/usbboot) and run:

       sudo ./rpiboot -d mass-storage-gadget64

3. The eMMC shows up as a new disk, for example `/dev/sda`. Check with `lsblk`, then:

       sudo dd if=cnode-image.img of=/dev/sdX bs=4M conv=fsync status=progress

4. Remove the jumper and power-cycle the board.

Login is `ukama` / `ukama`.

## Node ID

The image does not have a real node ID. Without one, bootstrap cannot find the node and it does not come online. The flash tool in the manufacturing repo takes a free ID from the factory and writes it to the board.
