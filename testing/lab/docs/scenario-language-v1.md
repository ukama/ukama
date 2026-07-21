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
- `set_sim_status`
- `promote_release`
- `software_update`
- `disconnect_nodes`
- `reconnect_nodes`
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
- `package_catalog_equals`
- `package_visible`
- `package_hidden`
- `package_name_available`
- `package_business_metrics`
- `sim_unallocated`
- `payment_equals`
- `payment_count`
- `kpi_value`
- `kpi_trend`
- `performance_report_cell`
- `node_state`
- `dashboard_loads`
- `node_version_equals`
- `node_health_ok`
- `release_unavailable`
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

- type: set_sim_status
  ues: all
  status: inactive
```

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
  timeout_seconds: 300
  poll_seconds: 5

- type: payment_equals
  ues: all
  package: next_plan
  status: settled
  currency: USD
  expected_value: 5.00
  tolerance: 0.001
```

`package_state` accepts `active`, `queued`, `inactive`, or `absent` and polls
the BFF because entitlement creation and transitions are asynchronous.
`settled` accepts the backend statuses `completed` and `success`.

Data-package GraphQL effect checks:

```yaml
- type: package_catalog_equals
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

`package_catalog_equals` compares BFF catalog fields with the scenario plan.
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
  expected_value: 500
  tolerance: 0.01
```

KPI and report checks poll until the expected value is visible or the timeout
expires. Monetary analytics values follow the BFF/console unit in the returned
cell or KPI (currently minor units when `unit` is `cents`).

`status_equals` for SIMs supports `active` and `inactive` based on active
package assignment. `set_sim_status` only validates the mutation path in this
build; runtime enforcement is tested separately.

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
