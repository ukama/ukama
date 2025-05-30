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
  useGetSitesQuery,
  useUpdateSiteMutation,
  NodeStateEnum,
  useGetNodesLazyQuery,
  NodeConnectivityEnum,
} from '@/client/graphql/generated';
import {
  Stats_Type,
  useGetSiteStatLazyQuery,
} from '@/client/graphql/generated/subscriptions';
import EditSiteDialog from '@/components/EditSiteDialog';
import SitesWrapper from '@/components/SitesWrapper';
import { useAppContext } from '@/context';
import { getUnixTime } from '@/utils';
import { AlertColor, Box, Paper, Stack, Typography } from '@mui/material';
import { useEffect, useRef, useState, useCallback, useMemo } from 'react';
import PubSub from 'pubsub-js';
import MetricStatBySiteSubscription from '@/lib/MetricStatBySiteSubscription';
import { useRouter } from 'next/navigation';
import { STAT_STEP_29 } from '@/constants';
import LoadingWrapper from '@/components/LoadingWrapper';
import colors from '@/theme/colors';

export default function Page() {
  const router = useRouter();
  const [sitesList, setSitesList] = useState<SiteDto[]>([]);
  const { setSnackbarMessage, network, user, env, subscriptionClient } =
    useAppContext();
  const [editSitedialogOpen, setEditSitedialogOpen] = useState(false);
  const [unassignedNodes, setUnassignedNodes] = useState<any[]>([]);

  const subscriptionsRef = useRef<Record<string, boolean>>({});
  const activeSubscriptionRef = useRef<any>(null);

  const [currentSite, setCurrentSite] = useState({
    siteName: '',
    siteId: '',
  });

  const cleanupSubscriptions = useCallback(() => {
    Object.keys(subscriptionsRef.current).forEach((topic) => {
      PubSub.unsubscribe(topic);
      delete subscriptionsRef.current[topic];
    });

    if (activeSubscriptionRef.current?.cancel) {
      activeSubscriptionRef.current.cancel();
      activeSubscriptionRef.current = null;
    }
  }, []);

  const {
    data: sitesData,
    refetch: refetchSites,
    loading: sitesLoading,
  } = useGetSitesQuery({
    fetchPolicy: 'network-only',
    skip: !network.id,
    variables: {
      data: { networkId: network.id },
    },
    onError: (error) => {
      setSitesList([]);
      setSnackbarMessage({
        id: 'sites-error',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });

  const [getNodes, { loading: nodesLoading }] = useGetNodesLazyQuery({
    fetchPolicy: 'network-only',
    onCompleted: (res) => {
      const allNodes = res.getNodes.nodes;
      const unknownNodes = allNodes.filter((node) => {
        const hasLocation = node.latitude !== 0 && node.longitude !== 0;
        return (
          (node.status.state === NodeStateEnum.Unknown &&
            node.status.connectivity == NodeConnectivityEnum.Online) ||
          !hasLocation
        );
      });

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

  useEffect(() => {
    if (sitesData?.getSites?.sites) {
      setSitesList(sitesData.getSites.sites);
    }
  }, [sitesData]);

  useEffect(() => {
    if (!network.id) return;

    setSitesList([]);
    setUnassignedNodes([]);
    cleanupSubscriptions();

    Promise.all([
      refetchSites(),
      getNodes({
        variables: {
          data: {
            state: NodeStateEnum.Unknown,
            connectivity: NodeConnectivityEnum.Online,
          },
        },
      }),
    ]);

    return cleanupSubscriptions;
  }, [network.id, refetchSites, getNodes, cleanupSubscriptions]);

  const [getSiteStatMetrics, { data: statData, loading: statLoading }] =
    useGetSiteStatLazyQuery({
      client: subscriptionClient,
      fetchPolicy: 'network-only',
      onCompleted: (data) => {
        if (data.getSiteStat.metrics.length > 0) {
          const from = getUnixTime() - STAT_STEP_29;
          const sKey = `stat-${user.orgName}-${user.id}-${Stats_Type.Site}-${from}`;

          if (!subscriptionsRef.current[sKey]) {
            subscriptionsRef.current[sKey] = true;

            if (activeSubscriptionRef.current?.cancel) {
              activeSubscriptionRef.current.cancel();
            }

            activeSubscriptionRef.current = MetricStatBySiteSubscription({
              key: sKey,
              nodeIds: [],
              userId: user.id,
              siteIds: sitesList.map((site) => site.id),
              url: env.METRIC_URL,
              orgName: user.orgName,
              type: Stats_Type.Site,
              from: from,
            });

            PubSub.subscribe(sKey, handleStatSubscription);
          }
        }
      },
    });

  const handleStatSubscription = useCallback((_: any, data: string) => {
    try {
      const parsedData = JSON.parse(data);
      const { value, type, success, siteId } =
        parsedData?.data?.getSiteMetricStatSub;
      if (!success || !siteId || !type) return;

      PubSub.publish(`stat-${type}-${siteId}`, value);
    } catch (error) {
      console.error('Error handling subscription data:', error);
    }
  }, []);

  const [updateSite, { loading: updateSiteLoading }] = useUpdateSiteMutation({
    onCompleted: () => {
      refetchSites();
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

  const handleSiteNameUpdate = useCallback(
    (siteId: string, siteName: string) => {
      setCurrentSite({
        siteId,
        siteName: siteName,
      });
      setEditSitedialogOpen(true);
    },
    [],
  );

  const handleSaveSiteName = useCallback(
    (siteId: string, siteName: string) => {
      updateSite({
        variables: {
          siteId: siteId,
          data: {
            name: siteName,
          },
        },
      });
    },
    [updateSite],
  );

  const closeEditSiteDialog = useCallback(() => {
    setEditSitedialogOpen(false);
  }, []);

  useEffect(() => {
    if (sitesList.length === 0 || !user.id) return;

    const to = getUnixTime();
    const from = to - STAT_STEP_29;

    getSiteStatMetrics({
      variables: {
        data: {
          to,
          from,
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
  }, [sitesList, user.id, user.orgName, getSiteStatMetrics]);

  const handleConfigureNode = useCallback(
    (nodeId: string) => {
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

      const p = new URLSearchParams();
      p.set('step', 'location');
      p.set('flow', 'ins');
      p.set('nid', nodeId);

      router.push(`/configure/check?${p.toString()}`);
    },
    [unassignedNodes, setSnackbarMessage, router],
  );

  useEffect(() => {
    return () => {
      cleanupSubscriptions();
    };
  }, [cleanupSubscriptions]);

  const memoizedStatData = useMemo(
    () => statData?.getSiteStat ?? { metrics: [] },
    [statData],
  );

  const isLoading = useMemo(
    () => nodesLoading || sitesLoading || statLoading,
    [nodesLoading, sitesLoading, statLoading],
  );

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
        <LoadingWrapper
          radius="small"
          width={'100%'}
          isLoading={isLoading}
          cstyle={{
            backgroundColor: false ? colors.white : 'transparent',
          }}
        >
          <Stack spacing={2} direction={'column'} height="100%">
            <Typography
              variant="h6"
              color="initial"
              sx={{ paddingLeft: '12px' }}
            >
              My sites
            </Typography>
            <SitesWrapper
              loading={nodesLoading || sitesLoading}
              sites={sitesList}
              siteMetricsStatData={memoizedStatData}
              handleSiteNameUpdate={handleSiteNameUpdate}
              handleConfigureNode={handleConfigureNode}
              unassignedNodes={unassignedNodes}
            />
          </Stack>
        </LoadingWrapper>
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
}
