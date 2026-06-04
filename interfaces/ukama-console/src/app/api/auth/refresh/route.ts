/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Clears the cached token cookie and bounces back to the app. Used to
 * self-heal when the gateway rejects a stale token (UNAUTHENTICATED): the
 * session cookie is still valid, so proxy.ts mints a fresh token on the
 * redirect. If the session is also gone, proxy.ts redirects to login.
 */
import { NextResponse } from 'next/server';
import { TOKEN_COOKIE } from '@/lib/auth/types';

export const runtime = 'nodejs';

export function GET(request: Request) {
  const redirectTo = new URL('/', request.url);
  const res = NextResponse.redirect(redirectTo);
  res.cookies.delete(TOKEN_COOKIE);
  return res;
}
