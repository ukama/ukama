# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

FROM node:24-bullseye-slim AS base

# pnpm pinned to the project's local version, with retry options
RUN npm config set registry https://registry.npmjs.org/ && \
    npm config set fetch-retry-mintimeout 20000 && \
    npm config set fetch-retry-maxtimeout 120000 && \
    npm config set fetch-timeout 300000 && \
    npm install -g pnpm@10.11.0

FROM base AS builder

WORKDIR /app

# Install dependencies first (cached while the lockfile is unchanged)
COPY package.json pnpm-lock.yaml ./

RUN pnpm config set fetch-retry-mintimeout 20000 && \
    pnpm config set fetch-retry-maxtimeout 120000 && \
    pnpm config set fetch-timeout 300000 && \
    pnpm install --frozen-lockfile --ignore-scripts

# Copy source files (src/proxy.ts is the Next 16 request interceptor)
COPY src ./src
COPY public ./public
COPY next.config.mjs tsconfig.json ./

# Build-time public env (NEXT_PUBLIC_* is inlined into the client bundle).
# Only the vars src/env.ts validates — extend both together.
ARG NEXT_PUBLIC_API_GW
ENV NEXT_PUBLIC_API_GW=${NEXT_PUBLIC_API_GW}
ARG NEXT_PUBLIC_METRIC_URL
ENV NEXT_PUBLIC_METRIC_URL=${NEXT_PUBLIC_METRIC_URL}
ARG NEXT_PUBLIC_API_GW_4SS
ENV NEXT_PUBLIC_API_GW_4SS=${NEXT_PUBLIC_API_GW_4SS}
ARG NEXT_PUBLIC_AUTH_APP_URL
ENV NEXT_PUBLIC_AUTH_APP_URL=${NEXT_PUBLIC_AUTH_APP_URL}
ARG NEXT_PUBLIC_SIM_TYPE
ENV NEXT_PUBLIC_SIM_TYPE=${NEXT_PUBLIC_SIM_TYPE}

RUN pnpm run build

FROM node:24-bullseye-slim AS runner

WORKDIR /app

ENV NODE_ENV=production
RUN addgroup --system --gid 1001 nodejs && \
    adduser --system --uid 1001 nextjs
USER nextjs

# `output: 'standalone'` (next.config.mjs): self-contained server, no pnpm needed
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static
COPY --from=builder --chown=nextjs:nodejs /app/public ./public

# Same runtime env (read server-side by src/env.ts)
ARG NEXT_PUBLIC_API_GW
ENV NEXT_PUBLIC_API_GW=${NEXT_PUBLIC_API_GW}
ARG NEXT_PUBLIC_METRIC_URL
ENV NEXT_PUBLIC_METRIC_URL=${NEXT_PUBLIC_METRIC_URL}
ARG NEXT_PUBLIC_API_GW_4SS
ENV NEXT_PUBLIC_API_GW_4SS=${NEXT_PUBLIC_API_GW_4SS}
ARG NEXT_PUBLIC_AUTH_APP_URL
ENV NEXT_PUBLIC_AUTH_APP_URL=${NEXT_PUBLIC_AUTH_APP_URL}
ARG NEXT_PUBLIC_SIM_TYPE
ENV NEXT_PUBLIC_SIM_TYPE=${NEXT_PUBLIC_SIM_TYPE}

# Parity with the legacy console deployment (it serves on 8080)
ENV PORT=8080
ENV HOSTNAME=0.0.0.0
EXPOSE 8080

CMD ["node", "server.js"]
