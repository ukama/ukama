/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import React from 'react';
import { Box, Paper, Stack, Typography, Button } from '@mui/material';
import colors from '@/theme/colors';

interface FeatureUsageCardProps {
  upTime: string;
  totalDataUsage: string;
  ActiveSubscriberCount: number;
  loading?: boolean;
}

const FeatureUsageCard: React.FC<FeatureUsageCardProps> = ({
  upTime,
  totalDataUsage,
  ActiveSubscriberCount,
  loading = false,
}) => {
  return (
    <Box>
      <Paper
        elevation={2}
        sx={{
          p: 2,
          borderRadius: '10px',
        }}
      >
        <Stack
          direction="row"
          justifyContent="space-between"
          alignItems="center"
          sx={{ mb: 2 }}
        >
          <Typography variant="h6">Feature usage</Typography>
        </Stack>

        <Stack direction={'column'} spacing={2}>
          <Stack direction={'row'} spacing={2} justifyContent={'space-between'}>
            <Typography variant="body2" sx={{ color: colors.vulcan }}>
              Up Time
            </Typography>
            <Typography variant="body2" sx={{ color: colors.vulcan }}>
              {upTime} %
            </Typography>
          </Stack>
          <Stack direction={'row'} spacing={2} justifyContent={'space-between'}>
            <Typography variant="body2" sx={{ color: colors.vulcan }}>
              Total Data Usage
            </Typography>
            <Typography variant="body2" sx={{ color: colors.vulcan }}>
              {totalDataUsage} GB
            </Typography>
          </Stack>
          <Stack direction={'row'} spacing={2} justifyContent={'space-between'}>
            <Typography variant="body2" sx={{ color: colors.vulcan }}>
              Active Subscribers
            </Typography>
            <Typography variant="body2" sx={{ color: colors.vulcan }}>
              {ActiveSubscriberCount}
            </Typography>
          </Stack>
        </Stack>
      </Paper>
    </Box>
  );
};

export default FeatureUsageCard;
