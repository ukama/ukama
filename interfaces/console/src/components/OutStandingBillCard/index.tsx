import React from 'react';
import {
  Box,
  Paper,
  Stack,
  Typography,
  Button,
  Divider,
  useMediaQuery,
  useTheme,
} from '@mui/material';
import colors from '@/theme/colors';
import { GetReportResDto } from '@/client/graphql/generated';
import { format } from 'date-fns';

interface OutStandingBillCardProps {
  reports?: GetReportResDto[];
  loading?: boolean;
  onPayAll?: () => void;
  onPaySingle?: (reportId: string) => void;
}

const OutStandingBillCard: React.FC<OutStandingBillCardProps> = ({
  reports = [],
  loading = false,
  onPayAll,
  onPaySingle,
}) => {
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('sm'));
  const outstandingReports = reports.filter((report) => !report.isPaid);

  const totalOutstandingAmount = outstandingReports.reduce(
    (total, report) => total + report.rawReport.totalAmountCents / 100,
    0,
  );

  if (outstandingReports.length === 0) {
    return null;
  }

  return (
    <Box>
      <Paper
        elevation={2}
        sx={{
          p: isMobile ? 2 : 4,
          borderRadius: '10px',
        }}
      >
        <Stack
          direction={isMobile ? 'column' : 'row'}
          justifyContent="space-between"
          alignItems={isMobile ? 'flex-start' : 'center'}
          spacing={isMobile ? 2 : 0}
          sx={{ mb: 2 }}
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
          {onPayAll && (
            <Button
              variant="contained"
              onClick={onPayAll}
              fullWidth={isMobile}
              sx={{
                py: isMobile ? 1 : 'auto',
                px: isMobile ? 2 : 'auto',
              }}
            >
              Pay all outstanding bills
            </Button>
          )}
        </Stack>
        <Typography
          variant="body2"
          sx={{
            color: colors.vulcan,
            mb: 2,
            fontSize: isMobile ? '0.8rem' : '0.875rem',
          }}
        >
          {`${outstandingReports.length} overdue bill${outstandingReports.length > 1 ? 's' : ''} for Ukama Console plan`}
        </Typography>
        <Divider sx={{ mb: 2 }} />

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
                variant="body2"
                sx={{
                  color: colors.vulcan,
                  fontSize: isMobile ? '0.8rem' : '0.875rem',
                  textAlign: isMobile ? 'left' : 'inherit',
                }}
              >
                {`Due: ${report.rawReport.totalAmountCurrency} ${(report.rawReport.totalAmountCents / 100).toFixed(2)} `}
                {report.period && `for ${report.period}`}
                {report.createdAt &&
                  ` (Created on ${format(new Date(report.createdAt), 'MM/dd/yy')})`}
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
            {index < outstandingReports.length - 1 && (
              <Divider sx={{ my: 2 }} />
            )}
          </React.Fragment>
        ))}

        <Divider sx={{ mt: 2, mb: 1 }} />
        <Stack
          direction={isMobile ? 'column' : 'row'}
          justifyContent="space-between"
          alignItems={isMobile ? 'flex-start' : 'center'}
        >
          <Typography
            variant="body1"
            fontWeight="bold"
            sx={{
              fontSize: isMobile ? '1rem' : '1.125rem',
            }}
          >
            Total Outstanding:{' '}
            {outstandingReports[0]?.rawReport?.totalAmountCurrency || '$'}
            {totalOutstandingAmount.toFixed(2)}
          </Typography>
        </Stack>
      </Paper>
    </Box>
  );
};

export default OutStandingBillCard;
