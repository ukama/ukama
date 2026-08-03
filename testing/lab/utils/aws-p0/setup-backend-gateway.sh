#!/usr/bin/env bash
# Configure one known-good EC2 instance as a simple NAT gateway for P0 workers.
# The instance already has the private route/tunnel needed to reach udev.

set -Eeuo pipefail
umask 077

SCRIPT_DIR="$(CDPATH= cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=lib.sh
. "$SCRIPT_DIR/lib.sh"
p0_load_config "$SCRIPT_DIR"

QUIET=0
if [[ "${1:-}" == "--quiet" ]]; then
    QUIET=1
elif (($#)); then
    p0_die "usage: $0 [--quiet]"
fi

p0_require_cmd aws
p0_require_cmd jq
p0_require_cmd ssh
p0_require_config AWS_REGION BACKEND_GATEWAY_INSTANCE_ID SUBNET_ID \
    SECURITY_GROUP_ID BACKEND_ROUTE_CIDR BACKEND_TEST_IP
p0_validate_simple_id BACKEND_GATEWAY_INSTANCE_ID "$BACKEND_GATEWAY_INSTANCE_ID"
p0_validate_simple_id SECURITY_GROUP_ID "$SECURITY_GROUP_ID"

instance_json="$(p0_aws ec2 describe-instances \
    --instance-ids "$BACKEND_GATEWAY_INSTANCE_ID" \
    --query 'Reservations[0].Instances[0]' \
    --output json)"

instance_state="$(jq -r '.State.Name // empty' <<<"$instance_json")"
if [[ "$instance_state" == "stopped" && \
      "${BACKEND_GATEWAY_START_IF_STOPPED:-true}" == "true" ]]; then
    ((QUIET)) || printf 'starting backend gateway %s\n' "$BACKEND_GATEWAY_INSTANCE_ID"
    p0_aws ec2 start-instances \
        --instance-ids "$BACKEND_GATEWAY_INSTANCE_ID" >/dev/null
    p0_aws ec2 wait instance-running \
        --instance-ids "$BACKEND_GATEWAY_INSTANCE_ID"
    instance_json="$(p0_aws ec2 describe-instances \
        --instance-ids "$BACKEND_GATEWAY_INSTANCE_ID" \
        --query 'Reservations[0].Instances[0]' \
        --output json)"
    instance_state="$(jq -r '.State.Name // empty' <<<"$instance_json")"
fi

[[ "$instance_state" == "running" ]] ||
    p0_die "backend gateway is not running: $BACKEND_GATEWAY_INSTANCE_ID ($instance_state)"

gateway_private_ip="$(jq -r '.PrivateIpAddress // empty' <<<"$instance_json")"
gateway_public_ip="$(jq -r '.PublicIpAddress // empty' <<<"$instance_json")"
gateway_vpc="$(jq -r '.VpcId // empty' <<<"$instance_json")"
gateway_subnet="$(jq -r '.SubnetId // empty' <<<"$instance_json")"
gateway_key_name="$(jq -r '.KeyName // empty' <<<"$instance_json")"
mapfile -t gateway_sgs < <(jq -r '.SecurityGroups[].GroupId' <<<"$instance_json")

[[ -n "$gateway_private_ip" ]] || p0_die 'backend gateway has no private IP'
((${#gateway_sgs[@]} > 0)) || p0_die 'backend gateway has no security group'

worker_subnet_json="$(p0_aws ec2 describe-subnets \
    --subnet-ids "$SUBNET_ID" \
    --query 'Subnets[0]' --output json)"
worker_vpc="$(jq -r '.VpcId // empty' <<<"$worker_subnet_json")"
worker_cidr="$(jq -r '.CidrBlock // empty' <<<"$worker_subnet_json")"

[[ "$gateway_vpc" == "$worker_vpc" ]] ||
    p0_die "gateway VPC $gateway_vpc differs from worker VPC $worker_vpc"
[[ -n "$worker_cidr" ]] || p0_die "cannot resolve CIDR for worker subnet $SUBNET_ID"

# A NAT instance must be allowed to forward packets not addressed to itself.
p0_aws ec2 modify-instance-attribute \
    --instance-id "$BACKEND_GATEWAY_INSTANCE_ID" \
    --no-source-dest-check

# Allow packets from the P0 worker security group into the gateway ENI.
gateway_sg="${BACKEND_GATEWAY_SECURITY_GROUP_ID:-${gateway_sgs[0]}}"
ip_permissions="$(jq -cn --arg worker_sg "$SECURITY_GROUP_ID" \
    '[{IpProtocol:"-1",UserIdGroupPairs:[{GroupId:$worker_sg,Description:"Ukama P0 worker forwarding"}]}]')"
set +e
ingress_error="$(p0_aws ec2 authorize-security-group-ingress \
    --group-id "$gateway_sg" \
    --ip-permissions "$ip_permissions" 2>&1)"
ingress_rc=$?
set -e
if ((ingress_rc != 0)) && ! grep -q 'InvalidPermission.Duplicate' <<<"$ingress_error"; then
    printf '%s\n' "$ingress_error" >&2
    p0_die "failed to allow P0 workers into gateway security group $gateway_sg"
fi

ssh_user="${BACKEND_GATEWAY_SSH_USER:-ubuntu}"
ssh_host="${BACKEND_GATEWAY_SSH_HOST:-}"
if [[ -z "$ssh_host" ]]; then
    if [[ "${BACKEND_GATEWAY_PREFER_PRIVATE_IP:-false}" == "true" ]]; then
        ssh_host="$gateway_private_ip"
    else
        ssh_host="${gateway_public_ip:-$gateway_private_ip}"
    fi
fi

ssh_key="${BACKEND_GATEWAY_SSH_KEY:-}"
if [[ -z "$ssh_key" && -n "$gateway_key_name" ]]; then
    for candidate in \
        "$HOME/.ssh/$gateway_key_name" \
        "$HOME/.ssh/$gateway_key_name.pem" \
        "$HOME/Downloads/$gateway_key_name.pem" \
        "$HOME/$gateway_key_name.pem"; do
        if [[ -r "$candidate" ]]; then
            ssh_key="$candidate"
            break
        fi
    done
fi

ssh_args=(
    -o BatchMode=yes
    -o ConnectTimeout="${BACKEND_GATEWAY_SSH_CONNECT_TIMEOUT:-20}"
    -o ServerAliveInterval=15
    -o ServerAliveCountMax=2
    -o StrictHostKeyChecking=accept-new
)
if [[ -n "$ssh_key" ]]; then
    [[ -r "$ssh_key" ]] ||
        p0_die "BACKEND_GATEWAY_SSH_KEY is not readable: $ssh_key"
    ssh_args+=(-i "$ssh_key")
fi

ssh_rule_added=0
ssh_cidr=""
remove_temporary_ssh_rule() {
    if ((ssh_rule_added)); then
        p0_aws ec2 revoke-security-group-ingress \
            --group-id "$gateway_sg" \
            --protocol tcp --port 22 --cidr "$ssh_cidr" >/dev/null 2>&1 || true
    fi
}
trap remove_temporary_ssh_rule EXIT INT TERM

if [[ "${BACKEND_GATEWAY_MANAGE_SSH_INGRESS:-true}" == "true" && \
      "$ssh_host" == "$gateway_public_ip" && -n "$gateway_public_ip" ]]; then
    controller_ip="$(curl -fsS --max-time 10 https://checkip.amazonaws.com 2>/dev/null | tr -d '[:space:]' || true)"
    if [[ "$controller_ip" =~ ^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
        ssh_cidr="$controller_ip/32"
        set +e
        ssh_error="$(p0_aws ec2 authorize-security-group-ingress \
            --group-id "$gateway_sg" --protocol tcp --port 22 \
            --cidr "$ssh_cidr" 2>&1)"
        ssh_rc=$?
        set -e
        if ((ssh_rc == 0)); then
            ssh_rule_added=1
        elif ! grep -q 'InvalidPermission.Duplicate' <<<"$ssh_error"; then
            printf '%s\n' "$ssh_error" >&2
            p0_die "failed to permit temporary SSH from $ssh_cidr"
        fi
    fi
fi

((QUIET)) || {
    printf 'backend gateway: %s (%s)\n' "$BACKEND_GATEWAY_INSTANCE_ID" "$gateway_private_ip"
    printf 'worker subnet:   %s (%s)\n' "$SUBNET_ID" "$worker_cidr"
    printf 'backend route:   %s via known-good host route\n' "$BACKEND_ROUTE_CIDR"
    printf 'ssh target:      %s@%s' "$ssh_user" "$ssh_host"
    [[ -z "$gateway_key_name" ]] || printf ' (EC2 key name: %s)' "$gateway_key_name"
    printf '\n'
}

remote_script="$(mktemp)"
cleanup_gateway_setup() {
    rm -f "$remote_script"
    remove_temporary_ssh_rule
}
trap cleanup_gateway_setup EXIT INT TERM
cat >"$remote_script" <<'REMOTE_EOF'
#!/usr/bin/env bash
set -Eeuo pipefail

export PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin

worker_cidr="$1"
backend_cidr="$2"
test_ip="$3"

if ((EUID != 0)); then
    echo 'gateway setup must run as root' >&2
    exit 2
fi

failure_report() {
    local rc=$?

    echo >&2
    echo "ukama-p0 gateway setup failed (exit=$rc)" >&2
    echo "route to backend test IP:" >&2
    ip -4 route get "$test_ip" >&2 || true
    echo "policy rules:" >&2
    ip -4 rule show >&2 || true
    echo "forwarding sysctl:" >&2
    sysctl net.ipv4.ip_forward >&2 || true
    echo "FORWARD rules:" >&2
    iptables -w -S FORWARD >&2 || true
    echo "POSTROUTING rules:" >&2
    iptables -w -t nat -S POSTROUTING >&2 || true
    exit "$rc"
}
trap failure_report ERR

if ! command -v iptables >/dev/null 2>&1; then
    if command -v apt-get >/dev/null 2>&1; then
        apt-get update
        DEBIAN_FRONTEND=noninteractive apt-get install -y iptables
    elif command -v dnf >/dev/null 2>&1; then
        dnf install -y iptables
    elif command -v yum >/dev/null 2>&1; then
        yum install -y iptables
    else
        echo 'iptables is required on the backend gateway' >&2
        exit 2
    fi
fi

cat >/etc/ukama-p0-gateway.env <<CONFIG_EOF
WORKER_CIDR=$(printf '%q' "$worker_cidr")
BACKEND_CIDR=$(printf '%q' "$backend_cidr")
TEST_IP=$(printf '%q' "$test_ip")
CONFIG_EOF
chmod 600 /etc/ukama-p0-gateway.env

cat >/usr/local/sbin/ukama-p0-gateway-apply <<'GATEWAY_EOF'
#!/usr/bin/env bash
set -Eeuo pipefail
export PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin

# shellcheck disable=SC1091
. /etc/ukama-p0-gateway.env

ROUTE_LINE="$(ip -4 route get "$TEST_IP")"
OUT_IF="$(awk '{for (i=1;i<=NF;i++) if ($i=="dev") {print $(i+1); exit}}' <<<"$ROUTE_LINE")"
ROUTE_TABLE="$(awk '{for (i=1;i<=NF;i++) if ($i=="table") {print $(i+1); exit}}' <<<"$ROUTE_LINE")"

[[ -n "$OUT_IF" ]] || {
    echo "cannot determine output interface for $TEST_IP: $ROUTE_LINE" >&2
    exit 1
}

# Verify that this host can reach udev before making it a gateway.
if command -v timeout >/dev/null 2>&1; then
    timeout 15 bash -c "</dev/tcp/$TEST_IP/443"
else
    bash -c "</dev/tcp/$TEST_IP/443"
fi

# If the tunnel uses a non-main policy table, forwarded traffic must use it too.
if [[ -n "$ROUTE_TABLE" && "$ROUTE_TABLE" != "main" ]]; then
    ip -4 rule show | grep -Eq "to $BACKEND_CIDR (lookup|table) $ROUTE_TABLE" ||
        ip -4 rule add priority 10000 to "$BACKEND_CIDR" lookup "$ROUTE_TABLE"
fi

sysctl -w net.ipv4.ip_forward=1 >/dev/null
sysctl -w net.ipv4.conf.all.rp_filter=0 >/dev/null || true
sysctl -w net.ipv4.conf.default.rp_filter=0 >/dev/null || true
sysctl -w net.ipv4.conf."$OUT_IF".rp_filter=0 >/dev/null || true

iptables -w -C FORWARD -s "$WORKER_CIDR" -d "$BACKEND_CIDR" -j ACCEPT 2>/dev/null ||
    iptables -w -I FORWARD 1 -s "$WORKER_CIDR" -d "$BACKEND_CIDR" -j ACCEPT

iptables -w -C FORWARD -s "$BACKEND_CIDR" -d "$WORKER_CIDR" \
    -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT 2>/dev/null ||
    iptables -w -I FORWARD 1 -s "$BACKEND_CIDR" -d "$WORKER_CIDR" \
        -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT

iptables -w -t nat -C POSTROUTING -s "$WORKER_CIDR" -d "$BACKEND_CIDR" \
    -o "$OUT_IF" -j MASQUERADE 2>/dev/null ||
    iptables -w -t nat -A POSTROUTING -s "$WORKER_CIDR" -d "$BACKEND_CIDR" \
        -o "$OUT_IF" -j MASQUERADE

iptables -w -t mangle -C FORWARD -s "$WORKER_CIDR" -d "$BACKEND_CIDR" \
    -o "$OUT_IF" -p tcp --tcp-flags SYN,RST SYN \
    -j TCPMSS --clamp-mss-to-pmtu 2>/dev/null ||
    iptables -w -t mangle -A FORWARD -s "$WORKER_CIDR" -d "$BACKEND_CIDR" \
        -o "$OUT_IF" -p tcp --tcp-flags SYN,RST SYN \
        -j TCPMSS --clamp-mss-to-pmtu 2>/dev/null || true

printf 'ukama-p0 gateway ready: %s -> %s via %s\n' \
    "$WORKER_CIDR" "$BACKEND_CIDR" "$OUT_IF"
GATEWAY_EOF
chmod 755 /usr/local/sbin/ukama-p0-gateway-apply

cat >/etc/sysctl.d/99-ukama-p0-gateway.conf <<SYSCTL_EOF
net.ipv4.ip_forward=1
net.ipv4.conf.all.rp_filter=0
net.ipv4.conf.default.rp_filter=0
SYSCTL_EOF
sysctl -w net.ipv4.ip_forward=1 >/dev/null
sysctl -w net.ipv4.conf.all.rp_filter=0 >/dev/null || true
sysctl -w net.ipv4.conf.default.rp_filter=0 >/dev/null || true

# Apply immediately and show the real failing command if anything is wrong.
# The previous implementation started through systemd first, which hid the
# useful error behind a generic "control process exited" message.
/usr/local/sbin/ukama-p0-gateway-apply

# Persistence is best effort. The gateway is already active for this batch.
if command -v systemctl >/dev/null 2>&1; then
    cat >/etc/systemd/system/ukama-p0-gateway.service <<'UNIT_EOF'
[Unit]
Description=Ukama P0 backend forwarding gateway
Wants=network-online.target
After=network-online.target

[Service]
Type=oneshot
ExecStart=/usr/local/sbin/ukama-p0-gateway-apply
RemainAfterExit=yes

[Install]
WantedBy=multi-user.target
UNIT_EOF
    systemctl daemon-reload
    systemctl reset-failed ukama-p0-gateway.service 2>/dev/null || true
    systemctl enable ukama-p0-gateway.service >/dev/null 2>&1 || true
fi
REMOTE_EOF

ssh "${ssh_args[@]}" "$ssh_user@$ssh_host" \
    "sudo bash -s -- $(printf '%q' "$worker_cidr") $(printf '%q' "$BACKEND_ROUTE_CIDR") $(printf '%q' "$BACKEND_TEST_IP")" \
    <"$remote_script"

# Cache the resolved next hop for worker.env. Preserve all existing state.
state_file="$P0_AWS_STATE"
state_tmp="$(mktemp)"
if [[ -r "$state_file" ]]; then
    grep -vE '^(BACKEND_GATEWAY_PRIVATE_IP|BACKEND_GATEWAY_SECURITY_GROUP_ID)=' \
        "$state_file" >"$state_tmp" || true
fi
printf 'BACKEND_GATEWAY_PRIVATE_IP=%s\n' "$gateway_private_ip" >>"$state_tmp"
printf 'BACKEND_GATEWAY_SECURITY_GROUP_ID=%s\n' "$gateway_sg" >>"$state_tmp"
mv "$state_tmp" "$state_file"
chmod 600 "$state_file"

((QUIET)) || printf 'backend gateway is ready: %s via %s\n' \
    "$BACKEND_ROUTE_CIDR" "$gateway_private_ip"
