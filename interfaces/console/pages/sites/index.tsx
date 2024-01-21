/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { commonData, snackbarMessage } from '@/app-recoil';
import { colors } from '@/styles/theme';
import LoadingWrapper from '@/ui/molecules/LoadingWrapper';
import { useState } from 'react';
// import Map from '@/ui/molecules/MapComponent';
import { Grid, Typography, AlertColor, Button } from '@mui/material';
import SiteCard from '@/ui/molecules/SiteCard';
import AddSiteDialog from '@/ui/molecules/AddSiteDialog';
import { useRecoilValue, useSetRecoilState } from 'recoil';
import { TCommonData, TSnackMessage } from '@/types';
import {
  useGetAllSitesQuery,
  useGetNetworksQuery,
  useAddSiteToNetworkMutation,
} from '@/generated';
import DeleteConfirmation from '@/ui/molecules/DeleteSiteDialog';

export default function Page() {
  const [isAddSiteDialogOpen, setIsAddSiteDialogOpen] = useState(false);
  const setSnackbarMessage = useSetRecoilState<TSnackMessage>(snackbarMessage);
  const [isConfirmationOpen, setIsConfirmationOpen] = useState(false);
  const [siteId, setSiteId] = useState<any>();
  const _commonData = useRecoilValue<TCommonData>(commonData);

  const { data: sitesData, loading: sitesLoading } = useGetAllSitesQuery({
    fetchPolicy: 'cache-and-network',
    variables: {
      networkId: _commonData.networkId,
    },
    onError: (err) => {
      setSnackbarMessage({
        id: 'nodes-msg',
        message: err.message,
        type: 'error',
        show: true,
      });
    },
  });
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
  const handleAddSite = async (values: any, network: string) => {
    addSite({
      variables: {
        networkId: network,
        data: {
          site: values.site,
        },
      },
    });
  };
  const handleDelete = async () => {};
  const handleCloseAction = () => setIsAddSiteDialogOpen(false);
  const handleDeleteSite = (siteId?: string) => {
    setIsConfirmationOpen(true);
    console.log('SITE ID', siteId);
    setSiteId(siteId);
  };
  const handleCancel = () => {
    setIsConfirmationOpen(false);
  };
  const handleAddSiteAction = () => setIsAddSiteDialogOpen(true);

  return (
    <>
      <LoadingWrapper
        radius="small"
        width={'100%'}
        isLoading={sitesLoading || netLoading || addSiteLoading}
        cstyle={{
          backgroundColor: false ? colors.white : 'transparent',
        }}
      >
        <Grid container spacing={2}>
          <Grid item xs={6}>
            <Typography variant="h6"> My sites</Typography>
          </Grid>
          <Grid item xs={6} container justifyContent={'flex-end'}>
            <Button variant="contained" onClick={handleAddSiteAction}>
              ADD SITE
            </Button>
          </Grid>
          {sitesData?.getAllSites.sites.map((site, index) => (
            <Grid item xs={12} key={index} md={6} lg={4}>
              <SiteCard sites={[site]} handleDeleteSite={handleDeleteSite} />
            </Grid>
          ))}
        </Grid>
        <AddSiteDialog
          isOpen={isAddSiteDialogOpen}
          title={'ADD SITE'}
          description=""
          handleCloseAction={handleCloseAction}
          networks={networkList?.getNetworks?.networks ?? []}
          handleAddSite={handleAddSite}
        />
        <DeleteConfirmation
          open={isConfirmationOpen}
          onDelete={handleDelete}
          onCancel={handleCancel}
          itemName={siteId}
          loading={false}
        />
      </LoadingWrapper>
    </>
  );
}
