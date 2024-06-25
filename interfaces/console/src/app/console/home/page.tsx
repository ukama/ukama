/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';
import { DataVolume, Throughput, UsersWithBG } from '@/../public/svg';
import { NodeStatusEnum } from '@/client/graphql/generated';
import EmptyView from '@/components/EmptyView';
import LoadingWrapper from '@/components/LoadingWrapper';
import { LabelOverlayUI, SitesTree } from '@/components/NetworkMap/OverlayUI';
import NetworkStatus from '@/components/NetworkStatus';
import { MONTH_FILTER, TIME_FILTER } from '@/constants';
import { useAppContext } from '@/context';
import NetworkIcon from '@mui/icons-material/Hub';
import { Box, Paper, Skeleton, Stack } from '@mui/material';
import Grid from '@mui/material/Unstable_Grid2';
import dynamic from 'next/dynamic';
import { useState } from 'react';
const DynamicMap = dynamic(() => import('@/components/NetworkMap/DynamicMap'), {
  ssr: false,
  loading: () => (
    <Skeleton
      variant="rectangular"
      sx={{
        borderRadius: '10px',
        height: 'calc(100vh - 332px)',
        width: '100%',
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

const networkLoading = false;
const networkNodesLoading = false;
export default function Page() {
  const { network, setSnackbarMessage } = useAppContext();
  const [filterState, setFilterState] = useState<NodeStatusEnum>(
    NodeStatusEnum.Undefined,
  );
  // const { data: networkRes, loading: networkLoading } = useGetSitesQuery({
  //   fetchPolicy: 'no-cache',
  //   variables: {
  //     networkId: network.id,
  //   },
  //   onError: (error) => {
  //     setSnackbarMessage({
  //       id: 'home-sites-err-msg',
  //       message: error.message,
  //       type: 'error' as AlertColor,
  //       show: true,
  //     });
  //   },
  // });

  // const { data: statsRes, loading: statsLoading } = useGetStatsMetricQuery({
  //   client: metricsClient,
  //   fetchPolicy: 'cache-and-network',
  // });

  // const { data: nodesLocationData, loading: nodesLocationLoading } =
  //   useGetNodesLocationQuery({
  //     fetchPolicy: 'cache-first',
  //     variables: {
  //       data: {
  //         nodeFilterState: filterState,
  //         networkId: network.id,
  //       },
  //     },
  //   });

  // const { data: networkNodes, loading: networkNodesLoading } =
  //   useGetNodesByNetworkQuery({
  //     fetchPolicy: 'cache-and-network',
  //     variables: {
  //       networkId: network.id,
  //     },
  //     onError: (error) => {
  //       setSnackbarMessage({
  //         id: 'home-network-nodes-err-msg',
  //         message: error.message,
  //         type: 'error' as AlertColor,
  //         show: true,
  //       });
  //     },
  //   });

  return (
    <Grid container rowSpacing={2} columnSpacing={2}>
      <Grid xs={12}>
        <NetworkStatus
          title={
            network.name
              ? `${network.name} is created.`
              : `No network selected.`
          }
          subtitle={network.name ? 'No node attached to this network.' : ''}
          loading={false}
          availableNodes={undefined}
          statusType="ONLINE"
          tooltipInfo="Network is online"
        />
      </Grid>
      <Grid xs={12}>
        <Stack direction={'row'}>
          <StatusCard
            option={''}
            subtitle2={''}
            Icon={UsersWithBG}
            subtitle1={`${0}`}
            options={TIME_FILTER}
            loading={networkLoading}
            title={'Active subscribers'}
            handleSelect={(value: string) => {}}
          />
          <Box p={1} />
          <StatusCard
            Icon={DataVolume}
            option={'usage'}
            subtitle2={`dBM`}
            subtitle1={`${0}`}
            options={TIME_FILTER}
            loading={networkLoading}
            title={'Data Volume'}
            handleSelect={(value: string) => {}}
          />
          <Box p={1} />
          <StatusCard
            option={'bill'}
            subtitle2={`bps`}
            subtitle1={`${0}`}
            Icon={Throughput}
            options={MONTH_FILTER}
            loading={networkLoading}
            title={'Average throughput'}
            handleSelect={(value: string) => {}}
          />
        </Stack>
      </Grid>
      <Grid xs={12}>
        <Paper
          sx={{
            borderRadius: '10px',
            height: 'calc(100vh - 332px)',
          }}
        >
          {network.id ? (
            <LoadingWrapper
              radius="small"
              width={'100%'}
              isLoading={networkNodesLoading}
            >
              <DynamicMap
                id="network-map"
                zoom={10}
                className="network-map"
                markersData={{ nodes: [], networkId: '' }}
              >
                {() => (
                  <>
                    <LabelOverlayUI name={network.name} />
                    <SitesTree
                      sites={[]}
                      // sites={structureNodeSiteDate(
                      //   networkNodes?.getNodesByNetwork.nodes || [],
                      // )}
                    />
                    {/* <SitesSelection
                      filterState={filterState}
                      handleFilterState={(value) => setFilterState(value)}
                    /> */}
                  </>
                )}
              </DynamicMap>
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
