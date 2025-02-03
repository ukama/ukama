/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';

import {
  Component_Type,
  SiteDto,
  useGetComponentsByUserIdLazyQuery,
  useGetSiteLazyQuery,
  useGetSitesQuery,
  useGetNodesByNetworkLazyQuery,
  useGetSubscribersByNetworkQuery,
  useRestartSiteMutation,
} from '@/client/graphql/generated';
import {
  Graphs_Type,
  MetricsRes,
  useGetMetricByTabLazyQuery,
} from '@/client/graphql/generated/subscriptions';
import SiteDetailsHeader from '@/components/SiteDetailsHeader';
import SiteOverallHealth from '@/components/SiteHealth';
import SiteInfo from '@/components/SiteInfos';
import SiteOverview from '@/components/SiteOverView';
import { useAppContext } from '@/context';
import colors from '@/theme/colors';
import { TMetricResDto, TSiteForm } from '@/types';
import { useFetchAddress } from '@/utils/useFetchAddress';
import GroupIcon from '@mui/icons-material/Group';
import { getUnixTime } from '@/utils';
import {
  AlertColor,
  Box,
  Grid,
  Paper,
  Skeleton,
  Typography,
} from '@mui/material';
import dynamic from 'next/dynamic';
import { useRouter } from 'next/navigation';
import React, { useEffect, useState, Suspense, useMemo } from 'react';
import MetricSubscription from '@/lib/MetricSubscription';

const SiteMapComponent = dynamic(
  () => import('@/components/SiteMapComponent'),
  {
    ssr: false,
    loading: () => (
      <Skeleton
        variant="rectangular"
        width="100%"
        height="100%"
        sx={{ borderRadius: '5px' }}
      />
    ),
  },
);

const SITE_INIT = {
  switch: '',
  power: '',
  access: '',
  backhaul: '',
  address: '',
  spectrum: '',
  siteName: '',
  latitude: NaN,
  longitude: NaN,
  network: '',
};

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
  const [site, setSite] = useState<TSiteForm>(SITE_INIT);
  const [activeSite, setActiveSite] = useState<SiteDto>(defaultSite);
  const [componentsList, setComponentsList] = useState<any[]>([]);
  const [metrics, setMetrics] = useState<MetricsRes>({ metrics: [] });
  const [metricFrom, setMetricFrom] = useState<number>(getUnixTime() - 140);
  const [graphType, setGraphType] = useState<Graphs_Type>(Graphs_Type.Battery);
  const [selectedSiteId, setSelectedSiteId] = useState<string | null>(id);
  const [sitesList, setSitesList] = useState<SiteDto[]>([]);
  const {
    user,
    setSnackbarMessage,
    network,
    setSelectedDefaultSite,
    env,
    subscriptionClient,
  } = useAppContext();
  const {
    address: CurrentSiteaddress,
    isLoading: CurrentSiteAddressLoading,
    error,
    fetchAddress,
  } = useFetchAddress();

  const router = useRouter();

  const { loading: sitesLoading } = useGetSitesQuery({
    skip: !network.id,
    variables: {
      networkId: network.id,
    },
    onCompleted: (res) => {
      setSitesList(res.getSites.sites);
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
  const [getSite, { loading: getSiteLoading }] = useGetSiteLazyQuery({
    onCompleted: (res) => {
      const siteData = res.getSite;
      setSite({
        power: siteData.powerId,
        siteName: siteData.name,
        switch: siteData.switchId,
        access: siteData.accessId,
        address: siteData.location,
        latitude: siteData.latitude,
        network: siteData.networkId,
        spectrum: siteData.spectrumId,
        backhaul: siteData.backhaulId,
        longitude: siteData.longitude,
      });
      setActiveSite(siteData);
    },
    onError: (error) => {
      setSnackbarMessage({
        id: 'sites-msg',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });

  const [getComponents] = useGetComponentsByUserIdLazyQuery({
    onCompleted: (res) => {
      if (res.getComponentsByUserId) {
        setComponentsList(res.getComponentsByUserId.components);
      }
    },
    onError: (error) => {
      setSnackbarMessage({
        id: 'components-msg',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });

  const [restartSite, { loading: restartSiteLoading }] = useRestartSiteMutation(
    {
      onCompleted: () => {
        setSnackbarMessage({
          id: 'restart-site-success',
          message: 'Site received restart command!',
          type: 'success' as AlertColor,
          show: true,
        });
      },
      onError: (error) => {
        setSnackbarMessage({
          id: 'restart-site-error',
          message: error.message,
          type: 'error' as AlertColor,
          show: true,
        });
      },
    },
  );

  const [fetchNode, { data: nodeData, loading: nodeLoading }] =
    useGetNodesByNetworkLazyQuery({});

  const { data: subscribers } = useGetSubscribersByNetworkQuery({
    variables: {
      networkId: activeSite.networkId,
    },
    fetchPolicy: 'network-only',
    onError: (error) => {
      setSnackbarMessage({
        id: 'subscriber-msg',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });

  useEffect(() => {
    if (activeSite?.networkId) {
      fetchNode({ variables: { networkId: activeSite.networkId } });
    }
  }, [activeSite, fetchNode]);

  useEffect(() => {
    getComponents({
      variables: {
        data: {
          category: Component_Type.All,
        },
      },
    });
  }, [getComponents]);

  useEffect(() => {
    if (activeSite.latitude && activeSite.longitude) {
      fetchAddress(activeSite.latitude, activeSite.longitude);
    }
    setSelectedDefaultSite(activeSite.name);
  }, [activeSite, fetchAddress, setSelectedDefaultSite]);

  const handleSiteChange = (newSiteId: string) => {
    setSelectedSiteId(newSiteId);
    router.push(`/console/sites/${newSiteId}`);
  };

  const handleSiteKpiChange = (type: Graphs_Type) => {
    setGraphType(type);
    setMetricFrom(getUnixTime() - 140);
  };

  const handleSiteRestart = () => {
    restartSite({
      variables: {
        data: {
          siteId: activeSite.id,
          networkId: activeSite.networkId,
        },
      },
    });
  };
  const [
    getSiteMetricByTab,
    { loading: siteMetricsLoading, variables: siteMetricsVariables },
  ] = useGetMetricByTabLazyQuery({
    client: subscriptionClient,
    fetchPolicy: 'network-only',
    onCompleted: (data) => {
      setMetrics(data.getMetricByTab);
      console.log('SiteOverallHealth metrics:', data);
    },
  });

  const handleNotification = (_: any, data: string) => {
    const parsedData: TMetricResDto = JSON.parse(data);
    const { msg, type, value, nodeId, success } =
      parsedData.data.getMetricByTabSub;
    if (success) {
      PubSub.publish(type, [Math.floor(value[0] ?? 0) * 1000, value[1]]);
    }
  };
  useEffect(() => {
    if (
      metricFrom > 0 &&
      activeSite.id &&
      siteMetricsVariables?.data?.from !== metricFrom
    ) {
      const psKey = `metric-${user.orgName}-${user.id}-${graphType}-${metricFrom}-${activeSite.id}`;
      getSiteMetricByTab({
        variables: {
          data: {
            nodeId: '',
            siteId: activeSite.id,
            userId: user.id,
            type: graphType,
            from: metricFrom,
            orgName: user.orgName,
            withSubscription: true,
          },
        },
      }).then(() => {
        MetricSubscription({
          nodeId: '',
          siteId: activeSite.id,
          key: psKey,
          type: graphType,
          userId: user.id,
          from: metricFrom,
          url: env.METRIC_URL,
          orgName: user.orgName,
        });
      });

      PubSub.subscribe(psKey, handleNotification);

      return () => {
        PubSub.unsubscribe(psKey);
      };
    }
  }, [
    metricFrom,
    graphType,
    activeSite.id,
    siteMetricsVariables?.data?.from,
    getSiteMetricByTab,
  ]);

  const backhaulComponent = useMemo(() => {
    return (
      componentsList.find(
        (component) =>
          component.id === activeSite?.backhaulId &&
          component.category === 'BACKHAUL',
      )?.type || ''
    );
  }, [componentsList, activeSite?.backhaulId]);

  useEffect(() => {
    if (id) {
      getSite({ variables: { siteId: id } });

      getSiteMetricByTab({
        variables: {
          data: {
            nodeId: '',
            siteId: id,
            userId: user.id,
            type: Graphs_Type.Battery,
            from: getUnixTime() - 140,
            orgName: user.orgName,
            withSubscription: true,
          },
        },
      });
    }
  }, [id, getSite, getSiteMetricByTab, user.id, user.orgName]);

  return (
    <Box>
      <SiteDetailsHeader
        siteList={sitesList || []}
        selectedSiteId={selectedSiteId}
        onSiteChange={handleSiteChange}
        isLoading={sitesLoading}
        onRestartSite={handleSiteRestart}
      />
      <Grid container spacing={2}>
        <Grid item xs={12} md={3}>
          <Paper sx={{ height: '250px', overflow: 'auto' }}>
            <SiteInfo
              selectedSite={activeSite}
              address={CurrentSiteaddress}
              nodes={nodeData?.getNodesByNetwork.nodes || []}
            />
          </Paper>
        </Grid>

        <Grid item xs={12} md={6}>
          <Paper sx={{ height: '250px', overflow: 'auto' }}>
            <SiteOverview metrics={metrics} loading={siteMetricsLoading} />
          </Paper>
        </Grid>
        <Grid item xs={12} md={3}>
          <Paper
            sx={{ height: '250px', overflow: 'hidden', position: 'relative' }}
          >
            <Box
              sx={{
                position: 'absolute',
                top: 206,
                left: 8,
                zIndex: 2,
                display: 'flex',
                alignItems: 'center',
                gap: 1,
                backgroundColor: colors.white,
                borderRadius: '4px',
                p: 1,
              }}
            >
              <GroupIcon fontSize="small" />
              <Typography variant="body2" fontWeight="medium">
                {subscribers?.getSubscribersByNetwork.subscribers.length || 0}
              </Typography>
            </Box>

            <Box sx={{ position: 'relative', zIndex: 1, height: '100%' }}>
              <Suspense
                fallback={
                  <Skeleton
                    variant="rectangular"
                    width="100%"
                    height="100%"
                    sx={{ borderRadius: '5px' }}
                  />
                }
              >
                <SiteMapComponent
                  posix={[activeSite.latitude, activeSite.longitude]}
                  address={CurrentSiteaddress}
                  height={'100%'}
                  mapStyle="satellite"
                />
              </Suspense>
            </Box>
          </Paper>
        </Grid>
        <Grid item xs={12} md={12}>
          <Paper
            elevation={3}
            sx={{
              p: 4,
              height: 'auto',
            }}
          >
            <SiteOverallHealth
              siteId={activeSite.id}
              metricFrom={metricFrom}
              metrics={metrics}
              loading={siteMetricsLoading}
              onGraphTypeChange={handleSiteKpiChange}
              nodes={nodeData?.getNodesByNetwork}
              backhaulComponent={backhaulComponent}
            />
          </Paper>
        </Grid>
      </Grid>
    </Box>
  );
};

export default Page;
