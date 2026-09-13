#!/usr/bin/env bash
# Direct worker VPN helpers. Source from worker.sh; no gateway or NAT setup.

p0_imds_role_check() {
    local token role
    token="$(curl --noproxy '*' -fsS --connect-timeout 3 --max-time 5 -X PUT \
        -H 'X-aws-ec2-metadata-token-ttl-seconds: 60' \
        http://169.254.169.254/latest/api/token)" || return 1
    role="$(curl --noproxy '*' -fsS --connect-timeout 3 --max-time 5 \
        -H "X-aws-ec2-metadata-token: $token" \
        http://169.254.169.254/latest/meta-data/iam/security-credentials/)" || return 1
    [[ "$role" =~ ^[A-Za-z0-9+=,.@_-]+$ ]] || return 1
    # Validate fresh role credentials without printing or caching them.
    curl --noproxy '*' -fsS --connect-timeout 3 --max-time 5 \
        -H "X-aws-ec2-metadata-token: $token" \
        "http://169.254.169.254/latest/meta-data/iam/security-credentials/$role" |
        jq -e '.Code == "Success" and (.AccessKeyId | length > 0) and
            (.SecretAccessKey | length > 0) and (.Token | length > 0)' >/dev/null
}

p0_preserve_imds_route() {
    local route dev gateway
    local args=(169.254.169.254/32)
    route="$(ip -4 route get 169.254.169.254)" || return 1
    dev="$(awk '{for(i=1;i<NF;i++) if($i=="dev") {print $(i+1);exit}}' <<<"$route")"
    gateway="$(awk '{for(i=1;i<NF;i++) if($i=="via") {print $(i+1);exit}}' <<<"$route")"
    [[ -n "$dev" && "$dev" != uvpn0 ]] || return 1
    [[ -z "$gateway" ]] || args+=(via "$gateway")
    args+=(dev "$dev")
    # Keep only EC2's link-local metadata endpoint on its original interface.
    # The AWS CLI remains free to refresh short-lived instance-role credentials.
    ip -4 route replace "${args[@]}" || return 1
    export NO_PROXY="${NO_PROXY:+$NO_PROXY,}169.254.169.254"
    export no_proxy="${no_proxy:+$no_proxy,}169.254.169.254"
    p0_imds_role_check
}

p0_vpn_start() {
    local profile="$1" deadline
    if ! command -v openvpn >/dev/null; then
        command -v apt-get >/dev/null || { echo 'OpenVPN missing; use the Ubuntu worker AMI' >&2; return 1; }
        timeout 180 apt-get update -qq || return 1
        timeout 180 env DEBIAN_FRONTEND=noninteractive apt-get install -y -qq openvpn || return 1
    fi
    p0_preserve_imds_route || { echo 'cannot preserve EC2 instance-role access' >&2; return 1; }
    cat >"$WORK_ROOT/vpn-dns.sh" <<'DNS'
#!/usr/bin/env bash
set -eu
command -v resolvectl >/dev/null || exit 0
dns=()
while IFS= read -r name; do
    value="${!name}"
    [[ "$value" != 'dhcp-option DNS '* ]] || dns+=("${value#dhcp-option DNS }")
done < <(compgen -A variable foreign_option_ || true)
if ((${#dns[@]})); then
    resolvectl dns "$dev" "${dns[@]}"
    resolvectl domain "$dev" '~udev.ukama.com'
fi
DNS
    chmod 700 "$WORK_ROOT/vpn-dns.sh"
    openvpn --config "$profile" --dev uvpn0 --dev-type tun --script-security 2 \
        --up "$WORK_ROOT/vpn-dns.sh" --up-restart --verb 3 --log "$WORK_ROOT/vpn.log" </dev/null &
    VPN_PID=$!
    deadline=$((SECONDS + ${VPN_CONNECT_TIMEOUT_SECONDS:-120}))
    while ! grep -q 'Initialization Sequence Completed' "$WORK_ROOT/vpn.log" 2>/dev/null; do
        kill -0 "$VPN_PID" 2>/dev/null || { echo 'OpenVPN exited during startup' >&2; return 1; }
        ((SECONDS < deadline)) || { echo 'OpenVPN connection timed out' >&2; return 1; }
        sleep 2
    done
    p0_imds_role_check || { echo 'EC2 instance-role credentials unavailable after VPN startup' >&2; return 1; }
    echo 'VPN connected; EC2 instance-role credential retrieval verified.'
}

p0_vpn_stop() {
    local n
    if [[ -n "${VPN_PID:-}" ]]; then
        kill -TERM "$VPN_PID" 2>/dev/null || true
        for ((n=0;n<15;n++)); do
            kill -0 "$VPN_PID" 2>/dev/null || break
            sleep 1
        done
        kill -KILL "$VPN_PID" 2>/dev/null || true
        wait "$VPN_PID" 2>/dev/null || true
        VPN_PID=''
    fi
    rm -f "$WORK_ROOT/client.ovpn"
}
