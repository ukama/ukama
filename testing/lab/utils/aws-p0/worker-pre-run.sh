#!/usr/bin/env bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

# Optional worker-local preparation. Keep this file idempotent.
# It is executed as root after both source archives are extracted and before
# run-p0-scenarios.sh starts.

set -Eeuo pipefail

# Start Open vSwitch using the service name available on the worker image.
if command -v systemctl >/dev/null 2>&1; then
    systemctl start openvswitch-switch 2>/dev/null ||
        systemctl start openvswitch 2>/dev/null || true
fi

# Fail clearly if the reusable AMI is missing the stable worker dependencies.
for command_name in aws jq curl tar gzip podman ovs-vsctl python3; do
    command -v "$command_name" >/dev/null 2>&1 || {
        printf 'missing worker AMI dependency: %s\n' "$command_name" >&2
        exit 1
    }
done

ovs-vsctl show >/dev/null
