/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import React from 'react';
import { Typography, Paper, Stack, Divider, Grid, Button } from '@mui/material';
import colors from '@/theme/colors';
import { GetReportResDto } from '@/client/graphql/generated';
import CustomLoadingSkeleton from '../CustomLoadingSkeleton';
import { format } from 'date-fns';

interface CurrentBillCardProps {
  currentBill?: GetReportResDto | null;
  isLoading?: boolean;
  onPay?: (billId?: string) => void;
}

const CurrentBillCard: React.FC<CurrentBillCardProps> = ({
  currentBill,
  isLoading = false,
  onPay,
}) => {
  const renderContent = <T,>(
    content: T,
    renderFn: (value: T) => React.ReactNode,
    skeletonProps?: { width?: number; height?: number },
  ) => {
    if (isLoading || !content) {
      return <CustomLoadingSkeleton {...skeletonProps} />;
    }
    return renderFn(content);
  };

  const handlePay = () => {
    if (onPay && currentBill?.id) {
      onPay(currentBill.id);
    }
  };

  return (
    <Grid container spacing={3}>
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
              {renderContent(
                currentBill?.createdAt,
                (date) => (
                  <Typography variant="body2" sx={{ color: colors.vulcan }}>
                    {format(new Date(date || ''), 'MMM dd, yyyy')}
                  </Typography>
                ),
                { width: 100, height: 20 },
              )}
            </Stack>

            <Stack direction="row" spacing={2} alignItems={'center'}>
              <Typography variant="body2" sx={{ color: colors.vulcan }}>
                Monthly SaaS fee for Ukama Console. Bill available
              </Typography>
              {renderContent(
                currentBill?.createdAt,
                (date) => (
                  <Typography variant="body2" sx={{ color: colors.vulcan }}>
                    {format(new Date(date || ''), 'MMM dd, yyyy')}
                  </Typography>
                ),
                { width: 100, height: 20 },
              )}
            </Stack>

            <Divider />
            <Stack direction="row" spacing={2} justifyContent={'space-between'}>
              {renderContent(
                currentBill?.rawReport?.amountCents ?? 0,
                (amountCents) => (
                  <Typography variant="body1">
                    Total due: $ {(amountCents / 100).toFixed(2)}
                  </Typography>
                ),
                { width: 100, height: 20 },
              )}
              <Button
                variant="contained"
                size="large"
                onClick={handlePay}
                disabled={!onPay || isLoading || !currentBill?.id}
              >
                Pay now
              </Button>
            </Stack>
          </Stack>
        </Paper>
      </Grid>
    </Grid>
  );
};

export default CurrentBillCard;
