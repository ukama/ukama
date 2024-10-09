#!/bin/bash

# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2023-present, Ukama Inc.

# script to install ARM musl-based or GNU-based toolchain -
# to build amplifier node image.

set -e  # Exit immediately if any command fails
set -o pipefail  # Exit if any part of a pipeline fails

# Define variables for the toolchains
MUSL_TOOLCHAIN_URL="https://musl.cc/arm-linux-musleabihf-cross.tgz"
MUSL_TOOLCHAIN_DIR="/opt/arm-linux-musleabihf-cross"

GNU_TOOLCHAIN_URL="https://releases.linaro.org/components/toolchain/binaries/latest-7/arm-linux-gnueabihf/gcc-linaro-7.5.0-2019.12-x86_64_arm-linux-gnueabihf.tar.xz"
GNU_TOOLCHAIN_DIR="/opt/gcc-linaro-7.5.0-2019.12-x86_64_arm-linux-gnueabihf"

check_sudo() {
    if ! sudo -v; then
        echo "You do not have sudo privileges or sudo is not configured correctly."
        exit 1
    fi
}

get_rc_file() {
    if [[ "$SHELL" == */zsh ]]; then
        echo "$HOME/.zshrc"
    elif [[ "$SHELL" == */bash ]]; then
        echo "$HOME/.bashrc"
    else
        echo "Unsupported shell ($SHELL). Please manually add the toolchain to your PATH."
        exit 1
    fi
}

install_musl_toolchain() {
    echo "Installing musl-based ARM toolchain..."

    # Check if the toolchain is already installed
    if [ -d "$MUSL_TOOLCHAIN_DIR" ]; then
        echo "Musl-based toolchain already installed at $MUSL_TOOLCHAIN_DIR"
    else
        # Download the ARM musl toolchain
        echo "Downloading ARM musl toolchain from $MUSL_TOOLCHAIN_URL ..."
        wget $MUSL_TOOLCHAIN_URL -O /tmp/arm-linux-musleabihf-cross.tgz

        # Extract the toolchain
        echo "Extracting toolchain to $MUSL_TOOLCHAIN_DIR ..."
        sudo mkdir -p $MUSL_TOOLCHAIN_DIR
        sudo tar -xzf /tmp/arm-linux-musleabihf-cross.tgz -C /opt/

        # Clean up downloaded file
        rm /tmp/arm-linux-musleabihf-cross.tgz
    fi
}

install_gnu_toolchain() {
    echo "Installing GNU-based ARM toolchain..."

    # Check if the toolchain is already installed
    if [ -d "$GNU_TOOLCHAIN_DIR" ]; then
        echo "GNU-based toolchain already installed at $GNU_TOOLCHAIN_DIR"
    else
        # Download the ARM GNU toolchain
        echo "Downloading ARM GNU toolchain from $GNU_TOOLCHAIN_URL ..."
        wget $GNU_TOOLCHAIN_URL -O /tmp/gcc-linaro-7.5.0-arm-linux-gnueabihf.tar.xz

        # Extract the toolchain
        echo "Extracting toolchain to $GNU_TOOLCHAIN_DIR ..."
        sudo mkdir -p $GNU_TOOLCHAIN_DIR
        sudo tar -xf /tmp/gcc-linaro-7.5.0-arm-linux-gnueabihf.tar.xz -C /opt/

        # Clean up downloaded file
        rm /tmp/gcc-linaro-7.5.0-arm-linux-gnueabihf.tar.xz
    fi
}

# Step 0: Check if the user has sudo permissions before continuing
check_sudo

# Step 1: Determine the correct RC file (bashrc or zshrc)
RC_FILE=$(get_rc_file)

# Step 2: Ask the user to choose between musl-based or GNU-based toolchain
echo "Select the ARM toolchain to install:"
echo "1) musl-based (arm-linux-musleabihf)"
echo "2) GNU-based (arm-linux-gnueabihf)"
read -p "Enter your choice [1 or 2]: " choice

if [ "$choice" == "1" ]; then
    install_musl_toolchain
    TOOLCHAIN_DIR=$MUSL_TOOLCHAIN_DIR
    COMPILER_BIN="arm-linux-musleabihf-gcc"
elif [ "$choice" == "2" ]; then
    install_gnu_toolchain
    TOOLCHAIN_DIR=$GNU_TOOLCHAIN_DIR
    COMPILER_BIN="arm-linux-gnueabihf-gcc"
else
    echo "Invalid choice. Exiting."
    exit 1
fi

# Step 3: Add the toolchain to PATH in the appropriate RC file
if ! grep -q "$TOOLCHAIN_DIR/bin" "$RC_FILE"; then
    echo "Adding toolchain to $RC_FILE ..."
    echo "export PATH=$TOOLCHAIN_DIR/bin:\$PATH" >> "$RC_FILE"
fi

# Reload the shell configuration
echo "Reloading $RC_FILE to apply the changes..."
source "$RC_FILE"

# Step 4: Verify the installation
echo "Verifying the installation ..."
if command -v $COMPILER_BIN &> /dev/null; then
    echo "ARM toolchain successfully installed!"
    $COMPILER_BIN --version
else
    echo "Installation failed: $COMPILER_BIN not found in PATH."
    exit 1
fi

echo "Installation complete."

