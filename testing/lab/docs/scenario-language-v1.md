# Scenario language v1

Top-level fields: `version`, `name`, `seed`, `world`, `packages`, `setup`,
`runtime`, optional `profiles`, `phases`, and `final_checks`.

The language is strict. Unknown event/check names must fail validation.

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

- `count`
- `node_ready`
- `ue_attached`
- `usage_per_sim`
- `usage_sample`
- `package_active`
- `package_remaining` (skipped until BFF exposes remaining balance)
- `node_state`
- `dashboard_loads`
- `balance_non_negative`
