# Phase-6 generator model

Phase-6 generated scenarios must be meaningful product cases, not blind Cartesian matrix combinations.

## Topologies

Use only these topology names:

- `simple`: 1 network, 1 site, 1 tower + 1 amplifier + 1 controller, 1 UE, 1 data package, 1GB traffic per UE.
- `medium`: 1 network, 3 sites, 1 tower + 1 amplifier + 1 controller per site, 30 UEs per tower, 5 data packages, 2GB traffic per UE.
- `large`: 1 network, 10 sites, 1 tower + 1 amplifier + 1 controller per site, 100 UEs per tower, 10 data packages, 5GB traffic per UE.

## Families

Use product family names only:

- `smoke`
- `backend`
- `usage`
- `sim_pool`
- `package`
- `subscriber`
- `lifecycle`
- `node_ops`
- `site_ops`
- `software_update`
- `failure`
- `scale`

Do not use `dashboard` for Phase-6 backend validation. UI dashboard validation is separate and later.

## Cases

Each family defines explicit named `cases`. A generated scenario must come from a real case name.

Example:

```yaml
family: sim_pool
cases:
  - name: sim_pool_empty_blocks_allocation
    topology: simple
    priority: p1
    tags: [sim_pool, negative]
    events: [allocate_sim]
    checks: [expected_failure, sim_pool_count_zero]
```

The old compatibility fields (`flows`, `topologies`, `scales`, `runtime`, `failures`, `verification`) remain only until the generator is switched to case-based generation.
