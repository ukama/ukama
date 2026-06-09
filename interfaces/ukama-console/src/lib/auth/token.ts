/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Server-only helpers for the gateway session token. Import only from
 * server contexts (proxy, route handlers, RSC) — never from client code.
 *
 * The gateway issues a token of the form `<base64-claims>.<hmac-signature>`.
 * The gateway is the security boundary and re-verifies the signature on
 * every request, so here we only decode the claims for UI/routing. We never
 * trust these claims for authorization decisions the gateway doesn't enforce.
 */
import { env } from '@/env';
import type { AuthUser } from './types';

/** Shape returned by the gateway `/get-user` endpoint. */
interface BffSessionResponse {
  orgId: string;
  orgName: string;
  role: string;
  name: string;
  email: string;
  userId: string;
  currency: string;
  country: string;
  token: string;
  isEmailVerified: boolean;
  isShowWelcome: boolean;
}

/**
 * Decodes the claims embedded in a session token without verifying its
 * signature (the gateway does that). Returns null for any malformed token.
 */
export function decodeUserFromToken(token: string): AuthUser | null {
  try {
    const payload = token.split('.')[0];
    if (!payload) return null;
    const decoded = Buffer.from(payload, 'base64').toString('utf8');
    const p = decoded.split(';');
    // orgId;orgName;userId;name;email;role;verified;welcome;country;currency;exp
    if (p.length < 10) return null;

    // Treat an expired token as absent so the proxy re-mints from the session
    // (the gateway also rejects expired tokens — this just avoids a 401 trip).
    const expRaw = p[10];
    if (expRaw) {
      const exp = Number.parseInt(expRaw, 10);
      if (!Number.isNaN(exp) && Math.floor(Date.now() / 1000) >= exp) {
        return null;
      }
    }

    return {
      orgId: p[0] ?? '',
      orgName: p[1] ?? '',
      id: p[2] ?? '',
      name: p[3] ?? '',
      email: p[4] ?? '',
      role: p[5] ?? '',
      isEmailVerified: (p[6] ?? '').includes('true'),
      isShowWelcome: (p[7] ?? '').includes('true'),
      country: p[8] ?? '',
      currency: p[9] ?? '',
    };
  } catch {
    return null;
  }
}

/**
 * Exchanges a session cookie for a freshly minted, signed token by calling
 * the gateway `/get-user`. Returns the token (to persist in a cookie) and
 * the decoded user, or null when the session is invalid/expired.
 */
export async function fetchSession(
  cookieHeader: string,
): Promise<{ user: AuthUser; token: string } | null> {
  try {
    const res = await fetch(`${env.NEXT_PUBLIC_API_GW_4SS}/get-user`, {
      method: 'GET',
      cache: 'no-store',
      headers: {
        cookie: cookieHeader,
        'content-type': 'application/json',
      },
    });
    if (!res.ok) return null;

    const data = (await res.json()) as BffSessionResponse;
    if (!data?.token) return null;

    return {
      token: data.token,
      user: {
        id: data.userId,
        name: data.name,
        email: data.email,
        role: data.role,
        orgId: data.orgId,
        orgName: data.orgName,
        country: data.country,
        currency: data.currency,
        isEmailVerified: data.isEmailVerified,
        isShowWelcome: data.isShowWelcome,
      },
    };
  } catch {
    return null;
  }
}
