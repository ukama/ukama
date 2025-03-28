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
import { SiteMetrics, TSiteForm } from '@/types';
import { getUnixTime } from '@/utils';
import { AlertColor, Box, Paper, Stack, Typography } from '@mui/material';
import { formatISO } from 'date-fns';
import { useEffect, useState } from 'react';
import PubSub from 'pubsub-js';
import MetricStatSubscription from '@/lib/MetricStatSubscription';

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

      sites.forEach((site) => {
        fetchSiteMetrics(site.id);
      });
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

  const [getSiteMetrics] = useGetSiteStatLazyQuery({
    client: subscriptionClient,
    fetchPolicy: 'network-only',
    onCompleted: (data) => {
      if (data.getSiteStat.metrics.length > 0) {
        const siteId = data.getSiteStat.metrics[0].siteId;
        const metrics = data.getSiteStat.metrics;

        metrics.forEach((metric) => {
          if (metric.type === 'site_uptime_seconds') {
            // Store the actual uptime value directly from the API
            setSiteMetrics((prev) => {
              const currentMetrics = prev[siteId] || {};

              return {
                ...prev,
                [siteId]: {
                  ...currentMetrics,
                  siteUptimeSeconds: metric.value,
                },
              };
            });
          }
        });

        const sKey = `stat-${user.orgName}-${user.id}-${Stats_Type.Site}-${siteId}`;
        MetricStatSubscription({
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
    },
    onError: (error) => {
      console.error('Error fetching site metrics:', error);
    },
  });

  const [getBatteryMetrics] = useGetSiteStatLazyQuery({
    client: subscriptionClient,
    fetchPolicy: 'network-only',
    onCompleted: (data) => {
      if (data.getSiteStat.metrics.length > 0) {
        const siteId = data.getSiteStat.metrics[0].siteId;
        const metrics = data.getSiteStat.metrics;

        metrics.forEach((metric) => {
          if (metric.type === 'battery_charge_percentage') {
            setSiteMetrics((prev) => {
              const currentMetrics = prev[siteId] || {};

              return {
                ...prev,
                [siteId]: {
                  ...currentMetrics,
                  batteryPercentage: metric.value,
                },
              };
            });
          }
        });

        const sKey = `stat-${user.orgName}-${user.id}-${Stats_Type.Battery}-${siteId}`;
        MetricStatSubscription({
          key: sKey,
          siteId: siteId,
          userId: user.id,
          url: env.METRIC_URL,
          orgName: user.orgName,
          type: Stats_Type.Battery,
          from: getUnixTime() - 40,
        });

        PubSub.subscribe(sKey, handleBatteryStatSubscription);
      }
    },
    onError: (error) => {
      console.error('Error fetching battery metrics:', error);
    },
  });

  const [getBackhaulMetrics] = useGetSiteStatLazyQuery({
    client: subscriptionClient,
    fetchPolicy: 'network-only',
    onCompleted: (data) => {
      if (data.getSiteStat.metrics.length > 0) {
        const siteId = data.getSiteStat.metrics[0].siteId;
        const metrics = data.getSiteStat.metrics;

        metrics.forEach((metric) => {
          if (metric.type === 'backhaul_speed') {
            setSiteMetrics((prev) => {
              const currentMetrics = prev[siteId] || {};

              return {
                ...prev,
                [siteId]: {
                  ...currentMetrics,
                  backhaulSpeed: metric.value,
                },
              };
            });
          }
        });

        const sKey = `stat-${user.orgName}-${user.id}-${Stats_Type.MainBackhaul}-${siteId}`;
        MetricStatSubscription({
          key: sKey,
          siteId: siteId,
          userId: user.id,
          url: env.METRIC_URL,
          orgName: user.orgName,
          type: Stats_Type.MainBackhaul,
          from: getUnixTime() - 40,
        });

        PubSub.subscribe(sKey, handleBackhaulStatSubscription);
      }
    },
    onError: (error) => {
      console.error('Error fetching backhaul metrics:', error);
    },
  });

  const handleSiteStatSubscription = (_: any, data: string) => {
    try {
      const parsedData = JSON.parse(data);
      const { value, type, success, siteId } = parsedData.data.getMetricStatSub;

      if (success && type === 'site_uptime_seconds') {
        setSiteMetrics((prev) => {
          const currentMetrics = prev[siteId] || {};

          return {
            ...prev,
            [siteId]: {
              ...currentMetrics,
              siteUptimeSeconds: value,
            },
          };
        });
      }
    } catch (error) {
      console.error('Error handling site stat subscription:', error);
    }
  };

  const handleBatteryStatSubscription = (_: any, data: string) => {
    try {
      const parsedData = JSON.parse(data);
      const { value, type, success, siteId } = parsedData.data.getMetricStatSub;

      if (success && type === 'battery_charge_percentage') {
        setSiteMetrics((prev) => {
          const currentMetrics = prev[siteId] || {};

          return {
            ...prev,
            [siteId]: {
              ...currentMetrics,
              batteryPercentage: value,
            },
          };
        });
      }
    } catch (error) {
      console.error('Error handling battery stat subscription:', error);
    }
  };

  const handleBackhaulStatSubscription = (_: any, data: string) => {
    try {
      const parsedData = JSON.parse(data);
      const { value, type, success, siteId } = parsedData.data.getMetricStatSub;

      if (success && type === 'backhaul_speed') {
        setSiteMetrics((prev) => {
          const currentMetrics = prev[siteId] || {};

          return {
            ...prev,
            [siteId]: {
              ...currentMetrics,
              backhaulSpeed: value,
            },
          };
        });
      }
    } catch (error) {
      console.error('Error handling backhaul stat subscription:', error);
    }
  };

  const fetchSiteMetrics = (siteId: string) => {
    const to = getUnixTime();
    const from = to - 40;

    getSiteMetrics({
      variables: {
        data: {
          to,
          siteId,
          from,
          userId: user.id,
          step: 300,
          orgName: user.orgName,
          withSubscription: true,
          type: Stats_Type.Site,
        },
      },
    });

    getBatteryMetrics({
      variables: {
        data: {
          to,
          siteId,
          from,
          userId: user.id,
          step: 300,
          orgName: user.orgName,
          withSubscription: true,
          type: Stats_Type.Battery,
        },
      },
    });

    getBackhaulMetrics({
      variables: {
        data: {
          to,
          siteId,
          from,
          userId: user.id,
          step: 300,
          orgName: user.orgName,
          withSubscription: true,
          type: Stats_Type.MainBackhaul,
        },
      },
    });
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
          message: 'Please create a network first.',
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

    return () => {
      sitesList.forEach((site) => {
        const siteKey = `stat-${user.orgName}-${user.id}-${Stats_Type.Site}-${site.id}`;
        const batteryKey = `stat-${user.orgName}-${user.id}-${Stats_Type.Battery}-${site.id}`;
        const backhaulKey = `stat-${user.orgName}-${user.id}-${Stats_Type.MainBackhaul}-${site.id}`;

        PubSub.unsubscribe(siteKey);
        PubSub.unsubscribe(batteryKey);
        PubSub.unsubscribe(backhaulKey);
      });
    };
  }, [sitesList]);

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
