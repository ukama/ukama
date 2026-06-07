/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

import { TOKEN_COOKIE } from '@/lib/auth/types';
import { cookieDomain, publicHost } from '@/lib/request-url';
import { readServerEnv } from '@/lib/runtime-env';
import { NextResponse } from 'next/server';

export const runtime = 'nodejs';

export async function POST(request: Request) {
  const cookieHeader = request.headers.get('cookie') ?? '';
  if (!cookieHeader) {
    return NextResponse.json({ ok: false }, { status: 401 });
  }

  const res = await fetch(`${readServerEnv().apiGw4ss}/welcome-seen`, {
    method: 'POST',
    cache: 'no-store',
    headers: { cookie: cookieHeader },
  }).catch(() => null);

  if (!res?.ok) {
    console.error(
      `[auth] /welcome-seen failed (${res ? res.status : 'unreachable'})`,
    );
    return NextResponse.json({ ok: false }, { status: 502 });
  }

  const response = NextResponse.json(
    { ok: true },
    { headers: { 'cache-control': 'no-store' } },
  );
  response.cookies.delete({
    name: TOKEN_COOKIE,
    path: '/',
    domain: cookieDomain(publicHost(request)),
  });
  return response;
}
