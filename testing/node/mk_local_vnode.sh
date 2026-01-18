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
VNODE_ID=""

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

usage() {
    cat <<EOF
Usage:
    $0 [--setup] [--node-id <id>|--node-id=<id>]

Options:
    --setup              Install build tools + runtime libs (Ubuntu).
    --node-id <id>       Virtual node ID (default: $DEFAULT_VNODE_ID)
EOF
}

log() {
    echo "[$(date +'%H:%M:%S')] $*"
}

die() {
    echo "ERROR: $*" >&2
    exit 1
}

have() {
    command -v "$1" >/dev/null 2>&1
}

RUN() {
    log "+ $*"
    "$@"
}

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

apt_has() {
    apt-cache show "$1" >/dev/null 2>&1
}

install_candidates() {
    local label="$1"
    shift

    local pkg
    for pkg in "$@"; do
        if apt_has "$pkg"; then
            log "Using $label: $pkg"
            APT_PKGS+=("$pkg")
            return 0
        fi
    done

    log "WARN: no candidate found for $label (tried: $*)"
}

ubuntu_setup() {
    have sudo || die "sudo not found"

    local codename="unknown"
    if [[ -f /etc/os-release ]]; then
        . /etc/os-release
        codename="${VERSION_CODENAME:-unknown}"
    fi

    log "Detected Ubuntu codename: $codename"

    RUN sudo apt-get update -y

    local BASE_TOOLS=(
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

    APT_PKGS=()
    APT_PKGS+=("${BASE_TOOLS[@]}" "${PODMAN_PKGS[@]}")
    APT_PKGS+=(libc6 libc-bin)

    install_candidates "zlib"    zlib1g
    install_candidates "zstd"    libzstd1
    install_candidates "nghttp2" libnghttp2-14
    install_candidates "psl"     libpsl5t64 libpsl5
    install_candidates "brotli"  libbrotli1

    install_candidates "curl-gnutls" libcurl3t64-gnutls libcurl3-gnutls
    install_candidates "libssh"      libssh-4
    install_candidates "librtmp"     librtmp1

    install_candidates "openssl"   libssl3t64 libssl3
    install_candidates "gnutls"    libgnutls30t64 libgnutls30
    install_candidates "nettle"    libnettle8t64 libnettle8
    install_candidates "hogweed"   libhogweed6t64 libhogweed6
    install_candidates "p11-kit"   libp11-kit0
    install_candidates "tasn1"     libtasn1-6
    install_candidates "idn2"      libidn2-0
    install_candidates "unistring" libunistring5 libunistring2
    install_candidates "ffi"       libffi8

    install_candidates "krb5"        libkrb5-3
    install_candidates "gssapi"      libgssapi-krb5-2
    install_candidates "k5crypto"    libk5crypto3
    install_candidates "krb5support" libkrb5support0
    install_candidates "keyutils"    libkeyutils1
    install_candidates "com-err"     libcom-err2

    install_candidates "ldap" libldap2 libldap-2.5-0
    install_candidates "lber" liblber2 liblber-2.5-0
    install_candidates "sasl" libsasl2-2

    install_candidates "sqlite"     libsqlite3-0
    install_candidates "microhttpd" libmicrohttpd12t64 libmicrohttpd12
    install_candidates "gmp"        libgmp10
    
    # ---- Build-time headers for Ulfius / GnuTLS stack ----
    install_candidates "gnutls dev headers"     libgnutls28-dev
    install_candidates "nettle dev headers"     nettle-dev
    install_candidates "p11-kit dev headers"    libp11-kit-dev
    install_candidates "tasn1 dev headers"      libtasn1-6-dev
    install_candidates "idn2 dev headers"       libidn2-0-dev
    install_candidates "unistring dev headers"  libunistring-dev
    install_candidates "microhttpd dev headers" libmicrohttpd-dev

    install_candidates "curl dev headers" libcurl4-gnutls-dev libcurl4-openssl-dev
    install_candidates "ssl dev headers"  libssl-dev

    RUN sudo apt-get install -y "${APT_PKGS[@]}"

    if [[ -d .git ]]; then
        RUN git submodule update --init --recursive
    fi

    log "Setup complete"
}

if [[ "$DO_SETUP" -eq 1 ]]; then
    ubuntu_setup
fi

have podman || die "podman not installed"

VERSION="$(make -s print-version)"
[[ -n "$VERSION" ]] || die "VERSION empty"

LOCAL_IMAGE="localhost/${REPO_SERVER_URL}/${REPO_NAME}:${VERSION}"
REG_IMAGE="${REGISTRY}/${REPO_SERVER_URL}/${REPO_NAME}:${VERSION}"
REGISTRY_NAME="ukama_local_registry"

cleanup() {
    podman rm -f "$REGISTRY_NAME" >/dev/null 2>&1 || true
}
trap cleanup EXIT

RUN make clean
RUN make container

podman rm -f "$REGISTRY_NAME" >/dev/null 2>&1 || true
RUN podman run -d --name "$REGISTRY_NAME" --network host registry:latest

RUN podman tag "$LOCAL_IMAGE" "$REG_IMAGE"
RUN podman push --tls-verify=false "$REG_IMAGE"

RUN podman run --network host --privileged -it \
    -v "${PWD}/ukama_${VERSION}.tgz:/ukama/ukama.tgz:ro" \
    --tls-verify=false \
    -e VNODE_METADATA="$VNODE_METADATA" \
    -e VNODE_ID="$VNODE_ID" \
    -e VNODE_RUN_TARGET="local" \
    -e REPO_SERVER_URL="$REPO_SERVER_URL" \
    -e REPO_NAME="$REPO_NAME" \
    -e BOOTSTRAP_SERVER="$BOOTSTRAP_SERVER" \
    "$REG_IMAGE"

OUT_REG_IMAGE="${REGISTRY}/${REPO_SERVER_URL}/${REPO_NAME}:${VNODE_ID}"
RUN podman pull --tls-verify=false "$OUT_REG_IMAGE"

OUT_LOCAL_IMAGE="${REPO_SERVER_URL}/${REPO_NAME}:${VNODE_ID}"
RUN podman tag "$OUT_REG_IMAGE" "$OUT_LOCAL_IMAGE"

RUN make clean

echo "Done. Image available: ${OUT_LOCAL_IMAGE}"
