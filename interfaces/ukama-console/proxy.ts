/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Next 16 request interceptor (replaces middleware.ts; Node runtime).
 * Phase-1 skeleton: pass-through. Phase 2 adds (BUILD-PLAN §13.6, §15.5):
 *  - ukama_session cookie → jose JWT verify (whoami fallback) + token cache
 *  - httpOnly token cookie + role/org headers
 *  - lens↔role route gating (/business, /network, /customer, /*\/manage)
 *  - per-request nonce + strict CSP (script-src/style-src) in production
 */
import { NextResponse } from 'next/server';

export default function proxy() {
  return NextResponse.next();
}

export const config = {
  matcher: ['/((?!_next/static|_next/image|favicon.ico|api/health|ping).*)'],
};
