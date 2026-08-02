#!/usr/bin/env bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

# Disposable EC2 worker. Inputs and outputs are regular files in S3.

set -Eeuo pipefail
umask 077

# cloud-init/systemd starts the worker without a login shell, so HOME may be
# absent. Podman network helpers use HOME for their per-user CNI paths.
export HOME="${HOME:-/root}"
mkdir -p "$HOME"

: "${BATCH_ID:?}"
: "${WORKER_ID:?}"
: "${S3_BUCKET:?}"
: "${S3_PREFIX:?}"
: "${AWS_REGION:?}"

export AWS_DEFAULT_REGION="$AWS_REGION"

WORK_ROOT="${WORK_ROOT:-/opt/ukama-p0}"
INPUT_URI="s3://${S3_BUCKET}/${S3_PREFIX%/}/${BATCH_ID}/input"
OUTPUT_URI="s3://${S3_BUCKET}/${S3_PREFIX%/}/${BATCH_ID}/output"
LOG_FILE="/var/log/ukama-p0-worker.log"
STATUS_FILE="$WORK_ROOT/status.tsv"
DONE_FILE="$WORK_ROOT/worker.done"
FAILED_FILE="$WORK_ROOT/worker.failed"
RESULT_ARCHIVE="$WORK_ROOT/${WORKER_ID}.tar.gz"
UPLOAD_STOP="$WORK_ROOT/.stop-upload"
UPLOADER_PID=""
FINALIZED=0
RUNNER_RC=125

mkdir -p "$WORK_ROOT"
touch "$LOG_FILE"
exec > >(tee -a "$LOG_FILE") 2>&1

now_utc() {
    date -u +%Y-%m-%dT%H:%M:%SZ
}

instance_id() {
    local token
    token="$(curl -fsS -X PUT \
        -H 'X-aws-ec2-metadata-token-ttl-seconds: 60' \
        http://169.254.169.254/latest/api/token 2>/dev/null || true)"
    if [[ -n "$token" ]]; then
        curl -fsS \
            -H "X-aws-ec2-metadata-token: $token" \
            http://169.254.169.254/latest/meta-data/instance-id 2>/dev/null || true
    fi
}

write_worker_status() {
    local state="$1"
    local current="${2:--}"
    local message="${3:-}"

    {
        printf 'worker\t%s\n' "$WORKER_ID"
        printf 'instance_id\t%s\n' "${INSTANCE_ID:-unknown}"
        printf 'state\t%s\n' "$state"
        printf 'total\t%s\n' "${TOTAL_SCENARIOS:-0}"
        printf 'completed\t%s\n' "${COMPLETED_SCENARIOS:-0}"
        printf 'passed\t%s\n' "${PASSED_SCENARIOS:-0}"
        printf 'failed\t%s\n' "${FAILED_SCENARIOS:-0}"
        printf 'skipped\t%s\n' "${SKIPPED_SCENARIOS:-0}"
        printf 'current\t%s\n' "$current"
        printf 'message\t%s\n' "$message"
        printf 'started_at\t%s\n' "${STARTED_AT:-$(now_utc)}"
        printf 'updated_at\t%s\n' "$(now_utc)"
    } >"$STATUS_FILE.tmp"
    mv "$STATUS_FILE.tmp" "$STATUS_FILE"
}

upload_status_once() {
    if [[ -f "$STATUS_FILE" ]]; then
        aws s3 cp "$STATUS_FILE" \
            "$OUTPUT_URI/${WORKER_ID}.status" \
            --only-show-errors || true
    fi
}

status_uploader() {
    local interval="${STATUS_INTERVAL_SECONDS:-10}"

    while [[ ! -e "$UPLOAD_STOP" ]]; do
        upload_status_once
        sleep "$interval"
    done
}

start_watchdog() {
    local minutes="${WORKER_TIMEOUT_MINUTES:-180}"

    if command -v systemd-run >/dev/null 2>&1; then
        systemd-run \
            --unit="ukama-p0-expiry-${WORKER_ID}" \
            --on-active="${minutes}m" \
            /usr/sbin/shutdown -P now >/dev/null 2>&1 || true
    else
        (
            sleep "$((minutes * 60))"
            /usr/sbin/shutdown -P now
        ) >/dev/null 2>&1 &
    fi
}

shutdown_worker() {
    if [[ "${P0_NO_SHUTDOWN:-0}" == "1" ]]; then
        return 0
    fi
    sync
    /usr/sbin/shutdown -P now || true
}

finalize() {
    local shell_rc=$?

    if ((FINALIZED)); then
        return
    fi
    FINALIZED=1
    trap - EXIT INT TERM

    touch "$UPLOAD_STOP"
    if [[ -n "$UPLOADER_PID" ]]; then
        wait "$UPLOADER_PID" 2>/dev/null || true
    fi

    set +e

    if [[ -d "$WORK_ROOT/results" ]]; then
        tar -C "$WORK_ROOT" -czf "$RESULT_ARCHIVE" results
        aws s3 cp "$RESULT_ARCHIVE" \
            "$OUTPUT_URI/${WORKER_ID}.tar.gz" \
            --only-show-errors
    fi

    upload_status_once
    if [[ -f "$LOG_FILE" ]]; then
        aws s3 cp "$LOG_FILE" \
            "$OUTPUT_URI/${WORKER_ID}.log" \
            --only-show-errors || true
    fi

    if [[ -f "$DONE_FILE" ]]; then
        aws s3 cp "$DONE_FILE" \
            "$OUTPUT_URI/${WORKER_ID}.done" \
            --only-show-errors
    else
        if [[ ! -f "$FAILED_FILE" ]]; then
            printf 'exit_code=%s\n' "$shell_rc" >"$FAILED_FILE"
        fi
        aws s3 cp "$FAILED_FILE" \
            "$OUTPUT_URI/${WORKER_ID}.failed" \
            --only-show-errors
    fi

    shutdown_worker
}

trap finalize EXIT INT TERM

STARTED_AT="$(now_utc)"
INSTANCE_ID="$(instance_id)"
write_worker_status STARTING - 'downloading batch input'
start_watchdog
status_uploader &
UPLOADER_PID=$!

command -v aws >/dev/null 2>&1 || {
    write_worker_status FAILED - 'aws CLI is missing from the worker AMI'
    exit 1
}
command -v tar >/dev/null 2>&1 || {
    write_worker_status FAILED - 'tar is missing from the worker AMI'
    exit 1
}

printf 'batch=%s worker=%s instance=%s\n' \
    "$BATCH_ID" "$WORKER_ID" "${INSTANCE_ID:-unknown}"

aws s3 cp "$INPUT_URI/worker.env" "$WORK_ROOT/worker.env" --only-show-errors
# shellcheck disable=SC1090
set -a
. "$WORK_ROOT/worker.env"
set +a

: "${SECRET_ID:?SECRET_ID is missing from worker.env}"

aws s3 cp "$INPUT_URI/ukama-lab.tar.gz" \
    "$WORK_ROOT/ukama-lab.tar.gz" --only-show-errors
aws s3 cp "$INPUT_URI/ukama.tar.gz" \
    "$WORK_ROOT/ukama.tar.gz" --only-show-errors
aws s3 cp "$INPUT_URI/checksums.sha256" \
    "$WORK_ROOT/checksums.sha256" --only-show-errors
aws s3 cp "$INPUT_URI/shards/${WORKER_ID}.txt" \
    "$WORK_ROOT/scenarios.txt" --only-show-errors

(
    cd "$WORK_ROOT"
    sha256sum -c checksums.sha256
)

SECRET_JSON="$(aws secretsmanager get-secret-value \
    --secret-id "$SECRET_ID" \
    --query SecretString --output text)"

eval "$(jq -r '
    to_entries[]
    | select(.key | test("^[A-Za-z_][A-Za-z0-9_]*$"))
    | "export \(.key)=\(.value | tostring | @sh)"
' <<<"$SECRET_JSON")"
unset SECRET_JSON

mkdir -p "$WORK_ROOT/ukama-lab" "$WORK_ROOT/ukama" "$WORK_ROOT/results"
tar -xzf "$WORK_ROOT/ukama-lab.tar.gz" -C "$WORK_ROOT/ukama-lab"
tar -xzf "$WORK_ROOT/ukama.tar.gz" -C "$WORK_ROOT/ukama"

LAB_ROOT="$WORK_ROOT/ukama-lab"
export UKAMA_REPO="$WORK_ROOT/ukama"
git config --global --add safe.directory "$UKAMA_REPO"
export LAB_BIN="$LAB_ROOT/bin/ukama-lab"
export SCENARIO_ROOT="scenarios/p0"
export P0_RUNS_DIR="$WORK_ROOT/results"
export P0_STATUS_FILE="$STATUS_FILE"

# Source archives intentionally exclude .git. Use one explicit version for all
# virtual-node and application builds in this disposable worker.
export UKAMA_APP_VERSION="${UKAMA_APP_VERSION:-v0.0.0-p0-${BATCH_ID}}"
printf 'source version: %s\n' "$UKAMA_APP_VERSION"

TOTAL_SCENARIOS="$(grep -Ev '^[[:space:]]*(#|$)' \
    "$WORK_ROOT/scenarios.txt" | wc -l | tr -d ' ')"
write_worker_status PREPARING - 'checking worker environment'

if [[ -x "$LAB_ROOT/utils/aws-p0/worker-pre-run.sh" ]]; then
    (
        cd "$LAB_ROOT"
        ./utils/aws-p0/worker-pre-run.sh
    )
fi

[[ -x "$LAB_BIN" ]] || {
    write_worker_status FAILED - "ukama-lab binary is not executable: $LAB_BIN"
    exit 1
}

write_worker_status RUNNING - 'starting scenario runner'
RUNNER_BATCH_ID="${BATCH_ID}-${WORKER_ID}"

set +e
(
    cd "$LAB_ROOT"
    ./utils/run-p0-scenarios.sh \
        --scenario-list "$WORK_ROOT/scenarios.txt" \
        --batch-id "$RUNNER_BATCH_ID" \
        --factory-nodes 0
)
RUNNER_RC=$?
set -e

REPORT_FILE="$WORK_ROOT/results/$RUNNER_BATCH_ID/batch-report.json"
if [[ -f "$REPORT_FILE" ]]; then
    COMPLETED_SCENARIOS="$(jq -r '.total // 0' "$REPORT_FILE")"
    PASSED_SCENARIOS="$(jq -r '.passed // 0' "$REPORT_FILE")"
    FAILED_SCENARIOS="$(jq -r '.failed // 0' "$REPORT_FILE")"
    SKIPPED_SCENARIOS="$(jq -r '.skipped // 0' "$REPORT_FILE")"
    write_worker_status DONE - "runner_exit_code=$RUNNER_RC"
    printf 'runner_exit_code=%s\n' "$RUNNER_RC" >"$DONE_FILE"
else
    write_worker_status FAILED - "runner exited $RUNNER_RC without batch-report.json"
    printf 'exit_code=%s\nreason=no_batch_report\n' \
        "$RUNNER_RC" >"$FAILED_FILE"
fi

# Scenario failures are represented in batch-report.json and are not worker
# infrastructure failures. The worker can terminate normally after upload.
exit 0
