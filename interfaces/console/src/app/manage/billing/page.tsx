/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';

import React, { useState, useEffect } from 'react';
import { AlertColor, Box, Grid, Paper } from '@mui/material';
import LoadingWrapper from '@/components/LoadingWrapper';
import {
  useGetReportsQuery,
  useGetPaymentsQuery,
} from '@/client/graphql/generated';
import { useAppContext } from '@/context';
import StripePaymentDialog from '@/components/StripePaymentDialog';
import { format } from 'date-fns';
import CurrentBillCard from '@/components/CurrentBillCard';
import FeatureUsageCard from '@/components/FeatureUsage';
import BillingOwnerDetailsCard from '@/components/BillingOwnerDetailsCard';
import BillingHistory from '@/components/BillHistoryTab';
import OutStandingBillCard from '@/components/OutStandingBillCard';

const BillingSettingsPage: React.FC = () => {
  const { setSnackbarMessage, user } = useAppContext();
  const [isPaymentDialogOpen, setIsPaymentDialogOpen] = useState(false);
  const [clientSecret, setClientSecret] = useState<string>('');

  const handlePaymentSuccess = () => {
    setSnackbarMessage({
      id: 'payment-success',
      message: 'Payment completed successfully',
      type: 'success' as AlertColor,
      show: true,
    });
    setIsPaymentDialogOpen(false);
  };

  const handlePaymentError = (error: any) => {
    setSnackbarMessage({
      id: 'payment-error',
      message: error.message || 'Payment failed',
      type: 'error' as AlertColor,
      show: true,
    });
  };

  const { data: reportsData, loading: reportsLoading } = useGetReportsQuery({
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
        id: 'reports-error',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });

  const { data: paymentsData, loading: paymentsLoading } = useGetPaymentsQuery({
    variables: {
      data: {
        paymentMethod: 'stripe',
        status: 'processing',
        type: 'invoice',
      },
    },
    fetchPolicy: 'network-only',
    onError: (error) => {
      setSnackbarMessage({
        id: 'reports-error',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });
  useEffect(() => {
    if (
      paymentsData?.getPayments?.payments &&
      paymentsData.getPayments.payments.length > 0
    ) {
      const firstPayment = paymentsData.getPayments.payments[0];
      const extractedSecret = firstPayment.extra ?? '';
      console.log('Extracted Client Secret:', extractedSecret);
      setClientSecret(extractedSecret);
    }
  }, [paymentsData]);

  const handleAddPayment = () => {
    setIsPaymentDialogOpen(true);
  };
  return (
    <LoadingWrapper
      width="100%"
      radius="medium"
      isLoading={paymentsLoading || reportsLoading}
      height="calc(100vh - 244px)"
    >
      <Grid container spacing={3}>
        <Grid item xs={12}>
          <CurrentBillCard
            amount={
              paymentsData?.getPayments?.payments &&
              paymentsData?.getPayments?.payments.length &&
              paymentsData?.getPayments?.payments[0]?.amount
                ? paymentsData?.getPayments?.payments[0]?.amount
                : '0'
            }
            startDate={
              paymentsData?.getPayments.payments[0]?.createdAt
                ? format(
                    new Date(paymentsData.getPayments.payments[0].createdAt),
                    'dd MMM yyyy',
                  )
                : 'N/A'
            }
            endDate={
              paymentsData?.getPayments.payments[0]?.createdAt
                ? format(
                    new Date(paymentsData.getPayments.payments[0].createdAt),
                    'dd MMM yyyy',
                  )
                : 'N/A'
            }
            onPay={handleAddPayment}
            isLoading={false}
          />
        </Grid>
        <Grid item xs={12}>
          <OutStandingBillCard totalAmount={''} />
        </Grid>
        <Grid item xs={12} md={6}>
          <FeatureUsageCard
            upTime={'12'}
            totalDataUsage={'12'}
            ActiveSubscriberCount={40}
            loading={false}
          />
        </Grid>
        <Grid item xs={12} md={6}>
          <BillingOwnerDetailsCard loading={false} name={user?.name} />
        </Grid>
        <Grid item xs={12}>
          <BillingHistory bills={[]} />
        </Grid>
      </Grid>

      <StripePaymentDialog
        open={isPaymentDialogOpen}
        onClose={() => setIsPaymentDialogOpen(false)}
        clientSecret={clientSecret}
        amount={parseFloat(
          paymentsData?.getPayments?.payments?.[0]?.amount ?? '0',
        )}
        onPaymentSuccess={handlePaymentSuccess}
        onPaymentError={handlePaymentError}
      />
    </LoadingWrapper>
  );
};

export default BillingSettingsPage;
