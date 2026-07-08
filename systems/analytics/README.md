# Ukama Analytics System (v2)

Windowed KPI pipeline. Design doc: [docs/analytics-v2-plan.md](docs/analytics-v2-plan.md).

```
source api-gateways ──pull──> INGEST ──window.ready──> ANALYSIS ──kpi.computed──> AGGREGATOR <──REST── api-gateway
                              raw_records              kpi_windows                kpi_rollups
```

| Service | Role |
|---|---|
| `schema` | Shared library: window grid, GORM models, spec types/validation, pipeline event payloads |
| `ingest` | Pulls data from other systems per source specs (`ingest/configs/sources/*.yaml`) on the epoch-aligned window grid; raw zone writer |
| `analysis` | Runs one algo per KPI per window (`analysis/configs/kpis/*.yaml` + `pkg/algos` registry); KPI zone writer |
| `aggregator` | Rolls KPI windows up to daily/weekly/monthly with SUM/AVG/MIN/MAX/COUNT/LAST + trends; serves the read gRPC API |
| `api-gateway` | REST facade over aggregator (fizz/tonic, OpenAPI) |

## MVP KPIs

- `SITES_ONLINE` (scope `network_id`) — sites with ≥1 online node, from `registry.node.list`.
- `ACTIVE_CUSTOMERS` (scope `network_id`) — subscribers with ≥1 active SIM, from `subscriber.registry.getByNetwork` (fans out per network via `for_each` over `registry.network.getAll`).

## Inter-service communication

- Fast path: internal msgbus events via msgClient — `...analytics.ingest.window.ready` and `...analytics.analysis.kpi.computed` (payloads are `structpb.Struct`, no custom protos).
- Source of truth: the `window_ledger` table. Every stage claims `(kind, key, org, window)` rows (`FOR UPDATE SKIP LOCKED`) before working — this dedupes pulls per window run and survives restarts/replicas. Sweepers in analysis (60s) and aggregator (120s) recover anything whose event was lost.

## Build

```bash
# 1. Generate aggregator gRPC code (requires protoc + protoc-gen-go(-grpc)):
cd aggregator && make gen

# 2. Tidy modules (resolves deps; go.sum files are created here):
for d in schema ingest analysis aggregator api-gateway; do (cd $d && go mod tidy); done

# 3. Build everything:
make build
```

Note: `pb/gen` for the aggregator and all `go.sum` files must be generated on a
machine with the Go toolchain + protoc (they are not committed by this change).

## Run

```bash
# shared infra (rabbitmq etc.) from systems/services first, then:
ORGNAME=<org> ORGID=<id> LOCAL_HOST_IP=<ip> docker-compose up --build
```

## Smoke test

```bash
curl 'http://localhost:8085/v1/analytics/kpis'
curl 'http://localhost:8085/v1/analytics/kpis/values?keys=SITES_ONLINE,ACTIVE_CUSTOMERS&span=daily'
curl 'http://localhost:8085/v1/analytics/kpis/timeseries?keys=ACTIVE_CUSTOMERS&span=daily&op=AVG'
curl 'http://localhost:8085/v1/analytics/kpis/breakdown?key=SITES_ONLINE&by=network_id&top=10'
```

## Adding a KPI

1. If it needs new data: add a pull to a source spec (or a new spec file) in `ingest/configs/sources/` — one static dataset key per endpoint.
2. Write the algo in `analysis/pkg/algos/` and register it in `Default()`.
3. Add the KPI spec yaml to `analysis/configs/kpis/` AND `aggregator/configs/kpis/` (aggregator reads rollup_ops/metadata from it).
4. Nothing else changes — aggregator, gateway and API are KPI-generic.

## MVP simplifications (vs. the full plan)

- Pull auth: BypassAuth inside ukama-net (plan decision #7).
- Partial-span trends use simple full-compare flagged `is_partial` (plan default `elapsed` mode lands in Phase 2).
- for_each iteration retries re-run the whole dataset window (idempotent via dedup); per-iteration retry is Phase 2.
- No TimescaleDB yet: plain Postgres (schema stays compatible; see plan §2).
- Admin RPCs (repull/recompute) and metrics push are stubs for Phase 2.
