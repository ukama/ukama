#!/bin/bash

# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2021-present, Ukama Inc.

# Script to generate ukamaOS bootable image with minimal rootfs

# Build busybox
# Build starter.d
# Build dhcpcd
# Build sysctl
# Copy all lib dependencies
# create cpio arch
 set -x
# Base parameters
UKAMA_OS=`realpath ../../.`
UKAMA_REPO=`realpath ../../../../.`
VENDOR_ROOT=${UKAMA_OS}/distro/vendor
VENDOR_BUILD=${VENDOR_ROOT}/build/
SYS_ROOT=${UKAMA_OS}/distro/system
SCRIPTS_ROOT=${UKAMA_OS}/distro/scripts
APPS_BUILDER_ROOT=${UKAMA_REPO}/nodes/builder
BB_ROOT=${VENDOR_ROOT}/busybox
STARTERD_ROOT=${SYS_ROOT}/starterd
DHCPCD_SRC_ROOT=${VENDOR_ROOT}/dhcpcd
COMPILER_PATH=/usr/bin

# Build config parameters
#BB_CONFIG=ukama_minimal_defconfig
BB_CONFIG=defconfig
# command line arguments
MIN_ARGS=2
DEF_ROOTFS=_ukama_os_rootfs/
DEF_CSPACE_ROOTFS=cspace_rootfs

# For os-release
OS_NAME="ukamaOS"
OS_ID="ukama"
OS_VERSION_ID="0.0.1"
OS_PRETTY_NAME="UkamaOS V0.0.1"
OS_HOME_URL="http://www.ukama.com/"

#Various network related parameters
HOSTNAME="localhost"

LIGHT_RED='\e[31m'
NO_COLOR='\e[0m'

# default target is local machine (gcc)
DEF_TARGET="local"
TARGET=${DEF_TARGET}

# default rootfs location is ${DEF_ROOTFS}
ROOTFS=`realpath ${DEF_ROOTFS}`
CSPACE_ROOTFS=${DEF_CSPACE_ROOTFS}

log_info() {
    echo "Info: $1"
}

log_error() {
    echo "${LIGHT_RED}Error:${NO_COLOR} $1"
}

copy_file() {

    SRC=$1
    DST=$2

    if [ -d "${DST}" ]; then
	    if [ -f "${SRC}" ]; then
	        cp $SRC $DST
	    else
	        log_error "${SRC} not found. Exiting"
	        exit 1
	    fi
    else
	    log_error "${DST} not found. Exiting"
	    exit 1
    fi
}

usage() {
    echo 'Usage: buildOS.sh -t <build_target> -p <path_for_rootfs>'
    exit
}

msg_usage() {

    echo "Usage:"
    echo "      buildOS.sh [options]"
    echo ""
    echo "Options:"
    echo "     -t target      # Target is local(default), cnode, anode, etc."
    echo "     -p string      # Path for minimal rootfs, e.g. _ukama_os_rootfs"
    echo "     -h             # Display this help message."
    exit 0
}

#
# copy_all_libs
#
copy_all_libs() {

    ARGS=$1

    log_info "Copying lib for: ${ARGS}"
    
    for BIN in ${ARGS}
    do
	    for lib in $(ldd ${BIN} | cut -d '>' -f2 | awk '{print $1}')
	    do
	        if [ -f "${lib}" ]; then
		        cp --parents "${lib}" ${ROOTFS}
		        cp "${lib}" ${ROOTFS}/lib
	        fi
	        patchelf --set-rpath /lib ${BIN}
	    done
    done
}

#
# Build ip utilies (iptables and ip)
#
build_ip_utils() {
    CWD=`pwd`

    # setup proper compiler option
    if [ "${TARGET}" != "local" ]
    then
	    XGCC_PATH=${COMPILER_PATH}/
	    XTARGET=${TARGET}
    else
	    XGCC_PATH=`which gcc | awk 'BEGIN{FS=OFS="/"}{NF--; print}'`
	    XTARGET="linux"
    fi

    # build and copy iptables
    cd ${VENDOR_ROOT}
    make TARGET=${XTARGET} XGCCPATH=${XGCC_PATH}/ DEPDIR=${ROOTFS} iptables
    make TARGET=${XTARGET} XGCCPATH=${XGCC_PATH}/ DEPDIR=${ROOTFS} iproute2

    #copy_file ${VENDOR_BUILD}/sbin/iptables $ROOTFS/sbin/
    # remove the link to busybox
    #rm $ROOTFS/sbin/ip
    sync
    rm -rf ${ROOTFS}/*.Po 
    rm -rf ${ROOTFS}/*.Plo

    #copy_file ${VENDOR_BUILD}/sbin/ip $ROOTFS/sbin/
    sync
    cd ${CWD}
}

#
# dhcpcd build
#
build_dhcpcd() {
    CWD=`pwd`
    
    # setup proper compiler option
    if [ "${TARGET}" != "local" ]
    then
        XGCC_PATH=${COMPILER_PATH}
    else
        XGCC_PATH=`which gcc | awk 'BEGIN{FS=OFS="/"}{NF--; print}'`
    fi

    # build and copy init micros
    cd ${DHCPCD_SRC_ROOT}
    make TARGET=${TARGET} XGCCPATH=${XGCC_PATH} ROOTFSPATH=${ROOTFS}
    cd ${CWD}

    log_info "dhcpcd successfully build"
}


build_dhcp() {
    CWD=`pwd`
    
    # setup proper compiler option
    if [ "${TARGET}" != "local" ]
    then
        XGCC_PATH=${COMPILER_PATH}
    else
        XGCC_PATH=`which gcc | awk 'BEGIN{FS=OFS="/"}{NF--; print}'`
    fi

    # build dhcp
    cd ${VENDOR_ROOT}
    make TARGET=${TARGET} dhcp 
    cd ${CWD}

    log_info "dhcpcd successfully build"
}

#
# Build starter.d
#
build_starterd() {
   CWD=`pwd`

   # setup proper compiler option
   if [ "${TARGET}" != "local" ]
   then
       XGCC_PATH=${COMPILER_PATH}
   else
       XGCC_PATH=`which gcc | awk 'BEGIN{FS=OFS="/"}{NF--; print}'`
   fi

   cd ${STARTERD_ROOT}
   make clean; make TARGET=${TARGET} XGCCPATH=${XGCC_PATH}/
   copy_file ${STARTERD_ROOT}/starter.d ${ROOTFS}/sbin

   # copy manifest file
   copy_file ${UKAMA_REPO}/nodes/manifest.json ${ROOTFS}

   # Go back and clean up
   cd ${STARTERD_ROOT}; make clean
   cd ${CWD}

   log_info "starter.d successfully build"
}

#
# Build busybox using the Ukama minimal configuration
#
build_busybox() {

    CWD=`pwd`
    cd ${BB_ROOT}

    # set the config file for BB build
    BB_CONFIG=defconfig
    #Execute make and copy conent of _ukamafs to rootfs

    mkdir -p ${BB_ROOTFS}

    # setup proper compiler option
    if [ "${TARGET}" != "local" ]
    then
	    XGCC_PATH=${COMPILER_PATH}
    else
	    XGCC_PATH=`which gcc | awk 'BEGIN{FS=OFS="/"}{NF--; print}'`
    fi

    make TARGET=${TARGET} XGCCPATH=${XGCC_PATH}/ BBCONFIG=${BB_CONFIG} \
	     ROOTFSPATH=${BB_ROOTFS}

    if [ $? -ne 0 ]
    then
        log_error "Busybox compliation failed"
        exit 1
    fi

    cd ${CWD}

    # Go back and clean up
    cd ${BB_ROOT}
    make clean XGCCPATH=${XGCC_PATH}/
    cd ${CWD}

    log_info "Busybox successfully build"
}

#
# Build apps and copy them to rootfs
#
build_apps() {

    # Steps are:
    # 1. Build builder
    # 2. Build capp pkgs using the builder
    # 3. Create /ukama/apps onto rootfs (pkgs, store, registry, rootfs)
    # 4. Copy pkgs

    CWD=`pwd`

    cd ${APPS_BUILDER_ROOT}

    # make TARGET=${TARGET} XGCCPATH=${XGCCPATH}/

    # if [ -d ${APPS_BUILDER_ROOT}/pkgs/ ]
    # then
	#     rm -rf ${APPS_BUILDER_ROOT}/pkgs/
    # fi

    # Compile the builder
    make clean; make

    # for each apps in systems
    apps=("bootstrap" "configd" "deviced" "lookoutd" "meshd" "metricsd" "noded" "notifyd" "rlog" "started" "wimcd") 
    for app in ${apps[@]}; do
        basename=$(basename "$app")
        #./builder --create --config ./configs/${basename}.toml
        make -C ${SYS_ROOT}/${basename} install TARGET=${TARGET} INSTALL_DIR=${ROOTFS} 
    done

    # create apps dir onto rootfs
    #DIRS="${ROOTFS}/ukama/apps/pkgs"
    #DIRS="${ROOTFS}/ukama/apps/registry ${DIRS}"
    #DIRS="${ROOTFS}/ukama/apps/rootfs   ${DIRS}"
    #mkdir -p ${DIRS}

    log_info "apps succesfully build"
}

copy_vendor_libs() {
    cd ${VENDOR_ROOT}
    cp -vrf ${VENDOR_ROOT}/build/* ${ROOTFS}
    sync 
    log_info "libs copied succesfully"
}

#
# Build the usr directory structure
#
build_usr_dirs() {

    DIRS="bin"
    DIRS="lib $DIRS"
    DIRS="local/bin   ${DIRS}"
    DIRS="local/lib   ${DIRS}"
    DIRS="local/share ${DIRS}"
    DIRS="sbin        ${DIRS}"
    DIRS="share/misc  ${DIRS}"

    cd ./usr
    mkdir -p ${DIRS}
    cd ../
}

#
# Build the etc directory structure
#
build_etc_dirs() {

    DIRS="network/if-down.d"
    DIRS="network/if-post-down.d ${DIRS}"
    DIRS="network/if-post-up.d   ${DIRS}"
    DIRS="network/if-pre-down.d  ${DIRS}"
    DIRS="network/if-pre-up.d    ${DIRS}"
    DIRS="network/if-up.d        ${DIRS}"

    cd ./etc
    mkdir -p ${DIRS}
    cd ../
}

#
# Build rootfs directory structure
#
build_rootfs_dirs() {

    DIRS="ukama"
    DIRS="boot ${DIRS}"
    DIRS="bin  ${DIRS}"
    DIRS="sbin ${DIRS}"
    DIRS="etc  ${DIRS}"
    DIRS="home ${DIRS}"
    DIRS="lib  ${DIRS}"
    DIRS="mnt  ${DIRS}"
    DIRS="tmp  ${DIRS}"
    DIRS="sys  ${DIRS}"
    DIRS="usr  ${DIRS}"
    DIRS="var  ${DIRS}"
    DIRS="dev  ${DIRS}"
    DIRS="conf ${DIRS}"
    DIRS="mnt  ${DIRS}"
    DIRS="proc ${DIRS}"
    DIRS="run  ${DIRS}"
    DIRS="var/log ${DIRS}"
    DIRS="var/run/netns ${DIRS}"

    cd ${ROOTFS}
    mkdir -p ${DIRS}

    mkdir -p /var/run/netns
    build_etc_dirs
    build_usr_dirs

    touch proc/mounts var/log/wtmp var/log/lastlog
    sync

    cd ../

    log_info "Building rootfs directory structure at: ${ROOTFS}"
}

#
# setup /etc content
#
setup_etc() {

    cd ${ROOTFS}/etc/

    printf "/dev/mmcblk0p1  /mmc  auto  errors=remount-ro  0  1" >> ./fstab

    printf "${OS_VERSION_ID}\n" > ./ukama-release
    printf "${HOSTNAME}\n"   > ./hostname

    printf "127.0.0.1  ${HOSTNAME} ${HOSTNAME}.localdomain\n" > ./hosts
    printf "::1        ${HOSTNAME} ${HOSTNAME}.localdomain\n" >> ./hosts

    printf "Welcome to ${OS_NAME}\n" > ./issue

    cat << EOF > ./motd
The programs included with the ${OS_NAME} Linux system are free software;
the exact distribution terms for each program are described in the
individual files in /usr/share/doc/*/copyright.
${OS_NAME} Linux comes with ABSOLUTELY NO WARRANTY, to the extent
permitted by applicable law.
EOF

    cat <<EOF > ./network/if-up.d/dad
#!/bin/sh

# Block ifup until DAD completion
# Copyright (c) 2016-2018 Kaarle Ritvanen

has_flag() {
	ip address show dev $IFACE | grep -q " $1 "
}

while has_flag tentative && ! has_flag dadfailed; do
	sleep 0.2
done
EOF

    cat <<EOF > ./os-release 
NAME=${OS_NAME}
ID=${OS_ID}
VERSION_ID=${OS_VERSION_ID}
PRETTY_NAME=${OS_PRETTY_NAME}
HOME_URL=${OS_HOME_URL}
EOF

    printf "root:x:0:0:root:/root:/bin/sh\nservice:x:1:1:service:/var/www/html:/usr/sbin/nologin \n" > ./passwd
    printf "ukama:x:1000:1000:Linux User,,,:/home/ukama:/bin/sh \n" >> ./passwd

    printf "root:mTteXHTdcIbEc:17743::::::\n" > ./shadow
    printf "ukama:mTteXHTdcIbEc:18585:0:99999:7:::\n" >> ./shadow

    printf "root:x:0:root\nservice:x:1:service\n" > ./group
    printf "ukama:x:1000:\n" >> ./group

    printf "/bin/sh\n" > ./shells

    cat <<EOF > ./securetty
console
tty1
tty2
tty3
tty4
tty5
tty6
tty7
tty8
tty9
tty10
tty11
EOF

    cp ../../files/protocols ./protocols
    cp ../../files/services  ./services

    cd ../../

    sync
}

#
# setup_device
#
setup_device() {

    CWD=`pwd`

    cd ${ROOTFS}
    sudo mknod ./dev/console c 5 1
    sudo mknod ./dev/tty c 5 0
    sync

    cd ${CWD}
}

#
# Script main
#
WD=`pwd`
while [ "$#" -gt 0 ]; do
    case $1 in
        -p|--path)
	        if [ -z "$2" ]
	        then
		        log_info "Missing rootfs parameter for -p"
		        log_info "Setting to default: ${DEF_ROOTFS}"
		        ROOTFS=${DEF_ROOTFS}
	        else
                ROOTFS=$2
                log_info "ukamaOS RootFS Path is: ${ROOTFS}"
                shift # Remove path from processing
	        fi
            shift
            ;;

        -h|--help)
            echo "Help message"
            msg_usage
            shift
            ;;

	    -t|--target)
	        if [ -z "$2" ]
	        then
		        log_info "Missing target parameter for -t"
		        log_info "Setting to default: ${DEF_TARGET}"
		        TARGET=${DEF_TARGET}
	        else
		        TARGET=$2
		        log_info "Target is: ${TARGET}"
		        shift
	        fi
	        shift
	        ;;

        *)
            log_error "Invalid args: ${1}"
            msg_usage
            shift # Remove generic argument from processing
            ;;
    esac
done

if [ -z ${ROOTFS} ]
then
    log_info "-p not defined. Setting to default: ${DEF_ROOTFS}"
    ROOTFS=${DEF_ROOTFS}
fi

#setup rootfs location
if [ -d "${ROOTFS}" ]
then
    log_error "Please remove existing copy of ${ROOTFS}"
    #exit
fi

mkdir -p ${ROOTFS}
BB_ROOTFS=${ROOTFS}

log_info "Setting up rootfs directory structure"
build_rootfs_dirs

log_info "Building busy box with Ukama minimal configuration"
build_busybox
sleep 2

log_info "Building dhcp"
build_dhcp
sleep 2

log_info "Building ip utils"
#build_ip_utils

sleep 2
log_info "Building apps"
build_apps

sleep 2
log_info "Setting up /etc contents under rootfs"
setup_etc

sleep 2 
log_info "Setting up /dev"
setup_device

log_info "Copying vendor libs"
copy_vendor_libs

# Change ownership and create archieve
log_info "Changing ownership, updating permission and creating cpio archive"
cd ${ROOTFS}

chmod 640  etc/shadow
chmod 664  var/log/lastlog var/log/wtmp
chmod 4755 bin/busybox
chmod 755  usr/sbin/nologin
chmod 644  etc/passwd etc/group etc/hostname etc/shells etc/hosts etc/fstab \
      etc/issue etc/motd

sudo chown root:root .
sudo chown -R root:root *

sync

# Building initramfs
log_info "Creating initrd.img"
IMG=${OS_NAME}_initrd_${TARGET}_${OS_VERSION_ID}
sudo find . | cpio --quiet -H newc -o | gzip -9 -n > ${WD}/${IMG}.img
sync
cd ${WD}

TOTAL_ROOTFS_SIZE=`du -chs ${ROOTFS} | awk '{print $1}' | uniq`
IMAGE_SIZE=`du -kh ${WD}/${IMG}.img | cut -f1`
IMAGE_LOC=`realpath ${WD}/${IMG}.img`

log_info "All done. Have fun!"
log_info "------------------"
log_info "  Rootfs loc:   ${ROOTFS}"
log_info "  Rootfs size:  ${TOTAL_ROOTFS_SIZE}"
log_info "  ukamaOS loc:  ${IMAGE_LOC}"
log_info "  ukamaOS size: ${IMAGE_SIZE}"

exit
