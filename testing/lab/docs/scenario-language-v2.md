# Scenario language v2: web-app foundation

Version 2 defines browser acceptance scenarios for the actual console-app.
Existing `scenarios/p0/console` and `console-kpi` are BFF tests and receive no
browser coverage credit. Version 1 behavior and files remain supported.

## Patch status

Patches 1–5 implement the contract, local Playwright worker, C runner and UI
network/site provisioning with local virtual-node runtime. The two operational
network examples are active; see [`webapp/provisioning.md`](webapp/provisioning.md).
Patch 5 adds browser operations, lock/recovery assertions and software updates;
see [`webapp/operations.md`](webapp/operations.md) for the action vocabulary,
scenario authoring and host preconditions.
`wip` and `skip` produce `outcome: SKIP`, `passed: false`, and zero checks without
launching the worker or backend setup. V2 `xfail` remains unsupported.

Version 2 is also used by the existing workload language. Workloads must declare
`kind: workload`; version alone does not select the workload runner.

The initial v2 vocabulary covers network/site fixtures, navigation, visible
text/count/action checks, and node disconnection/reconnection. Additional UI
operations are introduced with their adapter handlers; Patch 5 provides
`web_action` and `web_tab`.
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
succeeded. `create_via_webapp` instructs the executor to perform real UI
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
internal references. The ownership journal records UI creation receipts and factory runtime claims,
including partial setup. Unresolved submitted mutations fail cleanup and retain
resources for manual reconciliation; they are never automatically resubmitted.

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
| `web_field_equals` | `expected`, or `expected_ref` + `expected_property` | Exact visible field value. |
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

Patch 4 supports `expected_ref` for an existing node and `expected_property`
`id`, `model`, or `site_name`. Use both fields instead of literal `expected`.
C resolves these from the planned world and confirmed runtime identity, never
from the browser's displayed value. On `network_nodes`, add `nodes: <ref>` to
scope a field assertion to one visible card. Supported card labels are
`Serial #`, `Model type`, `Site`, and `Connectivity`. On node detail, omit the
card selector after `web_open` has selected the exact entity. Connectivity is
read from the visible dot's native tooltip (`title`), independently of lifecycle.

```yaml
- type: web_field_equals
  view: network_nodes
  nodes: tower-site-001-001
  label: "Serial #"
  expected_ref: tower-site-001-001
  expected_property: id
  requirement: WEB-NODE-002
```

The `webapp` block additionally accepts optional `switch_component`,
`backhaul_component`, and `power_component` visible option labels. Omit them
when the wizard auto-selects a single available component. Creation commands
and each runtime setup script receive at most 900 seconds, capped by the
remaining scenario budget; ordinary events/checks keep their declared budgets.

Requirement identifiers are syntax-checked by C. Catalog membership and
example-to-requirement links are checked by the Python contract tests. A
contract example does not count as implemented browser automation.

## Refresh and unit semantics

Network Home explicitly polls every 30 seconds in the supplied console. Node
detail now opts into visible polling (30 seconds, or 3 seconds during operations
and on Software). Other lists/site detail can still use the disabled default;
operation locks have separate busy/unverified/focus polling. Scenarios must state whether
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
route. It records 132 requirements; after Patch 5, 16 have complete automation and
all retain `verification: not_run`. Partial cases remain planned with explicit
gap notes. This is a scope inventory, not a passed release. Role permissions
require an agreed behavior matrix.

The network and operation YAML examples are active and have local fixture
qualification. Live target-app runs must still provide build/profile evidence. Mapping
an ID to a YAML file alone earns neither automation nor verified credit.

Implementation sequence and commit-message subjects are recorded in
`docs/webapp/patch-plan.md`.

## Patch 6 commerce extension (lab only)

`web_commerce` provides typed actions for plan/customer creation, exact SIM
allocation, SIM service toggles, cash top-up/cancellation and receipts.
`web_commerce_equals` asserts scoped visible fields, with independent literal
expectations or resolved `iccid`, `plan_name`, and `payment_id` identities.
World package durations use either days or canonical minutes; the UI permits
only 1/7/30 days (1440/10080/43200 minutes). Binary data units use 1024 MB/GB.

`start_ues` and `traffic` are permitted bounded runtime events after UI
allocation, with `ues` and optional `timeout_seconds` (1..900). Traffic requires
positive `amount_mb`; neither event is an acceptance check. Runtime actions
cannot mask failures. Browser checks observe their results separately.

See [webapp/commerce.md](webapp/commerce.md) for field tables, seven runnable
examples, local services, retained inventory/ledger records, and coverage gaps.
No console companion is required or included.

Patch 6 also supports plan-scoped `web_commerce_equals` on `business_packages`
for `Performance price`, `Performance sold`, `Performance revenue` and
`Performance share`. Use the world `package` reference and a literal `expected`.
Business headline totals use the existing `web_kpi_equals` checks. See wb-022
and wb-027 for independent purchase expectations and explicit reload checks.

## Patch 10 inventory and navigation

`web_inventory` and `web_inventory_equals` add closed navigation commands,
owned-world list memberships, detail identity, map selection and transient
scope observation. The C runner resolves identity expectations as strings or
complete JSON arrays and independently verifies returned values. See
[webapp/inventory.md](webapp/inventory.md) for action/label mappings, local
expectation inputs and the eighteen runnable scenarios. Controlled stale
preference and masked read-response scenarios never receive live coverage
credit. Missing list controls and geographic map placement remain partial.
