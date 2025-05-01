/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import React from 'react';
import { Box, Card, CardContent, Grid, Typography } from '@mui/material';
import { SiteDto } from '@/client/graphql/generated';
import { format } from 'date-fns';

interface SiteInfoProps {
  selectedSite: SiteDto;
  address?: string | null;
  nodeIds?: string[];
  createdDate?: string;
}

const SiteInfo: React.FC<SiteInfoProps> = ({
  selectedSite,
  address,
  nodeIds = [],
  createdDate,
}) => {
  const formattedDate = createdDate
    ? format(new Date(createdDate), 'MMMM d, yyyy')
    : selectedSite.installDate
      ? format(new Date(selectedSite.installDate), 'MMMM d, yyyy')
      : 'Not available';

  const formattedCoordinates =
    selectedSite.latitude && selectedSite.longitude
      ? `(${selectedSite.latitude}° N, ${selectedSite.longitude}° W)`
      : '';

  return (
    <Card
      sx={{
        borderRadius: 2,
        boxShadow: '0px 2px 6px rgba(0, 0, 0, 0.05)',
        height: '100%',
        display: 'flex',
        flexDirection: 'column',
      }}
    >
      <CardContent sx={{ padding: 4, flexGrow: 1 }}>
        <Typography variant="h6" sx={{ mb: 4 }}>
          Site information
        </Typography>

        <Grid container spacing={4}>
          <Grid item xs={12} md={4}>
            <Typography
              variant="body2"
              color="text.secondary"
              fontWeight="medium"
            >
              Nodes:
            </Typography>
          </Grid>
          <Grid item xs={12} md={8}>
            {nodeIds && nodeIds.length > 0 ? (
              <Box>
                {nodeIds.map((nodeId, index) => (
                  <Typography key={index} variant="body1" sx={{ mb: 1 }}>
                    {nodeId}
                  </Typography>
                ))}
              </Box>
            ) : (
              <Typography variant="body2">Not available</Typography>
            )}
          </Grid>

          <Grid item xs={12} md={4}>
            <Typography
              variant="body2"
              color="text.secondary"
              fontWeight="medium"
            >
              Date created:
            </Typography>
          </Grid>
          <Grid item xs={12} md={8}>
            <Typography variant="body2">{formattedDate}</Typography>
          </Grid>

          <Grid item xs={12} md={4}>
            <Typography
              variant="body2"
              color="text.secondary"
              fontWeight="medium"
            >
              Location:
            </Typography>
          </Grid>
          <Grid item xs={12} md={8}>
            <Typography variant="body2">
              {address || selectedSite.location || 'Not available'}
            </Typography>
            {formattedCoordinates && (
              <Typography variant="body2" sx={{ mt: 1 }}>
                {formattedCoordinates}
              </Typography>
            )}
          </Grid>
        </Grid>
      </CardContent>
    </Card>
  );
};

export default SiteInfo;
