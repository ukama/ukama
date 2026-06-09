/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Client-side auth helpers: browser-state cleanup and the logout flow.
 * Safe to import from client components only.
 */
import { publicEnv } from '@/lib/runtime-env';

/** localStorage keys owned by the app (cleared on logout). */
const PERSISTED_KEYS = ['uk-ui-prefs'];

/** Removes app-owned browser state. Apollo's in-memory cache dies on reload. */
export function clearClientData(): void {
  if (typeof window === 'undefined') return;
  try {
    for (const key of PERSISTED_KEYS) localStorage.removeItem(key);
    sessionStorage.clear();
  } catch {
    /* storage may be unavailable (private mode) — ignore */
  }
}

/**
 * Full logout: clear the token cookie + browser state, then hand off to the
 * auth app's logout route. `/user/logout` destroys the Kratos session
 * (ukama_session) — redirecting to /auth/login instead would leave the
 * session valid and bounce the user straight back into the console.
 */
export async function logout(): Promise<void> {
  if (typeof window === 'undefined') return;
  try {
    await fetch('/api/auth/logout', { method: 'POST' });
  } catch {
    /* best-effort — proceed with client cleanup regardless */
  }
  clearClientData();
  window.location.assign(`${publicEnv().authAppUrl}/user/logout`);
}
