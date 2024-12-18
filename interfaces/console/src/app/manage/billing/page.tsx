'use client';

import React, { useState, useEffect, useMemo } from 'react';
import { AlertColor, Box, Grid, Paper, Tabs, Tab } from '@mui/material';
import LoadingWrapper from '@/components/LoadingWrapper';
import {
  useGetReportsQuery,
  useGetPaymentsQuery,
  GetReportResDto,
} from '@/client/graphql/generated';
import { useAppContext } from '@/context';
import StripePaymentDialog from '@/components/StripePaymentDialog';
import CurrentBillCard from '@/components/CurrentBillCard';
import FeatureUsageCard from '@/components/FeatureUsage';
import BillingOwnerDetailsCard from '@/components/BillingOwnerDetailsCard';
import OutStandingBillCard from '@/components/OutStandingBillCard';
import DataTableWithOptions from '@/components/DataTableWithOptions';
import { BILLING_HISTORY_TABLE_MENU, BILLING_TABLE_COLUMNS } from '@/constants';
import SubscriberIcon from '@mui/icons-material/PeopleAlt';

const BillingSettingsPage: React.FC = () => {
  const { setSnackbarMessage, user } = useAppContext();
  const [isPaymentDialogOpen, setIsPaymentDialogOpen] = useState(false);
  const [clientSecret, setClientSecret] = useState<string>('');
  const [tabValue, setTabValue] = useState(0);

  const handleTabChange = (event: React.SyntheticEvent, newValue: number) => {
    setTabValue(newValue);
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

  const billingHistoryDataset = useMemo(() => {
    return reportsData?.getReports?.reports.map((report: GetReportResDto) => ({
      date: new Date(report.createdAt).toLocaleDateString(),
      amount: `${report.rawReport.totalAmountCurrency} ${(report.rawReport.totalAmountCents / 100).toFixed(2)}`,
      status: report.rawReport.paymentStatus || report.rawReport.status,
      period: report.period,
      id: report.id,
    }));
  }, [reportsData]);

  const handleMenuItemClick = (id: string, type: string) => {};
  const currentBill = React.useMemo(() => {
    if (!reportsData?.getReports || reportsData.getReports.reports.length === 0)
      return null;

    const unpaidBills = reportsData.getReports.reports
      .filter((b) => b.isPaid === false)
      .sort(
        (a, b) =>
          new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime(),
      );

    return unpaidBills[0] || null;
  }, [reportsData]);

  return (
    <LoadingWrapper
      width="100%"
      isLoading={paymentsLoading || reportsLoading}
      height="calc(100vh - 244px)"
    >
      <Box sx={{ width: '100%' }}>
        <Tabs
          value={tabValue}
          onChange={handleTabChange}
          aria-label="billing tabs"
          sx={{ borderBottom: 1, borderColor: 'divider' }}
        >
          <Tab label="Current Bill" />
          <Tab label="Billing History" />
        </Tabs>

        {tabValue === 0 && (
          <Box sx={{ py: 2 }}>
            <Grid container spacing={3}>
              <Grid item xs={12}>
                <CurrentBillCard
                  currentBill={currentBill}
                  onPay={handleAddPayment}
                  isLoading={reportsLoading || paymentsLoading}
                />
              </Grid>
              <Grid item xs={12}>
                <OutStandingBillCard
                  reports={reportsData?.getReports?.reports || []}
                  onPayAll={() => {}}
                  onPaySingle={(reportId) => {}}
                />
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
            </Grid>
          </Box>
        )}

        {tabValue === 1 && (
          <Box sx={{ py: 2 }}>
            <Paper
              sx={{
                height: '100%',
                borderRadius: '10px',
                px: { xs: 2, md: 3 },
                py: { xs: 2, md: 4 },
              }}
            >
              <DataTableWithOptions
                columns={BILLING_TABLE_COLUMNS}
                icon={SubscriberIcon}
                dataset={billingHistoryDataset}
                menuOptions={BILLING_HISTORY_TABLE_MENU}
                onMenuItemClick={handleMenuItemClick}
                emptyViewLabel="No billing history found"
                isRowClickable={false}
              />
            </Paper>
          </Box>
        )}
      </Box>

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
