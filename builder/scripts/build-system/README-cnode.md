# cnode image

`build-cnode-image.sh` builds the image for the controller node (cnode), a Raspberry Pi CM5 or CM4.

## What you need

- An x86 Ubuntu machine with about 40 GB free
- These tools:

      sudo apt install -y fdisk kpartx dosfstools rsync wget unzip jq binutils

- The repo with its submodules, because the apps are built from them:

      git clone --recurse-submodules https://github.com/ukama/ukama.git

  If you already have the repo, run `git submodule update --init --recursive`.

## Build the image

From the root of the repo:

    sudo builder/scripts/build-system/build-cnode-image.sh

It builds the apps, the OS and then the image. It takes a while. You get `builder/scripts/build-system/cnode-image.img`.

- It needs sudo because it builds in an Alpine chroot and mounts the image.
- When it ends, even if it fails, it deletes its temporary folders `buildenv/`, `rootfs/` and `mocksysfs/`, and gives the files it built back to you. Nothing is left owned by root.
- The apps come from `builder/boards/common.config` and `builder/boards/cnode.config`.
- The node finds its org through `dev.bootstrap.ukama.com`. To use another server:

      sudo BOOTSTRAP_SERVER=<server> builder/scripts/build-system/build-cnode-image.sh

### Log messages you can ignore

They come from the builder tools and do not break the image:

- `Error loading shared library libusys.so` and `Error relocating ...`: `ldd` cannot find `libusys.so` in the build chroot. The image gets it from `build/libs/vendor_libs.tgz`.
- `ERROR ... [build-misc] section parsing error`: logged for every app without a `[build-misc]` section, which is optional.
- `make: Circular aggregator <- aggregator dependency dropped`: a warning from the aggregator Makefile.
- `Unsupported architecture: aarch64`: the rootfs script ends with this on aarch64. The build script then checks the rootfs itself.
- `If you want to remove ... directory, run: ... destroy --remove`: the build script does this at the end.

## Check the image

    sudo builder/scripts/build-system/validate-cnode-image.sh builder/scripts/build-system/cnode-image.img

It checks the image without changing it. It reports 2 FAILs, `primary mock sysfs` and `passive mock sysfs`, because the flash tool adds the sysfs to each board (see Node ID).

## Flash it

In the factory, the flash tool in the manufacturing repo (`tools/flash/cnode`) writes the image and gives the board its node ID.

To flash a board by hand:

1. Put the boot jumper on and connect the board's USB port to your machine.
2. Build [usbboot](https://github.com/raspberrypi/usbboot) and run:

       sudo ./rpiboot -d mass-storage-gadget64

3. The eMMC shows up as a new disk, for example `/dev/sda`. Check with `lsblk`, then:

       sudo dd if=builder/scripts/build-system/cnode-image.img of=/dev/sdX bs=4M conv=fsync status=progress

4. Remove the jumper and power-cycle the board.

Login is `ukama` / `ukama`.

## Node ID

The image has no node ID. The flash tool takes a free ID from the factory and writes it to the board in `/ukama/mocksysfs/sys`. Without it, bootstrap cannot find the node and it does not come online.
