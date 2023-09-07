#!/bin/sh
PATH=/bin

echo -n "Mounting filesystems ... "
mount -t proc proc /proc
mount -t devtmpfs dev /dev
mount -t sysfs sys /sys
mount -t devpts pts /dev/pts
echo done

echo Hello, World
cd /root
./quicktest.sh
sh -i
