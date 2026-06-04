# Analytics Architecture

The analytics system now uses three clear runtime roles:

```text
console / bff
    |
    v
api-gateway-analytics
    |
    | thin REST -> gRPC gateway only
    v
analytics
    |
    | read-only KPI/domain queries
    v
postgres analytics database
    ^
    | writes snapshots, facts, intervals and rollups
    |
collector
```

## Services

### api-gateway-analytics

The API gateway is only a gateway. It owns HTTP routing, auth middleware,
request/response mapping, and calls downstream gRPC services.

It does not own KPI logic, repositories, schema access, rollups, events, or
refresh behavior.

### analytics

The `analytics` service is the single read-only analytics service behind the
gateway. It registers the business, customer and network gRPC APIs on one gRPC
server.

This removes the three separate read-service containers while keeping the
gateway clean and small.

The read domains remain separated at package/API level:

```text
business -> business KPIs and views
customer -> customer/SIM KPIs and views
network  -> network/site/node/radio/power views
```

The service must be treated as read-only. The collector owns schema creation and
all writes.

### collector

The collector is the worker/writer process. It owns:

- schema migration
- event consumption
- source refresh
- snapshot updates
- fact inserts
- interval updates
- dirty rollup marking
- scheduled rollup rebuilds

## Storage layers

The database is organized as four logical layers:

```text
raw_events   -> original event/audit trail
snapshots    -> current known state
facts        -> append-only activity records
rollups      -> precomputed hourly/daily KPI tables
```

Rules:

- raw events are immutable
- facts are append-only
- snapshots can be updated
- rollups can be rebuilt

This keeps analytics recoverable: rollups can be rebuilt from facts, and future
replay support can rebuild facts/snapshots from raw events.
