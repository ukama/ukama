# UkamaOS
[![build](https://github.com/ukama/ukamaOS/actions/workflows/ci.yaml/badge.svg)](https://github.com/ukama/ukamaOS/actions/workflows/ci.yaml)

Cloud-native and micro-services OS for Nodes

# UkamaDistro
UkamaDistro is lightweight busybox based distro build using musl-libc and few minimal services to boot up Ukama devices and make them functional.
Basic idea behind this distro is to have  what is needed.

# Supported archtectures:
- arm
- x86
- mips

## Firmware
Firmware consists fo the bootloader required to boot up Ukama devices. 
For x86 base device coreboot is used as rom boot loader and for arm based devices we have at91-bootstarp and u-boot as a boot loader.

### OS
Linux from the main line is used with few patches from SoC vendors and for Ukama Devices.

### Distro
Distro as mentioned above is leight weight busybox based on musl-libc.


# Build

## Prerequisites

```
sudo apt-get update
git submodule init
git submodule update
```

### Dependencies

```
sudo apt-get install bc build-essential git libncurses5-dev lzop perl libssl-dev gnat flex wget zlib1g-dev gcc-arm-linux-gnueabihf automake-1.15 bison python libelf-dev cmake curl libtool tcl pkg-config tcl pkg-config autopoint wget libisl-dev g++ texinfo texlive ghostscript
```

### Rust

```
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
cargo install cross
rustup target add armv7-unknown-linux-gnueabihf
rustup target add x86_64-unknown-linux-musl
```

### Buildah

```
. /etc/os-release
sudo sh -c "echo 'deb http://download.opensuse.org/repositories/devel:/kubic:/libcontainers:/stable/x${ID^}_${VERSION_ID}/ /' > /etc/apt/sources.list.d/devel:kubic:libcontainers:stable.list"
wget -nv https://download.opensuse.org/repositories/devel:kubic:libcontainers:stable/x${ID^}_${VERSION_ID}/Release.key -O Release.key
sudo apt-key add - < Release.key
sudo apt-get update -qq
sudo apt-get -qq -y install buildah
```

## Coreboot toolchain:

```
cd firmware/coreboot 
make crossgcc-i386 CPUS=$(nproc)
```

## Make

```
make TARGET=<cnode|anode|homenode>  
```

## Initramfs

```
make initramfs TARGET=<cnode|anode|homenode>
```

## Clean

```
make clean TARGET=<cnode|anode|homenode>
```

## Clean buid and toolchains used aswell

```
make distclean TARGET=<cnode|anode|homenode>
```

## Stand alone build
Each component could be build  individaully by providing target name to make.

### Distro/RootFS

```
make distro TARGET=<cnode|anode|homenode>
```

### Linux

```
make os TARGET=<cnode|anode|homenode>
```

### Firmware

```
make firmware TARGET=<cnode|anode|homenode>
```
