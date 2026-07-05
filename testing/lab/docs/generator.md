# Generated scenario catalog

Phase-6 generation is case-based, not Cartesian-product based.

The generator reads:

- `models/profiles/topology.yaml`
- `models/profiles/scale.yaml`
- `models/families/*.yaml`

It writes one scenario per explicit family case:

```sh
./bin/ukama-lab generate --model all --models models --out scenarios/generated
```

Generated output layout:

```text
scenarios/generated/
  backend/
  smoke/
  usage/
  sim_pool/
  package/
  subscriber/
  lifecycle/
  node_ops/
  site_ops/
  software_update/
  failure/
  scale/
  index.yaml
```

Topology names are meaningful and fixed:

- `simple`: 1 network, 1 site, 1 tower/amplifier/controller, 1 UE, 1 package, 1GB per UE
- `medium`: 1 network, 3 sites, 1 tower/amplifier/controller per site, 30 UEs per site, 5 packages, 2GB per UE
- `large`: 1 network, 10 sites, 1 tower/amplifier/controller per site, 100 UEs per site, 10 packages, 5GB per UE

Family files define explicit cases only. Example:

```yaml
family: usage
cases:
  - name: usage_simple_one_ue_1gb
    topology: simple
    priority: p1
    tags: [usage, runtime, simple]
    events: [traffic_1gb]
    checks: [usage_per_sim, balance_non_negative]
```

The generator marks unsupported cases as `status: wip` instead of inventing fake test behavior. WIP scenarios validate but are skipped by the runner.

Run examples:

```sh
./bin/ukama-lab validate scenarios/generated
./bin/ukama-lab run scenarios --suite generated --priority p0
./bin/ukama-lab run scenarios --tag usage
```
