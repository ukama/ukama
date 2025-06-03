/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { ReportDto } from '@/client/graphql/generated';
import colors from '@/theme/colors';
import {
  Box,
  Button,
  Divider,
  Paper,
  Stack,
  Typography,
  useMediaQuery,
  useTheme,
} from '@mui/material';
import { format } from 'date-fns';
import React from 'react';
interface OutStandingBillCardProps {
  reports?: ReportDto[];
  loading?: boolean;
  onPaySingle?: (reportId: string) => void;
}
const OutStandingBillCard: React.FC<OutStandingBillCardProps> = ({
  reports,
  onPaySingle,
}) => {
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('sm'));

  const outstandingReports = React.useMemo(() => {
    if (!reports) return [];

    const unpaidReports = reports.filter((report) => !report.isPaid);

    return unpaidReports.length >= 2 ? unpaidReports.slice(1) : [];
  }, [reports]);

  return (
    <Box>
      <Paper
        elevation={2}
        sx={{
          p: isMobile ? 2 : 4,
          borderRadius: '10px',
        }}
      >
        <Typography
          variant="h6"
          sx={{
            fontSize: isMobile ? '1.1rem' : '1.25rem',
            mb: isMobile ? 1 : 0,
          }}
        >
          Outstanding Bills
        </Typography>

        <Typography
          variant="body2"
          sx={{
            color: colors.vulcan,
            mb: 2,
            fontSize: isMobile ? '0.8rem' : '0.875rem',
          }}
        >
          {`${outstandingReports.length} overdue bill${outstandingReports.length > 1 ? 's' : ''} for Ukama Console plan.`}
        </Typography>
        {outstandingReports.length > 0 && <Divider sx={{ mb: 2 }} />}

        {outstandingReports.map((report, index) => (
          <React.Fragment key={report.id}>
            <Stack
              direction={isMobile ? 'column' : 'row'}
              spacing={2}
              justifyContent={'space-between'}
              alignItems={isMobile ? 'stretch' : 'center'}
              sx={{ mb: index < outstandingReports.length - 1 ? 2 : 0 }}
            >
              <Typography
                variant="body1"
                sx={{
                  color: colors.vulcan,
                  fontSize: isMobile ? '0.8rem' : '0.875rem',
                  textAlign: isMobile ? 'left' : 'inherit',
                }}
              >
                {report.createdAt && (
                  <>
                    {`${format(new Date(report.createdAt), 'MMMM')} `}
                    <span>overdue bill</span>
                    {`: $${Number(Number(report.rawReport.totalAmountCents) / 100).toFixed(2)} `}
                    <span style={{ color: colors.error }}>
                      {`Overdue on ${format(new Date(report.rawReport.paymentDueDate), 'dd/MM/yyyy')}`}
                    </span>
                  </>
                )}
              </Typography>

              {onPaySingle && (
                <Button
                  variant="contained"
                  onClick={() => onPaySingle(report.id)}
                  fullWidth={isMobile}
                  sx={{
                    py: isMobile ? 1 : 'auto',
                    px: isMobile ? 2 : 'auto',
                    mt: isMobile ? 1 : 0,
                  }}
                >
                  Pay now
                </Button>
              )}
            </Stack>
          </React.Fragment>
        ))}
      </Paper>
    </Box>
  );
};

export default OutStandingBillCard;
