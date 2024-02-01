/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { commonData, snackbarMessage } from '@/app-recoil';
import { metricsClient } from '@/client/ApolloClient';
import { MONTH_FILTER, TIME_FILTER } from '@/constants';
import {
  NodeStatusEnum,
  useGetNodesByNetworkQuery,
  useGetNodesLocationQuery,
  useGetSitesQuery,
} from '@/generated';
import { useGetStatsMetricQuery } from '@/generated/metrics';
import { DataBilling, DataUsage, UsersWithBG } from '@/public/svg';
import { TCommonData, TSnackMessage } from '@/types';
import StatusCard from '@/ui/components/StatusCard';
import EmptyView from '@/ui/molecules/EmptyView';
import LoadingWrapper from '@/ui/molecules/LoadingWrapper';
import {
  LabelOverlayUI,
  SitesSelection,
  SitesTree,
} from '@/ui/molecules/NetworkMap/OverlayUI';
import NetworkStatus from '@/ui/molecules/NetworkStatus';
import { structureNodeSiteDate } from '@/utils';
import NetworkIcon from '@mui/icons-material/Hub';
import { AlertColor, Paper } from '@mui/material';
import Grid from '@mui/material/Unstable_Grid2';
import dynamic from 'next/dynamic';
import { useState } from 'react';
import { useRecoilValue, useSetRecoilState } from 'recoil';
const DynamicMap = dynamic(
  () => import('../ui/molecules/NetworkMap/DynamicMap'),
  {
    ssr: false,
  },
);

export default function Page() {
  const _commonData = useRecoilValue<TCommonData>(commonData);
  const [filterState, setFilterState] = useState<NodeStatusEnum>(
    NodeStatusEnum.Undefined,
  );
  const setSnackbarMessage = useSetRecoilState<TSnackMessage>(snackbarMessage);
  const { data: networkRes, loading: networkLoading } = useGetSitesQuery({
    fetchPolicy: 'cache-and-network',
    variables: {
      networkId: _commonData?.networkId,
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

  const { data: statsRes, loading: statsLoading } = useGetStatsMetricQuery({
    client: metricsClient,
    fetchPolicy: 'cache-and-network',
  });

  const { data: networkNodes, loading: networkNodesLoading } =
    useGetNodesByNetworkQuery({
      fetchPolicy: 'cache-and-network',
      variables: {
        networkId: _commonData?.networkId,
      },
      onError: (error) => {
        setSnackbarMessage({
          id: 'home-network-nodes-err-msg',
          message: error.message,
          type: 'error' as AlertColor,
          show: true,
        });
      },
    });

  const { data: nodesLocationData, loading: nodesLocationLoading } =
    useGetNodesLocationQuery({
      fetchPolicy: 'cache-first',
      variables: {
        data: {
          nodeFilterState: filterState,
          networkId: _commonData?.networkId,
        },
      },
    });

  return (
    <>
      <Grid container spacing={2}>
        <Grid xs={12}>
          <NetworkStatus
            loading={false}
            availableNodes={4}
            statusType="ONLINE"
            tooltipInfo="Network is online"
          />
        </Grid>
        <Grid xs={12} md={6} lg={4}>
          <StatusCard
            Icon={UsersWithBG}
            title={'Active subscribers'}
            options={TIME_FILTER}
            subtitle1={`${statsRes?.getStatsMetric.activeSubscriber}` || '0'}
            subtitle2={''}
            option={''}
            loading={statsLoading}
            handleSelect={(value: string) => {}}
          />
        </Grid>
        <Grid xs={12} md={6} lg={4}>
          <StatusCard
            title={'Average signal strength'}
            subtitle1={
              `${statsRes?.getStatsMetric.averageSignalStrength}` || '0'
            }
            subtitle2={`dBM`}
            Icon={DataUsage}
            options={TIME_FILTER}
            option={'usage'}
            loading={statsLoading}
            handleSelect={(value: string) => {}}
          />
        </Grid>
        <Grid xs={12} md={6} lg={4}>
          <StatusCard
            title={'Average throughput'}
            subtitle1={`${statsRes?.getStatsMetric.averageThroughput}` || '0'}
            subtitle2={`bps`}
            Icon={DataBilling}
            options={MONTH_FILTER}
            loading={statsLoading}
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
            {_commonData.networkId ? (
              <LoadingWrapper
                radius="small"
                width={'100%'}
                isLoading={nodesLocationLoading || networkNodesLoading}
              >
                <DynamicMap
                  id="network-map"
                  zoom={10}
                  className="network-map"
                  markersData={nodesLocationData?.getNodesLocation}
                >
                  {() => (
                    <>
                      <LabelOverlayUI name={_commonData.networkName} />
                      <SitesTree
                        sites={structureNodeSiteDate(
                          networkNodes?.getNodesByNetwork.nodes || [],
                        )}
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
