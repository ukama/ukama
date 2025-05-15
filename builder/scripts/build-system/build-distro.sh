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

# set up resolv.conf
echo "nameserver 8.8.8.8" > /etc/resolv.conf

# set Up Package Repositories
echo "https://dl-cdn.alpinelinux.org/alpine/${MAJOR_VERSION}/main" > /etc/apk/repositories
echo "https://dl-cdn.alpinelinux.org/alpine/${MAJOR_VERSION}/community" >> /etc/apk/repositories

# Update Package Manager
apk update
apk upgrade

# ===== Install Essential System Packages =====
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

#build apps
${UKAMA_REPO}/builder/scripts/build-all-apps.sh ${UKAMA_REPO}
if [ $? -eq 0 ]; then
    echo "Apps build:"
    ls -ltr "${UKAMA_REPO}/build/pkgs/"*

    echo "Package vendor libs"
    cd "${UKAMA_REPO}/nodes/ukamaOS/distro/vendor" || exit 1
    ls -ltr build/lib/*
    mkdir -p "${UKAMA_REPO}/build/libs"
    tar -zcvf "${UKAMA_REPO}/build/libs/vendor_libs.tgz" build/lib/*
else
    exit 1
fi
