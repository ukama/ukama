/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/** Organization members — ported from the design prototype (`data.jsx`). */
import { z } from 'zod';

export const MemberSchema = z.object({
  id: z.string(),
  name: z.string(),
  email: z.string(),
  role: z.enum(['Owner', 'Admin', 'Network owner', 'Vendor']),
  status: z.enum(['active', 'pending']),
  last: z.string(),
});

export type Member = z.infer<typeof MemberSchema>;

export const ROLE_DESC: Record<Member['role'], string> = {
  Owner: 'Full access · billing',
  Admin: 'Manage network & members',
  'Network owner': 'Manages network operation',
  Vendor: 'Install & maintain hardware',
};
