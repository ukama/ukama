/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';

import React, { useState, useMemo, useEffect } from 'react';
import { Tabs, Tab, AlertColor, Paper } from '@mui/material';
import LoadingWrapper from '@/components/LoadingWrapper';
import CurrentBill from '@/components/CurrentBill';
import {
  useGetReportsQuery,
  useGetPaymentsQuery,
} from '@/client/graphql/generated';
import { useAppContext } from '@/context';
import BillingHistory from '@/components/BillHistoryTab';
import StripePaymentDialog from '@/components/StripePaymentDialog';
import { format } from 'date-fns';

const BillingSettingsPage: React.FC = () => {
  const [currentTab, setCurrentTab] = useState(0);
  const { setSnackbarMessage, user } = useAppContext();
  const [isPaymentDialogOpen, setIsPaymentDialogOpen] = useState(false);
  const [clientSecret, setClientSecret] = useState<string>('');
  const [paymentLoading, setPaymentLoading] = useState(false);

  const handleAddPaymentMethod = () => {
    const nextPaymentAmount =
      reportsData?.getReports?.reports[0]?.rawReport?.amountCents || 0;

    setPaymentLoading(true);
  };

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
    console.log('BRACKLEY :', error);
    setSnackbarMessage({
      id: 'payment-error',
      message: error.message || 'Payment failed',
      type: 'error' as AlertColor,
      show: true,
    });
  };

  const handleViewPaymentDetails = (paymentId: string) => {
    console.log(`Viewing details for payment ${paymentId}`);
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

  const totalAmountUSD = useMemo(() => {
    return (
      paymentsData?.getPayments.payments.reduce((total, payment) => {
        return total + parseFloat(payment.amount);
      }, 0) || 0
    );
  }, [paymentsData]);

  const handleTabChange = (_: React.SyntheticEvent, newValue: number) => {
    setCurrentTab(newValue);
  };
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
      <Tabs value={currentTab} onChange={handleTabChange}>
        <Tab label="Current Billing" />
        <Tab label="Billing History" />
      </Tabs>
      {currentTab === 0 && (
        <CurrentBill
          dataUsagePaid={0}
          notificationEmail={user.email}
          nextPaymentAmount={
            paymentsData?.getPayments?.payments &&
            paymentsData.getPayments.payments.length > 0
              ? Number(paymentsData.getPayments.payments[0]?.amount) || 0
              : 0
          }
          nextPaymentDate={
            paymentsData?.getPayments?.payments &&
            paymentsData.getPayments.payments.length > 0 &&
            paymentsData.getPayments.payments[0]?.createdAt
              ? format(
                  new Date(paymentsData.getPayments.payments[0].createdAt),
                  'dd MMM yyyy',
                )
              : 'N/A'
          }
          loading={paymentsLoading}
          onAddPaymentMethod={handleAddPaymentMethod}
          handleAddPayment={handleAddPayment}
        />
      )}

      {currentTab === 1 && (
        <Paper
          sx={{
            height: '100%',
            borderRadius: '10px',
            px: { xs: 2, md: 3 },
            py: { xs: 2, md: 4 },
          }}
        >
          <BillingHistory
            bills={reportsData?.getReports?.reports || []}
            loading={paymentsLoading}
            onViewDetails={handleViewPaymentDetails}
          />
        </Paper>
      )}
      <StripePaymentDialog
        open={isPaymentDialogOpen}
        onClose={() => setIsPaymentDialogOpen(false)}
        clientSecret={clientSecret}
        amount={parseFloat(
          paymentsData?.getPayments?.payments?.[0]?.amount ?? '0',
        )}
        // loading={paymentLoading}
        onPaymentSuccess={handlePaymentSuccess}
        onPaymentError={handlePaymentError}
      />
    </LoadingWrapper>
  );
};

export default BillingSettingsPage;
