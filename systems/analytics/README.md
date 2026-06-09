# Analytics System

Analytics owns every KPI formula consumed by the console. It ingests events from other systems via `msgclient-analytics`, backfills snapshots from source system API gateways, and serves pre-computed KPIs over REST (`/v1/analytics/*`) through its api-gateway. The console-bff is the first consumer; the UI and BFF never calculate KPIs.

See [`../../analytics-kpi-plan.md`](../../analytics-kpi-plan.md) for the full design and [`docs/schema.md`](docs/schema.md) for the DB contract.

## Services

| Service | Role |
|---|---|
| `api-gateway` | REST frontend (`/v1/analytics/business|customers|network/*`, collector admin routes). gRPC client to internal services. No KPI logic. |
| `business` | Business/operator KPIs: home, sales overview, package performance, billing summary, business sites, inventory readiness. Read-only DB access. |
| `customer` | Customer lifecycle & support: overview, list/search/detail, support diagnosis, SIMs, SIM pool. Read-only DB access. |
| `network` | Network ops: overview, topology, sites, nodes, node pool, radio, backhaul, power, alarms, metrics, events, support search. Read-only DB access. |
| `collector` | The only DB writer. Consumes events (idempotent), refreshes snapshots from source gateways, rebuilds rollups, owns migrations, demo seed. |

## Build

Each service:

```
make gen     # protoc + mockery (requires protoc, protoc-gen-go, protoc-gen-go-grpc, govalidators, mockery)
go mod tidy
make         # builds bin/<service>
make test
```

System: `make` at this level builds all services (subdirs pattern). `docker-compose build && docker-compose up` to run (requires services_ukama-net network, RabbitMQ, and the init system, same as other Ukama systems).

## Smoke test

```
curl http://localhost:8085/v1/analytics/business/home
curl http://localhost:8085/v1/analytics/customers/overview
curl http://localhost:8085/v1/analytics/network/overview
curl http://localhost:8085/v1/analytics/collector/state
```
