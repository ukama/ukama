/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Zod-validated environment (BUILD-PLAN §14). Server-only secrets must NEVER
 * use the NEXT_PUBLIC_ prefix. Endpoint names mirror the legacy console so existing deployments carry over.
 */
import { z } from 'zod';

const envSchema = z.object({
  NODE_ENV: z
    .enum(['development', 'test', 'production'])
    .default('development'),
  /** API gateway (GraphQL) */
  NEXT_PUBLIC_API_GW: z.string().url().default('http://localhost:8080'),
  /** Metrics endpoint (GraphQL + SSE) */
  NEXT_PUBLIC_METRIC_URL: z.string().url().default('http://localhost:8081'),
  /** Auth/session gateway (get-user, login redirects) */
  NEXT_PUBLIC_API_GW_4SS: z.string().url().default('http://localhost:8080'),
  /** Auth app (login redirect target) */
  NEXT_PUBLIC_AUTH_APP_URL: z.string().url().default('http://localhost:4455'),
});

/** Docker/CI pass unset build args through as empty strings — zod
 *  defaults only apply to `undefined`, so blank means "use default". */
const blank = (v: string | undefined) => (v === '' ? undefined : v);

export const env = envSchema.parse({
  NODE_ENV: process.env.NODE_ENV,
  NEXT_PUBLIC_API_GW: blank(process.env.NEXT_PUBLIC_API_GW),
  NEXT_PUBLIC_METRIC_URL: blank(process.env.NEXT_PUBLIC_METRIC_URL),
  NEXT_PUBLIC_API_GW_4SS: blank(process.env.NEXT_PUBLIC_API_GW_4SS),
  NEXT_PUBLIC_AUTH_APP_URL: blank(process.env.NEXT_PUBLIC_AUTH_APP_URL),
});
