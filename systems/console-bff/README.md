# Console-BFF

GraphQL backend-for-frontend for the Ukama Console. It authenticates console
sessions (Kratos), and proxies/aggregates the distributed backend systems
(nucleus, registry, dataplan, subscriber, node, …) behind a single `/graphql`
endpoint.

## How to run

Create a `.env` in the repo root (copy `.env.local.example` for host-run dev —
it uses the host-mapped backend ports; `.env.dev.example` uses docker-network
hostnames and only works inside the compose network).

```bash
# Consolidated API server (recommended) — one process, port 8090 by default
API_PORT=8090 ENABLE_INTROSPECTION=true yarn api-dev

# Legacy multi-process stack (federation gateway + 22 subgraphs)
yarn dev
```

Docker: `docker compose up -d --build bff` (runs the legacy stack until the
Phase C cutover — see below).

## Consolidation (2026-06) — what changed and what's left

The BFF historically ran as **23 processes**: an Apollo Federation gateway
(`gateway/`, port 8080) composing 22 per-module subgraph servers (ports
5042–5063) under supervisord. The modules share **no federation entity
references** (`@key`/`@external`), so the federation was pure namespacing.
It is being collapsed into **one Apollo server** (`server/`, `API_PORT`).
Full design + decision log: `CONSOLIDATION-DESIGN.md`.

### Current state

- `server/index.ts` — consolidated API server: shared express middleware
  (request-id, AsyncLocalStorage logging, helmet, rate-limit), one merged
  type-graphql schema (no federation), `/healthz` `/readyz` `/ping`
  `/get-user` `/set-theme`, graceful shutdown, introspection gated by
  `ENABLE_INTROSPECTION` (off in production by default).
- `server/context.ts` — unified `AppContext`: verified session token claims
  (`headers.orgId/userId/orgName`) + **named datasource slots**
  (`ctx.dataSources.org`, `.network`, …) instead of the old per-subgraph
  `ctx.dataSources.dataSource`.
- `server/schema.ts` — composes every module's resolvers into one schema.
- **All 22 modules are migrated.** Each module's `context/` and `index.ts`
  (subgraph bootstrap) were updated to the named key, so the legacy gateway
  stack still works during the transition.
- Naming fixes made for the merge (single schema = single namespace):
  `member` re-exports the `User*` types from `user`; subscriber's
  `SimsAPIResDto` → `SubscriberSimsAPIResDto`; planning-tool's `Site` →
  `PlanningSite` and `addSite/updateSite/deleteSite` → `addDraftSite/
  updateDraftSite/deleteDraftSite`.

### Planning-tool: parked (NOT in the merged schema)

`planning-tool/` (Prisma/PostgreSQL) is **intentionally excluded** from
`server/schema.ts` for phase 1 — instantiating its Prisma client crashes the
server when `PLANNING_TOOL_DB` isn't configured / `prisma generate` hasn't
run. The code is kept and its schema names are already collision-free.

**Re-enabling planning-tool** (when needed):
1. Configure `PLANNING_TOOL_DB` and run
   `yarn prisma-pre-build` (generates the Prisma client + pushes schema; move
   to `prisma migrate deploy` for production).
2. `server/schema.ts`: import `planningResolvers from "../planning-tool/modules"`
   and append `...planningResolvers` to `ALL_RESOLVERS`.
3. `server/context.ts`: add `prisma` to `AppContext` and pass it in the
   context built in `server/index.ts` (the singleton lives in
   `common/prisma`, exported as `prisma`).
4. Consider lazy-initializing the Prisma client so a missing DB degrades the
   planning queries instead of failing boot.

### Phase C — cutover (DONE, deployment)

- `Dockerfile` builds at image time (`yarn build`) and the default `CMD` runs
  `dist/server/index.js`; **supervisord removed**.
- `docker-compose.yml` runs **two containers from one image**: `api` (8080,
  consolidated server) and `subscriptions` (8081), shared env via a YAML
  anchor, `API_PORT=8080`.
- `package.json` scripts trimmed: `start` runs the consolidated server,
  `dev`→`api-dev`; removed the 22 `*-dev` scripts + `all-dev`/`all-start`.

Deploy: `docker compose up -d --build --remove-orphans` (the `--remove-orphans`
clears the old `bff` container).

### Phase C — remaining cleanup (run on host; build already ignores this code)

The image still compiles the legacy federation source (harmless — `CMD` runs
only `server/`). Remove it when convenient:

```bash
rm -f gateway/index.ts gateway/configureExpress.ts supervisord.conf
rm -rf common/apollo
# the 22 subgraph bootstraps (KEEP each module's context/ resolver(s)/ datasource/):
for m in org user init network site member invitation node package rate \
  subscriber sim controller health software component billing payment \
  report metric notification; do rm -f "$m/index.ts"; done
```

Keep `gateway/tests/`, `subscriptions/`, and `planning-tool/` (parked).
Optionally prune the 5042–5063 `*_PORT` consts in `common/configs`. Then
regenerate `ukama-console` codegen against the consolidated endpoint.

## Auth model (post-hardening)

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
