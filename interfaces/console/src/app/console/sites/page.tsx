/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';

import colors from '@/theme/colors';
import LoadingWrapper from '@/components/LoadingWrapper';
import SiteCard from '@/components/SiteCard';
import { Grid, Paper, Typography, Button } from '@mui/material';
import SiteConfigurationStepperDialog from '@/components/SiteConfigurationStepperDialog';
import { useState } from 'react';
import { Dummysites } from '@/constants/stubData';

const Sites = () => {
  const [open, setOpen] = useState(false);

  const handleOpen = () => {
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
  };

  const handleMenuClick = (siteId: string) => {
    console.log(`Menu clicked for siteId: ${siteId}`);
  };
  const handleFormDataSubmit = (formData: any) => {
    console.log('Form data submitted:', formData);
  };
  return (
    <>
      <LoadingWrapper
        radius="small"
        width={'100%'}
        isLoading={false}
        cstyle={{
          backgroundColor: false ? colors.white : 'transparent',
        }}
      >
        <Grid container spacing={0} sx={{ mt: 1 }}>
          <Grid item xs={12}>
            <Paper sx={{ p: 4 }}>
              <Grid container spacing={0} sx={{ mb: 2 }}>
                <Grid xs={6}>
                  <Typography variant="h6" color="initial">
                    My sites
                  </Typography>
                </Grid>
                <Grid
                  xs={6}
                  container
                  justifyItems={'center'}
                  justifyContent={'flex-end'}
                >
                  <Button
                    variant="contained"
                    color="primary"
                    onClick={handleOpen}
                  >
                    ADD SITE
                  </Button>
                </Grid>
              </Grid>

              <Grid container spacing={2}>
                {Dummysites.map((site, index) => (
                  <Grid item xs={12} md={4} lg={4} key={index}>
                    <SiteCard
                      siteId={site.siteId}
                      name={site.name}
                      address={site.address}
                      users={site.users}
                      status={site.status}
                      onClickMenu={handleMenuClick}
                    />
                  </Grid>
                ))}
              </Grid>
            </Paper>
          </Grid>
          <SiteConfigurationStepperDialog
            open={open}
            handleClose={handleClose}
            handleFormDataSubmit={handleFormDataSubmit}
          />
        </Grid>
      </LoadingWrapper>
    </>
  );
};
export default Sites;
