# Direct-VPN scenario runner

From the lab directory:

```bash
./utils/runner/run.sh p0 --dry-run
./utils/runner/run.sh resilience --workers 30 --dry-run
./utils/runner/run.sh resilience --workers 30 --vpn ~/Downloads/downloaded-dev-cluster-client-config.ovpn
./utils/runner/run.sh resilience allocation policy --workers 30
./utils/runner/run.sh p0 console/analytics --workers 30
```

The suite defaults to P0 for existing commands. Select `p0` or `resilience`
explicitly. Omit categories, or use `all`, to select the entire suite. Repeated
`--scenario PATH` options or `--scenario-list FILE` select exact YAML files;
paths must be repository-relative or absolute and belong to the selected suite.

Main's scheduling is retained: one worker per selected scenario directory,
with sequential execution within each worker and parallel EC2 workers.
`--group-by scenario` assigns one worker per scenario instead. Groups containing
matches from `exclusive.txt` are combined on one sequential worker. This worker
still runs alongside other workers and does not provide backend-wide isolation.

`--workers` is a ceiling, not a target or a round-robin shard count. If the plan
requires more workers, selection fails before launch. The default ceiling is
`P0_MAX_WORKERS`, or 100. Main's behavior ignores the older `MAX_WORKERS` and
`DEFAULT_WORKERS` settings. Preview the plan before a full run.

Keep the existing `config.env`, `credentials.env`, `.state.env`, AWS identity,
and worker AMI. Set `UKAMA_REPO` to the local Ukama source tree and build
`bin/ukama-lab`. Each worker establishes direct OpenVPN connectivity and verifies
S3 access before executing `utils/run-scenarios.sh` with the selected suite.
The supplied VPN profile is stored in the existing worker secret; omit `--vpn`
once that secret contains the profile. Keep the laptop VPN connected for local
factory preparation. No gateway setup is used.

`--dry-run` creates only the local plan. It does not prepare factory nodes,
package sources, call AWS, update secrets, or launch instances. Existing local
configuration and source prerequisites are still checked.

```bash
./utils/runner/run.sh --status BATCH_ID
./utils/runner/run.sh --resume BATCH_ID
./utils/runner/run.sh --collect BATCH_ID
./utils/runner/run.sh --cleanup BATCH_ID
```

The controller collects results and terminates remaining workers. Worker and
combined reports reconcile results against the submitted shards; missing,
duplicate, or unassigned results cannot produce a successful batch. Both suites
use the same report schema. No scenario retries are added.

This replacement contains tracked runner code only. It excludes local config,
credentials, state and VPN profiles. Remove the obsolete
`setup-backend-gateway.sh` and `worker-pre-run-final.sh` while resolving the merge.
