# Scenario generator

The generator creates normal scenario YAML files from simple product model files.
Generated scenarios run through the same `validate` path as handwritten scenarios.

## Command

```sh
ukama-lab generate --model sim --mode smoke --out scenarios/generated
```

Options:

- `--model <name|all>`: `org`, `network`, `site`, `node`, `sim`, `subscriber`, `package`, or `all`
- `--mode <name|all>`: `smoke`, `transition`, `negative`, `pairwise`, `full`, or `all`
- `--models <dir>`: model directory, default `models`
- `--templates <dir>`: template directory, default `templates/generated`
- `--out <dir>`: output directory, default `scenarios/generated`

## Models in this phase

- `org`
- `network`
- `site`
- `node`
- `sim`
- `subscriber`
- `package`

Later models: member/invite, payment/billing, notification/alarm.

## Modes in this phase

- `smoke`
- `transition`
- `negative`
- `pairwise`
- `full`

Later modes: fuzz and replay.

## Templates in this phase

Implemented:

- `state-transition`
- `blocked-transition`
- `lifecycle-cleanup`
- `permission-check`
- `retry-idempotency`
- `partial-failure`
- `wrong-org-network`
- `empty-state`
- `boundary-values`
- `backend-failure`
- `runtime-effect`
- `read-model-check`

Skipped for now:

- `relationship-check`
- `dashboard-view`
- `UI-handoff`
