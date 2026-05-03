#!/bin/sh
set -eu

ROOT="$(CDPATH= cd -- "$(dirname -- "$0")/../.." && pwd)"
PKG_DIR="/ukama/apps/pkgs"
DB_DIR="/ukama/apps/db"

if [ "${UKAMA_TEST_CONTAINER:-}" != "1" ]; then
    echo "ERROR: integration tests must run inside podman via make test"
    exit 1
fi

if [ ! -x "$ROOT/wimc.d" ]; then
    echo "SKIP: $ROOT/wimc.d not found. Build wimc.d first to run integration."
    exit 0
fi

rm -rf "$PKG_DIR" "$DB_DIR"
mkdir -p "$PKG_DIR" "$DB_DIR"

grep -q '^wimc[[:space:]]' /etc/services || \
    echo 'wimc 19079/tcp' >> /etc/services
grep -q '^wimc-agent-chunk[[:space:]]' /etc/services || \
    echo 'wimc-agent-chunk 19081/tcp' >> /etc/services
grep -q '^ukama[[:space:]]' /etc/services || \
    echo 'ukama 19080/tcp' >> /etc/services

python3 "$ROOT/test/tools/fake_hub.py" 19080 &
HUB_PID=$!

WIMC_FIXTURE_TARBALLS="$ROOT/test/fixtures/tarballs" \
WIMC_TEST_PKG_DIR="$PKG_DIR" \
python3 "$ROOT/test/tools/fake_agent.py" 19081 &
AGENT_PID=$!

cleanup() {
    kill "$HUB_PID" "$AGENT_PID" "${WIMC_PID:-}" 2>/dev/null || true
}
trap cleanup EXIT INT TERM

"$ROOT/wimc.d" -u http://127.0.0.1:19080 -l INFO &
WIMC_PID=$!

for i in $(seq 1 30); do
    if curl -fsS http://127.0.0.1:19079/v1/ping >/dev/null 2>&1; then
        break
    fi
    sleep 0.2
done

curl -fsS http://127.0.0.1:19079/v1/ping >/dev/null

curl -fsS \
    -X POST \
    -H 'Content-Type: application/json' \
    -d '{"hub":"http://127.0.0.1:19080"}' \
    http://127.0.0.1:19079/v1/apps/example/v1-abc \
    >/tmp/wimc-post.json

for i in $(seq 1 30); do
    curl -fsS \
        http://127.0.0.1:19079/v1/apps/example/v1-abc/status \
        >/tmp/wimc-status.json || true

    if grep -q '"status"[[:space:]]*:[[:space:]]*"available"' \
        /tmp/wimc-status.json; then
        break
    fi
    sleep 0.2
done

grep -q '"status"[[:space:]]*:[[:space:]]*"available"' \
    /tmp/wimc-status.json

test -f /ukama/apps/pkgs/example_v1-abc.tar.gz

echo "PASS: wimc.d operation test without real agent/casync"
