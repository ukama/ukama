# the lab

`ukama-lab` is a real-runtime E2E scenario runner for Ukama. It reads a
strict v1 YAML scenario, generates a deterministic test world, creates
control-plane objects through the Console BFF GraphQL API, starts real virtual
nodes/UEs through runtime adapters, generates traffic, tracks expected state,
and validates product-visible state through BFF queries.

## Build

```sh
make
```

## Dry run

```sh
bin/ukama-lab dry-run scenarios/smoke/usage-accumulation.yaml --print-world
```

## Validate

```sh
export UKAMA_LAB_BFF=http://localhost:4000/graphql
bin/ukama-lab validate scenarios/smoke/usage-accumulation.yaml
```
