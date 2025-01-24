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
import { getSiteTabTypeByIndex, getUnixTime } from '@/utils';
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
import React, { useEffect, useState, Suspense } from 'react';
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
  const [metricFrom, setMetricFrom] = useState<number>(0);
  const [graphType, setGraphType] = useState<Graphs_Type>(Graphs_Type.Power);
  const [selectedSiteId, setSelectedSiteId] = useState<string | null>(null);
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
  const [isDataReady, setIsDataReady] = useState(false);

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
  const [
    getSiteMetricByTab,
    { loading: siteMetricsLoading, variables: siteMetricsVariables },
  ] = useGetMetricByTabLazyQuery({
    client: subscriptionClient,
    fetchPolicy: 'network-only',
    onCompleted: (data) => {
      setMetrics(data.getMetricByTab);
    },
  });

  useEffect(() => {
    if (metricFrom > 0 && siteMetricsVariables?.data?.from !== metricFrom) {
      const psKey = `metric-${user.orgName}-${user.id}-${graphType}-${metricFrom}`;
      getSiteMetricByTab({
        variables: {
          data: {
            nodeId: id,
            userId: user.id,
            type: graphType,
            from: metricFrom,
            to: metricFrom + 120,
            orgName: user.orgName,
            withSubscription: true,
          },
        },
      }).then(() => {
        MetricSubscription({
          nodeId: id,
          key: psKey,
          type: graphType,
          userId: user.id,
          from: metricFrom,
          url: env.METRIC_URL,
          orgName: user.orgName,
        });
      });

      PubSub.subscribe(psKey, handleNotification);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [
    metricFrom,
    siteMetricsVariables?.data?.from,
    getSiteMetricByTab,
    graphType,
  ]); // Added all missing dependencies
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
  const handleNotification = (_: any, data: string) => {
    const parsedData: TMetricResDto = JSON.parse(data);
    const { msg, type, value, nodeId, success } =
      parsedData.data.getMetricByTabSub;
    if (success) {
      PubSub.publish(type, [Math.floor(value[0] ?? 0) * 1000, value[1]]);
      // setMetrics((prev) => {
      //   const updatedMetrics = prev.metrics.map((m) =>
      //     m.type === type ? { ...m, values: [...m.values, value] } : m,
      //   );
      //   return { ...prev, metrics: updatedMetrics };
      // });
    }
  };

  const [restartSite, { loading: restartSiteLoading }] = useRestartSiteMutation(
    {
      onCompleted: (data) => {
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
  useEffect(() => {
    const value = 0; // Define the value variable
    setGraphType(getSiteTabTypeByIndex(value) ?? Graphs_Type.Power);
    setMetricFrom(() => getUnixTime() - 120);
  }, [metrics]);
  const handleOverviewSectionChange = (type: Graphs_Type) => {
    setGraphType(type);
    setMetricFrom(() => getUnixTime() - 120);
  };

  const [getSite, { loading: getSiteLoading }] = useGetSiteLazyQuery({
    onCompleted: (res) => {
      setSite({
        power: res.getSite.powerId,
        siteName: res.getSite.name,
        switch: res.getSite.switchId,
        access: res.getSite.accessId,
        address: res.getSite.location,
        latitude: res.getSite.latitude,
        network: res.getSite.networkId,
        spectrum: res.getSite.spectrumId,
        backhaul: res.getSite.backhaulId,
        longitude: res.getSite.longitude,
      });
      setActiveSite({
        id: res.getSite.id,
        accessId: res.getSite.accessId,
        backhaulId: res.getSite.backhaulId,
        createdAt: res.getSite.createdAt,
        installDate: res.getSite.installDate,
        isDeactivated: res.getSite.isDeactivated,
        latitude: res.getSite.latitude,
        location: res.getSite.location,
        longitude: res.getSite.longitude,
        name: res.getSite.name,
        networkId: res.getSite.networkId,
        powerId: res.getSite.powerId,
        spectrumId: res.getSite.spectrumId,
        switchId: res.getSite.switchId,
      });
      checkDataReadiness();
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
  const [fetchNode, { data: nodeData, loading: nodeLoading }] =
    useGetNodesByNetworkLazyQuery();

  const { data: subscribers } = useGetSubscribersByNetworkQuery({
    variables: {
      networkId: activeSite.networkId,
    },
    fetchPolicy: 'network-only',
    nextFetchPolicy: 'network-only',
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
    if (activeSite?.id) {
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
  }, []);

  useEffect(() => {
    getSite({ variables: { siteId: id } });
  }, []);

  useEffect(() => {
    if (id) {
      setSelectedSiteId(id);
    } else if (sitesList.length > 0) {
      setSelectedSiteId(sitesList[0].id);
    }
  }, [id, sitesList]);

  const handleSiteChange = (newSiteId: string) => {
    setSelectedSiteId(newSiteId);
    router.push(`/console/sites/${newSiteId}`);
  };

  useEffect(() => {
    if (selectedSiteId) {
      getSite({ variables: { siteId: selectedSiteId } });
    }
  }, [selectedSiteId]);

  useEffect(() => {
    const handleFetchAddress = async () => {
      setSnackbarMessage({
        id: 'fetching-address',
        type: 'info',
        show: true,
        message: 'Fetching address with coordinates',
      });
      await fetchAddress(activeSite.latitude, activeSite.longitude);
    };

    setSelectedDefaultSite(activeSite.name);

    if (activeSite && activeSite.latitude && activeSite.longitude) {
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

  const checkDataReadiness = () => {
    if (
      activeSite.id &&
      CurrentSiteaddress &&
      !getSiteLoading &&
      !CurrentSiteAddressLoading
    ) {
      setIsDataReady(true);
    }
  };

  useEffect(() => {
    checkDataReadiness();
  }, [
    activeSite,
    CurrentSiteaddress,
    getSiteLoading,
    CurrentSiteAddressLoading,
  ]);

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
  const handleSiteRestart = () => {
    setSnackbarMessage({
      id: 'site-restart',
      type: 'info',
      show: true,
      message: 'Restarting site',
    });
    restartSite({
      variables: {
        data: {
          siteId: activeSite.id,
          networkId: activeSite.networkId,
        },
      },
    });
  };

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
            <SiteOverview
              inputPower="120W"
              solarStorage="80%"
              consumption="40W"
            />
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
        <Grid
          item
          xs={12}
          md={12}
          sx={{
            height: 'auto',
          }}
        >
          <Paper
            elevation={3}
            sx={{
              p: 4,
              height: 'auto',
              display: 'flex',
              flexDirection: 'column',
            }}
          >
            <SiteOverallHealth
              batteryInfo={[]}
              solarHealth={'good'}
              nodeHealth={'good'}
              switchHealth={'good'}
              controllerHealth={'good'}
              batteryHealth={'good'}
              backhaulHealth={'good'}
              nodes={nodeData?.getNodesByNetwork.nodes || []}
              nodeId={activeSite.id}
              metricFrom={0}
              topic="Site Overall Health"
              title="Site Overall Health"
              initData={[]}
              hasData={false}
              loading={false}
              tabSection={Graphs_Type.Power}
            />
          </Paper>
        </Grid>
      </Grid>
    </Box>
  );
};

export default Page;
