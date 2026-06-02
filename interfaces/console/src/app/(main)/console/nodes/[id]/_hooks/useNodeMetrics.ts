/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

import {
  Graphs_Type,
  Stats_Type,
} from '@/client/graphql/generated/subscriptions';
import {
  activeGraphTypeVar,
  activeNodeTabVar,
  nodeMetricsVar,
} from '@/client/vars';
import { NODE_KPIS } from '@/constants';
import { useEnvContext, useUserContext } from '@/context';
import MetricStatSubscription from '@/features/subscriptions/MetricStatSubscription';
import { TMetricResDto } from '@/types';
import { useReactiveVar } from '@apollo/client';
import { useCallback, useState } from 'react';

interface UseNodeMetricsParams {
  id: string;
  nodeType: string;
  getMetricStat: (options: { variables: { data: object } }) => void;
  cleanupSubscription: () => void;
  subscribe: (topic: string, cleanup: () => void) => void;
}

export function useNodeMetrics({
  id,
  nodeType,
  getMetricStat: _getMetricStat,
  cleanupSubscription,
  subscribe,
}: UseNodeMetricsParams) {
  const [metricFrom, setMetricFrom] = useState<number>(0);
  const [nodeUptime, setNodeUptime] = useState<number>(0);

  // Reactive vars — shared state, no prop drilling needed
  const graphType = useReactiveVar(activeGraphTypeVar);
  const selectedTab = useReactiveVar(activeNodeTabVar);
  const metrics = useReactiveVar(nodeMetricsVar);

  const { env } = useEnvContext();
  const { user } = useUserContext();

  const handleStatSubscription = useCallback(
    (_: unknown, data: string) => {
      const parsedData: TMetricResDto = JSON.parse(data);
      const { value, type, success } = parsedData.data.getMetricStatSub;
      if (success) {
        if (type === (NODE_KPIS.NODE_UPTIME as Record<string, { id: string }[]>)[nodeType]?.[0]?.id) {
          setNodeUptime(Math.floor(value[1]));
        }
        PubSub.publish(`stat-${type}`, value);
      }
    },
    [nodeType],
  );

  const startStatSubscription = useCallback(
    (from: number, statFrom: number, subscriptionKeyRef: React.MutableRefObject<string | null>, subscriptionControllerRef: React.MutableRefObject<{ cancel: () => void } | null>) => {
      const sKey = `stat-${user.orgName}-${user.id}-${Stats_Type.AllNode}-${from}`;
      cleanupSubscription();
      subscriptionKeyRef.current = sKey;

      const controller = MetricStatSubscription({
        key: sKey,
        nodeId: id,
        userId: user.id,
        url: env.METRIC_URL,
        orgName: user.orgName,
        type: Stats_Type.AllNode,
        from: statFrom,
      });

      subscriptionControllerRef.current = controller;
      PubSub.subscribe(sKey, handleStatSubscription);
      subscribe(sKey, () => PubSub.unsubscribe(sKey));
    },
    [id, user.id, user.orgName, env.METRIC_URL, cleanupSubscription, handleStatSubscription, subscribe],
  );

  return {
    metricFrom,
    setMetricFrom,
    graphType,
    setGraphType: (type: Graphs_Type) => activeGraphTypeVar(type),
    nodeUptime,
    setNodeUptime,
    selectedTab,
    setSelectedTab: (tab: number) => activeNodeTabVar(tab),
    metrics,
    setMetrics: (data: { metrics: unknown[] }) => nodeMetricsVar(data as Parameters<typeof nodeMetricsVar>[0]),
    handleStatSubscription,
    startStatSubscription,
  };
}
