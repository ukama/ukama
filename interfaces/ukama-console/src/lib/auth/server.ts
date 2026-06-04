/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Server-side accessor for the authenticated user. proxy.ts attaches the
 * user (base64-encoded JSON) to a request header; any server component or
 * route handler can read it here without re-fetching.
 */
import { headers } from 'next/headers';
import { USER_HEADER, type AuthUser } from './types';

/** Returns the current user from the request header, or null if absent. */
export async function getCurrentUser(): Promise<AuthUser | null> {
  const raw = (await headers()).get(USER_HEADER);
  if (!raw) return null;
  try {
    return JSON.parse(Buffer.from(raw, 'base64').toString('utf8')) as AuthUser;
  } catch {
    return null;
  }
}
