/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import React from 'react';
import { Box, Skeleton, Grid } from '@mui/material';
import CurrentBillCard from '@/components/CurrentBillCard';
import PaymentMethodCard from '@/components/PaymentMethodCard';
import DataUsageComponent from '@/components/DataUsageComponent';
import NotificationEmailSettings from '@/components/NotificationEmailSettings';

interface CurrentBillProps {
  dataUsagePaid: number;
  nextPaymentAmount: number;
  nextPaymentDate: string;
  notificationEmail: string;
  loading?: boolean;
  onAddPaymentMethod: () => void;
  handleAddPayment: () => void;
}

const CurrentBill: React.FC<CurrentBillProps> = ({
  dataUsagePaid,
  nextPaymentAmount,
  nextPaymentDate,
  notificationEmail,
  loading,
  handleAddPayment,
  onAddPaymentMethod,
}) => {
  return (
    <Box sx={{ py: 2 }}>
      {(loading || !nextPaymentAmount) && (
        <Grid container spacing={2}>
          <Grid item xs={12}>
            <Skeleton
              variant="rectangular"
              width="100%"
              height={200}
              sx={{ borderRadius: 2, mb: 2 }}
            />
          </Grid>
          <Grid item xs={12}>
            <Skeleton
              variant="rectangular"
              width="100%"
              height={150}
              sx={{ borderRadius: 2, mb: 2 }}
            />
          </Grid>
          <Grid item xs={12} md={6}>
            <Skeleton
              variant="rectangular"
              width="100%"
              height={200}
              sx={{ borderRadius: 2 }}
            />
          </Grid>
          <Grid item xs={12} md={6}>
            <Skeleton
              variant="rectangular"
              width="100%"
              height={200}
              sx={{ borderRadius: 2 }}
            />
          </Grid>
        </Grid>
      )}

      {!loading && nextPaymentAmount > 0 && (
        <Grid container spacing={2}>
          <Grid item xs={12}>
            <CurrentBillCard
              amount={nextPaymentAmount.toString()}
              startDate={nextPaymentDate}
              endDate={nextPaymentDate}
              onPay={handleAddPayment}
            />
          </Grid>
          <Grid item xs={12}>
            <PaymentMethodCard onAddPaymentMethod={onAddPaymentMethod} />
          </Grid>
          <Grid item xs={12} md={6}>
            <DataUsageComponent
              dataUsagePaid={dataUsagePaid}
              subscriberCount={0}
            />
          </Grid>
          <Grid item xs={12} md={6}>
            <NotificationEmailSettings
              primaryEmail={notificationEmail}
              additionalEmails={[]}
            />
          </Grid>
        </Grid>
      )}
    </Box>
  );
};

export default CurrentBill;
