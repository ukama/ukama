'use client';

import React from 'react';
import {
  Tabs,
  Tab,
  Paper,
  Box,
  Typography,
  Divider,
  Stack,
  TextField,
} from '@mui/material';
import LoadingWrapper from '@/components/LoadingWrapper';
import colors from '@/theme/colors';
import PaymentCard from '@/components/PaymentCard';
import { globalUseStyles } from '@/styles/global';
const paymentMock = [
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
const handlePaymentMethodChange = () => {};
const CurrentBilling = () => {
  const gclasses = globalUseStyles();

  return (
    <Box sx={{ py: 2 }}>
      <PaymentCard
        amount="$20.00"
        startDate="06/14/22"
        endDate="07/14/22"
        paymentMethod={paymentMock[0].label}
        onChangePaymentMethod={handlePaymentMethodChange}
        paymentMethods={[
          'American Express - ending in 1234',
          'Visa - ending in 5678',
        ]}
      />
      <Paper
        elevation={2}
        sx={{
          padding: 2,
          marginTop: 2,
          borderRadius: '10px',
          backgroundColor: colors.white,
        }}
      >
        <Typography variant="h6">Data usage</Typography>
        <Stack direction={'row'} justifyContent={'space-between'}>
          <Typography variant="body2" sx={{ color: colors.black54 }}>
            Data usage paid for by subscribers.
          </Typography>
          <Typography variant="body1">$ 0</Typography>
        </Stack>
      </Paper>
      <Paper
        elevation={2}
        sx={{
          padding: 2,
          marginTop: 2,
          borderRadius: '10px',
          backgroundColor: colors.white,
        }}
      >
        <Typography variant="h6">Notification settings</Typography>
        <Stack direction={'row'} justifyContent={'space-between'}>
          <Stack direction={'column'} spacing={2}>
            <Typography variant="body2" sx={{ color: colors.black54 }}>
              All entered emails will receive receipts for the monthly bill.{' '}
            </Typography>
            <TextField
              fullWidth
              label="Email"
              InputLabelProps={{ shrink: true }}
              InputProps={{
                classes: { input: gclasses.inputFieldStyle },
              }}
            />
          </Stack>
        </Stack>
      </Paper>
    </Box>
  );
};

const BillingHistory = () => {
  return (
    <Box sx={{ p: 3 }}>
      <Typography variant="h6">Billing History</Typography>
    </Box>
  );
};

const BillingSettingsPage: React.FC = () => {
  const [currentTab, setCurrentTab] = React.useState(0);

  const handleTabChange = (event: React.SyntheticEvent, newValue: number) => {
    setCurrentTab(newValue);
  };

  return (
    <LoadingWrapper
      width={'100%'}
      radius="medium"
      isLoading={false}
      height={'calc(100vh - 244px)'}
    >
      <Tabs value={currentTab} onChange={handleTabChange}>
        <Tab label="Current Billing" />
        <Tab label="Billing History" />
      </Tabs>
      {currentTab === 0 && <CurrentBilling />}
      {currentTab === 1 && <BillingHistory />}
    </LoadingWrapper>
  );
};

export default BillingSettingsPage;
