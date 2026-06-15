/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/** SIM pool + node pool inventory — ported from the prototype (`data.jsx`). */
import { z } from 'zod';

export const SimBatchSchema = z.object({
  id: z.string(),
  batch: z.string(),
  type: z.enum(['Physical', 'eSIM']),
  qty: z.number(),
  assigned: z.number(),
  uploaded: z.string(),
});

export type SimBatch = z.infer<typeof SimBatchSchema>;

export const NodePoolItemSchema = z.object({
  id: z.string(),
  serial: z.string(),
  type: z.enum(['Tower node', 'Amplifier node']),
  status: z.enum(['available', 'assigned', 'rma']),
  site: z.string().optional(),
  added: z.string(),
});

export type NodePoolItem = z.infer<typeof NodePoolItemSchema>;
