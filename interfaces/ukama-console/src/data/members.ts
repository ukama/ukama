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

export const MEMBERS: Member[] = [
  { id: 'm1', name: 'Joseph Mulenga', email: 'joseph@kwacha.co', role: 'Owner', status: 'active', last: 'Active now' },
  { id: 'm2', name: 'Grace Tembo', email: 'grace@kwacha.co', role: 'Admin', status: 'active', last: '2h ago' },
  { id: 'm3', name: 'Daniel Phiri', email: 'daniel@kwacha.co', role: 'Network owner', status: 'active', last: 'Yesterday' },
  { id: 'm4', name: 'Ruth Mwanza', email: 'ruth@kwacha.co', role: 'Vendor', status: 'active', last: '3d ago' },
  { id: 'm5', name: 'Peter Banda', email: 'peter@kwacha.co', role: 'Admin', status: 'pending', last: 'Invited 1d ago' },
];

export const ROLE_DESC: Record<Member['role'], string> = {
  Owner: 'Full access · billing',
  Admin: 'Manage network & members',
  'Network owner': 'Operate a single network',
  Vendor: 'Install & maintain hardware',
};
