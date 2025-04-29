/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';
import { useGetNodesQuery, useGetSitesQuery } from '@/client/graphql/generated';
import EmptyView from '@/components/EmptyView';
import LoadingWrapper from '@/components/LoadingWrapper';
import { SitesTree } from '@/components/NetworkMap/OverlayUI';
import NetworkStatus from '@/components/NetworkStatus';
import { MONTH_FILTER, NODE_KPIS, TIME_FILTER } from '@/constants';
import { useAppContext } from '@/context';
import { colors } from '@/theme';
import { structureNodeSiteDate } from '@/utils';
import DataVolume from '@mui/icons-material/DataSaverOff';
import GroupPeople from '@mui/icons-material/Group';
import NetworkIcon from '@mui/icons-material/Hub';
import SalesIcon from '@mui/icons-material/MonetizationOn';
import { AlertColor, Paper, Skeleton, Stack } from '@mui/material';
import Grid from '@mui/material/Grid2';
import dynamic from 'next/dynamic';
import { useState } from 'react';
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

const getFrom = () => Math.floor(Date.now() / 1000) - 10000; //past 2.45 hours

export default function Page() {
  const kpiConfig = NODE_KPIS.HOME.stats;
  const [networkStats, setNetworkStats] = useState({
    uptime: 0,
    sales: 0,
    dataVolume: 0,
    activeSubscribers: 0,
  });
  const { user, network, setSnackbarMessage, subscriptionClient } =
    useAppContext();

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

  return (
    <Grid container rowSpacing={2} columnSpacing={2}>
      <Grid size={12}>
        <NetworkStatus
          title={
            network.name
              ? `${network.name} is created. Network is online for `
              : `No network selected.`
          }
          subtitle={network.name ? networkStats.uptime : 0}
          loading={false}
          availableNodes={undefined}
          statusType="ONLINE"
          tooltipInfo="Network is online"
        />
      </Grid>
      <Grid size={12}>
        <Stack direction={'row'} spacing={{ xs: 0.5, md: 1 }}>
          <StatusCard
            option={'bill'}
            subtitle2={`$`}
            subtitle1={`${networkStats.sales}`}
            Icon={SalesIcon}
            options={MONTH_FILTER}
            loading={false}
            iconColor={colors.beige}
            title={'Sales'}
            handleSelect={(value: string) => {}}
          />
          <StatusCard
            Icon={DataVolume}
            option={'usage'}
            subtitle2={`GBs`}
            subtitle1={`${networkStats.dataVolume}`}
            options={TIME_FILTER}
            loading={false}
            title={'Data Volume'}
            iconColor={colors.secondaryMain}
            handleSelect={(value: string) => {}}
          />
          <StatusCard
            option={''}
            subtitle2={''}
            Icon={GroupPeople}
            subtitle1={`${networkStats.activeSubscribers}`}
            options={TIME_FILTER}
            loading={false}
            title={'Active subscribers'}
            iconColor={colors.primaryMain}
            handleSelect={(value: string) => {}}
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
