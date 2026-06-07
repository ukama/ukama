/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Acknowledges the first-visit welcome page. Forwards the request (with the
 * signed token cookie) to the gateway's /welcome-seen, which records the
 * acknowledgment per user. On success the cached token cookie is cleared so
 * the next navigation re-mints a token with isShowWelcome=false.
 */
import { NextResponse } from 'next/server';
import { readServerEnv } from '@/lib/runtime-env';
import { TOKEN_COOKIE } from '@/lib/auth/types';

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
  response.cookies.delete(TOKEN_COOKIE);
  return response;
}
