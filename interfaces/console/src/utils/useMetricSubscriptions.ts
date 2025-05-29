/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { useEffect, useRef, useState, useCallback } from 'react';
import PubSub from 'pubsub-js';
import { ApolloClient, NormalizedCacheObject } from '@apollo/client';
import {
  Graphs_Type,
  Stats_Type,
  useGetMetricBySiteLazyQuery,
  useGetSiteStatLazyQuery,
  MetricsRes,
} from '@/client/graphql/generated/subscriptions';
import { getUnixTime } from '@/utils';
import { METRIC_RANGE_10800, STAT_STEP_29 } from '@/constants';
import MetricStatBySiteSubscription from '@/lib/MetricStatBySiteSubscription';
import { TMetricResDto } from '@/types';

interface MetricSubscriptionsProps {
  siteId: string;
  userId: string;
  orgName: string;
  metricUrl: string;
  subscriptionClient: ApolloClient<NormalizedCacheObject>;
  activeGraphType: Graphs_Type;
  nodeIds: string[];
  nodesFetched: boolean;
}

export const useMetricSubscriptions = ({
  siteId,
  userId,
  orgName,
  metricUrl,
  subscriptionClient,
  activeGraphType,
  nodeIds,
  nodesFetched,
}: MetricSubscriptionsProps) => {
  const subscriptionsRef = useRef<Record<string, boolean>>({});
  const [metrics, setMetrics] = useState<MetricsRes>({ metrics: [] });
  const [metricFrom, setMetricFrom] = useState<number>(0);
  const [isInitialized, setIsInitialized] = useState(false);

  const cleanupSubscriptions = useCallback(() => {
    Object.keys(subscriptionsRef.current).forEach((topic) => {
      PubSub.unsubscribe(topic);
      delete subscriptionsRef.current[topic];
    });
  }, []);

  const [
    getMetricBySite,
    { loading: metricsLoading, variables: metricsVariables },
  ] = useGetMetricBySiteLazyQuery({
    client: subscriptionClient,
    fetchPolicy: 'network-only',
    onCompleted: (data) => {
      setMetrics(data.getMetricBySite);
    },
  });

  const fetchMetrics = useCallback(() => {
    if (!siteId || !userId || !orgName || metricFrom <= 0 || !activeGraphType) {
      return;
    }

    const topic = `${userId}/${activeGraphType}/${metricFrom}`;
    subscriptionsRef.current[topic] = true;

    getMetricBySite({
      variables: {
        data: {
          step: 30,
          siteId,
          userId,
          type: activeGraphType,
          from: metricFrom,
          orgName,
          withSubscription: true,
          to: metricFrom + METRIC_RANGE_10800,
        },
      },
    });
  }, [siteId, userId, orgName, metricFrom, activeGraphType, getMetricBySite]);

  useEffect(() => {
    if (
      isInitialized &&
      metricFrom > 0 &&
      activeGraphType &&
      siteId &&
      metricsVariables?.data?.from !== metricFrom
    ) {
      fetchMetrics();
    }
  }, [
    metricFrom,
    activeGraphType,
    siteId,
    fetchMetrics,
    isInitialized,
    metricsVariables,
  ]);

  const [
    getSiteMetricStat,
    { data: statData, loading: statLoading, variables: statVar },
  ] = useGetSiteStatLazyQuery({
    client: subscriptionClient,
    fetchPolicy: 'network-only',
    onCompleted: (data) => {
      const sKey = `stat-${orgName}-${userId}-${Stats_Type.Site}-${
        statVar?.data.from ?? 0
      }`;

      if (data.getSiteStat.metrics.length > 0) {
        subscriptionsRef.current[sKey] = true;

        MetricStatBySiteSubscription({
          key: sKey,
          nodeIds: [],
          siteIds: [siteId],
          userId,
          url: metricUrl,
          orgName,
          type: Stats_Type.Site,
          from: statVar?.data.from ?? 0,
        });

        PubSub.subscribe(sKey, handleSiteStatSubscription);
      }
    },
  });

  const handleSiteStatSubscription = useCallback((_: any, data: string) => {
    try {
      const parsedData: TMetricResDto = JSON.parse(data);
      if (parsedData?.data?.getSiteMetricStatSub) {
        const { type, success, nodeId, value } =
          parsedData.data.getSiteMetricStatSub;

        if (success) {
          if (nodeId) {
            PubSub.publish(`stat-${type}-${nodeId}`, value);
          }
          PubSub.publish(`stat-${type}`, value);
        }
      }
    } catch (error) {
      console.error('Error in handleSiteStatSubscription:', error);
    }
  }, []);

  useEffect(() => {
    if (siteId && userId && orgName) {
      setIsInitialized(false);
      cleanupSubscriptions();
      setMetrics({ metrics: [] });

      const newMetricFrom = getUnixTime() - METRIC_RANGE_10800;
      setMetricFrom(newMetricFrom);

      setIsInitialized(true);
    }

    return () => {
      cleanupSubscriptions();
    };
  }, [siteId, userId, orgName, cleanupSubscriptions]);

  useEffect(() => {
    if (!nodesFetched || !siteId) return;

    const to = getUnixTime();
    const from = to - STAT_STEP_29;
    const sKey = `stat-${orgName}-${userId}-${Stats_Type.Site}-${from}`;

    Object.keys(subscriptionsRef.current).forEach((topic) => {
      if (topic.startsWith(`stat-${orgName}-${userId}-${Stats_Type.Site}`)) {
        PubSub.unsubscribe(topic);
        delete subscriptionsRef.current[topic];
      }
    });

    subscriptionsRef.current[sKey] = true;

    getSiteMetricStat({
      variables: {
        data: {
          to,
          from,
          userId,
          step: STAT_STEP_29,
          orgName,
          withSubscription: true,
          type: Stats_Type.Site,
          siteIds: [siteId],
          nodeIds: [],
        },
      },
    });

    return () => {
      PubSub.unsubscribe(sKey);
      delete subscriptionsRef.current[sKey];
    };
  }, [siteId, nodesFetched, userId, orgName, getSiteMetricStat]);

  const resetMetrics = useCallback(() => {
    setMetrics({ metrics: [] });
    setMetricFrom(getUnixTime() - METRIC_RANGE_10800);
  }, []);

  useEffect(() => {
    return () => {
      cleanupSubscriptions();
    };
  }, [cleanupSubscriptions]);

  return {
    metrics,
    metricFrom,
    metricsLoading,
    statData,
    statLoading,
    resetMetrics,
    cleanupSubscriptions,
  };
};
