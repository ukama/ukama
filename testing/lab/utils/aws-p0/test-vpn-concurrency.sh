#!/usr/bin/env bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
# Copyright (c) 2026-present, Ukama Inc.
#
# Standalone diagnostic; does not run scenarios or modify the existing runner.
# Place beside utils/aws-p0/config.env and lib.sh, then run from your laptop:
#   bash ./test-vpn-concurrency.sh --vpn ~/Downloads/downloaded-dev-cluster-client-config.ovpn
# Keep your existing laptop VPN connection running throughout the test.
# Two EC2 instances are billed until terminated. No gateway/route-table/IAM edits.
# The profile is temporarily stored in the existing private S3 bucket, encrypted
# at rest; its exact object version is deleted on exit. Never run with bash -x.
# Revision 3: use jq-compatible variable names for schedule and worker reports.

set +x
set -Eeuo pipefail
umask 077

aws_call() {
    aws --no-cli-pager --cli-connect-timeout 10 --cli-read-timeout 20 "$@"
}

die() {
    printf 'error: %s\n' "$*" >&2
    exit 2
}

write_schedule() {
    jq -n --argjson connect "$1" --argjson start "$2" --argjson window_end "$3" \
        '{connect:$connect,start:$start,end:$window_end}'
}

write_worker_result() {
    jq -n --arg worker "$WORKER_ID" --arg reason "$reason" \
        --arg s3_over_vpn "$s3_over_vpn" --arg s3_error "$s3_error" \
        --argjson rc "$rc" --argjson start "$start" --argjson window_end "$end" \
        --argjson checks "$checks" --argjson failures "$failures" \
        '{worker:$worker,passed:($rc==0),reason:$reason,start:$start,end:$window_end,
          checks:$checks,failures:$failures,s3_over_vpn:$s3_over_vpn,s3_error:$s3_error}'
}

json_preflight() (
    # Execute the real JSON writers with the installed jq BEFORE paying for EC2.
    # This also catches older jq parsers that reject keyword variable names.
    WORKER_ID=preflight reason=preflight s3_over_vpn=NOT_TESTED s3_error=''
    rc=1 start=0 end=0 checks=0 failures=0
    write_worker_result | jq -e '.passed == false and .end == 0' >/dev/null
    write_schedule 1 2 3 | jq -e '.connect == 1 and .start == 2 and .end == 3' >/dev/null
)

# Resolve on every sample, pin curl to that address, and check its actual route.
# HTTP 4xx still demonstrates transport access; this is not an API health test.
probe() {
    local address route code
    address="$(timeout 8 getent ahostsv4 "$CHECK_HOST" | awk 'NR==1 {print $1}')" || return 1
    [[ -n "$address" ]] || return 1
    route="$(ip -4 route get "$address")" || return 1
    [[ " $route " == *" dev $CHECK_DEVICE "* ]] || return 1
    code="$(curl --noproxy '*' -sS -o /dev/null --connect-timeout 5 --max-time 8 \
        --resolve "$CHECK_HOST:$CHECK_PORT:$address" -w '%{http_code}' \
        "$CHECK_URL" 2>/dev/null)" || return 1
    [[ "$code" =~ ^[1-5][0-9][0-9]$ ]]
}

network_snapshot() {
    # Routes/DNS only; never request or print instance-role credentials.
    ip -4 route show
    ip -4 rule show
    ip -4 route get 169.254.169.254
    if command -v resolvectl >/dev/null; then timeout 5 resolvectl status; fi
    if [[ -n "${AWS_REGION:-}" ]]; then
        timeout 8 getent ahostsv4 "s3.$AWS_REGION.amazonaws.com" || true
    fi
}

worker_main() {
    : "${TEST_URI:?}" "${WORKER_ID:?}" "${CHECK_URL:?}" "${CHECK_HOST:?}" "${CHECK_PORT:?}"
    # EXIT also runs after errexit unwinds this function: cleanup state must
    # remain available outside the function's local scope.
    work=/run/ukama-vpn-check logs=/var/log/ukama-vpn-check
    vpn_pid='' reason=bootstrap_failed rc=1 checks=0 failures=0
    connect=0 start=0 end=0 first=0 last=0
    s3_over_vpn=NOT_TESTED s3_error=''
    mkdir -p "$work" "$logs"
    exec >"$logs/bootstrap.log" 2>&1
    export CHECK_DEVICE=uvpn0 AWS_RETRY_MODE=standard AWS_MAX_ATTEMPTS=2

    worker_finish() {
        local shell_rc=$? archive=/run/ukama-vpn-check-results.tar.gz
        trap - EXIT INT TERM
        set +e
        if [[ -n "$vpn_pid" ]]; then
            kill -TERM "$vpn_pid" 2>/dev/null
            for ((n=0; n<15; n++)); do
                kill -0 "$vpn_pid" 2>/dev/null || break
                sleep 1
            done
            kill -KILL "$vpn_pid" 2>/dev/null
            wait "$vpn_pid" 2>/dev/null
        fi
        rm -f "$work/client.ovpn"
        if ((shell_rc != 0)); then rc=1; fi
        write_worker_result >"$logs/result.json"
        tar -C "$logs" -czf "$archive" .
        # Stop VPN before final uploads, so broken tunnel routing cannot hide
        # the diagnostic. Publish result last, after the archive is durable.
        if aws_call s3 cp "$archive" "$TEST_URI/output/$WORKER_ID.tar.gz" --only-show-errors; then
            aws_call s3 cp "$logs/result.json" "$TEST_URI/output/$WORKER_ID.json" --only-show-errors
        fi
        shutdown -h now
    }
    trap worker_finish EXIT
    trap 'reason=interrupted; exit 1' INT TERM

    # Bootstrap watchdog is also installed by user-data before this script.
    command -v apt-get >/dev/null || { reason=ubuntu_required; exit 1; }
    if ! command -v openvpn >/dev/null; then
        reason=openvpn_install_failed
        timeout 180 apt-get update -qq
        timeout 180 env DEBIAN_FRONTEND=noninteractive apt-get install -y -qq openvpn
    fi
    for name in jq curl ip getent timeout tar; do
        command -v "$name" >/dev/null || { reason="missing_$name"; exit 1; }
    done
    reason=json_preflight_failed
    json_preflight
    reason=profile_download_failed
    aws_call s3 cp "$TEST_URI/input/client.ovpn" "$work/client.ovpn" --only-show-errors
    network_snapshot >"$logs/network-before-vpn.txt" 2>&1 || true

    # Honour server-pushed DNS on Ubuntu without changing global DNS settings.
    # iproute/tunnel setup remains OpenVPN's responsibility.
    cat >"$work/dns-up.sh" <<'DNS'
#!/usr/bin/env bash
set -eu
command -v resolvectl >/dev/null || exit 0
dns=()
while IFS= read -r name; do
    value="${!name}"
    if [[ "$value" == 'dhcp-option DNS '* ]]; then dns+=("${value#dhcp-option DNS }"); fi
done < <(compgen -A variable foreign_option_ || true)
if ((${#dns[@]})); then
    resolvectl dns "$dev" "${dns[@]}"
    resolvectl domain "$dev" '~udev.ukama.com'
fi
DNS
    chmod 700 "$work/dns-up.sh"

    # Coordinate while both EC2 instances still use their ordinary AWS network.
    # A VPN that prevents S3/IMDS access must not prevent measuring VPN sessions.
    reason=mailbox_unreachable_before_vpn
    printf '{"prepared":true}\n' >"$work/ready.json"
    aws_call s3 cp "$work/ready.json" "$TEST_URI/output/$WORKER_ID.ready" --only-show-errors
    reason=waiting_for_second_worker
    ready_deadline=$((SECONDS + 600))
    while ! aws_call s3 cp "$TEST_URI/input/start.json" "$work/start.json" --only-show-errors \
        >"$logs/schedule-download.log" 2>&1; do
        ((SECONDS < ready_deadline)) || { reason=start_timeout; exit 1; }
        sleep 3
    done
    connect="$(jq -er '.connect' "$work/start.json")"
    start="$(jq -er '.start' "$work/start.json")"
    end="$(jq -er '.end' "$work/start.json")"
    (($(date +%s) <= connect)) || { reason=connect_window_missed; exit 1; }
    while (($(date +%s) < connect)); do sleep 1; done

    reason=vpn_start_failed
    openvpn --config "$work/client.ovpn" --dev "$CHECK_DEVICE" --dev-type tun \
        --script-security 2 --up "$work/dns-up.sh" --up-restart \
        --log "$logs/openvpn.log" --verb 3 </dev/null &
    vpn_pid=$!
    ready_deadline=$((SECONDS + 120))
    while ! grep -q 'Initialization Sequence Completed' "$logs/openvpn.log" 2>/dev/null; do
        kill -0 "$vpn_pid" 2>/dev/null || exit 1
        ((SECONDS < ready_deadline)) || { reason=vpn_connect_timeout; exit 1; }
        sleep 2
    done
    reason=backend_probe_failed
    probe
    ip -4 addr show dev "$CHECK_DEVICE" >"$logs/tunnel.txt"
    (($(date +%s) <= start)) || { reason=start_window_missed; exit 1; }
    while (($(date +%s) < start)); do
        kill -0 "$vpn_pid" 2>/dev/null && probe || { reason=vpn_lost_before_window; exit 1; }
        [[ "$(grep -c 'Initialization Sequence Completed' "$logs/openvpn.log")" == 1 ]] || {
            reason=vpn_reconnected_before_window; exit 1;
        }
        sleep 3
    done
    reason=monitoring
    printf 'epoch\tpassed\n' >"$logs/samples.tsv"
    while :; do
        stamp="$(date +%s)"
        ((stamp < end)) || break
        count="$(grep -c 'Initialization Sequence Completed' "$logs/openvpn.log" || true)"
        ((first != 0)) || first=$stamp
        if ((last != 0 && stamp-last > 25)); then failures=$((failures+1)); fi
        last=$stamp
        checks=$((checks+1))
        if kill -0 "$vpn_pid" 2>/dev/null && [[ "$count" == 1 ]] && probe; then
            printf '%s\t1\n' "$stamp" >>"$logs/samples.tsv"
        else
            failures=$((failures+1))
            printf '%s\t0\n' "$stamp" >>"$logs/samples.tsv"
        fi
        sleep 5
    done
    count="$(grep -c 'Initialization Sequence Completed' "$logs/openvpn.log" || true)"
    if ((checks >= (end-start)/25 && failures == 0 && first <= start+25 && last >= end-25)) &&
        [[ "$count" == 1 ]] && kill -0 "$vpn_pid" 2>/dev/null &&
        ! grep -Eq 'SIGUSR1|AUTH_FAILED|Restart pause' "$logs/openvpn.log"; then
        reason=stable_during_window; rc=0
    else
        reason=connectivity_lost_or_vpn_reconnected
    fi

    # Measure the runner's S3 path separately, after the VPN observation window.
    # Preserve the actual AWS error: the earlier label alone could not tell a
    # routing/DNS failure from an instance-role credential retrieval failure.
    if kill -0 "$vpn_pid" 2>/dev/null; then
        network_snapshot >"$logs/network-during-vpn.txt" 2>&1 || true
        printf '{"probe":true}\n' >"$work/mailbox.json"
        if timeout 45 aws --no-cli-pager --cli-connect-timeout 10 --cli-read-timeout 20 \
            s3 cp "$work/mailbox.json" "$TEST_URI/output/$WORKER_ID.mailbox-probe" \
            --only-show-errors >"$logs/s3-over-vpn.log" 2>&1; then
            s3_over_vpn=PASS
        else
            s3_rc=$?
            s3_over_vpn=FAIL
            s3_error="exit=$s3_rc $(head -c 1500 "$logs/s3-over-vpn.log")"
        fi
    fi
    exit 0
}

if [[ "${1:-}" == --worker ]]; then
    worker_main
    exit 0
fi

usage() {
    cat <<'USAGE'
usage: bash test-vpn-concurrency.sh --vpn FILE [options]

  --vpn FILE             Self-contained OpenVPN profile (inline ca/cert/key).
  --seconds N            Common observation window, 60..600 seconds (default 300).
  --url URL              Private backend HTTP(S) URL (default BFF GraphQL).
  --instance-type TYPE   EC2 type (default from config.env).
  --dry-run             Validate local inputs and print plan; no AWS calls.

Run on your Linux laptop with its existing VPN connected. Uses the same
config.env/.state.env, AWS profile, AMI, subnet, security group, role and S3
bucket as aws-p0. Does not call setup.sh or configure the old gateway.
Credentials are read from FILE at run time, never embedded in this script.
Revision 3 fixes jq compatibility; S3-over-VPN remains a separate check.
USAGE
}

vpn='' duration=300 url=https://bff.udev.ukama.com/gateway/graphql type='' dry=0
while (($#)); do
    case "$1" in
        --vpn|--seconds|--url|--instance-type)
            (($# >= 2)) || die "$1 requires a value"
            case "$1" in
                --vpn) vpn=$2;; --seconds) duration=$2;; --url) url=$2;; --instance-type) type=$2;;
            esac
            shift 2;;
        --dry-run) dry=1; shift;;
        -h|--help) usage; exit 0;;
        *) die "unknown argument: $1";;
    esac
done
[[ -r "$vpn" ]] || die 'provide --vpn with your readable .ovpn file'
[[ "$duration" =~ ^[0-9]+$ ]] && ((duration >= 60 && duration <= 600)) || die '--seconds must be 60..600'
script_dir="$(CDPATH= cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
[[ -r "$script_dir/lib.sh" ]] || die 'place this script in utils/aws-p0 beside lib.sh'
# shellcheck source=/dev/null
. "$script_dir/lib.sh"
p0_load_config "$script_dir"
p0_require_config AWS_REGION S3_BUCKET S3_PREFIX AMI_ID SUBNET_ID SECURITY_GROUP_ID INSTANCE_PROFILE_NAME
type="${type:-${INSTANCE_TYPE:-t3.small}}"
for cmd in python3 jq curl ip getent timeout; do command -v "$cmd" >/dev/null || die "missing command: $cmd"; done
json_preflight
printf 'Local JSON checks passed (%s).\n' "$(jq --version)"

# Validate only configuration structure; never print the certificate/key.
metadata="$(python3 - "$vpn" "$url" <<'PY'
import json,pathlib,re,shlex,sys,urllib.parse
s=pathlib.Path(sys.argv[1]).read_text()
for tag in ('ca','cert','key'):
    if not re.search(r'(?ms)^<'+tag+r'>\s*.*?^</'+tag+r'>\s*$',s):
        sys.exit('profile must contain inline ca, cert and key blocks')
plain=re.sub(r'(?ms)^<(\w+)>\s*.*?^</\1>\s*$', '', s)
allowed={'client','dev','dev-type','proto','remote','remote-random-hostname',
 'resolv-retry','nobind','remote-cert-tls','cipher','data-ciphers','data-ciphers-fallback',
 'verb','reneg-sec','persist-key','persist-tun','auth','auth-nocache','connect-retry',
 'connect-timeout','connect-retry-max','explicit-exit-notify','setenv','route','route-ipv6',
 'route-delay','route-metric','route-nopull','pull-filter','redirect-gateway','dhcp-option',
 'ping','ping-restart','keepalive','tun-mtu','mssfix','tls-version-min','tls-cipher',
 'tls-ciphersuites','verify-x509-name','remote-random','nobind','float'}
for line in plain.splitlines():
    words=shlex.split(line,comments=True)
    if words and words[0] not in allowed:
        sys.exit('profile directive needs review before unattended use: '+words[0])
u=urllib.parse.urlsplit(sys.argv[2])
if u.scheme not in ('http','https') or not u.hostname or u.username or u.password or u.query or u.fragment:
    sys.exit('--url must be an HTTP(S) endpoint without credentials, query or fragment')
if ':' in u.hostname: sys.exit('this diagnostic requires an IPv4 endpoint')
print(json.dumps({'host':u.hostname,'port':u.port or (443 if u.scheme=='https' else 80)}))
PY
)"
export CHECK_HOST CHECK_PORT CHECK_URL="$url"
CHECK_HOST="$(jq -r .host <<<"$metadata")"
CHECK_PORT="$(jq -r .port <<<"$metadata")"
printf 'VPN diagnostic revision 3\n'
printf 'Plan: two %s instances in %s, %ss shared VPN observation window.\n' "$type" "$AWS_REGION" "$duration"
printf 'Backend probe: %s\n' "$CHECK_URL"
if ((dry)); then printf 'Dry run: no AWS calls or VPN connections made.\n'; exit 0; fi
command -v aws >/dev/null || die 'missing command: aws'
aws_call sts get-caller-identity >/dev/null
export AWS_RETRY_MODE=standard AWS_MAX_ATTEMPTS=2

# Check the laptop is actually using its existing tunnel, not a public route.
laptop_ip="$(timeout 8 getent ahostsv4 "$CHECK_HOST" | awk 'NR==1 {print $1}')"
[[ -n "$laptop_ip" ]] || die 'backend DNS failed on your laptop'
CHECK_DEVICE="$(ip -4 route get "$laptop_ip" | awk '{for(i=1;i<NF;i++) if($i=="dev") {print $(i+1);exit}}')"
export CHECK_DEVICE
[[ -r "/sys/class/net/$CHECK_DEVICE/tun_flags" ]] || die 'backend route on laptop is not a TUN device; connect your OpenVPN first'
probe || die 'backend probe failed on your laptop before launch'

# This is the bucket already configured by setup.sh. Check its privacy before
# temporarily placing the credential there; do not change bucket settings.
aws_call s3api get-public-access-block --bucket "$S3_BUCKET" |
    jq -e '.PublicAccessBlockConfiguration | .BlockPublicAcls and .IgnorePublicAcls and .BlockPublicPolicy and .RestrictPublicBuckets' >/dev/null ||
    die 'existing S3 bucket must have all four Block Public Access settings enabled'

test_id="vpn-check-$(date -u +%Y%m%dT%H%M%SZ)-$(python3 -c 'import secrets; print(secrets.token_hex(4))')"
key="${S3_PREFIX%/}/$test_id"
uri="s3://$S3_BUCKET/$key"
lab_root="$(CDPATH= cd -- "$script_dir/../.." && pwd)"
out="$lab_root/runs/p0-aws/$test_id"
mkdir -p "$out/output"
ids=() monitor_pid='' uploaded=0 launch_attempted=0 version='' cleanup_failed=0

cleanup_test() {
    local saved_rc=$? found id head_rc
    trap - EXIT INT TERM
    set +e
    if [[ -n "$monitor_pid" ]]; then kill "$monitor_pid" 2>/dev/null; wait "$monitor_pid" 2>/dev/null; fi
    if ((launch_attempted)); then
        # Tags recover instances even if run-instances succeeded but its reply
        # was lost. A query error is not mistaken for an empty instance list.
        found="$(aws_call ec2 describe-instances --filters "Name=tag:UkamaP0VpnTest,Values=$test_id" \
            'Name=instance-state-name,Values=pending,running,stopping,stopped' \
            --query 'Reservations[].Instances[].InstanceId' --output text)"
        if (($? != 0)); then cleanup_failed=1; fi
        for id in $found; do [[ "$id" == i-* ]] && ids+=("$id"); done
        if ((${#ids[@]})); then
            mapfile -t ids < <(printf '%s\n' "${ids[@]}" | sort -u)
            aws_call ec2 terminate-instances --instance-ids "${ids[@]}" >/dev/null || cleanup_failed=1
            aws_call ec2 wait instance-terminated --instance-ids "${ids[@]}" || cleanup_failed=1
        fi
    fi
    if ((uploaded)); then
        if [[ -z "$version" ]]; then
            aws_call s3api head-object --bucket "$S3_BUCKET" --key "$key/input/client.ovpn" \
                >"$out/profile-head.json" 2>"$out/profile-head-error.txt"
            head_rc=$?
            if ((head_rc == 0)); then
                version="$(jq -r '.VersionId // empty' "$out/profile-head.json")"
            elif ! grep -Eq '404|Not Found|NoSuchKey' "$out/profile-head-error.txt"; then
                cleanup_failed=1
            fi
        fi
        # Exact version deletion avoids leaving private keys in version history.
        if [[ -n "$version" ]]; then
            aws_call s3api delete-object --bucket "$S3_BUCKET" --key "$key/input/client.ovpn" \
                --version-id "$version" >/dev/null || cleanup_failed=1
        else
            aws_call s3api delete-object --bucket "$S3_BUCKET" --key "$key/input/client.ovpn" >/dev/null || cleanup_failed=1
        fi
    fi
    if ((cleanup_failed)); then
        printf 'CLEANUP INCOMPLETE: inspect EC2 tag UkamaP0VpnTest=%s and S3 %s/input/client.ovpn\n' "$test_id" "$uri" >&2
        saved_rc=1
    else
        printf 'Cleanup complete: test instances terminated; temporary VPN profile removed.\n'
    fi
    printf 'Diagnostics: %s\n' "$out"
    exit "$saved_rc"
}
trap cleanup_test EXIT
trap 'exit 130' INT
trap 'exit 143' TERM

printf 'Test ID: %s\nKeep your laptop VPN connected.\n' "$test_id"
aws_call s3 cp "${BASH_SOURCE[0]}" "$uri/input/test.sh" --only-show-errors
# Mark attempted upload before the API call so an interrupted response still
# triggers cleanup. Resolve its version via head-object if needed below.
uploaded=1
aws_call s3api put-object --bucket "$S3_BUCKET" --key "$key/input/client.ovpn" \
    --body "$vpn" --server-side-encryption AES256 >"$out/profile-upload.json"
version="$(jq -r '.VersionId // empty' "$out/profile-upload.json")"

root_device="$(aws_call ec2 describe-images --image-ids "$AMI_ID" --query 'Images[0].RootDeviceName' --output text)"
[[ "$root_device" == /dev/* ]] || die 'cannot determine the AMI root device'
jq -n --arg device "$root_device" \
    '[{DeviceName:$device,Ebs:{DeleteOnTermination:true,Encrypted:true}}]' >"$out/block-device.json"

(
    printf 'epoch\tpassed\n' >"$out/laptop.tsv"
    while :; do
        stamp="$(date +%s)"
        if probe; then printf '%s\t1\n' "$stamp"; else printf '%s\t0\n' "$stamp"; fi
        sleep 5
    done >>"$out/laptop.tsv"
) &
monitor_pid=$!

jq -n --arg subnet "$SUBNET_ID" --arg sg "$SECURITY_GROUP_ID" \
    --argjson public "${ASSOCIATE_PUBLIC_IP:-true}" \
    '[{DeviceIndex:0,SubnetId:$subnet,Groups:[$sg],AssociatePublicIpAddress:$public,DeleteOnTermination:true}]' >"$out/network.json"
for worker in worker-01 worker-02; do
    {
        printf '#!/usr/bin/env bash\nset -Eeuo pipefail\numask 077\n'
        printf 'shutdown -h +35\n'
        printf "trap 'shutdown -h now' EXIT\n"
        printf 'export AWS_DEFAULT_REGION=%q AWS_REGION=%q AWS_RETRY_MODE=standard AWS_MAX_ATTEMPTS=2\n' "$AWS_REGION" "$AWS_REGION"
        printf 'export TEST_URI=%q WORKER_ID=%q CHECK_URL=%q CHECK_HOST=%q CHECK_PORT=%q\n' \
            "$uri" "$worker" "$CHECK_URL" "$CHECK_HOST" "$CHECK_PORT"
        printf 'aws --cli-connect-timeout 10 --cli-read-timeout 20 s3 cp %q /run/vpn-check.sh --only-show-errors\n' "$uri/input/test.sh"
        printf 'bash /run/vpn-check.sh --worker\n'
    } >"$out/$worker-user-data.sh"
    launch_attempted=1
    id="$(aws_call ec2 run-instances --image-id "$AMI_ID" --instance-type "$type" \
        --count 1 --client-token "$test_id-$worker" \
        --iam-instance-profile "Name=$INSTANCE_PROFILE_NAME" \
        --network-interfaces "file://$out/network.json" \
        --block-device-mappings "file://$out/block-device.json" \
        --metadata-options 'HttpEndpoint=enabled,HttpTokens=required' \
        --instance-initiated-shutdown-behavior terminate \
        --user-data "file://$out/$worker-user-data.sh" \
        --tag-specifications "ResourceType=instance,Tags=[{Key=Name,Value=$test_id-$worker},{Key=UkamaP0VpnTest,Value=$test_id}]" \
        --query 'Instances[0].InstanceId' --output text)"
    [[ "$id" == i-* ]] || die "unexpected EC2 launch result for $worker"
    ids+=("$id")
    printf '%s\t%s\n' "$worker" "$id" | tee -a "$out/instances.tsv"
done

deadline=$((SECONDS + 720))
window_started=0 last_progress=0
while :; do
    aws_call s3 sync "$uri/output/" "$out/output/" --only-show-errors || die 'cannot read worker status from S3'
    ready=0 finished=0 early_failure=0
    for worker in worker-01 worker-02; do
        [[ ! -f "$out/output/$worker.ready" ]] || ready=$((ready+1))
        if [[ -f "$out/output/$worker.json" ]]; then
            finished=$((finished+1))
            jq -e '.passed == true' "$out/output/$worker.json" >/dev/null || early_failure=1
        fi
    done
    if ((finished == 2 || early_failure)); then break; fi
    if ((ready == 2 && !window_started)); then
        connect=$(($(date +%s)+30)); start=$((connect+150)); end=$((start+duration))
        write_schedule "$connect" "$start" "$end" >"$out/window.json"
        aws_call s3 cp "$out/window.json" "$uri/input/start.json" --only-show-errors
        window_started=1
        deadline=$((SECONDS + duration + 330))
        printf 'Both workers prepared. VPN connections start in 30s; common %ss window starts in 180s.\n' "$duration"
    fi
    if ((SECONDS-last_progress >= 30)); then
        printf 'Workers prepared=%s/2 finished=%s/2; laptop monitoring active.\n' "$ready" "$finished"
        last_progress=$SECONDS
    fi
    ((SECONDS < deadline)) || { printf 'Test timed out waiting for worker reports.\n' >&2; break; }
    sleep 5
done
kill "$monitor_pid" 2>/dev/null || true
wait "$monitor_pid" 2>/dev/null || true
monitor_pid=''
python3 - "$out" <<'PY' | tee "$out/summary.txt"
import csv,json,pathlib,sys
p=pathlib.Path(sys.argv[1]); ok=True; s3_failed=[]
window=json.loads((p/'window.json').read_text()) if (p/'window.json').exists() else None
print('VPN concurrency diagnostic')
for name in ('worker-01','worker-02'):
    f=p/'output'/(name+'.json')
    if not f.exists():
        print(name+': INCOMPLETE (no final report)');ok=False;continue
    d=json.loads(f.read_text())
    passed=bool(window and d.get('passed') and d.get('start')==window['start'] and d.get('end')==window['end'] and d.get('checks',0)>0 and d.get('failures')==0)
    print(f"{name}: {'PASS' if passed else 'FAIL'} — {d.get('reason','unknown')}; samples={d.get('checks',0)}, failures={d.get('failures',0)}")
    if d.get('s3_over_vpn')=='FAIL':
        s3_failed.append(name)
        print(f"  S3 during VPN: FAIL — {' '.join(d.get('s3_error','').split())}")
    elif d.get('s3_over_vpn')=='PASS':
        print('  S3 during VPN: PASS')
    ok &= passed
rows=list(csv.DictReader((p/'laptop.tsv').open(),delimiter='\t'))
samples=[r for r in rows if window and window['start']<=int(r['epoch'])<window['end']]
times=[int(r['epoch']) for r in samples]
# Endpoint probes have an 8s DNS limit, 8s HTTP limit and 5s sampling interval.
# Reject sparse/missing observations and any failed laptop sample, including
# disconnects caused by workers before the common observation window began.
coverage=bool(window and times and times[0]<=window['start']+25 and times[-1]>=window['end']-25 and all(b-a<=25 for a,b in zip(times,times[1:])))
failed=sum(r['passed']!='1' for r in rows)
laptop_ok=coverage and failed==0
laptop_state='PASS' if laptop_ok else ('FAIL' if failed else 'NOT TESTED / INCOMPLETE')
print(f"laptop: {laptop_state} — window samples={len(samples)}, total failed samples={failed}")
ok &= laptop_ok
print('VPN CONCURRENCY: '+('PASS — both workers and laptop retained sampled backend access in the same window.' if ok else 'FAIL / INCOMPLETE — inspect the reports before changing the runner.'))
if s3_failed: print('RUNNER FOLLOW-UP: S3 access during VPN failed on '+', '.join(s3_failed)+'. See the AWS errors above.')
print('Scope: VPN session stability and host HTTP reachability; not 29-worker capacity or a full lab scenario.')
sys.exit(0 if ok else 1)
PY
