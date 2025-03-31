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
  useGetSitesQuery,
  useUpdateSiteMutation,
} from '@/client/graphql/generated';
import {
  Stats_Type,
  useGetSiteStatLazyQuery,
} from '@/client/graphql/generated/subscriptions';
import ConfigureSiteDialog from '@/components/ConfigureSiteDialog';
import EditSiteDialog from '@/components/EditSiteDialog';
import SitesWrapper from '@/components/SitesWrapper';
import { useAppContext } from '@/context';
import { SiteMetrics, TMetricResDto, TSiteForm } from '@/types';
import { getUnixTime } from '@/utils';
import { AlertColor, Box, Paper, Stack, Typography } from '@mui/material';
import { formatISO } from 'date-fns';
import { useEffect, useState } from 'react';
import PubSub from 'pubsub-js';
import MetricStatBySiteSubscription from '@/lib/MetricStatBySiteSubscription';

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

export default function Page() {
  const [sitesList, setSitesList] = useState<SiteDto[]>([]);
  const [componentsList, setComponentsList] = useState<any[]>([]);
  const { setSnackbarMessage, network, user, env, subscriptionClient } =
    useAppContext();
  const [openSiteConfig, setOpenSiteConfig] = useState(false);
  const [site, setSite] = useState<TSiteForm>(SITE_INIT);
  const [editSitedialogOpen, setEditSitedialogOpen] = useState(false);
  const [currentSite, setCurrentSite] = useState({
    siteName: '',
    siteId: '',
  });

  const [siteMetrics, setSiteMetrics] = useState<
    Record<string, Partial<SiteMetrics>>
  >({});

  const { refetch: refetchSites, loading: sitesLoading } = useGetSitesQuery({
    skip: !network.id,
    variables: {
      networkId: network.id,
    },
    onCompleted: (res) => {
      const sites = res.getSites.sites;
      setSitesList(sites);
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

  const [getSiteStatsMetrics] = useGetSiteStatLazyQuery({
    client: subscriptionClient,
    fetchPolicy: 'network-only',
    onError: (error) => {
      console.error('Error fetching site metrics:', error);
    },
  });

  const handleSiteStatSubscription = (_: any, data: string) => {
    try {
      const parsedData: TMetricResDto = JSON.parse(data);
      const { value, type, success, siteId } =
        parsedData.data.getSiteMetricStatSub;

      if (success && siteId) {
        setSiteMetrics((prev) => {
          const currentMetrics = prev[siteId] ? { ...prev[siteId] } : {};
          switch (type) {
            case 'site_uptime_seconds':
              currentMetrics.siteUptimeSeconds = value[1];
              break;
            case 'battery_charge_percentage':
              currentMetrics.batteryPercentage = value[1];
              break;
            case 'backhaul_speed':
              currentMetrics.backhaulSpeed = value[1];
              break;
          }
          return {
            ...prev,
            [siteId]: currentMetrics,
          };
        });
      }
    } catch (error) {
      console.error('Error handling site stat subscription:', error);
    }
  };

  const fetchSiteMetrics = async () => {
    const to = getUnixTime();
    const from = to - 40; // Example time range: last 40 seconds

    for (const site of sitesList) {
      try {
        const res = await getSiteStatsMetrics({
          variables: {
            data: {
              siteId: site.id,
              type: Stats_Type.Site,
              from,
              to,
              step: 30,
              userId: user.id,
              orgName: user.orgName,
              withSubscription: true,
            },
          },
        });

        if (
          res.data?.getSiteStat?.metrics &&
          res.data.getSiteStat.metrics.length > 0
        ) {
          const siteId = res.data.getSiteStat.metrics[0].siteId;
          const metrics = res.data.getSiteStat.metrics;

          setSiteMetrics((prev) => {
            const currentMetrics = prev[siteId] || {};
            metrics.forEach((metric) => {
              switch (metric.type) {
                case 'site_uptime_seconds':
                  currentMetrics.siteUptimeSeconds = metric.value;
                  break;
                case 'battery_charge_percentage':
                  currentMetrics.batteryPercentage = metric.value;
                  break;
                case 'backhaul_speed':
                  currentMetrics.backhaulSpeed = metric.value;
                  break;
              }
            });
            return {
              ...prev,
              [siteId]: currentMetrics,
            };
          });

          const sKey = `stat-${user.orgName}-${user.id}-${Stats_Type.Site}-${siteId}`;
          MetricStatBySiteSubscription({
            key: sKey,
            siteId: siteId,
            userId: user.id,
            url: env.METRIC_URL,
            orgName: user.orgName,
            type: Stats_Type.Site,
            from: getUnixTime() - 40,
          });

          PubSub.subscribe(sKey, handleSiteStatSubscription);
        }
      } catch (error) {
        console.error(`Error fetching metrics for site ${site.id}:`, error);
      }
    }
  };

  const [addSite, { loading: addSiteLoading }] = useAddSiteMutation({
    onCompleted: () => {
      refetchSites().then((res) => {
        setSitesList(res.data.getSites.sites);
      });
      setSnackbarMessage({
        id: 'add-site-success',
        message: 'Site added successfully!',
        type: 'success' as AlertColor,
        show: true,
      });
    },
    onError: (error) => {
      setSnackbarMessage({
        id: 'add-site-error',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });

  const [updateSite, { loading: updateSiteLoading }] = useUpdateSiteMutation({
    onCompleted: () => {
      refetchSites().then((res) => {
        setSitesList(res.data.getSites.sites);
      });
      setSnackbarMessage({
        id: 'update-site-success',
        message: 'Site updated successfully!',
        type: 'success' as AlertColor,
        show: true,
      });
    },
    onError: (error) => {
      setSnackbarMessage({
        id: 'update-site-error',
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

  const { data: networks, loading: networksLoading } = useGetNetworksQuery({
    fetchPolicy: 'cache-and-network',
    onCompleted: (res) => {
      if (res.getNetworks.networks.length === 0) {
        setSnackbarMessage({
          id: 'no-network-msg',
          message: ' CASA create a network first.',
          type: 'warning' as AlertColor,
          show: true,
        });
      }
    },
    onError: (error) => {
      setSnackbarMessage({
        id: 'networks-msg',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });

  useEffect(() => {
    if (!network.id)
      setSite({
        ...site,
        network: network.id,
      });
    getComponents({
      variables: {
        data: {
          category: Component_Type.All,
        },
      },
    });

    if (sitesList.length > 0) {
      fetchSiteMetrics();
    }

    return () => {
      sitesList.forEach((site) => {
        const sKey = `stat-${user.orgName}-${user.id}-${Stats_Type.Site}-${site.id}`;
        PubSub.unsubscribe(sKey);
      });
    };
  }, [sitesList]);

  useEffect(() => {
    const setupSubscriptions = async () => {
      const to = getUnixTime();
      const from = to - 30; // Last 30 seconds

      // Clear existing subscriptions
      sitesList.forEach((site) => {
        const sKey = `stat-${user.orgName}-${user.id}-${Stats_Type.Site}-${site.id}`;
        PubSub.unsubscribe(sKey);
      });

      // Set up new subscriptions for each site
      sitesList.forEach((site) => {
        const sKey = `stat-${user.orgName}-${user.id}-${Stats_Type.Site}-${site.id}`;

        // Initial subscription
        MetricStatBySiteSubscription({
          key: sKey,
          siteId: site.id,
          userId: user.id,
          url: env.METRIC_URL,
          orgName: user.orgName,
          type: Stats_Type.Site,
          from,
        });

        // Subscribe to updates
        PubSub.subscribe(sKey, handleSiteStatSubscription);

        // Set up interval for periodic updates
        const interval = setInterval(() => {
          MetricStatBySiteSubscription({
            key: sKey,
            siteId: site.id,
            userId: user.id,
            url: env.METRIC_URL,
            orgName: user.orgName,
            type: Stats_Type.Site,
            from: getUnixTime() - 30,
          });
        }, 30000); // Update every 30 seconds

        // Store interval ID for cleanup
        return () => {
          clearInterval(interval);
          PubSub.unsubscribe(sKey);
        };
      });
    };

    if (sitesList.length > 0) {
      setupSubscriptions();
    }

    // Cleanup function
    return () => {
      sitesList.forEach((site) => {
        const sKey = `stat-${user.orgName}-${user.id}-${Stats_Type.Site}-${site.id}`;
        PubSub.unsubscribe(sKey);
      });
    };
  }, [sitesList, user.id, user.orgName, env.METRIC_URL]);

  const handleCloseSiteConfig = () => {
    setSite(SITE_INIT);
    setOpenSiteConfig(false);
  };

  const handleSiteConfiguration = (values: TSiteForm) => {
    setSite(values);
    setOpenSiteConfig(false);
    addSite({
      variables: {
        data: {
          name: values.siteName,
          power_id: values.power,
          location: values.address,
          access_id: values.access,
          switch_id: values.switch,
          latitude: values.latitude,
          network_id: values.network,
          longitude: values.longitude,
          backhaul_id: values.backhaul,
          spectrum_id: values.spectrum,
          install_date: formatISO(new Date()),
        },
      },
    });
  };

  const handleSiteNameUpdate = (siteId: string, siteName: string) => {
    setCurrentSite((prevState) => ({
      ...prevState,
      siteId,
      siteName: siteName,
    }));
    setEditSitedialogOpen(true);
  };

  const handleSaveSiteName = (siteId: string, siteName: string) => {
    updateSite({
      variables: {
        siteId: siteId,
        data: {
          name: siteName,
        },
      },
    });
  };

  const closeEditSiteDialog = () => {
    setEditSitedialogOpen(false);
  };

  return (
    <Box mt={2}>
      <Paper
        sx={{
          overflow: 'auto',
          padding: '20px',
          borderRadius: '10px',
          height: 'calc(100vh - 212px)',
        }}
      >
        <Stack spacing={2} direction={'column'} height="100%">
          <Typography variant="h6" color="initial" sx={{ paddingLeft: '12px' }}>
            My sites
          </Typography>
          <SitesWrapper
            loading={sitesLoading || networksLoading}
            sites={sitesList}
            siteMetrics={siteMetrics}
            handleSiteNameUpdate={handleSiteNameUpdate}
          />
        </Stack>
      </Paper>
      <ConfigureSiteDialog
        site={site}
        open={openSiteConfig}
        addSiteLoading={addSiteLoading}
        onClose={handleCloseSiteConfig}
        components={componentsList || []}
        networks={networks?.getNetworks?.networks || []}
        handleSiteConfiguration={handleSiteConfiguration}
      />
      <EditSiteDialog
        open={editSitedialogOpen}
        siteId={currentSite.siteId}
        currentSiteName={currentSite.siteName}
        onClose={closeEditSiteDialog}
        onSave={handleSaveSiteName}
        updateSiteLoading={updateSiteLoading}
      />
    </Box>
  );
}
