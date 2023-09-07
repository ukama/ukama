#!/bin/sh

ARG=$0

preInit () {
echo "Staring custom preInit."

echo "Preparing tmpfs in /mnt/"

mount -t tmpfs rootfs /mnt

echo "Copying."
cp -ar / /mnt

sync;

mv /mnt/preInit.sh /mnt/Init.sh

sync

cd /mnt

echo "Setting on /mnt as /"
mount --move /mnt/ /

echo "Change root."
exec /bin/busybox chroot . /Init.sh

}

Init () {
echo "Starting Init."

echo "Starting dhcpcd"
/sbin/dhcpcd -b -f /etc/dhcpcd.conf

echo "Mount proc sys and dev file systems."
mount -t proc none /proc
mount -t sysfs none /sys
mount -t devtmpfs none /dev

echo "Enabling ip forwarding."
echo 1 > /proc/sys/net/ipv4/ip_forward

echo "For starting lxce:"
echo "  Run \`exec /boot/init\` on shell"

echo "Starting /bin/sh."
exec /bin/sh

}

echo "Preparing ukamaOS."
if [ ${ARG} == "/preInit.sh" ]; then
	echo "Stage 1 Init."
	preInit
else
	echo "Stage 2 Init."
	Init
fi
