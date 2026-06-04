/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Networks — sample data ported from the design prototype (`data.jsx`).
 * Pattern: Zod schema → z.infer type → typed dataset (BUILD-PLAN §10/§13.5).
 */
import { z } from 'zod';

export const NetworkSchema = z.object({
  id: z.string(),
  name: z.string(),
  region: z.string(),
  status: z.enum(['online', 'degraded', 'offline']),
});

export type Network = z.infer<typeof NetworkSchema>;

export const NETWORKS: Network[] = [
  {
    id: 'kwacha',
    name: 'Kwacha Mobile',
    region: 'Zambia · Lusaka Province',
    status: 'online',
  },
  {
    id: 'copperbelt',
    name: 'Copperbelt Rural',
    region: 'Zambia · Copperbelt',
    status: 'degraded',
  },
  { id: 'demo', name: 'Demo network', region: 'Sandbox', status: 'online' },
  { id: 'maiko', name: 'Maiko', region: 'DRC · Kinshasa', status: 'online' },
];
