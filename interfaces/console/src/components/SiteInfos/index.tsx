/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import React from 'react';
import { Paper, Stack, Typography } from '@mui/material';
import { SiteDto } from '@/client/graphql/generated';

interface SiteInfoProps {
  selectedSite: SiteDto;
  address: string | null;
}

const SiteInfo: React.FC<SiteInfoProps> = ({ selectedSite, address }) => {
  return (
    <Paper
      elevation={3}
      sx={{
        p: 2,
        flex: 1,
        height: '100%',
        borderRadius: '5px',
        position: 'relative',
      }}
    >
      <Stack direction="column" spacing={2}>
        <Typography variant="h6">Site Information</Typography>
        <Stack direction="row" spacing={4} justifyItems={'center'}>
          <Typography variant="subtitle1">Location:</Typography>
          <Typography variant="subtitle1">{selectedSite.location}</Typography>
        </Stack>
        <Stack direction="row" spacing={4} justifyItems={'center'}>
          <Typography variant="subtitle1">Coordinates:</Typography>
          <Typography variant="subtitle1">
            ( {selectedSite.latitude}, {selectedSite.longitude} )
          </Typography>
        </Stack>
        <Stack direction="row" spacing={4} justifyItems={'center'}>
          <Typography variant="subtitle1">Address:</Typography>
          <Typography variant="subtitle1">{address}</Typography>
        </Stack>
      </Stack>
    </Paper>
  );
};

export default SiteInfo;
