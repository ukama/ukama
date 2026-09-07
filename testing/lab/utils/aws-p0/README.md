# Direct-VPN P0 runner

Extract the replacement archive **from the lab repository root**:

```bash
tar -xzf ~/Downloads/aws-p0-direct-vpn.tar.gz
cd utils/aws-p0
./run.sh --dry-run
./run.sh --vpn ~/Downloads/downloaded-dev-cluster-client-config.ovpn
```

Run the controller as your normal user, using your existing AWS CLI identity.
Keep `UKAMA_REPO` set to your local Ukama source checkout and build
`bin/ukama-lab` before starting. Keep your laptop VPN connected if using local
factory preparation. OpenVPN runs as root on the EC2 workers automatically.

The archive contains full replacement files; it does not contain or replace
`config.env`, `credentials.env`, `.state.env`, your VPN profile, or your source
checkout. Existing Ubuntu worker AMIs work if their normal lab dependencies are
installed; a worker installs OpenVPN if necessary. New AMIs built with
`install-worker-base.sh` include it. No gateway setup or new AMI is required just
to use direct VPN.

After extracting over an older runner, remove the obsolete gateway-only files:

```bash
rm -f utils/aws-p0/setup-backend-gateway.sh \
      utils/aws-p0/worker-pre-run-final.sh
```

`diagnose-p0-network.sh` may also be removed if it was added solely for the old
middle-node investigation. Old `BACKEND_GATEWAY_*`, `BACKEND_ROUTE_CIDR`,
`BACKEND_TEST_IP`, `MAX_WORKERS`, and `DEFAULT_WORKERS` entries in `config.env`
are ignored and can be deleted. Set `P0_MAX_WORKERS=28` only if you want this
checkout's current 28-worker plan to be a hard ceiling.

## Assignment

The default is one EC2 instance per directory containing selected scenario YAML
files. Scenarios in a directory run sequentially; all directory workers launch
without waiting for other workers to finish. Startup times differ; there is no
shared start barrier. Directories containing a deployment-wide scenario matched
by `exclusive.txt` are merged onto one final sequential worker. The uploaded
source contains **146 scenario files in 30 directory groups**; three protected
groups collapse to one, producing **28 workers**. Your local checkout may differ;
`--dry-run` prints the exact assignment without AWS calls or factory writes.

```bash
# Two subcategories, each on its own instance:
./run.sh --vpn ~/Downloads/downloaded-dev-cluster-client-config.ovpn billing/payment usage/accounting

# All groups, with an explicit maximum of 28 instances for this checkout:
./run.sh --workers 28

# Optional: one instance per non-exclusive scenario; protected scenarios remain
# together on one final worker (131 workers in the supplied checkout):
./run.sh --group-by scenario --workers 200
```

`--workers` is a ceiling, not a request to redistribute scenarios. The default
ceiling is 100, configurable through `P0_MAX_WORKERS`. If the planned worker count
exceeds the ceiling, planning fails before launching anything. `exclusive.txt`
retains its original shell-style path patterns; any matched directory is kept on
the final sequential worker. Legacy `MAX_WORKERS`, `DEFAULT_WORKERS`, and
`BACKEND_GATEWAY_*` settings are unused. Top-level categories, nested
subcategories, `--scenario PATH` and `--scenario-list FILE` all remain supported.

## VPN and credentials

The first `--vpn FILE` stores the self-contained inline CA/certificate/key profile
as `ULAB_VPN_CONFIG` in your existing `SECRET_ID`, preserving backend credentials.
Later runs can omit `--vpn`. Updating it requires your local AWS identity to have
`secretsmanager:GetSecretValue` and `secretsmanager:PutSecretValue` permission on
that secret. Workers use the existing instance profile to read the secret and
access the configured S3 prefix. `setup.sh` now preserves the stored VPN profile
when refreshing backend credentials; you do not need to rerun setup for an
already configured installation.

Each worker downloads input first, preserves the original host route to EC2
metadata at `169.254.169.254/32`, and opens its own OpenVPN connection. It then
checks fresh IMDSv2 instance-role credentials and an authenticated S3 upload
before running any scenario. This addresses the likely metadata-routing cause
of the diagnostic's “Unable to locate credentials”; the live checks verify the
actual result rather than assuming the cause. Credentials remain managed by the
AWS CLI and can refresh throughout long runs.

Pushed VPN DNS servers are registered with systemd-resolved for
`udev.ukama.com`. PAUTH, BFF, factory, warehouse and node-gateway HTTP reachability
are checked before scenarios start. A received HTTP error still proves transport
reachability; scenario assertions determine application correctness. Workers
never configure a middle gateway or route all `10.0.0.0/8` through another EC2.

The VPN profile stays in the existing secret and a mode-0600 temporary worker
file, outside results. It is not embedded in EC2 user-data or output archives.
Workers disconnect and remove the temporary file before final result upload.
Source packaging excludes `*.ovpn` and the supplied profile's basename.

## Results and cleanup

The controller prints its batch ID and continuously displays worker status.
Workers periodically copy status and partial results to S3. Final results include
scenario reports, logs, and a VPN log. A completion marker is written only after
the batch report accounts for every assigned scenario exactly once and the final
archive uploads successfully. Product failures are distinct from incomplete
workers. Collected reports include explicit `MISSING` results for unfinished
scenarios and return a nonzero exit code on product or infrastructure failures.

```bash
./run.sh --status BATCH_ID
./run.sh --resume BATCH_ID
./run.sh --collect BATCH_ID
./run.sh --cleanup BATCH_ID
```

Local output is under `runs/p0-aws/BATCH_ID`: `summary.txt`, `batch-report.json`,
`combined.tsv`, `infrastructure-failures.tsv`, and `workers/`. Input group lists
are under `input/groups.tsv` and `input/shards/`.

Workers terminate themselves after upload and have an independent shutdown
deadline. The controller terminates tagged batch instances on completion,
interruption or launch failure. AWS cleanup errors are reported as failures and
print a recovery command. If your laptop disappears abruptly, workers still
finish or expire; `--collect` retrieves persisted output later. S3 results are
retained unless successful collection and `DELETE_S3_AFTER_COLLECT=true` permit
removal. There are no automatic scenario retries.

## Validation scope

Local checks exercised the real scenario discovery, selectors and limits;
controller launch/collection and cleanup, including a lost launch response;
worker VPN/S3 gating, missing metadata credentials, scenario failure, missing
reports, incomplete reports, failed archive upload, profile removal and merged
report outcomes using simulated AWS/network commands. Shell syntax and Python
compilation were checked. No EC2 instances or full lab scenarios were run from
the development environment.

Your diagnostic established simultaneous backend reachability for two EC2 VPN
clients plus your laptop for five minutes. Full group capacity, virtual-node
traffic and backend load still require your planned live run. EC2 isolation does
not isolate shared backend mutations, such as global software-release promotion;
concurrently selected scenarios must use compatible backend state.

Implementation references: [AWS instance-role credentials](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/instance-metadata-security-credentials.html),
[OpenVPN options](https://openvpn.net/community-docs/community-articles/openvpn-2-6-manual.html).
