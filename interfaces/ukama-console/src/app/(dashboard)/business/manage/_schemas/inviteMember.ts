/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/** Invite-member form schema (colocated per BUILD-PLAN §6). */
import { z } from 'zod';

export const inviteMemberSchema = z.object({
  email: z.string().email('Enter a valid email'),
  role: z.enum(['Owner', 'Administrator', 'Vendor', 'Network owner'], {
    message: 'Select a role',
  }),
});

export type InviteMemberValues = z.infer<typeof inviteMemberSchema>;
