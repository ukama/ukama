/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import React from 'react';
import {
  Typography,
  Paper,
  Stack,
  Divider,
  Grid,
  Button,
  Skeleton,
} from '@mui/material';
import colors from '@/theme/colors';

interface BillCardProps {
  amount?: string;
  startDate?: string;
  endDate?: string;
  isLoading?: boolean;
  onPay?: () => void;
}

const CurrentBillCard: React.FC<BillCardProps> = ({
  amount = '',
  startDate = '',
  endDate = '',
  isLoading = false,
  onPay,
}) => {
  if (isLoading) {
    return (
      <Grid container spacing={3}>
        <Grid item xs={12}>
          <Paper
            elevation={3}
            sx={{
              padding: 3,
              borderRadius: 2,
              backgroundColor: colors.white,
              boxShadow: '0 4px 6px rgba(0,0,0,0.1)',
            }}
          >
            <Stack spacing={2}>
              <Stack
                direction="row"
                spacing={2}
                justifyContent={'space-between'}
              >
                <Skeleton variant="text" width="30%" />
                <Skeleton variant="text" width="20%" />
              </Stack>
              <Skeleton variant="text" width="100%" />
              <Divider />
              <Stack
                direction="row"
                spacing={2}
                justifyContent={'space-between'}
              >
                <Skeleton variant="text" width="40%" />
                <Skeleton variant="rectangular" width={120} height={40} />
              </Stack>
            </Stack>
          </Paper>
        </Grid>
      </Grid>
    );
  }

  return (
    <Grid container spacing={3}>
      <Grid item xs={12}>
        <Paper
          elevation={3}
          sx={{
            padding: 3,
            borderRadius: 2,
            backgroundColor: colors.white,
            boxShadow: '0 4px 6px rgba(0,0,0,0.1)',
          }}
        >
          <Stack spacing={2}>
            <Stack direction="row" spacing={2} justifyContent={'space-between'}>
              <Typography variant="h6">Next bill </Typography>
              <Typography variant="body2">
                {startDate} - {endDate}
              </Typography>
            </Stack>
            <Typography variant="body2" sx={{ color: colors.black54 }}>
              Monthly SaaS fee for Ukama Console. Due EOD {startDate}
            </Typography>
            <Divider />
            <Stack direction="row" spacing={2} justifyContent={'space-between'}>
              <Typography variant="h5">Total due: $ {amount}</Typography>
              <Button
                variant="contained"
                size="large"
                onClick={onPay}
                disabled={!onPay}
              >
                Pay now
              </Button>
            </Stack>
          </Stack>
        </Paper>
      </Grid>
    </Grid>
  );
};

export default CurrentBillCard;
