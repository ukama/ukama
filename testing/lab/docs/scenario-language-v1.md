# Scenario language v1

Top-level fields: `version`, `name`, `seed`, optional `suite`, `priority`,
`tags`, `status`, optional `provider`, then `world`, `packages`, `setup`,
`runtime`, optional `profiles`, `phases`, and `final_checks`.

The language is strict. Unknown event/check names must fail validation.
Packages use `duration_days`.

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
- `remove_package_from_sim`
- `set_sim_status`
- `software_update`
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
- `node_state`
- `dashboard_loads`
- `node_version_equals`
- `node_health_ok`
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

