#!/bin/sh
# Cleanup leaked ukama-lab Podman networks/containers and node image tags.
# Keeps virtualnode-base:* and git-sha virtualnode images.

set -eu

echo "cleanup-stale: remove old ukama-lab networks"

podman network ls --format '{{.Name}}' |
    grep '^ukama-lab-' |
    while read -r net; do
        echo "cleanup-stale: network $net"

        containers="$(podman ps -a --filter "network=$net" --format '{{.Names}}' 2>/dev/null || true)"
        for c in $containers; do
            echo "cleanup-stale: rm container $c"
            podman rm -f "$c" >/dev/null 2>&1 || true
        done

        podman network rm "$net" >/dev/null 2>&1 || true
    done

echo "cleanup-stale: remove old node image tags"

podman images --format '{{.Repository}}:{{.Tag}}' |
    grep -E '^(localhost/)?testing/virtualnode:uk-sa.*-(t|c|a)node-|^localhost:5000/testing/virtualnode:uk-sa.*-(t|c|a)node-' |
    while read -r img; do
        echo "cleanup-stale: rmi tag $img"
        podman rmi "$img" >/dev/null 2>&1 || true
    done

echo "cleanup-stale: done"
