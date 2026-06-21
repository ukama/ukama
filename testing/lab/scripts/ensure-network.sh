#!/bin/sh
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

set -eu

if [ "$#" -ne 1 ]; then
    echo "usage: $0 <run-dir>" >&2
    exit 2
fi

RUN_DIR="$1"
RUN_ID="$(basename "$RUN_DIR")"
SAFE_RUN_ID="$(printf "%s" "$RUN_ID" | tr -c 'A-Za-z0-9-' '-' | sed 's/^-*//;s/-*$//')"
LAB_NET="ukama-lab-$SAFE_RUN_ID"
STATE_DIR="$RUN_DIR/runtime-net"
STATE_FILE="$STATE_DIR/net.env"
PROBE_IMAGE="${ULAB_NET_PROBE_IMAGE:-alpine:3.20}"

need_cmd() {
    if ! command -v "$1" >/dev/null 2>&1; then
        echo "missing required command: $1" >&2
        exit 1
    fi
}

cni_config_path() {
    for p in \
        "$HOME/.config/cni/net.d/$LAB_NET.conflist" \
        "/etc/cni/net.d/$LAB_NET.conflist"; do
        if [ -f "$p" ]; then
            printf "%s\n" "$p"
            return 0
        fi
    done
    return 1
}

network_backend() {
    podman info --format '{{.Host.NetworkBackend}}' 2>/dev/null || echo "unknown"
}

normalize_cni_config() {
    cfg="$(cni_config_path || true)"

    if [ -z "$cfg" ]; then
        return 0
    fi

    # Some rootless Podman installs create cniVersion 1.0.0 while the installed
    # firewall CNI plugin only supports older CNI spec versions. 0.4.0 works
    # with the standard bridge/portmap/firewall/tuning plugins and avoids
    # failing later during node startup.
    if grep -q '"cniVersion"[[:space:]]*:[[:space:]]*"1\.0\.0"' "$cfg"; then
        echo "network: normalize CNI config $cfg cniVersion 1.0.0 -> 0.4.0"
        cp "$cfg" "$cfg.bak" 2>/dev/null || true
        sed -i 's/"cniVersion"[[:space:]]*:[[:space:]]*"1\.0\.0"/"cniVersion": "0.4.0"/' "$cfg"
    fi
}

remove_network() {
    podman network rm "$LAB_NET" >/dev/null 2>&1 || true
    rm -f "$HOME/.config/cni/net.d/$LAB_NET.conflist" 2>/dev/null || true
}

create_network() {
    if podman network exists "$LAB_NET" >/dev/null 2>&1; then
        echo "network: exists $LAB_NET"
    else
        echo "network: create $LAB_NET"
        podman network create --disable-dns "$LAB_NET" >/dev/null
    fi

    normalize_cni_config
}

probe_network() {
    out="$(mktemp "${TMPDIR:-/tmp}/ukama-net-probe.XXXXXX")"

    if ! podman image exists "$PROBE_IMAGE" >/dev/null 2>&1; then
        echo "network: pull probe image $PROBE_IMAGE"
        if ! podman pull "$PROBE_IMAGE" >"$out" 2>&1; then
            echo "network: cannot pull probe image $PROBE_IMAGE" >&2
            cat "$out" >&2 || true
            rm -f "$out"
            return 1
        fi
    fi

    if podman run --rm --network "$LAB_NET" "$PROBE_IMAGE" /bin/sh -c 'ip -4 addr >/dev/null' >"$out" 2>&1; then
        rm -f "$out"
        return 0
    fi

    echo "network: probe failed for $LAB_NET" >&2
    cat "$out" >&2 || true
    rm -f "$out"
    return 1
}

need_cmd podman
need_cmd sed
need_cmd grep
need_cmd mktemp

mkdir -p "$STATE_DIR"

BACKEND="$(network_backend)"
echo "network: backend=$BACKEND"

create_network

if ! probe_network; then
    echo "network: retry after clean recreate $LAB_NET" >&2
    remove_network
    create_network

    if ! probe_network; then
        echo "network: failed to create usable podman network $LAB_NET" >&2
        echo "network: backend=$BACKEND" >&2
        echo "network: inspect output:" >&2
        podman network inspect "$LAB_NET" >&2 || true
        cfg="$(cni_config_path || true)"
        if [ -n "$cfg" ]; then
            echo "network: cni config $cfg:" >&2
            cat "$cfg" >&2 || true
        fi
        exit 1
    fi
fi

cat > "$STATE_FILE" <<STATE
LAB_NET=$LAB_NET
PODMAN_NETWORK_BACKEND=$BACKEND
STATE

cfg="$(cni_config_path || true)"
if [ -n "$cfg" ]; then
    printf 'CNI_CONFIG=%s\n' "$cfg" >> "$STATE_FILE"
fi

echo "network-ready name=$LAB_NET backend=$BACKEND"
