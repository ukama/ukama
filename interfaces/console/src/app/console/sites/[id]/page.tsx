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
import { SectionData } from '@/constants/index';
import SiteDetailsHeader from '@/components/SiteDetailsHeader';
import SiteInfo from '@/components/SiteInfos';
import SiteOverview from '@/components/SiteOverView';
import { SITE_KPIS } from '@/constants';
import { useAppContext } from '@/context';
import { ActiveView, KPIType } from '@/types';
import { useFetchAddress } from '@/utils/useFetchAddress';
import { graphTypeToSection, kpiToGraphType } from '@/utils';
import { AlertColor, Box, Grid, Skeleton } from '@mui/material';
import dynamic from 'next/dynamic';
import React, { useEffect, useState, useCallback, useMemo } from 'react';
import { useRouter } from 'next/navigation';
import { useMetricSubscriptions } from '@/utils/useMetricSubscriptions';

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
  latitude: 0,
  location: '',
  longitude: 0,
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

const Page: React.FC<SiteDetailsProps> = ({ params }) => {
  const { id } = params;
  const router = useRouter();
  const [activeSite, setActiveSite] = useState<SiteDto>(defaultSite);
  const [nodeIds, setNodeIds] = useState<string[]>([]);
  const [nodesFetched, setNodesFetched] = useState(false);
  const [isDataReady, setIsDataReady] = useState(false);
  const [activeView, setActiveView] = useState<ActiveView>({
    graphType: Graphs_Type.Solar,
    kpi: 'solar',
  });

  const {
    setSnackbarMessage,
    network,
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
  useEffect(() => {
    return () => {
      cleanupSubscriptions();
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
    skip: !network.id,
    nextFetchPolicy: 'network-only',
    fetchPolicy: 'network-only',
    variables: {
      data: { networkId: network.id },
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
    fetchPolicy: 'network-only',
    onCompleted: (data) => {
      const nodeIds = data.getNodes.nodes
        .filter(
          (node) =>
            node.latitude !== 0 &&
            node.longitude !== 0 &&
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
        graphType: Graphs_Type.Solar,
        kpi: 'solar',
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
        setSnackbarMessage({
          id: 'fetching-address',
          type: 'info',
          show: true,
          message: 'Fetching address with coordinates',
        });
        await fetchAddress(activeSite.latitude, activeSite.longitude);
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

  if (!isDataReady) {
    return (
      <Grid container columnSpacing={2}>
        {[1, 2].map((item) => (
          <Grid item xs={6} key={item}>
            <Skeleton
              variant="rectangular"
              height={158}
              width={'100%'}
              sx={{ borderRadius: '5px' }}
            />
          </Grid>
        ))}
      </Grid>
    );
  }

  return (
    <Box
      sx={{
        overflowY: 'auto',
        overflowX: 'hidden',
        borderRadius: '10px',
        width: '100%',
      }}
    >
      <SiteDetailsHeader
        siteList={siteData?.getSites.sites || []}
        selectedSiteId={activeSite.id}
        onSiteChange={handleSiteChange}
        isLoading={sitesLoading || statLoading}
        siteStatMetrics={statData?.getSiteStat ?? { metrics: [] }}
      />

      <Grid
        container
        spacing={2}
        sx={{
          mt: 1,
          height: 'calc(50vh - 50px)',
        }}
      >
        <Grid item sx={{ height: '100%' }} xs={12} sm={6} md={4}>
          <SiteInfo
            selectedSite={activeSite}
            address={CurrentSiteaddress}
            nodeIds={nodeIds}
          />
        </Grid>
        <Grid item sx={{ height: '100%' }} xs={12} sm={6} md={5}>
          <SiteOverview
            installationDate={new Date(activeSite.installDate)}
            isLoading={statLoading}
            siteId={activeSite.id}
            siteStatMetrics={statData?.getSiteStat ?? { metrics: [] }}
          />
        </Grid>
        <Grid item sx={{ height: '100%' }} xs={12} sm={6} md={3}>
          <SiteMapComponent
            posix={[activeSite.latitude, activeSite.longitude]}
            address={CurrentSiteaddress}
            height={'100%'}
            mapStyle="satellite"
            showUserCount={true}
            userCount={0}
          />
        </Grid>
      </Grid>

      <Box sx={{ mt: 4, mb: 4 }}>
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
        />
      </Box>
    </Box>
  );
};
export default Page;
