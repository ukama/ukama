/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Next 16 request interceptor (replaces middleware.ts; Node runtime).
 *
 * Standard reactive auth gate, evaluated per navigation:
 *  - no ukama_session         → logout: clear token cookie, go to auth app
 *  - session, no usable token → mint a signed token via /get-user and cache
 *                               it in an httpOnly cookie; if the session is
 *                               invalid/unreachable, redirect to the auth app
 *  - session + token          → forward decoded user to RSC; refresh cookie
 *
 * Stale tokens are handled lazily: the gateway returns 401, the Apollo error
 * link bounces through /api/auth/refresh, and the next navigation re-mints
 * from the still-valid session. No background polling.
 */
import { env } from '@/env';
import { readServerEnv } from '@/lib/runtime-env';
import { publicUrl } from '@/lib/request-url';
import { decodeUserFromToken, fetchSession } from '@/lib/auth/token';
import {
  SESSION_COOKIE,
  TOKEN_COOKIE,
  TOKEN_MAX_AGE_SECONDS,
  USER_HEADER,
  type AuthUser,
} from '@/lib/auth/types';
import { NextResponse, type NextRequest } from 'next/server';

const encodeUser = (user: AuthUser): string =>
  Buffer.from(JSON.stringify(user), 'utf8').toString('base64');

/** Clears the token cookie and sends the user to the auth app (logout). */
const logoutRedirect = (): NextResponse => {
  const res = NextResponse.redirect(
    new URL('/auth/login', readServerEnv().authAppUrl),
  );
  res.cookies.delete(TOKEN_COOKIE);
  return res;
};

export default async function proxy(
  request: NextRequest,
): Promise<NextResponse> {
  const { pathname } = request.nextUrl;
  // Auth API routes must run even when the token cookie is missing/stale;
  // /unauthorized must render without a resolvable user (it's the landing
  // page for exactly that case).
  if (pathname.startsWith('/api/auth') || pathname === '/unauthorized') {
    return NextResponse.next();
  }

  // No session → logout flow.
  const session = request.cookies.get(SESSION_COOKIE)?.value;
  if (!session) return logoutRedirect();

  // Resolve the user: reuse the cached token, else mint one from the session.
  const cachedToken = request.cookies.get(TOKEN_COOKIE)?.value ?? '';
  let user = cachedToken ? decodeUserFromToken(cachedToken) : null;
  let freshToken: string | null = null;

  if (!user) {
    const result = await fetchSession(request.headers.get('cookie') ?? '');
    // Session exists but doesn't resolve to a complete console identity
    // (no ukama user / no org / no role / incomplete claims — the BFF
    // returned 401 with the failing step). Don't bounce to the auth app
    // (the Kratos session is valid — that would loop); land on
    // /unauthorized, where the user can only log out or contact support.
    if (!result) {
      const res = NextResponse.redirect(publicUrl(request, '/unauthorized'));
      res.cookies.delete(TOKEN_COOKIE);
      return res;
    }
    user = result.user;
    freshToken = result.token;
  }

  // Write the Set-Cookie header raw: response.cookies.set() percent-encodes
  // the value (+ / = in the base64 payload), and the gateway verifies the
  // HMAC over the exact payload bytes — an encoded token fails with 401.
  // Raw is safe here: base64 + base64url contain no ';' or whitespace.
  const attachFreshToken = (res: NextResponse): NextResponse => {
    if (!freshToken) return res;
    const parts = [
      `${TOKEN_COOKIE}=${freshToken}`,
      'Path=/',
      'HttpOnly',
      'SameSite=Lax',
      `Max-Age=${TOKEN_MAX_AGE_SECONDS}`,
    ];
    if (env.NODE_ENV === 'production') parts.push('Secure');
    res.headers.append('set-cookie', parts.join('; '));
    return res;
  };

  // First-visit welcome gate (page navigations only): eligible users land on
  // /welcome until they acknowledge it; everyone else is kept out of it.
  if (!pathname.startsWith('/api')) {
    const onWelcome = pathname === '/welcome';
    if (user.isShowWelcome && !onWelcome) {
      return attachFreshToken(
        NextResponse.redirect(publicUrl(request, '/welcome')),
      );
    }
    if (!user.isShowWelcome && onWelcome) {
      return attachFreshToken(NextResponse.redirect(publicUrl(request, '/')));
    }
  }

  // Forward the user to server components (getCurrentUser) via a header.
  const requestHeaders = new Headers(request.headers);
  requestHeaders.set(USER_HEADER, encodeUser(user));

  const response = NextResponse.next({ request: { headers: requestHeaders } });
  return attachFreshToken(response);
}

export const config = {
  matcher: ['/((?!_next/static|_next/image|favicon.ico|api/health|ping).*)'],
};
