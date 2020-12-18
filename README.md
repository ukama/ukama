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
sudo apt-get install bc build-essential git libncurses5-dev lzop perl libssl-dev gnat flex wget zlib1g-dev gcc-arm-linux-gnueabi
```
## Coreboot toolchain:
cd firmware/coreboot 
```
make crossgcc-i386 CPUS=$(nproc)
```
## Make
```
make TARGETBOARD=<cnode|anode|homenode>  
```
## Initramfs
```
make initramfs TARGETBOARD=<cnode|anode|homenode>
```
## Clean
```
make clean TARGETBOARD=<cnode|anode|homenode>
```
## Stand alone build
Each component could be build  individaully by providing target name to make.
### Distro/RootFS
```
make distro TARGETBOARD=<cnode|anode|homenode>
```
### Linux
```
make os TARGETBOARD=<cnode|anode|homenode>
```
### Firmware
```
make firmware TARGETBOARD=<cnode|anode|homenode>
```
