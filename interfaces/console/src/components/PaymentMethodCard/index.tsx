/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */ import React from 'react';
import { Paper, Typography, Stack, Button } from '@mui/material';
import AddCardIcon from '@mui/icons-material/AddCard';
import colors from '@/theme/colors';

interface PaymentMethodCardProps {
  onAddPaymentMethod?: () => void;
}

const PaymentMethodCard: React.FC<PaymentMethodCardProps> = ({
  onAddPaymentMethod,
}) => {
  return (
    <Paper
      elevation={2}
      sx={{
        p: 4,
        mt: 2,
        borderRadius: '10px',
        bgcolor: colors.white,
      }}
    >
      <Stack
        direction="row"
        justifyContent="space-between"
        alignItems="center"
        sx={{ mb: 2 }}
      >
        <Typography variant="h6">Payment methods</Typography>
      </Stack>

      <Typography variant="body2" sx={{ color: colors.black54 }}>
        Payments for monthly SaaS fees will be made with default card. There are
        currently no integrations for auto-payment.{' '}
      </Typography>
      <Button
        variant="outlined"
        color="primary"
        startIcon={<AddCardIcon />}
        onClick={onAddPaymentMethod}
        disabled={!onAddPaymentMethod}
        sx={{ mt: 2 }}
      >
        Add Payment Method
      </Button>
    </Paper>
  );
};

export default PaymentMethodCard;
