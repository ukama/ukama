/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import React from 'react';
import {
  Box,
  Paper,
  Stack,
  Typography,
  Button,
  Divider,
  useMediaQuery,
  useTheme,
} from '@mui/material';
import colors from '@/theme/colors';

interface OutStandingBillCardProps {
  totalAmount: string;
  loading?: boolean;
}

const OutStandingBillCard: React.FC<OutStandingBillCardProps> = ({
  loading = false,
}) => {
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('sm'));

  return (
    <Box>
      <Paper
        elevation={2}
        sx={{
          p: isMobile ? 2 : 4,
          borderRadius: '10px',
        }}
      >
        <Stack
          direction={isMobile ? 'column' : 'row'}
          justifyContent="space-between"
          alignItems={isMobile ? 'flex-start' : 'center'}
          spacing={isMobile ? 2 : 0}
          sx={{ mb: 2 }}
        >
          <Typography
            variant="h6"
            sx={{
              fontSize: isMobile ? '1.1rem' : '1.25rem',
              mb: isMobile ? 1 : 0,
            }}
          >
            Outstanding bills
          </Typography>
          <Button
            variant="contained"
            fullWidth={isMobile}
            sx={{
              py: isMobile ? 1 : 'auto',
              px: isMobile ? 2 : 'auto',
            }}
          >
            Pay all outstanding bills
          </Button>
        </Stack>
        <Typography
          variant="body2"
          sx={{
            color: colors.vulcan,
            mb: 2,
            fontSize: isMobile ? '0.8rem' : '0.875rem',
          }}
        >
          Overdue bills for Ukama Console plan.
        </Typography>
        <Divider sx={{ mb: 2 }} />
        <Stack
          direction={isMobile ? 'column' : 'row'}
          spacing={2}
          justifyContent={'space-between'}
          alignItems={isMobile ? 'stretch' : 'center'}
        >
          <Typography
            variant="body2"
            sx={{
              color: colors.vulcan,
              fontSize: isMobile ? '0.8rem' : '0.875rem',
              textAlign: isMobile ? 'left' : 'inherit',
            }}
          >
            Total due: $20.00 Overdue on 10/05/25
          </Typography>

          <Button
            variant="contained"
            fullWidth={isMobile}
            sx={{
              py: isMobile ? 1 : 'auto',
              px: isMobile ? 2 : 'auto',
              mt: isMobile ? 1 : 0,
            }}
          >
            Pay now
          </Button>
        </Stack>
      </Paper>
    </Box>
  );
};

export default OutStandingBillCard;
