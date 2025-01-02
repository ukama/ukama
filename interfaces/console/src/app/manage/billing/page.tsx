/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';
import React, { useState, useCallback } from 'react';
import { Box, Grid } from '@mui/material';
import LoadingWrapper from '@/components/LoadingWrapper';
import {
  useGetReportsQuery,
  useGetPaymentsQuery,
  useUpdatePaymentMutation,
  ReportDto,
} from '@/client/graphql/generated';
import { useAppContext } from '@/context';
import StripePaymentDialog from '@/components/StripePaymentDialog';
import CurrentBillCard from '@/components/CurrentBillCard';
import BillingOwnerDetailsCard from '@/components/BillingOwnerDetailsCard';
import OutStandingBillCard from '@/components/OutStandingBillCard';
import BillingHistoryTable from '@/components/BillingHistory';

const BillingSettingsPage: React.FC = () => {
  const { setSnackbarMessage, user } = useAppContext();
  const [isPaymentDialogOpen, setIsPaymentDialogOpen] = useState(false);
  const [extraKey, setExtraKey] = useState<string>('');
  const [myBill, setMyBill] = useState<ReportDto>();

  const {
    data: reportsData,
    loading: reportsLoading,
    refetch: refetchReports,
  } = useGetReportsQuery({
    variables: {
      data: {
        report_type: 'invoice',
        count: 100,
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
        type: 'error',
        show: true,
      });
    },
  });

  const {
    data: paymentsData,
    loading: paymentsLoading,
    refetch: refetchPayments,
  } = useGetPaymentsQuery({
    variables: {
      data: {
        type: 'invoice',
      },
    },
    fetchPolicy: 'network-only',
    onError: (error) => {
      setSnackbarMessage({
        id: 'payments-error',
        message: error.message,
        type: 'error',
        show: true,
      });
    },
  });

  const [updatePayment] = useUpdatePaymentMutation({
    onError: (error) => {
      setSnackbarMessage({
        id: 'update-payment-error',
        message: error.message,
        type: 'error',
        show: true,
      });
    },
  });

  const handleAddPayment = useCallback(
    async (billId: string) => {
      const currentPayment = paymentsData?.getPayments?.payments.find(
        (payment) => payment.itemId === billId,
      );
      const myBill = reportsData?.getReports?.reports.find(
        (bill) => bill.id === billId,
      );
      setMyBill(myBill);

      if (!currentPayment) {
        setSnackbarMessage({
          id: 'payment-error',
          message: 'Payment not found for this bill.',
          type: 'error',
          show: true,
        });
        return;
      }

      try {
        if (currentPayment?.extra) {
          setExtraKey(currentPayment.extra);
          setIsPaymentDialogOpen(true);
        } else {
          const currentPaymentId = currentPayment?.id;
          if (!currentPaymentId) {
            setSnackbarMessage({
              id: 'payment-error',
              message: 'Payment ID is missing.',
              type: 'error',
              show: true,
            });
            return;
          }

          const result = await updatePayment({
            variables: {
              data: {
                id: currentPaymentId,
                paymentMethod: 'stripe',
                payerEmail: user?.email,
                payerName: user?.name,
              },
            },
          });

          if (result.errors) {
            setSnackbarMessage({
              id: 'payment-error',
              message: result.errors[0].message,
              type: 'error',
              show: true,
            });
            return;
          }

          if (result.data) {
            const updatedPaymentsData = await refetchPayments();

            const updatedPaymentSection =
              updatedPaymentsData?.data?.getPayments?.payments.find(
                (payment) => payment.itemId === billId,
              );

            if (updatedPaymentSection?.extra) {
              setExtraKey(updatedPaymentSection.extra);
              setIsPaymentDialogOpen(true);
            }
          }
        }
      } catch (error) {
        setSnackbarMessage({
          id: 'payment-error',
          message: error instanceof Error ? error.message : 'Unknown error',
          type: 'error',
          show: true,
        });
      }
    },
    [
      paymentsData,
      reportsData,
      updatePayment,
      refetchPayments,
      user,
      setSnackbarMessage,
    ],
  );

  const handlePaymentSuccess = async () => {
    try {
      setSnackbarMessage({
        id: 'payment-success',
        message: 'Payment completed successfully',
        type: 'success',
        show: true,
      });

      await Promise.all([refetchReports(), refetchPayments()]);

      setIsPaymentDialogOpen(false);
      setExtraKey('');
      setMyBill(undefined);
    } catch (error) {
      setSnackbarMessage({
        id: 'payment-error',
        message: error instanceof Error ? error.message : 'Unknown error',
        type: 'error',
        show: true,
      });
    }
  };

  const billingHistoryDataset = reportsData?.getReports?.reports.map(
    (report) => ({
      id: report.id,
      posted: new Date(report.createdAt).toLocaleDateString(),
      billing: report.rawReport?.subscriptions
        .map(
          (sub) =>
            `${new Date(sub.startedAt).toLocaleDateString()} - ${
              sub.terminatedAt
                ? new Date(sub.terminatedAt).toLocaleDateString()
                : 'Present'
            }`,
        )
        .join(', '),
      payment: report.isPaid ? 'Paid' : 'Unpaid',
      description: report.rawReport?.fees[0].item.name,
      pdf: report.rawReport?.fileUrl,
    }),
  );

  return (
    <LoadingWrapper
      width="100%"
      isLoading={reportsLoading || paymentsLoading}
      height="calc(100vh - 244px)"
    >
      <Box sx={{ width: '100%' }}>
        <Grid container spacing={2}>
          <Grid item xs={12}>
            <CurrentBillCard
              bills={reportsData?.getReports?.reports.filter(
                (report) => !report.isPaid,
              )}
              onPay={handleAddPayment}
              isLoading={reportsLoading || paymentsLoading}
            />
          </Grid>
          <Grid item xs={12}>
            <OutStandingBillCard
              reports={reportsData?.getReports?.reports.filter(
                (report) => !report.isPaid,
              )}
              onPaySingle={handleAddPayment}
            />
          </Grid>
          <Grid item xs={12}>
            <BillingOwnerDetailsCard email={user?.email} />
          </Grid>
          <Grid item xs={12}>
            {billingHistoryDataset && (
              <BillingHistoryTable data={billingHistoryDataset} />
            )}
          </Grid>
        </Grid>
      </Box>
      {extraKey && (
        <StripePaymentDialog
          open={isPaymentDialogOpen}
          onClose={() => {
            setIsPaymentDialogOpen(false);
            setExtraKey('');
            setMyBill(undefined);
          }}
          extraKey={extraKey}
          bill={myBill}
          onPaymentSuccess={handlePaymentSuccess}
        />
      )}
    </LoadingWrapper>
  );
};

export default BillingSettingsPage;
