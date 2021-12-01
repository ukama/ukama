#!/bin/sh
# Copyright (c) 2021-present, Ukama Inc.
# All rights reserved.

# Script to generate minimal rootfs for Ukama Contained spaces and apps.

# Base parameters
UKAMA_OS=../../
SYS_ROOT=${UKAMA_ROOT}/distro/
BB_ROOT=${UKAMA_OS}/distro/system/busybox
BB_CONFIG=ukama_minimal_defconfig

# command line arguments
MIN_ARGS=2
DEF_ROOTFS=./_ukama_minimal_rootfs/

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

log_info() {
    echo "Info: $1"
}

log_error() {
    echo "${LIGHT_RED}Error:${NO_COLOR} $1"
}

usage() {
    echo 'Usage: mk_minimal_rootfs.sh -p <path_for_rootfs>'
    exit
}

msg_usage() {
    echo "Usage:"
    echo "      mk_minimal_rootfs.sh [options]"
    echo ""
    echo "Options:"
    echo "     -t target # Target is local(default), cnode, anode, etc."
    echo "     -p string # Path for minimal rootfs, e.g. ./_ukama_capp_rootfs/"
    echo "     -h        # Display this help message."
    exit
}

#
# Build busybox using the Ukama minimal configuration
#
build_busybox() {
    CWD=`pwd`
    cd ${BB_ROOT}/

    # set the config file for BB build
    BB_CONFIG=ukama_minimal_defconfig
    BB_ROOTFS=_ukamafs
    #Execute make and copy conent of _ukamafs to rootfs

    # setup proper compiler option.
    if [ "${TARGET}" != "local" ]
    then
	XGCC_PATH=${UKAMA_OS}/distro/tools/musl-cross-make/output/bin/
    else
	XGCC_PATH=`which gcc | awk 'BEGIN{FS=OFS="/"}{NF--; print}'`
    fi

    make XGCCPATH=${XGCC_PATH}/
    cd ${CWD}
    cp -rf ${BB_ROOT}/${BB_ROOTFS}/* $ROOTFS

    # Go back and clean up
    cd ${BB_ROOT}
    make clean XGCCPATH=${XGCC_PATH}/
    cd ${CWD}

    log_info "Busybox successfully build."
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
    
    log_info "Building rootfs directory structure at ${ROOTFS}"
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

    cp ../../files/protocols ./protocols
    cp ../../files/services  ./services

    cd ../../
}

# Script main.

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
                log_info "Ukama RootFS Path is: ${ROOTFS}"
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
            log_error "Invalid args: ${1}."
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
    log_info "Removing existing copy of ${ROOTFS}"
    rm -rf ${ROOTFS}
fi

mkdir -p ${ROOTFS}

log_info "Building busy box with Ukama minimal configuration"
build_busybox

log_info "Setting up rootfs directory structure"
build_rootfs_dirs

log_info "Setting up /etc contents under rootfs"
setup_etc

TOTAL_SIZE=`du -chs ${ROOTFS} | awk '{print $1}' | uniq`

log_info "All done."

log_info "Rootfs is located at: ${ROOTFS}"
log_info "Rootfs size is: ${TOTAL_SIZE}"

exit
