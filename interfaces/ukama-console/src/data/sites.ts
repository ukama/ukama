/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/** Site shape (zod schema + inferred type). */
import { z } from 'zod';

export const SiteSchema = z.object({
  id: z.string(),
  name: z.string(),
  area: z.string(),
  status: z.enum(['online', 'degraded', 'offline']),
  subs: z.number(),
  nodes: z.number(),
  uptime: z.number(),
  battery: z.number(),
  signal: z.number().nullable(),
  data: z.string(),
  /** dummy coordinates within the operating region (real geo at API phase) */
  lat: z.number(),
  lng: z.number(),
  plan: z.string(),
  issue: z.string().optional(),
});

export type Site = z.infer<typeof SiteSchema>;
