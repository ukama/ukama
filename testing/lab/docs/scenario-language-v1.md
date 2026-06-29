# Scenario language v1

Top-level fields: `version`, `name`, `seed`, optional `suite`, `priority`,
`tags`, `status`, optional `provider`, then `world`, `packages`, `setup`,
`runtime`, optional `profiles`, `phases`, and `final_checks`.

The language is strict. Unknown event/check names must fail validation.
Packages use `duration_days`.

Provider block is optional. Missing provider defaults to `virtual`.

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
- `wait_nodes_ready`
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
- `balance_non_negative`


Event expected failure:

```yaml
- type: restart_nodes
  nodes: all
  expect:
    result: failure
    error_contains: "script failed"
```

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
