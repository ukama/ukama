#!/usr/bin/env bash
# Store the supplied VPN profile in the existing worker secret; preserve its
# other fields. The profile never enters a source archive or EC2 user-data.
set +x
set -Eeuo pipefail
umask 077
SCRIPT_DIR="$(CDPATH= cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
. "$SCRIPT_DIR/lib.sh"
p0_load_config "$SCRIPT_DIR"
check=0
if [[ "${1:-}" == --check ]]; then check=1; shift; fi
[[ $# == 1 && -r "$1" ]] || p0_die 'usage: configure-vpn.sh [--check] CLIENT.ovpn'
python3 - "$1" <<'PY'
import pathlib,re,shlex,sys
s=pathlib.Path(sys.argv[1]).read_text()
for tag in ('ca','cert','key'):
    if not re.search(r'(?ms)^<'+tag+r'>\s*.*?^</'+tag+r'>\s*$',s):
        sys.exit('VPN profile must contain inline ca, cert and key blocks')
plain=re.sub(r'(?ms)^<(\w+)>\s*.*?^</\1>\s*$', '', s)
for line in plain.splitlines():
    words=shlex.split(line,comments=True)
    if words and words[0] in {'auth-user-pass','askpass','auth-federate','config','up','down','route-up','plugin','daemon','log','log-append','writepid','ca','cert','key','pkcs12'}:
        sys.exit('VPN profile needs unattended-worker review: '+words[0])
PY
if ((check)); then exit 0; fi
p0_require_config SECRET_ID AWS_REGION
tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT
p0_aws secretsmanager get-secret-value --secret-id "$SECRET_ID" \
    --query SecretString --output text >"$tmp/current.json"
jq --rawfile vpn "$1" '. + {ULAB_VPN_CONFIG:$vpn}' "$tmp/current.json" >"$tmp/updated.json"
p0_aws secretsmanager put-secret-value --secret-id "$SECRET_ID" \
    --secret-string "file://$tmp/updated.json" >/dev/null
echo 'VPN profile saved to the existing worker secret.'
