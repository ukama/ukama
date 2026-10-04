# Scenario language v2: web-app foundation

Version 2 defines browser acceptance scenarios for the actual console-app.
Existing `scenarios/p0/console` and `console-kpi` are BFF tests and receive no
browser coverage credit. Version 1 behavior and files remain supported.

## Patch status

Patches 1–3 implement the contract, local Playwright worker and C runner.
`run`/`validate` now execute **active, zero-fixture authenticated scenarios**.
See [`webapp/runner.md`](webapp/runner.md) for the executable Members example,
setup, lifecycle, journals, cleanup and testing commands.

Network/site provisioning and node runtime setup remain gated until Patch 4
supplies their UI handlers. The two network examples remain WIP. `wip` and
`skip` produce `outcome: SKIP`, `passed: false`, and zero checks without launching
the worker or backend setup. V2 `xfail` execution is rejected so infrastructure
failures cannot be converted to a passing scenario.

Version 2 is also used by the existing workload language. Workloads must declare
`kind: workload`; version alone does not select the workload runner.

The initial v2 vocabulary covers network/site fixtures, navigation, visible
text/count/action checks, and node disconnection/reconnection. Additional UI
operations are introduced with their adapter handlers in later patches.
Unknown operations are rejected rather than treated as no-ops.

## Offline validation

With a full lab binary:

```sh
./bin/ukama-lab lint scenarios/webapp
./bin/ukama-lab lint scenarios/p0 scenarios/smoke
```

Without the external ukamaOS build tree or a lab binary:

```sh
./utils/lint-scenarios.sh scenarios/webapp
python3 -m unittest discover -s tests/webapp -p test_contract.py -v
```

The helper needs a C compiler and uses the same C parser/validators as the lab.
It compiles only into a temporary directory and creates no backend/runtime
resources. The existing fixed-size `scenario_t` is about 1.1 GiB; the lint
helper inherits that memory requirement. It reuses one scenario allocation and
runs files sequentially. This patch does not refactor the existing structure.

`lint` recursively checks `.yaml`/`.yml` files in deterministic order. It prints
`LINT OK` for a valid contract and a summary with no execution/coverage credit.
It exits nonzero on an error or if no scenario files were found. Workload YAML,
model YAML, and generated catalog indexes are not scenario documents; pass the
scenario files/directories you intend to validate.

Environment substitution remains explicit: `${NAME}` must exist. Offline
validation does not invent software-version values. The regression tests use
fixed placeholder version strings only to validate legacy scenario contracts;
they do not claim a published software release or run an update.

## Top-level contract

Start a scenario with `version: 2`, `name`, `suite: webapp`, `priority: p0|p1|p2`,
and `status: active|wip|skip|xfail`. `seed`, `description`, and `tags` retain their
version 1 meanings. The provider is `virtual` for this foundation.

```yaml
webapp:
  base_url: http://localhost:3000
  browser: chromium
  headless: false
  auth_state: .auth/owner.json
  action_timeout_seconds: 30
  check_timeout_seconds: 30
  scenario_timeout_seconds: 3600
```

- `base_url` and `auth_state` are required. The latter is a path to real saved
  Playwright authentication state; the executor will resolve it relative to
  the lab working directory. Lint checks the contract, not file existence,
  credentials, endpoint health, or a browser installation.
- Supported browser identifiers are `chromium`, `firefox`, and `webkit`.
  Patch 2 initially qualifies Chromium; a valid identifier is not a claim of
  an installed/tested browser.
- `headless` defaults to `true`; booleans accept only `true` and `false`.
- Action/check defaults are 30 seconds. Per-step `timeout_seconds` overrides
  the corresponding default. Step deadlines are 1..900 seconds and cannot
  exceed the whole-scenario deadline of 1..86400 seconds (default 3600).
- Browser worker, browser, and C runner execute on the host. The configuration
  does not select an AWS/distributed runner. Local execution is separate from
  the chosen app/auth/BFF endpoint locations.
- Authentication state and secrets are local configuration and must not be
  included in scenarios, reports, or patch bundles.
- The standalone Patch 2 worker accepts a `base_url` origin without a path
  prefix; the supplied console's routes are rooted at `/`. It rejects a
  configured subpath explicitly even though the foundation linter accepts it.

## World, setup, and resource ownership

```yaml
world:
  networks: 1
  sites_per_network: 2
  nodes_per_site:
    tower: 1
    amplifier: 1
    controller: 1
runtime:
  start: [nodes]
  wait: [nodes_ready]
setup:
  create_via_webapp: [networks, sites]
```

`world` is the desired fixture inventory. It does not prove that any UI action
succeeded. `create_via_webapp` instructs the future executor to perform real UI
workflows, log their substeps, verify visible completion, and bind actual
identities back to `net-NNN`, `site-NNN`, and node references.

The local runtime prepares/registers virtual node bundles and checks physical
readiness. It must not pre-create the networks/sites whose UI creation is under
test. The site fixture has exactly one tower, amplifier, and controller. Both
runtime start and readiness wait are required. Network/site fixture counts
must agree with `create_via_webapp`; BFF provisioning flags are forbidden in v2.
The foundation bounds fixture dimensions at 100 networks and 100 sites per
network; this is a validation bound, not a recommended host workload.

Zero-network scenarios are allowed for organization/account screens and have
no network/site setup. Customer, plan, SIM, UE, and traffic fixture contracts
are added with patch 6. They currently fail with an explicit unsupported-scope
message, so they cannot silently use a BFF shortcut.

Patch 3 generates short, deterministic network/site UI names independently of
internal references. The ownership journal is ready for UI provisioning handlers
to record each confirmed created resource in Patch 4, including partial creation.

## Events

| Type | Required fields | Optional fields | Meaning |
|---|---|---|---|
| `web_open` | `view`; one `networks: net-NNN` for network-dependent views | `timeout_seconds`; matching `sites` or `nodes` detail selector | Navigate using the UI to the selected semantic view/entity. |
| `web_select_network` | One `networks: net-NNN` | `timeout_seconds` | Use the network switcher and confirm selection. |
| `web_reload` | None | `timeout_seconds` | Explicit browser reload; never inserted by a check. |
| `disconnect_nodes` | Existing node reference or `nodes: all` | None | Use the existing host fault controls. |
| `reconnect_nodes` | Existing node reference or `nodes: all` | None | Restore the selected node connections. |

Node references use the existing shape, e.g. `tower-site-002-001`. Detail
navigation must name a node/site belonging to the selected network. Invalid
references, `networks: all`, unsupported fields, and masked failures are errors.
Expected UI rejection will use an explicit visible-error assertion when its
handler is added. A worker crash, timeout, or missing locator is never an
expected product rejection.

Accepted semantic views:

- `business_home`, `business_revenue`, `business_customers`, `business_packages`,
  `business_data_plans`, `business_members`, `business_sim_pool`,
  `business_support`, `business_settings`.
- `network_home`, `network_sites`, `network_site_detail`, `network_nodes`,
  `network_node_detail`, `network_customers`, `network_node_pool`,
  `network_sim_pool`, `network_support`, `network_settings`.
- `customer_customers`, `customer_data_plans`, `customer_settings`.
- `welcome`, `unauthorized`.

Settings, members, inventory-pool, welcome, and unauthorized views may omit a
network selector. Other views require one. The setup executor owns the
multi-step onboarding routes; they are not generic navigation aliases here.

## Checks

Every web check requires `view`, a visible `label`, and a `requirement: WEB-*`
identifier. Labels are located within the expected view, using accessible
roles/labels or an explicit test ID contract. The adapter verifies the current
view; a check does not navigate or reload to repair state.

| Type | Additional required field | Observation |
|---|---|---|
| `web_kpi_equals` | `expected` | Exact rendered KPI value, including units/formatting. |
| `web_field_equals` | `expected` | Exact visible field value. |
| `web_table_count_equals` | `expected_count` | Matching visible data rows, excluding header/loading rows. |
| `web_action_available` | `available: true|false` | Visible action exists and is enabled/disabled as expected. |

`timeout_seconds` is optional. All web checks are immediate, required, and
non-skippable by default; `immediate`, `poll_seconds`, `required`, and BFF check
fields are intentionally not accepted in this namespace. A count of zero,
`available: false`, and an explicitly empty expected string are valid expected
states. An unavailable action must be visibly disabled; a missing action is a
different requirement and must not be silently accepted as disabled.

Table row count means the current visible table page. Total result count and
all-page reconciliation require their own assertions when pagination handlers
are introduced. Do not interpret three visible rows as proof of total inventory.

A phase runs its events, then its checks, before the next phase. Event-only
phases are allowed. `final_checks` run before cleanup. BFF assertions cannot be
used as v2 acceptance checks. Browser assertions compare scenario expectations
to rendered content; GraphQL responses can support diagnostics/identity binding,
but cannot replace a visible assertion or supply both actual and expected.

The supplied Nodes and Sites screens render card lists, not tables. The node
example checks the visible `Nodes count` field with `web_field_equals`.
The Patch 2 adapter rejects table checks on those card views. Its supported
component locators and current per-screen limitations are listed in `worker.md`.

Requirement identifiers are syntax-checked by C. Catalog membership and
example-to-requirement links are checked by the Python contract tests. A
contract example does not count as implemented browser automation.

## Refresh and unit semantics

Network Home explicitly polls every 30 seconds in the supplied console. Node
and site list/detail queries commonly use the shared helper's disabled default;
operation locks have separate busy/focus polling. Scenarios must state whether
they expect automatic convergence or use an explicit `web_reload` event.
A failing automatic-update assertion must not be repaired by a hidden reload.

Longer analytics deadlines are observation budgets, not a license to repeat
creation, payment, restart, or other mutation actions. The scenario-wide deadline
includes setup, navigation, convergence, and final checks. Cleanup has its
own bounded budget and runs after failure/cancellation.

The UI's day-based validity and the backend's minute-based validity represent
one contract: 1/7/30 days correspond to 1440/10080/43200 minutes. Browser tests
enter days and verify the equivalent day/week/month presentation through edit,
allocation, and reload. This is a conversion requirement, not a missing feature.

## Coverage and rollout

`docs/webapp/coverage.json` contains the initial source-derived inventory and
all supplied page-route files, including redirects and the disabled Billing
route. It records 132 requirements, all `planned`/`not_run`. It is a starting
inventory for scope review, not an assertion of exhaustive coverage or a passed
release. Role permissions require an agreed behavior matrix.

The first two YAML files are `wip` contract examples. Later patches must add
handlers, necessary assertions, real run evidence, and activate them. Mapping
an ID to a YAML file alone earns neither automation nor verified credit.

Implementation sequence and commit-message subjects are recorded in
`docs/webapp/patch-plan.md`.
