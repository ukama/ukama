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
import { useEffect, useState, useRef, useMemo } from 'react';
import PubSub from 'pubsub-js';
import MetricStatBySiteSubscription, {
  cancelSubscription,
} from '@/lib/MetricStatBySiteSubscription';
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

const initializedSites = new Set();

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

  const pubSubTokens = useRef<Record<string, string>>({});

  const [isInitializingSubscriptions, setIsInitializingSubscriptions] =
    useState(false);

  const { refetch: refetchSites, loading: sitesLoading } = useGetSitesQuery({
    skip: !network.id,
    variables: {
      data: { networkId: network.id },
    },
    onCompleted: (res) => {
      const sites = res.getSites.sites;
      setSitesList(sites);
    },
  });

  const [getSiteMetrics, { loading: metricsLoading }] = useGetSiteStatLazyQuery(
    {
      client: subscriptionClient,
      fetchPolicy: 'network-only',
      onCompleted: (data) => {
        if (!data?.getSiteStat?.metrics?.length) return;

        const siteId = data.getSiteStat.metrics[0].siteId;

        data.getSiteStat.metrics.forEach((metric) => {
          if (metric && metric.type && metric.value !== undefined) {
            PubSub.publish(`site-metrics-${siteId}`, {
              type: metric.type,
              value: metric.value,
            });
          }
        });
      },
      onError: (error) => {
        console.error('Error fetching site metrics:', error);
      },
    },
  );

  const initializeSiteMetrics = (siteId: string) => {
    if (initializedSites.has(siteId)) {
      return;
    }

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

    setupSiteSubscription(siteId, from);

    initializedSites.add(siteId);
  };

  const setupSiteSubscription = (siteId: string, from: number) => {
    const subscriptionKey = `stat-${user.orgName}-${user.id}-${Stats_Type.Site}-${siteId}`;

    cancelSubscription(subscriptionKey);
    PubSub.unsubscribe(subscriptionKey);

    MetricStatBySiteSubscription({
      url: env.METRIC_URL,
      key: subscriptionKey,
      from,
      siteId,
      userId: user.id,
      orgName: user.orgName,
      type: Stats_Type.Site,
    });

    const token = PubSub.subscribe(subscriptionKey, (msg, data) => {
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

    pubSubTokens.current[subscriptionKey] = token;
  };

  const initializeAllSiteMetrics = () => {
    if (!sitesList.length || !network.id || isInitializingSubscriptions) return;

    setIsInitializingSubscriptions(true);

    const batchSize = 3;
    const processBatch = (startIndex: number) => {
      const endIndex = Math.min(startIndex + batchSize, sitesList.length);

      for (let i = startIndex; i < endIndex; i++) {
        initializeSiteMetrics(sitesList[i].id);
      }

      if (endIndex < sitesList.length) {
        setTimeout(() => processBatch(endIndex), 500);
      } else {
        setIsInitializingSubscriptions(false);
      }
    };

    processBatch(0);
  };

  const cleanupAllSubscriptions = () => {
    Object.keys(pubSubTokens.current).forEach((key) => {
      const token = pubSubTokens.current[key];
      if (token) {
        PubSub.unsubscribe(token);
      }
      cancelSubscription(key);
    });

    pubSubTokens.current = {};
  };

  useEffect(() => {
    if (!network.id) return;

    cleanupAllSubscriptions();
    initializedSites.clear();

    refetchSites().then(() => {});

    getComponents({
      variables: {
        data: {
          category: Component_Type.All,
        },
      },
    });

    return cleanupAllSubscriptions;
  }, [network.id]);

  useEffect(() => {
    if (sitesList.length && network.id) {
      initializeAllSiteMetrics();
    }
  }, [sitesList]);

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

  const isLoading = useMemo(() => {
    return (
      sitesLoading ||
      networksLoading ||
      metricsLoading ||
      isInitializingSubscriptions
    );
  }, [
    sitesLoading,
    networksLoading,
    metricsLoading,
    isInitializingSubscriptions,
  ]);

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
            loading={isLoading}
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
