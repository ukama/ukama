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
import { TOKEN_COOKIE } from '@/lib/auth/types';
import { cookieDomain, publicHost, publicUrl } from '@/lib/request-url';
import { NextResponse } from 'next/server';

export const runtime = 'nodejs';

export function GET(request: Request) {
  const res = NextResponse.redirect(publicUrl(request, '/'));
  // Delete with the same Domain the cookie was set with, or it won't clear.
  res.cookies.delete({
    name: TOKEN_COOKIE,
    path: '/',
    domain: cookieDomain(publicHost(request)),
  });
  return res;
}
