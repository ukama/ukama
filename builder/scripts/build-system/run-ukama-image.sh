#!/bin/bash

# This script mounts the rootfs from a partitioned disk image, creates a Docker image from it,
# and launches an interactive shell inside a container using that image.

set -euo pipefail

### === CONFIGURATION ===
IMG_FILE="ukama-com-image.img"
ROOTFS_PART_NUM=2
MOUNT_POINT="/mnt/ukama-com-rootfs"
TARBALL="/tmp/ukama-com-rootfs.tar.gz"
DOCKER_IMAGE_NAME="ukama-com-musl"
DOCKER_CONTAINER_NAME="ukama-com-musl"

if [ ! -f "$IMG_FILE" ]; then
    echo "[ERROR] Image file '$IMG_FILE' not found."
    exit 1
fi

if ! command -v losetup >/dev/null || ! command -v docker >/dev/null; then
    echo "[ERROR] Required tools 'losetup' and/or 'docker' are missing."
    exit 1
fi

### === CLEANUP EXISTING CONTAINER ===
if docker ps -a --format '{{.Names}}' | grep -q "^${DOCKER_CONTAINER_NAME}$"; then
    echo "[INFO] Removing existing container named ${DOCKER_CONTAINER_NAME}..."
    docker rm -f "$DOCKER_CONTAINER_NAME"
fi

### === CLEANUP EXISTING IMAGE ===
if docker images -q "$DOCKER_IMAGE_NAME" >/dev/null; then
    echo "[INFO] Removing existing Docker image: $DOCKER_IMAGE_NAME"
    docker rmi -f "$DOCKER_IMAGE_NAME"
fi

### === MOUNT IMAGE PARTITION ===
echo "[INFO] Attaching loop device for $IMG_FILE..."
LOOPDEV=$(sudo losetup --find --partscan --show "$IMG_FILE")
PART_DEV="${LOOPDEV}p${ROOTFS_PART_NUM}"

if [ ! -b "$PART_DEV" ]; then
    echo "[ERROR] Partition device $PART_DEV not found. Aborting."
    sudo losetup -d "$LOOPDEV"
    exit 1
fi

echo "[INFO] Mounting rootfs partition $PART_DEV..."
sudo mkdir -p "$MOUNT_POINT"
sudo mount "$PART_DEV" "$MOUNT_POINT"

### === CREATE TAR.GZ ROOTFS ===
echo "[INFO] Creating tarball from rootfs..."
sudo tar -czf "$TARBALL" -C "$MOUNT_POINT" .

### === CLEANUP MOUNT & LOOP ===
echo "[INFO] Cleaning up mount and loop device..."
sudo umount "$MOUNT_POINT"
sudo losetup -d "$LOOPDEV"

### === IMPORT ROOTFS INTO DOCKER ===
echo "[INFO] Importing rootfs into Docker as image: $DOCKER_IMAGE_NAME"
cat "$TARBALL" | docker import - "$DOCKER_IMAGE_NAME"

### === RUN DOCKER CONTAINER IN BACKGROUND ===
echo "[INFO] Starting container in background from image: $DOCKER_IMAGE_NAME"
docker run -it -d \
    --name "$DOCKER_CONTAINER_NAME" \
    "$DOCKER_IMAGE_NAME" /bin/sh

echo "[INFO] Container '$DOCKER_CONTAINER_NAME' is running."

echo
echo "✅ To connect:    docker exec -it $DOCKER_CONTAINER_NAME /bin/sh"
echo "✅ To attach:     docker attach $DOCKER_CONTAINER_NAME"
echo "✅ To stop:       docker stop $DOCKER_CONTAINER_NAME"
