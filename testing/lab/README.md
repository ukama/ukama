# the lab

`ukama-lab` is a real-runtime E2E scenario runner for Ukama. It reads a
strict v1 YAML scenario, generates a deterministic test world, creates
control-plane objects through the Console BFF GraphQL API, starts nodes through a provider layer
(currently `virtual`) and starts UEs through runtime adapters, generates traffic, tracks expected state,
and validates product-visible state through backend/BFF queries.

## Build

```sh
make
```

## Validate

```sh
export UKAMA_LAB_BFF=http://localhost:4000/graphql
bin/ukama-lab validate scenarios/smoke/usage-accumulation.yaml
```


## Scenario language additions

Phase-2 adds expected event failures and backend-oriented checks:

```text
backend_count
list_contains / list_excludes
status_equals
traffic_allowed / traffic_blocked
```

`count` is kept as an alias for `backend_count`.

## Scenario generator

Generate normal scenario YAMLs from product models:

```sh
ukama-lab generate --model sim --mode smoke --out scenarios/generated
ukama-lab generate --model all --mode full --out scenarios/generated
```

Models live in `models/`. Templates live in `templates/generated/`. Generated scenarios use the same `validate` path as handwritten scenarios.

