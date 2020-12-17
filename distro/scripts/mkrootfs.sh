#!/bin/sh
set -e
set +x
#**************************************
# UkamaOS make rootfs script.  
#**************************************

#DEFINE
CORES=$(getconf _NPROCESSORS_ONLN)
DISTRO_NAME="UkamaOS"
DISTRO_VERSION="0.0.1"
WD=${PWD}
ANODE=anode
CNODE=cnode
EXARGS=4
#DEFINE END

echo "Info:: Starting mkrootfs.sh in ${WD}."

#****** User should be root ***********
if [ $(id -u) -ne 0 ]; then
  echo "Err:: Run as root"; exit 1
fi
#**************************************

#****** Usage option ***********
usage ()
{
  echo 'Usage : mkrootfs -p <rootfs path>'
  exit;
}

msg_usage() {
        echo "Usage:"
        echo "          mkrootfs.sh [options]         "
        echo ""
        echo "Options:"
	echo "          -p string (/ukamafs)                     Path for rootfs created using make"
	echo "          -u string (anode|cnode)                  Unit type"
        echo "          -h                                       Display this help message."
	exit ;
}


#Argument check
if [ "$#" -ne ${EXARGS} ]
then
	echo "Err:: $# provided expected ${EXARGS} args (key value pair)" 
  	msg_usage
fi

while [ "$#" -gt 0 ]; do
    case $1 in
        -p|--path)
		ROOTFSPATH=${2}  
		echo "Info:: Ukama RootFS Path ${ROOTFSPATH}"
        	shift # Remove path from processing
		shift
        	;;
	-u|--unit)
                UNITTYPE=${2}  
                echo "Info:: Unit Type ${UNITTYPE}"
                shift # Remove path from processing
                shift
                ;;
        -h|--help)
        	echo "Help message:"
		msg_usage
        	shift
        	;;
        *)
        	echo "Err: Invalid args ${1}."
		msg_usage
        	shift # Remove generic argument from processing
        	;;
    esac
done

echo "Info:: Adding system configs and user config to ${DISTRO_NAME} rootfs."
if [ ! -d "${ROOTFSPATH}" ]
then
	echo "Err:: RootFS directory ${ROOTFSPATH} doesn't exist."
fi



#************ Adding rootfs dir structure **************
cd ${ROOTFSPATH}
echo "Info:: Adding directory struct to rootfs."
mkdir -p dev lib lib64 run mnt/root proc sys mmc history \
         tmp home var/log usr/share/udhcpc usr/local/bin \
         var/spool/cron/crontabs etc/init.d etc/rc.d var/run \
         var/www/html etc/network/if-down.d etc/network/if-post-down.d \
         etc/network/if-pre-up.d etc/network/if-up.d run \
         etc/cron/daily etc/cron/hourly etc/cron/monthly etc/cron/weekly
sync

#********* Device configuration***************
# config
HOST="ukama"
echo "Info:: Adding Devices to /dev"
mknod dev/console c 5 1
mknod dev/tty c 5 0
printf ${HOST} > etc/hostname
printf "root:x:0:0:root:/root:/bin/sh\nservice:x:1:1:service:/var/www/html:/usr/sbin/nologin" > etc/passwd
printf "ukama:x:1000:1000:Linux User,,,:/home/ukama:/bin/sh" >> etc/passwd
echo "root:mTteXHTdcIbEc:17743::::::" > etc/shadow
echo "ukama:mTteXHTdcIbEc:18585:0:99999:7:::" >> etc/shadow
echo "root:x:0:root\nservice:x:1:service" > etc/group
echo "ukama:x:1000:" >> etc/group
echo "/bin/sh" > etc/shells
echo "127.0.0.1	 localhost $host" > etc/hosts
echo "<html><h1>It Works!!</h1></html>" > var/www/html/index.html
#echo "UUID=$uuid  /  ext4  defaults,errors=remount-ro  0  1" > etc/fstab
echo "/dev/mmcblk0p1  /mmc  auto  errors=remount-ro  0  1" >> etc/fstab
#[ $swap_size -gt 0 ] && echo "UUID=$swap_uuid  none  swap  sw  0  0" >> etc/fstab

# path, prompt and aliases
echo "Info:: Adding profile"
cat << EOF > etc/profile
uname -snrvm
export PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin
PS1="\\u@\\h:\\w\\$ "
[ \$(id -u) -eq 0 ] && PS1="\\u@\\h:\\w# "
alias vim=vi
alias su="su -l"
alias locate=which
alias whereis=which
alias logout=exit
EOF

# banner
echo "Info:: Adding UkamaOS Banner"
printf "\n${DISTRO_NAME} Linux ${DISTRO_VERSION} \n" | tee -a etc/issue usr/share/infoban >/dev/null
cat << EOF >> etc/issue
|||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||
|||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||
||||||||  |||||  ||||  |||  |||||         ||||    ||||||    ||||         ||||         ||||         ||||||||||
||||||||  |||||  ||||  ||  ||||||  |||||  ||||  |  ||||  |  ||||  |||||  ||||  |||||  ||||  |||||  ||||||||||
||||||||  |||||  ||||  |  |||||||  |||||  ||||  ||  ||  ||  ||||  |||||  ||||  |||||  ||||  |||||||||||||||||
||||||||  |||||  ||||    ||||||||         ||||  |||    |||  ||||         ||||  |||||  ||||         ||||||||||
||||||||  |||||  ||||  |  |||||||  |||||  ||||  ||||||||||  ||||  |||||  ||||  |||||  |||||||||||  ||||||||||
||||||||  |||||  ||||  ||  ||||||  |||||  ||||  ||||||||||  ||||  |||||  ||||  |||||  ||||  |||||  ||||||||||
||||||||         ||||  |||  |||||  |||||  ||||  ||||||||||  ||||  |||||  ||||         ||||         ||||||||||
|||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||
|||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||
EOF
echo "cp /usr/share/infoban /etc/issue" > sbin/disban

# legal
echo "Info:: Adding copyright"
cat << EOF > etc/motd
The programs included with the $DISTRO_NAME Linux system are free software;
the exact distribution terms for each program are described in the
individual files in /usr/share/doc/*/copyright.
$DISTRO_NAME Linux comes with ABSOLUTELY NO WARRANTY, to the extent
permitted by applicable law.
EOF

# inittab
echo "Info:: Adding inittab"
cat << EOF > etc/inittab
ttyS0::respawn:/sbin/getty -L 115200 ttyS0
::sysinit:/sbin/swapon -a
::sysinit:/bin/hostname -F /etc/hostname
::sysinit:/etc/init.d/rcS
::ctrlaltdel:/sbin/reboot
::shutdown:/bin/echo SHUTTING DOWN
::shutdown:/sbin/swapoff -a
::shutdown:/etc/init.d/rcK
::shutdown:/bin/umount -a -r
EOF

echo "Info:: Adding networking interface"
# networking
cat << EOF > etc/network/interfaces
auto lo
iface lo inet loopback
auto eth0
iface eth0 inet static
	address 10.5.5.101
	netmask 255.255.255.0
EOF

# init
echo "Info:: Adding init script"
cat << EOF > init
#!/bin/busybox sh
/bin/busybox --install -s
export PATH=/bin:/sbin:/usr/bin:/usr/sbin
mountpoint -q proc || mount -t proc proc proc
mountpoint -q sys || mount -t sysfs sys sys
mknod /dev/null c 1 3
if ! mountpoint -q dev
then
  mount -t tmpfs -o size=64k,mode=0755 tmpfs dev
  mount -t tmpfs -o mode=1777 tmpfs tmp
  mkdir -p dev/pts
  mdev -s
fi
echo 0 > /proc/sys/kernel/printk
sleep 1

echo "Init:: UkamaOS init started."
sleep 1
echo "Init:: Install ethernet driver."
modprobe igb;
sleep 5;
echo "Try eth0 interface up."
ifconfig eth0 up;

mount -t tmpfs run /run -o mode=0755,nosuid,nodev
if [ ! -d /mnt/root/bin ] ; then
for i in bin etc lib root sbin usr home var boot share mmc history; do
  cp -r -p /\$i /mnt/root
done
mkdir /mnt/root/mnt
fi

for i in run tmp dev proc sys; do
  [ -d /mnt/root/\$i ] || mkdir /mnt/root/\$i
  mount -o bind /\$i /mnt/root/\$i
done

echo "chroot into /mnt/root"
mount -t devpts none /mnt/root/dev/pts
rm -r /bin /etc /sbin /usr
exec /mnt/root/bin/busybox chroot /mnt/root /sbin/init
EOF

# nologin
echo "Info:: Adding no-login message"
rm -rf usr/sbin/nologin
printf "#!/bin/sh
echo 'This account is currently not available.'
sleep 3
exit 1" > usr/sbin/nologin

# halt
#Make sure halt is not present
echo "Info:: Adding halt script"
rm -rf sbin/halt
cat << EOF > sbin/halt
#!/bin/sh
if [ \$1 ] && [ \$1 = '-p' ] ; then
    /bin/busybox poweroff
    return 0
fi
/bin/busybox halt
EOF

echo "Info:: Adding man pages"
# mini man pages
cat << EOF > sbin/man
#!/bin/sh
if [ -z "\$(busybox \$1 --help 2>&1 | head -1 | grep 'applet not found')" ]
then
  clear
  head="\$(echo \$1 | tr 'a-z' 'A-Z')(1)\\t\\t\\tManual page\\n"
  body="\$(busybox \$1 --help 2>&1 | tail -n +2)\\n\\n"
  printf "\$head\$body" | more
  exit 0
fi
echo "No manual entry for \$1"
EOF

#copy init-functions from LFS
cp -v ${WD}/distro/scripts/linux/init-functions etc/init.d/init-functions

# rcS & rcK
echo "Info:: Adding init and kill scripts"
#printf "#!/bin/sh
#. /etc/init.d/init-functions " > etc/init.d/rcS
#ln -s /etc/init.d/rcS etc/init.d/rcK

cat << EOF > etc/init.d/rcS
#!/bin/sh

# Start all init scripts in /etc/init.d
# executing them in numerical order.
#

# allow scripts to continue running after rcS completes
trap "" SIGHUP

echo "Init::rcS:: Checking keys for dropbear sshd daemon."
mkdir -p /etc/dropbear
if [ ! -f "/etc/dropbear/dropbear_dss_host_key" ] ; then
	dropbearkey -t dss -f /etc/dropbear/dropbear_dss_host_key
fi

if [ ! -f "/etc/dropbear/dropbear_rsa_host_key" ] ; then
	dropbearkey -t rsa -f /etc/dropbear/dropbear_rsa_host_key
fi

echo "Init:rcS:: Setting a temp sysfs directory"
if [ -d "/etc/ukamaEDR/sys" ]; then
	cp -rf /etc/ukamaEDR/sys /tmp/sys
fi

for i in /etc/rc.d/??*.sh ;do
     echo "Init::rcS:: looking for \$i ."
     # Ignore dangling symlinks (if any).
     [ ! -f "\$i" ] && continue

     case "\$i" in
	*.sh)
	    # Source shell script for speed.
	    (
		trap - INT QUIT TSTP
		set start
		. \$i
	    )
	    ;;
	*)
	    # No sh extension, so fork subprocess.
	    \$i start
	    ;;
    esac
done
EOF

# default crontabs
echo "Info:: Adding crontabs"
cat << EOF > var/spool/cron/crontabs/root
15  * * * *   cd / && run-parts /etc/cron/hourly
23  6 * * *   cd / && run-parts /etc/cron/daily
47  6 * * 0   cd / && run-parts /etc/cron/weekly
33  5 1 * *   cd / && run-parts /etc/cron/monthly
EOF

# logrotate
echo "Info:: Adding logrotate"
cat << EOF > etc/cron/daily/logrotate
#!/bin/sh
maxsize=512
dir=/var/log
for log in messages lastlog; do
  size=\$(du "\$dir/\$log" | tr -s '\t' ' ' | cut -d' ' -f1)
  if [ "\$size" -gt "\$maxsize" ] ; then
    tsp=\$(date +%s)
    mv "\$dir/\$log" "\$dir/\$log.\$tsp"
    touch "\$dir/\$log"
    gzip "\$dir/\$log.\$tsp"
  fi
done
EOF

# init scripts installer
echo "Info:: Adding add-rc.d"
cat << EOF > usr/bin/add-rc.d
#!/bin/sh
if [ -f /etc/init.d/\$1 ] && [ "\$2" -gt 0 ] ; then
ln -s /etc/init.d/\$1 /etc/rc.d/\$2\$1
echo "added \$1 to init."
else
echo "
  ** $distro_name add-rc.d ussage:
  add-rc.d [init.d script name] [order number]
  examples:
  add-rc.d httpd 40
  add-rc.d ftpd 40
  add-rc.d telnetd 50
"
fi
EOF

#Init scripts
echo "Adding init scripts."
INITDATA="
networking|network|30|/sbin/ifup|-a|/sbin/ifdown|-a
telnetd|telnet |80|/usr/sbin/telnetd|-p 23
cron|cron jobs|20|/usr/sbin/crond
syslogd|syslog|10|/sbin/syslogd
httpd|http server||/usr/sbin/httpd|-vvv -f -u service -h /var/www/html||httpd.log
ftpd|ftp services||/usr/bin/tcpsvd|-u service -vE 0.0.0.0 21 ftpd -S /var/www/html/
dropbear| ssh |50|/sbin/dropbear
dhcpcd| dhcpcd |40|/sbin/dhcpcd| -u service -b -f /etc/dhcpcd.conf ||dhcpcd.log"

OIFS=$IFS
IFS='
'
for i in $INITDATA; do
IFS='|'
set -- $i
echo "Adding $1 services."

if [ -z "$6" ] ; then
	KILLCOMMAND="pidof $1 | xargs kill"
else
	KILLCOMMAND="$6 $7"
fi

cat << EOF > etc/init.d/$1
#!/bin/sh
#
# Start the $1 service
#
echo "$1 service for $2"
case "\$1" in
  start)
    echo "Starting $1 service"
    $4 $5
    sleep 2
    ;;
  stop)
    echo -n "Stopping $1 service"
    $KILLCOMMAND
    ;;
  restart|reload)
    "\$0" stop
    "\$0" start
    ;;
  *)
    echo $"Usage: \$0 {start|stop|restart}"
    exit 1
esac

exit $?
EOF

chmod 744 etc/init.d/$1
[ $1 = 'telnetd' ] && [ "$telnetd" = false ] && continue;
[ "$3" ] && ln -s ../init.d/$1 etc/rc.d/$3$1.sh
done

#Adding daemon services for the Ukama DM
echo "Adding custom daemon services."

if [ ${UNITTYPE} = ${ANODE} ] ; then
	SYSTEMDBPATH=anode-systemdb
elif [ ${UNITTYPE} = ${CNODE} ] ; then
	SYSTEMDBPATH=cnode-systemdb
else
	echo "${UNITTYPE} is not added."
	exit 1;
fi

SERVDATA="
ukamaEDR|UBSP|90|/sbin/ukamaEDR| --p /etc/ukamaEDR/property.json --s /etc/ukamaEDR/sys/${SYSTEMDBPATH}
client|DMClient|91|/sbin/client| -4"

OIFS=$IFS
IFS='
'
for i in $SERVDATA; do
IFS='|'
set -- $i
echo "Adding $1 daemon services."

cat << EOF > etc/init.d/$1
#!/bin/sh
#
# Start the $1 service
#
PROG="$4"
ARGS="$5"
PID="/var/run/$1.pid"
LOG="/var/log/$1.log"
echo "$1 service for $2"
case "\$1" in
  start)
    echo "Starting $1 service"
    sleep 5
    start-stop-daemon --start -b --make-pidfile --pidfile=\$PID --startas /bin/sh -- -c "exec \$PROG \$ARGS >> \$LOG 2>&1"
    sleep 20
    ;;
  stop)
    echo -n "Stopping $1 service"
    start-stop-daemon --stop --exec \$PROG
    ;;
  restart|reload)
    "\$0" stop
    "\$0" start
    ;;
  *)
    echo $"Usage: \$0 {start|stop|restart}"
    exit 1
esac

exit $?
EOF

chmod 744 etc/init.d/$1
[ $1 = 'client' ] && continue;
[ "$3" ] && ln -s ../init.d/$1 etc/rc.d/$3$1.sh
#Client i.e Lwm2m2 service is only added not started on bootup.
done

# permissions
echo "Info:: Updating permissions"
touch proc/mounts var/log/wtmp var/log/lastlog
chmod 640  etc/shadow etc/inittab
chmod 664  var/log/lastlog var/log/wtmp
chmod 4755 bin/busybox
chmod 600  var/spool/cron/crontabs/root
chmod 755  usr/sbin/nologin sbin/disban init sbin/man etc/init.d/rcS\
           etc/cron/daily/logrotate usr/bin/add-rc.d sbin/halt
chmod 644  etc/passwd etc/group etc/hostname etc/shells etc/hosts etc/fstab\
           etc/issue etc/motd etc/network/interfaces etc/profile
sync

# Building initramfs
echo "Info:: Creating initrd.img"
find . | cpio --quiet -H newc -o | gzip -9 -n > ../${DISTRO_NAME}_initrd_${DISTRO_VERSION}.img
sync
echo "Info:: Initramfs image ${DISTRO_NAME}_initrd_${UNITTYPE}_${DISTRO_VERSION}.img is ready and available at ${WD}."
