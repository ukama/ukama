/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/** Data plans — ported from the design prototype (`data.jsx`). */
import { z } from 'zod';

export const PlanSchema = z.object({
  id: z.string(),
  name: z.string(),
  price: z.number(),
  data: z.string(),
  days: z.number(),
  subs: z.number(),
  color: z.string(),
});

export type Plan = z.infer<typeof PlanSchema>;

export const PLANS: Plan[] = [
  { id: 'p1', name: 'Starter', price: 5, data: '5 GB', days: 30, subs: 248, color: 'var(--uk-ink-3)' },
  { id: 'p2', name: 'Standard', price: 12, data: '20 GB', days: 30, subs: 712, color: 'var(--uk-ac)' },
  { id: 'p3', name: 'Unlimited', price: 25, data: 'Unlimited', days: 30, subs: 296, color: 'var(--uk-secondary)' },
  { id: 'p4', name: 'Day pass', price: 1, data: '1 GB', days: 1, subs: 88, color: 'var(--uk-orange)' },
];
