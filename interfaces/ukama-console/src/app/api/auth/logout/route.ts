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
import { cookies } from 'next/headers';
import { TOKEN_COOKIE } from '@/lib/auth/types';
import { env } from '@/env';

export const runtime = 'nodejs';

export async function POST() {
  (await cookies()).delete(TOKEN_COOKIE);
  return NextResponse.json(
    { ok: true, redirect: `${env.NEXT_PUBLIC_AUTH_APP_URL}/auth/login` },
    { headers: { 'cache-control': 'no-store' } },
  );
}
