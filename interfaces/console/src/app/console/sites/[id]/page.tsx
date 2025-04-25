/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';

import {
  SiteDto,
  useGetNodesForSiteLazyQuery,
  useGetSitesQuery,
  useToggleInternetSwitchMutation,
} from '@/client/graphql/generated';
import {
  Graphs_Type,
  MetricsRes,
  Stats_Type,
  useGetMetricBySiteLazyQuery,
  useGetSiteStatLazyQuery,
} from '@/client/graphql/generated/subscriptions';
import SiteComponents from '@/components/SiteComponents';
import { SectionData, STAT_STEP_29 } from '@/constants/index';
import SiteDetailsHeader from '@/components/SiteDetailsHeader';
import SiteInfo from '@/components/SiteInfos';
import SiteOverview from '@/components/SiteOverView';
import { SITE_KPIS } from '@/constants';
import { METRIC_RANGE_10800 } from '@/constants';
import { useAppContext } from '@/context';
import { TMetricResDto } from '@/types';
import { useFetchAddress } from '@/utils/useFetchAddress';
import { getSectionFromKPI, getUnixTime } from '@/utils';
import { AlertColor, Box, Grid, Skeleton } from '@mui/material';
import dynamic from 'next/dynamic';
import React, {
  useEffect,
  useState,
  useCallback,
  useRef,
  useMemo,
} from 'react';
import MetricStatBySiteSubscription from '@/lib/MetricStatBySiteSubscription';
import { useRouter } from 'next/navigation';
import PubSub from 'pubsub-js';
import { SITE_KPI_TYPES } from '@/constants';

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
  const subscriptionsRef = useRef<Record<string, boolean>>({});

  const [activeSite, setActiveSite] = useState<SiteDto>(defaultSite);
  const [selectedSiteId, setSelectedSiteId] = useState<string | null>(null);
  const [sitesList, setSitesList] = useState<SiteDto[]>([]);
  const [nodeIds, setNodeIds] = useState<string[]>([]);
  const [nodesFetched, setNodesFetched] = useState(false);
  const [isDataReady, setIsDataReady] = useState(false);
  const [metricFrom, setMetricFrom] = useState<number>(0);
  const [graphType, setGraphType] = useState<Graphs_Type>(Graphs_Type.Solar);

  const [metrics, setMetrics] = useState<MetricsRes>({ metrics: [] });
  const [activeSection, setActiveSection] = useState<string>('SOLAR');
  const [activeKPI, setActiveKPI] = useState<string>('solar');

  const {
    setSnackbarMessage,
    network,
    setSelectedDefaultSite,
    user,
    env,
    subscriptionClient,
  } = useAppContext();

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
    if (!metricFrom || !id) return;

    const topic = `${user.id}/${graphType}/${metricFrom}`;
    subscriptionsRef.current[topic] = true;

    getMetricBySite({
      variables: {
        data: {
          step: 30,
          siteId: id,
          userId: user.id,
          type: graphType,
          from: metricFrom,
          orgName: user.orgName,
          withSubscription: true,
          to: metricFrom + METRIC_RANGE_10800,
        },
      },
    });

    return () => {
      delete subscriptionsRef.current[topic];
    };
  }, [user.id, graphType, metricFrom, id]);
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

  const { loading: sitesLoading } = useGetSitesQuery({
    skip: !network.id,
    nextFetchPolicy: 'network-only',
    fetchPolicy: 'network-only',
    variables: {
      data: { networkId: network.id },
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

  const [fetchNodesForSite] = useGetNodesForSiteLazyQuery({
    onCompleted: (res) => {
      if (res.getNodesForSite?.nodes) {
        const ids = res.getNodesForSite.nodes.map((node) => node.id);
        setNodeIds(ids);
      } else {
        setNodeIds([]);
      }
      setNodesFetched(true);
    },
    onError: (error) => {
      setSnackbarMessage({
        id: 'nodes-msg',
        message: error.message,
        type: 'error',
        show: true,
      });
      setNodesFetched(true);
    },
  });

  const handleSectionChange = useCallback((section: string): void => {
    setActiveSection(section);
    setMetrics({ metrics: [] });
    if (section !== 'NODE') {
      let newGraphType: Graphs_Type;
      switch (section) {
        case 'SOLAR':
          newGraphType = Graphs_Type.Solar;
          break;
        case 'BATTERY':
          newGraphType = Graphs_Type.Battery;
          break;
        case 'CONTROLLER':
          newGraphType = Graphs_Type.Controller;
          break;
        case 'MAIN_BACKHAUL':
          newGraphType = Graphs_Type.MainBackhaul;
          break;
        case 'SWITCH':
          newGraphType = Graphs_Type.Switch;
          break;
        default:
          newGraphType = Graphs_Type.Solar;
      }

      setGraphType(newGraphType);
      setMetricFrom(() => getUnixTime() - METRIC_RANGE_10800);
    }
  }, []);

  const handleComponentClick = useCallback(
    (kpiType: string) => {
      setActiveKPI(kpiType);
      const sectionName = getSectionFromKPI(kpiType);
      handleSectionChange(sectionName);
      setMetricFrom(() => getUnixTime() - METRIC_RANGE_10800);
    },
    [handleSectionChange],
  );
  const checkDataReadiness = useCallback(() => {
    if (activeSite.id && CurrentSiteaddress && !CurrentSiteAddressLoading) {
      setIsDataReady(true);
    }
  }, [activeSite.id, CurrentSiteaddress, CurrentSiteAddressLoading]);
  const filterActiveSite = useCallback(
    (siteId: string) => {
      const foundSite = sitesList.find((site) => site.id === siteId);
      if (foundSite) {
        setActiveSite(foundSite);
        checkDataReadiness();
      } else {
        setSnackbarMessage({
          id: 'site-not-found',
          message: 'Site not found in available sites',
          type: 'error',
          show: true,
        });
      }
    },
    [sitesList, checkDataReadiness, setSnackbarMessage],
  );
  useEffect(() => {
    cleanupSubscriptions();
    filterActiveSite(id);
    setMetrics({ metrics: [] });
    if (id) {
      setMetricFrom(getUnixTime() - METRIC_RANGE_10800);
      setActiveSection('SOLAR');
      setActiveKPI('solar');
      setGraphType(Graphs_Type.Solar);
    }

    return () => {
      cleanupSubscriptions();
    };
  }, [id, filterActiveSite, cleanupSubscriptions]);

  useEffect(() => {
    if (id) {
      setSelectedSiteId(id);
    } else if (sitesList.length > 0) {
      setSelectedSiteId(sitesList[0].id);
    }
  }, [id, sitesList]);

  const handleSiteChange = useCallback(
    (newSiteId: string) => {
      setSelectedSiteId(newSiteId);
      router.push('/console/sites/' + newSiteId);
    },
    [router],
  );

  useEffect(() => {
    if (selectedSiteId) {
      filterActiveSite(selectedSiteId);
    }
  }, [selectedSiteId, filterActiveSite]);

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

    if (id && activeSite.latitude && activeSite.longitude) {
      handleFetchAddress();
    }
  }, [
    activeSite,
    setSnackbarMessage,
    fetchAddress,
    setSelectedDefaultSite,
    id,
  ]);

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
    if (activeSite) {
      setNodesFetched(false);
      fetchNodesForSite({ variables: { siteId: id } });
    }
  }, [activeSite, fetchNodesForSite, id]);

  const [getSiteMetricStat, { loading: statLoading, variables: statVar }] =
    useGetSiteStatLazyQuery({
      client: subscriptionClient,
      fetchPolicy: 'network-only',
      onCompleted: (data) => {
        const sKey = `stat-${user.orgName}-${user.id}-${Stats_Type.Site}-${
          statVar?.data.from ?? 0
        }`;

        if (data.getSiteStat.metrics.length > 0) {
          subscriptionsRef.current[sKey] = true;

          MetricStatBySiteSubscription({
            key: sKey,
            nodeIds: nodeIds,
            siteIds: [id],
            userId: user.id,
            url: env.METRIC_URL,
            orgName: user.orgName,
            type: Stats_Type.Site,
            from: statVar?.data.from ?? 0,
          });

          PubSub.subscribe(sKey, handleSiteStatSubscription);
        }
      },
    });

  useEffect(() => {
    if (!nodesFetched) return;

    if (id) {
      const to = getUnixTime();
      const from = to - STAT_STEP_29;
      const sKey = `stat-${user.orgName}-${user.id}-${Stats_Type.Site}-${from}`;

      Object.keys(subscriptionsRef.current).forEach((topic) => {
        if (
          topic.startsWith(`stat-${user.orgName}-${user.id}-${Stats_Type.Site}`)
        ) {
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
            userId: user.id,
            step: STAT_STEP_29,
            orgName: user.orgName,
            withSubscription: true,
            type: Stats_Type.Site,
            siteIds: [id],
            nodeIds: nodeIds,
          },
        },
      });

      MetricStatBySiteSubscription({
        key: sKey,
        siteIds: [id],
        userId: user.id,
        url: env.METRIC_URL,
        orgName: user.orgName,
        type: Stats_Type.Site,
        from,
        nodeIds: nodeIds,
      });

      return () => {
        PubSub.unsubscribe(sKey);
        delete subscriptionsRef.current[sKey];
      };
    }
  }, [
    id,
    nodeIds,
    nodesFetched,
    user.id,
    user.orgName,
    env.METRIC_URL,
    getSiteMetricStat,
  ]);

  useEffect(() => {
    if (metricFrom > 0 && id) {
      getMetricBySite({
        variables: {
          data: {
            step: 30,
            siteId: id,
            userId: user.id,
            type: graphType,
            from: metricFrom,
            orgName: user.orgName,
            withSubscription: true,
            to: metricFrom + METRIC_RANGE_10800,
          },
        },
      });
    }
  }, [metricFrom, getMetricBySite, graphType, id, user.id, user.orgName]);

  const handleSiteStatSubscription = useCallback(
    (_: any, data: string) => {
      try {
        const parsedData: TMetricResDto = JSON.parse(data);
        if (parsedData?.data?.getSiteMetricStatSub) {
          const metric = parsedData.data.getSiteMetricStatSub;
          const { type, success, siteId, nodeId, value } = metric;

          if (success && siteId === id) {
            if (!nodeId) {
              if (type === SITE_KPI_TYPES.SITE_UPTIME) {
                PubSub.publish(`stat-site-uptime`, Math.floor(value[1]));
              } else if (type === SITE_KPI_TYPES.SITE_UPTIME_PERCENTAGE) {
                PubSub.publish(
                  `stat-site-uptime-percentage`,
                  Math.floor(value[1]),
                );
              }
            }

            if (nodeId && type === 'unit_uptime') {
              PubSub.publish(
                `stat-node-uptime-${nodeId}`,
                Math.floor(value[1] || 1),
              );
            }
          }

          if (type && type !== 'unit_uptime') {
            PubSub.publish(`stat-${type}`, value);
          }
        }
      } catch (error) {
        console.error('Error in handleSiteStatSubscription:', error);
      }
    },
    [id],
  );

  useEffect(() => {
    return () => {
      cleanupSubscriptions();
    };
  }, [cleanupSubscriptions]);

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
        siteList={sitesList || []}
        selectedSiteId={selectedSiteId}
        onSiteChange={handleSiteChange}
        isLoading={sitesLoading || statLoading}
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
          key={`${activeKPI}-${metricFrom}`}
          siteId={selectedSiteId || ''}
          metrics={metrics}
          sections={sections}
          activeKPI={activeKPI}
          activeSection={activeSection}
          metricFrom={metricFrom}
          metricsLoading={metricsLoading}
          onComponentClick={handleComponentClick}
          onSwitchChange={handleSwitchChange}
        />
      </Box>
    </Box>
  );
};

export default Page;
