/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/** Subscribers / customers — ported from the design prototype (`data.jsx`). */
import { z } from 'zod';

export const SubscriberSchema = z.object({
  id: z.string(),
  name: z.string(),
  email: z.string().optional(),
  phone: z.string(),
  site: z.string(),
  plan: z.string(),
  usage: z.number(),
  cap: z.number().nullable(),
  sim: z.enum(['active', 'inactive', 'suspended']),
  iccid: z.string(),
  /** SIM record id (for SIM-scoped actions like top-up); undefined in mocks. */
  simId: z.string().optional(),
  seen: z.string(),
});

export type Subscriber = z.infer<typeof SubscriberSchema>;
