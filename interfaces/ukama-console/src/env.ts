/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Zod-validated environment (BUILD-PLAN §14). Server-only secrets must NEVER
 * use the NEXT_PUBLIC_ prefix. Grows in Phase 2 (API_GW, auth, metrics).
 */
import { z } from 'zod';

const envSchema = z.object({
  NODE_ENV: z
    .enum(['development', 'test', 'production'])
    .default('development'),
});

export const env = envSchema.parse({
  NODE_ENV: process.env.NODE_ENV,
});
