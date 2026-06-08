/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Server side of logout: clears the httpOnly token cookie. The ukama_session
 * cookie lives on the auth app's domain (Kratos) and is cleared there; the
 * client redirects to the auth app after calling this. Browser-stored UI
 * state is cleared client-side (see lib/auth/client logout()).
 */
import { NextResponse } from 'next/server';
import { TOKEN_COOKIE } from '@/lib/auth/types';
import { readServerEnv } from '@/lib/runtime-env';
import { cookieDomain, publicHost } from '@/lib/request-url';

export const runtime = 'nodejs';

export function POST(request: Request) {
  const res = NextResponse.json(
    { ok: true, redirect: `${readServerEnv().authAppUrl}/user/logout` },
    { headers: { 'cache-control': 'no-store' } },
  );
  // Clear with the same Domain it was set with so it actually goes away.
  res.cookies.delete({
    name: TOKEN_COOKIE,
    path: '/',
    domain: cookieDomain(publicHost(request)),
  });
  return res;
}
