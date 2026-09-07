# Distributed scenario runner

Run from the ukama-lab repository root:

```bash
./utils/runner/run.sh p0 --workers 20
./utils/runner/run.sh resilience --workers 20
./utils/runner/run.sh resilience --workers 4 allocation policy
./utils/runner/run.sh resilience --workers 4 --scenario-list resilience.txt
./utils/runner/run.sh resilience --workers 4 --dry-run
```

List files contain one repository-relative YAML path per line, for example
`scenarios/resilience/policy/r058-idle-session-must-not-consume-payload.yaml`.
Use repeated `--scenario PATH` options to select individual files instead.
All selected files must belong to the selected suite. Omit categories or use
`all` to select the complete suite. Existing commands without a suite use P0.

The controller packages the current source trees, creates scenario shards,
and launches disposable EC2 workers in parallel. Each worker invokes
`utils/run-scenarios.sh` with the selected suite and runs its shard sequentially.
Status and results are transferred through S3; the controller collects reports
and terminates remaining workers. `--workers` is a maximum, capped by available
work and the existing `MAX_WORKERS` configuration.

Keep the existing `config.env`, `credentials.env`, `.state.env`, AMI, network,
and AWS resources. Existing P0-named configuration variables and resource tags
remain compatible. The selected suite is recorded in each batch manifest and
worker environment. Generated batch IDs include the suite.

`--dry-run` creates the local plan and source archives without preparing factory
nodes, uploading files, or launching instances. AWS credentials are still needed
for the existing configuration checks.

```bash
./utils/runner/run.sh --status BATCH_ID
./utils/runner/run.sh --resume BATCH_ID
./utils/runner/run.sh --collect BATCH_ID
./utils/runner/run.sh --cleanup BATCH_ID
```

Patterns in `exclusive.txt` group shared-state scenarios onto one worker, where
they run sequentially. That worker runs alongside the other workers; these
patterns do not provide batch-wide or cross-batch isolation. Run scenarios that
require the entire backend to be idle in a separate batch after other runs end.
Do not run P0 and resilience software-update batches concurrently against the
same release catalog.
