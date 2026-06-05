/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Activation state for the onboarding flow (ACTIVATION-PLAN).
 *
 * State is DERIVED server-side (BFF onboardingStatus) from real
 * networks/sites/nodes — never a stored flag — so it is always correct for
 * multi-user orgs and after deletions. The console only acts on *known*
 * state: while loading or on error, nothing gates and no banner shows.
 */
'use client';

import { useOnboardingStatusQuery } from '@/client/graphql/onboarding-status.generated';

/**
 * 'hybrid'    — dashboard fully usable; setup alert bar shows until activated.
 * 'hard-gate' — every dashboard route redirects to /configure until activated.
 * Flip this single value to change the product behavior.
 */
export const ACTIVATION_MODE: 'hybrid' | 'hard-gate' = 'hybrid';

/** Activation bar: a network AND a site must exist (SIMs never block). */
const isActivated = (s: { hasNetwork: boolean; hasSite: boolean }): boolean =>
  s.hasNetwork && s.hasSite;

export interface Activation {
  /** True until the first onboardingStatus result arrives. */
  loading: boolean;
  /** Derived state, undefined while loading or on error ("unknown"). */
  status?: {
    hasNetwork: boolean;
    hasSite: boolean;
    hasNode: boolean;
    networkId: string | null;
    networkName: string | null;
  };
  /** True only when state is KNOWN and complete. */
  isActivated: boolean;
  /** True only when state is KNOWN and incomplete (drives banner/gate). */
  needsSetup: boolean;
  refetch: () => void;
}

export function useActivation(): Activation {
  const { data, loading, refetch } = useOnboardingStatusQuery();
  const s = data?.onboardingStatus;
  return {
    loading,
    status: s
      ? {
          hasNetwork: s.hasNetwork,
          hasSite: s.hasSite,
          hasNode: s.hasNode,
          networkId: s.networkId ?? null,
          networkName: s.networkName ?? null,
        }
      : undefined,
    isActivated: s ? isActivated(s) : false,
    needsSetup: s ? !isActivated(s) : false,
    refetch: () => void refetch(),
  };
}

/**
 * Where "Continue setup" (alert bar / welcome page / hard gate) should land.
 *
 * 1. Saved in-flow progress always wins — /configure steps self-guard, so a
 *    stale URL (e.g. network created meanwhile) auto-forwards correctly.
 * 2. No network yet → network creation (an independent, self-guarding step).
 * 3. Network but no site → install-site step.
 * 4. Activated → dashboard.
 */
export function resolveResumeUrl(
  status: Activation['status'],
  lastConfigureUrl: string | null,
): string {
  if (lastConfigureUrl?.startsWith('/configure')) return lastConfigureUrl;
  if (!status || !status.hasNetwork) return '/configure/network';
  if (!status.hasSite) {
    const q = status.networkId
      ? `?networkid=${encodeURIComponent(status.networkId)}`
      : '';
    return `/configure/install${q}`;
  }
  return '/';
}
