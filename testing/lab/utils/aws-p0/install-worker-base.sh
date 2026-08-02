#!/usr/bin/env bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

# Install everything required by an AWS P0 worker, including the native
# dependencies used by Ukama testing/node virtual-node builds.

set -Eeuo pipefail

if ((EUID != 0)); then
    printf 'run as root: sudo %s\n' "$0" >&2
    exit 2
fi

if command -v apt-get >/dev/null 2>&1; then
    export DEBIAN_FRONTEND=noninteractive

    apt-get update
    apt-get install -y --no-install-recommends \
        awscli \
        ca-certificates \
        curl \
        wget \
        git \
        gzip \
        tar \
        xz-utils \
        rsync \
        file \
        jq \
        python3 \
        iproute2 \
        iptables \
        openvswitch-switch \
        podman \
        build-essential \
        libc6-dev \
        binutils \
        pkg-config \
        autoconf \
        automake \
        libtool \
        libtool-bin \
        m4 \
        bison \
        flex \
        cmake \
        meson \
        ninja-build \
        patchelf \
        golang-go \
        libmicrohttpd-dev \
        libgnutls28-dev \
        nettle-dev \
        libp11-kit-dev \
        zlib1g-dev \
        libcurl4-openssl-dev \
        libssl-dev \
        libjansson-dev \
        uuid-dev \
        libsqlite3-dev \
        libzstd-dev \
        libgmp-dev

elif command -v dnf >/dev/null 2>&1; then
    # The current Ukama worker AMI is Ubuntu. Keep the generic RPM path, but
    # fail clearly rather than pretending its native build set is equivalent.
    dnf install -y \
        awscli ca-certificates curl gcc gcc-c++ glibc-devel git gzip \
        iproute iptables jq make openvswitch podman python3 tar

    printf '\nThe complete Ukama native build dependency set is currently\n' >&2
    printf 'defined for the Ubuntu worker AMI only.\n' >&2
    exit 2
else
    printf 'unsupported package manager\n' >&2
    exit 2
fi

if command -v systemctl >/dev/null 2>&1; then
    systemctl enable --now openvswitch-switch 2>/dev/null ||
        systemctl enable --now openvswitch 2>/dev/null || true
fi

printf '\nVerifying compiler startup files...\n'
crti="$(gcc -print-file-name=crti.o)"
if [[ "$crti" == "crti.o" || ! -f "$crti" ]]; then
    printf 'missing crti.o after package installation\n' >&2
    exit 1
fi
printf 'found %s\n' "$crti"

printf '\nVerifying external libraries used by Ukama Makefiles...\n'
tmp_dir="$(mktemp -d)"
trap 'rm -rf "$tmp_dir"' EXIT

cat >"$tmp_dir/link-check.c" <<'EOF_C'
int main(void) {
    return 0;
}
EOF_C

gcc "$tmp_dir/link-check.c" \
    -Wl,--no-as-needed \
    -lmicrohttpd \
    -lgnutls \
    -lnettle \
    -lhogweed \
    -lp11-kit \
    -lz \
    -lcurl \
    -lssl \
    -lcrypto \
    -ljansson \
    -lpthread \
    -lrt \
    -lm \
    -luuid \
    -lsqlite3 \
    -lzstd \
    -lgmp \
    -o "$tmp_dir/link-check"

"$tmp_dir/link-check"
