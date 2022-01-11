#!/bin/sh
# Copyright (c) 2021-present, Ukama Inc.
# All rights reserved.

# Script to generate rootfs for ukama capps.
# Used by the capp utility program

# Base parameters
UKAMA_OS=`realpath ../../../.`
SYS_ROOT=${UKAMA_OS}/distro/
SCRIPTS_ROOT=${SYS_ROOT}/scripts/
BB_ROOT=${UKAMA_OS}/distro/system/busybox
BB_CONFIG=ukama_minimal_defconfig

DEF_ROOTFS=_ukama_minimal_rootfs/

# For os-release
OS_NAME="ukamaOS"
OS_ID="ukama"
OS_VERSION_ID="0.0.1"
OS_PRETTY_NAME="UkamaOS V0.0.1"
OS_HOME_URL="http://www.ukama.com/"

#Various network related parameters
HOSTNAME="localhost"

# default target is local machine (gcc)
DEF_TARGET="local"
TARGET=${DEF_TARGET}

# default rootfs location is ${DEF_ROOTFS}
ROOTFS=${DEF_ROOTFS}

#
# Build busybox using the Ukama's minimal configuration
#
build_busybox() {
    CWD=`pwd`
    cd ${BB_ROOT}

    # set the config file for BB build
    BB_CONFIG=ukama_minimal_defconfig
    #Execute make and copy conent of _ukamafs to rootfs

    mkdir -p ${BB_ROOT}/${BB_ROOTFS}

    # setup proper compiler option.
    if [ "${TARGET}" != "local" ]
    then
	XGCC_PATH=${UKAMA_OS}/distro/tools/musl-cross-make/output/bin
    else
	XGCC_PATH=`which gcc | awk 'BEGIN{FS=OFS="/"}{NF--; print}'`
    fi

    make XGCCPATH=${XGCC_PATH}/ BBCONFIG=${BB_CONFIG} \
	 ROOTFSPATH=${BB_ROOT}/${BB_ROOTFS}

    if [ $? -ne 0 ]
    then
       exit 1
    fi

    cd ${CWD}
    cp -rf ${BB_ROOT}/${BB_ROOTFS}/* $ROOTFS

    # Go back and clean up
    cd ${BB_ROOT}
    make clean XGCCPATH=${XGCC_PATH}/
    cd ${CWD}
}

#
# Build the app at given src path and cmd
#
build_app() {

    SRC=$1
    CMD=$2
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

    DIRS="bin"
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

    cd ${ROOTFS}
    mkdir -p ${DIRS}
    build_etc_dirs
    build_usr_dirs

    cd ../
}

#
# setup /etc content
#
setup_etc() {

    cd ${ROOTFS}/etc/

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

    cp ${SCRIPTS_ROOT}/files/protocols ./protocols
    cp ${SCRIPTS_ROOT}/files/services  ./services

    cd ../../
}

# main

#remove existing copy of rootfs
if [ -d "${ROOTFS}" ]
then
    rm -rf ${ROOTFS}
fi

mkdir -p ${ROOTFS}
BB_ROOTFS=${ROOTFS}

# Action can be 'build', 'cp' and 'mkdir'
ACTION=$1

case "$ACTION" in
    #case 1
    "build")
	if [ "$2" = "app" ]
	then
	    build_app $3 $4
	elif [ "$2" = "busybox" ]
	then
	     build_busybox
	     build_rootfs_dirs
	     setup_etc
	fi
	;;
    "cp")
	cp $4 ${ROOTFS}/$4
	;;
    "mkdir")
	mkdir ${ROOTFS}/$4
	;;
esac

exit
