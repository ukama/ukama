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
  useGetNetworksQuery,
  useGetSitesQuery,
  useUpdateSiteMutation,
  useGetNodeStateLazyQuery,
  useGetNodesQuery,
  NodeStateEnum,
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
import MetricStatBySiteSubscription from '@/lib/MetricStatBySiteSubscription';
import { usePathname, useRouter, useSearchParams } from 'next/navigation';

import {
  CHECK_SITE_FLOW,
  INSTALLATION_FLOW,
  NETWORK_FLOW,
  ONBOARDING_FLOW,
  STAT_STEP_29,
} from '@/constants';
import { setQueryParam } from '@/utils';

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
  const router = useRouter();
  const [sitesList, setSitesList] = useState<SiteDto[]>([]);
  const [componentsList, setComponentsList] = useState<any[]>([]);
  const { setSnackbarMessage, network, user, env, subscriptionClient } =
    useAppContext();
  const [openSiteConfig, setOpenSiteConfig] = useState(false);
  const [site, setSite] = useState<TSiteForm>(SITE_INIT);
  const [editSitedialogOpen, setEditSitedialogOpen] = useState(false);
  const [selectedNodeId, setSelectedNodeId] = useState<string | null>(null);
  const [unassignedNodes, setUnassignedNodes] = useState<any[]>([]);
  const searchParams = useSearchParams();
  const flow = searchParams.get('flow') ?? INSTALLATION_FLOW;
  const pathname = usePathname();

  const [currentSite, setCurrentSite] = useState({
    siteName: '',
    siteId: '',
  });

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
  const { loading: nodesLoading } = useGetNodesQuery({
    onCompleted: async (res) => {
      const allNodes = res.getNodes.nodes;
      const unknownNodes = [];

      for (const node of allNodes) {
        const { data } = await getNodeState({
          variables: { getNodeStateId: node.id },
        });

        const hasLocation = node.latitude !== 0 && node.longitude !== 0;

        if (
          (data?.getNodeState.currentState === NodeStateEnum.Unknown &&
            (node.site.siteId === '' || node.site.siteId == null)) ||
          !hasLocation
        ) {
          unknownNodes.push(node);
        }
      }

      setUnassignedNodes(unknownNodes);
    },
    onError: (error) => {
      setSnackbarMessage({
        id: 'unassigned-nodes-error',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });
  const [
    getSiteStatMetrics,
    { data: statData, loading: statLoading, variables: statVar },
  ] = useGetSiteStatLazyQuery({
    client: subscriptionClient,
    fetchPolicy: 'network-only',
    onCompleted: (data) => {
      if (data.getSiteStat.metrics.length > 0) {
        data.getSiteStat.metrics.forEach((m) => {
          console.log(m.type);
          console.log(m.value);
        });

        const sKey = `stat-${user.orgName}-${user.id}-${Stats_Type.Site}-${statVar?.data.from ?? 0}`;

        MetricStatBySiteSubscription({
          key: sKey,
          nodeIds: [],
          userId: user.id,
          siteIds: sitesList.map((site) => site.id),
          url: env.METRIC_URL,
          orgName: user.orgName,
          type: Stats_Type.Site,
          from: statVar?.data.from ?? 0,
        });
        PubSub.subscribe(sKey, handleStatSubscription);
      }
    },
  });
  const handleStatSubscription = (_: any, data: string) => {
    try {
      const parsedData = JSON.parse(data);

      const { msg, value, type, success, siteId } =
        parsedData.data.getSiteMetricStatSub;

      if (success) {
        PubSub.publish(`stat-${type}`, { value, siteId });

        if (siteId) {
          const siteTopic = `site-metrics-${siteId}`;
          console.log(`Publishing to site topic ${siteTopic}`);
          PubSub.publish(siteTopic, {
            metrics: [{ type, value }],
          });
        }
      }
    } catch (error) {
      console.error('Error handling subscription data:', error);
    }
  };
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

  useEffect(() => {
    if (sitesList.length === 0) return;

    const to = getUnixTime();
    const from = to - STAT_STEP_29;
    const sKey = `stat-${user.orgName}-${user.id}-${Stats_Type.Site}-${from ?? 0}`;

    getSiteStatMetrics({
      variables: {
        data: {
          to: to,
          from: from,
          userId: user.id,
          nodeIds: [],
          siteIds: sitesList.map((site) => site.id),
          step: STAT_STEP_29,
          orgName: user.orgName,
          withSubscription: true,
          type: Stats_Type.Site,
        },
      },
    });

    const currentSKey = sKey;

    return () => {
      console.log('Unsubscribing from:', currentSKey);
      PubSub.unsubscribe(currentSKey);
    };
  }, [sitesList, user.id, user.orgName]);
  const [getNodeState] = useGetNodeStateLazyQuery({
    fetchPolicy: 'network-only',
    onCompleted: (data) => {},
  });
  const handleConfigureNode = async (nodeId: string) => {
    const node = unassignedNodes.find((n) => n.id === nodeId);

    if (!node) {
      setSnackbarMessage({
        id: 'node-not-found',
        message: 'Node not found',
        type: 'error' as AlertColor,
        show: true,
      });
      return;
    }

    try {
      const result = await getNodeState({
        variables: {
          getNodeStateId: nodeId,
        },
      });

      if (result.data?.getNodeState.currentState === NodeStateEnum.Unknown) {
        let p = setQueryParam(
          'lat',
          node.latitude?.toString() || '0',
          searchParams.toString(),
          pathname,
        );
        p.set('lng', node.longitude?.toString() || '0');
        p.set(
          'flow',
          flow === NETWORK_FLOW
            ? ONBOARDING_FLOW
            : flow === CHECK_SITE_FLOW
              ? INSTALLATION_FLOW
              : flow,
        );
        p.delete('nid');
        router.push(`/configure/node/${node.id}?${p.toString()}`);
      }
    } catch (error) {
      console.error('Error checking node state:', error);
      setSnackbarMessage({
        id: 'node-state-error',
        message: 'Error checking node state',
        type: 'error' as AlertColor,
        show: true,
      });
    }
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
            loading={statLoading}
            sites={sitesList}
            siteMetricsStatData={statData?.getSiteStat ?? { metrics: [] }}
            handleSiteNameUpdate={handleSiteNameUpdate}
            handleConfigureNode={handleConfigureNode}
            unassignedNodes={unassignedNodes}
          />
        </Stack>
      </Paper>
      <ConfigureSiteDialog
        site={site}
        open={openSiteConfig}
        addSiteLoading={false}
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
