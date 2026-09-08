/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/**
 * Timestamp of the last completed fetch for a query, for the header's
 * "Updated ..." label.
 *
 * Apollo skips `onCompleted` when a poll resolves to data identical to what
 * the cache already holds, which is the common case for a KPI that has not
 * moved — so the label would keep counting up while requests were in fact
 * succeeding. `networkStatus` transitions on every round trip regardless of
 * the response, so it is what the timestamp tracks. Pass the status from a
 * query declared with `notifyOnNetworkStatusChange: true`; without that flag
 * the status does not change on a poll.
 */
import { useState } from 'react';
import { NetworkStatus } from '@apollo/client';

export function useLastFetched(networkStatus: NetworkStatus): Date {
  const [fetchedAt, setFetchedAt] = useState(() => new Date());
  const [seen, setSeen] = useState(networkStatus);

  if (seen !== networkStatus) {
    setSeen(networkStatus);
    if (networkStatus === NetworkStatus.ready) setFetchedAt(new Date());
  }

  return fetchedAt;
}
