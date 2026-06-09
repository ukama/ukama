# Analytics API Gateway

REST gateway for the Ukama analytics system. Translates `/v1/analytics/*` HTTP
calls into gRPC calls against the internal `business`, `customer`, `network`
and `collector` services. Contains no KPI logic.

## Ports

| Port | Purpose |
|------|---------|
| 8080 | REST API (`Server.Port`) |
| 10250 | Prometheus metrics (common default, `Metrics`) |

## Configuration (env)

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVICES_TIMEOUT` | `20s` | gRPC call timeout |
| `SERVICES_BUSINESS` | `business:9090` | business service gRPC endpoint |
| `SERVICES_CUSTOMER` | `customer:9090` | customer service gRPC endpoint |
| `SERVICES_NETWORK` | `network:9090` | network service gRPC endpoint |
| `SERVICES_COLLECTOR` | `collector:9090` | collector service gRPC endpoint |
| `AUTH_*` | — | auth host config (`LoadAuthHostConfig`), incl. `BYPASSAUTHMODE` |
| `DEBUGMODE` | `false` | gin debug mode; also enables the `/collector/seed-demo` route |

## Routes

Common query params: `period` (today|week|month|custom), `from`/`to`
(RFC3339), `timezone`; list endpoints also take `page`, `page_size`.

### Business
- `GET /v1/analytics/business/home`
- `GET /v1/analytics/business/sales/overview`
- `GET /v1/analytics/business/sales/packages` (alias: `GET /v1/analytics/business/packages`)
- `GET /v1/analytics/business/billing`
- `GET /v1/analytics/business/sites`
- `GET /v1/analytics/business/sites/:site_id`
- `GET /v1/analytics/business/inventory`

### Customers
- `GET /v1/analytics/customers/overview`
- `GET /v1/analytics/customers/list`
- `GET /v1/analytics/customers/search?q=...`
- `GET /v1/analytics/customers/sims`
- `GET /v1/analytics/customers/sim-pool`
- `GET /v1/analytics/customers/:customer_id`
- `GET /v1/analytics/customers/:customer_id/support`

Note: gin gives static segments (`overview`, `list`, `search`, `sims`,
`sim-pool`) priority over the `:customer_id` param sibling, so customer ids
(UUIDs) never collide with the static routes.

### Network
- `GET /v1/analytics/network/overview`
- `GET /v1/analytics/network/topology`
- `GET /v1/analytics/network/sites`
- `GET /v1/analytics/network/sites/:site_id`
- `GET /v1/analytics/network/nodes`
- `GET /v1/analytics/network/nodes/:node_id`
- `GET /v1/analytics/network/node-pool`
- `GET /v1/analytics/network/radio`
- `GET /v1/analytics/network/backhaul`
- `GET /v1/analytics/network/power`
- `GET /v1/analytics/network/alarms`
- `GET /v1/analytics/network/metrics`
- `GET /v1/analytics/network/events`
- `GET /v1/analytics/network/support/search?q=...`

### Collector
- `POST /v1/analytics/collector/refresh` `{ "source": "registry|subscriber|dataplan|metrics|node|inventory|billing|all" }`
- `GET /v1/analytics/collector/state`
- `POST /v1/analytics/collector/rollups/rebuild` `{ "family": "business|customer|network|all", "from": "...", "to": "..." }`
- `POST /v1/analytics/collector/seed-demo` (debug mode only)

Also exposed: `GET /ping`, `GET /openapi.json`, Swagger UI at `/swagger`.

## Development

```sh
make server   # run locally
make test     # unit tests
make build    # static linux binary in bin/
```
