#!/bin/sh
# Only meant for test purpose. This temporary test file will be replaced by ukamaOS virtual test suit which
# will mock devices.   

#Copying file
echo "Copying init files..!!"
sudo cp -v ../../init/target/x86_64-unknown-linux-musl/release/init ./microfs/
sudo cp -v ../../rcinit/target/x86_64-unknown-linux-musl/release/rcinit ./microfs/bin/rc.init
sudo cp -v ../../sysinit/target/x86_64-unknown-linux-musl/release/sysinit ./microfs/usr/bin/usysinit

#Sync files
sync

#Creating archive
cd ./microfs/
echo "Creating cpio archive for test."
sudo find . | cpio --quiet -H newc -o | gzip -9 -n > ../microfs.img; 

#Help to start qemu
echo "Use command: sudo qemu-system-x86_64 -m 512M -kernel ./virt-kernel  -initrd ./microfs.img -append \"console=ttyAMA0 console=tty0 console=ttyS0 rdinit=/init\" -serial stdio"

