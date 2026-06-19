#!/bin/sh
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

set -eu

if [ "$#" -lt 3 ]; then
    echo "usage: $0 <base-image> <node-id> <output-image> [runtime]" >&2
    exit 2
fi

BASE_IMAGE="$1"
NODE_ID="$2"
OUT_IMAGE="$3"
NODE_RUNTIME="${4:-starter}"

safe_name() {
    printf "%s" "$1" | tr -c 'A-Za-z0-9_.-' '-'
}

need_cmd() {
    if ! command -v "$1" >/dev/null 2>&1; then
        echo "missing required command: $1" >&2
        exit 1
    fi
}

node_type() {
    case "$1" in
        *-tnode-*) echo "tnode" ;;
        *-cnode-*) echo "cnode" ;;
        *-anode-*) echo "anode" ;;
        *)
            echo "cannot determine node type from node id: $1" >&2
            exit 1
            ;;
    esac
}

entrypoint_change() {
    case "$1" in
        starter)
            printf '%s\n' 'ENTRYPOINT ["/sbin/starter.d"]'
            ;;
        supervisor)
            printf '%s\n' 'ENTRYPOINT ["/usr/bin/supervisord","-c","/etc/supervisor.conf"]'
            ;;
        *)
            echo "invalid runtime: $1" >&2
            exit 1
            ;;
    esac
}

need_cmd podman
need_cmd mktemp

NODE_TYPE="$(node_type "$NODE_ID")"
SAFE_NODE_ID="$(safe_name "$NODE_ID")"
TMP_CONTAINER="ukama-stamp-$SAFE_NODE_ID-$$"
ENTRYPOINT_CHANGE="$(entrypoint_change "$NODE_RUNTIME")"

cleanup() {
    podman rm -f "$TMP_CONTAINER" >/dev/null 2>&1 || true
}
trap cleanup EXIT INT TERM

if podman image exists "$OUT_IMAGE"; then
    echo "stamp-node-image: image exists $OUT_IMAGE"
    exit 0
fi

if ! podman image exists "$BASE_IMAGE"; then
    echo "stamp-node-image: missing base image $BASE_IMAGE" >&2
    exit 1
fi

echo "stamp-node-image: base=$BASE_IMAGE out=$OUT_IMAGE node=$NODE_ID type=$NODE_TYPE"

podman rm -f "$TMP_CONTAINER" >/dev/null 2>&1 || true

podman run -d \
    --name "$TMP_CONTAINER" \
    --privileged \
    --entrypoint /bin/sh \
    "$BASE_IMAGE" \
    -lc 'sleep 3600' >/dev/null

podman exec \
    -e "NODE_ID=$NODE_ID" \
    -e "NODE_TYPE=$NODE_TYPE" \
    "$TMP_CONTAINER" \
    sh -eu -c '
        echo "$NODE_ID" > /ukama/nodeid

        rm -rf /tmp/sys /ukama/mocksysfs/sys
        mkdir -p /ukama/mocksysfs /tmp

        case "$NODE_TYPE" in
            anode)
                /sbin/mock-sysfs-anode.sh --clean >/dev/null 2>&1 || true
                /sbin/mock-sysfs-anode.sh

                cd /
                /sbin/genSchema --u "$NODE_ID" \
                    --n ctrl --m UK-8001-RFC-1102 --f /mfgdata/schema/ctrl.json \
                    --n fe1  --m UK-8001-RFE-1103 --f /mfgdata/schema/fe1.json \
                    --n fe2  --m UK-8001-RFE-1104 --f /mfgdata/schema/fe2.json

                /sbin/genInventory \
                    -n ctrl -m UK-8001-RFC-1102 -f /mfgdata/schema/ctrl.json \
                    -n fe1  -m UK-8001-RFE-1103 -f /mfgdata/schema/fe1.json \
                    -n fe2  -m UK-8001-RFE-1104 -f /mfgdata/schema/fe2.json
                ;;

            tnode)
                /sbin/prepare_env.sh --clean >/dev/null 2>&1 || true
                /sbin/prepare_env.sh -u tnode -u anode

                cd /
                /sbin/genSchema --u "$NODE_ID" \
                    --n com  --m UK-SA9001-COM-A1-1103 --f /mfgdata/schema/com.json \
                    --n trx  --m UK-SA9001-TRX-A1-1103 --f /mfgdata/schema/trx.json \
                    --n mask --m UK-SA9001-MSK-A1-1103 --f /mfgdata/schema/mask.json

                /sbin/genInventory \
                    --n com  --m UK-SA9001-COM-A1-1103 --f /mfgdata/schema/com.json \
                    --n trx  --m UK-SA9001-TRX-A1-1103 --f /mfgdata/schema/trx.json \
                    --n mask --m UK-SA9001-MSK-A1-1103 --f /mfgdata/schema/mask.json
                ;;

            cnode)
                /sbin/mock-sysfs-cnode.sh --clean >/dev/null 2>&1 || true
                /sbin/mock-sysfs-cnode.sh

                cd /
                /sbin/genSchema --u "$NODE_ID" \
                    --n cm4 --m UK-SA2602-CM4-1102 --f /mfgdata/schema/cnode.json

                /sbin/genInventory \
                    --n cm4 --m UK-SA2602-CM4-1102 --f /mfgdata/schema/cnode.json
                ;;

            *)
                echo "unknown node type: $NODE_TYPE" >&2
                exit 1
                ;;
        esac

        if [ ! -d /tmp/sys ]; then
            echo "/tmp/sys not found after sysfs generation" >&2
            exit 1
        fi

        rm -rf /ukama/mocksysfs/sys
        mkdir -p /ukama/mocksysfs
        ( cd /tmp && tar -cf - sys ) | ( cd /ukama/mocksysfs && tar -xpf - )

        rm -rf /tmp/sys
        ln -s /ukama/mocksysfs/sys /tmp/sys
    '

podman commit \
    --change "$ENTRYPOINT_CHANGE" \
    --change 'CMD []' \
    "$TMP_CONTAINER" \
    "$OUT_IMAGE" >/dev/null

podman rm -f "$TMP_CONTAINER" >/dev/null
trap - EXIT INT TERM

echo "stamp-node-image: created $OUT_IMAGE"
