#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2024-present, Ukama Inc.

# Script to build and package ukamaOS app

set -e

MAJOR_VERSION=${1:-v3.17}
UKAMA_REPO=${2:-/ukamarepo}
export BUILD_MODE=release

# set up resolv.conf
echo "nameserver 8.8.8.8" > /etc/resolv.conf

# set Up Package Repositories
echo "https://dl-cdn.alpinelinux.org/alpine/${MAJOR_VERSION}/main" > /etc/apk/repositories
echo "https://dl-cdn.alpinelinux.org/alpine/${MAJOR_VERSION}/community" >> /etc/apk/repositories

# Update Package Manager
apk update
apk upgrade

# install essential system packages
apk add alpine-base openrc busybox bash sudo shadow tzdata

apk add acpid busybox-openrc busybox-extras busybox-mdev-openrc

apk add readline bash autoconf automake libmicrohttpd-dev gnutls-dev \
    openssl-dev iptables libuuid sqlite dhcpcd protobuf iproute2 zlib \
    curl-dev nettle libcap libidn2 libmicrohttpd gnutls openssl-dev \
    curl-dev linux-headers bsd-compat-headers tree libtool sqlite-dev \
    openssl-dev readline cmake autoconf automake alpine-sdk build-base \
    git tcpdump ethtool iperf3 htop vim doas \
    libunistring-dev \
    patchelf

# build apps
${UKAMA_REPO}/builder/scripts/build-all-apps.sh ${UKAMA_REPO}
if [ $? -eq 0 ]; then
    echo "Apps build:"
    ls -ltr "${UKAMA_REPO}/build/pkgs/"*

    echo "Package vendor libs and platform lib"
    mkdir -p "${UKAMA_REPO}/build/libs"

    VENDOR_LIB_DIR="${UKAMA_REPO}/nodes/ukamaOS/distro/vendor/build/lib"
    PLATFORM_LIB="${UKAMA_REPO}/nodes/ukamaOS/distro/platform/build/libusys.so"

    # Get list of *.a and *.so files (flat only)
    FILES=$(cd "$VENDOR_LIB_DIR" && ls *.a *.so* 2>/dev/null)

    tar -zcvf "${UKAMA_REPO}/build/libs/vendor_libs.tgz" \
        -C "$VENDOR_LIB_DIR" $FILES \
        -C "$(dirname "$PLATFORM_LIB")" "$(basename "$PLATFORM_LIB")"
else
    exit 1
fi

# Temporary - mocksysfs
cwd=$(pwd)
cd "${UKAMA_REPO}/nodes/ukamaOS/distro/system/noded"
rm -rf /tmp/sys/
rm -rf "${cwd}/mocksysfs"

# genSchema and genInventory are only run once and on host machine
# update the rpath so it can find the right libs.
UKAMAOS_ROOT="${UKAMA_REPO}/nodes/ukamaOS"
VENDOR_BUILD="${UKAMAOS_ROOT}/distro/vendor/build"
VENDOR_LIB="${VENDOR_BUILD}/lib"
VENDOR_LIB64="${VENDOR_BUILD}/lib64"
PLATFORM_LIB="${UKAMAOS_ROOT}/distro/platform/build"

RPATH_PATHS="${PLATFORM_LIB}:${VENDOR_LIB}:${VENDOR_LIB64}"

patchelf --set-rpath "${RPATH_PATHS}" "./build/genSchema"
patchelf --set-rpath "${RPATH_PATHS}" "./build/genInventory"

./utils/prepare_env.sh -u tnode -u anode
./build/genSchema --u UK-SA9001-HNODE-A1-1103 \
                  --n com --m UK-SA9001-COM-A1-1103  \
                  --f mfgdata/schema/com.json --n trx \
                  --m UK-SA9001-TRX-A1-1103  \
                  --f mfgdata/schema/trx.json --n mask \
                  --m UK-SA9001-MSK-A1-1103\
                  --f mfgdata/schema/mask.json
./build/genInventory --n com --m UK-SA9001-COM-A1-1103 \
                     --f mfgdata/schema/com.json -n trx \
                     --m UK-SA9001-TRX-A1-1103 \
                     --f mfgdata/schema/trx.json \
                     --n mask -m UK-SA9001-MSK-A1-1103 \
                     --f mfgdata/schema/mask.json
cp -rf /tmp/sys "${cwd}/mocksysfs"
cd "${cwd}"

exit 0
