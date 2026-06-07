/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Resolves the public origin for same-origin redirects.
 *
 * Behind an ingress/proxy the pod's own `Host` can be its listen address
 * (e.g. 0.0.0.0:8080), so `request.url` is unsafe for building redirect
 * targets — it would send the browser to the internal address. The proxy
 * sets X-Forwarded-Host/Proto to the real public host; prefer those, then
 * the Host header, and only fall back to the request URL's origin.
 */
export function publicOrigin(req: { headers: Headers; url: string }): string {
  const first = (v: string | null): string =>
    (v ?? '').split(',')[0]?.trim() ?? '';

  const host = first(req.headers.get('x-forwarded-host')) || first(req.headers.get('host'));
  if (!host) {
    try {
      return new URL(req.url).origin;
    } catch {
      return '';
    }
  }
  let proto = first(req.headers.get('x-forwarded-proto'));
  if (!proto) {
    try {
      proto = new URL(req.url).protocol.replace(':', '');
    } catch {
      proto = 'https';
    }
  }
  return `${proto}://${host}`;
}

/** Builds an absolute same-origin URL safe to use as a redirect target. */
export function publicUrl(
  req: { headers: Headers; url: string },
  pathname: string,
): URL {
  return new URL(pathname, publicOrigin(req));
}
