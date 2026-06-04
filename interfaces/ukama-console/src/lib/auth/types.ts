/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Shared auth types and constants (no runtime side-effects, safe to import
 * from both server and client code). The signed token itself is never part
 * of AuthUser — it stays in an httpOnly cookie and is sent to the gateway
 * automatically; only non-sensitive claims are exposed to the UI.
 */

/** Identity claims surfaced to the application from the session token. */
export interface AuthUser {
  id: string;
  name: string;
  email: string;
  role: string;
  orgId: string;
  orgName: string;
  country: string;
  currency: string;
  isEmailVerified: boolean;
  isShowWelcome: boolean;
}

/** Cookie set by the auth app (Kratos) on successful login. */
export const SESSION_COOKIE = 'ukama_session';

/** httpOnly cookie holding the gateway-signed session token. */
export const TOKEN_COOKIE = 'token';

/** Internal request header carrying the (base64-JSON) user to RSC pages. */
export const USER_HEADER = 'x-ukama-user';

/** Token cookie lifetime (24h), matched by the gateway token semantics. */
export const TOKEN_MAX_AGE_SECONDS = 60 * 60 * 24;
