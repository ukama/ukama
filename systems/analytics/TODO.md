# Analytics TODO

Items intentionally deferred from this pass.

## Phase 2

- Enforce database write ownership with separate DB users:
  - `analytics_writer` for `collector` only.
  - `analytics_reader` for the `analytics` read service only.
- Split public analytics routes from internal/admin collector routes:
  - public: `/v1/analytics/*`
  - admin/internal: `/v1/admin/analytics/*`
- Move source refresh code behind formal source adapters:
  - registry
  - subscriber
  - dataplan
  - metrics
  - node
  - inventory
  - billing
- Add a formal KPI contract document for every KPI:
  - source table/event
  - formula
  - window rule
  - timezone rule
  - freshness rule
  - limitations
- Add first-class replay support for analytics events/facts.
- Add retention and partitioning for high-volume tables:
  - raw metric samples
  - event logs
  - fact tables
- Introduce a platform-wide event envelope:
  - event_id
  - event_type
  - producer
  - aggregate_type
  - aggregate_id
  - occurred_at
  - schema_version
  - sequence/version
  - payload

## Rollup gaps

The scheduler is now automatic for the rollups that already have rebuild
functions. Add rebuild functions for the remaining rollups:

- business_site_daily
- business_inventory_daily
- network_health_hourly
- site_health_hourly
- node_health_hourly
- radio_hourly
- backhaul_hourly
- power_hourly

## Event processing hardening

The next correctness pass should make event handling fully transactional:

```text
receive event
validate payload
BEGIN
  insert analytics_event_logs idempotency row
  if duplicate: COMMIT and ACK
  update snapshots/facts/intervals
  mark dependent rollups dirty
COMMIT
ACK
```

Malformed events should be ACKed and recorded in `analytics_event_errors`.
Processing failures should not be silently ACKed; they should retry or go to a
DLQ/quarantine path.
