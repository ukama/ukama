#!/usr/bin/env bash

# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2023-present, Ukama Inc.

# Script to create ukama's virtual node locally
#
# Behavior:
# - --setup installs build tools + runtime libs via apt on Ubuntu
# - BOOTSTRAP_SERVER defaults to dev.bootstrap.ukama.com (env override allowed)
# - VNODE_METADATA is hard-coded (env override allowed)
# - VNODE_ID comes from:
#     1) --node-id <id> or --node-id=<id>
#     2) env VNODE_ID
#     3) default hard-coded value

# Example run:
# existing host
# $  ./mk_local_vnode.sh --node-id uk-sa2602-tnode-v0-344c
#
# for new OS/machine: require sudo
# $  ./mk_local_vnode.sh --setup --node-id uk-sa2602-tnode-v0-344c

set -euo pipefail

REPO_SERVER_URL="${REPO_SERVER_URL:-testing}"
REPO_NAME="${REPO_NAME:-virtualnode}"
REGISTRY="${REGISTRY:-localhost:5000}"
BOOTSTRAP_SERVER="${BOOTSTRAP_SERVER:-dev.bootstrap.ukama.com}"
DEFAULT_VNODE_ID="${DEFAULT_VNODE_ID:-uk-sa2601-tnode-v0-62f1}"

DO_SETUP=0

# Always overwrite VNODE_METADATA with known-good JSON
VNODE_METADATA=$(cat <<'JSON'
{
  "nodeInfo": {
    "type": "hnode",
    "partNumber": "",
    "skew": "",
    "mac": "",
    "swVersion": "",
    "mfgSwVersion": "",
    "assemblyDate": "2022-05-09T14:08:02.985079028-07:00",
    "oem": "",
    "mfgTestStatus": "pending",
    "status": "LabelGenerated"
  },
  "nodeConfig": [
    {
      "moduleID": "ukma-sa2219-trx-m0-e479",
      "type": "trx",
      "partNumber": "",
      "hwVersion": "",
      "mac": "",
      "swVersion": "",
      "mfgSwVersion": "",
      "mfgDate": "2022-05-09T14:08:02.985112609-07:00",
      "mfgName": "",
      "status": "AssemblyCompleted"
    }
  ]
}
JSON
)
export VNODE_METADATA

VNODE_ID=""

usage() {
  cat <<EOF
Usage:
  $0 [--setup] [--node-id <id>|--node-id=<id>]

Options:
  --setup              Install build tools + required runtime libs (Ubuntu apt).
  --node-id <id>       Virtual node ID (default: $DEFAULT_VNODE_ID)
EOF
}

log()  { echo "[$(date +'%H:%M:%S')] $*"; }
die()  { echo "ERROR: $*" >&2; exit 1; }
have() { command -v "$1" >/dev/null 2>&1; }
RUN()  { log "+ $*"; "$@"; }

while [[ $# -gt 0 ]]; do
    case "$1" in
        --setup)
            DO_SETUP=1
            shift
            ;;
        --node-id)
            [[ $# -ge 2 ]] || die "--node-id requires a value"
            VNODE_ID="$2"
            shift 2
            ;;
        --node-id=*)
            VNODE_ID="${1#*=}"
            shift
            ;;
        -h|--help)
            usage
            exit 0
            ;;
        *)
            die "Unknown argument: $1"
            ;;
    esac
done

: "${VNODE_ID:=$DEFAULT_VNODE_ID}"
export VNODE_ID BOOTSTRAP_SERVER

ubuntu_setup() {
    have sudo || die "sudo not found; can't run --setup"

    local BUILD_TOOLS=(
        git
        build-essential
        pkg-config
        autoconf
        automake
        libtool
        m4
        cmake
        meson
        ninja-build
        ca-certificates
        curl
    )

    local PODMAN_PKGS=(
        podman
        uidmap
        slirp4netns
        fuse-overlayfs
  )

    local RUNTIME_LIBS=(
        libc6
        libc-bin
        zlib1g
        libzstd1
        libnghttp2-14
        libpsl5
        libcurl3-gnutls
        libssh-4
        librtmp1
        libssl3
        libgnutls30
        libnettle8
        libhogweed6
        libp11-kit0
        libtasn1-6
        libidn2-0
        libunistring2
        libbrotli1
        libkrb5-3
        libgssapi-krb5-2
        libk5crypto3
        libkrb5support0
        libkeyutils1
        libldap-2.5-0
        liblber-2.5-0
        libsasl2-2
        libffi8
        libgmp10
        libcom-err2
        libsqlite3-0
        libmicrohttpd12
    )

    log "Running apt-get install for build tools, podman, and runtime libs..."
    RUN sudo apt-get update -y
    RUN sudo apt-get install -y "${BUILD_TOOLS[@]}" "${PODMAN_PKGS[@]}" "${RUNTIME_LIBS[@]}"

    # Ensure submodules are available (safe no-op if already initialized)
    if [[ -d .git ]]; then
        log "Initializing/updating git submodules..."
        RUN git submodule update --init --recursive
    else
        log "No .git directory here; skipping submodule init (run in repo root)."
    fi

    log "Setup complete."
}

if [[ "$DO_SETUP" -eq 1 ]]; then
    if [[ -f /etc/os-release ]] && grep -qi ubuntu /etc/os-release; then
        ubuntu_setup
    else
        die "--setup currently supports Ubuntu via apt-get only."
    fi
fi

### Main

VERSION="$(make -s print-version)"
[[ -n "$VERSION" ]] || die "make -s print-version returned empty VERSION"
have podman         || die "podman not found. Run with --setup or install podman."

LOCAL_IMAGE="localhost/${REPO_SERVER_URL}/${REPO_NAME}:${VERSION}"
REG_IMAGE="${REGISTRY}/${REPO_SERVER_URL}/${REPO_NAME}:${VERSION}"

REGISTRY_NAME="ukama_local_registry"

cleanup() {
    set +e
    podman rm -f "$REGISTRY_NAME" >/dev/null 2>&1 || true
}
trap cleanup EXIT

# Build the container image locally (podman build via Makefile)
RUN make clean
RUN make container

# Start (or restart) local registry on port 5000 (podman-only)
podman rm -f "$REGISTRY_NAME" >/dev/null 2>&1 || true
RUN podman run -d --name "$REGISTRY_NAME" --network host registry:latest

# Push the base image to the local registry
RUN podman tag "${LOCAL_IMAGE}" "${REG_IMAGE}"
RUN podman push --tls-verify=false "${REG_IMAGE}"

# Run incubator from the registry ref
RUN podman run --network host --privileged -it \
    -v "${PWD}/ukama_${VERSION}.tgz:/ukama/ukama.tgz:ro" \
    --tls-verify=false \
    -e VNODE_METADATA="$VNODE_METADATA" \
    -e VNODE_ID="$VNODE_ID" \
    -e VNODE_RUN_TARGET="local" \
    -e REPO_SERVER_URL="$REPO_SERVER_URL" \
    -e REPO_NAME="$REPO_NAME" \
    -e BOOTSTRAP_SERVER="$BOOTSTRAP_SERVER" \
    "${REG_IMAGE}"

# Pull output image from registry
OUT_REG_IMAGE="${REGISTRY}/${REPO_SERVER_URL}/${REPO_NAME}:${VNODE_ID}"
RUN podman pull --tls-verify=false "${OUT_REG_IMAGE}"

# Retag pulled image into local short name so inspect works
OUT_LOCAL_IMAGE="${REPO_SERVER_URL}/${REPO_NAME}:${VNODE_ID}"
RUN podman tag "${OUT_REG_IMAGE}" "${OUT_LOCAL_IMAGE}"

# cleanup
RUN make clean

echo "Done. Image available: ${OUT_LOCAL_IMAGE}"
