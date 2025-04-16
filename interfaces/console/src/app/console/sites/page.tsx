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
import { TSiteForm } from '@/types';
import { getUnixTime } from '@/utils';
import { AlertColor, Box, Paper, Stack, Typography } from '@mui/material';
import { formatISO } from 'date-fns';
import { useEffect, useState, useRef } from 'react';
import PubSub from 'pubsub-js';
import MetricStatBySiteSubscription from '@/lib/MetricStatBySiteSubscription';
import { STAT_STEP_29 } from '@/constants';

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

  const subscriptionsInitialized = useRef<Record<string, boolean>>({});

  const { refetch: refetchSites } = useGetSitesQuery({
    skip: !network.id,
    variables: { data: { networkId: network.id } },
    onCompleted: (res) => {
      const sites = res.getSites.sites;
      sites.forEach((site) => {
        if (!subscriptionsInitialized.current[site.id]) {
          fetchSiteMetrics(site.id);
          subscriptionsInitialized.current[site.id] = true;
        }
      });
    },
  });

  const [getSiteMetrics] = useGetSiteStatLazyQuery({
    client: subscriptionClient,
    fetchPolicy: 'network-only',
    onCompleted: (data) => {
      if (data.getSiteStat.metrics.length === 0) return;
      const siteId = data.getSiteStat.metrics[0].siteId;
      data.getSiteStat.metrics.forEach((metric) => {
        PubSub.publish(`site-metrics-${siteId}`, {
          type: metric.type,
          value: metric.value,
        });
      });
    },
  });
  const setupSubscriptions = (siteId: string) => {
    const key = `stat-${user.orgName}-${user.id}-${Stats_Type.Site}-${siteId}`;
    PubSub.unsubscribe(key);

    MetricStatBySiteSubscription({
      url: env.METRIC_URL,
      key,
      from: getUnixTime() - STAT_STEP_29,
      siteId,
      userId: user.id,
      orgName: user.orgName,
      type: Stats_Type.Site,
    });

    PubSub.subscribe(key, (msg, data) => {
      try {
        const parsedData = JSON.parse(data);
        const { value, type, success, siteId } =
          parsedData.data.getSiteMetricStatSub;
        if (success) {
          PubSub.publish(`site-metrics-${siteId}`, { type, value: value[1] });
        }
      } catch (error) {
        console.error('Error handling metric update:', error);
      }
    });
  };
  const fetchSiteMetrics = (siteId: string) => {
    const to = getUnixTime();
    const from = to - STAT_STEP_29;
    getSiteMetrics({
      variables: {
        data: {
          to,
          siteId,
          from,
          userId: user.id,
          step: STAT_STEP_29,
          orgName: user.orgName,
          withSubscription: true,
          type: Stats_Type.Site,
        },
      },
    });
    setupSubscriptions(siteId);
  };

  useEffect(() => {
    const metricRequestToken = PubSub.subscribe('request-metrics-*', (msg) => {
      const siteId = msg.split('-').pop() || '';
      if (!subscriptionsInitialized.current[siteId]) {
        fetchSiteMetrics(siteId);
        subscriptionsInitialized.current[siteId] = true;
      }
    });

    return () => {
      PubSub.unsubscribe(metricRequestToken);
    };
  }, []);

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
    if (!network.id) return;

    refetchSites().then((res) => {
      const sites = res.data.getSites.sites;
      setSitesList(sites);
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
        const key = `stat-${user.orgName}-${user.id}-${Stats_Type.Site}-${site.id}`;
        PubSub.unsubscribe(key);
      });
    };
  }, [network.id]);

  useEffect(() => {
    sitesList.forEach((site) => {
      if (!subscriptionsInitialized.current[site.id]) {
        fetchSiteMetrics(site.id);
        subscriptionsInitialized.current[site.id] = true;
      }
    });
  }, [sitesList]);
  useEffect(() => {
    return () => {
      Object.keys(subscriptionsInitialized.current).forEach((siteId) => {
        PubSub.unsubscribe(
          `stat-${user.orgName}-${user.id}-${Stats_Type.Site}-${siteId}`,
        );
      });
    };
  }, []);

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
            loading={networksLoading}
            sites={sitesList}
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
