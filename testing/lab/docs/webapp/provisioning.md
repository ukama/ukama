# Local UI provisioning and operational journeys — Patch 4

Apply the lab archive after Patch 3. Apply the console companion archive to the
uploaded console source and rebuild/restart it. Both archives contain changed
files directly at their repository root plus a unified diff and commit message.
The companion fix carries the chosen tower ID from Node pool through the
select-network, network-creation fallback and install steps. Selecting an
amplifier/controller anchors its associated tower. The worker checks that ID
and network again on the visible site-settings screen before submission.

## Host prerequisites

Use an authenticated owner who already reaches the dashboard (initial account
onboarding is not part of Patch 4). The console, auth, BFF, registry and supporting
services must refer to the same lab organization. The test runner, Playwright
browser and virtual-node Podman runtime execute on the host. No AWS runner is
created. Endpoints are configurable; existing factory/auth defaults can still
point at development services, so set them for an entirely local stack.

Install/build the worker as described in `worker.md`. Use the existing virtual
node prerequisites: Podman, curl, jq, an Ukama checkout and a configured Factory
with complete unprovisioned tower/amplifier/controller bundles. The factory
organization must match the browser/BFF organization. Every trio needs valid
coordinates, access inventory for its tower, spectrum inventory, and switch,
backhaul and power components. Registration/readiness occurs through the normal
runtime. The test never creates operator networks/sites through the BFF.

Capture `.auth/owner.json` with `utils/webapp-worker.sh auth`. Configure backend
teardown credentials independently using `UKAMA_SESSION_TOKEN` +
`UKAMA_BFF_TOKEN`, or `UKAMA_IDENTIFIER` + `UKAMA_PASSWORD` + `PAUTH_URL`.
Browser state is not copied into the cleanup client. The same organization and
teardown permissions are required; a failed deletion fails the run.

From the lab root, after configuring those values locally:

```sh
export FACTORY_ORG=your-lab-org
export PAUTH_URL=http://localhost:8081

./bin/ukama-lab lint scenarios/webapp
./bin/ukama-lab run scenarios/webapp/p0/network/wb-001-sites-online-recovery.yaml \
  --repo /path/to/ukama --scripts scripts --out runs \
  --bff http://localhost:8080/graphql --factory-url http://localhost:8082
./bin/ukama-lab run scenarios/webapp/p0/network/wb-002-node-list-detail.yaml \
  --repo /path/to/ukama --scripts scripts --out runs \
  --bff http://localhost:8080/graphql --factory-url http://localhost:8082
```

The port numbers are examples; use your configured endpoints. Keep secret values
out of committed YAML. `--factory-url` (default from `UKAMA_LAB_FACTORY_URL`) overrides `FACTORY_URL` for setup scripts.
Runtime repository and run-directory paths currently cannot contain spaces or
shell metacharacters because the inherited runtime API serializes arguments.
Scripts execute as argv, without shell evaluation of scenario data.

The site wizard uses its default when a component category has exactly one
choice. For multiple choices, use exact visible labels in the YAML profile:

```yaml
webapp:
  base_url: http://localhost:3000
  auth_state: .auth/owner.json
  switch_component: Your switch description
  backhaul_component: Your backhaul description
  power_component: Your power description
```

## Execution and evidence

1. Create networks using the switcher and Add network dialog; verify the visible
   selected name. Names are unique deterministic lowercase names within 40 chars.
2. Reserve/build/start one local virtual trio, map its factory IDs to world refs,
   and wait for runtime readiness. Select that tower's Configure action in Node
   pool and choose the new network using the radio control.
3. Confirm installation, name the site, verify its visible tower identity, choose
   components, and press Create site once. Finish setup and verify the new site
   card and its matching detail URL and heading. Repeat for the next site.
4. Execute the scenario's browser events and assertions. Recovery stays on
   Network Home while host controls disconnect/reconnect the second tower. Its
   three checks require 2/2, 1/2 and 2/2 without a reload. The node example checks
   all three cards and details against world IDs, models, site names and Online
   connectivity (25 field assertions). Search/filter/lifecycle/action coverage
   belongs to later patches and is not credited here.
5. Close the browser, recover receipts, reconnect runtime nodes, delete owned
   backend nodes/sites/networks, and remove runtime containers/network.

Browser setup commands and each setup script have a 900-second maximum capped
by the scenario deadline. Recovery checks allow up to 900 seconds each because
backend connectivity/analytics may converge slowly. Ordinary navigation uses
the profile action timeout. Cleanup budgets are fresh: reconnect 30s, resources
120s total, runtime teardown 15s. Large worlds may need future configurable
cleanup budgets; these examples use one/two sites.

`webapp-resources.json` records creation intents before dispatch. The worker
atomically fsyncs `browser/<run_id>/creation-<command_id>.json` before submission
and when passive UI response metadata reveals a matching ID. Site polling from
the UI can also establish that identity after an uncertain transport response.
These receipts establish ownership only; passing still requires visible UI
completion. C validates receipts against the intended command/ref/name and
recovers them even when the command failed or its response was lost.

A submitted creation without an identifiable receipt is uncertain. The runner
fails, retains owned backend resources, stops its local runtime and records
incomplete cleanup. Use the planned names, command journal, private trace and
backend operator tools to reconcile it manually. It never retries Create/Try
again or guesses an ID for deletion. SIGKILL/host loss still requires manual
recovery from saved evidence; automatic cross-run resume is not provided.

The runtime writes the exclusively reserved, assigned trio state before build
or start, allowing teardown after a partial runtime failure. Factory reservation
IDs are also recorded in `runtime-sites/factory-claims.tsv`. Existing factory
reservations are consumed, not automatically returned to the unprovisioned pool;
no release API is assumed. An interrupted reservation/organization assignment
may require factory-side reconciliation from that file. This is separate from
backend resource and Podman cleanup.

## Verification scope

Run `tests/webapp/provisioning-browser.mjs` with the compiled C binary and worker.
It uses Chromium with a source-derived DOM, runtime-script doubles and a cleanup
BFF fixture. It exercises UI mutations, all three node mappings, automatic KPI
updates, partial failures, durable ownership and strict deletion results. The
runtime double never supplies browser assertions directly.

Live console/auth/BFF/Factory/Podman journeys require the configured host stack
and remain **not run** in this delivery. Fixture passes do not count as verified
product coverage. The inventory marks five requirements as implemented and
retains live verification as `not_run`; no 90–100% claim is made. The existing
BFF `p0/console` suite remains separate. Day/minute translation and commerce
journeys remain scheduled for Patch 6.
