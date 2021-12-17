#!/bin/sh

# Only meant for test purpose. This temporary test file will be replaced by ukamaOS virtual test suit which
# will mock devices.   

#Run script with sudo
if [ "$EUID" -ne 0 ]; then
    echo "Please run as root"
    exit
fi

#Cleaning files
echo "Cleaning microfs...!"
sudo rm -rf ./microfs/preInit ./microfs/init ./microfs/sbin/sysInit ./microfs/sbin/init
sync

#Copying file
echo "Copying init files...!"
sudo cp -v ../../../../init ./microfs/sbin/init
sudo cp -v ../../preInit/target/x86_64-unknown-linux-musl/release/preInit ./microfs/preInit
sudo cp -v ../../sysInit/target/x86_64-unknown-linux-musl/release/sysInit ./microfs/sbin/sysInit
sync

# create missing dir
sudo mkdir ./microfs/tmp ./microfs/var/log ./microfs/mnt ./microfs/sys

#Creating archive
cd ./microfs/
# Setting ownership in microfs
sudo chown -R root:root *
sync

echo "Creating cpio archive for test."
sudo find . | cpio --quiet -H newc -o | gzip -9 -n > ../microfs.img; 
sync

#Help to start qemu
echo "Use command: sudo qemu-system-x86_64 -m 512M -kernel ./virt-kernel  -initrd ./microfs.img -append \"console=ttyAMA0 console=tty0 console=ttyS0 rdinit=/preInit\" -serial stdio"
