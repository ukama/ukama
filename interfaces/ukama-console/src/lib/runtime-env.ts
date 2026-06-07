/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Runtime public config.
 *
 * Next inlines `process.env.NEXT_PUBLIC_*` (dot access) into the bundle at
 * BUILD time, so those values are frozen in the image and can't be changed
 * by the deployment (Helm configmap). To make a single image configurable
 * at deploy time we instead:
 *   - read the container env at REQUEST time on the server via *bracket*
 *     access (`process.env['NEXT_PUBLIC_X']`), which Next does NOT inline; and
 *   - inject those values into the page (window.__UK_ENV__) so the browser
 *     reads the deployed values too.
 *
 * Build-time `@/env` defaults remain the fallback for local dev.
 */

export interface PublicEnv {
  /** GraphQL API gateway (console-bff API). */
  apiGw: string;
  /** Metrics / subscriptions endpoint (console-bff subscriptions). */
  metricUrl: string;
  /** Auth/session gateway (get-user, token mint). */
  apiGw4ss: string;
  /** Auth app (login redirect target). */
  authAppUrl: string;
}

/** Local-dev fallbacks (mirror src/env.ts defaults). */
const DEFAULTS: PublicEnv = {
  apiGw: 'http://localhost:8080',
  metricUrl: 'http://localhost:8081',
  apiGw4ss: 'http://localhost:8080',
  authAppUrl: 'http://localhost:4455',
};

/** Global key the server injects and the client reads. */
export const RUNTIME_ENV_KEY = '__UK_ENV__';

/**
 * SERVER ONLY. Reads the live container env at request time. Bracket access
 * is deliberate — it prevents Next from inlining the value at build, so the
 * Helm-provided env wins. Empty/unset falls back to the dev defaults.
 */
export function readServerEnv(): PublicEnv {
  const pick = (key: string, fallback: string): string => {
    const v = process.env[key];
    return v && v.length > 0 ? v : fallback;
  };
  return {
    apiGw: pick('NEXT_PUBLIC_API_GW', DEFAULTS.apiGw),
    metricUrl: pick('NEXT_PUBLIC_METRIC_URL', DEFAULTS.metricUrl),
    apiGw4ss: pick('NEXT_PUBLIC_API_GW_4SS', DEFAULTS.apiGw4ss),
    authAppUrl: pick('NEXT_PUBLIC_AUTH_APP_URL', DEFAULTS.authAppUrl),
  };
}

/**
 * Client-safe accessor. Reads the values the server injected on the page; on
 * the server (or before injection) it reads the live env directly.
 */
export function publicEnv(): PublicEnv {
  if (typeof window !== 'undefined') {
    const injected = (window as unknown as Record<string, PublicEnv>)[
      RUNTIME_ENV_KEY
    ];
    return injected ?? DEFAULTS;
  }
  return readServerEnv();
}

/** Serialized injection script body (safe against </script> breakout). */
export function runtimeEnvScript(): string {
  const json = JSON.stringify(readServerEnv()).replace(/</g, '\\u003c');
  return `window.${RUNTIME_ENV_KEY}=${json};`;
}
