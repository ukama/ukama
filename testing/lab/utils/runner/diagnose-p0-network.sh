#!/usr/bin/env bash
# Collect local, AWS VPC, and disposable-worker network diagnostics for P0.

set -Eeuo pipefail
umask 077

SCRIPT_DIR="$(CDPATH= cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=lib.sh
. "$SCRIPT_DIR/lib.sh"
p0_load_config "$SCRIPT_DIR"

usage() {
    cat <<'USAGE'
usage: diagnose-p0-network.sh [options]

Options:
  --no-launch            Collect local and AWS configuration only.
  --keep-instance        Do not terminate the diagnostic EC2 instance.
  --instance-type TYPE   Diagnostic instance type (default: t3.micro).
  --timeout SECONDS      Wait limit for remote report (default: 420).
  --output DIR           Local report directory.
  -h, --help             Show this help.

The script never prints or uploads UKAMA_PASSWORD or AWS access keys.
USAGE
}

NO_LAUNCH=0
KEEP_INSTANCE=0
DIAG_INSTANCE_TYPE="${DIAGNOSTIC_INSTANCE_TYPE:-t3.micro}"
WAIT_SECONDS="${DIAGNOSTIC_TIMEOUT_SECONDS:-420}"
OUTPUT_DIR=""

while (($#)); do
    case "$1" in
        --no-launch)
            NO_LAUNCH=1
            shift
            ;;
        --keep-instance)
            KEEP_INSTANCE=1
            shift
            ;;
        --instance-type)
            (($# >= 2)) || p0_die '--instance-type requires a value'
            DIAG_INSTANCE_TYPE="$2"
            shift 2
            ;;
        --instance-type=*)
            DIAG_INSTANCE_TYPE="${1#*=}"
            shift
            ;;
        --timeout)
            (($# >= 2)) || p0_die '--timeout requires seconds'
            WAIT_SECONDS="$2"
            shift 2
            ;;
        --timeout=*)
            WAIT_SECONDS="${1#*=}"
            shift
            ;;
        --output)
            (($# >= 2)) || p0_die '--output requires a directory'
            OUTPUT_DIR="$2"
            shift 2
            ;;
        --output=*)
            OUTPUT_DIR="${1#*=}"
            shift
            ;;
        -h|--help)
            usage
            exit 0
            ;;
        *)
            p0_die "unknown option: $1"
            ;;
    esac
done

[[ "$WAIT_SECONDS" =~ ^[1-9][0-9]*$ ]] ||
    p0_die '--timeout must be a positive integer'

for cmd in aws jq curl sed awk grep sort tar sha256sum; do
    p0_require_cmd "$cmd"
done

p0_require_config AWS_REGION S3_BUCKET S3_PREFIX AMI_ID SUBNET_ID \
    SECURITY_GROUP_ID INSTANCE_PROFILE_NAME

CREDENTIALS_FILE="${P0_AWS_CREDENTIALS:-$SCRIPT_DIR/credentials.env}"
[[ -r "$CREDENTIALS_FILE" ]] ||
    p0_die "credentials file not found: $CREDENTIALS_FILE"

# Load endpoint values only. The password is never printed or written out.
set -a
# shellcheck disable=SC1090
. "$CREDENTIALS_FILE"
set +a

: "${PAUTH_URL:?set PAUTH_URL in credentials.env}"
: "${BFF_BASE_URL:?set BFF_BASE_URL in credentials.env}"

DIAG_ID="network-$(date -u +%Y%m%dt%H%M%Sz)"
if [[ -z "$OUTPUT_DIR" ]]; then
    LAB_ROOT="$(git -C "$SCRIPT_DIR" rev-parse --show-toplevel 2>/dev/null ||
        CDPATH= cd -- "$SCRIPT_DIR/../.." && pwd)"
    OUTPUT_DIR="$LAB_ROOT/runs/p0-network/$DIAG_ID"
fi
mkdir -p "$OUTPUT_DIR"

LOCAL_REPORT="$OUTPUT_DIR/local.txt"
AWS_REPORT="$OUTPUT_DIR/aws.txt"
REMOTE_REPORT="$OUTPUT_DIR/remote.txt"
SUMMARY="$OUTPUT_DIR/summary.txt"
USER_DATA="$OUTPUT_DIR/user-data.sh"
ENDPOINTS_FILE="$OUTPUT_DIR/endpoints.env"
INSTANCE_FILE="$OUTPUT_DIR/instance-id"
S3_ROOT="s3://${S3_BUCKET}/${S3_PREFIX%/}/diagnostics/${DIAG_ID}"
INSTANCE_ID=""

redact_url() {
    # Keep scheme, host, port, and path. Remove URL user-info and query values.
    printf '%s' "$1" |
        sed -E 's#(https?://)[^/@]+@#\1#; s#\?.*$##'
}

add_endpoint() {
    local name="$1"
    local url="${2:-}"

    [[ -n "$url" ]] || return 0
    printf '%s\t%s\n' "$name" "$(redact_url "$url")" >>"$ENDPOINTS_FILE"
}

: >"$ENDPOINTS_FILE"
add_endpoint PAUTH "$PAUTH_URL"
add_endpoint BFF "${BFF_BASE_URL%/}/gateway/graphql"
add_endpoint WAREHOUSE "${UKAMA_LAB_WAREHOUSE_URL:-}"
add_endpoint FACTORY "${UKAMA_LAB_FACTORY_URL:-}"
add_endpoint FACTORY_SEED "${ULAB_FACTORY_SEED_URL:-}"
add_endpoint NODE_GATEWAY "${ULAB_UKAMA_AGENT_NODE_GW_URL:-}"
add_endpoint API_GATEWAY "${ULAB_UKAMA_AGENT_API_GW_URL:-}"

# Remove the password from this process as soon as the endpoints are loaded.
unset UKAMA_PASSWORD

section() {
    printf '\n===== %s =====\n' "$1"
}

host_from_url() {
    local value="$1"

    value="${value#*://}"
    value="${value%%/*}"
    value="${value%%:*}"
    printf '%s' "$value"
}

scheme_from_url() {
    printf '%s' "${1%%://*}"
}

port_from_url() {
    local value="$1"
    local authority
    local scheme

    scheme="$(scheme_from_url "$value")"
    authority="${value#*://}"
    authority="${authority%%/*}"
    if [[ "$authority" == *:* ]]; then
        printf '%s' "${authority##*:}"
    elif [[ "$scheme" == "https" ]]; then
        printf '443'
    else
        printf '80'
    fi
}

probe_url() {
    local name="$1"
    local url="$2"
    local host
    local port
    local scheme
    local result
    local rc

    host="$(host_from_url "$url")"
    port="$(port_from_url "$url")"
    scheme="$(scheme_from_url "$url")"
    printf '\n[%s]\nurl=%s\nhost=%s\n' "$name" "$url" "$host"

    if command -v getent >/dev/null 2>&1; then
        printf '%s\n' '-- DNS (getent ahosts) --'
        getent ahosts "$host" 2>&1 || true
    fi
    if command -v dig >/dev/null 2>&1; then
        printf '%s\n' '-- DNS (dig +short) --'
        dig +time=3 +tries=1 +short "$host" A 2>&1 || true
    fi

    printf '%s\n' '-- Route to resolved IPv4 addresses --'
    if command -v ip >/dev/null 2>&1; then
        { getent ahostsv4 "$host" 2>/dev/null || true; } |
            awk '{print $1}' | sort -u |
            while IFS= read -r address; do
                [[ -n "$address" ]] || continue
                printf 'address=%s\n' "$address"
                ip route get "$address" 2>&1 || true
            done
    fi

    printf '%s\n' '-- HTTPS probe --'
    set +e
    result="$(curl -sS -L -o /dev/null \
        --connect-timeout 8 --max-time 20 \
        -w 'http_code=%{http_code}\nremote_ip=%{remote_ip}\nlocal_ip=%{local_ip}\ntime_namelookup=%{time_namelookup}\ntime_connect=%{time_connect}\ntime_appconnect=%{time_appconnect}\ntime_total=%{time_total}\n' \
        "$url" 2>&1)"
    rc=$?
    set -e
    printf 'curl_exit=%s\n%s\n' "$rc" "$result"

    if [[ "$scheme" == "https" ]] &&
       command -v openssl >/dev/null 2>&1; then
        printf '%s\n' '-- TLS handshake --'
        timeout 12 openssl s_client \
            -connect "${host}:${port}" -servername "$host" -brief \
            </dev/null 2>&1 | sed -n '1,30p' || true
    fi
}

collect_local() {
    {
        printf 'diagnostic_id=%s\n' "$DIAG_ID"
        printf 'collected_at=%s\n' "$(date -u +%Y-%m-%dT%H:%M:%SZ)"
        printf 'hostname=%s\n' "$(hostname 2>/dev/null || true)"
        printf 'aws_region=%s\n' "$AWS_REGION"

        section 'LOCAL AWS IDENTITY'
        p0_aws sts get-caller-identity --output json 2>&1 || true

        section 'LOCAL ADDRESSES'
        ip -brief address 2>&1 || true

        section 'LOCAL ROUTES'
        ip route show table all 2>&1 || true
        ip rule show 2>&1 || true

        section 'LOCAL DNS'
        cat /etc/resolv.conf 2>&1 || true
        if command -v resolvectl >/dev/null 2>&1; then
            resolvectl status 2>&1 || true
        fi

        section 'ACTIVE VPN / NETWORK CONNECTIONS'
        if command -v nmcli >/dev/null 2>&1; then
            nmcli -t -f NAME,TYPE,DEVICE connection show --active 2>&1 || true
        fi
        ip -brief link 2>&1 || true

        section 'LOCAL PUBLIC IP'
        curl -sS --connect-timeout 5 --max-time 10 \
            https://checkip.amazonaws.com 2>&1 || true

        section 'LOCAL ENDPOINT PROBES'
        while IFS=$'\t' read -r name url; do
            probe_url "$name" "$url"
        done <"$ENDPOINTS_FILE"
    } >"$LOCAL_REPORT" 2>&1
}

aws_try() {
    local title="$1"
    shift

    section "$title"
    p0_aws "$@" --output json 2>&1 || true
}

collect_aws() {
    local vpc_id
    local main_route_table

    vpc_id="$(p0_aws ec2 describe-subnets \
        --subnet-ids "$SUBNET_ID" \
        --query 'Subnets[0].VpcId' --output text)"

    {
        printf 'diagnostic_id=%s\n' "$DIAG_ID"
        printf 'collected_at=%s\n' "$(date -u +%Y-%m-%dT%H:%M:%SZ)"
        printf 'aws_region=%s\n' "$AWS_REGION"
        printf 'ami_id=%s\n' "$AMI_ID"
        printf 'subnet_id=%s\n' "$SUBNET_ID"
        printf 'vpc_id=%s\n' "$vpc_id"
        printf 'security_group_id=%s\n' "$SECURITY_GROUP_ID"
        printf 'associate_public_ip=%s\n' "${ASSOCIATE_PUBLIC_IP:-true}"

        aws_try 'CALLER IDENTITY' sts get-caller-identity
        aws_try 'AMI' ec2 describe-images --image-ids "$AMI_ID"
        aws_try 'VPC' ec2 describe-vpcs --vpc-ids "$vpc_id"
        aws_try 'VPC DNS SUPPORT' ec2 describe-vpc-attribute \
            --vpc-id "$vpc_id" --attribute enableDnsSupport
        aws_try 'VPC DNS HOSTNAMES' ec2 describe-vpc-attribute \
            --vpc-id "$vpc_id" --attribute enableDnsHostnames
        aws_try 'SUBNET' ec2 describe-subnets --subnet-ids "$SUBNET_ID"
        aws_try 'EXPLICIT SUBNET ROUTE TABLE' ec2 describe-route-tables \
            --filters "Name=association.subnet-id,Values=$SUBNET_ID"

        main_route_table="$(p0_aws ec2 describe-route-tables \
            --filters \
                "Name=vpc-id,Values=$vpc_id" \
                'Name=association.main,Values=true' \
            --query 'RouteTables[0].RouteTableId' \
            --output text 2>/dev/null || true)"
        if [[ -n "$main_route_table" && "$main_route_table" != "None" ]]; then
            aws_try 'MAIN VPC ROUTE TABLE' ec2 describe-route-tables \
                --route-table-ids "$main_route_table"
        fi

        aws_try 'NETWORK ACL' ec2 describe-network-acls \
            --filters "Name=association.subnet-id,Values=$SUBNET_ID"
        aws_try 'SECURITY GROUP' ec2 describe-security-groups \
            --group-ids "$SECURITY_GROUP_ID"
        aws_try 'INTERNET GATEWAYS' ec2 describe-internet-gateways \
            --filters "Name=attachment.vpc-id,Values=$vpc_id"
        aws_try 'NAT GATEWAYS' ec2 describe-nat-gateways \
            --filter "Name=vpc-id,Values=$vpc_id"
        aws_try 'VIRTUAL PRIVATE GATEWAYS' ec2 describe-vpn-gateways \
            --filters "Name=attachment.vpc-id,Values=$vpc_id"
        aws_try 'TRANSIT GATEWAY ATTACHMENTS' \
            ec2 describe-transit-gateway-attachments \
            --filters "Name=resource-id,Values=$vpc_id"
        aws_try 'VPC PEERING CONNECTIONS' \
            ec2 describe-vpc-peering-connections \
            --filters \
                "Name=requester-vpc-info.vpc-id,Values=$vpc_id" \
                "Name=accepter-vpc-info.vpc-id,Values=$vpc_id"
        aws_try 'VPC ENDPOINTS' ec2 describe-vpc-endpoints \
            --filters "Name=vpc-id,Values=$vpc_id"

        section 'PRIVATE HOSTED ZONES ASSOCIATED WITH VPC'
        p0_aws route53 list-hosted-zones-by-vpc \
            --vpc-id "$vpc_id" --vpc-region "$AWS_REGION" \
            --output json 2>&1 || true
    } >"$AWS_REPORT" 2>&1
}

make_user_data() {
    local endpoint_lines

    endpoint_lines="$(sed "s/'/'\\''/g" "$ENDPOINTS_FILE")"
    cat >"$USER_DATA" <<EOF_USER_DATA
#!/usr/bin/env bash
set -uo pipefail
umask 077

REPORT=/tmp/ukama-p0-network-remote.txt
DONE=/tmp/ukama-p0-network-remote.done
S3_ROOT='$S3_ROOT'
AWS_REGION='$AWS_REGION'

export AWS_DEFAULT_REGION="\$AWS_REGION"
exec > >(tee -a "\$REPORT") 2>&1

section() {
    printf '\n===== %s =====\n' "\$1"
}

host_from_url() {
    local value="\$1"
    value="\${value#*://}"
    value="\${value%%/*}"
    value="\${value%%:*}"
    printf '%s' "\$value"
}

scheme_from_url() {
    printf '%s' "\${1%%://*}"
}

port_from_url() {
    local value="\$1"
    local authority
    local scheme

    scheme="\$(scheme_from_url "\$value")"
    authority="\${value#*://}"
    authority="\${authority%%/*}"
    if [[ "\$authority" == *:* ]]; then
        printf '%s' "\${authority##*:}"
    elif [[ "\$scheme" == "https" ]]; then
        printf '443'
    else
        printf '80'
    fi
}

probe_url() {
    local name="\$1"
    local url="\$2"
    local host
    local port
    local scheme
    local result
    local rc

    host="\$(host_from_url "\$url")"
    port="\$(port_from_url "\$url")"
    scheme="\$(scheme_from_url "\$url")"
    printf '\n[%s]\nurl=%s\nhost=%s\n' "\$name" "\$url" "\$host"

    printf '%s\n' '-- DNS (getent ahosts) --'
    getent ahosts "\$host" 2>&1 || true
    if command -v dig >/dev/null 2>&1; then
        printf '%s\n' '-- DNS (dig +short) --'
        dig +time=3 +tries=1 +short "\$host" A 2>&1 || true
    fi

    printf '%s\n' '-- Route to resolved IPv4 addresses --'
    { getent ahostsv4 "\$host" 2>/dev/null || true; } |
        awk '{print \$1}' | sort -u |
        while IFS= read -r address; do
            [[ -n "\$address" ]] || continue
            printf 'address=%s\n' "\$address"
            ip route get "\$address" 2>&1 || true
        done

    printf '%s\n' '-- HTTPS probe --'
    set +e
    result="\$(curl -sS -L -o /dev/null \\
        --connect-timeout 8 --max-time 20 \\
        -w 'http_code=%{http_code}\\nremote_ip=%{remote_ip}\\nlocal_ip=%{local_ip}\\ntime_namelookup=%{time_namelookup}\\ntime_connect=%{time_connect}\\ntime_appconnect=%{time_appconnect}\\ntime_total=%{time_total}\\n' \\
        "\$url" 2>&1)"
    rc=\$?
    set -e
    printf 'curl_exit=%s\n%s\n' "\$rc" "\$result"

    if [[ "\$scheme" == "https" ]] &&
       command -v openssl >/dev/null 2>&1; then
        printf '%s\n' '-- TLS handshake --'
        timeout 12 openssl s_client \\
            -connect "\${host}:\${port}" -servername "\$host" -brief \\
            </dev/null 2>&1 | sed -n '1,30p' || true
    fi
}

printf 'diagnostic_id=%s\n' '$DIAG_ID'
printf 'collected_at=%s\n' "\$(date -u +%Y-%m-%dT%H:%M:%SZ)"
printf 'hostname=%s\n' "\$(hostname 2>/dev/null || true)"

section 'INSTANCE METADATA'
token="\$(curl -fsS -X PUT \\
    -H 'X-aws-ec2-metadata-token-ttl-seconds: 300' \\
    http://169.254.169.254/latest/api/token 2>/dev/null || true)"
if [[ -n "\$token" ]]; then
    curl -fsS -H "X-aws-ec2-metadata-token: \$token" \\
        http://169.254.169.254/latest/dynamic/instance-identity/document || true
    printf '\n'
fi

section 'AWS IDENTITY'
aws sts get-caller-identity --output json 2>&1 || true

section 'ADDRESSES'
ip -brief address 2>&1 || true

section 'ROUTES'
ip route show table all 2>&1 || true
ip rule show 2>&1 || true

section 'DNS'
cat /etc/resolv.conf 2>&1 || true
if command -v resolvectl >/dev/null 2>&1; then
    resolvectl status 2>&1 || true
fi

section 'PUBLIC IP'
curl -sS --connect-timeout 5 --max-time 10 \\
    https://checkip.amazonaws.com 2>&1 || true

section 'ENDPOINT PROBES'
while IFS=\$'\\t' read -r name url; do
    probe_url "\$name" "\$url"
done <<'EOF_ENDPOINTS'
$endpoint_lines
EOF_ENDPOINTS

section 'UPLOAD TEST'
printf 'done_at=%s\n' "\$(date -u +%Y-%m-%dT%H:%M:%SZ)" >"\$DONE"
aws s3 cp "\$REPORT" "\$S3_ROOT/remote.txt" --only-show-errors
aws s3 cp "\$DONE" "\$S3_ROOT/remote.done" --only-show-errors
sync
shutdown -P now
EOF_USER_DATA
    chmod 600 "$USER_DATA"
}

cleanup_instance() {
    local state

    [[ -n "$INSTANCE_ID" ]] || return 0
    ((KEEP_INSTANCE)) && return 0

    state="$(p0_aws ec2 describe-instances \
        --instance-ids "$INSTANCE_ID" \
        --query 'Reservations[0].Instances[0].State.Name' \
        --output text 2>/dev/null || true)"
    case "$state" in
        pending|running|stopping|stopped)
            printf 'terminating diagnostic instance %s\n' "$INSTANCE_ID" >&2
            p0_aws ec2 terminate-instances \
                --instance-ids "$INSTANCE_ID" >/dev/null || true
            ;;
    esac
}

launch_remote() {
    local associate
    local network_json
    local elapsed=0

    make_user_data
    associate="${ASSOCIATE_PUBLIC_IP:-true}"
    case "$associate" in
        true|TRUE|1|yes|YES)
            associate=true
            ;;
        *)
            associate=false
            ;;
    esac

    network_json="$(jq -cn \
        --arg subnet "$SUBNET_ID" \
        --arg sg "$SECURITY_GROUP_ID" \
        --argjson public "$associate" \
        '[{DeviceIndex:0,SubnetId:$subnet,Groups:[$sg],
           AssociatePublicIpAddress:$public,
           DeleteOnTermination:true}]')"

    printf 'launching diagnostic instance in subnet %s\n' "$SUBNET_ID"
    INSTANCE_ID="$(p0_aws ec2 run-instances \
        --image-id "$AMI_ID" \
        --instance-type "$DIAG_INSTANCE_TYPE" \
        --iam-instance-profile "Name=$INSTANCE_PROFILE_NAME" \
        --network-interfaces "$network_json" \
        --instance-initiated-shutdown-behavior terminate \
        --metadata-options \
            'HttpEndpoint=enabled,HttpTokens=required,HttpPutResponseHopLimit=2' \
        --user-data "file://$USER_DATA" \
        --tag-specifications \
            "ResourceType=instance,Tags=[
              {Key=Name,Value=ukama-p0-network-diagnostic},
              {Key=Application,Value=ukama-lab-p0},
              {Key=UkamaP0Diagnostic,Value=$DIAG_ID},
              {Key=Ephemeral,Value=true}]" \
        --query 'Instances[0].InstanceId' --output text)"
    printf '%s\n' "$INSTANCE_ID" >"$INSTANCE_FILE"
    printf 'instance=%s\n' "$INSTANCE_ID"

    while ((elapsed < WAIT_SECONDS)); do
        if p0_aws s3api head-object \
            --bucket "$S3_BUCKET" \
            --key "${S3_PREFIX%/}/diagnostics/${DIAG_ID}/remote.done" \
            >/dev/null 2>&1; then
            p0_aws s3 cp "$S3_ROOT/remote.txt" "$REMOTE_REPORT" \
                --only-show-errors
            return 0
        fi
        printf '\rwaiting for remote report: %ss/%ss' \
            "$elapsed" "$WAIT_SECONDS"
        sleep 5
        elapsed=$((elapsed + 5))
    done
    printf '\n'

    {
        printf 'remote diagnostic did not upload a report within %ss\n' \
            "$WAIT_SECONDS"
        p0_aws ec2 describe-instances --instance-ids "$INSTANCE_ID" \
            --output json 2>&1 || true
        p0_aws ec2 get-console-output --latest \
            --instance-id "$INSTANCE_ID" --output json 2>&1 || true
    } >"$REMOTE_REPORT"
    return 1
}

make_summary() {
    local local_pauth_ip
    local remote_pauth_ip
    local local_curl
    local remote_curl

    local_pauth_ip="$(awk '
        /^\[PAUTH\]/{in_pauth=1; next}
        /^\[/{if (in_pauth) exit}
        in_pauth && /^[0-9]/ {print $1; exit}
    ' "$LOCAL_REPORT" 2>/dev/null || true)"
    remote_pauth_ip="$(awk '
        /^\[PAUTH\]/{in_pauth=1; next}
        /^\[/{if (in_pauth) exit}
        in_pauth && /^[0-9]/ {print $1; exit}
    ' "$REMOTE_REPORT" 2>/dev/null || true)"
    local_curl="$(awk '
        /^\[PAUTH\]/{in_pauth=1; next}
        /^\[/{if (in_pauth) exit}
        in_pauth && /^curl_exit=/ {print; exit}
    ' "$LOCAL_REPORT" 2>/dev/null || true)"
    remote_curl="$(awk '
        /^\[PAUTH\]/{in_pauth=1; next}
        /^\[/{if (in_pauth) exit}
        in_pauth && /^curl_exit=/ {print; exit}
    ' "$REMOTE_REPORT" 2>/dev/null || true)"

    {
        printf 'P0 network diagnostic\n\n'
        printf 'diagnostic_id=%s\n' "$DIAG_ID"
        printf 'aws_region=%s\n' "$AWS_REGION"
        printf 'ami_id=%s\n' "$AMI_ID"
        printf 'subnet_id=%s\n' "$SUBNET_ID"
        printf 'security_group_id=%s\n' "$SECURITY_GROUP_ID"
        printf 'instance_id=%s\n' "${INSTANCE_ID:-not-launched}"
        printf '\nPAUTH comparison\n'
        printf 'local_dns_first_ip=%s\n' "${local_pauth_ip:-unknown}"
        printf 'remote_dns_first_ip=%s\n' "${remote_pauth_ip:-unknown}"
        printf 'local_%s\n' "${local_curl:-curl_exit=unknown}"
        printf 'remote_%s\n' "${remote_curl:-curl_exit=unknown}"
        printf '\nFiles\n'
        printf 'local=%s\n' "$LOCAL_REPORT"
        printf 'aws=%s\n' "$AWS_REPORT"
        printf 'remote=%s\n' "$REMOTE_REPORT"
        printf '\nInterpretation\n'
        if [[ -n "$local_pauth_ip" && -n "$remote_pauth_ip" &&
              "$local_pauth_ip" != "$remote_pauth_ip" ]]; then
            printf '%s\n' \
                'Local and EC2 DNS differ: investigate split DNS / Route 53 Resolver.'
        elif [[ "$local_curl" == 'curl_exit=0' &&
                "$remote_curl" != 'curl_exit=0' ]]; then
            printf '%s\n' \
                'Local reaches PAUTH but EC2 does not: inspect EC2 route/allow-list.'
        else
            printf '%s\n' \
                'Use local.txt, aws.txt, and remote.txt for the exact failure point.'
        fi
    } >"$SUMMARY"
}

trap cleanup_instance EXIT INT TERM

printf 'Collecting local endpoint and VPN information...\n'
collect_local
printf 'Collecting AWS VPC, subnet, route, ACL, and security-group data...\n'
collect_aws

remote_rc=0
if ((NO_LAUNCH)); then
    printf 'remote launch disabled (--no-launch)\n' >"$REMOTE_REPORT"
else
    launch_remote || remote_rc=$?
fi

make_summary

tar -C "$OUTPUT_DIR" -czf "$OUTPUT_DIR/diagnostic.tar.gz" \
    local.txt aws.txt remote.txt summary.txt endpoints.env \
    ${INSTANCE_ID:+instance-id} 2>/dev/null || true

printf '\nDiagnostic complete.\n'
printf 'summary: %s\n' "$SUMMARY"
printf 'bundle:  %s\n' "$OUTPUT_DIR/diagnostic.tar.gz"
printf '\n'
cat "$SUMMARY"

if ((remote_rc != 0)); then
    exit "$remote_rc"
fi
