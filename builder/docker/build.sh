#!/bin/bash -x

# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2024-present, Ukama Inc.

REPO_DIR=/workspace/ukama
VENDOR_DIR=${REPO_DIR}/nodes/ukamaOS/distro/vendor

is_empty() {
    if [ -z "$(ls -A "$1")" ]; then
        return 0
    else
        return 1
    fi
}

# clone Ukama code and update the submodule
git clone http://github.com/ukama/ukama
cd ukama
git submodule init
git submodule update

if [ ! -d "$VENDOR_DIR" ]; then
    echo "The specified directory does not exist: ${VENDOR_DIR}"
    exit 1
fi

# Loop through first-level directories within the base directory
for dir in "$VENDOR_DIR"/*; do
    if [ -d "$dir" ]; then
        if is_empty "$dir"; then
            echo "Directory $dir is empty: ${dir}"
            exit 1
        fi
    fi
done

# install needed package for build
apt-get update && apt-get install -y \
    software-properties-common \
    && add-apt-repository universe \
    && apt-get update && apt-get install -y \
    build-essential \
    git \
    wget \
    autoconf \
    automake \
    libtool \
    pkg-config \
    libssl-dev \
    texinfo \
    cmake \
    tcl \
    zlib1g-dev \
    texlive \
    texlive-latex-extra \
    ghostscript \
    gperf \
    gtk-doc-tools \
    libev-dev

rm -rf /var/lib/apt/lists/*
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

# Build the builder
cd ${REPO_DIR}/builder && make

echo "All Done."

exit 0


