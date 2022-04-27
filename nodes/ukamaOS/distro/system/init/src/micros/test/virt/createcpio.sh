#!/bin/sh
# Copyright (c) 2021-present, Ukama Inc.
# All rights reserved.

# Only meant for test purpose. This temporary test file will be replaced by
# ukamaOS virtual test suit, mocking devices.

UKAMA_OS=`realpath ../../../../../../../`
INIT_DIR=${UKAMA_OS}/distro/system/init
LXCE_DIR=${UKAMA_OS}/distro/system/lxce
SCRIPTS_DIR=${UKAMA_OS}/distro/scripts
PREINIT_REL_DIR=../../preInit/target/x86_64-unknown-linux-musl/release
SYSINIT_REL_DIR=../../sysInit/target/x86_64-unknown-linux-musl/release

#Run script with sudo
if ! [ $(id -u) = 0 ]; then
   echo "I am special. Run me as root!"
   exit 1
fi

#Cleaning files
echo "Cleaning microfs..."
rm -rf ./microfs/preInit ./microfs/init ./microfs/sbin/sysInit
rm -rf ./microfs/sbin/init ./microfs/conf ./microfs/lib
sync

# create missing dir
mkdir -p ./microfs/tmp ./microfs/var/log ./microfs/mnt ./microfs/sys
mkdir -p ./microfs/lib ./microfs/conf

#Copying file
echo "Copying init files..."
cp ${INIT_DIR}/init ./microfs/sbin/init
cp ${PREINIT_REL_DIR}/preInit ./microfs/preInit
cp ${SYSINIT_REL_DIR}/sysInit ./microfs/sbin/sysInit
cp ${LXCE_DIR}/lxce.d ./microfs/sbin/lxce.d
cp ./lxce/* ./microfs/conf/
cp ${SCRIPTS_DIR}/setup_space_network.sh ./microfs/sbin
sync

# copy all dependency lib
echo "Copying lib files..."
for lib in $(ldd ./microfs/sbin/lxce.d | cut -d '>' -f2 | awk '{print $1}')
do
    if [ -f "${lib}" ]; then
	sudo cp --parents "${lib}" "./microfs/"
	sudo cp "${lib}" "./microfs/lib/"
    fi
done

#Creating archive
cd ./microfs/
# Setting ownership in microfs
sudo chown -R root:root *
sync

echo "Creating cpio archive for test..."
sudo find . | cpio --quiet -H newc -o | gzip -9 -n > ../microfs.img; 
sync

echo "I'm done"

#Help to start qemu
echo "Now run following command to start QEMU:"
echo "sudo qemu-system-x86_64 -m 512M -kernel ./virt-kernel  -initrd ./microfs.img -append \"console=ttyAMA0 console=tty0 console=ttyS0 rdinit=/preInit\" -serial stdio"
