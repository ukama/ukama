#!/bin/bash

# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2023-present, Ukama Inc.

# script to install ARM musl-based toolchain - use to build
# amplifier node image.

set -e
set -o pipefail

# Define variables for the toolchain
TOOLCHAIN_URL="https://musl.cc/arm-linux-musleabihf-cross.tgz"
TOOLCHAIN_DIR="/opt/arm-linux-musleabihf-cross"

# Function to check for sudo permissions
check_sudo() {
    if ! sudo -v; then
        echo "You do not have sudo privileges or sudo is not configured correctly."
        exit 1
    fi
}

# Function to determine which shell is being used and return the right rc file
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

# Step 0: Check if the user has sudo permissions before continuing
check_sudo

# Step 1: Determine the correct RC file (bashrc or zshrc)
RC_FILE=$(get_rc_file)

# Check if the toolchain is already installed
if [ -d "$TOOLCHAIN_DIR" ]; then
    echo "Toolchain already installed at $TOOLCHAIN_DIR"
else
    # Step 2: Download the ARM musl toolchain
    echo "Downloading ARM musl toolchain from $TOOLCHAIN_URL ..."
    wget $TOOLCHAIN_URL -O /tmp/arm-linux-musleabihf-cross.tgz

    # Step 3: Extract the toolchain
    echo "Extracting toolchain to $TOOLCHAIN_DIR ..."
    sudo mkdir -p $TOOLCHAIN_DIR
    sudo tar -xzf /tmp/arm-linux-musleabihf-cross.tgz -C /opt/

    # Step 4: Clean up downloaded file
    rm /tmp/arm-linux-musleabihf-cross.tgz
fi

# Step 5: Add the toolchain to PATH in the appropriate RC file
if ! grep -q "$TOOLCHAIN_DIR/bin" "$RC_FILE"; then
    echo "Adding toolchain to $RC_FILE ..."
    echo "export PATH=\"$TOOLCHAIN_DIR/bin:\$PATH\"" >> "$RC_FILE"
fi

# Reload the shell configuration
echo "Reloading $RC_FILE to apply the changes..."
source "$RC_FILE"

# Step 6: Verify the installation
echo "Verifying the installation ..."
if command -v arm-linux-musleabihf-gcc &> /dev/null; then
    echo "ARM musl toolchain successfully installed!"
    arm-linux-musleabihf-gcc --version
else
    echo "Installation failed: arm-linux-musleabihf-gcc not found in PATH."
    exit 1
fi

echo "Installation complete."
