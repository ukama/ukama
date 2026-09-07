#!/usr/bin/env bash
# Worker-local preparation. Runs as root after source extraction and before P0.

set -Eeuo pipefail

# systemd/cloud-init workers may not receive a login-shell HOME.
export HOME="${HOME:-/root}"
mkdir -p "$HOME"

if command -v systemctl >/dev/null 2>&1; then
    systemctl start openvswitch-switch 2>/dev/null ||
        systemctl start openvswitch 2>/dev/null || true
fi

for command_name in aws jq curl tar gzip podman ovs-vsctl python3 ip git make; do
    command -v "$command_name" >/dev/null 2>&1 || {
        printf 'missing worker AMI dependency: %s\n' "$command_name" >&2
        exit 1
    }
done

ovs-vsctl show >/dev/null

# The controller intentionally excludes the large .git directory from the
# Ukama source archive. Some Ukama makefiles still derive VERSION with git.
# Create disposable local metadata so those builds receive a valid, stable
# version without shipping repository history to every worker.
: "${UKAMA_REPO:?UKAMA_REPO is not set}"

if ! git -C "$UKAMA_REPO" rev-parse --is-inside-work-tree >/dev/null 2>&1; then
    printf 'source metadata: creating disposable Git metadata in %s\n' \
        "$UKAMA_REPO"

    git -C "$UKAMA_REPO" init -q
    git -C "$UKAMA_REPO" config user.name "Ukama P0 Worker"
    git -C "$UKAMA_REPO" config user.email "p0-worker@ukama.local"

    GIT_AUTHOR_DATE='2026-01-01T00:00:00Z' \
    GIT_COMMITTER_DATE='2026-01-01T00:00:00Z' \
        git -C "$UKAMA_REPO" commit -q --allow-empty \
            -m 'Disposable P0 worker source snapshot'

    git -C "$UKAMA_REPO" tag -f ukama-p0-source >/dev/null
fi

source_version="$(
    git -C "$UKAMA_REPO" describe --always --dirty=-dirty 2>/dev/null ||
        git -C "$UKAMA_REPO" rev-parse --short HEAD
)"
[[ -n "$source_version" ]] || {
    printf 'unable to establish Ukama source version\n' >&2
    exit 1
}
printf 'source version: %s\n' "$source_version"

# The udev endpoints resolve to private 10/8 addresses that are reachable from
# the known-good gateway host, not directly from ordinary EC2 workers.
if [[ -n "${BACKEND_GATEWAY_PRIVATE_IP:-}" ]]; then
    : "${BACKEND_ROUTE_CIDR:?BACKEND_ROUTE_CIDR is required when a gateway is configured}"
    : "${BACKEND_TEST_IP:?BACKEND_TEST_IP is required when a gateway is configured}"

    gateway_route="$(ip -4 route get "$BACKEND_GATEWAY_PRIVATE_IP")"
    gateway_dev="$(
        awk '{
            for (i = 1; i <= NF; i++) {
                if ($i == "dev") {
                    print $(i + 1)
                    exit
                }
            }
        }' <<<"$gateway_route"
    )"
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

    code="$(
        curl -k -sS -o /dev/null \
            --connect-timeout "${P0_CONNECT_TIMEOUT_SECONDS:-15}" \
            --max-time "${P0_HTTP_TIMEOUT_SECONDS:-25}" \
            -w '%{http_code}' "$url"
    )" || {
        printf '%s is unreachable from worker: %s\n' "$label" "$url" >&2
        return 1
    }

    printf '%s reachable: HTTP %s (%s)\n' "$label" "$code" "$url"
}

check_url PAUTH "${PAUTH_URL:-}"
check_url BFF "${BFF_BASE_URL:-}"
