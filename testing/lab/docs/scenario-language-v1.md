# Scenario language v1

Top-level fields: `version`, `name`, `seed`, optional `suite`, `priority`,
`tags`, `status`, then `world`, `packages`, `setup`, `runtime`, optional
`profiles`, `phases`, and `final_checks`.

The language is strict. Unknown event/check names must fail validation.
Packages use `duration_days` only.

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

- `model_count` (`count` is accepted as a legacy alias)
- `bff_count` (reserved/skipped in Phase 0)
- `node_ready`
- `ue_attached`
- `usage_per_sim`
- `usage_sample`
- `package_active`
- `package_remaining` (skipped until BFF exposes remaining balance)
- `node_state`
- `dashboard_loads`
- `balance_non_negative`
