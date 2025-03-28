/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import React from 'react';
import { Paper, Typography, Stack, Skeleton } from '@mui/material';
import colors from '@/theme/colors';

interface BillingOwnerDetailsCardProps {
  email: string;
}

const BillingOwnerDetailsCard: React.FC<BillingOwnerDetailsCardProps> = ({
  email,
}) => {
  return (
    <Paper
      elevation={2}
      sx={{
        p: 2,
        borderRadius: '10px',
      }}
    >
      <Typography variant="h6" sx={{ color: colors.vulcan, mb: 2 }}>
        Billing owner{' '}
      </Typography>
      <Typography variant="body2" sx={{ color: colors.vulcan, mb: 2 }}>
        Billing owner is responsible for monthly payment.
      </Typography>
      <Stack direction={'column'} spacing={2}>
        <Typography variant="body2" sx={{ color: colors.black54 }}>
          OWNER EMAIL
        </Typography>
        {!email ? (
          <Skeleton variant="text" width="100%" />
        ) : (
          <Typography variant="body2" sx={{ color: colors.vulcan }}>
            {email}
          </Typography>
        )}
      </Stack>
    </Paper>
  );
};

export default BillingOwnerDetailsCard;
