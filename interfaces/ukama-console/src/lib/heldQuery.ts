/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Reads an Apollo query result so a screen only ever renders a response it
 * has actually been given.
 *
 * Two states hand a screen an empty result with `loading: false`, and both
 * mean "no answer yet" rather than "nothing there":
 *
 *  - Cache expiry. Every entity carries a TTL (`invalidationPolicies` in
 *    client/apollo.ts — 5 minutes by default, 1 for NodeDto, 2 for SiteDto).
 *    On expiry the entities are evicted and the watched query re-reads an
 *    emptied result while it refetches in the background.
 *  - A skipped query. One gated on `skip: !networkId` sits idle while
 *    `useNetworkId` validates the id; callers pass `useNetworkQueryPending()`
 *    as `pending` to cover that.
 *
 * `data` therefore falls back to the last delivered response, and `loading`
 * reports whether a first response is still owed — so an empty state means
 * the backend said empty, and a skeleton shows only before anything has
 * arrived, never during a background refetch.
 */
export interface QueryResultLike<T> {
  data?: T;
  previousData?: T;
  loading: boolean;
}

export interface HeldQuery<T> {
  /** The current response, or the last one delivered while it is unavailable. */
  data: T | undefined;
  /** True until the first response lands: render a skeleton, not an empty state. */
  loading: boolean;
}

export function heldQuery<T>(
  result: QueryResultLike<T>,
  pending = false,
): HeldQuery<T> {
  const data = result.data ?? result.previousData;

  return { data, loading: data === undefined && (result.loading || pending) };
}
