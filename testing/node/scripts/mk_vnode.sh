#!/usr/bin/env bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2022-present, Ukama Inc.

set -euo pipefail

log() {
    local level="$1"; shift
    printf '[%s] %s\n' "$level" "$*"
}

die() {
    log "ERROR" "$*"
    exit 1
}

DEF_BUILD_DIR="./build"
DEF_TARGET="local"
TARGET="${DEF_TARGET}"

# "Namespace" (repo path) and image name. These are NOT registry hostnames.
IMAGE_NS="${REPO_SERVER_URL:-testing}"      # e.g. testing
IMAGE_NAME="${REPO_NAME:-virtualnode}"      # e.g. virtualnode

# Local registry host:port for pushes
LOCAL_REGISTRY="${LOCAL_REGISTRY:-localhost:5000}"

# Will be set by update_ukama_os_env
BUILD_ENV=""        # "local" or "container"
UKAMA_OS="${UKAMA_OS:-}"
NODED_ROOT=""       # ${UKAMA_OS}/distro/system/noded

# Optional extra args for schema/inventory generation (may be empty)
VNODE_SCHEMA_ARGS="${VNODE_SCHEMA_ARGS:-}"

CWD="$(pwd)"
export CWD

# Ensure build dir exists before realpath
mkdir -p "${DEF_BUILD_DIR}"
BUILD_DIR="$(realpath "${DEF_BUILD_DIR}")"

detect_env() {
    BUILD_ENV="local"

    # Common markers
    if [ -f "/.dockerenv" ] || [ -f "/run/.containerenv" ]; then
        BUILD_ENV="container"
        return
    fi

    # Some runtimes set env vars
    if [ -n "${container:-}" ] || [ -n "${CONTAINER:-}" ]; then
        BUILD_ENV="container"
        return
    fi

    # cgroup markers (best-effort)
    if grep -qaE '(docker|podman|containerd|kubepods|lxc)' /proc/1/cgroup 2>/dev/null; then
        BUILD_ENV="container"
        return
    fi
}

update_ukama_os_env() {
    detect_env

    # Allow explicit override via UKAMA_OS env var
    if [ -z "${UKAMA_OS:-}" ]; then
        if [ "${BUILD_ENV}" = "local" ]; then
            UKAMA_OS="$(realpath ../../nodes/ukamaOS)"
        else
            UKAMA_OS="/tmp/virtnode/ukama/nodes/ukamaOS"
        fi
    fi

    [ -d "${UKAMA_OS}" ] || die "Failed to find ukamaOS at: ${UKAMA_OS} (BUILD_ENV=${BUILD_ENV})"

    NODED_ROOT="${UKAMA_OS}/distro/system/noded"
    [ -d "${NODED_ROOT}" ] || die "Failed to find noded root at: ${NODED_ROOT}"

    log "INFO" "UKAMA_OS set to ${UKAMA_OS} (BUILD_ENV=${BUILD_ENV})"
}

build_utils() {
    mkdir -p "${BUILD_DIR}/utils"
    update_ukama_os_env

    log "INFO" "Building utils in ${NODED_ROOT}"

    pushd "${NODED_ROOT}" >/dev/null

    make genSchema
    [ -f "${NODED_ROOT}/build/genSchema" ] || die "Error building genSchema."
    cp -f "${NODED_ROOT}/build/genSchema" "${BUILD_DIR}/utils/"

    make genInventory
    [ -f "${NODED_ROOT}/build/genInventory" ] || die "Error building genInventory."
    cp -f "${NODED_ROOT}/build/genInventory" "${BUILD_DIR}/utils/"

    popd >/dev/null

    log "SUCCESS" "Utils built into ${BUILD_DIR}/utils"
}

build_sysfs() {
    local node_type="${1:-}"
    local node_uuid="${2:-}"

    [ -n "${node_type}" ] || die "build_sysfs requires NODE_TYPE as arg1"
    [ -n "${node_uuid}" ] || die "build_sysfs requires NODE_UUID as arg2"

    update_ukama_os_env

    log "INFO" "Preparing sysfs (type=${node_type}, uuid=${node_uuid})"

    "${NODED_ROOT}/utils/prepare_env.sh" --clean
    "${NODED_ROOT}/utils/prepare_env.sh" -u tnode -u anode

    # Copy schema + mfgdata locally
    mkdir -p "${BUILD_DIR}/schemas"
    cp -f "${NODED_ROOT}/mfgdata/schema/"*.json "${BUILD_DIR}/schemas/" 2>/dev/null || true
    cp -rf "${NODED_ROOT}/mfgdata" "${BUILD_DIR}/"

    pushd "${BUILD_DIR}" >/dev/null

    "${BUILD_DIR}/utils/genSchema" --u "${node_uuid}" \
                                   --n com --m UK-SA9001-COM-A1-1103  \
                                   --f mfgdata/schema/com.json --n trx \
                                   --m UK-SA9001-TRX-A1-1103  \
                                   --f mfgdata/schema/trx.json --n mask \
                                   --m UK-SA9001-MSK-A1-1103\
                                   --f mfgdata/schema/mask.json

    "${BUILD_DIR}/utils/genInventory" --n com --m UK-SA9001-COM-A1-1103 \
                                      --f mfgdata/schema/com.json -n trx \
                                      --m UK-SA9001-TRX-A1-1103 \
                                      --f mfgdata/schema/trx.json \
                                      --n mask -m UK-SA9001-MSK-A1-1103 \
                                      --f mfgdata/schema/mask.json

    popd >/dev/null

    # Copy sysfs to build dir
    [ -d /tmp/sys ] || die "/tmp/sys not found after prepare_env/genSchema steps"
    rm -rf "${BUILD_DIR}/sys"
    cp -rf /tmp/sys "${BUILD_DIR}/ukama/"
    mv "${BUILD_DIR}/ukama/sys" "${BUILD_DIR}/ukama/mocksysfs"
    rm -rf /tmp/sys

    log "SUCCESS" "Sysfs built at ${BUILD_DIR}/ukama/mocksysfs"
}

setup_ukama_dirs() {
    local nodeid="${1:-unknown}"
    local bootstrap_server="${2:-}"

    log "INFO" "Creating Ukama directories..."

    : "${BUILD_DIR:?BUILD_DIR not set}"
    : "${UKAMA_OS:?UKAMA_OS not set}"

    mkdir -p "${BUILD_DIR}/ukama"/{configs,apps/lib,apps/pkgs,apps/rootfs,apps/registry}

    # Metadata
    echo "${nodeid}" > "${BUILD_DIR}/ukama/nodeid"
    echo "${bootstrap_server}" > "${BUILD_DIR}/ukama/bootstrap"
    touch "${BUILD_DIR}/ukama/apps.log"

    # Copy all the apps configs.
    if [ -d "${UKAMA_OS}/../configs/apps" ]; then
        cp -r "${UKAMA_OS}/../configs/apps/." "${BUILD_DIR}/ukama/configs/"
    else
        log "WARN" "Apps config directory not found: ${UKAMA_OS}/../configs/apps"
    fi

    log "SUCCESS" "Ukama directories created at ${BUILD_DIR}/ukama"
}

build_image() {
    local file="${1:-}"
    local uuid="${2:-}"

    [ -n "${file}" ] || die "build_image requires ContainerFile path as arg1"
    [ -n "${uuid}" ] || die "build_image requires UUID as arg2"
    [ -f "${file}" ] || die "ContainerFile not found: ${file}"

    local name_tag
    name_tag="$(echo "${uuid}" | awk '{print tolower($0)}')"

    log "INFO" "Building image ${IMAGE_NS}/${IMAGE_NAME}:${name_tag}"

    # copy capp's sbin, conf and lib to /sbin, /conf and /lib
    mkdir -p "${BUILD_DIR}/"{sbin,lib,conf,tmp,bin}

    # Safer copy of capps content: avoid failing if glob doesn't match
    shopt -s nullglob
    for d in "${BUILD_DIR}"/capps/*; do
        [ -d "$d" ] || continue
        [ -d "$d/sbin" ] && cp -rf "$d/sbin" "${BUILD_DIR}/"
        [ -d "$d/conf" ] && cp -rf "$d/conf" "${BUILD_DIR}/"
        [ -d "$d/lib"  ] && cp -rf "$d/lib"  "${BUILD_DIR}/"
    done
    shopt -u nullglob

    cp -f ./scripts/runme.sh     "${BUILD_DIR}/bin/"
    cp -f ./scripts/waitfor.sh   "${BUILD_DIR}/bin/"
    cp -f ./scripts/kickstart.sh "${BUILD_DIR}/bin/"

    buildah bud -f "${file}" -t "${IMAGE_NS}/${IMAGE_NAME}:${name_tag}" .

    log "SUCCESS" "Buildah created image ${IMAGE_NS}/${IMAGE_NAME}:${name_tag}"
}

push_image_to_repo() {
    local uuid="${1:-}"
    local target="${2:-}"

    [ -n "${uuid}" ] || die "push_image_to_repo requires UUID as arg1"
    [ -n "${target}" ] || die "push_image_to_repo requires TARGET as arg2"

    local tag
    tag="$(echo "${uuid}" | awk '{print tolower($0)}')"

    if [ "${target}" != "remote" ]; then
        log "INFO" "Pushing to local registry ${LOCAL_REGISTRY}"
        buildah push --tls-verify=false \
                "${IMAGE_NS}/${IMAGE_NAME}:${tag}" \
                "${LOCAL_REGISTRY}/${IMAGE_NS}/${IMAGE_NAME}:${tag}"
        log "SUCCESS" "Image pushed to ${LOCAL_REGISTRY}/${IMAGE_NS}/${IMAGE_NAME}:${tag}"
        return
    fi

    # Remote push
    : "${REMOTE_REGISTRY:?REMOTE_REGISTRY must be set for remote push (e.g. ECR registry hostname)}"

    # If ECR, use AWS login if available
    if command -v aws >/dev/null 2>&1; then
        log "INFO" "Attempting AWS ECR login to ${REMOTE_REGISTRY}"
        local pass
        pass="$(aws ecr get-login-password)"
        buildah login --username "AWS" --password "${pass}" "${REMOTE_REGISTRY}"
    else
        : "${DOCKER_USER:?DOCKER_USER must be set for remote push if aws is not available}"
        : "${DOCKER_PASS:?DOCKER_PASS must be set for remote push if aws is not available}"
        log "INFO" "Logging into ${REMOTE_REGISTRY} as ${DOCKER_USER}"
        buildah login --username "${DOCKER_USER}" --password "${DOCKER_PASS}" "${REMOTE_REGISTRY}"
    fi

    log "INFO" "Pushing to remote registry ${REMOTE_REGISTRY}"
    buildah push \
            "${IMAGE_NS}/${IMAGE_NAME}:${tag}" \
            "${REMOTE_REGISTRY}/${IMAGE_NS}/${IMAGE_NAME}:${tag}"

    log "SUCCESS" "Image pushed to ${REMOTE_REGISTRY}/${IMAGE_NS}/${IMAGE_NAME}:${tag}"
}

# Main

ACTION="${1:-}"
shift || true

case "${ACTION}" in
    init)
        build_utils
        ;;
    sysfs)
        build_sysfs "${1:-}" "${2:-}"
        ;;
    ukamadirs)
        update_ukama_os_env
        setup_ukama_dirs "${1:-}" "${2:-}"
        ;;
    build)
        build_image "${1:-}" "${2:-}"
        ;;
    push)
        push_image_to_repo "${1:-}" "${2:-}"
        ;;
    cp)
        [ -n "${1:-}" ] || die "cp requires source as arg1"
        [ -n "${2:-}" ] || die "cp requires destination (relative under BUILD_DIR) as arg2"
        cp -- "${1}" "${BUILD_DIR}/${2}"
        ;;
    clean)
        update_ukama_os_env
        rm -f ContainerFile supervisor.conf
        [ -n "${1:-}" ] && buildah rmi -f "localhost/${1}" || true
        pushd "${NODED_ROOT}" >/dev/null
        make clean
        popd >/dev/null
        log "SUCCESS" "Clean complete"
        ;;
    *)
        echo "Usage:"
        echo "  $0 init"
        echo "  $0 sysfs <NODE_TYPE> <NODE_UUID>"
        echo "  $0 ukamadirs <NODE_ID> <BOOTSTRAP_SERVER>"
        echo "  $0 build <ContainerFile> <UUID>"
        echo "  $0 push <UUID> <remote|local>"
        echo "  $0 cp <src> <dest-relative-under-build-dir>"
        echo "  $0 clean [image_tag_or_name]"
        echo
        echo "Env vars (optional):"
        echo "  UKAMA_OS, VNODE_SCHEMA_ARGS, REPO_SERVER_URL, REPO_NAME, LOCAL_REGISTRY"
        echo "  For remote push: REMOTE_REGISTRY (and optionally DOCKER_USER/DOCKER_PASS)"
        exit 2
        ;;
esac

exit 0
