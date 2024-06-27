#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2024-present, Ukama Inc.

REPO_DIR=/workspace/ukama
PLATFORM_DIR=${REPO_DIR}/nodes/ukamaOS/distro/platform
VENDOR_DIR=${REPO_DIR}/nodes/ukamaOS/distro/vendor
VENDOR_MAKEFILE=${VENDOR_DIR}/Makefile
VENDOR_LIBRARIES=

set -e

is_empty() {
    # Check if the directory contains anything other than .git
    if [ -z "$(ls -A "$1" | grep -v '.git')" ]; then
        return 0
    else
        return 1
    fi
}

spinner() {

    local pid=$!
    local delay=0.1
    local spinstr='|/-\'

    while [ "$(ps a | awk '{print $1}' | grep -w $pid)" ]; do
        local temp=${spinstr#?}
        printf " [%c]  " "$spinstr"
        local spinstr=$temp${spinstr%"$temp"}
        sleep $delay
        printf "\b\b\b\b\b\b"
    done
    printf "    \b\b\b\b"
}

build_libraries() {

    VENDOR_LIBRARIES=$(grep -E "^LIST \+= " "$VENDOR_MAKEFILE" | awk '{print $3}')

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

create_ukama_build_file() {

    local original_file="ukama.json"
    local new_file="ukama_build.json"

    # Use sed to replace everything up to the last directory with /workspace
    sed -E 's|/[^/]+(/[^/]+)*|/workspace|g' "$original_file" > "$new_file"

    if [[ -f "$new_file" ]]; then
        echo "New file created: $new_file"
    else
        echo "Failed to create the new file."
        return 1
    fi
}

#default branch is main
BRANCH=${1:-main}

# clone Ukama code and update the submodule
git clone https://github.com/ukama/ukama
cd ukama
if [ "$BRANCH" != "main" ]; then
    git fetch
    git checkout "$BRANCH"
fi
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
export LD_LIBRARY_PATH=${VENDOR_DIR}/build/lib:${LD_LIBRARY_PATH}
export LD_LIBRARY_PATH=${VENDOR_DIR}/build/lib64:${LD_LIBRARY_PATH}
export LD_LIBRARY_PATH=${PLATFORM_DIR}/build:${LD_LIBRARY_PATH}
cd ${REPO_DIR}/builder && make

# Build the nodes
create_ukama_build_file
./builder nodes build --config-file ./ukama_build.json
rm ./ukama_build.json

# copy the created init, kernel and img file
cp ${REPO_DIR}/builder/script/vmlinuz* ${REPO_DIR}/..
cp ${REPO_DIR}/builder/script/initrd*  ${REPO_DIR}/..
cp ${REPO_DIR}/builder/script/*.img    ${REPO_DIR}/..

echo "All Done."

exit 0
