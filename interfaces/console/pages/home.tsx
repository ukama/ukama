/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';
import { MONTH_FILTER, TIME_FILTER } from '@/constants';
import { useAppContext } from '@/context';
import { NodeStatusEnum } from '@/generated';
import { DataBilling, DataUsage, UsersWithBG } from '@/public/svg';
import EmptyView from '@/ui/molecules/EmptyView';
import LoadingWrapper from '@/ui/molecules/LoadingWrapper';
import {
  LabelOverlayUI,
  SitesSelection,
  SitesTree,
} from '@/ui/molecules/NetworkMap/OverlayUI';
import NetworkStatus from '@/ui/molecules/NetworkStatus';
import StatusCard from '@/ui/molecules/StatusCard';
import NetworkIcon from '@mui/icons-material/Hub';
import { Paper } from '@mui/material';
import Grid from '@mui/material/Unstable_Grid2';
import dynamic from 'next/dynamic';
import { useState } from 'react';
const DynamicMap = dynamic(
  () => import('../ui/molecules/NetworkMap/DynamicMap'),
  {
    ssr: false,
  },
);

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
    <>
      <Grid container spacing={2}>
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
        <Grid xs={12} md={6} lg={4}>
          <StatusCard
            Icon={UsersWithBG}
            title={'Active subscribers'}
            options={TIME_FILTER}
            subtitle1={`${0}`}
            subtitle2={''}
            option={''}
            loading={networkLoading}
            handleSelect={(value: string) => {}}
          />
        </Grid>
        <Grid xs={12} md={6} lg={4}>
          <StatusCard
            title={'Average signal strength'}
            subtitle1={`${0}`}
            subtitle2={`dBM`}
            Icon={DataUsage}
            options={TIME_FILTER}
            option={'usage'}
            loading={networkLoading}
            handleSelect={(value: string) => {}}
          />
        </Grid>
        <Grid xs={12} md={6} lg={4}>
          <StatusCard
            title={'Average throughput'}
            subtitle1={`${0}`}
            subtitle2={`bps`}
            Icon={DataBilling}
            options={MONTH_FILTER}
            loading={networkLoading}
            option={'bill'}
            handleSelect={(value: string) => {}}
          />
        </Grid>
        <Grid xs={12}>
          <Paper
            sx={{
              borderRadius: '5px',
              height: 'calc(100vh - 310px)',
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
                      <SitesSelection
                        filterState={filterState}
                        handleFilterState={(value) => setFilterState(value)}
                      />
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
    </>
  );
}
