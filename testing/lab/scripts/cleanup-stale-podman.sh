#!/bin/sh
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

# Cleanup leaked ukama-lab Podman networks and node image tags.
# Keeps base images and git-sha virtualnode images.

set -eu

echo "cleanup-stale: removing old ukama-lab networks"

podman network ls --format '{{.Name}}' |
    grep '^ukama-lab-' |
    while read -r net; do
        echo "remove stale network: $net"

        containers="$(podman ps -a --filter "network=$net" --format '{{.Names}}' 2>/dev/null || true)"
        for c in $containers; do
            echo "remove stale container: $c"
            podman rm -f "$c" >/dev/null 2>&1 || true
        done

        podman network rm "$net" >/dev/null 2>&1 || true
    done

echo "cleanup-stale: removing old node image tags"

podman images --format '{{.Repository}}:{{.Tag}}' |
    grep -E '^(localhost/)?testing/virtualnode:uk-sa.*-(t|c|a)node-|^localhost:5000/testing/virtualnode:uk-sa.*-(t|c|a)node-' |
    while read -r img; do
        echo "remove stale node image tag: $img"
        podman rmi "$img" >/dev/null 2>&1 || true
    done

echo "cleanup-stale: done"
