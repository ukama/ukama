#!/usr/bin/env bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

# Shared helpers for the plain EC2/S3 P0 runner.

set -o pipefail

p0_die() {
    printf 'error: %s\n' "$*" >&2
    exit 2
}

p0_note() {
    printf '%s\n' "$*"
}

p0_require_cmd() {
    command -v "$1" >/dev/null 2>&1 || p0_die "missing command: $1"
}

p0_script_dir() {
    CDPATH= cd -- "$(dirname -- "${BASH_SOURCE[1]}")" && pwd
}

p0_load_config() {
    local script_dir="$1"
    local config_file="${P0_AWS_CONFIG:-$script_dir/config.env}"
    local state_file="${P0_AWS_STATE:-$script_dir/.state.env}"

    [[ -r "$config_file" ]] || p0_die "config file not found: $config_file (copy config.env.example to config.env)"

    # These are trusted, user-owned shell configuration files.
    set -a
    # shellcheck disable=SC1090
    . "$config_file"
    if [[ -r "$state_file" ]]; then
        # shellcheck disable=SC1090
        . "$state_file"
    fi
    set +a

    export P0_AWS_CONFIG="$config_file"
    export P0_AWS_STATE="$state_file"
    export AWS_REGION="${AWS_REGION:-${AWS_DEFAULT_REGION:-}}"
    export AWS_DEFAULT_REGION="$AWS_REGION"
    if [[ -n "${AWS_PROFILE:-}" ]]; then
        export AWS_PROFILE
    fi
}

p0_require_config() {
    local name
    for name in "$@"; do
        [[ -n "${!name:-}" ]] || p0_die "required config value is empty: $name"
    done
}

p0_validate_simple_id() {
    local name="$1"
    local value="$2"

    [[ "$value" =~ ^[A-Za-z0-9._:/+=,@-]+$ ]] ||
        p0_die "$name contains unsupported characters: $value"
}

p0_aws() {
    aws --no-cli-pager "$@"
}

p0_s3_batch_root() {
    local batch_id="$1"
    printf 's3://%s/%s/%s' "$S3_BUCKET" "${S3_PREFIX%/}" "$batch_id"
}

p0_local_batch_root() {
    local lab_root="$1"
    local batch_id="$2"
    local runs_root="${LOCAL_RUNS_DIR:-runs/p0-aws}"

    if [[ "$runs_root" = /* ]]; then
        printf '%s/%s' "${runs_root%/}" "$batch_id"
    else
        printf '%s/%s/%s' "${lab_root%/}" "${runs_root%/}" "$batch_id"
    fi
}

p0_instance_ids() {
    local batch_id="$1"

    p0_aws ec2 describe-instances \
        --filters \
            "Name=tag:UkamaP0Batch,Values=$batch_id" \
            'Name=instance-state-name,Values=pending,running,stopping,stopped' \
        --query 'Reservations[].Instances[].InstanceId' \
        --output text 2>/dev/null |
        tr '\t' '\n' |
        sed '/^$/d'
}

p0_cleanup_batch() {
    local batch_id="$1"
    local ids=()

    mapfile -t ids < <(p0_instance_ids "$batch_id")
    if ((${#ids[@]} == 0)); then
        return 0
    fi

    printf 'terminating %s remaining EC2 worker(s) for batch %s\n' \
        "${#ids[@]}" "$batch_id" >&2

    p0_aws ec2 terminate-instances \
        --instance-ids "${ids[@]}" >/dev/null || return 1

    p0_aws ec2 wait instance-terminated \
        --instance-ids "${ids[@]}" || return 1
}
