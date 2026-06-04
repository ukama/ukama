# Analytics Network Service

Read-only gRPC service serving network analytics KPIs for the Ukama console
(via the analytics api-gateway and console-bff).

- **Module**: `github.com/ukama/ukama/systems/analytics/network`
- **Database**: shared `analytics` database; the collector service is the
  only writer and the only service running AutoMigrate. This service mirrors
  the model structs it reads (`pkg/db/model.go`) and never writes.
- **MsgBus**: client is created for service lifecycle registration only; no
  listener routes (read-only, no event handling).

## RPCs

`GetOverview`, `GetTopology`, `GetSites`, `GetSite`, `GetNodes`, `GetNode`,
`GetNodePool`, `GetRadio`, `GetBackhaul`, `GetPower`, `GetAlarms`,
`GetMetrics`, `GetEvents`, `SupportSearch` — see `pb/network.proto`.

## Derivations

- Network status (`pkg/server/status.go`): critical if any open critical
  alarm or >20% sites offline; degraded if any open alarm or any site
  degraded/offline; else healthy.
- Support recommendation: escalate if offline >30min; restart if
  needs_attention/degraded/offline; else none.
- Flags from config thresholds: `NetworkLatencyThresholdMs` (default 100),
  `BatteryCriticalPercent` (default 20), `TelemetryFreshSeconds` (default 600).

## Build

```
make gen   # generate pb/gen from pb/network.proto (requires protoc)
make test
make build
```
