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
