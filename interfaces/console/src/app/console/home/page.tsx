/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';
import {
  useGetNodesLocationQuery,
  useGetSitesQuery,
} from '@/client/graphql/generated';
import EmptyView from '@/components/EmptyView';
import LoadingWrapper from '@/components/LoadingWrapper';
import { SitesTree } from '@/components/NetworkMap/OverlayUI';
import NetworkStatus from '@/components/NetworkStatus';
import { MONTH_FILTER, TIME_FILTER } from '@/constants';
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
  const { network, setSnackbarMessage } = useAppContext();
  const { data: sitesRes, loading: sitesLoading } = useGetSitesQuery({
    fetchPolicy: 'no-cache',
    variables: {
      networkId: network.id,
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

  const { data: nodesData, loading: nodesLoading } = useGetNodesLocationQuery({
    fetchPolicy: 'cache-and-network',
    onError: (error) => {
      setSnackbarMessage({
        id: 'home-nodes-err-msg',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });

  // const { data: statsRes, loading: statsLoading } = useGetStatsMetricQuery({
  //   client: getMetricClient("", ""),
  //   fetchPolicy: 'cache-and-network',
  // });

  return (
    <Grid container rowSpacing={2} columnSpacing={2}>
      <Grid size={12}>
        <NetworkStatus
          title={
            network.name
              ? `${network.name} is created.`
              : `No network selected.`
          }
          subtitle={network.name ? '' : ''}
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
            subtitle1={`${0}`}
            Icon={SalesIcon}
            options={MONTH_FILTER}
            loading={networkLoading}
            iconColor={colors.beige}
            title={'Sales'}
            handleSelect={(value: string) => {}}
          />
          <StatusCard
            Icon={DataVolume}
            option={'usage'}
            subtitle2={`GBs`}
            subtitle1={`${0}`}
            options={TIME_FILTER}
            loading={networkLoading}
            title={'Data Volume'}
            iconColor={colors.secondaryMain}
            handleSelect={(value: string) => {}}
          />
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
              isLoading={networkNodesLoading || sitesLoading || nodesLoading}
            >
              <NetworkMap
                id="network-map"
                zoom={10}
                className="network-map"
                markersData={{
                  nodes:
                    sitesRes && sitesRes?.getSites.sites.length > 0
                      ? nodesData?.getNodesLocation.nodes || []
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
                              ? nodesData?.getNodesLocation.nodes || []
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
