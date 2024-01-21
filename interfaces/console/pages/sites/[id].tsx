/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { colors } from '@/styles/theme';
import LoadingWrapper from '@/ui/molecules/LoadingWrapper';
// import Map from '@/ui/molecules/MapComponent';
import SiteHeader from '@/ui/molecules/SiteHeader';
import SiteOverallHealth from '@/ui/molecules/SiteHealth';
import { AlertColor, Grid, Paper, Stack, Typography } from '@mui/material';
import { RoundedCard } from '@/styles/global';
import { SitePowerStatus } from '@/ui/molecules/SvgIcons';
import SiteDetailsCard from '@/ui/molecules/SiteDetailsCard';
import GroupIcon from '@mui/icons-material/Group';
import { useRouter } from 'next/router';
import AddSiteDialog from '@/ui/molecules/AddSiteDialog';
import {
  useGetNetworksQuery,
  useGetSingleSiteQuery,
  useGetAllSitesQuery,
  useGetNodesForSiteLazyQuery,
  useAddSiteToNetworkMutation,
} from '@/generated';
import { TCommonData, TSnackMessage } from '@/types';
import { commonData, snackbarMessage } from '@/app-recoil';
import { useRecoilValue, useSetRecoilState } from 'recoil';
import { useState } from 'react';

export default function Page() {
  const router = useRouter();
  const setSnackbarMessage = useSetRecoilState<TSnackMessage>(snackbarMessage);
  const _commonData = useRecoilValue<TCommonData>(commonData);
  const [numberOfNodesForSite, setNumberOfNodesForSite] = useState<number>();
  const [siteId, setSiteId] = useState<string>();
  const [isAddSiteDialogOpen, setIsAddSiteDialogOpen] = useState(false);
  const handleCloseAction = () => setIsAddSiteDialogOpen(false);
  const [addSite, { loading: addSiteLoading }] = useAddSiteToNetworkMutation({
    onCompleted: () => {
      setSnackbarMessage({
        id: 'site-added-success',
        message: 'Site added successfully!',
        type: 'success' as AlertColor,
        show: true,
      });
      setIsAddSiteDialogOpen(false);
    },
    onError: (error) => {
      setSnackbarMessage({
        id: 'site-added-error',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });

  const handleAddSite = (values: any, network: string) => {
    addSite({
      variables: {
        networkId: network,
        data: {
          site: values.site,
        },
      },
    });
  };
  const handleSiteRestart = () => {
    // Logic to restart a site
  };

  const batteryInfo = [
    { label: 'Model number', value: 'V1234' },
    { label: 'Current', value: '10 A' },
    { label: 'Charge', value: '80 %' },
    { label: 'Power', value: '100 W' },
    { label: 'Voltage', value: '12 V' },
  ];

  const { data: getSiteData, loading: getSiteLoading } = useGetSingleSiteQuery({
    fetchPolicy: 'cache-and-network',
    variables: {
      siteId: (router.query['id'] as string) ?? siteId,
      networkId: _commonData.networkId,
    },
    onError: (err) => {
      setSnackbarMessage({
        id: 'node-msg',
        message: err.message,
        type: 'error',
        show: true,
      });
    },
  });

  const { data: getAllSite, loading: getSitesLoading } = useGetAllSitesQuery({
    fetchPolicy: 'cache-and-network',
    variables: {
      networkId: _commonData.networkId,
    },
    onError: (err) => {
      setSnackbarMessage({
        id: 'site-msg',
        message: err.message,
        type: 'error',
        show: true,
      });
    },
  });
  const [getNodesForSite, { loading: NodesForSiteLoading }] =
    useGetNodesForSiteLazyQuery({
      fetchPolicy: 'cache-and-network',
      onCompleted: (data) => {
        setNumberOfNodesForSite(data.getNodesForSite.nodes.length);
      },
      onError: (err) => {
        setSnackbarMessage({
          id: 'nodesForSite-msg',
          message: err.message,
          type: 'error',
          show: true,
        });
      },
    });
  const { data: networkList, loading: netLoading } = useGetNetworksQuery({
    fetchPolicy: 'cache-and-network',

    onError: (error) => {
      setSnackbarMessage({
        id: 'networks-msg',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });
  const handleAddSiteAction = () => setIsAddSiteDialogOpen(true);

  const handleOnSiteSelect = (siteId: string) => {
    setSiteId(siteId);
    getNodesForSite({
      variables: {
        siteId: siteId ?? '',
      },
    });
  };
  return (
    <>
      <LoadingWrapper
        radius="small"
        width={'100%'}
        isLoading={
          getSiteLoading ||
          getSitesLoading ||
          NodesForSiteLoading ||
          addSiteLoading ||
          netLoading
        }
        cstyle={{
          backgroundColor: false ? colors.white : 'transparent',
        }}
      >
        <SiteHeader
          sites={getAllSite?.getAllSites.sites}
          addSiteAction={handleAddSiteAction}
          restartSiteAction={handleSiteRestart}
          onSiteSelect={handleOnSiteSelect}
        />

        <Grid container spacing={2} sx={{ mt: 1 }}>
          <SiteDetailsCard
            dateCreated={getSiteData?.getSingleSite.createdAt || ''}
            location={`1000 Nelson Way`}
            numberOfNodes={numberOfNodesForSite || 0}
          />

          <Grid item xs={6}>
            <RoundedCard>
              <Typography variant="h6" sx={{ py: 1 }}>
                Site overview
              </Typography>
              <Stack direction="row" spacing={2} alignItems={'center'}>
                <Stack direction="row" spacing={1} alignItems={'center'}>
                  <SitePowerStatus />
                  <Typography variant="body1">Input power</Typography>
                  <SitePowerStatus />
                  <Typography variant="body1">Storage</Typography>
                  <SitePowerStatus />
                  <Typography variant="body1">Consumption</Typography>
                </Stack>
              </Stack>
            </RoundedCard>
          </Grid>
          <Grid item xs={3}>
            <RoundedCard>
              <Stack direction="row" spacing={1}>
                <GroupIcon />
                <Typography variant="body1">
                  {numberOfNodesForSite || 0}
                </Typography>
              </Stack>
            </RoundedCard>
          </Grid>

          <Grid item xs={12}>
            <Paper sx={{ p: 2 }}>
              <Typography variant="h6" gutterBottom>
                Site components
              </Typography>

              <SiteOverallHealth
                solarHealth={'warning'}
                nodeHealth={'good'}
                switchHealth={'good'}
                controllerHealth={'good'}
                batteryHealth={'good'}
                backhaulHealth={'good'}
                batteryInfo={batteryInfo}
              />
              <AddSiteDialog
                isOpen={isAddSiteDialogOpen}
                title={'ADD SITE'}
                description=""
                handleCloseAction={handleCloseAction}
                networks={networkList?.getNetworks?.networks ?? []}
                handleAddSite={handleAddSite}
              />
            </Paper>
          </Grid>
        </Grid>
      </LoadingWrapper>
    </>
  );
}
