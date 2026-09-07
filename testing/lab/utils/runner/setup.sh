#!/usr/bin/env bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

# One-time AWS setup for the plain EC2/S3 P0 runner.

set -Eeuo pipefail

SCRIPT_DIR="$(CDPATH= cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=lib.sh
. "$SCRIPT_DIR/lib.sh"
p0_load_config "$SCRIPT_DIR"

CREDENTIALS_FILE="${P0_AWS_CREDENTIALS:-$SCRIPT_DIR/credentials.env}"
[[ -r "$CREDENTIALS_FILE" ]] ||
    p0_die "credentials file not found: $CREDENTIALS_FILE (copy credentials.env.example)"

p0_require_cmd aws
p0_require_cmd jq
p0_require_config AWS_REGION S3_BUCKET S3_PREFIX SECRET_ID \
    IAM_ROLE_NAME INSTANCE_PROFILE_NAME SECURITY_GROUP_NAME

p0_aws sts get-caller-identity >/dev/null
ACCOUNT_ID="$(p0_aws sts get-caller-identity --query Account --output text)"

printf 'AWS account: %s\n' "$ACCOUNT_ID"
printf 'AWS region:  %s\n' "$AWS_REGION"

if ! p0_aws s3api head-bucket --bucket "$S3_BUCKET" >/dev/null 2>&1; then
    printf 'creating S3 bucket %s\n' "$S3_BUCKET"
    if [[ "$AWS_REGION" == "us-east-1" ]]; then
        p0_aws s3api create-bucket --bucket "$S3_BUCKET" >/dev/null
    else
        p0_aws s3api create-bucket \
            --bucket "$S3_BUCKET" \
            --create-bucket-configuration "LocationConstraint=$AWS_REGION" \
            >/dev/null
    fi
fi

p0_aws s3api put-public-access-block \
    --bucket "$S3_BUCKET" \
    --public-access-block-configuration \
        'BlockPublicAcls=true,IgnorePublicAcls=true,BlockPublicPolicy=true,RestrictPublicBuckets=true'

p0_aws s3api put-bucket-encryption \
    --bucket "$S3_BUCKET" \
    --server-side-encryption-configuration \
        '{"Rules":[{"ApplyServerSideEncryptionByDefault":{"SSEAlgorithm":"AES256"},"BucketKeyEnabled":true}]}'

if [[ "${S3_RETENTION_DAYS:-0}" =~ ^[1-9][0-9]*$ ]]; then
    lifecycle_file="$(mktemp)"
    jq -n \
        --arg prefix "${S3_PREFIX%/}/" \
        --argjson days "$S3_RETENTION_DAYS" \
        '{Rules:[{ID:"ExpireUkamaP0Batches",Status:"Enabled",Filter:{Prefix:$prefix},Expiration:{Days:$days}}]}' \
        >"$lifecycle_file"
    p0_aws s3api put-bucket-lifecycle-configuration \
        --bucket "$S3_BUCKET" \
        --lifecycle-configuration "file://$lifecycle_file"
    rm -f "$lifecycle_file"
fi

# The credentials file is trusted local input. Export it only long enough to
# create the secret JSON, then clear the password from the shell environment.
set -a
# shellcheck disable=SC1090
. "$CREDENTIALS_FILE"
set +a

: "${UKAMA_IDENTIFIER:?set UKAMA_IDENTIFIER in credentials.env}"
: "${UKAMA_PASSWORD:?set UKAMA_PASSWORD in credentials.env}"
: "${PAUTH_URL:?set PAUTH_URL in credentials.env}"
: "${BFF_BASE_URL:?set BFF_BASE_URL in credentials.env}"

secret_file="$(mktemp)"
jq -n \
    --arg UKAMA_IDENTIFIER "$UKAMA_IDENTIFIER" \
    --arg UKAMA_PASSWORD "$UKAMA_PASSWORD" \
    --arg PAUTH_URL "$PAUTH_URL" \
    --arg BFF_BASE_URL "$BFF_BASE_URL" \
    --arg ULAB_UKAMA_AGENT_NODE_GW_URL "${ULAB_UKAMA_AGENT_NODE_GW_URL:-}" \
    --arg ULAB_UKAMA_AGENT_API_GW_URL "${ULAB_UKAMA_AGENT_API_GW_URL:-}" \
    --arg ULAB_PAYMENT_FAILURE_ON_CMD "${ULAB_PAYMENT_FAILURE_ON_CMD:-}" \
    --arg ULAB_PAYMENT_FAILURE_OFF_CMD "${ULAB_PAYMENT_FAILURE_OFF_CMD:-}" \
    --arg ULAB_SOFTWARE_FAILURE_ON_CMD "${ULAB_SOFTWARE_FAILURE_ON_CMD:-}" \
    --arg ULAB_SOFTWARE_FAILURE_OFF_CMD "${ULAB_SOFTWARE_FAILURE_OFF_CMD:-}" \
    '$ARGS.named' >"$secret_file"

if p0_aws secretsmanager describe-secret \
    --secret-id "$SECRET_ID" >/dev/null 2>&1; then
    printf 'updating Secrets Manager secret %s\n' "$SECRET_ID"
    p0_aws secretsmanager put-secret-value \
        --secret-id "$SECRET_ID" \
        --secret-string "file://$secret_file" >/dev/null
else
    printf 'creating Secrets Manager secret %s\n' "$SECRET_ID"
    p0_aws secretsmanager create-secret \
        --name "$SECRET_ID" \
        --secret-string "file://$secret_file" >/dev/null
fi
rm -f "$secret_file"
unset UKAMA_PASSWORD

SECRET_ARN="$(p0_aws secretsmanager describe-secret \
    --secret-id "$SECRET_ID" --query ARN --output text)"

trust_file="$(mktemp)"
cat >"$trust_file" <<'JSON'
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {"Service": "ec2.amazonaws.com"},
      "Action": "sts:AssumeRole"
    }
  ]
}
JSON

if ! p0_aws iam get-role --role-name "$IAM_ROLE_NAME" >/dev/null 2>&1; then
    printf 'creating IAM role %s\n' "$IAM_ROLE_NAME"
    p0_aws iam create-role \
        --role-name "$IAM_ROLE_NAME" \
        --assume-role-policy-document "file://$trust_file" >/dev/null
else
    p0_aws iam update-assume-role-policy \
        --role-name "$IAM_ROLE_NAME" \
        --policy-document "file://$trust_file"
fi
rm -f "$trust_file"

policy_file="$(mktemp)"
jq -n \
    --arg bucket "$S3_BUCKET" \
    --arg prefix "${S3_PREFIX%/}" \
    --arg secret "$SECRET_ARN" \
    '{
      Version:"2012-10-17",
      Statement:[
        {
          Sid:"ListBatchObjects",
          Effect:"Allow",
          Action:["s3:ListBucket"],
          Resource:[("arn:aws:s3:::"+$bucket)],
          Condition:{StringLike:{"s3:prefix":[($prefix+"/*")]}}
        },
        {
          Sid:"ReadWriteBatchObjects",
          Effect:"Allow",
          Action:["s3:GetObject","s3:PutObject","s3:AbortMultipartUpload"],
          Resource:[("arn:aws:s3:::"+$bucket+"/"+$prefix+"/*")]
        },
        {
          Sid:"ReadUdevSecret",
          Effect:"Allow",
          Action:["secretsmanager:GetSecretValue"],
          Resource:[$secret]
        }
      ]
    }' >"$policy_file"

p0_aws iam put-role-policy \
    --role-name "$IAM_ROLE_NAME" \
    --policy-name UkamaP0WorkerAccess \
    --policy-document "file://$policy_file"
rm -f "$policy_file"

if ! p0_aws iam get-instance-profile \
    --instance-profile-name "$INSTANCE_PROFILE_NAME" >/dev/null 2>&1; then
    printf 'creating instance profile %s\n' "$INSTANCE_PROFILE_NAME"
    p0_aws iam create-instance-profile \
        --instance-profile-name "$INSTANCE_PROFILE_NAME" >/dev/null
fi

profile_roles="$(p0_aws iam get-instance-profile \
    --instance-profile-name "$INSTANCE_PROFILE_NAME" \
    --query 'InstanceProfile.Roles[].RoleName' --output text)"
if ! grep -qw -- "$IAM_ROLE_NAME" <<<"$profile_roles"; then
    p0_aws iam add-role-to-instance-profile \
        --instance-profile-name "$INSTANCE_PROFILE_NAME" \
        --role-name "$IAM_ROLE_NAME"
fi

if [[ -z "${VPC_ID:-}" ]]; then
    VPC_ID="$(p0_aws ec2 describe-vpcs \
        --filters Name=is-default,Values=true \
        --query 'Vpcs[0].VpcId' --output text)"
fi
[[ -n "$VPC_ID" && "$VPC_ID" != "None" ]] ||
    p0_die 'VPC_ID is empty and no default VPC exists'

if [[ -z "${SUBNET_ID:-}" ]]; then
    SUBNET_ID="$(p0_aws ec2 describe-subnets \
        --filters "Name=vpc-id,Values=$VPC_ID" \
        --query 'sort_by(Subnets,&AvailabilityZone)[0].SubnetId' \
        --output text)"
fi
[[ -n "$SUBNET_ID" && "$SUBNET_ID" != "None" ]] ||
    p0_die 'SUBNET_ID is empty and no subnet was discovered'

if [[ -z "${SECURITY_GROUP_ID:-}" ]]; then
    SECURITY_GROUP_ID="$(p0_aws ec2 describe-security-groups \
        --filters \
            "Name=vpc-id,Values=$VPC_ID" \
            "Name=group-name,Values=$SECURITY_GROUP_NAME" \
        --query 'SecurityGroups[0].GroupId' --output text)"
fi

if [[ -z "$SECURITY_GROUP_ID" || "$SECURITY_GROUP_ID" == "None" ]]; then
    printf 'creating outbound-only security group %s\n' "$SECURITY_GROUP_NAME"
    SECURITY_GROUP_ID="$(p0_aws ec2 create-security-group \
        --vpc-id "$VPC_ID" \
        --group-name "$SECURITY_GROUP_NAME" \
        --description 'Disposable Ukama P0 workers; no inbound access' \
        --query GroupId --output text)"
    p0_aws ec2 create-tags \
        --resources "$SECURITY_GROUP_ID" \
        --tags Key=Name,Value="$SECURITY_GROUP_NAME" Key=Application,Value=ukama-lab-p0
fi

cat >"$P0_AWS_STATE" <<EOF_STATE
# Generated by setup.sh. Safe to recreate.
ACCOUNT_ID=$ACCOUNT_ID
VPC_ID=$VPC_ID
SUBNET_ID=$SUBNET_ID
SECURITY_GROUP_ID=$SECURITY_GROUP_ID
SECRET_ARN=$SECRET_ARN
EOF_STATE
if [[ -n "${AMI_ID:-}" && "$AMI_ID" != "ami-REPLACE_ME" ]]; then
    printf 'AMI_ID=%s\n' "$AMI_ID" >>"$P0_AWS_STATE"
fi
chmod 600 "$P0_AWS_STATE"

printf '\nAWS setup complete.\n'
printf 'state:          %s\n' "$P0_AWS_STATE"
printf 'bucket:         s3://%s/%s/\n' "$S3_BUCKET" "${S3_PREFIX%/}"
printf 'instance role:  %s\n' "$IAM_ROLE_NAME"
printf 'instance profile: %s\n' "$INSTANCE_PROFILE_NAME"
printf 'subnet:         %s\n' "$SUBNET_ID"
printf 'security group: %s\n' "$SECURITY_GROUP_ID"
printf '\nSet AMI_ID in config.env before running a batch.\n'
