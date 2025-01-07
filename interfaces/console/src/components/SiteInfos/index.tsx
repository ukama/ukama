/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import React from 'react';
import { Box, Paper, Stack, Typography } from '@mui/material';
import { SiteDto } from '@/client/graphql/generated';
import { format } from 'date-fns';

interface SiteInfoProps {
  selectedSite: SiteDto;
  address: string | null;
  nodeId: string;
}

const SiteInfo: React.FC<SiteInfoProps> = ({
  selectedSite,
  address,
  nodeId,
}) => {
  const formattedDate = selectedSite.createdAt
    ? format(new Date(selectedSite.createdAt), 'MMMM d, yyyy')
    : 'N/A';
  return (
    <Box
      sx={{
        px: 3,
        borderRadius: '5px',
        py: 2,
      }}
    >
      <Stack direction="column" spacing={2}>
        <Typography variant="body1" fontWeight="Bold">
          Site Information
        </Typography>

        <Stack direction="row" spacing={2} justifyItems="center">
          <Typography variant="subtitle1">Date Created:</Typography>
          <Typography variant="subtitle1">{formattedDate}</Typography>
        </Stack>

        <Stack direction="row" spacing={2} justifyItems="center">
          <Typography variant="subtitle1">Location:</Typography>
          <Typography variant="subtitle1">
            {address ||
              `${selectedSite.location} (${selectedSite.latitude}, ${selectedSite.longitude})` ||
              'N/A'}
            ({selectedSite.latitude}, {selectedSite.longitude})
          </Typography>
        </Stack>

        <Stack direction="row" spacing={2} justifyItems="center">
          <Typography variant="subtitle1">Node:</Typography>
          <Typography variant="subtitle1">{nodeId}</Typography>
        </Stack>
      </Stack>
    </Box>
  );
};

export default SiteInfo;
