/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * v1 freshness mechanism = screen-scoped polling, paused on hidden tabs
 * (BUILD-PLAN §5.1·4). No WebSockets/subscriptions in version 1.
 *
 * POLLING IS CURRENTLY DISABLED GLOBALLY: queries fetch on mount (subject to
 * the Apollo TTL cache) and on demand (refetch). The plumbing stays in place
 * at every call site — enable per screen by passing `enabled: true`, or
 * globally via NEXT_PUBLIC_ENABLE_POLLING=true, when we turn screens live
 * case by case.
 *
 * Usage: useFooQuery({ variables, ...visiblePoll(POLL_OVERVIEW_MS) })
 * Opt-in later: useFooQuery({ variables, ...visiblePoll(POLL_OVERVIEW_MS, true) })
 */
export const POLL_LIVE_MS = 30_000; // node/site detail health & KPIs
export const POLL_OVERVIEW_MS = 60_000; // home/list status screens

const POLLING_ENABLED = process.env.NEXT_PUBLIC_ENABLE_POLLING === 'true';

export const visiblePoll = (
  pollInterval: number,
  enabled: boolean = POLLING_ENABLED
): { pollInterval?: number; skipPollAttempt?: () => boolean } => {
  if (!enabled) return {};
  return {
    pollInterval,
    /** Apollo skips the tick entirely while the tab is hidden. */
    skipPollAttempt: () =>
      typeof document !== 'undefined' && document.visibilityState === 'hidden',
  };
};
