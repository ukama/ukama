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
  useGetSitesQuery,
  useGetNodesQuery,
  useGetSubscribersByNetworkQuery,
  useUpdateSiteMutation,
} from '@/client/graphql/generated';
import EditSiteDialog from '@/components/EditSiteDialog';
import SitesWrapper from '@/components/SitesWrapper';
import { useAppContext } from '@/context';
import { TSiteForm } from '@/types';
import { AlertColor, Box, Paper, Typography } from '@mui/material';
import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import {
  Graphs_Type,
  MetricsRes,
  useGetMetricByTabLazyQuery,
} from '@/client/graphql/generated/subscriptions';
import { getUnixTime } from '@/utils';
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

const Sites = () => {
  const router = useRouter();
  const [sitesList, setSitesList] = useState<SiteDto[]>([]);
  const [componentsList, setComponentsList] = useState<any[]>([]);
  const { setSnackbarMessage, user, subscriptionClient, network } =
    useAppContext();
  const [site, setSite] = useState<TSiteForm>(SITE_INIT);
  const [editSitedialogOpen, setEditSitedialogOpen] = useState(false);
  const [unnamedNodes, setUnnamedNodes] = useState<any[]>([]);
  const [siteMetrics, setSiteMetrics] = useState<Record<string, MetricsRes>>(
    {},
  );

  const [currentSite, setCurrentSite] = useState({
    siteName: '',
    siteId: '',
  });
  const { refetch: refetchSites, loading: sitesLoading } = useGetSitesQuery({
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
  }, []);

  const [getSiteMetricByTab] = useGetMetricByTabLazyQuery({
    client: subscriptionClient,
    fetchPolicy: 'network-only',
    onCompleted: (data) => {
      const siteId = data.getMetricByTab.metrics[0]?.siteId;
      if (siteId) {
        setSiteMetrics((prev) => ({
          ...prev,
          [siteId]: {
            ...prev[siteId],
            metrics: [
              ...(prev[siteId]?.metrics || []),
              ...data.getMetricByTab.metrics,
            ],
          },
        }));
      }
    },
  });

  useEffect(() => {
    if (sitesList.length > 0) {
      sitesList.forEach((site) => {
        getSiteMetricByTab({
          variables: {
            data: {
              nodeId: '',
              siteId: site.id,
              userId: user.id,
              type: Graphs_Type.Battery,
              from: getUnixTime() - 140,
              orgName: user.orgName,
              withSubscription: true,
            },
          },
        });

        getSiteMetricByTab({
          variables: {
            data: {
              nodeId: '',
              siteId: site.id,
              userId: user.id,
              type: Graphs_Type.Backhaul,
              from: getUnixTime() - 140,
              orgName: user.orgName,
              withSubscription: true,
            },
          },
        });
      });
    }
  }, [sitesList, user.id, user.orgName]);

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
  const { data: subscribers } = useGetSubscribersByNetworkQuery({
    variables: {
      networkId: network.id,
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
  const { data: nodes } = useGetNodesQuery({
    fetchPolicy: 'cache-and-network',
    onError: (error) => {
      setSnackbarMessage({
        id: 'nodes-msg',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });
  useEffect(() => {
    if (nodes) {
      const unnamedNodes = nodes.getNodes.nodes.filter(
        (node) => !node.site.siteId,
      );
      setUnnamedNodes(unnamedNodes);
    }
  }, [nodes]);
  const handleSiteConfig = () => {
    router.push(`/configure/network?step=1`);
  };
  return (
    <Box mt={2}>
      <Paper
        sx={{
          p: 4,
          overflow: 'auto',
          borderRadius: '10px',
          height: 'calc(100vh - 212px)',
        }}
      >
        <Typography variant="h6" color="initial" sx={{ mb: 2 }}>
          My sites
        </Typography>
        <SitesWrapper
          loading={sitesLoading}
          sites={sitesList}
          handleSiteNameUpdate={handleSiteNameUpdate}
          subscriberCount={
            subscribers?.getSubscribersByNetwork.subscribers.length
          }
          unnamedNodes={unnamedNodes}
          handleConfigureSite={handleSiteConfig}
          siteMetrics={siteMetrics}
        />
      </Paper>
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
};

export default Sites;
