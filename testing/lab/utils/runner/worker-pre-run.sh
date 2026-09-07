#!/usr/bin/env bash
# Worker-local preparation. Runs as root after source extraction and before P0.

set -Eeuo pipefail

if command -v systemctl >/dev/null 2>&1; then
    systemctl start openvswitch-switch 2>/dev/null ||
        systemctl start openvswitch 2>/dev/null || true
fi

for command_name in aws jq curl tar gzip podman ovs-vsctl python3 ip; do
    command -v "$command_name" >/dev/null 2>&1 || {
        printf 'missing worker AMI dependency: %s\n' "$command_name" >&2
        exit 1
    }
done

ovs-vsctl show >/dev/null

# Turn connectivity failures into one worker infrastructure failure instead of
# making every scenario in the shard look like a product failure.
check_url() {
    local label="$1"
    local url="$2"
    local code

    [[ -n "$url" ]] || {
        printf '%s URL is empty\n' "$label" >&2
        return 1
    }

    code="$(curl --noproxy '*' -sS -o /dev/null \
        --connect-timeout "${P0_CONNECT_TIMEOUT_SECONDS:-15}" \
        --max-time "${P0_HTTP_TIMEOUT_SECONDS:-25}" \
        -w '%{http_code}' "$url")" || {
        printf '%s is unreachable from worker: %s\n' "$label" "$url" >&2
        return 1
    }

    printf '%s reachable: HTTP %s (%s)\n' "$label" "$code" "$url"
}

check_url PAUTH "${PAUTH_URL:-}"
check_url BFF "${BFF_BASE_URL:-}"

check_url FACTORY "${ULAB_FACTORY_SEED_URL:-}"
check_url WAREHOUSE "${UKAMA_LAB_WAREHOUSE_URL:-}"
