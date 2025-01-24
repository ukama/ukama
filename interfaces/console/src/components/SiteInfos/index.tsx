/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import React from 'react';
import { Box, Stack, Typography } from '@mui/material';
import { SiteDto } from '@/client/graphql/generated';
import { format } from 'date-fns';

interface SiteInfoProps {
  selectedSite: SiteDto;
  address: string | null;
  nodes: any[];
}

const SiteInfo: React.FC<SiteInfoProps> = ({
  selectedSite,
  address,
  nodes,
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
          <Typography variant="subtitle2">Date Created:</Typography>
          <Typography variant="subtitle2">{formattedDate}</Typography>
        </Stack>

        <Stack direction="row" spacing={2} justifyItems="center">
          <Typography variant="subtitle2">Location:</Typography>
          <Typography variant="subtitle2">
            {address ||
              `${selectedSite.location} (${selectedSite.latitude}, ${selectedSite.longitude})` ||
              'N/A'}
            ({selectedSite.latitude}, {selectedSite.longitude})
          </Typography>
        </Stack>

        <Stack direction="row" spacing={2} justifyItems="center">
          <Typography variant="subtitle2">Node:</Typography>
          <Stack direction={'column'} spacing={1}>
            {nodes.map((n, index) => (
              <Typography key={index} variant="subtitle2">
                #{n.id}
              </Typography>
            ))}
          </Stack>
        </Stack>
      </Stack>
    </Box>
  );
};

export default SiteInfo;