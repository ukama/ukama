/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/** Create-plan form schema (BUILD-PLAN §13.4 — zod schema colocated). */
import { z } from 'zod';

export const createPlanSchema = z
  .object({
    name: z.string().min(1, 'Plan name is required'),
    price: z.coerce.number({ message: 'Enter a price' }).positive('Must be > 0'),
    data: z.coerce.number({ message: 'Enter a volume' }).positive('Must be > 0'),
    unit: z.enum(['GB', 'MB']),
    // Stored as days; the form offers preset durations (daily/weekly/monthly).
    days: z.coerce.number().int().refine((d) => [1, 7, 30].includes(d), {
      message: 'Select a validity period',
    }),
    // Org-wide plans carry no network; otherwise a specific network is required.
    availableWithinOrg: z.boolean(),
    networkId: z.string().optional(),
  })
  .refine((v) => v.availableWithinOrg || !!v.networkId, {
    message: 'Select a network',
    path: ['networkId'],
  });

export type CreatePlanValues = z.infer<typeof createPlanSchema>;

/** Validity presets — value is the day count stored on the plan. */
export const VALIDITY_OPTIONS = [
  { value: '1', label: 'Daily (1 day)' },
  { value: '7', label: 'Weekly (7 days)' },
  { value: '30', label: 'Monthly (30 days)' },
] as const;
