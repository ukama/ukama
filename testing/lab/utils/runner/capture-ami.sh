#!/usr/bin/env bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

# Capture a prepared, stopped/running EC2 instance as the reusable worker AMI.

set -Eeuo pipefail

SCRIPT_DIR="$(CDPATH= cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=lib.sh
. "$SCRIPT_DIR/lib.sh"
p0_load_config "$SCRIPT_DIR"

INSTANCE_ID="${1:-}"
[[ -n "$INSTANCE_ID" ]] || p0_die "usage: $0 INSTANCE_ID [IMAGE_NAME]"
IMAGE_NAME="${2:-ukama-lab-p0-worker-$(date -u +%Y%m%dt%H%M%Sz)}"

p0_require_cmd aws
p0_validate_simple_id INSTANCE_ID "$INSTANCE_ID"
p0_aws sts get-caller-identity >/dev/null

state="$(p0_aws ec2 describe-instances \
    --instance-ids "$INSTANCE_ID" \
    --query 'Reservations[0].Instances[0].State.Name' --output text)"
[[ "$state" != "None" ]] || p0_die "instance not found: $INSTANCE_ID"

printf 'creating AMI %s from %s (state=%s)\n' \
    "$IMAGE_NAME" "$INSTANCE_ID" "$state"
AMI_ID="$(p0_aws ec2 create-image \
    --instance-id "$INSTANCE_ID" \
    --name "$IMAGE_NAME" \
    --description 'Ukama lab disposable P0 worker' \
    --query ImageId --output text)"

p0_aws ec2 wait image-available --image-ids "$AMI_ID"
p0_aws ec2 create-tags \
    --resources "$AMI_ID" \
    --tags Key=Name,Value="$IMAGE_NAME" Key=Application,Value=ukama-lab-p0

state_tmp="$(mktemp)"
if [[ -r "$P0_AWS_STATE" ]]; then
    grep -v '^AMI_ID=' "$P0_AWS_STATE" >"$state_tmp" || true
else
    : >"$state_tmp"
fi
printf 'AMI_ID=%s\n' "$AMI_ID" >>"$state_tmp"
mv "$state_tmp" "$P0_AWS_STATE"
chmod 600 "$P0_AWS_STATE"

printf 'AMI ready: %s\n' "$AMI_ID"
printf 'saved to: %s\n' "$P0_AWS_STATE"
