/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import React from 'react';
import { Box, Stack, Typography } from '@mui/material';
import { SiteOverViewPlaceHolder } from '../../../public/svg';

interface SiteOverviewProps {
  inputPower?: string;
  solarStorage?: string;
  consumption?: string;
}

const SiteOverview: React.FC<SiteOverviewProps> = ({
  inputPower = 'N/A',
  solarStorage = 'N/A',
  consumption = 'N/A',
}) => {
  return (
    <Box sx={{ p: 2 }}>
      <Typography variant="body1" fontWeight="Bold" sx={{ mb: 1 }}>
        Site Overview
      </Typography>

      <Stack direction="row" spacing={4} alignItems="center" sx={{ mb: 2 }}>
        <Box
          sx={{
            display: 'flex',
            alignItems: 'center',
            gap: 1,
          }}
        >
          <Box
            sx={{
              width: 10,
              height: 10,
              backgroundColor: 'grey.300',
              borderRadius: '2px',
            }}
          />
          <Typography variant="caption">Input Power: {inputPower}</Typography>
        </Box>

        <Box
          sx={{
            display: 'flex',
            alignItems: 'center',
            gap: 1,
          }}
        >
          <Box
            sx={{
              width: 10,
              height: 10,
              backgroundColor: 'grey.300',
              borderRadius: '2px',
            }}
          />
          <Typography variant="caption">
            Solar Storage: {solarStorage}
          </Typography>
        </Box>

        <Box
          sx={{
            display: 'flex',
            alignItems: 'center',
            gap: 1,
          }}
        >
          <Box
            sx={{
              width: 10,
              height: 10,
              backgroundColor: 'grey.300',
              borderRadius: '2px',
            }}
          />
          <Typography variant="caption">Consumption: {consumption}</Typography>
        </Box>
      </Stack>

      <Box sx={{ display: 'flex', justifyContent: 'center', mt: 2 }}>
        <SiteOverViewPlaceHolder />
      </Box>
    </Box>
  );
};

export default SiteOverview;
