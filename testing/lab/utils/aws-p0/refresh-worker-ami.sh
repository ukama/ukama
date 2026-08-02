#!/usr/bin/env bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

# Launch a temporary AWS P0 worker builder, install/verify dependencies,
# capture a new AMI, update .state.env, and terminate the builder.

set -Eeuo pipefail

SCRIPT_DIR="$(CDPATH= cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"

# shellcheck source=lib.sh
. "$SCRIPT_DIR/lib.sh"
p0_load_config "$SCRIPT_DIR"

INSTALLER="$SCRIPT_DIR/install-worker-base.sh"
CHECKER="$SCRIPT_DIR/check-worker-ami.sh"
CAPTURE="$SCRIPT_DIR/capture-ami.sh"

SSH_USER="${P0_BUILDER_SSH_USER:-ubuntu}"
SSH_KEY="${P0_BUILDER_SSH_KEY:-/home/kashif/VirtualNodePEM.pem}"
KEY_NAME="${P0_BUILDER_KEY_NAME:-VirtualNodePEM}"
BUILDER_TYPE="${P0_BUILDER_INSTANCE_TYPE:-${INSTANCE_TYPE:-c7i.xlarge}}"
BUILDER_VOLUME_GB="${ROOT_VOLUME_GB:-250}"
BUILDER_VOLUME_TYPE="${ROOT_VOLUME_TYPE:-gp3}"

BUILDER_ID=""
BUILDER_HOST=""
CAPTURE_COMPLETE=0

cleanup() {
    local rc=$?

    if [[ -n "$BUILDER_ID" ]]; then
        printf '\nTerminating temporary builder %s...\n' "$BUILDER_ID"
        p0_aws ec2 terminate-instances \
            --instance-ids "$BUILDER_ID" >/dev/null 2>&1 || true
        p0_aws ec2 wait instance-terminated \
            --instance-ids "$BUILDER_ID" >/dev/null 2>&1 || true
    fi

    if ((rc != 0)); then
        printf '\nBuilder refresh failed with exit code %d.\n' "$rc" >&2
    fi

    exit "$rc"
}
trap cleanup EXIT INT TERM

p0_require_cmd aws
p0_require_cmd ssh
p0_require_cmd scp
p0_require_cmd grep

p0_require_config \
    AMI_ID \
    AWS_REGION \
    SUBNET_ID \
    SECURITY_GROUP_ID \
    INSTANCE_PROFILE_NAME

[[ -f "$INSTALLER" ]] ||
    p0_die "missing installer: $INSTALLER"
[[ -f "$CHECKER" ]] ||
    p0_die "missing checker: $CHECKER"
[[ -x "$CAPTURE" ]] ||
    p0_die "missing executable capture script: $CAPTURE"
[[ -r "$SSH_KEY" ]] ||
    p0_die "SSH key not readable: $SSH_KEY"

# Prevent accidentally using the old generic-only installer.
grep -q 'libmicrohttpd-dev' "$INSTALLER" ||
    p0_die "install-worker-base.sh is still the old generic-only version"

chmod 600 "$SSH_KEY"

printf 'Current AMI:       %s\n' "$AMI_ID"
printf 'Builder type:      %s\n' "$BUILDER_TYPE"
printf 'Builder disk:      %s GB %s\n' \
    "$BUILDER_VOLUME_GB" "$BUILDER_VOLUME_TYPE"
printf 'Subnet:            %s\n' "$SUBNET_ID"
printf 'Security group:    %s\n' "$SECURITY_GROUP_ID"
printf 'SSH key pair:      %s\n' "$KEY_NAME"

p0_aws sts get-caller-identity >/dev/null

ROOT_DEVICE="$(
    p0_aws ec2 describe-images \
        --image-ids "$AMI_ID" \
        --query 'Images[0].RootDeviceName' \
        --output text
)"

[[ -n "$ROOT_DEVICE" && "$ROOT_DEVICE" != "None" ]] ||
    p0_die "could not determine root device for $AMI_ID"

BLOCK_DEVICE_MAPPING="$(
    printf '%s' \
        "DeviceName=$ROOT_DEVICE,Ebs={VolumeSize=$BUILDER_VOLUME_GB," \
        "VolumeType=$BUILDER_VOLUME_TYPE,DeleteOnTermination=true}"
)"

printf '\nLaunching temporary builder...\n'

BUILDER_ID="$(
    p0_aws ec2 run-instances \
        --image-id "$AMI_ID" \
        --instance-type "$BUILDER_TYPE" \
        --subnet-id "$SUBNET_ID" \
        --security-group-ids "$SECURITY_GROUP_ID" \
        --key-name "$KEY_NAME" \
        --iam-instance-profile "Name=$INSTANCE_PROFILE_NAME" \
        --associate-public-ip-address \
        --block-device-mappings "$BLOCK_DEVICE_MAPPING" \
        --tag-specifications \
            'ResourceType=instance,Tags=[{Key=Name,Value=ukama-p0-ami-builder},{Key=Application,Value=ukama-lab-p0}]' \
        --query 'Instances[0].InstanceId' \
        --output text
)"

[[ "$BUILDER_ID" == i-* ]] ||
    p0_die "AWS returned an invalid builder instance ID: $BUILDER_ID"

printf 'Builder ID: %s\n' "$BUILDER_ID"

p0_aws ec2 wait instance-running \
    --instance-ids "$BUILDER_ID"

p0_aws ec2 wait instance-status-ok \
    --instance-ids "$BUILDER_ID"

BUILDER_HOST="$(
    p0_aws ec2 describe-instances \
        --instance-ids "$BUILDER_ID" \
        --query \
            'Reservations[0].Instances[0].[PublicDnsName,PublicIpAddress]' \
        --output text |
        awk '{ if ($1 != "None" && $1 != "") print $1; else print $2 }'
)"

[[ -n "$BUILDER_HOST" && "$BUILDER_HOST" != "None" ]] ||
    p0_die "builder has no public DNS name or IP address"

printf 'Builder host: %s\n' "$BUILDER_HOST"

SSH_OPTIONS=(
    -i "$SSH_KEY"
    -o BatchMode=yes
    -o ConnectTimeout=10
    -o StrictHostKeyChecking=accept-new
)

printf '\nWaiting for SSH...\n'
ssh_ready=0
for attempt in $(seq 1 60); do
    if ssh "${SSH_OPTIONS[@]}" \
            "$SSH_USER@$BUILDER_HOST" true 2>/dev/null; then
        ssh_ready=1
        break
    fi

    printf 'SSH not ready yet (%d/60)\n' "$attempt"
    sleep 10
done

((ssh_ready == 1)) ||
    p0_die "SSH did not become ready on $BUILDER_HOST"

printf '\nUploading worker installer and checker...\n'
scp "${SSH_OPTIONS[@]}" \
    "$INSTALLER" \
    "$CHECKER" \
    "$SSH_USER@$BUILDER_HOST:/tmp/"

printf '\nInstalling and verifying worker dependencies...\n'
ssh "${SSH_OPTIONS[@]}" \
    "$SSH_USER@$BUILDER_HOST" \
    'set -Eeuo pipefail
     sudo chmod +x \
         /tmp/install-worker-base.sh \
         /tmp/check-worker-ami.sh
     sudo /tmp/install-worker-base.sh

     # Keep this explicit so an older local installer cannot miss pigz.
     sudo apt-get update
     sudo DEBIAN_FRONTEND=noninteractive \
         apt-get install -y --no-install-recommends pigz

     command -v pigz
     pigz --version
     test -f "$(gcc -print-file-name=crti.o)"
     sudo /tmp/check-worker-ami.sh
     df -hT /
     sudo rm -f \
         /tmp/install-worker-base.sh \
         /tmp/check-worker-ami.sh
     sudo apt-get clean
     sudo rm -rf /var/lib/apt/lists/*'

AMI_NAME="ukama-lab-p0-worker-$(date -u +%Y%m%dt%H%M%Sz)"

printf '\nCapturing AMI %s...\n' "$AMI_NAME"
"$CAPTURE" "$BUILDER_ID" "$AMI_NAME"
CAPTURE_COMPLETE=1

NEW_AMI="$(
    awk -F= '$1 == "AMI_ID" { value=$2 } END { print value }' \
        "$P0_AWS_STATE"
)"

[[ "$NEW_AMI" == ami-* ]] ||
    p0_die "capture completed but no valid AMI_ID was written to $P0_AWS_STATE"

printf '\n========================================\n'
printf 'New worker AMI: %s\n' "$NEW_AMI"
printf 'Saved in:       %s\n' "$P0_AWS_STATE"
printf 'Builder:        %s (will now be terminated)\n' "$BUILDER_ID"
printf '========================================\n'
printf '\nRun all P0 scenarios with:\n'
printf '  ./utils/aws-p0/run.sh --workers 20\n'
