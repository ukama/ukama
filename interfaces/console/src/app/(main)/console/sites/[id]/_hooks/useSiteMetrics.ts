/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

import { Node, SiteDto } from '@/client/graphql/generated';
import { activeGraphTypeVar, siteActiveSubscribersVar } from '@/client/vars';
import { SITE_KPI_TYPES } from '@/constants';
import { useEnvContext, useUserContext } from '@/context';
import { ActiveView } from '@/types';
import { extractMetricValue } from '@/utils';
import { useSubscriptionManager } from '@/features/subscriptions/useSubscriptionManager';
import { useMetricSubscriptions } from '@/utils/useMetricSubscriptions';
import { useReactiveVar } from '@apollo/client';
import { useCallback, useEffect, useMemo, useRef } from 'react';
import PubSub from 'pubsub-js';

type SiteMetric = {
  type: string;
  success: boolean;
  siteId?: string;
  nodeId?: string;
  value: number | number[];
};

function getSiteActiveSubscribers(
  metricsData: { metrics?: SiteMetric[] } | null | undefined,
  siteId: string,
): number | null {
  if (!metricsData?.metrics || !siteId) return null;

  const subscriberMetrics = metricsData.metrics.filter(
    (m) =>
      m.type === SITE_KPI_TYPES.ACTIVE_SUBSCRIBERS &&
      m.success === true &&
      m.siteId === siteId,
  );

  if (subscriberMetrics.length === 0) return null;

  return subscriberMetrics.reduce((total, metric) => {
    const value = extractMetricValue(metric.value);
    return total + (value || 0);
  }, 0);
}

export function getInitialNodeUptimesFromMetrics(
  metrics: SiteMetric[] | undefined,
): Record<string, number> {
  if (!metrics?.length) return {};

  return metrics.reduce<Record<string, number>>((acc, metric) => {
    if (
      metric.type === SITE_KPI_TYPES.NODE_UPTIME &&
      metric.nodeId &&
      metric.success
    ) {
      acc[metric.nodeId] = typeof metric.value === 'number' ? metric.value : 0;
    }
    return acc;
  }, {});
}

export function getSiteUptimeFromMetrics(
  metrics: SiteMetric[] | undefined,
  siteId: string,
): number {
  if (!metrics?.length || !siteId) return 0;

  const siteMetrics = metrics.filter((m) => m.siteId === siteId && m.success);
  const uptimeMetric = siteMetrics.find(
    (m) => m.type === SITE_KPI_TYPES.SITE_UPTIME,
  );

  if (uptimeMetric?.value !== undefined) {
    const v = uptimeMetric.value;
    const num = typeof v === 'number' ? v : parseFloat(String(v));
    return Math.floor(num);
  }
  return 0;
}

export function useSiteMetrics(
  id: string,
  activeSite: SiteDto,
  nodes: Node[],
  nodesFetched: boolean,
  activeView: ActiveView,
) {
  const { env, subscriptionClient } = useEnvContext();
  const { user } = useUserContext();
  const { subscribe } = useSubscriptionManager();

  // Reactive var — no useState/prop needed, any component can read this directly
  const activeSubscribers = useReactiveVar(siteActiveSubscribersVar);
  const subscribersSubscriptionRef = useRef<string | null>(null);

  // Sync activeView.graphType into the reactive var so other consumers stay in sync
  useEffect(() => {
    activeGraphTypeVar(activeView.graphType);
  }, [activeView.graphType]);

  const {
    metrics,
    metricFrom,
    metricsLoading,
    statData,
    statLoading,
    resetMetrics,
    cleanupSubscriptions,
  } = useMetricSubscriptions({
    siteId: id,
    userId: user.id,
    orgName: user.orgName,
    metricUrl: env.METRIC_URL,
    subscriptionClient: subscriptionClient!,
    activeGraphType: activeView.graphType,
    nodeIds: nodes.map((node) => node.id),
    nodesFetched,
  });

  const handleSubscribersUpdate = useCallback((_: unknown, data: unknown) => {
    if (data !== null && data !== undefined) {
      const value =
        Array.isArray(data) && data.length > 1
          ? extractMetricValue(data[1])
          : extractMetricValue(data);
      if (value !== null) siteActiveSubscribersVar(value);
    }
  }, []);

  useEffect(() => {
    if (!id || !activeSite.id) return;

    // Reset stale count when site changes
    siteActiveSubscribersVar(0);

    if (subscribersSubscriptionRef.current) {
      PubSub.unsubscribe(subscribersSubscriptionRef.current);
      subscribersSubscriptionRef.current = null;
    }

    const topic = `stat-${SITE_KPI_TYPES.ACTIVE_SUBSCRIBERS}-${id}`;
    const token = PubSub.subscribe(topic, handleSubscribersUpdate);
    subscribersSubscriptionRef.current = token;
    subscribe(topic, () => PubSub.unsubscribe(token));

    return () => {
      if (subscribersSubscriptionRef.current) {
        PubSub.unsubscribe(subscribersSubscriptionRef.current);
        subscribersSubscriptionRef.current = null;
      }
    };
  }, [id, activeSite.id, handleSubscribersUpdate]);

  useEffect(() => {
    if (statData?.getSiteStat && activeSite.id) {
      const count = getSiteActiveSubscribers(
        statData.getSiteStat as { metrics?: SiteMetric[] },
        activeSite.id,
      );
      if (count !== null) siteActiveSubscribersVar(count);
    }
  }, [statData, activeSite.id]);

  useEffect(
    () => () => {
      cleanupSubscriptions();
      if (subscribersSubscriptionRef.current) {
        PubSub.unsubscribe(subscribersSubscriptionRef.current);
        subscribersSubscriptionRef.current = null;
      }
    },
    [cleanupSubscriptions],
  );

  const siteMetrics = statData?.getSiteStat?.metrics as
    | SiteMetric[]
    | undefined;

  const initialNodeUptimes = useMemo(
    () => getInitialNodeUptimesFromMetrics(siteMetrics),
    [siteMetrics],
  );

  const siteUptime = useMemo(
    () => getSiteUptimeFromMetrics(siteMetrics, activeSite.id),
    [siteMetrics, activeSite.id],
  );

  return {
    metrics,
    metricFrom,
    metricsLoading,
    statData,
    statLoading,
    resetMetrics,
    activeSubscribers,
    initialNodeUptimes,
    siteUptime,
  };
}
