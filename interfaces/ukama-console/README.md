# Ukama Console (reborn)

Greenfield rebuild of the Ukama operator console — one app, three lenses
(Business / Network / Customer). See `ukama-console-BUILD-PLAN.md` for the full
build plan and decisions.

## Run (development)

```bash
pnpm install
pnpm dev:review
```

`dev:review` runs three processes concurrently:

| Process | Output |
|---------|--------|
| `next dev` | http://localhost:3000 + `.logs/dev.log` |
| `tsc --noEmit --watch` | `.logs/tsc.log` (full-project type errors) |
| `scripts/lint-watch.mjs` | `.logs/lint.log` (ESLint, re-run on file change) |

The `.logs/*` files exist so the AI pair (Claude) can read type/lint errors
without terminal access. Plain `pnpm dev` works too.

## Scripts

- `pnpm typecheck` — one-shot full type-check
- `pnpm lint` — ESLint
- `pnpm prettier` — format
- `pnpm check:headers` — verify the MPL license header on every source file

## Conventions

- Every `.ts/.tsx/.js/.jsx/.mjs` file starts with the MPL-2.0 header.
- Colocation: route-specific code lives in `_components`/`_schemas` next to its
  route; shared code is promoted only when a second route needs it.
- Pages stay thin — business logic lives in `src/features/*`.
- Static/seed data lives only in `src/data/`.
