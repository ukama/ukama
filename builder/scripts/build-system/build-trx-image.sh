#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2025-present, Ukama Inc.

set -e

STAGE="init"
DIR="$(pwd)"
UKAMA_OS=$(realpath ../../../nodes/ukamaOS)
UKAMA_ROOT=$(realpath ../../../)
UKAMA_STACK_REPO=""
MODE=""

trap cleanup EXIT

log() {
    local type="$1"
    local message="$2"
    local color
    case "$type" in
        "INFO") color="\033[1;34m";;
        "SUCCESS") color="\033[1;32m";;
        "ERROR") color="\033[1;31m";;
        *) color="\033[1;37m";;
    esac
    echo -e "${color}${type}: ${message}\033[0m"
}

check_status() {
    if [ $1 -ne 0 ]; then
        log "ERROR" "Script failed at stage: $3"
        exit 1
    fi
    log "SUCCESS" "$2"
}

cleanup() {
    cd "$UKAMA_STACK_REPO"
    make clean || true
    cd "$DIR"
    log "INFO" "Cleanup completed."
}

check_required_packages() {
    local packages=(
        bc libncurses5-dev lzop libssl-dev gnat flex zlib1g-dev gcc automake-1.15
        bison libelf-dev cmake curl libtool tcl pkg-config autopoint wget libisl-dev
        g++ texinfo texlive ghostscript libnftnl-dev libmnl-dev patchelf libsctp-dev
        openvpn libmicrohttpd-dev libcurl4-gnutls-dev libgnutls28-dev libgcrypt20-dev
        libsystemd-dev libcurl4 python3 bison cpp flex gettext cvs patch patchutils
        libncurses5-dev u-boot-tools python3-xcbgen automake mtd-utils gawk
        bsdmainutils rpm2cpio docbook docbook-utils sharutils zlib1g-dev libc6-i386
        libc6-dev-i386 python2
    )

    log "INFO" "Checking required packages..."
    local missing=()

    for pkg in "${packages[@]}"; do
        if ! dpkg -s "$pkg" &>/dev/null; then
            missing+=("$pkg")
        fi
    done

    if [ "${#missing[@]}" -ne 0 ]; then
        log "ERROR" "Missing packages: ${missing[*]}"
        echo "Run: sudo apt-get install ${missing[*]}"
        exit 1
    fi

    check_status 0 "All required packages are installed." "Package Check"
}

check_git_repo_and_submodules() {
    local repo_dir="$1"

    if [ ! -d "$repo_dir/.git" ]; then
        log "ERROR" "Directory '$repo_dir' is not a valid Git repository."
        exit 1
    fi

    log "INFO" "Found valid Git repository at '$repo_dir'."

    if [ -f "$repo_dir/.gitmodules" ]; then
        log "INFO" "Initializing and updating submodules..."
        (cd "$repo_dir" && git submodule init && git submodule update)
        check_status $? "Git submodules initialized and updated." "Git Submodule Update"
    else
        log "INFO" "No submodules found in '$repo_dir'."
    fi
}

if [ -z "$1" ]; then
    log "ERROR" "Usage: $0 <lte-stack-repo> [FDD|TDD]"
    exit 1
fi

UKAMA_STACK_REPO=$(realpath "$1")

MODE="${2:-FDD}"
if [[ "$MODE" != "FDD" && "$MODE" != "TDD" ]]; then
    log "ERROR" "Invalid mode '$MODE'. Expected 'FDD' or 'TDD'."
    exit 1
fi

check_required_packages
check_git_repo_and_submodules "$UKAMA_STACK_REPO"

# Build the toolchain and stack
TOOLS_DIR="$UKAMA_STACK_REPO/tools"
DISTRO_DIR="$UKAMA_STACK_REPO/distro"
BUILD_DIR="$TOOLS_DIR/build"
CROSSTOOL_DIR="$TOOLS_DIR/crosstool-ng"
HOST_BIN_DIR="$DISTRO_DIR/host/bin"

log "INFO" "Starting toolchain build..."

cd "$CROSSTOOL_DIR"
./bootstrap
check_status $? "Bootstrap completed." "crosstool-ng bootstrap"

mkdir -p "$BUILD_DIR"
./configure --prefix="$BUILD_DIR"
check_status $? "Configure completed." "crosstool-ng configure"

make
check_status $? "Build completed." "crosstool-ng make"

make install
check_status $? "Install completed." "crosstool-ng make install"

export PATH="$BUILD_DIR/bin:$PATH"

log "SUCCESS" "Toolchain built and installed."

# Setup build environment
cd "$DISTRO_DIR"
source ./env-setup OCTEON_CNF71XX_PASS1_1
check_status $? "Environment setup sourced." "env-setup"

export PATH="$HOST_BIN_DIR:$PATH"

# Ensure python2 symlink
if ! command -v python &> /dev/null; then
    sudo apt install -y python2
    sudo ln -sf /usr/bin/python2 /usr/bin/python
    log "INFO" "python2 installed and symlinked as python"
fi

# Clean up dtc parser (optional)
cd "$DISTRO_DIR"
if [ -f "dtc-lexer.l" ] && [ -f "dtc-parser.y" ]; then
    flex -o dtc-lexer.lex.c dtc-lexer.l
    bison -d -o dtc-parser.tab.c dtc-parser.y
    log "INFO" "DTC parser generated."
fi

# Build firmware and images
make fw
check_status $? "Firmware image built (fw.bin)." "make fw"

make kernel-deb
check_status $? "Kernel image built (lsm_os.gz)." "make kernel-deb"

make rootfs
check_status $? "Rootfs image built (lsm_rd.gz)." "make rootfs"

make kernel
check_status $? "Combined kernel+rootfs image built." "make kernel"

# build stack
cd "$UKAMA_STACK_REPO"
sudo make stack TYPE="$MODE"
check_status $? "Build stack image. mode: $MODE" "make stack"

log "SUCCESS" "TRX board images creation complete!"
