# Console-BFF

GraphQL backend-for-frontend for the Ukama Console. It authenticates console
sessions (Kratos) and proxies/aggregates the distributed backend systems
(nucleus, registry, dataplan, subscriber, node, …) behind a single `/graphql`
endpoint.

It runs as **two processes**: the consolidated API server (`server/`, port
8080) and the real-time subscriptions service (`subscriptions/`, port 8081).

## How to run

Create a `.env` in the repo root (copy `.env.local.example` for host-run dev —
it uses the host-mapped backend ports; `.env.dev.example` uses docker-network
hostnames and only works inside the compose network).

```bash
# API server (one process, all modules) — port 8090 by default in dev
API_PORT=8090 ENABLE_INTROSPECTION=true pnpm dev

# Subscriptions service (real-time metrics/notifications)
pnpm subscriptions-dev
```

Docker (two containers from one image — `api`:8080 and `subscriptions`:8081):

```bash
docker compose up -d --build --remove-orphans
```

## Architecture (post-consolidation, 2026-06)

The BFF historically ran as **23 processes**: an Apollo Federation gateway
composing 22 per-module subgraph servers (ports 5042–5063) under supervisord.
The modules shared **no federation entity references**, so the federation was
pure namespacing — it has been collapsed into **one Apollo server**. Design +
decision log: `CONSOLIDATION-DESIGN.md`.

- `server/index.ts` — the API server: shared express middleware (request-id,
  AsyncLocalStorage logging, helmet, rate-limit), one merged type-graphql
  schema (no federation), `/healthz` `/readyz` `/ping` `/get-user`
  `/set-theme`, graceful shutdown, introspection gated by
  `ENABLE_INTROSPECTION` (off in production by default).
- `server/context.ts` — unified `AppContext`: verified session-token claims
  (`headers.orgId/userId/orgName`) + **named datasource slots**
  (`ctx.dataSources.org`, `.network`, …), instantiated per request.
- `server/schema.ts` — composes every module's resolvers into one schema.
- Each module keeps its `resolver(s)/`, `datasource/` and `context/`; the old
  per-module bootstraps, the federation gateway, supervisord and the
  `@apollo/gateway`/`@apollo/subgraph` deps have been removed.

## Planning-tool: parked

`planning-tool/` (Prisma/PostgreSQL) is excluded from the schema **and from
build/lint** (`tsconfig.json` include list, `.eslintignore`) — its Prisma
client requires a configured `PLANNING_TOOL_DB`. The code is kept and its
schema names are already de-conflicted (`PlanningSite`, `*DraftSite`).

**Re-enabling planning-tool:**
1. Configure `PLANNING_TOOL_DB`; run `pnpm prisma-pre-build` (use
   `prisma migrate deploy` in production).
2. Add `"planning-tool"` back to the `include` list in `tsconfig.json` and
   remove it from `.eslintignore`. The standalone `planning-tool/index.ts`
   bootstrap is dead (it predates consolidation and imports the removed
   `@apollo/subgraph`) — delete it or leave it excluded.
3. `server/schema.ts`: import `planningResolvers from "../planning-tool/modules"`
   and append `...planningResolvers` to `ALL_RESOLVERS`.
4. `server/context.ts`: add `prisma` to `AppContext` and pass it from
   `server/index.ts` (singleton exported from `common/prisma`). Consider
   lazy-initializing the client so a missing DB degrades planning queries
   instead of failing boot.

## Auth model

- `/get-user` validates the Kratos `ukama_session`, resolves the user's
  org/role, and issues an HMAC-signed token (`JWT_SECRET`) with an `exp`
  claim. Every `/graphql` request must carry the session cookie + token;
  signature and expiry are verified at the server entry, claims are decoded
  into `ctx.headers`. Introspection-only queries may pass without a session
  when `ENABLE_INTROSPECTION=true` (schema only, never data).

## Useful endpoints

| Path | Purpose |
|------|---------|
| `/graphql` | the API |
| `/healthz` / `/readyz` | liveness / readiness probes |
| `/get-user` | session → signed token exchange |
| `/ping` | legacy alias (also checks subscriptions service) |

## Testing & CI

- `pnpm test:unit` — dependency-free unit suite (`common/tests`: token,
  storage TTL, env validation, request context) with a coverage gate
  (`jest.unit.config.ts`, ≥60% on the security-critical helpers). Runs in CI
  (`.github/workflows/bff.yaml`) and needs no backend systems.
- `pnpm test` — the integration suite (`gateway/tests`); requires live backend
  systems + a `TOKEN`, so it's run locally/in an integrated env, not plain CI.
- `pnpm audit:ci` — fails CI on **critical** dependency advisories.
- `load/load-test.js` — k6 load/SLO test. Run:
  `BASE_URL=http://localhost:8080 UKAMA_SESSION=<c> TOKEN=<c> k6 run load/load-test.js`
  (health probes run without creds; the GraphQL query is skipped if creds are
  absent). Thresholds are starting SLOs — tune to observed baselines.

## Dependencies & TypeScript notes

- `type-graphql` remains on `2.0.0-rc.2` — upgrade to a stable 2.x once
  published/verified (`pnpm add type-graphql@<version>` + `pnpm build` +
  `pnpm test:unit`).
- `skipLibCheck` is on: newer transitive deps ship `.d.ts` files targeting
  newer TS lib types; our own code is still fully checked.
- `strictPropertyInitialization` is intentionally **off**: type-graphql DTOs
  declare `@Field() x: string;` without initializers by design.
