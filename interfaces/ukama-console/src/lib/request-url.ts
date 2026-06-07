/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
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

export function publicUrl(
  req: { headers: Headers; url: string },
  pathname: string,
): URL {
  return new URL(pathname, publicOrigin(req));
}

export function publicHost(req: { headers: Headers; url: string }): string {
  const first = (v: string | null): string =>
    (v ?? '').split(',')[0]?.trim() ?? '';
  let host =
    first(req.headers.get('x-forwarded-host')) ||
    first(req.headers.get('host'));
  if (!host) {
    try {
      host = new URL(req.url).host;
    } catch {
      host = '';
    }
  }
  return host.split(':')[0] ?? '';
}

export function cookieDomain(host: string): string | undefined {
  const h = host.split(':')[0] ?? '';
  if (!h || h === 'localhost' || /^[0-9.]+$/.test(h) || h.includes(':')) {
    return undefined;
  }
  const parts = h.split('.');
  if (parts.length < 2) return undefined;
  const parent = parts.length > 2 ? parts.slice(1).join('.') : h;
  return `.${parent}`;
}
