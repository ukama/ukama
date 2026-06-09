/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Typed URL state for /configure (ACTIVATION-PLAN §3). One zod schema
 * parses the search params; steps never read raw params. The URL carries
 * in-step context only — which step is legitimate is decided by each step's
 * self-guard against server data.
 */
'use client';

import { usePathname, useSearchParams } from 'next/navigation';
import { z } from 'zod';

export const FLOWS = ['onboarding', 'install-site', 'add-network'] as const;

const paramsSchema = z.object({
  flow: z.enum(FLOWS).catch('onboarding'),
  networkid: z.string().catch(''),
  nid: z.string().catch(''),
  /** Carried from the name step into the settings step. */
  sitename: z.string().catch(''),
  location: z.string().catch(''),
});

export type ConfigureParams = z.infer<typeof paramsSchema>;

export function useConfigureParams(): ConfigureParams {
  const sp = useSearchParams();
  return paramsSchema.parse({
    flow: sp.get('flow') ?? undefined,
    networkid: sp.get('networkid') ?? undefined,
    nid: sp.get('nid') ?? undefined,
    sitename: sp.get('sitename') ?? undefined,
    location: sp.get('location') ?? undefined,
  });
}

/** Builds a /configure step URL, carrying forward only non-empty params. */
export function stepUrl(
  path: string,
  params: Partial<ConfigureParams>,
): string {
  const q = new URLSearchParams();
  if (params.flow && params.flow !== 'onboarding') q.set('flow', params.flow);
  if (params.networkid) q.set('networkid', params.networkid);
  if (params.nid) q.set('nid', params.nid);
  if (params.sitename) q.set('sitename', params.sitename);
  if (params.location) q.set('location', params.location);
  const qs = q.toString();
  return qs ? `${path}?${qs}` : path;
}

/** Current full /configure URL (path + query) — the resume point. */
export function useCurrentConfigureUrl(): string {
  const pathname = usePathname();
  const sp = useSearchParams();
  const qs = sp.toString();
  return qs ? `${pathname}?${qs}` : pathname;
}
