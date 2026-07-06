/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/**
 * Console-side consumption of the BFF's operation-lock status (see
 * OPERATION-LOCK-PLAN.md). These hooks own only the *polling policy* — the
 * BFF owns topology and aggregation. Policy:
 *   - hydrate once on mount (also survives reload / another user's op),
 *   - poll (~3s) ONLY while busy, and only while the tab is visible,
 *   - optimistic `markBusy()` flips busy instantly on click, bounded by a
 *     timeout so it can never hang and cleared once the server confirms,
 *   - refetch on tab focus.
 * The backend's conflict-on-start remains the safety net, so a momentarily
 * stale read can never cause a double-execute.
 */
import { useCallback, useEffect, useRef, useState } from 'react';

import {
  useGetNodeOperationStatusQuery,
  useGetSiteOperationStatusQuery,
} from '@/client/graphql/operation.generated';

const POLL_MS = 3000;
/** Max time an optimistic "busy" survives without server confirmation. */
const OPTIMISTIC_MS = 10000;

type Poller = {
  startPolling: (ms: number) => void;
  stopPolling: () => void;
  refetch: () => unknown;
};

/** True while the document is visible (SSR-safe default: true). */
function useVisible(): boolean {
  const [visible, setVisible] = useState(true);
  useEffect(() => {
    const on = () => setVisible(document.visibilityState === 'visible');
    on();
    document.addEventListener('visibilitychange', on);
    window.addEventListener('focus', on);
    return () => {
      document.removeEventListener('visibilitychange', on);
      window.removeEventListener('focus', on);
    };
  }, []);
  return visible;
}

/**
 * Optimistic busy with a race guard: stays true for OPTIMISTIC_MS after
 * markBusy() regardless of an early not-busy poll, and is dropped as soon as
 * the server confirms the lock (server then owns the truth).
 *
 * `optimisticUntil` is a timestamp. We reset it to 0 during render once the
 * server confirms busy (React's supported "adjust state while rendering"
 * pattern — converges immediately, no effect needed), and an async tick
 * forces a re-render when the window elapses.
 */
function useOptimisticBusy(serverBusy: boolean) {
  const [optimistic, setOptimistic] = useState(false);
  const timer = useRef<ReturnType<typeof setTimeout> | undefined>(undefined);

  // Server owns the truth now — consume the optimistic window (React's
  // supported "adjust state while rendering" pattern; converges immediately).
  if (serverBusy && optimistic) {
    setOptimistic(false);
  }

  const markBusy = useCallback(() => {
    setOptimistic(true);
    clearTimeout(timer.current);
    // Bounded fallback: if no lock ever appears (op finished instantly, or
    // failed to start), drop optimism instead of hanging. Runs in a timeout
    // callback, so it never fires synchronously during render/effect.
    timer.current = setTimeout(() => setOptimistic(false), OPTIMISTIC_MS);
  }, []);

  useEffect(() => () => clearTimeout(timer.current), []);

  return { busy: serverBusy || optimistic, markBusy };
}

/** Poll only while busy AND visible; refetch when the tab regains focus. */
function usePollWhileBusy(poller: Poller, busy: boolean) {
  const visible = useVisible();
  const { startPolling, stopPolling, refetch } = poller;

  useEffect(() => {
    if (busy && visible) startPolling(POLL_MS);
    else stopPolling();
    return () => stopPolling();
  }, [busy, visible, startPolling, stopPolling]);

  const wasVisible = useRef(visible);
  useEffect(() => {
    if (visible && !wasVisible.current) refetch();
    wasVisible.current = visible;
  }, [visible, refetch]);
}

export interface NodeOperationView {
  busy: boolean;
  operation:
    | { id: string; type: string; status: string; requestedBy?: string | null; startedAt?: string | null }
    | undefined;
  markBusy: () => void;
  refetch: () => void;
}

export function useNodeOperationStatus(nodeId: string | null | undefined): NodeOperationView {
  const q = useGetNodeOperationStatusQuery({
    variables: { nodeId: nodeId ?? '' },
    skip: !nodeId,
    fetchPolicy: 'cache-and-network',
    errorPolicy: 'all', // fail-open: a read error leaves busy=false
    notifyOnNetworkStatusChange: true,
  });

  const status = q.data?.getNodeOperationStatus;
  const serverBusy = !!status?.busy;
  const { busy, markBusy } = useOptimisticBusy(serverBusy);

  usePollWhileBusy(q, busy);

  return {
    busy,
    operation: status?.operation ?? undefined,
    markBusy,
    refetch: () => void q.refetch(),
  };
}

type ActionAvailability = { available: boolean; reason?: string | null };

export interface SiteOperationView {
  busy: boolean;
  degraded: boolean;
  actions: {
    restartSite: ActionAvailability;
    rf: ActionAvailability;
    service: ActionAvailability;
  };
  nodes: NonNullable<
    ReturnType<typeof useGetSiteOperationStatusQuery>['data']
  >['getSiteOperationStatus']['nodes'];
  /** Operation blocking the site (first busy node's op), for labels. */
  activeOperation:
    | { type: string; requestedBy?: string | null; startedAt?: string | null }
    | undefined;
  markBusy: () => void;
  refetch: () => void;
}

const AVAILABLE: ActionAvailability = { available: true };

export function useSiteOperationStatus(siteId: string | null | undefined): SiteOperationView {
  const q = useGetSiteOperationStatusQuery({
    variables: { siteId: siteId ?? '' },
    skip: !siteId,
    fetchPolicy: 'cache-and-network',
    errorPolicy: 'all',
    notifyOnNetworkStatusChange: true,
  });

  const status = q.data?.getSiteOperationStatus;
  const serverBusy = !!status?.busy;
  const { busy, markBusy } = useOptimisticBusy(serverBusy);

  usePollWhileBusy(q, busy);

  const actions = status?.actions;
  const activeOperation = status?.nodes.find((n) => n.busy && n.operation)?.operation ?? undefined;

  return {
    busy,
    degraded: !!status?.degraded,
    actions: {
      restartSite: actions?.restartSite ?? AVAILABLE,
      rf: actions?.rf ?? AVAILABLE,
      service: actions?.service ?? AVAILABLE,
    },
    nodes: status?.nodes ?? [],
    activeOperation: activeOperation ?? undefined,
    markBusy,
    refetch: () => void q.refetch(),
  };
}
