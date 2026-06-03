/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/** SIM pool + node pool inventory — ported from the prototype (`data.jsx`). */
import { z } from 'zod';

export const SIMS_SUMMARY = {
  total: 2000,
  assigned: 1344,
  available: 611,
  suspended: 31,
  faulty: 14,
} as const;

export const SimBatchSchema = z.object({
  id: z.string(),
  batch: z.string(),
  type: z.enum(['Physical', 'eSIM']),
  qty: z.number(),
  assigned: z.number(),
  uploaded: z.string(),
});

export type SimBatch = z.infer<typeof SimBatchSchema>;

export const SIM_BATCHES: SimBatch[] = [
  { id: 'b1', batch: 'KW-2024-Q4-A', type: 'Physical', qty: 1000, assigned: 842, uploaded: '12 Oct 2025' },
  { id: 'b2', batch: 'KW-2024-Q4-B', type: 'eSIM', qty: 500, assigned: 312, uploaded: '12 Oct 2025' },
  { id: 'b3', batch: 'KW-2025-Q1-A', type: 'Physical', qty: 500, assigned: 190, uploaded: '04 Jan 2026' },
];

export const NodePoolItemSchema = z.object({
  id: z.string(),
  serial: z.string(),
  type: z.enum(['Tower node', 'Amplifier node']),
  status: z.enum(['available', 'assigned', 'rma']),
  site: z.string().optional(),
  added: z.string(),
});

export type NodePoolItem = z.infer<typeof NodePoolItemSchema>;

export const NODE_POOL: NodePoolItem[] = [
  { id: 'np1', serial: 'uk-tnode-a06-2010', type: 'Tower node', status: 'available', added: '04 Jan 2026' },
  { id: 'np2', serial: 'uk-tnode-a06-2011', type: 'Tower node', status: 'available', added: '04 Jan 2026' },
  { id: 'np3', serial: 'uk-anode-a06-2012', type: 'Amplifier node', status: 'assigned', site: 'Kabwe Central', added: '04 Jan 2026' },
  { id: 'np4', serial: 'uk-tnode-a06-2013', type: 'Tower node', status: 'rma', added: '18 Dec 2025' },
];
