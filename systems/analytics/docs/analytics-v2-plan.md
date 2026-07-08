# Ukama Analytics System — v2 Design Plan

**Status:** Proposal (rev 3) · **Date:** 2026-07-07
**Scope:** Ground-up redesign of `systems/analytics` as a windowed KPI pipeline: **Ingest → Analysis → Aggregator**, fronted by one api-gateway. Existing code is reference-only.

---

## 1. Goals & Principles

- Deep, per-resource insight for org operators (Console via console-bff) and Ukama platform teams: orgs, users, members, networks, sites, nodes, subscribers, SIMs, packages, payments, invoices, data usage, alarms, telemetry.
- KPIs are the product: every output is `<KPI_KEY>: <value>` + self-describing metadata (from/to, type, unit, symbol).
- Rollups at **daily, weekly (ISO-8601), monthly (calendar)** spans — span set is config, so quarterly/annual are additive later.
- **Deterministic:** same input windows → same KPIs → same rollups, on any host, after any restart or replay.
- **Modular & abstracted:** each stage is spec/config-driven; adding a source or KPI is YAML + (at most) one registered Go function — no engine changes.
- Fits Ukama conventions: Go gRPC services + api-gateway (fizz/tonic), msgClient/RabbitMQ, GORM/Postgres, protoc+mockery Makefiles, docker-compose, MPL-2.0 headers.

Non-goals now: sub-second streaming, ML forecasting (Phase 4+), replacing `systems/metrics` (Prometheus stays the ops-monitoring path; analytics samples summarized KPIs from it).

---

## 2. Architecture

```
 source systems (registry, subscriber, data-plan, billing, nucleus,
                 node, messaging, inventory, ukama-agent, metrics)
        │ REST pulls via system api-gateways (common/rest/client)
        ▼
 ┌─────────────────────────── INGEST (writer of raw zone) ───────────────────────────┐
 │ source specs (YAML): what to pull, cadence, params, mapping                       │
 │ pulls each eligible window per spec → epoch-aligned grid → raw_records            │
 │ publishes analytics.ingest.window.ready ; updates window_ledger                   │
 └───────────────┬───────────────────────────────────────────────────────────────────┘
                 ▼
 ┌─────────────────────────── ANALYSIS (writer of KPI zone) ─────────────────────────┐
 │ KPI specs (YAML) + algo registry (Go, versioned)                                  │
 │ consumes window.ready → runs algo on that window's raw data                       │
 │ → kpi_windows (value + components: sum/count/min/max) + metadata                  │
 │ publishes analytics.analysis.kpi.computed ; updates window_ledger                 │
 └───────────────┬───────────────────────────────────────────────────────────────────┘
                 ▼
 ┌─────────────────────────── AGGREGATOR (writer of rollup zone, READ API) ──────────┐
 │ span configs: daily/weekly/monthly (org tz) ; ops: SUM AVG MIN MAX COUNT LAST     │
 │ consumes kpi.computed → weighted rollups from components → kpi_rollups            │
 │ computes TRENDS vs previous span (direction + %); serves all read RPCs            │
 └───────────────┬───────────────────────────────────────────────────────────────────┘
                 ▼
        api-gateway (REST, auth, org scoping) → console-bff → Console
```

| Service | Writes | Reads | Runtime |
|---|---|---|---|
| `ingest` | `raw_records`, `window_ledger`, `source_cursors`, `ingest_errors` | source APIs | gRPC (admin) + scheduler (+ msgClient publish) |
| `analysis` | `kpi_windows`, `analysis_errors` | `raw_records`, `window_ledger` | gRPC (admin) + msgClient consumer |
| `aggregator` | `kpi_rollups` | `kpi_windows`, `window_ledger` | gRPC (read + admin) + msgClient consumer |
| `api-gateway` | — | aggregator gRPC | HTTP (fizz/tonic, OpenAPI) |

One Postgres instance; each service has its own DB role restricted to its zone (writer isolation = modularity + blast-radius control). All tables carry `org_id` (multi-tenant); gateway enforces org scoping from auth context.

**TimescaleDB — scoped, optional optimization (not an architectural dependency).** Since the KPI engine does all computation, Timescale contributes nothing to compute; its value is storage lifecycle on the one high-volume table, `raw_records`:
- *chunked time-partitioning* — Analysis always reads "one window"; chunk pruning keeps that fast at any total size;
- *columnar compression* (~90%+) — makes 180-day raw retention affordable;
- *retention as policy* — `drop_chunks` instead of DELETE+vacuum churn.

Cost: non-vanilla image (`timescale/timescaledb`), hypertable calls in migrations. Rules: schema stays vanilla-Postgres-compatible; Timescale setup lives in guarded migration steps; only `raw_records` becomes a hypertable (plus `kpi_windows` later if volume warrants); `kpi_rollups` never needs it. Fallback without Timescale: native daily range partitioning on `raw_records` + partition-drop job (no compression). Decision driver: ≥ ~1M raw rows/day (e.g. 10k SIMs at 1-min windows) → use Timescale; small deployments run stock Postgres unchanged.

**Pipeline coordination — events + ledger (never events alone):**
- Fast path: internal msgbus events `analytics.ingest.window.ready {window_id, source, org}` and `analytics.analysis.kpi.computed {window_id, kpi_key, org}` trigger the next stage within seconds.
- Source of truth: `window_ledger(window_id, source/kpi, org, stage, status, attempt, updated_at)` — a state machine `pulled → analyzed → aggregated` (+ `failed`, `dirty`). Each service also runs a low-frequency sweeper over the ledger to pick up anything whose event was lost. This yields at-least-once triggering with exactly-once effect (all writes are idempotent upserts keyed by window).

---

## 3. Window Model (shared contract)

Defined once in the `schema` module; identical for all three services.

- **Base window `W`** (config; **default 5min** — decision #2, §11): the atomic unit of the pipeline. Boundaries are **epoch-aligned**: window `N` = `[N·W, (N+1)·W)` UTC. Window ID = `N` (deterministic everywhere; never derived from timers or process start).
- **Eligibility (simplified — watermark removed):** every window is pulled as soon as it closes; there is no separate watermark lag. *(Revised decision: an earlier draft had a lag `L` for late source data; dropped for simplicity.)*
- **Late data:** any data a source finalizes after its window was pulled is handled by marking the window `dirty` in the ledger → Analysis recomputes it → Aggregator recomputes the containing day/week/month. The dirty-cascade is the single mechanism for late data, corrections, backfill, and algo upgrades.
- **Pull cadence is ingest-level config, not spec-level:** all pulls run on the global base window `W` (from the shared pipeline config). Specs declare *what* to pull, never *when* — keeps every dataset on one deterministic grid and specs purely declarative. (Per-source cadence multiples are a possible future config extension, not in specs.)
- **Config governance:** `W` lives in one shared config (env-injected to all services from compose/helm) and are recorded in `pipeline_config` with a version; changing `W` starts a new **grid epoch** — old windows stay valid under their epoch, new windows use the new grid (no silent invalidation).
- Spans (day/week/month) are defined in **org timezone**, computed by Aggregator over UTC base windows; span set is a config list (`daily, weekly, monthly` now; `quarterly, annual` later — zero code change).

---

## 4. Service 1 — Ingest

Generic engine driven by per-source YAML specs (`configs/sources/*.yaml`). **API pulling only** (v2 decision): Analysis needs bulk window data, and pulls provide it directly; msgbus event ingestion is deliberately out of scope. The spec format reserves a future `events:` source type — adding it later is an engine extension, not a redesign.

```yaml
version: 1
source: subscriber
system: subscriber              # logical name → resolved via initclient to apiGwURL
# base_url: http://...          # optional override, skips initclient (dev/direct-only systems)

pulls:
  - key: ukamaagent.usage.getForPeriod   # static dataset key — unique across ALL specs
    client: ukamaagent          # common/rest/client package
    endpoint: /v1/usages
    strategy: window            # window | incremental | full_snapshot
    params: { from: "{{.WindowStart}}", to: "{{.WindowEnd}}" }
    rate_limit: 10rps
    map: { sim_id: $.sim_id, bytes: $.total_bytes, at: $.end_time }
  - key: registry.network.getAll
    client: registry
    endpoint: /v1/networks
    strategy: full_snapshot
    map: { network_id: $.id, name: $.name, status: $.status }
```

Pull interval/window logic lives in **ingest service config** (shared pipeline config: `W`), not in specs — specs stay purely declarative (*what*, not *when*). Pagination is intentionally absent: source api-gateways don't paginate yet; when they do, it becomes an engine-level capability (spec opt-in), not a per-spec invention.

- **Raw zone:** `raw_records(id, org_id, dataset_key, window_id, event_time, payload jsonb, fields jsonb, ingested_at)` — hypertable on `event_time`. `payload` = untouched original (bronze, replayable); `fields` = mapped/normalized columns per spec. Records are assigned to windows by **event time**, not arrival time.
- **Address resolution (initclient):** the spec's `system:` field is a logical name — never a URL. Before pulling, Ingest resolves it through the init system, following the repo-wide pattern (see `analytics/collector/cmd/server/main.go:129`, `notification/event-notify/cmd/server/main.go:96`):
  1. `ic.NewInitClient(cfg.Http.InitClient)` — init api-gateway host from config (`Http.InitClient`, env `HTTP_INITCLIENT`, default `api-gateway-init:8080`).
  2. `ic.GetHostAddress(client, ic.CreateHostString(org, spec.System), &org)` → `GET /v1/orgs/{org}/systems/{system}` on init/lookup → `SystemIPInfo{ApiGwUrl, ApiGwIp, ApiGwPort, …}`. `GetHostAddress` prefers the registered `ApiGwUrl` (stable hostname), falling back to `ApiGwIp:ApiGwPort`.
  3. Full pull address = `{resolvedBaseURL}{spec.endpoint}` fed into the `common/rest/client` typed client.

  Two deviations from the existing startup-time pattern (both fix known weaknesses — see the TODO in event-notify's main.go):
  - **Resolve lazily with a TTL cache**, not once at boot: resolution happens per (org, system) at first pull, cached (default 10 min), so Ingest survives systems registering late, address changes, and per-org fan-out without restarts.
  - **Re-resolve on connection failure**: a refused/unreachable pull invalidates the cache entry and retries once with a fresh lookup before counting as a pull failure in the ledger.

  A spec may set `base_url:` as an explicit override (skips initclient — used for local dev and for systems reachable only by direct address, like collector's current `InventoryClient` config).
- **Shared datasets & pull dedup (no double-fetching):** the invariant is **KPIs never trigger fetches** — pulls exist only in source specs, each identified by a **static dataset key** with convention `<system>.<resource>.<operation>` (e.g. `registry.network.getAll`, `registry.site.getAll`, `ukamaagent.usage.getForPeriod`). The key is the single handle everywhere: KPI specs reference it in `inputs`, `raw_records` rows are tagged with it, and the ledger tracks it. Dedup in two layers:
  1. *Static:* spec validation (CI + startup) enforces global key uniqueness — one endpoint, one key, one definition; every spec/KPI needing that data references the existing key.
  2. *Runtime (pulled-list per run):* per window tick the engine builds a pull plan of due dataset keys, coalesces duplicates, and claims each `(dataset_key, org, window_id)` row in `window_ledger` (`SELECT … FOR UPDATE SKIP LOCKED`) before executing — already-`pulled` or in-flight keys are skipped and the stored dataset reused. The "pulled list for a run" is thus the ledger itself (persistent), so dedup survives restarts, replays, and multiple ingest replicas — never an in-memory set.
- **Dependent pulls — `for_each` DAG:** a pull may iterate over the rows of a parent dataset, binding row fields into its endpoint/params. Chains of these form a per-window **DAG**, executed in topological stages: parents complete → children fan out one call per parent row. Example — the full chain for `NETWORK_UPTIME`:

  ```yaml
  pulls:
    # bulk list endpoint with filters — one call enumerates all nodes with lineage ids
    - key: registry.node.list
      endpoint: /v1/nodes/list
      strategy: full_snapshot
      map: { node_id: $.node_id, site_id: $.site_id, network_id: $.network_id, state: $.state }

    # for_each only where the API is genuinely per-entity
    - key: metrics.node.uptime
      endpoint: /v1/metrics/nodes/{{.node_id}}/uptime
      params: { from: "{{.WindowStart}}", to: "{{.WindowEnd}}" }
      for_each: { dataset: registry.node.list, bind: [node_id] }
      map: { uptime_s: $.uptime }
  ```

  Deep chains (networks → for-each sites → for-each nodes) remain expressible for per-entity-only APIs, but bulk **list endpoints with filters** (like `/v1/nodes/list?network_id=&site_id=&state=…`) collapse enumeration to a single pull and should always be preferred — the fan-out here is one call per node for the metric, not three levels of listing.

  Semantics:
  - **Lineage propagation:** bound fields plus the parent's own lineage are stamped onto every child row (`metrics.node.uptime` rows carry `node_id, site_id, network_id`). The `network_uptime@v1` algo then just groups one dataset by `network_id` — no joins, no extra fetches.
  - **Validation:** DAG must be acyclic and reference existing dataset keys (CI + startup).
  - **Iteration-level idempotency & retry:** each iteration's dedup key includes its bound params, and the ledger tracks per-dataset iteration counts (`done/total`). A partially failed fan-out retries **only the failed iterations**; the dataset (and thus `window.ready`) isn't marked `pulled` until all iterations land or the retry budget is exhausted (`failed`, with per-iteration errors in `ingest_errors`).
  - **Fan-out control:** N networks × M sites × K nodes calls per window — bounded worker pool per source system + spec `rate_limit`; chain latency should fit within one window `W` (sizing input for `W`). Guardrail metric: iterations per window per dataset.
  - **Prefer bulk endpoints:** `for_each` is for APIs that only serve per-entity calls. If a bulk endpoint exists (`/v1/nodes?network=…` or plain getAll), use it — one pull beats N. Spec review should treat every new `for_each` as a question: "does the source system have (or should it grow) a bulk endpoint?"
- **Pull strategies:** `window` (from/to templated per window — preferred), `incremental` (persisted cursor per (source, org) in `source_cursors`), `full_snapshot` (dimension-style reference data: sims, packages, nodes — stored as **change-log + state-as-of**, decision #4 §11: a row persists only when its content-hash changes; per-window high-water marks make deletions detectable; Analysis reads *state as of window N* = latest change ≤ N).
- **Dedup/idempotency:** natural key `(dataset_key, org, dedup_key)` unique index; `dedup_key` = source row id for windowed/incremental records, content-hash of mapped fields for snapshots (which is what implements the change-log). Re-pulls are no-ops.
- **Auth to source gateways:** BypassAuth inside `ukama-net` for now (decision #7 §11), matching existing compose `BYPASS_AUTH_MODE` usage; service credentials when the platform grows them.
- **Failure handling:** per-source retry with backoff inside the eligibility period; a window that can't complete is ledgered `failed` with error in `ingest_errors` (never silently skipped); admin RPC `RepullWindow(source, window_range)`.
- **Per-org fan-out** via initclient (same as existing collector's refresher pattern).
- Spec hygiene: JSON-schema validation at CI + startup; golden-file tests (sample API payload → expected raw_records).

## 5. Service 2 — Analysis

Runs one **algo** per KPI per window. Driven by KPI specs (`configs/kpis/*.yaml`) + versioned Go algo registry.

```yaml
kpi: USAGE_BY_NETWORK
domain: subscriber              # console grouping tag
algo: usage_by_network@v1
scope: [network_id]             # one output row per network per window
inputs:                          # reference datasets by static key (§4)
  usage: { dataset: ukamaagent.usage.getForPeriod, fields: [sim_id, bytes] }
  sims:  { dataset: subscriber.sim.getAll, fields: [sim_id, network_id] }
output: { type: float, unit: mb, symbol: MB }
rollup_ops: [SUM, AVG, MAX]     # ops Aggregator may apply (validated at load)
# lookback: 30d                 # optional (decision #6 §11): engine feeds prior windows'
#                               # datasets / kpi_windows to the algo (MTTR, churn, cohorts, runway)
```

- Trigger: `window.ready` for any input dataset of the KPI (plus ledger sweeper). A KPI computes when **all** its input dataset keys are `pulled` for that window; the ledger tracks partial readiness.
- Algo contract (pure function → deterministic, unit-testable):
  `func(win Window, in Datasets) ([]KpiResult, error)` where each result is
  `{kpi, scope{...}, value, components{sum, count, min, max}, metadata}`.
- **KPI zone:** `kpi_windows(kpi_key, org_id, scope jsonb, window_id, value, sum, count, min, max, value_type, unit, symbol, algo_version, computed_at)` — unique on `(kpi_key, org_id, scope, window_id)`; recompute upserts. **Components are mandatory**: they are what make higher-span rollups exact (weighted, not avg-of-avg).
- Metadata envelope served downstream:

```json
{
  "USAGE_BY_NETWORK": 2048,
  "metadata": {
    "from": "2026-07-06T10:30:00Z", "to": "2026-07-06T10:30:30Z",
    "type": "float", "unit": "mb", "symbol": "MB",
    "scope": {"network_id": "…"}, "algo": "usage_by_network@v1"
  }
}
```

- Algo upgrades: bump version → mark historical windows dirty (ranged) → recompute cascade; `algo_version` on every row records provenance.
- Errors: malformed/insufficient input → `analysis_errors` + ledger `failed`; never a partial silent write (whole window per KPI is one transaction).

## 6. Service 3 — Aggregator (+ read API)

Rolls `kpi_windows` up to spans and serves all reads.

- **Spans:** `daily`, `weekly` (ISO), `monthly` (calendar) in org tz — config list, extensible to quarterly/annual.
- **Operations:** `SUM, AVG, MIN, MAX, COUNT, LAST` computed **from components**, exactly:
  - `SUM = Σ sum` · `COUNT = Σ count` · `AVG = Σ sum / Σ count` (weighted — never average of window averages) · `MIN/MAX = min/max of window min/max` · `LAST = value of latest window`.
  - Each KPI's `rollup_ops` whitelist is enforced at spec load (summing a percentage or averaging a counter is rejected up front).
- Trigger: `kpi.computed` events (debounced per span bucket) + ledger sweeper; dirty base windows dirty their containing spans transitively.
- **Rollup zone:** `kpi_rollups(kpi_key, org_id, scope jsonb, span, span_start, span_end, op, value, value_type, unit, symbol, is_partial, prev_value, change_abs, change_pct, trend, computed_at)` — unique on `(kpi_key, org, scope, span, span_start, op)`. Current (incomplete) day/week/month is computed and served with `is_partial: true`.

**Trends (first-class output).** For every rollup row, Aggregator compares against the previous same-span row (same kpi/org/scope/op): day N vs day N−1, ISO week vs previous week, month vs previous month.

- `change_abs = value − prev_value` · `change_pct = change_abs / |prev_value| × 100`
- `trend = up | down | flat | new | na`
  - `flat` when `|change_pct| < flat_threshold` (config, default 1%; absolute-epsilon fallback for near-zero values)
  - `new` when no previous span exists (first day/week/month of a KPI or scope)
  - `na` when `prev_value = 0` and value ≠ 0 (percent undefined — `change_abs` still served) or previous span is missing/failed
- For **partial** current spans, naive comparison misleads (Tuesday-noon vs full Monday). Config `partial_trend_mode`: `elapsed` (default — compare against the same elapsed fraction of the previous span, recomputed from base windows, so "so far today vs same time yesterday") or `full` (simple compare, flagged by `is_partial`).
- Trend direction is polarity-aware per KPI: spec field `positive_direction: up|down` (churn going down is *good*); stored trend stays factual (`up/down`), the polarity flag lets the console color it.
- Recompute cascade covers trends automatically: a corrected previous span dirties the *following* span's trend row.

API surface: every value response includes `trend: {direction, change_pct, change_abs, prev_value, prev_span_start}` — e.g. `("up", +30%)`.
- **Read RPCs** (generic — no per-KPI code; this is what api-gateway exposes):
  - `ListKpis()` — registry: keys, domains, units, symbols, spans, scopes, ops
  - `GetKpis(keys[], span, from, to, scope_filters)` — latest values + metadata envelope
  - `GetKpiTimeSeries(keys[], span, from, to, scope_filters)`
  - `GetKpiBreakdown(key, scope_dimension, span, window, top_n)`
- Fine-grained (base-window) KPI reads can be exposed later without redesign — the data is already in `kpi_windows`.

## 7. API Gateway

```
GET  /v1/analytics/kpis                                   # registry (self-describing)
GET  /v1/analytics/kpis/values?keys=USAGE_BY_NETWORK&span=daily&from=&to=&network_id=
GET  /v1/analytics/kpis/timeseries?keys=…&span=weekly&…
GET  /v1/analytics/kpis/breakdown?key=USAGE_BY_NETWORK&by=network_id&span=monthly&top=10
GET  /v1/analytics/kpis/export?format=csv
POST /v1/analytics/admin/repull|recompute|rollup|reconcile   # admin role
```
Org scoping from authenticated context only; platform-domain KPIs gated by platform-admin role; fizz/tonic OpenAPI per repo convention. Console-bff consumes these REST routes.

---

## 8. Source & KPI Catalog

### 8.1 Sources (each = one spec file; all via `common/rest/client` + initclient resolution)
| Source | Pulls (record types) |
|---|---|
| nucleus | orgs, users (snapshots) |
| registry | networks, sites, nodes, members (snapshots + state) |
| subscriber | subscribers, sims, sim-pool (snapshots) |
| data-plan | packages, base-rates, markup (snapshots) |
| billing | invoices/reports (windowed by Period), payment statuses |
| ukama-agent | usage via `GetUsageForPeriod` (windowed from/to — the USAGE_* input) |
| node | node state, health reports |
| inventory | components, accounting |
| metrics | sampled KPI set from metrics api-gateway (throughput, active UEs, latency, battery/solar) |
| notification | notifications/alarms (windowed, by severity) |

Sources without windowed query endpoints (only "current state") are pulled as `full_snapshot`; KPIs over them are state-style (counts, statuses) and their history accrues from successive snapshots taken every base window `W`.

### 8.2 KPIs (each = one spec + one algo)
**Business:** GROSS_REVENUE, MRR, ARPU, PAYMENT_SUCCESS_RATE, PACKAGE_SALES, PACKAGE_ATTACH_RATE, COLLECTION_RATE, REVENUE_BY_{PACKAGE|NETWORK|SITE}, DATA_SOLD_VS_CONSUMED.
**Subscriber & SIM:** ACTIVE_SUBSCRIBERS, NEW_SUBSCRIBERS, CHURNED_SUBSCRIBERS, CHURN_RATE, RETENTION_COHORT, USAGE_BY_NETWORK, USAGE_PER_SUBSCRIBER (avg), TOP_CONSUMERS, SIM_POOL_UTILIZATION, SIM_POOL_RUNWAY, SIMS_BY_STATUS, PACKAGE_EXPIRIES_UPCOMING.
**Network:** NODE_AVAILABILITY, FLEET_ONLINE_RATIO, NODE_STATE_TRANSITIONS, MTBF, MTTR, ALARMS_BY_SEVERITY, HEALTH_SCORE, THROUGHPUT_{DL|UL}, BACKHAUL_LATENCY (avg/min/max), ACTIVE_UES, BATTERY_LEVEL, LOW_POWER_INCIDENTS, CAPACITY_UTILIZATION.
**Platform (cross-org — future, separate Ukama-operated instance per decision #1 §11):** ORG_GROWTH, USER_GROWTH, NETWORKS_DEPLOYED, NODES_DEPLOYED, FLEET_AVAILABILITY, PLATFORM_REVENUE, ADOPTION_FUNNEL (org → network → node online → first subscriber → first payment).

Percentile variants (p95 latency/usage) are deferred to Phase 4 (decision #5 §11). Money: integer cents in **org currency** (decision #9 §11 — converted at ingest with original amount/currency/fx retained for audit; sources store float64). Bytes: bigint. Timestamps: `timestamptz` UTC.

---

## 9. Production-Readiness Checklist

- **Determinism:** epoch-aligned windows; pure algos; component-based rollups; idempotent upserts everywhere; replay from `raw_records` reproduces identical `kpi_windows` and `kpi_rollups`.
- **Reliability:** events+ledger dual triggering; sweeper recovery; retries with backoff; `*_errors` tables + admin repull/recompute; a simulated 24h outage self-heals via ledger backlog draining (Phase 3 exit test).
- **Reconciliation:** daily count-diff of dims vs source systems (via `common/rest/client`), mismatches alert + auto-correct snapshots (audited in `reconciliation_log`).
- **Observability:** per-stage lag metrics (`newest_eligible_window - newest_processed_window`), events/s, error rates, freshness (`as_of` in every API response), pushed to pushgateway per repo convention.
- **Scalability:** stages scale independently; work is partitioned naturally by (org, window) → safe horizontal sharding; Timescale compression on `raw_records`/`kpi_windows`.
- **Retention:** raw 180d, kpi_windows 2y, kpi_rollups indefinitely; via Timescale policies (compress @7d/@30d) or partition-drop jobs on stock Postgres (§2).
- **Security:** three DB roles (zone-scoped writers, aggregator read-only on kpi zone); org scoping at gateway; no PII in raw `fields` beyond spec-mapped needs.
- **Testing:** golden-file tests per source spec and per KPI algo; ledger state-machine unit tests; integration test = docker-compose pipeline with fixture events + pulls, asserting exact rollup values.
- **Future event ingestion:** spec format reserves an `events:` source type; routing-key quirks documented for that day (payments under system `payments`; CDR keys under `ukamaagent`, absent from `common/events/event.go`).

## 10. Delivery Phases

| Phase | Scope | Exit criteria |
|---|---|---|
| **0 — Foundation** | `schema` module (window model, ledger, zones, migrations — vanilla-compatible, guarded Timescale steps), Postgres, `ingest` engine (grid scheduler, pull strategies, dedup, ledger, window.ready), service/gateway skeletons, compose + CI | one source spec lands raw windows end-to-end; ledger + events observable; every planned source's endpoints verified against its spec `strategy` (decision #8 §11) |
| **1 — Analysis + first KPIs (daily)** | algo registry, `kpi_windows`, analysis triggering; `aggregator` SUM/AVG/COUNT with daily span; business + subscriber KPI set; gateway + console-bff wiring | USAGE_BY_NETWORK, GROSS_REVENUE, ACTIVE_SUBSCRIBERS correct vs manual SQL; adding a KPI = spec+algo only |
| **2 — Weekly/Monthly + Network + Trends** | tz spans, is_partial, MIN/MAX/LAST, **trend engine** (direction/change_pct, elapsed-mode partial comparison, polarity), dirty-cascade recompute; node/alarm/telemetry sources + network KPI set | all KPIs at 3 spans with trends; late-data cascade verified |
| **3 — Hardening** | reconciler, replay tooling, DQ alerts, retention/compression, DB roles | 24h-outage self-heal test passes; replay reproduces identical rollups |
| **4 — Advanced** | cohorts/retention, SIM runway forecast, percentile KPIs (t-digest or daily-grain), CSV export, threshold alerts, quarterly/annual spans; platform instance (cross-org KPIs) | — |

## 11. Resolved Decisions (2026-07-07)

| # | Question | Decision |
|---|---|---|
| 1 | Deployment model | **Per-org deployment** (repo convention: `ORGNAME` env, per-org compose/helm). `org_id` stays on every table (cheap, keeps schema portable); platform cross-org KPIs run later in a separate Ukama-operated instance — out of scope for v1. |
| 2 | Window sizing | **`W = 5min`** (pipeline config default). Snapshots pulled 288×/day; `for_each` chains should complete within one window — monitored via chain-latency metric. *(Revised: the separate watermark `L` was dropped — windows are pulled at close; late data goes through the dirty-recompute path.)* |
| 3 | Historical backfill | **None — start fresh.** Trends read `new` for the first day/week/month. Admin `RepullWindow` remains for later selective backfill of windowed sources. |
| 4 | Snapshot storage | **Change-log + state-as-of.** A `full_snapshot` row is stored only when its content-hash changes (`dedup_key = hash(fields)`); each dataset also records "seen in window N" high-water marks so deletions are detectable. Analysis reads *state as of window N* = latest change ≤ N. Resolves the dedup contradiction with minimal storage. |
| 5 | Percentile KPIs | **Deferred to Phase 4** (with t-digest sketches or daily-grain compute — decided then). v1 ships AVG/MIN/MAX for latency/usage KPIs; BACKHAUL_LATENCY_P95 and USAGE_PER_SUBSCRIBER p95 removed from the v1 catalog. |
| 6 | Cross-window KPIs | **KPI spec `lookback` field** (`lookback: 30d` or N windows). Engine feeds prior raw windows and/or prior `kpi_windows` to the algo alongside the current window — still pure and deterministic; dirty-cascade recompute unchanged. Enables MTBF/MTTR, CHURN_RATE, RETENTION_COHORT, SIM_POOL_RUNWAY. |
| 7 | Pull auth | **BypassAuth for now** — analytics runs inside `ukama-net` and calls source gateways directly (matching existing compose `BYPASS_AUTH_MODE` usage). Revisit when a platform-wide service-credential mechanism exists. |
| 8 | Endpoint verification | *Remains a Phase 0 exit task:* verify each source spec's `strategy` against endpoint reality (billing reports are period-based; check notification/health time filters). Mismatches fall back to `full_snapshot` + change-log. |
| 9 | Multi-currency revenue | **Normalize to org currency.** Org currency declared in service config; ingest converts amounts at pull time, storing `amount_org_cents` + original `amount, currency, fx_rate, fx_at` for audit. Until a mixed-currency source actually appears, conversion is identity (sources already bill in org currency); an FX-rate source is only introduced when needed — its absence for a non-org currency is an ingest error, not a silent pass-through. |

## 12. Repository Layout

```
systems/analytics/
  schema/                    # shared module: window model, models, migrations, spec schemas, KPI registry types
  ingest/                    # cmd, pb, pkg/engine (grid, pullers, event consumers), pkg/server, configs/sources/*.yaml
  analysis/                  # cmd, pb, pkg/algos (registry), pkg/server, configs/kpis/*.yaml
  aggregator/                # cmd, pb (read API), pkg/rollup, pkg/server
  api-gateway/               # cmd, pkg/rest, pkg/client
  docs/                      # this plan, schema.md, runbook.md
  docker-compose.yml         # postgres (timescale image optional, §2), 3 services, gateway, msgclient-analytics
```

Each service: standard Ukama layout (cmd/server, pb/gen, pkg/config embedding `uconf.BaseConfig`, mocks, Dockerfile, Int.Dockerfile, Makefile with protoc+mockery `gen`, MPL-2.0 headers). Shared window/span config env-injected identically to all three services from compose/helm.
