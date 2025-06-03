/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';
import {
  useGetCurrencySymbolQuery,
  useGetNodesQuery,
  useGetSitesQuery,
} from '@/client/graphql/generated';
import {
  Stats_Type,
  useGetMetricsStatLazyQuery,
} from '@/client/graphql/generated/subscriptions';
import EmptyView from '@/components/EmptyView';
import LoadingWrapper from '@/components/LoadingWrapper';
import { SitesTree } from '@/components/NetworkMap/OverlayUI';
import { MONTH_FILTER, NODE_KPIS, TIME_FILTER } from '@/constants';
import { useAppContext } from '@/context';
import MetricStatSubscription from '@/lib/MetricStatSubscription';
import { colors } from '@/theme';
import { TMetricResDto } from '@/types';
import { formatBytesToGB, getUnixTime, structureNodeSiteDate } from '@/utils';
import DataVolume from '@mui/icons-material/DataSaverOff';
import GroupPeople from '@mui/icons-material/Group';
import NetworkIcon from '@mui/icons-material/Hub';
import SalesIcon from '@mui/icons-material/MonetizationOn';
import { AlertColor, Paper, Skeleton, Stack } from '@mui/material';
import Grid from '@mui/material/Grid2';
import dynamic from 'next/dynamic';
import { useCallback, useEffect, useRef } from 'react';
const NetworkMap = dynamic(() => import('@/components/NetworkMap'), {
  ssr: false,
  loading: () => (
    <Skeleton
      variant="rectangular"
      sx={{
        width: '100%',
        borderRadius: '10px',
        height: { xs: 'calc(100vh - 280px)', md: 'calc(100vh - 318px)' },
      }}
    />
  ),
});
const StatusCard = dynamic(() => import('@/components/StatusCard'), {
  ssr: false,
  loading: () => (
    <Skeleton
      variant="rectangular"
      sx={{ borderRadius: '10px', height: '90px', width: '100%', m: 0 }}
    />
  ),
});

export default function Page() {
  const kpiConfig = NODE_KPIS.HOME.stats;
  const { env, user, network, setSnackbarMessage, subscriptionClient } =
    useAppContext();
  const subscriptionKeyRef = useRef<string | null>(null);
  const subscriptionControllerRef = useRef<AbortController | null>(null);

  const cleanupSubscription = useCallback(() => {
    if (subscriptionKeyRef.current) {
      PubSub.unsubscribe(subscriptionKeyRef.current);
      subscriptionKeyRef.current = null;
    }
    if (subscriptionControllerRef.current) {
      subscriptionControllerRef.current.abort();
      subscriptionControllerRef.current = null;
    }
  }, []);

  const { data: sitesRes, loading: sitesLoading } = useGetSitesQuery({
    skip: !network?.id,
    fetchPolicy: 'no-cache',
    variables: {
      data: {
        networkId: network.id,
      },
    },
    onError: (error) => {
      setSnackbarMessage({
        id: 'home-sites-err-msg',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });

  const { data: nodesData, loading: nodesLoading } = useGetNodesQuery({
    fetchPolicy: 'cache-and-network',
    variables: { data: {} },
    onError: (error) => {
      setSnackbarMessage({
        id: 'home-nodes-err-msg',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });

  const { data: currencyData } = useGetCurrencySymbolQuery({
    skip: !user.currency,
    fetchPolicy: 'cache-first',
    variables: {
      code: user.currency,
    },
    onError: (error) => {
      setSnackbarMessage({
        id: 'currency-info-error',
        message: error.message,
        type: 'error',
        show: true,
      });
    },
  });

  const [getMetricStat, { loading: statLoading, variables: statVar }] =
    useGetMetricsStatLazyQuery({
      client: subscriptionClient,
      fetchPolicy: 'network-only',
      onCompleted: async (data) => {
        if (data.getMetricsStat.metrics.length > 0) {
          const sKey = `stat-${user.orgName}-${user.id}-${Stats_Type.Home}-${statVar?.data.from ?? 0}`;
          cleanupSubscription();
          subscriptionKeyRef.current = sKey;

          if (statVar?.data.withSubscription) {
            const controller = await MetricStatSubscription({
              key: sKey,
              userId: user.id,
              url: env.METRIC_URL,
              networkId: network.id,
              orgName: user.orgName,
              type: Stats_Type.Home,
              from: statVar?.data.from ?? 0,
            });

            data.getMetricsStat.metrics.forEach((metric) => {
              PubSub.publish(
                `${metric.type}-${network.id}`,
                metric.type === kpiConfig[2].id
                  ? formatBytesToGB(metric.value)
                  : metric.value,
              );
            });

            subscriptionControllerRef.current = controller;
            PubSub.subscribe(sKey, handleStatSubscription);
          }
        }
      },
    });

  useEffect(() => {
    const to = getUnixTime();
    const from = to;
    if (network.id) {
      cleanupSubscription();

      getMetricStat({
        variables: {
          data: {
            to: to,
            step: 1,
            from: from,
            userId: user.id,
            operation: 'sum',
            networkId: network.id,
            orgName: user.orgName,
            withSubscription: true,
            type: Stats_Type.Home,
          },
        },
      });
    }
  }, [network.id, user.id, user.orgName, cleanupSubscription]);

  const handleStatSubscription = (_: any, data: string) => {
    const parsedData: TMetricResDto = JSON.parse(data);
    const { value, type, success } = parsedData.data.getMetricStatSub;

    if (success && value.length === 2) {
      PubSub.publish(
        `${type}-${network.id}`,
        type === kpiConfig[2].id ? formatBytesToGB(value[1]) : value[1],
      );
    }
  };

  return (
    <Grid container rowSpacing={2} columnSpacing={2}>
      <Grid size={12}>
        {/* TODO: Need more discussion
        <NetworkStatus
          title={
            network.name
              ? `${network.name} is created. Network is online for `
              : `No network selected.`
          }
          subtitle={network.name ? networkStats.uptime : 0}
          loading={statLoading}
          availableNodes={undefined}
          statusType="ONLINE"
          tooltipInfo="Network is online"
        /> */}
      </Grid>
      <Grid size={12}>
        <Stack direction={'row'} spacing={{ xs: 0.5, md: 1 }}>
          <StatusCard
            option={'bill'}
            subtitle2={currencyData?.getCurrencySymbol.symbol ?? ''}
            title={'Sales'}
            Icon={SalesIcon}
            loading={statLoading}
            options={MONTH_FILTER}
            iconColor={colors.beige}
            topic={`${kpiConfig[1].id}-${network.id}`}
            handleSelect={() => {}}
          />
          <StatusCard
            Icon={DataVolume}
            option={'usage'}
            subtitle2={'GBs'}
            options={TIME_FILTER}
            loading={statLoading}
            title={'Data Volume'}
            iconColor={colors.secondaryMain}
            topic={`${kpiConfig[2].id}-${network.id}`}
            handleSelect={() => {}}
          />
          <StatusCard
            option={''}
            subtitle2={''}
            Icon={GroupPeople}
            options={TIME_FILTER}
            loading={statLoading}
            title={'Active subscribers'}
            iconColor={colors.primaryMain}
            topic={`${kpiConfig[3].id}-${network.id}`}
            handleSelect={() => {}}
          />
        </Stack>
      </Grid>
      <Grid size={12}>
        <Paper
          sx={{
            borderRadius: '10px',
            height: { xs: 'calc(100vh - 280px)', md: 'calc(100vh - 318px)' },
          }}
        >
          {network.id ? (
            <LoadingWrapper
              radius="small"
              width={'100%'}
              isLoading={sitesLoading || nodesLoading}
            >
              <NetworkMap
                id="network-map"
                zoom={10}
                className="network-map"
                markersData={{
                  nodes:
                    sitesRes && sitesRes?.getSites.sites.length > 0
                      ? nodesData?.getNodes.nodes || []
                      : [],
                }}
              >
                {() => (
                  <>
                    <SitesTree
                      sites={structureNodeSiteDate(
                        {
                          nodes:
                            sitesRes && sitesRes?.getSites.sites.length > 0
                              ? nodesData?.getNodes.nodes || []
                              : [],
                        },
                        { sites: sitesRes?.getSites.sites || [] },
                      )}
                    />
                    {/* <SitesSelection
                      filterState={filterState}
                      handleFilterState={(value) => setFilterState(value)}
                    /> */}
                  </>
                )}
              </NetworkMap>
            </LoadingWrapper>
          ) : (
            <EmptyView
              title="No network selected"
              icon={NetworkIcon}
              size="medium"
            />
          )}
        </Paper>
      </Grid>
    </Grid>
  );
}
