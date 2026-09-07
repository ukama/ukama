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

# The udev endpoints resolve to private 10/8 addresses that are reachable from
# the known-good gateway host, not directly from ordinary EC2 workers.
if [[ -n "${BACKEND_GATEWAY_PRIVATE_IP:-}" ]]; then
    : "${BACKEND_ROUTE_CIDR:?BACKEND_ROUTE_CIDR is required when a gateway is configured}"
    : "${BACKEND_TEST_IP:?BACKEND_TEST_IP is required when a gateway is configured}"

    gateway_route="$(ip -4 route get "$BACKEND_GATEWAY_PRIVATE_IP")"
    gateway_dev="$(awk '{for (i=1;i<=NF;i++) if ($i=="dev") {print $(i+1); exit}}' \
        <<<"$gateway_route")"
    [[ -n "$gateway_dev" ]] || {
        printf 'cannot determine interface for backend gateway %s\n' \
            "$BACKEND_GATEWAY_PRIVATE_IP" >&2
        exit 1
    }

    ip -4 route replace "$BACKEND_ROUTE_CIDR" \
        via "$BACKEND_GATEWAY_PRIVATE_IP" dev "$gateway_dev"

    printf 'backend route: %s via %s dev %s\n' \
        "$BACKEND_ROUTE_CIDR" "$BACKEND_GATEWAY_PRIVATE_IP" "$gateway_dev"
    ip -4 route get "$BACKEND_TEST_IP"
fi

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

    code="$(curl -k -sS -o /dev/null \
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
