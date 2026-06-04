/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/** Create-plan form schema (BUILD-PLAN §13.4 — zod schema colocated). */
import { z } from 'zod';

export const createPlanSchema = z.object({
  name: z.string().min(1, 'Plan name is required'),
  price: z.coerce.number({ message: 'Enter a price' }).positive('Must be > 0'),
  data: z.coerce.number({ message: 'Enter a volume' }).positive('Must be > 0'),
  unit: z.enum(['GB', 'MB', 'Unlimited']),
  days: z.coerce
    .number({ message: 'Enter validity days' })
    .int()
    .positive('Must be > 0'),
});

export type CreatePlanValues = z.infer<typeof createPlanSchema>;
