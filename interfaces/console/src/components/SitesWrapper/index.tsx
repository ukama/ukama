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
  sitesStatus?: Record<
    string,
    {
      status: string;
      batteryStatus: string;
      signalStrength: string;
    }
  >;
  handleAddSite?: () => void;
  handleSiteNameUpdate: (siteId: string, siteName: string) => void;
}

const SitesWrapper: React.FC<SitesWrapperProps> = ({
  sites,
  loading,
  sitesStatus = {},
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
        }}
      >
        <Typography variant="h6" color="textSecondary">
          No sites found
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
            const siteStatus = sitesStatus[site.id] || {
              status: 'Online',
              batteryStatus: 'Charged',
              signalStrength: 'Strong',
            };

            return (
              <Grid item xs={12} md={4} lg={4} key={site.id}>
                <SiteCard
                  siteId={site.id}
                  name={site.name}
                  address={site.location}
                  connectionStatus={siteStatus.status}
                  batteryStatus={siteStatus.batteryStatus}
                  signalStrength={siteStatus.signalStrength}
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
