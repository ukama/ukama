# Local console-app worker — Patches 2–4

Patch 2 adds a standalone Node/TypeScript Playwright worker inside ukama-lab.
The worker and browser run on your host. It drives the console's visible
navigation and reads rendered values. It never calls a BFF endpoint directly,
injects application state, or reloads a page to make an assertion pass.

Patch 3 connects this worker to the C scenario lifecycle; see [runner.md](runner.md).
Active scenarios with no provisioned world can now execute. The network/site
examples now execute through Patch 4 UI setup (see provisioning.md). No live requirement in
the coverage inventory has been marked verified from fixture runs.

## Install on the host

Use Node.js 22 or newer and npm. From the **lab directory**:

```sh
npm --prefix adapters/webapp ci
npm --prefix adapters/webapp run build
npm --prefix adapters/webapp run install:browser
```

The package lock pins Playwright and the TypeScript build dependencies. On a
Linux host missing browser libraries, install them with the package's local
Playwright CLI (this command may need administrator privileges):

```sh
./adapters/webapp/node_modules/.bin/playwright install --with-deps chromium
```

The host must reach the configured console, auth service and the endpoints used
by the console. Running the worker locally does not relocate those services.
Use your local console URL in the profile. Its existing runtime configuration
must point to your intended auth/BFF services, with working cookie domains and
cross-origin configuration. The worker does not start or reconfigure the app.

## Capture a real login once

From an interactive terminal with a graphical display:

```sh
./utils/webapp-worker.sh auth \
  --base-url http://localhost:3000 \
  --out .auth/owner.json
```

A browser opens. Sign in through the actual auth flow and reach the dashboard,
then press Enter in the terminal. The helper verifies the dashboard at the
requested origin before writing the authentication state. It does not collect
your password or fabricate a session. Complete any first-login welcome/setup
flow manually for this foundation test account.

State is written atomically with owner-only file permissions. `.auth/` is
ignored by git. Regenerate the state when it expires or the app/auth hostname
changes. The saved state is input only: scenario workers never overwrite it.
Each worker creates a fresh browser context from it. Separate contexts isolate
browser storage, not shared backend resources; start with one worker/account.

## Run a visible smoke check

Copy `adapters/webapp/profile.example.json` to
`adapters/webapp/profile.local.json` and adjust `base_url`/`auth_state`.
All relative paths resolve from the **lab working directory**, including the
state path inside the profile. `utils/webapp-worker.sh` consistently sets that
working directory even if invoked elsewhere.

Choose an existing network and an expectation you know independently:

```sh
./utils/webapp-worker.sh smoke \
  --profile adapters/webapp/profile.local.json \
  --network lab-network \
  --view network_home \
  --label 'Sites online' \
  --expected '2/2'
```

This example assumes a network named `lab-network` with two fully online sites;
it creates no fixtures. Change the name and expectation to your test inventory.
The check uses the visible network switcher and Network/Home navigation. It
compares the rendered KPI to `2/2`, without deriving the expectation from the
same backend response as the page.

For the visible node count:

```sh
./utils/webapp-worker.sh smoke \
  --profile adapters/webapp/profile.local.json \
  --network lab-network \
  --view network_nodes \
  --kind field --label 'Nodes count' --expected '3'
```

`--kind` accepts `kpi`, `field`, `action` and `table`. Action expectations must
be literal `true`/`false`; table expectations must be a nonnegative integer.
The smoke CLI covers a single check on a non-detail view; detail commands are
available through the worker protocol. Every smoke run has a unique run ID,
JSON result lines on stdout, concise progress on stderr, and exit code 0/1.
An unsuccessful check cannot become PASS when browser cleanup succeeds.

## Implemented adapter scope

Navigation supports the 22 dashboard views listed in `src/console-app.ts`:
Business, Network and Customer lenses, their sidebar pages, and node/site
detail pages opened by clicking the corresponding visible card. Navigation
uses the lens controls, network switcher and sidebar; only initial session
loading uses `goto`. Only `web_reload` performs an explicit page reload.

The existing source has no stable test IDs on the shared components, so the
adapter uses accessible roles/text and the following existing component classes.
It does not use generated Emotion classes or mutate the DOM. A duplicate
visible match is an error, not permission to select the first element.

| Operation | Supplied console component / selector contract |
|---|---|
| Dashboard ready | `DashboardShell`, `header.topbar .viewseg`, `aside.sidebar`, `main.main` |
| Lens | `LensSegment`, visible Business/Network/Customer buttons |
| Network | `NetSwitch`, `button.netswitch .nm`, menu item containing exact network name |
| Page navigation | `Sidebar`, exact link name in the visible sidebar |
| Detail navigation | `NodesScreen`/`SitesScreen`, `.ecard[role="button"]` containing resolved entity text; resulting URL must match resolved ID |
| KPI value | `Kpi`, one `.MuiCard-root` with the exact label; its second direct div contains value/unit |
| Detail field | `KV`, one visible `.kv-row` with exact label and value span `.tnum` |
| Nodes/Sites count | `PageHeader`, `.pagetitle .cnt`; explicit `No nodes yet`/`No sites yet` proves zero |
| Action availability | Exact button name within `main.main`, visible and enabled/disabled as expected |
| Table count | One actual visible table in the named list view, visible `tbody > tr` rows, no skeletons |

Nodes and Sites are card lists. The Patch 1 node example has been corrected to
check the `Nodes count` field. That count alone does not satisfy the complete
requirement to reconcile every card and detail; Patch 4 now adds those assertions.

Fields rendered outside the shared KV structure need explicit handlers in
later patches. The current table handler refuses card lists and ambiguous
tables. A missing table, error, skeleton, or replacement empty-state panel is
not silently accepted as zero. Actual zero-row tables are accepted. Per-screen
empty-table panels and pagination need later handlers.

The worker checks the expected view URL and visible network selection on every
observation. Checks never navigate, switch networks, click Retry, call an API,
or resubmit an action. A missing action cannot satisfy `available: false`.
Whitespace is normalized like rendered text; units/case/value remain exact.
An explicit empty, rendered KV value can satisfy `expected: ""`.

`welcome` and `unauthorized` remain recognized scenario-language names but have
no Patch 2 navigation handler: the worker rejects them explicitly. Login,
welcome, role/unauthorized acceptance cases belong to later patches. Creating
networks/sites, restart/update actions, commerce flows and cleanup of backend
resources are also not implemented here.

## JSONL protocol, version 1

Start `./utils/webapp-worker.sh worker` with dedicated stdin/stdout pipes. Send
one UTF-8 JSON object per line and wait for its response before sending the
next command. Logs go to stderr. Maximum command size is 1 MiB.

Each request contains exactly:

```json
{
  "protocol": 1,
  "run_id": "local-run-001",
  "command_id": 1,
  "action": "init",
  "deadline_ms": 1791137400000,
  "inputs": {
    "profile": {
      "base_url": "http://localhost:3000",
      "auth_state": ".auth/owner.json"
    },
    "artifacts_dir": "runs/webapp-worker"
  }
}
```

The timestamp above only illustrates the shape. The caller must calculate a
fresh absolute Unix-millisecond deadline, e.g. `Date.now() + 30000`, for each
new command. It must be in the next 1..900000 ms. After receipt the worker uses
a monotonic budget. The profile's scenario deadline also bounds the whole run,
including idle time. Action/check timeout profile values are defaults for the
caller; its per-command deadline supplies any scenario step override.

| Action | `inputs` |
|---|---|
| `init` | `profile` object, `artifacts_dir`; first command, once per worker |
| `web_open` | `view`, resolved `network_name` when scoped; detail views also require `entity: {ref, id, text}` |
| `web_select_network` | Resolved `network_name`; follow with `web_open` before checking a view |
| `web_reload` | Empty object |
| `web_kpi_equals`, `web_field_equals` | `view`, `label`, `requirement: WEB-*`, `expected` string |
| `web_table_count_equals` | `view`, `label`, `requirement`, `expected_count` integer |
| `web_action_available` | `view`, `label`, `requirement`, `available` boolean |
| `close` | `{}` or optional `failed: boolean`, `reason: UPPERCASE_CODE`; required even after failure |

The runner resolves YAML references to visible names/identities. For example,
`nodes: tower-site-001-001` becomes an entity object with that `ref`, the actual
node `id` and a distinctive visible serial/name in `text`. A detail click
returns a URL-observed binding only after the requested ID is in the URL.
This is identity evidence, not a substitute for visible field assertions.

Every response contains matching `protocol`, `run_id`, `command_id`, `action`,
`status: ok|error`, `run_status: running|passed|failed`, `duration_ms`,
`expected`, `actual`, `bindings`, `artifacts`, and an `error` code/message when
applicable. Invalid envelopes which cannot be parsed have null correlation
fields. Assertion timeouts retain the most recent observed value when one was
available. Missing elements produce null, not an invented zero or false.

Run IDs accept letters, digits, hyphen and underscore (1..80 chars). New command
IDs increase from 1 up to 10000. Repeating the exact same command returns its
cached response and does not execute it again; keep the original deadline on
a replay. Reusing an ID with different content is a terminal error. The cache
is per process, not durable recovery: after a crash, never blindly resubmit a
mutation to a new worker. The C runner journals requests before sending and never replays an uncertain action.

A failure stops browser work, captures evidence and closes the browser. It
continues accepting `close` to finish the protocol, and that response preserves
`run_status: failed`. Closing successfully means worker commands completed;
the C runner owns the final scenario result and coverage accounting.
`close` with `failed: true` preserves a C/runtime-originated failure and retains
browser evidence even if every browser command passed. `reason` accepts 1..64
uppercase letters/underscores; keep raw error details in the C report.
EOF without explicit close, SIGINT/SIGTERM, broken output and scenario timeout
are failures. Evidence has a separate five-second capture budget, followed by
two seconds each for context/browser close. The runner uses a bounded worker shutdown followed by process-group termination
for remaining helpers.

## Failure evidence

The artifact directory is `runs/webapp-worker/<run_id>/` by default. Existing
run directories are rejected instead of overwritten.

- `results.jsonl`: ordered command results, expected/actual values and errors.
- `worker-summary.json`: terminal state/reason, including EOF or cancellation.
- `failure.png`: screenshot when a page remains available.
- `trace.zip`: Playwright trace on failure; successful runs discard it.
- `diagnostics.json`: page URL without query/fragment/credentials, console
  event levels, request-failure metadata and HTTP error statuses. No headers,
  response bodies, cookie values or raw console arguments.

Directories are created owner-only; files are owner read/write. Screenshots and
traces can contain account/customer data; raw traces can contain authentication
material. Keep them local/private, outside source control and shareable reports.
They are not redacted just because the diagnostics JSON is sanitized. State
files are never copied into reports. An abrupt browser crash can make a
screenshot or trace unavailable; the result still remains a failure.

To inspect a private trace locally:

```sh
./adapters/webapp/node_modules/.bin/playwright show-trace runs/webapp-worker/RUN_ID/trace.zip
```

## Verification and limits

```sh
npm --prefix adapters/webapp test
python3 -m unittest discover -s tests/webapp -v
./utils/lint-scenarios.sh scenarios/webapp
```

Worker tests start a local HTTP fixture with DOM structures derived from the
supplied console. A real Chromium browser exercises the protocol, navigation,
observations, auth-state loading, traces, failure/exit behavior, replay and
deadlines. The fixture is not the console or BFF; those passes do not qualify
product behavior or the 90–100% coverage target.

The console/backend stack and a real login were not available in the delivery
environment. Run the smoke command above on the host to validate the adapter
against that stack before integrating it. The delivery verification record
states the browser/version actually used. Firefox/WebKit identifiers are
accepted and can be installed explicitly, but were not qualified by this patch.

For a host with an existing compatible browser, the optional environment
variable `ULAB_WEBAPP_EXECUTABLE_PATH` selects an absolute executable path.
Normally leave it unset so Playwright uses its matching installed browser.
Always record overrides with test evidence; they are not equivalent to testing
every Playwright browser build. Viewport is fixed to 1440×1000, locale to en-US,
and timezone to UTC for this foundation.

Playwright API references used for this implementation:
[Browser contexts](https://playwright.dev/docs/api/class-browsercontext),
[tracing](https://playwright.dev/docs/api/class-tracing),
[library usage](https://playwright.dev/docs/library).


## Patch 4 extension

The C setup lifecycle sends `web_create_network` and `web_create_site`; these
are internal setup commands, not arbitrary YAML events. Network inputs contain
`ref` and planned `name`. Site inputs additionally contain `network_name`,
`network_id`, `tower_id` and optional exact component labels in `components`.
The worker only clicks/fills real UI controls. It passively observes matching
mutation responses (or the site's UI-triggered scoped polling) to write durable
creation receipts. It does not issue direct create, GraphQL or acceptance reads.
See [provisioning.md](provisioning.md) for receipt recovery and uncertain outcomes.

`web_field_equals` accepts an optional resolved `node_id` on `network_nodes` to
scope Serial #, Model type, Site and Connectivity to one card. Detail navigation
matches a node card's serial line, then verifies the exact requested URL.
C resolves YAML `expected_ref`/`expected_property` before dispatch; the worker
never computes the expected value from the displayed value. The connectivity
dot's native tooltip is read only when its element is visible. Lifecycle states
remain separate and are not inferred from connectivity.


## Patch 5 extension

`web_action` accepts the resolved detail `view`, `network_name`, `entity`,
semantic `action` and action-specific `value` or `app`/`tag`. `web_tab` accepts
`tab: primary|secondary`; each tab has an independent navigation scope. Operation
acknowledgements contain `actual.executed: true` and no resource bindings.
`web_field_equals` and `web_action_available` accept `app` for software scope.
A reason check can use `match: contains`; every other value remains exact. C
verifies the acknowledgement and expected/actual match independently.

See [operations.md](operations.md) for the complete supported vocabulary and
seven scenarios. No direct operation API requests or application-store reads
are used by these handlers. Two tabs share authenticated browser context state;
this is not a different-user test.

## Patch 6 extension

The worker accepts `web_commerce`, `web_commerce_equals` and the runner-owned
`web_import_sims` setup command. See `commerce.md` for typed fields, UI locators,
world references, units and limits. Commerce submissions have durable private
receipts and per-run mutation guards; duplicate command replay returns the
cached response. Observation remains separate from submission.

Patches 12/13 extend this contract with `rename_plan`, `allocate_auto`, scoped
commerce read faults, `failed_top_up`, `top_up_rapid` and `download_receipt`.
See `lifecycle.md` and `payments.md`. Browser contexts now accept downloads so the
requested receipt can be copied to private artifacts; temporary browser downloads
are removed when the context closes. PDF validation requires Poppler on PATH.
An intercepted payment rejection permits one explicit new submission, while
unknown or ambiguous outcomes retain the existing duplicate guard and journal.

Patch 6 also removes the operations adapter's dependency on console companion
changes. Current locators target the unchanged console source. A stale or
optimistic app value is not silently corrected by the adapter.

## Staged onboarding (Patch 9)

`session_mode: onboarding` initializes an isolated saved-auth context without
opening a dashboard. `web_onboard` and `web_onboard_equals` drive the configure
wizard with staged ownership receipts. Manual auth capture additionally accepts
`--landing configure`. See `onboarding.md` for the closed command contract, local
prerequisites, controlled-evidence limits and original console gaps.

## Operation observation and distinct sessions (Patch 11)

See [operation-recovery.md](operation-recovery.md) for WB110–117. `web_action`
adds port controls, read-only restart/timeout watches and explicitly classified
status-read faults. Only `watch_restart` accepts C-resolved `nodes`; YAML cannot
supply this identity array. First `web_tab` with `tab: peer` requires a separate
saved `auth_state`; its context is isolated and both contexts retain failure
traces. Existing primary/secondary tabs still share their original context.
