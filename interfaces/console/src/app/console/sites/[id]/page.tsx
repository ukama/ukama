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
  useAddSiteMutation,
  useGetComponentsByUserIdLazyQuery,
  useGetNetworksQuery,
  useGetNodesForSiteLazyQuery,
  useGetSiteLazyQuery,
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
import ConfigureSiteDialog from '@/components/ConfigureSiteDialog';
import SiteComponents from '@/components/SiteComponents';
import { SectionData, STAT_STEP_29 } from '@/constants/index';
import SiteDetailsHeader from '@/components/SiteDetailsHeader';
import SiteInfo from '@/components/SiteInfos';
import SiteOverview from '@/components/SiteOverView';
import { SITE_KPIS } from '@/constants';
import { METRIC_RANGE_10800 } from '@/constants';
import { useAppContext } from '@/context';
import { TMetricResDto, TSiteForm } from '@/types';
import { useFetchAddress } from '@/utils/useFetchAddress';
import { getSectionFromKPI, getUnixTime } from '@/utils';
import { AlertColor, Box, Grid, Skeleton } from '@mui/material';
import { formatISO } from 'date-fns';
import dynamic from 'next/dynamic';
import React, { useEffect, useState } from 'react';
import MetricStatBySiteSubscription from '@/lib/MetricStatBySiteSubscription';
import { useRouter } from 'next/navigation';

const SiteMapComponent = dynamic(
  () => import('@/components/SiteMapComponent'),
  {
    ssr: false,
  },
);

const NODE_UPTIME_KEY = 'unit_uptime';

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
  const [openSiteConfig, setOpenSiteConfig] = useState(false);
  const [componentsList, setComponentsList] = useState<any[]>([]);
  const [selectedSiteId, setSelectedSiteId] = useState<string | null>(null);
  const [sitesList, setSitesList] = useState<SiteDto[]>([]);
  const [nodeIds, setNodeIds] = useState<string[]>([]);
  const [siteUptime, setSiteUptime] = useState<number>(0);
  const [siteUptimePercentage, setSiteUptimePercentage] = useState<number>(0);
  const [nodeUptimes, setNodeUptimes] = useState<Record<string, number>>({});
  const [nodesFetched, setNodesFetched] = useState(false);
  const router = useRouter();
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
  const [isDataReady, setIsDataReady] = useState(false);

  const [metricFrom, setMetricFrom] = useState<number>(0);
  const [graphType, setGraphType] = useState<Graphs_Type>(Graphs_Type.Solar);
  const [metrics, setMetrics] = useState<MetricsRes>({ metrics: [] });
  const [activeSection, setActiveSection] = useState<string>('SOLAR');
  const [activeKPI, setActiveKPI] = useState<string>('solar');
  const sections: SectionData = {
    SOLAR: SITE_KPIS.SOLAR.metrics,
    BATTERY: SITE_KPIS.BATTERY.metrics,
    CONTROLLER: SITE_KPIS.CONTROLLER.metrics,
    MAIN_BACKHAUL: SITE_KPIS.MAIN_BACKHAUL.metrics,
    SWITCH: SITE_KPIS.SWITCH.metrics,
  };

  const handleSiteConfigOpen = () => {
    setOpenSiteConfig(true);
  };

  const handleCloseSiteConfig = () => {
    setOpenSiteConfig(false);
  };

  const [
    getMetricBySite,
    { loading: metricsLoading, variables: metricsVariables },
  ] = useGetMetricBySiteLazyQuery({
    client: subscriptionClient,
    fetchPolicy: 'network-only',
    onCompleted: (data) => {
      setMetrics(data.getMetricBySite);
    },
    onError: (err) => {
      setSnackbarMessage({
        id: 'site-metrics-err-msg',
        message: err.message,
        type: 'error',
        show: true,
      });
    },
  });

  const [updateSwitchPort, { loading: updatePortLoading }] =
    useToggleInternetSwitchMutation({
      onError: (err) => {
        setSnackbarMessage({
          id: 'update-node-err-msg',
          message: err.message,
          type: 'error',
          show: true,
        });
      },
    });

  const handleSwitchChange = async (
    portNumber: number,
    currentStatus: boolean,
  ) => {
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
          message: `Port ${portNumber} status updated successfully to ${newStatus ? 'On' : 'Off'}`,
          type: 'success',
          show: true,
        });

        if (metricFrom > 0) {
          getMetricBySite({
            variables: {
              data: {
                step: 30,
                siteId: id,
                userId: user.id,
                type: graphType,
                from: metricFrom,
                orgName: user.orgName,
                withSubscription: false,
                to: metricFrom + METRIC_RANGE_10800,
              },
            },
          });
        }
      } else {
        setSnackbarMessage({
          id: 'update-switch-error',
          message: `Failed to update Port ${portNumber} status. Please try again.`,
          type: 'error',
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
  };

  const [addSite, { loading: addSiteLoading }] = useAddSiteMutation({
    onCompleted: (res) => {
      setSnackbarMessage({
        id: 'add-site-success',
        message: 'Site added successfully!',
        type: 'success' as AlertColor,
        show: true,
      });
    },
    onError: (error) => {
      setSnackbarMessage({
        id: 'add-subscriber-error',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });

  const { loading: sitesLoading } = useGetSitesQuery({
    skip: !network.id,
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

  const handleSiteConfiguration = async (data: any) => {
    const variables = {
      access_id: data.access,
      backhaul_id: data.backhaul,
      install_date: formatISO(new Date()),
      latitude: data.coordinates.lat,
      location: data.location,
      longitude: data.coordinates.lng,
      name: data.siteName,
      network_id: data.selectedNetwork,
      power_id: data.power,
      spectrum_id: data.spectrumId || '',
      switch_id: data.switch,
      is_deactivated: data.is_deactivated || false,
    };

    try {
      await addSite({ variables: { data: variables } });
    } catch (error) {
      console.error('Error submitting site configuration:', error);
    }
  };

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

  const { data: networks } = useGetNetworksQuery({
    fetchPolicy: 'cache-and-network',
    onError: (error) => {
      setSnackbarMessage({
        id: 'networks-msg',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });

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
        type: 'error',
        show: true,
      });
    },
  });

  const [fetchNodesForSite] = useGetNodesForSiteLazyQuery({
    onCompleted: (res) => {
      if (res.getNodesForSite?.nodes) {
        const ids = res.getNodesForSite.nodes.map((node) => node.id);
        setNodeIds(ids);
        console.log('Node IDs:', ids);
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

  useEffect(() => {
    if (id) {
      setMetricFrom(getUnixTime() - METRIC_RANGE_10800);
    }
  }, [selectedSiteId]);

  useEffect(() => {
    if (metricFrom > 0) {
      getMetricBySite({
        variables: {
          data: {
            step: 30,
            siteId: id,
            userId: user.id,
            type: graphType,
            from: metricFrom,
            orgName: user.orgName,
            withSubscription: false,
            to: metricFrom + METRIC_RANGE_10800,
          },
        },
      });
    }
  }, [metricFrom, graphType, getMetricBySite, id, user.id, user.orgName]);

  const handleSectionChange = (section: string): void => {
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
      setMetricFrom(getUnixTime() - METRIC_RANGE_10800);
    }
  };

  const handleComponentClick = (kpiType: string) => {
    setActiveKPI(kpiType);

    const sectionName = getSectionFromKPI(kpiType);
    handleSectionChange(sectionName);
  };

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
    getSite({ variables: { siteId: id } });
  }, [id, getSite]);

  useEffect(() => {
    if (id) {
      setSelectedSiteId(id);
    } else if (sitesList.length > 0) {
      setSelectedSiteId(sitesList[0].id);
    }
  }, [id, sitesList]);

  const handleSiteChange = (newSiteId: string) => {
    setSelectedSiteId(newSiteId);
    router.push('/console/sites/' + newSiteId);
  };

  useEffect(() => {
    if (selectedSiteId) {
      getSite({ variables: { siteId: selectedSiteId } });
    }
  }, [selectedSiteId, getSite]);

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

  useEffect(() => {
    if (activeSite) {
      setNodesFetched(false);
      fetchNodesForSite({ variables: { siteId: id } });
    }
  }, [activeSite, fetchNodesForSite, id]);

  const [getSiteMetricStat, { variables: statSiteVar }] =
    useGetSiteStatLazyQuery({
      client: subscriptionClient,
      fetchPolicy: 'network-only',
      onCompleted: (data) => {
        if (data.getSiteStat.metrics.length > 0) {
          data.getSiteStat.metrics.forEach((m) => {
            if (m.type === 'site_uptime_seconds') {
              setSiteUptime(m.value);
            }
            if (m.type === 'site_uptime_percentage') {
              setSiteUptimePercentage(m.value);
            }

            if (m.nodeId && m.type === NODE_UPTIME_KEY) {
              setNodeUptimes((prev) => ({
                ...prev,
                [m.nodeId]: m.value,
              }));
            }
          });

          const sKey = `stat-${user.orgName}-${user.id}-${Stats_Type.Site}-${
            statSiteVar?.data.from ?? 0
          }`;

          MetricStatBySiteSubscription({
            key: sKey,
            siteId: id,
            userId: user.id,
            url: env.METRIC_URL,
            orgName: user.orgName,
            type: Stats_Type.Site,
            from: statSiteVar?.data.from ?? 0,
            nodeIds: nodeIds,
          });

          PubSub.subscribe(sKey, handleSiteStatSubscription);
        }
      },
      onError: (err) => {
        console.error('Error fetching site/node metrics:', err);
        setSnackbarMessage({
          id: 'site-metrics-err-msg',
          message: err.message,
          type: 'error',
          show: true,
        });
      },
    });

  const handleMetricSubscription = (_: any, data: string) => {
    try {
      const parsedData: TMetricResDto = JSON.parse(data);
      if (!parsedData?.data?.getSiteMetricStatSub) {
        console.error('Invalid metric subscription data format');
        return;
      }

      const { value, type, success } = parsedData.data.getSiteMetricStatSub;

      if (type === graphType && success) {
        setMetrics((prevMetrics) => ({
          ...prevMetrics,
          metrics: prevMetrics.metrics.map((metric) =>
            metric.type === type
              ? {
                  ...metric,
                  values: [...metric.values, value],
                }
              : metric,
          ),
        }));
      }
    } catch (error) {
      console.error('Error handling metric subscription data:', error);
    }
  };

  useEffect(() => {
    const subscriptionKey = `metric-${user.id}-${id}-${graphType}-${metricFrom}`;
    PubSub.subscribe(subscriptionKey, handleMetricSubscription);

    return () => {
      PubSub.unsubscribe(subscriptionKey);
    };
  }, [id, user.id, graphType, metricFrom]);

  const handleSiteStatSubscription = (_: any, data: string) => {
    try {
      const parsedData: TMetricResDto = JSON.parse(data);

      if (!parsedData?.data?.getSiteMetricStatSub) {
        console.error('Invalid site subscription data format');
        return;
      }

      const { value, type, success, nodeId } =
        parsedData.data.getSiteMetricStatSub;
      if (success) {
        if (type === 'site_uptime_seconds') {
          setSiteUptime(() => Math.floor(value[1]));
        }
        if (type === 'site_uptime_percentage') {
          setSiteUptimePercentage(() => Math.floor(value[1]));
        }

        if (nodeId && type === NODE_UPTIME_KEY) {
          setNodeUptimes((prev) => ({
            ...prev,
            [nodeId]: Math.floor(value[1]),
          }));
        }

        if (nodeId) {
          PubSub.publish(`stat-${type}-${nodeId}`, value);
        } else {
          PubSub.publish(`stat-${type}`, value);
        }
      }
    } catch (error) {
      console.error('Error handling site subscription data:', error);
    }
  };

  useEffect(() => {
    if (!nodesFetched) {
      return;
    }

    const to = getUnixTime();
    const from = to - STAT_STEP_29;
    if (!id) {
      setSnackbarMessage({
        id: 'site-not-found-msg',
        message: 'Site not found.',
        type: 'error',
        show: true,
      });
    } else if (id) {
      getSiteMetricStat({
        variables: {
          data: {
            to: to,
            siteId: id,
            from: from,
            userId: user.id,
            step: STAT_STEP_29,
            orgName: user.orgName,
            withSubscription: true,
            type: Stats_Type.Site,
            nodeIds: nodeIds,
          },
        },
      });
    }

    return () => {
      const sKey = `stat-${user.orgName}-${user.id}-${Stats_Type.Site}-${
        from ?? 0
      }`;
      PubSub.unsubscribe(sKey);
    };
  }, [
    id,
    user.id,
    user.orgName,
    getSiteMetricStat,
    setSnackbarMessage,
    nodeIds,
    nodesFetched,
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
        addSite={handleSiteConfigOpen}
        siteList={sitesList || []}
        selectedSiteId={selectedSiteId}
        onSiteChange={handleSiteChange}
        isLoading={sitesLoading}
        siteUpTime={siteUptime}
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
            siteUptimeSeconds={siteUptime}
            uptimePercentage={siteUptimePercentage}
            installationDate={new Date(activeSite.installDate)}
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
          nodeUptimes={nodeUptimes}
        />
      </Box>

      <ConfigureSiteDialog
        site={site}
        open={openSiteConfig}
        addSiteLoading={addSiteLoading}
        onClose={handleCloseSiteConfig}
        components={componentsList || []}
        networks={networks?.getNetworks?.networks || []}
        handleSiteConfiguration={handleSiteConfiguration}
      />
    </Box>
  );
};

export default Page;
