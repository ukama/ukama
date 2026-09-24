#!/bin/sh
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

set -eu

if [ "$#" -lt 2 ]; then
    echo "usage: $0 <check|install|start|stop|remove> <container> [run-dir]" >&2
    exit 2
fi

ACTION="$1"
CONTAINER_NAME="$2"
UNIT="ulab-${CONTAINER_NAME}.service"
RESTART_DELAY="${ULAB_VNODE_RESTART_DELAY_SEC:-15}"

case "$CONTAINER_NAME" in
    ''|*[!A-Za-z0-9_.-]*) echo "node-host: invalid container name" >&2; exit 2 ;;
esac

host_systemctl() {
    if [ "$(id -u)" -eq 0 ]; then
        systemctl "$@"
    else
        systemctl --user "$@"
    fi
}

check_manager() {
    case "$RESTART_DELAY" in
        ''|*[!0-9]*)
            echo "node-host: ULAB_VNODE_RESTART_DELAY_SEC must be whole seconds" >&2
            exit 2
            ;;
    esac
    command -v systemctl >/dev/null 2>&1 || {
        echo "node-host: systemctl is required for delayed node restarts" >&2
        exit 1
    }
    if ! host_systemctl show --property=Version --value >/dev/null; then
        echo "node-host: no systemd manager for the current Podman user" >&2
        exit 1
    fi
    podman generate systemd --help >/dev/null
}

case "$ACTION" in
    check)
        check_manager
        exit 0
        ;;
    install)
        [ "$#" -eq 3 ] || { echo "node-host: install requires run-dir" >&2; exit 2; }
        check_manager
        UNIT_DIR="$(CDPATH= cd -- "$3" && pwd)/runtime-systemd"
        mkdir -p "$UNIT_DIR"
        UNIT_FILE="$UNIT_DIR/$UNIT"
        TEMP_FILE="$UNIT_FILE.tmp.$$"
        trap 'rm -f "$TEMP_FILE"' EXIT

        # Without --new, the generated unit starts the SAME container and
        # supervises conmon. Its writable filesystem survives each reboot.
        podman generate systemd --name --container-prefix=ulab \
            --restart-policy=always \
            "$CONTAINER_NAME" > "$TEMP_FILE"
        # Older Podman versions lack --restart-sec. Set the systemd directive
        # last so it also overrides any delay emitted by newer generators.
        printf '\n[Service]\nRestartSec=%ss\n' "$RESTART_DELAY" >> "$TEMP_FILE"
        mv "$TEMP_FILE" "$UNIT_FILE"
        host_systemctl link --runtime "$UNIT_FILE"
        host_systemctl daemon-reload
        host_systemctl start "$UNIT"
        echo "node-host: unit=$UNIT restart_delay=${RESTART_DELAY}s"
        exit 0
        ;;
    start|stop|remove) ;;
    *) echo "node-host: unknown action $ACTION" >&2; exit 2 ;;
esac

MANAGED="$(podman inspect --format \
    '{{ index .Config.Labels "io.ukama.lab.restart-manager" }}' \
    "$CONTAINER_NAME" 2>/dev/null || true)"
LOAD_STATE=""
if command -v systemctl >/dev/null 2>&1; then
    LOAD_STATE="$(host_systemctl show --property=LoadState --value \
        "$UNIT" 2>/dev/null || true)"
fi

if [ -n "$LOAD_STATE" ] && [ "$LOAD_STATE" != "not-found" ]; then
    case "$ACTION" in
        start)
            host_systemctl reset-failed "$UNIT"
            host_systemctl start "$UNIT"
            ;;
        stop)
            # Also cancels a pending restart while the container is down.
            host_systemctl stop "$UNIT"
            ;;
        remove)
            host_systemctl stop "$UNIT"
            host_systemctl disable --runtime "$UNIT"
            host_systemctl daemon-reload
            ;;
    esac
elif [ "$MANAGED" = "systemd" ]; then
    if [ "$ACTION" != "remove" ] || [ "$LOAD_STATE" != "not-found" ]; then
        echo "node-host: cannot $ACTION managed container; service $UNIT unavailable" >&2
        exit 1
    fi
else
    # Containers from older runs still use Podman's restart policy.
    case "$ACTION" in
        start) podman start "$CONTAINER_NAME" >/dev/null ;;
        stop) podman stop "$CONTAINER_NAME" >/dev/null ;;
    esac
fi
