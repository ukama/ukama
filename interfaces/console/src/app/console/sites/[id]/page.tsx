/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';

import {
  NodeConnectivityEnum,
  NodeStateEnum,
  SiteDto,
  useGetNodesLazyQuery,
  useGetSitesQuery,
  useToggleInternetSwitchMutation,
} from '@/client/graphql/generated';
import { Graphs_Type } from '@/client/graphql/generated/subscriptions';
import SiteComponents from '@/components/SiteComponents';
import SiteDetailsHeader from '@/components/SiteDetailsHeader';
import SiteInfo from '@/components/SiteInfos';
import SiteOverview from '@/components/SiteOverView';
import { SITE_KPI_TYPES, SITE_KPIS } from '@/constants';
import { SectionData } from '@/constants/index';
import { useAppContext } from '@/context';
import { ActiveView, KPIType } from '@/types';
import {
  extractMetricValue,
  graphTypeToSection,
  kpiToGraphType,
} from '@/utils';
import { useFetchAddress } from '@/utils/useFetchAddress';
import { useMetricSubscriptions } from '@/utils/useMetricSubscriptions';
import { AlertColor, Box, Grid2, Skeleton } from '@mui/material';
import dynamic from 'next/dynamic';
import { useRouter } from 'next/navigation';
import PubSub from 'pubsub-js';
import React, {
  useCallback,
  useEffect,
  useMemo,
  useRef,
  useState,
} from 'react';

const SiteMapComponent = dynamic(
  () => import('@/components/SiteMapComponent'),
  {
    ssr: false,
  },
);

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

interface SiteDetailsProps {
  params: {
    id: string;
  };
}

const getSiteActiveSubscribers = (
  metricsData: any,
  siteId: string,
): number | null => {
  if (!metricsData || !metricsData.metrics || !siteId) return null;

  const subscriberMetrics = metricsData.metrics.filter(
    (m: any) =>
      m.type === SITE_KPI_TYPES.ACTIVE_SUBSCRIBERS &&
      m.success === true &&
      m.siteId === siteId,
  );

  if (subscriberMetrics.length === 0) return null;

  const totalSubscribers = subscriberMetrics.reduce(
    (total: number, metric: any) => {
      const value = extractMetricValue(metric.value);
      return total + (value || 0);
    },
    0,
  );

  return totalSubscribers;
};

const Page: React.FC<SiteDetailsProps> = ({ params }) => {
  const { id } = params;
  const router = useRouter();
  const [activeSite, setActiveSite] = useState<SiteDto>(defaultSite);
  const [nodeIds, setNodeIds] = useState<string[]>([]);
  const [nodesFetched, setNodesFetched] = useState(false);
  const [isDataReady, setIsDataReady] = useState(false);
  const [activeSubscribers, setActiveSubscribers] = useState<number>(0);
  const [activeView, setActiveView] = useState<ActiveView>({
    graphType: Graphs_Type.Solar,
    kpi: 'node',
  });

  const subscribersSubscriptionRef = useRef<string | null>(null);

  const {
    setSnackbarMessage,
    setSelectedDefaultSite,
    user,
    env,
    subscriptionClient,
  } = useAppContext();

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
    nodeIds,
    nodesFetched,
  });

  const {
    address: CurrentSiteaddress,
    isLoading: CurrentSiteAddressLoading,
    error,
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

  const getSectionName = useCallback((graphType: Graphs_Type): string => {
    return graphTypeToSection[graphType] || 'SOLAR';
  }, []);

  const [updateSwitchPort] = useToggleInternetSwitchMutation({
    onError: (err) => {
      setSnackbarMessage({
        id: 'update-node-err-msg',
        message: err.message,
        type: 'error',
        show: true,
      });
    },
  });

  const handleSubscribersUpdate = useCallback((_: any, data: any) => {
    if (data !== null && data !== undefined) {
      const value =
        Array.isArray(data) && data.length > 1
          ? extractMetricValue(data[1])
          : extractMetricValue(data);

      if (value !== null) {
        setActiveSubscribers(value);
      }
    }
  }, []);

  useEffect(() => {
    if (!id || !activeSite.id) return;

    if (subscribersSubscriptionRef.current) {
      PubSub.unsubscribe(subscribersSubscriptionRef.current);
      subscribersSubscriptionRef.current = null;
    }

    const subscribersTopic = `stat-${SITE_KPI_TYPES.ACTIVE_SUBSCRIBERS}-${id}`;
    subscribersSubscriptionRef.current = PubSub.subscribe(
      subscribersTopic,
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
      const initialSubscribers = getSiteActiveSubscribers(
        statData.getSiteStat,
        activeSite.id,
      );
      if (initialSubscribers !== null) {
        setActiveSubscribers(initialSubscribers);
      }
    }
  }, [statData, activeSite.id]);

  useEffect(() => {
    return () => {
      cleanupSubscriptions();
      if (subscribersSubscriptionRef.current) {
        PubSub.unsubscribe(subscribersSubscriptionRef.current);
        subscribersSubscriptionRef.current = null;
      }
    };
  }, [cleanupSubscriptions]);

  const handleViewChange = useCallback(
    (kpiType: string): void => {
      const graphType = kpiToGraphType[kpiType] || Graphs_Type.Solar;

      setActiveView({
        graphType,
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
            data: {
              port: portNumber,
              siteId: id,
              status: newStatus,
            },
          },
        });

        if (result.data?.toggleInternetSwitch?.success) {
          setSnackbarMessage({
            id: 'update-switch-success',
            message: `Port ${portNumber} status updated successfully to ${
              newStatus ? 'On' : 'Off'
            }`,
            type: 'success',
            show: true,
          });
        }
      } catch (error) {
        const errorMessage =
          error instanceof Error ? error.message : 'An unknown error occurred';
        setSnackbarMessage({
          id: 'update-site-error',
          message: errorMessage,
          type: 'error',
          show: true,
        });
      }
    },
    [id, updateSwitchPort, setSnackbarMessage],
  );

  const { data: siteData, loading: sitesLoading } = useGetSitesQuery({
    fetchPolicy: 'cache-first',
    nextFetchPolicy: 'cache-and-network',
    variables: {
      data: {},
    },
    onError: (error) => {
      setSnackbarMessage({
        id: 'fetching-sites-msg',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });

  const [fetchNodesForSite] = useGetNodesLazyQuery({
    fetchPolicy: 'cache-first',
    onCompleted: (data) => {
      const nodeIds = data.getNodes.nodes
        .filter(
          (node) =>
            node.latitude !== null &&
            node.site.siteId === activeSite.id &&
            node.longitude !== null &&
            node.status.connectivity === NodeConnectivityEnum.Online &&
            node.status.state === NodeStateEnum.Configured,
        )
        .map((node) => node.id);
      setNodeIds(nodeIds);
      setNodesFetched(true);
    },
  });

  const checkDataReadiness = useCallback(() => {
    if (activeSite.id && CurrentSiteaddress && !CurrentSiteAddressLoading) {
      setIsDataReady(true);
    }
  }, [activeSite.id, CurrentSiteaddress, CurrentSiteAddressLoading]);

  const filterActiveSite = useCallback(
    (siteId: string) => {
      const foundSite = siteData?.getSites.sites.find(
        (site) => site.id === siteId,
      );
      if (foundSite) {
        setActiveSite(foundSite);
        checkDataReadiness();
      }
    },
    [siteData, checkDataReadiness],
  );

  useEffect(() => {
    if (id && user.id && user.orgName) {
      filterActiveSite(id);
      setActiveView({
        graphType: Graphs_Type.NodeHealth,
        kpi: 'node',
      });
    }
  }, [id, user.id, user.orgName, filterActiveSite]);

  useEffect(() => {
    if (siteData?.getSites?.sites) {
      const foundSite = siteData.getSites.sites.find((site) => site.id === id);

      if (foundSite) {
        setActiveSite(foundSite);
      } else if (siteData.getSites.sites.length > 0 && id === '') {
        const firstSite = siteData.getSites.sites[0];
        router.push('/console/sites/' + firstSite.id);
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
  }, [activeSite, setSnackbarMessage, fetchAddress, setSelectedDefaultSite]);

  useEffect(() => {
    if (error) {
      setSnackbarMessage({
        id: 'error-fetching-address',
        type: 'error',
        show: true,
        message: 'Error fetching address from coordinates',
      });
    }
  }, [error, setSnackbarMessage]);

  useEffect(() => {
    if (activeSite.id) {
      setNodesFetched(false);
      fetchNodesForSite({
        variables: {
          data: {
            state: NodeStateEnum.Configured,
            siteId: activeSite.id,
          },
        },
      });
    }
  }, [activeSite.id, fetchNodesForSite]);

  const getInitialNodeUptimes = (): Record<string, number> => {
    if (!statData?.getSiteStat?.metrics || !nodeIds || nodeIds.length === 0) {
      return {};
    }

    const nodeUptimes: Record<string, number> = {};
    statData.getSiteStat.metrics.forEach((metric: any) => {
      if (
        metric.type === SITE_KPI_TYPES.NODE_UPTIME &&
        metric.nodeId &&
        metric.success
      ) {
        nodeUptimes[metric.nodeId] = metric.value;
      }
    });

    return nodeUptimes;
  };

  const initialNodeUptimes = getInitialNodeUptimes();

  if (!isDataReady) {
    return (
      <Grid2 container columnSpacing={2} rowSpacing={2}>
        {[1, 2, 3].map((item) => (
          <Grid2 size={4} key={item}>
            <Skeleton
              height={164}
              width={'100%'}
              variant="rectangular"
              sx={{ borderRadius: '5px' }}
            />
          </Grid2>
        ))}
        <Grid2 size={12}>
          <Skeleton
            height={300}
            width={'100%'}
            variant="rectangular"
            sx={{ borderRadius: '5px' }}
          />
        </Grid2>
      </Grid2>
    );
  }

  return (
    <Box
      sx={{
        overflowY: 'auto',
        overflowX: 'hidden',
        borderRadius: '10px',
        width: '100%',
        height: 'calc(100vh - 164px)',
      }}
    >
      <SiteDetailsHeader
        siteList={siteData?.getSites.sites || []}
        selectedSiteId={activeSite.id}
        onSiteChange={handleSiteChange}
        isLoading={sitesLoading || statLoading}
        siteStatMetrics={statData?.getSiteStat ?? { metrics: [] }}
      />

      <Grid2
        container
        spacing={2}
        alignItems="stretch"
        sx={{
          mt: 1,
          height: 'max-content',
        }}
      >
        <Grid2
          size={{ xs: 12, sm: 6, md: 4 }}
          sx={{ height: 'auto', display: 'flex' }}
        >
          <SiteInfo
            selectedSite={activeSite}
            address={CurrentSiteaddress}
            nodeIds={nodeIds}
          />
        </Grid2>
        <Grid2
          size={{ xs: 12, sm: 6, md: 5 }}
          sx={{ height: '100%', display: 'flex' }}
        >
          <SiteOverview
            installationDate={new Date(activeSite.installDate)}
            isLoading={statLoading}
            siteId={activeSite.id}
            siteStatMetrics={statData?.getSiteStat ?? { metrics: [] }}
          />
        </Grid2>
        <Grid2
          size={{ xs: 12, sm: 6, md: 3 }}
          sx={{ height: 'auto', display: 'flex', minHeight: 200 }}
        >
          <SiteMapComponent
            id="site-map"
            zoom={15}
            posix={[activeSite.latitude ?? '0', activeSite.longitude ?? '0']}
            address={CurrentSiteaddress}
            height={'100%'}
            mapStyle="satellite"
            showUserCount={true}
            userCount={activeSubscribers}
          />
        </Grid2>
        <Grid2 size={12}>
          <SiteComponents
            key={`${activeView.kpi}-${metricFrom}`}
            siteId={activeSite.id}
            metrics={metrics}
            sections={sections}
            activeKPI={activeView.kpi}
            activeSection={getSectionName(activeView.graphType)}
            metricFrom={metricFrom}
            metricsLoading={metricsLoading}
            onComponentClick={handleViewChange}
            onSwitchChange={handleSwitchChange}
            nodeIds={nodeIds}
            initialNodeUptimes={initialNodeUptimes}
          />
        </Grid2>
      </Grid2>
    </Box>
  );
};

export default Page;
