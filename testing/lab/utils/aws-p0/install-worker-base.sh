#!/usr/bin/env bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

# Install the stable, generic worker packages on a fresh EC2 builder.
# Project-specific Ukama virtual-node build dependencies may still be needed.

set -Eeuo pipefail

if ((EUID != 0)); then
    printf 'run as root: sudo %s\n' "$0" >&2
    exit 2
fi

if command -v apt-get >/dev/null 2>&1; then
    apt-get update
    DEBIAN_FRONTEND=noninteractive apt-get install -y \
        awscli ca-certificates curl gcc git gzip iproute2 iptables jq make \
        openvswitch-switch podman python3 tar
elif command -v dnf >/dev/null 2>&1; then
    dnf install -y \
        awscli ca-certificates curl gcc git gzip iproute iptables jq make \
        openvswitch podman python3 tar
else
    printf 'unsupported package manager; install worker dependencies manually\n' >&2
    exit 2
fi

if command -v systemctl >/dev/null 2>&1; then
    systemctl enable --now openvswitch-switch 2>/dev/null ||
        systemctl enable --now openvswitch 2>/dev/null || true
fi

printf '\nGeneric worker packages installed.\n'
printf 'Now install any dependencies required by Ukama testing/node builds,\n'
printf 'run one normal scenario, then run check-worker-ami.sh.\n'
