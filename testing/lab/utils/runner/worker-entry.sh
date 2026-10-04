#!/usr/bin/env bash
# Runs only on disposable EC2 workers, under cloud-init as root.
set -Eeuo pipefail
umask 077
export HOME=/root DEBIAN_FRONTEND=noninteractive
work=/opt/ukama-runner
export AWS_DEFAULT_REGION="$ULAB_RUNNER_REGION" AWS_PAGER=""
input="s3://$ULAB_RUNNER_BUCKET/inputs/$ULAB_RUNNER_BATCH"

aws s3 cp "$input/$ULAB_RUNNER_WORKER.json" "$work/worker.json" --only-show-errors

mapfile -t extra_packages < <(python3 - "$work/worker.json" <<'PY'
import json, sys
print('\n'.join(json.load(open(sys.argv[1]))['config']['EXTRA_APT_PACKAGES']))
PY
)
packages=(podman netavark aardvark-dns crun uidmap dbus-user-session
          openvpn dnsmasq-base iproute2 iptables jq curl ca-certificates
          python3 python3-yaml git make gcc g++ pkg-config rsync unzip)
for package in "${extra_packages[@]}"; do
    [[ -n "$package" ]] && packages+=("$package")
done
printf 'UKAMA_RUNNER_BOOTSTRAP stage=worker-packages\n'
printf 'Worker setup: installing Podman, VPN, DNS and build dependencies\n'
apt-get -o DPkg::Lock::Timeout=300 -o Acquire::Retries=3 install -y "${packages[@]}"
printf 'Worker setup: dependency installation complete\n'
python3 "$work/bootstrap/worker.py"
