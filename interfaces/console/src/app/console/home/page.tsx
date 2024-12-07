/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';
import EmptyView from '@/components/EmptyView';
import LoadingWrapper from '@/components/LoadingWrapper';
import { LabelOverlayUI, SitesTree } from '@/components/NetworkMap/OverlayUI';
import NetworkStatus from '@/components/NetworkStatus';
import { MONTH_FILTER, TIME_FILTER } from '@/constants';
import { useAppContext } from '@/context';
import { colors } from '@/theme';
import DataVolume from '@mui/icons-material/DataSaverOff';
import GroupPeople from '@mui/icons-material/Group';
import NetworkIcon from '@mui/icons-material/Hub';
import Throughput from '@mui/icons-material/NetworkCheck';
import { Box, Paper, Skeleton, Stack } from '@mui/material';
import Grid from '@mui/material/Grid2';
import dynamic from 'next/dynamic';
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

const networkLoading = false;
const networkNodesLoading = false;
export default function Page() {
  const { network } = useAppContext();
  // const [filterState, setFilterState] = useState<NodeStatusEnum>(
  //   NodeStatusEnum.Undefined,
  // );
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
  //   client: getMetricClient("", ""),
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
      <Grid size={12}>
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
      <Grid size={12}>
        <Stack direction={'row'}>
          <StatusCard
            option={''}
            subtitle2={''}
            Icon={GroupPeople}
            subtitle1={`${0}`}
            options={TIME_FILTER}
            loading={networkLoading}
            title={'Active subscribers'}
            iconColor={colors.primaryMain}
            handleSelect={(value: string) => {}}
          />
          <Box p={{ xs: 0.5, md: 1 }} />
          <StatusCard
            Icon={DataVolume}
            option={'usage'}
            subtitle2={`dBM`}
            subtitle1={`${0}`}
            options={TIME_FILTER}
            loading={networkLoading}
            title={'Data Volume'}
            iconColor={colors.secondaryMain}
            handleSelect={(value: string) => {}}
          />
          <Box p={{ xs: 0.5, md: 1 }} />
          <StatusCard
            option={'bill'}
            subtitle2={`bps`}
            subtitle1={`${0}`}
            Icon={Throughput}
            options={MONTH_FILTER}
            loading={networkLoading}
            iconColor={colors.black54}
            title={'Average throughput'}
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
              isLoading={networkNodesLoading}
            >
              <NetworkMap
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
                      //   networkNodes?.getNodesByNetwork.nodes ?? [],
                      // )}
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
