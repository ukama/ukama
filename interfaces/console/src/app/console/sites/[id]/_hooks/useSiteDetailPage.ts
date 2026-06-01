/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

import {
  Node,
  NodeStateEnum,
  NodeTypeEnum,
  SiteDto,
  Timeframe_Filter,
  useGetHealthReportQuery,
  useGetNodesForSiteLazyQuery,
  useGetSitesQuery,
  useToggleInternetSwitchMutation,
  useToggleRfStatusMutation,
  useToggleServiceMutation,
} from '@/client/graphql/generated';
import { Graphs_Type } from '@/client/graphql/generated/subscriptions';
import {
  NODE_ACTIONS_ENUM,
  SITE_ACTIONS_BUTTONS,
  SITE_KPI_TYPES,
  SITE_KPIS,
} from '@/constants';
import { SectionData } from '@/constants/index';
import { useEnvContext, useUserContext, useNetworkContext, useUIContext } from '@/context';
import { ActiveView, KPIType, TSiteActionToggle, TStatusBarObj } from '@/types';
import {
  extractMetricValue,
  graphTypeToSection,
  kpiToGraphType,
  stringToBoolean,
} from '@/utils';
import { useFetchAddress } from '@/utils/useFetchAddress';
import { useMetricSubscriptions } from '@/utils/useMetricSubscriptions';
import { AlertColor } from '@mui/material';
import { useRouter } from 'next/navigation';
import PubSub from 'pubsub-js';
import { useCallback, useEffect, useMemo, useRef, useState } from 'react';

type SiteMetric = {
  type: string;
  success: boolean;
  siteId?: string;
  nodeId?: string;
  value: number | number[];
};

const defaultSite: SiteDto = {
  id: '',
  accessId: '',
  backhaulId: '',
  createdAt: '',
  installDate: '',
  isDeactivated: false,
  latitude: '',
  location: '',
  longitude: '',
  name: '',
  networkId: '',
  powerId: '',
  spectrumId: '',
  switchId: '',
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

function getInitialNodeUptimesFromMetrics(
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

function getSiteUptimeFromMetrics(
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

export function useSiteDetailPage(id: string) {
  const router = useRouter();
  const [activeSite, setActiveSite] = useState<SiteDto>(defaultSite);
  const [nodes, setNodes] = useState<Node[]>([]);
  const [nodesFetched, setNodesFetched] = useState(false);
  const [isDataReady, setIsDataReady] = useState(false);
  const [activeSubscribers, setActiveSubscribers] = useState<number>(0);
  const [siteActionData, setSiteActionData] = useState<TSiteActionToggle[]>([]);
  const [activeView, setActiveView] = useState<ActiveView>({
    graphType: Graphs_Type.Solar,
    kpi: 'node',
  });

  const subscribersSubscriptionRef = useRef<string | null>(null);

  const { env, subscriptionClient } = useEnvContext();
  const { user } = useUserContext();
  const { setSelectedDefaultSite } = useNetworkContext();
  const { setSnackbarMessage } = useUIContext();

  const notify = (msgId: string, message: string, type: string | AlertColor) =>
    setSnackbarMessage({ id: msgId, message, type, show: true });

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

  const {
    address: currentSiteAddress,
    isLoading: currentSiteAddressLoading,
    error: addressError,
    fetchAddress,
  } = useFetchAddress();

  const sections: SectionData = useMemo(
    () => ({
      SOLAR: SITE_KPIS.SOLAR.metrics,
      BATTERY: SITE_KPIS.BATTERY.metrics,
      CONTROLLER: SITE_KPIS.CONTROLLER.metrics,
      MAIN_BACKHAUL: SITE_KPIS.MAIN_BACKHAUL.metrics,
      SWITCH: SITE_KPIS.SWITCH.metrics,
    }),
    [],
  );

  const getSectionName = useCallback(
    (graphType: Graphs_Type): string =>
      graphTypeToSection[graphType] || 'SOLAR',
    [],
  );

  const [updateSwitchPort] = useToggleInternetSwitchMutation({
    onError: (err) => notify('update-node-err-msg', err.message, 'error'),
  });

  const handleSubscribersUpdate = useCallback((_: unknown, data: unknown) => {
    if (data !== null && data !== undefined) {
      const value =
        Array.isArray(data) && data.length > 1
          ? extractMetricValue(data[1])
          : extractMetricValue(data);
      if (value !== null) setActiveSubscribers(value);
    }
  }, []);

  useEffect(() => {
    if (!id || !activeSite.id) return;

    if (subscribersSubscriptionRef.current) {
      PubSub.unsubscribe(subscribersSubscriptionRef.current);
      subscribersSubscriptionRef.current = null;
    }

    const topic = `stat-${SITE_KPI_TYPES.ACTIVE_SUBSCRIBERS}-${id}`;
    subscribersSubscriptionRef.current = PubSub.subscribe(
      topic,
      handleSubscribersUpdate,
    );

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
      if (count !== null) setActiveSubscribers(count);
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

  const handleViewChange = useCallback(
    (kpiType: string): void => {
      setActiveView({
        graphType: kpiToGraphType[kpiType] || Graphs_Type.Solar,
        kpi: kpiType as KPIType,
      });
      resetMetrics();
    },
    [resetMetrics],
  );

  const handleSwitchChange = useCallback(
    async (portNumber: number, currentStatus: boolean) => {
      const newStatus = !currentStatus;
      try {
        const result = await updateSwitchPort({
          variables: {
            data: { port: portNumber, siteId: id, status: newStatus },
          },
        });
        if (result.data?.toggleInternetSwitch?.success) {
          notify(
            'update-switch-success',
            `Port ${portNumber} status updated to ${newStatus ? 'On' : 'Off'}`,
            'success',
          );
        }
      } catch (err) {
        notify(
          'update-site-error',
          err instanceof Error ? err.message : 'Unknown error',
          'error',
        );
      }
    },
    [id, updateSwitchPort, setSnackbarMessage],
  );

  const { data: siteData } = useGetSitesQuery({
    fetchPolicy: 'cache-and-network',
    variables: { data: {} },
    onError: (err) =>
      notify('fetching-sites-msg', err.message, 'error' as AlertColor),
  });

  const [fetchNodesForSite] = useGetNodesForSiteLazyQuery({
    fetchPolicy: 'cache-first',
    onCompleted: (data) => {
      const filtered = data.getNodesForSite.nodes.filter(
        (node) =>
          node.site.siteId === activeSite.id &&
          node.status.state === NodeStateEnum.Configured,
      ) as Node[];
      setNodes(filtered);
      setNodesFetched(true);
    },
  });

  const [toggleRFStatus, { loading: toggleRFStatusLoading }] =
    useToggleRfStatusMutation({
      fetchPolicy: 'network-only',
      onCompleted: (_, ctx) => {
        notify(
          'toggle-rf-status-success-msg',
          `RF status turned ${ctx?.variables?.data?.status ? 'On' : 'Off'} successfully.`,
          'success',
        );
      },
      onError: (_, ctx) => {
        notify(
          'toggle-rf-status-error-msg',
          `Failed to turn RF status ${ctx?.variables?.data?.status ? 'On' : 'Off'}.`,
          'error',
        );
      },
    });

  const [toggleService, { loading: toggleServiceLoading }] =
    useToggleServiceMutation({
      fetchPolicy: 'network-only',
      onCompleted: (_, ctx) => {
        notify(
          'toggle-service-status-success-msg',
          `Service status turned ${ctx?.variables?.data?.status ? 'On' : 'Off'} successfully.`,
          'success',
        );
      },
      onError: (_, ctx) => {
        notify(
          'toggle-service-status-error-msg',
          `Failed to turn service status ${ctx?.variables?.data?.status ? 'On' : 'Off'}.`,
          'error',
        );
      },
    });

  const { loading: healthLoading } = useGetHealthReportQuery({
    variables: {
      data: {
        id: '',
        timestamp: '',
        timeframe: Timeframe_Filter.Latest,
        nodeId:
          nodes.find((node) => node.id.includes(NodeTypeEnum.Tnode))?.id || '',
      },
    },
    onCompleted: (data) => {
      if (data.getHealthReport.system.length > 0) {
        const actions: TSiteActionToggle[] = [];
        data.getHealthReport.system.forEach((system) => {
          if (system.name === 'radio') {
            actions.push({
              id: NODE_ACTIONS_ENUM.TOGGLE_RADIO,
              key: 'radio',
              value: stringToBoolean(system.value),
            });
          }
          if (system.name === 'service') {
            actions.push({
              id: NODE_ACTIONS_ENUM.TOGGLE_SERVICE,
              key: 'service',
              value: stringToBoolean(system.value),
            });
          }
        });
        setSiteActionData(actions);
      }
    },
    onError: (err) =>
      notify('fetching-health-report-msg', err.message, 'error'),
  });

  const checkDataReadiness = useCallback(() => {
    if (activeSite.id && currentSiteAddress && !currentSiteAddressLoading) {
      setIsDataReady(true);
    }
  }, [activeSite.id, currentSiteAddress, currentSiteAddressLoading]);

  const filterActiveSite = useCallback(
    (siteId: string) => {
      const found = siteData?.getSites.sites.find((s) => s.id === siteId);
      if (found) {
        setActiveSite(found);
        checkDataReadiness();
      }
    },
    [siteData, checkDataReadiness],
  );

  useEffect(() => {
    if (id && user.id && user.orgName) {
      filterActiveSite(id);
      setActiveView({ graphType: Graphs_Type.NodeHealth, kpi: 'node' });
    }
  }, [id, user.id, user.orgName, filterActiveSite]);

  useEffect(() => {
    if (siteData?.getSites?.sites) {
      const found = siteData.getSites.sites.find((s) => s.id === id);
      if (found) {
        setActiveSite(found);
      } else if (siteData.getSites.sites.length > 0 && id === '') {
        router.push('/console/sites/' + siteData.getSites.sites[0].id);
      }
    }
  }, [id, siteData, router]);

  const handleSiteChange = useCallback(
    (newSiteId: string) => {
      router.push('/console/sites/' + newSiteId);
    },
    [router],
  );

  useEffect(() => {
    const handleFetchAddress = async () => {
      if (activeSite.latitude && activeSite.longitude) {
        await fetchAddress(
          activeSite.latitude.toString(),
          activeSite.longitude.toString(),
        );
      }
    };
    setSelectedDefaultSite(activeSite.name);
    if (activeSite.id && activeSite.latitude && activeSite.longitude) {
      handleFetchAddress();
    }
  }, [activeSite, fetchAddress, setSelectedDefaultSite]);

  useEffect(() => {
    if (addressError) {
      notify(
        'error-fetching-address',
        'Error fetching address from coordinates',
        'error',
      );
    }
  }, [addressError]);

  useEffect(() => {
    if (activeSite.id) {
      setNodesFetched(false);
      fetchNodesForSite({ variables: { siteId: activeSite.id } });
    }
  }, [activeSite.id, fetchNodesForSite]);

  const handleActionClick = useCallback(
    (actionId: string, value: boolean) => {
      const tnodeId =
        nodes.find((node) => node.id.includes(NodeTypeEnum.Tnode))?.id ?? '';
      switch (actionId) {
        case NODE_ACTIONS_ENUM.TOGGLE_RADIO:
          toggleRFStatus({
            variables: { data: { nodeId: tnodeId, status: value } },
          });
          break;
        case NODE_ACTIONS_ENUM.TOGGLE_SERVICE:
          toggleService({
            variables: { data: { nodeId: tnodeId, status: value } },
          });
          break;
      }
      setSiteActionData((prev) =>
        prev.map((item) => (item.id === actionId ? { ...item, value } : item)),
      );
    },
    [nodes, siteActionData, toggleRFStatus, toggleService],
  );

  const siteMetrics = statData?.getSiteStat?.metrics as
    | SiteMetric[]
    | undefined;
  const initialNodeUptimes = getInitialNodeUptimesFromMetrics(siteMetrics);
  const siteUptime = getSiteUptimeFromMetrics(siteMetrics, activeSite.id);

  return {
    id,
    activeSite,
    nodes,
    isDataReady,
    activeSubscribers,
    siteActionData,
    activeView,
    sections,
    metrics,
    metricFrom,
    metricsLoading,
    statData,
    statLoading,
    siteData,
    currentSiteAddress,
    healthLoading,
    toggleRFStatusLoading,
    toggleServiceLoading,
    initialNodeUptimes,
    siteUptime,
    actionOptions: SITE_ACTIONS_BUTTONS,
    getSectionName,
    handleViewChange,
    handleSwitchChange,
    handleSiteChange,
    handleActionClick,
    handleSelected: (obj: TStatusBarObj) => handleSiteChange(obj.id),
  };
}
