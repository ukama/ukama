#!/bin/bash
set -euo pipefail

# CONFIGURATION
ARCH="aarch64"
ALPINE_VERSION="v3.19"
MINIROOTFS_VERSION="3.19.1"
BASE_URL="https://dl-cdn.alpinelinux.org/alpine/${ALPINE_VERSION}/releases/${ARCH}"
TARBALL="alpine-minirootfs-${MINIROOTFS_VERSION}-${ARCH}.tar.gz"
INITRAMFS_IMG="initramfs-${ARCH}.img.gz"
TMP_DIR="/tmp/alpine-initramfs"

# LOGGING
log() {
    local type="$1"; shift
    echo -e "\033[1;32m[$type]\033[0m $*"
}

error() {
    echo -e "\033[1;31m[ERROR]\033[0m $*" >&2
    exit 1
}

# CLEANUP
cleanup() {
    log "INFO" "Cleaning up temporary directory"
    sudo rm -rf "${TMP_DIR}"
}
trap cleanup EXIT

# MAIN
log "INFO" "Creating temp dir: ${TMP_DIR}"
mkdir -p "${TMP_DIR}"
cd "${TMP_DIR}"

log "INFO" "Downloading Alpine Minirootfs: ${TARBALL}"
wget -q "${BASE_URL}/${TARBALL}" || error "Failed to download minirootfs"

log "INFO" "Extracting minirootfs"
mkdir rootfs
sudo tar -xzf "${TARBALL}" -C rootfs

log "INFO" "Creating initramfs image: ${INITRAMFS_IMG}"
cd rootfs
sudo find . -print0 | sudo cpio --null -ov --format=newc | gzip -9 > "../${INITRAMFS_IMG}"
cd ..

log "SUCCESS" "Initramfs created: ${TMP_DIR}/${INITRAMFS_IMG}"
