# Local and AWS scenario execution

This replacement is based on the clean `ukama-lab.tar(20261003-163604).gz`
baseline. It changes `utils/run-scenarios.sh` and adds the files in this
directory plus `tests/aws_runner/test_runner.py`. It makes no changes to the
scenario YAML files, C implementation, or `scripts/` directory.

## Install

Extract the drop-in archive from the lab directory:

```bash
tar -xzf /path/to/ukama-lab-local-aws-drop-in.tar.gz
```

Local mode remains the default. It does not require AWS, a VPN profile file,
or the AWS JSON configuration:

```bash
./utils/run-scenarios.sh p0
./utils/run-scenarios.sh p0 --mode local
```

Both commands use the existing selection, sequential execution, command
arguments, scenario reports, exit codes, and batch report logic. Your local
environment exports remain applicable.

## Configure AWS once

```bash
# First installation only; keep an existing aws.json.
cp -n utils/runner/aws.example.json utils/runner/aws.json
chmod 600 utils/runner/aws.json
```

Edit `VPN_CONFIG_FILE` to point to your working OpenVPN `.ovpn` profile.
`REGION` defaults to `us-east-1`. Paths in the JSON may be absolute, use `~`,
or be relative to the JSON file. Keep the real configuration untracked.

Continue exporting your existing `UKAMA_IDENTIFIER`, `UKAMA_PASSWORD`, service
URLs, `ULAB_*` scenario settings, `UKAMA_REPO`, and software-update versions.
The runner copies scenario environment values to workers. It rewrites the
repository, binary, and results paths to worker locations. Optional
`ENV` entries in the JSON override exported scenario settings; `{}` uses your
current exports. It does not copy your AWS access keys to workers.

The controller requires Linux x86_64, Python 3.10+, Bash, the AWS CLI, curl,
ldd, and access to the local Podman account containing your starter images.
Run it as the same user you normally use for ukama-lab/Podman. Use your normal
AWS CLI authentication. The AWS identity must be allowed to deploy the
included CloudFormation template, create/pass its EC2 role, launch and
terminate instances, read EC2 console output and regional S3 prefix lists
(`ec2:DescribePrefixLists`), and read/write its S3 bucket. The template is available
for review in `infrastructure.json`.

Your existing compiled `bin/ukama-lab` must work locally. The supplied source
archive did not include that binary or the external Ukama repository; this
runner packages them from your working machine when you launch a batch.

AWS mode packages the same kubectl executable used locally and sets
`ULAB_KUBECTL` to its installed worker path. The existing `stop-node.sh` then
deletes only the mesh pod matching that node's factory ID, using the existing
VPN. Keep your normal `export ULAB_KUBECTL="$HOME/kubectl"`; if unset, the runner
looks for `$HOME/kubectl` and then `kubectl` on PATH. Linux x86_64 binaries and
self-contained executable scripts are supported.

If your local kubectl uses an existing kubeconfig, its current context is
exported locally with `config view --raw --flatten --minify` and transferred
privately with the runtime. No Kubernetes login, authentication probe, or
cluster-readiness gate is added. When no kubeconfig exists, none is required.
Optional `KUBECTL_FILE` and `KUBECONFIG_FILE` in `aws.json` select explicit local
paths; normally the existing tool and configuration are discovered automatically.
Keep any existing exec-plugin dependencies available if your config uses them;
the runner does not create new credentials. Local mode is unchanged.

AWS prints timestamped instance creation/state changes, observed setup stages,
scenario starts, result collection, and termination requests. It confirms EC2
termination for up to 60 seconds after collection; an unconfirmed termination
is reported honestly without changing scenario verdicts. These messages are
saved in the batch's `aws-lifecycle.log`. The final interactive terminal summary
uses green PASS and red FAIL; saved reports remain plain text. Requested and
actual worker counts are shown separately (one selected scenario uses one worker).

Each FAIL in the final AWS summary includes its recorded failure reason:

```text
FAIL worker-001 console/nodes/tc-024-list-detail.yaml "14:32:58 ERROR runtime infra cleanup had 3 failed step(s)"
```

Failed checks/events come from the scenario report; cleanup and setup errors
come from its log. The final totals line and diagnostic-collection errors are
excluded. Multiple recorded failures are separated by semicolons. The same
text is saved as `failure_reason` on each failed result in `batch-report.json`.
Missing details are stated explicitly. Scenario verdicts and original artifacts
are unchanged; PASS stays green and FAIL red in interactive terminal output.

## Run

Start with one AWS worker to establish the environment against your VPN:

```bash
./utils/run-scenarios.sh p0 --mode aws --workers 1
```

Then use the desired number of workers:

```bash
./utils/run-scenarios.sh p0 --mode aws --workers 10
./utils/run-scenarios.sh p0/billing --mode aws --workers 3
./utils/run-scenarios.sh resilience --mode aws --workers 10
./utils/run-scenarios.sh p0 --mode aws --workers 10 --list
```

`--list` only lists the existing scenario selection; it does not contact AWS.
`--aws-config /path/to/aws.json` selects another configuration.
`--scenario-list FILE` retains its existing meaning. For AWS, selected files
must be within the lab directory and `SCENARIO_ROOT` so their paths can be
reproduced on workers. Set `SCENARIO_ROOT=scenarios` for a mixed-suite list.

No separate build-artifacts or AMI-refresh command is required.

## Consolidated correction (3 October 2026)

Use the full `ukama-lab-local-aws-drop-in.tar.gz` replacement. It includes all
previous bootstrap, VPN and kubectl fixes. Keep your existing `aws.json`.
No patch ordering, AMI refresh or manual cache deletion is required.

The latest supplied batch, `20261003t194607z-2b9da4b2`, stopped before any
scenario ran. All 20 workers failed loading the image archive. The exported
`.tar.gz` was 93 bytes and expanded to zero bytes: the previous exporter gave
the source tar and compressed output the same temporary filename. Opening
the output truncated the source. A checksum of that broken gzip could still
match, so download checksum checking alone did not catch it. This was a
runner defect. VPN route-up completed in that batch; later connectivity
checks had not yet run.

This replacement:

- Uses separate temporary tar and gzip files, with cleanup on failure.
- Validates the exported manifest, expected image tags and configuration IDs,
  and every referenced layer's checksum before creating AWS workers.
- Verifies the compressed archive expands to the exact original tar bytes.
- Discards the old unvalidated cache automatically and exports your existing
  images again. It does not rebuild starter images that are already present.
- Checks compressed and expanded checksums and disk space on each worker,
  loads a plain tar with Podman, then verifies all loaded image IDs.
- Starts worker-001 first and launches the remaining workers only after its
  existing infrastructure checks succeed and it publishes `ready.json`.
  This is not a scenario PASS gate. The first worker runs only its assigned
  scenarios, and a scenario FAIL does not prevent infrastructure admission.
- Records startup stages and actual Podman error output, and serializes
  heartbeat/readiness status publication to avoid concurrent file replacement.

Apply from your lab directory, then start a new batch:

```bash
tar -xzf /path/to/ukama-lab-local-aws-drop-in.tar.gz
./utils/run-scenarios.sh p0 --mode aws --workers 20
```

The failed batch cannot be repaired with `--resume`: its workers already
received the old code and empty archive. If old instances remain, collect
and terminate only that batch with:

```bash
./utils/run-scenarios.sh p0 --mode aws --cleanup runs/p0-batches/20261003t194607z-2b9da4b2
```

## What happens

1. The existing script selects the scenarios once. AWS mode partitions that
   list across the requested workers, capped at the number of selected files.
2. It reuses `localhost/testing/virtualnode-base:anode-starter`,
   `cnode-starter`, and `tnode-starter` from your local Podman store. A missing
   base is built by the existing `scripts/build-node.sh`, using an available
   Factory node ID. This lookup does not reserve or provision the node. Your
   laptop needs backend access when a missing image needs that lookup.
3. It exports the starter images and probe image together. Existing UE/media
   images are also included when present. The export is cached locally under
   the runs directory and in S3 by image content. The first export checks the
   manifest, tag/configuration identities and layer contents. Cache reuse checks
   the compressed checksum and a complete decompression round-trip. Unchanged
   exports are reused on subsequent batches.
4. It packages the current lab code, its working native executable and linked
   libraries, and the external repository build contexts. By default those
   contexts are `testing/node` and `testing/ue`, the paths used by the current
   scripts. `REPO_PATHS` can include additional relative paths or `"."` for the
   complete repository. Build output outside those contexts is unnecessary
   for stamping nodes from transferred starter bases.
5. CloudFormation creates or reuses a dedicated VPC/subnet with Internet
   access, an outbound-only security group, an S3 gateway endpoint, a private
   encrypted bucket, and an EC2 instance role. Instances have public IPv4
   addresses for outgoing Internet/VPN connections. There is no inbound SSH
   requirement and no SSH key-pair dependency.
6. EC2 starts the first Ubuntu 24.04 x86_64 worker. With an empty `AMI_ID`, the
   controller resolves Canonical's image through SSM and records the exact
   AMI in the batch state. `AMI_ID` can pin a compatible Ubuntu 24.04 image.
   Host tools are installed by cloud-init; a custom AMI bake is not required.
7. Each worker configures its host, connects VPN, loads the same starter images
   into root's local Podman store, and uses the system systemd manager. Root
   is supported by the existing node-host-control script. Startup checks the
   lab binary, systemd/Podman availability, and
   host/container DNS and HTTP(S) connectivity. The worker publishes a durable
   readiness marker through S3. Only then does the controller launch the rest
   of the requested workers; each performs these same infrastructure checks.
   The controller stops expansion if the first worker has an infrastructure
   failure. No extra scenario is selected or run.
8. The worker runs `run-scenarios.sh --mode local --scenario-list ...` with a
   unique batch ID. It runs one complete scenario at a time. All nodes and
   UEs belonging to that scenario stay on that worker. Existing scripts stamp
   node IDs from the cached bases and perform normal setup and cleanup.
9. After each scenario, its original reports, logs and progress are uploaded.
   A worker completion marker is published only after its final result upload
   succeeds. The worker then shuts down; EC2 is configured to terminate on
   shutdown. The controller also terminates workers whose completion marker
   it observes and collects the original reports into the local batch folder.

Starter images remain available for all scenarios on a worker. The existing
UE/media scripts reuse their images when present and build when absent or
explicitly forced. Their behavior is unchanged. Exporting images does not
transfer all intermediate Podman build cache.

## VPN and DNS

This implementation supports noninteractive OpenVPN profiles. It creates
one tunnel per worker. The worker copy embeds certificate/key references and
replaces laptop-specific up/down hooks. The original profile is not edited.
`VPN_AUTH_FILE` supplies a two-line username/password file when required.
Interactive MFA, encrypted keys needing a prompt, nested configuration files,
and VPN products other than OpenVPN require a compatible worker profile.

For servers allowing one connection per client identity, supply distinct
profiles in `VPN_CONFIG_FILES`, one per worker. A single `VPN_CONFIG_FILE` is
reused when your VPN server permits concurrent sessions with that identity.

The worker preserves IPv4 routing from your profile and the VPN server. Both
full-tunnel `redirect-gateway` and pushed backend-specific routes are accepted.
`VPN_ROUTES` can supply additional backend CIDRs when needed. Before VPN
connection and again on reconnect, the worker pins regional S3, VPC, instance
metadata and AWS DNS traffic to its original AWS gateway. The regional S3
prefixes come from AWS's S3 prefix list, not hard-coded addresses. These more
specific routes keep image downloads, reporting and instance-role credentials
on the AWS path even with the full-tunnel VPN active. The Windows-only pushed
`block-outside-dns` option is ignored on Linux.

The worker runs a small DNS forwarder. Queries for `VPN_DOMAINS` go to DNS
servers supplied by OpenVPN or explicitly listed in `VPN_DNS_SERVERS`.
Other domains use the AWS resolver. Containers use this same forwarder at
the worker's reachable private address, not a container loopback address.
The default backend DNS domain is `udev.ukama.com`.

`VPC_CIDR` and `PODMAN_CIDR` must not overlap specific backend VPN routes.
Full-tunnel default routes are permitted. Startup rejects other overlaps with
a diagnostic instead of starting scenarios on broken
routing. Container traffic is masqueraded into the tunnel, with TCP MSS
clamping. S3/instance metadata remain on the AWS path. No network namespace
wrapper is placed around ukama-lab or its systemd-managed containers.

Connectivity checks use GET requests to endpoint roots and accept any HTTP
status. They check DNS, TCP and TLS, not application responses, analytics
data, KPIs, or scenario verdicts. No application setup or gate is added.
TLS verification stays enabled. `EXTRA_PROBE_URLS` can name additional
HTTP(S) endpoints. With `NETWORK_MODE=direct`, the managed VPC must already
be able to reach the selected endpoints without a client VPN.

## Reports and lifecycle

The controller prints the batch directory, normally
`runs/p0-batches/<batch-id>/`. It contains:

| Path | Contents |
| --- | --- |
| `state.json` | Account, AMI, instance IDs and exact assignments; no scenario credentials |
| `batch-report.json`, `batch-report.txt` | Combined outcomes and unfinished assignments |
| `workers/worker-NNN/batch/...` | Original worker batch reports, scenario reports and logs |
| `workers/worker-NNN/infrastructure/` | Routes, DNS, VPN journal, Podman and disk diagnostics |
| `workers/worker-NNN/bootstrap.log`, `worker.log`, `image-load.log` | Startup, sequential execution and Podman image-load logs |
| `workers/worker-NNN/ready.json` | Durable infrastructure readiness, independent of scenario verdicts |
| `worker-NNN-console.log`, `worker-NNN-console.json` | EC2 console diagnostics collected by the controller, including failures before S3 becomes usable |

A scenario failure remains a scenario failure. A worker interruption or
failure before execution is an infrastructure error; unrun assignments are
listed as unfinished. Neither infrastructure failure nor `--resume` causes
automatic scenario retries. Original report files are not rewritten; only
paths in the combined report are mapped to downloaded files.

Ctrl-C detaches the controller. Existing workers continue, upload their
results, and terminate on completion. Reattach or explicitly stop a batch:

```bash
./utils/run-scenarios.sh p0 --mode aws --resume runs/p0-batches/BATCH_ID
./utils/run-scenarios.sh p0 --mode aws --cleanup runs/p0-batches/BATCH_ID
```

Resume requires only the saved state and AWS access. It monitors existing
workers; it does not launch previously unlaunched workers after a partial
launch failure. If you detach while the first worker is being prepared, the
remaining workers are not launched; keep the controller attached through
startup to create the full requested pool. Cleanup downloads what has already reached S3, then requests
termination of that batch's remaining instances. Active scenarios can be
interrupted by explicit cleanup, so use resume for normal collection.

`--fail-fast` preserves the existing nonzero-exit trigger. Workers finish
their current scenarios and cleanup, then stop before the next assignment
after observing the stop request. Simultaneous completions can occur before
all workers see that request.

Optional `--factory-nodes auto/N` preparation runs once on the controller.
Workers disable repeated preparation. Default preparation remains off.
`--prepare-only` launches no EC2 instances. Runtime node allocation remains
inside the existing scenario scripts.

`BOOT_TIMEOUT_MINUTES` bounds initial preparation while attached.
`MAX_WORKER_HOURS` schedules shutdown independently on each worker before
package installation starts. `UPLOAD_RETRY_MINUTES` controls result-upload
retries; an unsuccessful upload does not publish a successful completion.
The maximum lifetime is a hard limit and can interrupt a long scenario.
`INSTANCE_TYPE`, `DISK_GB`, and `MIN_FREE_DISK_GB` are configurable for the
largest scenario you intend to run.

The VPC, role and bucket are reused across batches. EC2 instances and root
disks are disposable. S3 inputs expire after two days, results/control files
after `RESULT_RETENTION_DAYS` (default 30), and image cache after 90 days.
The bucket is retained if the CloudFormation stack is deleted. Inputs include
the credentials needed by scenarios; S3 is private and encrypted and the
worker role is restricted to this runner bucket. Code packaging excludes
common credential locations, VPN profiles, `.env*`, `.pem`, `.key`, local
results, and the laptop `setup` file. Avoid embedding secrets in other source
files. Original verbose logs can contain application request details and are
stored with the same private access policy.

## Validation included

```bash
python3 -m unittest discover -s tests/aws_runner -v
```

106 offline tests passed. They cover local-default behavior, exact selection,
fail-fast, worker interruption and upload failure, unique assignments,
launch configuration, archive safety, executable/library packaging, VPN/DNS
configuration, result collection, packaged kubectl cleanup, lifecycle logging,
termination confirmation, terminal summary colors, scenario failure reasons, execution
of generated bootstrap Bash with isolated command stubs, installer failures,
and detection of startup failures without a worker heartbeat. VPN-hook tests
execute the actual Python entry point with a restricted PATH and harmless
kernel-command fixtures, checking sbin command lookup, readiness/error files,
full/split-tunnel routes, overlap rejection, and AWS route installation.
Image pipeline tests exercise real tar/gzip bytes and the complete export,
cache and worker-decompression functions. They reproduce the exact temporary
path collision, test automatic replacement of that broken cache, verify all
image tags/configurations/layers, and reject truncated or corrupted inputs
before Podman is invoked. Admission tests verify that no second worker
launches before infrastructure readiness, that early infrastructure failures
stop expansion, and that scenario failures do not become admission gates.
External AWS/Podman/system services remain substituted in these offline tests.
The original and updated scripts list
the same 146 P0 scenarios and 321 resilience scenarios. The absent
`resilience-hooks` tree produces the same error in both versions.

Bash syntax and Python compilation checks passed. No EC2 deployment, actual
Podman/VPN networking, or live scenario run was performed in the delivery
environment. The first real AWS batch validates those environment-specific
parts; offline tests do not establish the cloud PASS rate.

The official AWS CLI installer was additionally downloaded and executed in an
isolated local directory; the installed `aws --version` returned CLI 2.37.9.
This verifies the installer path, not EC2/VPN execution.

Implementation references:

- https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/user-data.html
- https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html
- https://ubuntu.com/aws/docs/aws-how-to/instances/find-ubuntu-images/
- https://docs.podman.io/en/stable/markdown/podman-save.1.html
- https://github.com/containers/common/blob/main/docs/containers.conf.5.md
- https://openvpn.net/community-docs/community-articles/openvpn-2-6-manual.html
