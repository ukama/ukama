/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import React from 'react';
import { Box, Grid, Typography } from '@mui/material';
import SiteCard from '@/components/SiteCard';
import { SiteDto } from '@/client/graphql/generated';
import LoadingWrapper from '@/components/LoadingWrapper';

interface SitesWrapperProps {
  sites: SiteDto[];
  loading: boolean;
  handleAddSite?: () => void;
  handleSiteNameUpdate: (siteId: string, siteName: string) => void;
}

const SitesWrapper: React.FC<SitesWrapperProps> = ({
  sites,
  loading,
  handleSiteNameUpdate,
}) => {
  if (sites?.length === 0 && !loading) {
    return (
      <Box
        sx={{
          height: '100%',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          flexDirection: 'column',
          gap: 2,
          padding: '0 20px',
          textAlign: 'center',
        }}
      >
        <Typography variant="h6" color="textSecondary">
          No sites yet!
        </Typography>
        <Typography
          variant="body1"
          color="textSecondary"
          sx={{ maxWidth: '450px' }}
        >
          A site is a complete connection point to the network, made up of your
          Ukama node, and the power and backhaul components. Install a site to
          get started.
        </Typography>
      </Box>
    );
  }

  return (
    <LoadingWrapper isLoading={loading} height="100%">
      <Box
        sx={{
          height: '100%',
          overflowY: 'auto',
          padding: '10px',
        }}
      >
        <Grid container spacing={2}>
          {sites?.map((site) => {
            return (
              <Grid item xs={12} md={4} lg={4} key={site.id}>
                <SiteCard
                  siteId={site.id}
                  name={site.name}
                  address={site.location}
                  loading={loading}
                  handleSiteNameUpdate={handleSiteNameUpdate}
                />
              </Grid>
            );
          })}
        </Grid>
      </Box>
    </LoadingWrapper>
  );
};

export default SitesWrapper;
