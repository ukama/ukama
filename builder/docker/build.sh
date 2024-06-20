#!/bin/bash -x

# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2024-present, Ukama Inc.

REPO_DIR=/workspace/ukama
VENDOR_DIR=${REPO_DIR}/nodes/ukamaOS/distro/vendor
VENDOR_MAKEFILE=${VENDOR_DIR}/Makefile
VENDOR_LIBRARIES=$(grep -E "^LIST \+= " "$VENDOR_MAKEFILE" | awk '{print $3}')

is_empty() {
    # Check if the directory contains anything other than .git
    if [ -z "$(ls -A "$1" | grep -v '.git')" ]; then
        return 0
    else
        return 1
    fi
}

build_libraries() {

    BUILD_DIR=${VENDOR_DIR}/build
    for lib in $VENDOR_LIBRARIES; do
        if [ "$lib" == "zlib" ]; then
            LIB_PATH="${BUILD_DIR}/lib/libz.a"
        elif [ "$lib" == "libmicrohttpd" ]; then
            LIB_PATH="${BUILD_DIR}/lib/libmicrohttpd.a"
        elif [ "$lib" == "libuuid" ]; then
            LIB_PATH="${BUILD_DIR}/lib/libuuid.a"
        elif [ "$lib" == "openssl" ]; then
            LIB_PATH="${BUILD_DIR}/lib/libssl.a"
        else
            LIB_PATH="${BUILD_DIR}/lib/lib${lib}.a"
        fi

        if [ ! -f "$LIB_PATH" ]; then
            make $lib
        fi
    done
}

# clone Ukama code and update the submodule
git clone https://github.com/ukama/ukama
cd ukama
git fetch
git checkout build-update
git submodule init
git submodule update

if [ ! -d "$VENDOR_DIR" ]; then
    echo "The specified directory does not exist: ${VENDOR_DIR}"
    exit 1
fi

# Loop through first-level directories within submodules
for dir in "$VENDOR_DIR"/*; do
    if [ -d "$dir" ]; then
        if is_empty "$dir"; then
            echo "Directory $dir is empty: ${dir}"
            exit 1
        fi
    fi
done

cd /workspace

# Install some additional packages needed for building vendor
wget http://ftp.gnu.org/gnu/gettext/gettext-0.21.tar.gz \
    && tar -xvf gettext-0.21.tar.gz \
    && cd gettext-0.21 \
    && ./configure \
    && make \
    && make install \
    && cd .. \
    && rm -rf gettext-0.21.tar.gz \
    && rm -rf gettext-0.21

wget http://ftpmirror.gnu.org/libtool/libtool-2.4.7.tar.gz \
    && tar -xzf libtool-2.4.7.tar.gz \
    && cd libtool-2.4.7 \
    &&  ./configure \
    && make \
    && make install \
    && cd .. \
    && rm -rf libtool-2.4.7.tar.gz \
    && rm -rf libtool-2.4.7

# Build the vendor libraries
cd ${VENDOR_DIR} && build_libraries

# Build the builder
cd ${REPO_DIR}/builder && make

echo "All Done."

exit 0
