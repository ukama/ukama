# Distributed P0 runner

This directory implements a deliberately small distributed test runner:

- the local machine is the controller;
- S3 stores input and output files;
- each EC2 instance receives one text file of scenario paths;
- each worker runs its scenarios sequentially;
- workers upload status and results, then terminate;
- the local controller displays status and merges the reports.

There is no queue, scheduler, database, daemon, web service, Lambda, Batch,
or cross-worker communication.

## Files

- `run.sh` — local controller and terminal monitor
- `worker.sh` — disposable EC2 worker
- `report.sh` — local report merger
- `setup.sh` — one-time S3/IAM/secret/security-group setup
- `cleanup.sh` — terminate a batch by ID
- `install-worker-base.sh` — install generic worker packages on a builder
- `check-worker-ami.sh` — verify the builder before AMI capture
- `capture-ami.sh` — capture a prepared EC2 instance as the worker AMI
- `worker-pre-run.sh` — small, editable worker preparation hook
- `exclusive.txt` — scenarios kept on one sequential worker
- `config.env` — non-secret configuration, created from the example
- `credentials.env` — local credentials, stored in Secrets Manager by setup
- `.state.env` — IDs discovered or created by setup

## Worker AMI

The worker AMI is the one reusable machine image. Prepare one EC2 instance
that can run a normal ukama-lab scenario, then capture it.

The image must provide at least:

```text
bash
aws CLI
jq
curl
tar
gzip
python3
podman
Open vSwitch and ovs-vsctl
systemd
```

It must also contain any stable build/runtime dependencies required by
`$UKAMA_REPO/testing/node/mk_local_vnode.sh` and the Ukama virtual-node
containers. Do not put Ukama credentials, AWS access keys, source code, or
scenario results in the AMI.

On a fresh builder, the generic packages can be installed with:

```bash
sudo ./utils/aws-p0/install-worker-base.sh
./utils/aws-p0/check-worker-ami.sh
```

The generic installer cannot know every dependency used by the current Ukama
`testing/node` tree. Install those project-specific dependencies and validate
one normal scenario before capture:

```bash
./utils/aws-p0/capture-ami.sh i-0123456789abcdef0
```

The script writes the new `AMI_ID` to `.state.env`.

## One-time setup

```bash
cd /path/to/ukama-lab

cp utils/aws-p0/config.env.example utils/aws-p0/config.env
cp utils/aws-p0/credentials.env.example utils/aws-p0/credentials.env
chmod 600 utils/aws-p0/credentials.env

$EDITOR utils/aws-p0/config.env
$EDITOR utils/aws-p0/credentials.env

./utils/aws-p0/setup.sh
```

`setup.sh` creates or updates:

- one private S3 bucket/prefix;
- one Secrets Manager secret;
- one EC2 IAM role and instance profile;
- one outbound-only security group;
- `.state.env` containing resolved IDs.

AWS access uses the normal AWS CLI credential chain. `AWS_PROFILE` may be set
in `config.env`. No AWS access key is copied to a worker.

## Run all P0 scenarios

Build `bin/ukama-lab` locally first. The controller packages the current local
working trees, including modified and untracked files.

```bash
cd /path/to/ukama-lab

export UKAMA_REPO=/home/kashif/work/ukama/repos/ukama

./utils/aws-p0/run.sh --workers 10
```

The command remains attached and redraws a status table until all workers
finish. It then downloads the worker archives, prints one combined report,
and terminates any worker still alive.

Run selected categories:

```bash
./utils/aws-p0/run.sh --workers 6 console usage software-update
```

Preview packaging and shard assignment without uploading or launching:

```bash
./utils/aws-p0/run.sh --workers 10 --dry-run
```

## Reconnect to a batch

A worker uploads status to S3, so closing the local terminal does not lose the
batch state.

Display once:

```bash
./utils/aws-p0/run.sh --status 20260801t230000z
```

Resume monitoring, collection, reporting, and cleanup:

```bash
./utils/aws-p0/run.sh --resume 20260801t230000z
```

Download and report without waiting:

```bash
./utils/aws-p0/run.sh --collect 20260801t230000z
```

Terminate remaining workers:

```bash
./utils/aws-p0/run.sh --cleanup 20260801t230000z
```

## Local result layout

```text
runs/p0-aws/<batch-id>/
├── input/
├── instances.tsv
├── raw-output/
├── workers/
├── batch-report.json
├── combined.tsv
├── failed.txt
├── infrastructure-failures.tsv
└── summary.txt
```

## Cost controls

Every worker is launched with:

- instance-initiated shutdown behavior set to `terminate`;
- termination protection disabled;
- an encrypted root volume with `DeleteOnTermination=true`;
- a worker-side hard lifetime watchdog;
- a local-controller cleanup trap.

The S3 lifecycle configured by `S3_RETENTION_DAYS` deletes old batch objects.
Set `DELETE_S3_AFTER_COLLECT=true` to remove a batch immediately after local
collection.

## Scenario distribution

Normal scenarios are sorted and distributed round-robin. Shell patterns in
`exclusive.txt` are placed together on one final worker and run sequentially.
This is intended for deployment-wide failure controls or other scenarios that
must not overlap ordinary tests.

There is intentionally no automatic retry. A product failure remains a
failure. An EC2/bootstrap failure appears under `infrastructure-failures.tsv`.
