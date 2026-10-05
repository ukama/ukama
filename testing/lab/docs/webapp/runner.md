# Local browser scenario runner — Patch 4

The C lab runner now launches one persistent Playwright worker per active v2
scenario. It runs each phase's events and checks in order, then final checks,
then bounded cleanup. Checks observe the current page. Only an explicit
`web_reload` event reloads it. All components run locally on the host.

Patch 4 also executes UI network/site provisioning and local virtual-node
runtime setup. The two network examples are active. See
[provisioning.md](provisioning.md) for host prerequisites, the required console
unchanged-console constraints, creation receipts, component selection and operational runs.

## Run the first scenario

Apply Patches 1 and 2 first and rebuild the lab with your normal local build.
The Makefile already includes `src/*.c`, so it picks up the three new C modules.
Install/build the worker and save an authenticated **owner** session as described
in [worker.md](worker.md). Start your local console and its auth/BFF dependencies.
The session must have permission to see the Members invitation control.

From the lab directory:

```sh
npm --prefix adapters/webapp ci
npm --prefix adapters/webapp run build
npm --prefix adapters/webapp run install:browser

./bin/ukama-lab lint scenarios/webapp/p0/session/wb-000-authenticated-members.yaml
./bin/ukama-lab run scenarios/webapp/p0/session/wb-000-authenticated-members.yaml \
  --repo /path/to/ukama --out runs
```

The example uses `http://localhost:3000` and `.auth/owner.json`. Edit a local copy
for another URL/auth-state path. `${ENV_VAR}` substitution is supported by the
scenario parser if you prefer explicit environment configuration. Relative auth
paths resolve from the lab invocation directory. Invoke from the lab root, or
pass an absolute `--webapp-worker /path/to/ukama-lab/utils/webapp-worker.sh`.
`UKAMA_LAB_WEBAPP_WORKER` sets its default. The launcher must accept `worker` as
its first argument, write only protocol JSONL to stdout, and log to stderr.
The existing CLI still requires `--repo`, even for this read-only journey.

The scenario opens Business → Members, checks that **Invite member** is visible
and enabled, explicitly reloads, and checks again. It sends no invitation.
It contributes only a partial smoke check for `WEB-TEAM-002`; the complete
invitation/role requirement remains planned and live verification remains not run.

## Result and artifact contract

The C runner is the authority for the scenario result. Worker success cannot
hide an earlier event failure, a timeout, cancellation, or cleanup failure.
V2 `xfail` is rejected. `wip`/`skip` create a report with `outcome: SKIP`,
`passed: false`, zero checks, and no browser/backend execution.

Each active run has a new, private directory. A supplied `--run-id` must start
with a letter/digit and contain 1..80 letters/digits/hyphens/underscores. An
existing directory is rejected, including after an interrupted run. There is
no automatic resume, deletion of previous fixtures, or replay of uncertain
mutations. Generated run IDs are short; generated network/site UI names are
stable per run and fit the form's lowercase 3..40-character constraint.

| File | Purpose |
|---|---|
| `report.json`, `report.txt` | One scenario result; ordered events and browser checks with phase, requirement, label, expected/actual values and artifact links |
| `world.json` | Planned world and confirmed identity bindings, refreshed after cleanup |
| `webapp-commands.jsonl` | Request persisted before send, then correlated response; no automatic action replay |
| `webapp-resources.json` | Owned resources, observations, cleanup state, worker/runtime shutdown state and run result |
| `webapp-worker.log` | Worker stderr kept separate from protocol stdout |
| `browser/<run_id>/` | Command results, atomic worker summary, and failure screenshot/trace/diagnostics when available |

Requests carry protocol version, run ID, increasing command ID, action, resolved
inputs and an absolute deadline. C enforces correlation, bounded message size,
structured replies and matching expected/actual values before accepting PASS.
Text checks normalize the same whitespace as the worker; values and units remain
exact. Browser resource names and IDs come from resolved world references.
Selectors must resolve to the requested network and detail entity; a detail URL
alone does not establish ownership or replace a visible assertion.

## Ownership and cleanup

The resource journal writes through a temporary file, fsyncs it, atomically
renames it and fsyncs the containing directory. The UI provisioning handler
explicitly records `created=1` with a matching planned name and UI-result ID.
Navigation records observations only. Conflicting IDs or unknown references fail.
Only resources marked owned can reach the deletion hook; name prefixes and
observed IDs never imply ownership. Cleanup proceeds nodes → sites → networks,
reverse creation order within each kind, records failures and continues while
budget remains. A failed cleanup retains the exact resource IDs for diagnosis.

The browser close handshake has a separate budget. C-originated failure is sent
with `failed: true`, preserving failure evidence. If the channel is lost, C first
signals the worker to capture evidence, then forcibly terminates remaining
process-group members and reaps the worker. Missing shutdown acknowledgement is
reported as unconfirmed cleanup, never a successful run. Screenshots/traces are
best effort after a browser crash or forced termination.

Runtime reconnection has a separate 30-second budget. Owned-resource cleanup
has a fresh 120-second total budget; final runtime cleanup
has a separate 15-second budget. SIGINT/SIGTERM stop acceptance work, but do not
cancel cleanup. Runtime/cleanup hooks execute in bounded child process groups;
callbacks must not depend on child memory changes reaching the parent. The
legacy `ULAB_NO_CLEANUP` switch does not disable this browser lifecycle.

An assertion/action receives the smaller of its step deadline and the remaining
scenario budget. C allows up to ten further seconds to receive failure evidence,
but never accepts a late successful assertion. Cleanup can extend elapsed time
beyond the acceptance deadline. SIGKILL/host loss cannot execute cleanup; this
patch retains journals but provides no automatic crash recovery command.

## Verification

The tests require a C compiler, Jansson/libyaml development headers/libraries,
Python 3, Node.js, the worker dependencies and a locally installed Chromium.
Build the complete lab CLI first. For a binary outside `bin/ukama-lab`, set
`ULAB_TEST_BINARY` to its absolute path. Nonstandard dependency locations can be
supplied with `ULAB_TEST_CFLAGS` and `ULAB_TEST_LIBS` for the journal probe.

```sh
python3 -m unittest discover -s tests/webapp -v
npm --prefix adapters/webapp test
node --test --test-concurrency=1 tests/webapp/runner-browser.mjs tests/webapp/provisioning-browser.mjs
```

The C tests cover the real parser and runner, workload/webapp dispatch, command
order, correlation failures, invalid JSON, crashes, false PASS replies, timeouts,
cancellation, private/exclusive artifacts, SKIP semantics, entity binding,
observed-vs-owned resources, partial cleanup failure, and a hung deletion hook.
The end-to-end harness runs the real C CLI and worker against a source-derived
DOM fixture with real Chromium, including failure and active cancellation.
These are infrastructure tests. They do not claim a real console/BFF journey,
Firefox/WebKit qualification, or 90–100% product coverage.

## Patch 7 expansion

See [expanded.md](expanded.md) for semantic UI interactions and 14 added scenarios,
and [matrix.md](matrix.md) for serial browser execution and conservative evidence
reporting. Existing direct-run commands remain supported.
