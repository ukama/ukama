/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/** Node shape (zod schema + inferred type). */
import { z } from 'zod';

export const NodeSchema = z.object({
  id: z.string(),
  serial: z.string(),
  /** Human-friendly node name from the registry; undefined in mocks. */
  name: z.string().optional(),
  /** Raw connectivity (Online/Offline/Unknown) for the status dot. */
  connectivity: z.string().optional(),
  /** Raw operational state (Operational/Configured/Faulty/Unknown). */
  state: z.string().optional(),
  type: z.enum(['Tower node', 'Amplifier node', 'Controller node', 'Home node']),
  site: z.string(),
  status: z.enum(['online', 'degraded', 'offline', 'configuring']),
  cpu: z.number(),
  mem: z.number(),
  temp: z.number().nullable(),
  fw: z.string(),
  up: z.string(),
  note: z.string().optional(),
});

export type UkamaNode = z.infer<typeof NodeSchema>;
