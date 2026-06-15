/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/** Auth role code → display label. */
const ROLE_LABELS: Record<string, string> = {
  ROLE_OWNER: 'Owner',
  ROLE_ADMIN: 'Admin',
  ROLE_NETWORK_OWNER: 'Network owner',
  ROLE_VENDOR: 'Vendor',
  ROLE_USER: 'Member',
};

/** Human label for an auth role code; humanizes unknown codes, '' when absent. */
export const roleLabel = (role?: string): string => {
  if (!role) return '';
  if (ROLE_LABELS[role]) return ROLE_LABELS[role];
  const c = role.replace(/^ROLE_/, '').replace(/_/g, ' ').toLowerCase();
  return c ? c[0]!.toUpperCase() + c.slice(1) : role;
};
