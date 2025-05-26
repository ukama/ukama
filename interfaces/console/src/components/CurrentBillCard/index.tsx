/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import React from 'react';
import {
  Typography,
  Paper,
  Stack,
  Divider,
  Grid,
  Button,
  Alert,
} from '@mui/material';
import colors from '@/theme/colors';
import { ReportDto } from '@/client/graphql/generated';
import CustomLoadingSkeleton from '../CustomLoadingSkeleton';
import { format, isValid } from 'date-fns';

interface CurrentBillCardProps {
  isLoading?: boolean;
  onPay: (billId: string) => void;
  bills?: ReportDto[];
}

const CurrentBillCard: React.FC<CurrentBillCardProps> = ({
  bills,
  isLoading = false,
  onPay,
}) => {
  const formatDate = (dateString: string | null | undefined) => {
    if (!dateString) return '-';
    const date = new Date(dateString);
    return isValid(date) ? format(date, 'MMM dd, yyyy') : '-';
  };

  const renderContent = <T,>(
    content: T,
    renderFn: (value: T) => React.ReactNode,
    skeletonProps?: { width?: number; height?: number },
  ) => {
    if (isLoading || content === undefined || content === null) {
      return <CustomLoadingSkeleton {...skeletonProps} />;
    }
    return renderFn(content);
  };

  const currentBill =
    bills && bills?.length > 0
      ? [...bills].sort(
        (a, b) =>
          new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime(),
      )[0]
      : null;

  const hasMultipleUnpaidBills =
    bills &&
    bills?.filter(
      (bill) =>
        bill.isPaid === false &&
        Number(bill.rawReport?.totalAmountCents || 0) > 100,
    ).length > 1;

  const handlePay = () => {
    if (onPay && currentBill?.id) {
      onPay(currentBill.id);
    }
  };

  const totalAmountCents = Number(
    currentBill?.rawReport?.totalAmountCents ?? 0,
  );

  const getBillingMessage = (dueDate: string | null | undefined) => {
    if (!bills || bills.length === 0) {
      return 'Monthly SaaS fee for Ukama Console. Bill available at the end of month, and will be due within 30 days.';
    }

    if (currentBill && totalAmountCents < 1) {
      return 'Monthly SaaS fee for Ukama Console. Free trial active.';
    }

    if (dueDate) {
      return `Monthly SaaS fee for Ukama Console. Bill available on ${formatDate(dueDate)} and will be due on ${formatDate(dueDate)}`;
    }

    return 'Monthly SaaS fee for Ukama Console. Bill available at the end of month, and will be due within 30 days.';
  };

  return (
    <Grid container spacing={2}>
      {hasMultipleUnpaidBills && (
        <Grid item xs={12}>
          <Alert severity="error">
            <Typography variant="body2" sx={{ color: colors.vulcan }}>
              Service will be paused unless you pay your outstanding bills.
            </Typography>
          </Alert>
        </Grid>
      )}

      <Grid item xs={12}>
        <Paper
          elevation={3}
          sx={{
            padding: 2,
            borderRadius: 2,
            backgroundColor: colors.white,
            boxShadow: '0 4px 6px rgba(0,0,0,0.1)',
          }}
        >
          <Stack spacing={2}>
            <Stack direction="row" spacing={2} justifyContent={'space-between'}>
              <Typography variant="h6" sx={{ color: colors.vulcan }}>
                Ukama plan
              </Typography>

              <Stack direction={'row'} spacing={1} alignItems={'center'}>
                {renderContent(
                  currentBill?.rawReport?.subscriptions[0]?.startedAt,
                  (date) => (
                    <Typography variant="body2" sx={{ color: colors.vulcan }}>
                      {formatDate(date)}
                    </Typography>
                  ),
                  { width: 100, height: 20 },
                )}

                <Typography variant="body2" sx={{ color: colors.vulcan }}>
                  -
                </Typography>

                {renderContent(
                  currentBill?.rawReport?.subscriptions[0]?.terminatedAt,
                  (date) => (
                    <Typography variant="body2" sx={{ color: colors.vulcan }}>
                      {formatDate(date)}
                    </Typography>
                  ),
                  { width: 100, height: 20 },
                )}
              </Stack>
            </Stack>

            <Typography variant="body2" sx={{ color: colors.vulcan }}>
              {getBillingMessage(currentBill?.rawReport?.paymentDueDate)}
            </Typography>

            <Divider />

            <Stack direction="row" spacing={2} justifyContent={'space-between'}>
              {renderContent(
                currentBill,
                (bill) => (
                  <Typography variant="body1">
                    {Number(bill?.rawReport?.totalAmountCents) < 1
                      ? 'Free trial $0'
                      : `Total due: $${(Number(bill?.rawReport?.totalAmountCents || 0) / 100).toFixed(2)}`}
                  </Typography>
                ),
                { width: 100, height: 20 },
              )}

              {currentBill && totalAmountCents >= 1 && (
                <Button
                  variant="contained"
                  onClick={handlePay}
                  disabled={!onPay || isLoading || !currentBill?.id}
                >
                  Pay now
                </Button>
              )}
            </Stack>
          </Stack>
        </Paper>
      </Grid>
    </Grid>
  );
};

export default CurrentBillCard;
