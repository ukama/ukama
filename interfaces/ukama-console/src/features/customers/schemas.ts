/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/** Add-customer form schema (BUILD-PLAN §13.4 — zod schema colocated). */
import { z } from 'zod';

export const addCustomerSchema = z.object({
  first: z.string().min(1, 'First name is required'),
  last: z.string().optional(),
  mobile: z
    .string()
    .min(7, 'Enter a valid mobile number')
    .regex(/^[+0-9 ()-]+$/, 'Digits, spaces and + only'),
  email: z.string().email('Enter a valid email').optional().or(z.literal('')),
  planId: z.string().optional(),
  sim: z.string().optional(),
});

export type AddCustomerValues = z.infer<typeof addCustomerSchema>;
