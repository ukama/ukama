# micorInit for UkamaOS
 
# microInit
microInit is a very minimal init required to boot and prepare system for microCE to run under UkamaOS.  
It has three major components:
- ### init
 init basically checks filesystem we are booting in like ramfs/tmpfs or ext2 and then do the exec to busybox init.
 
- ### rcinit
rcinit is similar to rc scripts. It scans through the /etc/init.d/ and execute all the executables in sequential manner. 
rcinit also does all the necessary mounts and module installation required by system.

- ### sysinit
sysinit starts the init bundles (using crun ) and preparing system to run microCE.  

# Build

## Prerequisites
Install some basic packages required for armv7.
```
sudo apt install binutils-arm-linux-gnueabihf
sudo apt-get install gcc-arm-linux-gnueabihf
```

Install rust and armv7 target.
```
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
rustup update
rustup target add armv7-unknown-linux-musleabihf

```

## Make 
```
cd distro/system/init
make TARGETBOARD=<cnode|anode|homenode>  
```

## Clean
```
cd distro/system/init
make clean TARGETBOARD=<cnode|anode|homenode>
```

## Create cpio archive (Only suported for x86 currently)
#### Command:
```
cd distro/system/init/microInit/test/virt
sudo ./createcpio.sh
```
#### Output:
```
Cleaning microfs...!
Copying init files...!
'../../init/target/x86_64-unknown-linux-musl/release/init' -> './microfs/init'
'../../rcinit/target/x86_64-unknown-linux-musl/release/rcinit' -> './microfs/bin/rc.init'
'../../sysinit/config/microInit.toml' -> './microfs/etc/microInit.toml'
'../../sysinit/target/x86_64-unknown-linux-musl/release/sysinit' -> './microfs/usr/bin/usysinit'
Creating cpio archive for test.
Use command: sudo qemu-system-x86_64 -m 512M -kernel ./virt-kernel  -initrd ./microfs.img -append "console=ttyAMA0 console=tty0 console=ttyS0 rdinit=/init" -serial stdio
```
# Test
Testing is done by spinning up a virtual enviornment using qemu. Create cpio archive step prepares the rootfs image for the test.
It includes busybox, our recently build init binaries and a couple of basic init bundles like sysctl and dhcpcd. Once these are executed it starts a dummy ukamaCE and enable shell to play with.
Command:
```
cd distro/system/init/microInit/test/virt
sudo qemu-system-x86_64 -m 512M -kernel ./virt-kernel  -initrd ./microfs.img -append "console=ttyAMA0 console=tty0 console=ttyS0 rdinit=/init" -serial stdio
```

