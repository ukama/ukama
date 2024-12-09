/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';

import React, { useState } from 'react';
import {
  Tabs,
  Tab,
  Paper,
  Box,
  Typography,
  Stack,
  TextField,
  AlertColor,
} from '@mui/material';
import LoadingWrapper from '@/components/LoadingWrapper';
import colors from '@/theme/colors';
import PaymentCard from '@/components/PaymentCard';
import { globalUseStyles } from '@/styles/global';
import {
  useGetReportsQuery,
  useGetPaymentsQuery,
} from '@/client/graphql/generated';
import { useAppContext } from '@/context';

const PAYMENT_METHODS = [
  {
    value: 'no_payment_method_Set',
    label: 'No payment method set',
  },
  {
    value: 'stripe',
    label: 'Stripe',
  },
  {
    value: 'paypal',
    label: 'PayPal',
  },
];

const CREDIT_CARDS = [
  'American Express - ending in 1234',
  'Visa - ending in 5678',
];

interface currentBillProps {
  packagePaid: number;
}

const CurrentBilling: React.FC<currentBillProps> = () => {
  const gclasses = globalUseStyles();

  const handlePaymentMethodChange = () => {
    // Implement payment method change logic
  };

  return (
    <Box sx={{ py: 2 }}>
      <PaymentCard
        amount="$20.00"
        startDate="06/14/22"
        endDate="07/14/22"
        paymentMethod={PAYMENT_METHODS[0].label}
        onChangePaymentMethod={handlePaymentMethodChange}
        paymentMethods={CREDIT_CARDS}
      />

      <Paper
        elevation={2}
        sx={{
          p: 4,
          mt: 2,
          borderRadius: '10px',
          bgcolor: colors.white,
        }}
      >
        <Typography variant="h6">Data usage</Typography>
        <Stack
          direction="row"
          justifyContent="space-between"
          alignItems="center"
        >
          <Typography variant="body2" sx={{ color: colors.black54 }}>
            Data usage paid for by subscribers
          </Typography>
          <Typography variant="body1">$ 0</Typography>
        </Stack>
      </Paper>

      <Paper
        elevation={2}
        sx={{
          p: 4,
          mt: 2,
          borderRadius: '10px',
          bgcolor: colors.white,
        }}
      >
        <Typography variant="h6">Notification settings</Typography>
        <Stack direction="column" spacing={2}>
          <Typography variant="body2" sx={{ color: colors.black54 }}>
            All entered emails will receive receipts for the monthly bill
          </Typography>
          <TextField
            fullWidth
            label="PRIMARY EMAIL"
            InputLabelProps={{ shrink: true }}
            InputProps={{
              classes: { input: gclasses.inputFieldStyle },
            }}
          />
        </Stack>
      </Paper>
    </Box>
  );
};

const BillingHistory: React.FC = () => (
  <Box sx={{ p: 3 }}>
    <Typography variant="h6">Billing History</Typography>
  </Box>
);

const BillingSettingsPage: React.FC = () => {
  const [currentTab, setCurrentTab] = useState(0);
  const { setSnackbarMessage, network, env } = useAppContext();

  const handleTabChange = (_: React.SyntheticEvent, newValue: number) => {
    setCurrentTab(newValue);
  };

  const { data } = useGetPaymentsQuery({
    variables: {
      data: {
        paymentMethod: 'stripe',
        status: 'completed',
        type: 'package',
      },
    },
    fetchPolicy: 'network-only',
    onError: (error) => {
      setSnackbarMessage({
        id: 'sims-error-msg',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });

  const { data: reports } = useGetReportsQuery({
    variables: {
      data: {
        isPaid: false,
        report_type: 'invoice',
        count: 0,
        networkId: '',
        ownerId: '',
        ownerType: '',
        sort: false,
      },
    },
    fetchPolicy: 'network-only',
    onError: (error) => {
      setSnackbarMessage({
        id: 'reports-error-msg',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });

  const totalAmountUSD = data?.getPayments.payments.reduce((total, payment) => {
    return total + parseFloat(payment.amount) / 100;
  }, 0);

  return (
    <LoadingWrapper
      width="100%"
      radius="medium"
      isLoading={false}
      height="calc(100vh - 244px)"
    >
      <Tabs value={currentTab} onChange={handleTabChange}>
        <Tab label="Current Billing" />
        <Tab label="Billing History" />
      </Tabs>

      {currentTab === 0 && (
        <CurrentBilling
          packagePaid={
            totalAmountUSD ? parseFloat(totalAmountUSD.toFixed(2)) : 0
          }
        />
      )}
      {currentTab === 1 && <BillingHistory />}
    </LoadingWrapper>
  );
};

export default BillingSettingsPage;
