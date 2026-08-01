# Scenario language v1

Top-level fields: `version`, `name`, `seed`, optional `suite`, `priority`,
`tags`, `status`, optional `provider`, then `world`, `packages`, `setup`,
`runtime`, optional `profiles`, `phases`, and `final_checks`.

The language is strict. Unknown event/check names must fail validation.
Packages accept either `duration_days` or `duration_minutes`; exactly one is
required. Both forms remain first-class scenario input. The BFF currently
accepts duration in minutes, so the runner converts days to `days * 1440` at
the BFF boundary and passes minute values unchanged.
Packages default to `scope: network`. A network-scoped package may name an
explicit `network`; otherwise one package instance is created per scenario
network. `scope: organization` creates one organization package and must not
specify a network.

Provider block is optional. Missing provider defaults to `virtual`.


Scenario scalar values support explicit environment substitution using
`${NAME}`. A missing variable is a scenario-load error. This is useful for
exact build versions that change with the repository commit:

```yaml
tag: ${ULAB_SOFTWARE_TARGET_VERSION}
```

```yaml
provider:
  type: virtual
```

Only `virtual` is supported in this build.

Scenario status values:

- `active` runs normally
- `wip` is skipped by default
- `skip` is skipped by default
- `xfail` may fail without failing the command

Supported events:

- `traffic`
- `traffic_by_profile`
- `create_ues` (reserved/disabled in v1.0)
- `start_ues`
- `wait_ues_attached`
- `restart_nodes`
- `wait_node_connectivity`
- `wait_nodes_ready`
- `add_package_to_sim`
- `purchase_package`
- `purchase_packages_parallel`
- `allocate_sim`
- `create_invalid_package`
- `wait_package_boundary`
- `set_package_active`
- `remove_package_from_sim`
- `toggle_sim_status` (`set_sim_status` remains an alias)
- `wait_sim_status`
- `toggle_service`
- `toggle_radio`
- `toggle_internet_switch`
- `restart_site`
- `promote_release`
- `software_update`
- `disconnect_nodes`
- `reconnect_nodes`
- `failure_control`
- `check`

Supported checks:

- `backend_count` (`count` is accepted as an alias)
- `list_contains`
- `list_excludes`
- `status_equals`
- `traffic_allowed`
- `traffic_blocked`
- `node_ready`
- `ue_attached`
- `usage_per_sim`
- `usage_sample`
- `package_active`
- `package_remaining` (skipped until BFF exposes remaining balance)
- `package_state`
- `package_assignment_count`
- `package_assignment_chain`
- `package_fields_equal` (`package_catalog_equals` remains an alias)
- `package_visible`
- `package_hidden`
- `package_name_available`
- `package_business_metrics`
- `sim_unallocated`
- `payment_equals`
- `payment_count`
- `kpi_value`
- `kpi_trend`
- `kpi_contract`
- `kpi_rollup_consistency`
- `performance_report_cell`
- `performance_report_row`
- `revenue_summary`
- `subscriber_billing_summary`
- `payment_entitlement_reconciles`
- `package_dashboard_metric`
- `network_overview_metric`
- `console_inventory_reconciles`
- `usage_aggregate`
- `node_state`
- `dashboard_loads`
- `node_version_equals`
- `node_health_ok`
- `release_unavailable`
- `list_count_equals`
- `entity_fields_equal`
- `entity_reconciles`
- `node_status_equals`
- `software_status_equals`
- `software_count_equals`
- `node_operation_status_equals`
- `site_operation_status_equals`
- `kpi_state_equals`
- `kpi_timeseries`
- `balance_non_negative`


Event expected failure:

```yaml
- type: restart_nodes
  nodes: all
  expect:
    result: failure
    error_contains: "script failed"
```

Node restart transition:

```yaml
- type: wait_node_connectivity
  type_selector: tower
  count_per_network: 1
  connectivity: Online
  seconds: 120

- type: restart_nodes
  type_selector: tower
  count_per_network: 1

- type: wait_node_connectivity
  type_selector: tower
  count_per_network: 1
  connectivity: Offline
  seconds: 120

- type: wait_node_connectivity
  type_selector: tower
  count_per_network: 1
  connectivity: Online
  seconds: 300
```

`wait_node_connectivity` polls BFF `getNode.status.connectivity`. It succeeds only
after every selected node has been observed with the requested connectivity.
The default timeout is 180 seconds.
`ULAB_NODE_CONNECTIVITY_POLL_SEC` controls the polling interval and defaults to
two seconds.

Virtual node network outage:

```yaml
- type: disconnect_nodes
  type_selector: controller
  count_per_network: 1

- type: wait_node_connectivity
  type_selector: controller
  count_per_network: 1
  connectivity: Offline
  seconds: 180

- type: reconnect_nodes
  type_selector: controller
  count_per_network: 1
```

`disconnect_nodes` removes the selected running Podman container from the lab
network without deleting it. `reconnect_nodes` reconnects the same container,
preserving its writable state across the simulated backhaul outage.

A software update can be interrupted without a special combined event:

```yaml
- type: software_update
  type_selector: controller
  count_per_network: 1
  app: example
  tag: ${ULAB_SOFTWARE_TARGET_VERSION}

- type: wait
  seconds: 2

- type: disconnect_nodes
  type_selector: controller
  count_per_network: 1
```

The update mutation is asynchronous, so the following `wait` and
`disconnect_nodes` events provide deterministic orchestration while keeping
each action independently visible in the run log.

Backend count:

```yaml
- type: backend_count
  target: sims
  expected: from_world
```

List/status/runtime checks:

```yaml
- type: list_contains
  view: sims
  ref: ue-000001

- type: status_equals
  entity: sim
  ref: ue-000001
  status: active

- type: traffic_allowed
  ues: all
  amount_mb: 1
```


BFF lifecycle events:

```yaml
- type: remove_package_from_sim
  ues: all

- type: add_package_to_sim
  ues: all
  package: daily_1gb

- type: toggle_sim_status
  ues: all
  status: inactive

- type: wait_sim_status
  ues: all
  status: inactive
  seconds: 120
```

`toggle_sim_status` calls the same mutation as the console and treats
`success: false` as an event failure, preserving the BFF message for
`expect.error_contains`. The backend applies SIM status asynchronously;
`wait_sim_status` polls `getSim.status` until every selected SIM reaches the
requested state. `ULAB_SIM_STATUS_POLL_SEC` controls the polling interval and
defaults to two seconds. `set_sim_status` remains a compatible alias.

Console site actions:

```yaml
- type: restart_site
  sites: site-001-001

- type: toggle_service
  sites: site-001-001
  state: on

- type: toggle_radio
  sites: site-001-001
  state: off

- type: toggle_internet_switch
  sites: site-001-001
  port: 1
  state: on
```

`restart_site` uses the console BFF `restartSite` mutation with both the site
and network IDs. For compatibility with older generated scenarios, it also
accepts a `nodes:` selector and derives the unique sites from those nodes.
`toggle_internet_switch` requires a positive port number. Controller action
responses include the BFF failure message, so operation-lock scenarios can
assert the reason with `expect.error_contains`.

Subsequent cash package sale:

```yaml
- type: purchase_package
  ues: all
  package: next_plan
  amount: 5.00
  currency: USD
```

The first package is mandatory during SIM allocation and its payment is
created internally by the backend. `purchase_package` is only for subsequent
purchases on an already allocated SIM. It calls BFF `addPayment` with
`itemType: package`, `paymentMethod: cash`, the selected SIM UUID, and the
package UUID. An `idempotency_key` is parsed but rejected at execution time
until BFF exposes such a field; this makes the missing contract explicit.

Data-package edge-case events:

```yaml
- type: allocate_sim
  ues: all
  package: initial_plan
  expect:
    result: any

- type: purchase_packages_parallel
  ues: all
  package: plan_a
  other_package: plan_b

- type: create_invalid_package
  package: baseline_plan
  variant: negative_price
  expect:
    result: failure

- type: wait_package_boundary
  ues: all
  package: initial_plan
  offset_seconds: 1
```

`allocate_sim` exercises an explicit allocation attempt using a SIM selected
from the factory pool. `result: any` is useful for retry/idempotency scenarios:
the mutation may return the original allocation or reject the duplicate, while
the following GraphQL checks remain authoritative. `purchase_packages_parallel`
submits the two cash sales concurrently. `create_invalid_package` supports
the `allowance`, `duration`, `price`, and `currency` variants, which submit a
negative allowance, zero duration, negative price, and empty currency.
`wait_package_boundary` obtains the entitlement end date through
`getPackagesForSim` and waits relative to that server-provided boundary; a
negative offset means before expiration.

Entitlement and payment checks:

```yaml
- type: package_state
  ues: all
  package: next_plan
  expected: queued
  timeout_seconds: 30
  poll_seconds: 5

- type: payment_equals
  ues: all
  package: next_plan
  status: settled
  payment_method: cash
  currency: USD
  expected_value: 5.00
  tolerance: 0.001
```

`package_state` accepts `active`, `queued`, `inactive`, or `absent` and polls
the BFF because entitlement creation and transitions are asynchronous. Its
default timeout is 30 seconds with a 5-second poll.
`package_business_metrics` defaults to 60 seconds with a 10-second poll because
the dashboard read model may converge more slowly. Both defaults can be
overridden with `timeout_seconds` and `poll_seconds`.
`settled` accepts the backend statuses `completed` and `success`.
`payment_method` optionally verifies the GraphQL `paymentMethod` value.

Data-package GraphQL effect checks:

```yaml
- type: package_fields_equal
  package: minute_plan

- type: package_assignment_chain
  ues: all
  package: initial_plan
  other_package: next_plan
  expected_count: 2

- type: package_business_metrics
  package: next_plan
  expected_value: 5.00
  expected_count: 1

- type: package_visible
  package: organization_plan
  network: net-001
```

`package_fields_equal` compares the direct `getPackage` fields with the scenario plan. `package_catalog_equals` remains accepted as a temporary alias for existing scenarios.
`package_assignment_chain` reads `getPackagesForSim`, verifies the expected
assignment count and ordering, and rejects overlapping ranges; `expected: all`
checks the full returned chain. `package_business_metrics` polls the BFF
packages dashboard and compares package revenue (`expected_value`) and active
attachment count (`expected_count`). `package_visible` and `package_hidden`
verify catalog scope from a selected network. `package_name_available` proves a
rejected invalid mutation did not reserve its attempted name. `sim_unallocated`
uses the BFF SIM list to verify that a rejected initial allocation left the SIM
outside a subscriber/network assignment.

Checks may set `immediate: true`. Such checks execute at their phase position
instead of being deferred until after traffic reconciliation. This is intended
for transition-boundary assertions whose state could change again while CDRs
are being processed.

Reusable console checks use the same direct BFF calls as the console:

```yaml
- type: list_count_equals
  target: nodes
  networks: net-001
  expected_count: 3

- type: entity_fields_equal
  entity: site
  ref: site-001

- type: entity_reconciles
  entity: node
  ref: node-001

- type: node_status_equals
  type_selector: controller
  connectivity: Online
  state: Operational
```

`list_count_equals` counts the full current BFF list instead of counting only
resources already known to the generated world. Supported targets are
`networks`, `sites`, `nodes`, `customers`/`subscribers`, `plans`/`packages`,
and `sims`. For every target except `networks`, select exactly one network
because a console list is always scoped to one selected network.

`entity_fields_equal` reads the entity's direct detail query and compares the
visible identity fields with the generated world. Supported entities are
`network`, `site`, `node`, `customer`/`subscriber`, and `plan`/`package`.
`entity_reconciles` reads both the direct list and detail queries and verifies
that their common fields agree.

`node_status_equals` keeps connectivity and lifecycle state separate. Either
field may be omitted, or both may be asserted together. It polls until all
selected nodes match or `timeout_seconds` expires.

Software and operation-state checks:

```yaml
- type: software_status_equals
  type_selector: controller
  app: example
  status: update_in_progress
  desired_version: ${ULAB_SOFTWARE_TARGET_VERSION}

- type: software_count_equals
  type_selector: controller
  expected_count: 4

- type: node_operation_status_equals
  type_selector: controller
  busy: true
  operation_type: software_update
  operation_status: RUNNING

- type: site_operation_status_equals
  sites: site-001
  busy: true
  degraded: false
  restart_available: false
  rf_available: true
  service_available: false
```

`software_status_equals` reads `getSoftwares`, which is the console-facing
software status source. It can validate `status`, `current_version`, and
`desired_version`. `software_count_equals` checks the full software-row count
returned for each selected node. Runtime version and health should still be confirmed with
`node_version_equals` and `node_health_ok`, which use `getApps`.

`node_operation_status_equals` reads `getNodeOperationStatus`.
`site_operation_status_equals` reads `getSiteOperationStatus` and can validate
site busy/degraded state, per-action availability, and an active node
operation's type/status. A repeated read is how a scenario represents a
console refresh during an operation.

KPI state and chart checks:

```yaml
- type: kpi_state_equals
  key: REVENUE
  networks: net-001
  span: daily
  value_state: fresh
  max_age_seconds: 120

- type: kpi_timeseries
  key: REVENUE
  networks: net-001
  span: daily
  from: 2026-07-01T00:00:00Z
  to: 2026-08-01T00:00:00Z
  expected_count: 31
  expected_value: 25.00
  comparator: equals
  require_computed_at: true
```

`kpi_state_equals` supports `missing`, `unavailable`, `no_data`, `present`,
`available`, `zero`, `non_zero`, `fresh`, and `stale`. Fresh/stale checks use
`max_age_seconds` and the BFF `computedAt` timestamp.

`kpi_timeseries` reads `getKpiTimeSeries`, verifies ordered bucket boundaries,
and can assert bucket count, summed value, computed timestamps, and one of the
states `empty`, `present`, `all_zero`, or `non_zero`. It intentionally validates
BFF data rather than chart pixels.

Console analytics checks always use BFF `getKpiValues` or
`getPerformanceReport`; scenarios never call the analytics backend directly:

```yaml
- type: kpi_value
  key: PACKAGE_SALES
  span: daily
  op: SUM
  networks: all
  scope_key: package_id
  package: next_plan
  expected_value: 1
  timeout_seconds: 600
  poll_seconds: 5

- type: performance_report_cell
  report: package_performance
  span: daily
  networks: all
  package: next_plan
  column: revenue
  expected_value: 5.00
  tolerance: 0.01
```

KPI and report checks poll until the expected value is visible or the timeout
expires. Monetary analytics values are in major currency units, like every
other monetary value on the BFF surface (detailed below).

The P0 billing, console, and usage checks also stay entirely on the BFF
GraphQL contract:

- `revenue_summary` selects a field from `RevenueOverview` (`total_paid`,
  `total_pending`, `month_paid`, `previous_month_paid`, or `mom_pct`).
- `subscriber_billing_summary` checks settled count and amount from
  `SubscriberDetail.billing`; `payment_entitlement_reconciles` compares
  `GetPayments` with `getPackagesForSim`.
- `package_dashboard_metric` checks `mrr`, `arpu`, `revenue`, or
  `attach_count`; `network_overview_metric` checks the subscriber, site, and
  node summary fields shown by `NetworkHome`.
- `console_inventory_reconciles` compares the console's composite GraphQL
  views with one another for `component`, `node`, or `sim` inventory.
- `usage_aggregate` compares the model total with the scoped KPI for a `sim`,
  `subscriber`, `package`, `site`, or `network`, converting the expected value
  to the unit returned by GraphQL.
- `kpi_contract` can require `expected_partial`, `computedAt`, scope, and
  trend consistency. `kpi_rollup_consistency` compares daily, weekly, and
  monthly values. `performance_report_row` verifies row identity attributes,
  status, the plan `active` attribute, and optional before/after ordering.

`performance_report_row` distinguishes two fields that are easy to confuse:

- `status:` matches the row's sales-performance label. For
  `package_performance` the server vocabulary is `Inactive`, `No sales`,
  `Low sales`, `Active`, computed from the report spec's status rules
  (`active == false` -> `Inactive`, `sold == 0` -> `No sales`, `sold < 25` ->
  `Low sales`, otherwise `Active`). The comparison is case sensitive, and
  `Active` means "active plan with enough sales", not "the plan flag is on".
- `active:` matches the row's `active` attribute (`true`/`false`), which is the
  plan flag from the catalog. The comparison is case insensitive.

Monetary values are in major currency units everywhere on the BFF surface.
The analytics pipeline stores money as integer minor units (cents) for
arithmetic exactness, and console-bff converts at the presentation boundary
(`systems/console-bff/analytics/money.ts`), tagging the converted value with
the org currency code (e.g. `unit: "usd"`). So a 7.00 USD plan sold once is
`expected_value: 7.00` in `performance_report_cell`, matching
`revenue_summary`, `package_business_metrics` and `purchase_package amount:`.

Note that `span:` is inert for `package_performance`: the aggregator accepts
`daily|weekly|monthly` but the composer only honours the rolling tokens
(`last_24h`, `last_7d`, `last_30d`). Any calendar span falls back to the
configured report window and the response echoes that window's label (`8w`).

`status_equals` for SIMs continues to support `active` and `inactive` based
on active package assignment for existing scenarios. Use `wait_sim_status`
when validating the physical SIM status displayed by the console.

## Controlled service failures

`failure_control` enables deployment-specific failure injection without
hard-coding Kubernetes, Podman, or Compose commands into ukama-lab:

```yaml
- type: failure_control
  target: payment
  state: on

- type: purchase_package
  ues: all
  package: next_plan
  expect:
    result: failure

- type: failure_control
  target: payment
  state: off
```

Supported targets are `payment` and `software`. Configure the commands used by
the environment before running the scenario:

```sh
export ULAB_PAYMENT_FAILURE_ON_CMD='kubectl ...'
export ULAB_PAYMENT_FAILURE_OFF_CMD='kubectl ...'
export ULAB_SOFTWARE_FAILURE_ON_CMD='kubectl ...'
export ULAB_SOFTWARE_FAILURE_OFF_CMD='kubectl ...'
```

The commands are executed by `scripts/test-control.sh`. Active controls are
restored automatically during runtime cleanup, including after a failed
scenario. If the required command is not configured, the event fails clearly
instead of pretending a service failure was injected.

## Generated scenarios

Generated scenarios are plain scenario YAML files. Use:

```sh
ukama-lab generate --model sim --mode smoke --out scenarios/generated
```

Supported generator models in this phase: `org`, `network`, `site`, `node`, `sim`, `subscriber`, `package`.

Supported generator modes in this phase: `smoke`, `transition`, `negative`, `pairwise`, `full`.

Generator docs: `docs/generator.md`.

## Software update smoke scenarios

A software-only scenario may set `ues_per_site: 0` and omit packages,
subscribers, SIMs, and UEs. The virtual site bundle is still started, while
the selector chooses the node under test.

Before any virtual infrastructure or backend resources are created,
`ukama-lab` scans the scenario for `software_update` events and verifies that
each exact app/version has a non-empty `tar.gz` artifact in Hub. The Hub URL
is selected with `--hub-url` or `UKAMA_LAB_HUB_URL`.

The scenario should use exact current and target versions supplied through
environment variables:

```sh
export ULAB_SOFTWARE_CURRENT_VERSION='v0.1.0-e9aa34489'
export ULAB_SOFTWARE_TARGET_VERSION='2.0.0-lab.ge9aa34489'
```

The first phase verifies that the app is running and that the running binary
reports the expected current version:

```yaml
- name: verify_app_before_update
  checks:
    - type: node_health_ok
      type_selector: controller
      app: example
    - type: node_version_equals
      type_selector: controller
      app: example
      version: ${ULAB_SOFTWARE_CURRENT_VERSION}
```

The second phase calls BFF `updateSoftware`, then polls BFF until the running
binary reports the exact target version and the node/app are healthy again:

```yaml
- name: update_app
  events:
    - type: software_update
      type_selector: controller
      app: example
      tag: ${ULAB_SOFTWARE_TARGET_VERSION}
  checks:
    - type: node_version_equals
      type_selector: controller
      app: example
      version: ${ULAB_SOFTWARE_TARGET_VERSION}
    - type: node_health_ok
      type_selector: controller
      app: example
```

`node_version_equals` matches only the app's `version` returned by BFF
`getApps`; the package `tag` alone is not accepted. This prevents a false pass
when the package metadata changed but the old process is still running.

`node_health_ok` with `app` requires BFF `getNode.status.connectivity` to be
`Online` and the selected app's BFF `getApps.status` to be `running` or
`active`.

The default check timeout is 180 seconds with a five-second polling interval.
They can be adjusted with `ULAB_SOFTWARE_UPDATE_TIMEOUT_SEC` and
`ULAB_SOFTWARE_UPDATE_POLL_SEC`.
