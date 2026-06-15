/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/** Attention feed / notifications — ported from the design prototype (`data.jsx`). */
import { z } from 'zod';

export const AlertSchema = z.object({
  id: z.string(),
  sev: z.enum(['critical', 'warning', 'info']),
  icon: z.string(),
  title: z.string(),
  detail: z.string(),
  site: z.string().optional(),
  action: z.string(),
  age: z.string(),
});

export type Alert = z.infer<typeof AlertSchema>;
