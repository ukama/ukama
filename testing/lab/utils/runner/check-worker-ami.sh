#!/usr/bin/env bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

# Run on the prepared EC2 builder before capturing the reusable worker AMI.

set -Eeuo pipefail

missing=0
for command_name in aws bash curl gzip jq ovs-vsctl podman python3 tar; do
    if command -v "$command_name" >/dev/null 2>&1; then
        printf 'ok       %s\n' "$command_name"
    else
        printf 'missing  %s\n' "$command_name" >&2
        missing=1
    fi
done

if ((missing)); then
    exit 1
fi

podman info >/dev/null
ovs-vsctl show >/dev/null

printf '\nworker base checks passed\n'
printf 'A successful real ukama-lab scenario is still required before capture.\n'
