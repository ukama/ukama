import React, { useState } from 'react';
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
import ChevronLeftIcon from '@mui/icons-material/ChevronLeft';
import ChevronRightIcon from '@mui/icons-material/ChevronRight';

interface Bill {
  id: string;
  amount: number;
  dueDate: string;
  plan: string;
}

interface OutStandingBillCardProps {
  bills: Bill[];
  loading?: boolean;
}

const OutStandingBillCard: React.FC<OutStandingBillCardProps> = ({
  bills = [],
  loading = false,
}) => {
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('sm'));
  const [currentBillIndex, setCurrentBillIndex] = useState(0);

  if (bills.length === 0) {
    return (
      <Box>
        <Paper
          elevation={2}
          sx={{
            p: isMobile ? 2 : 4,
            borderRadius: '10px',
          }}
        >
          <Typography variant="body1">No outstanding bills</Typography>
        </Paper>
      </Box>
    );
  }

  const currentBill = bills[currentBillIndex];

  const handleNextBill = () => {
    setCurrentBillIndex((prev) => (prev + 1) % bills.length);
  };

  const handlePrevBill = () => {
    setCurrentBillIndex((prev) => (prev - 1 + bills.length) % bills.length);
  };

  const totalDueAmount = bills.reduce((sum, bill) => sum + bill.amount, 0);

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
            Outstanding bills
          </Typography>
          <Button
            variant="contained"
            fullWidth={isMobile}
            sx={{
              py: isMobile ? 1 : 'auto',
              px: isMobile ? 2 : 'auto',
            }}
          >
            Pay all outstanding bills
          </Button>
        </Stack>
        <Typography
          variant="body2"
          sx={{
            color: colors.vulcan,
            mb: 2,
            fontSize: isMobile ? '0.8rem' : '0.875rem',
          }}
        >
          Overdue bills for {currentBill.plan} plan
        </Typography>
        <Divider sx={{ mb: 2 }} />

        {bills.length > 1 && (
          <Stack
            direction="row"
            justifyContent="space-between"
            alignItems="center"
            sx={{ mb: 2 }}
          >
            <Button
              onClick={handlePrevBill}
              disabled={bills.length <= 1}
              startIcon={<ChevronLeftIcon />}
            >
              Previous
            </Button>
            <Typography variant="body2">
              Bill {currentBillIndex + 1} of {bills.length}
            </Typography>
            <Button
              onClick={handleNextBill}
              disabled={bills.length <= 1}
              endIcon={<ChevronRightIcon />}
            >
              Next
            </Button>
          </Stack>
        )}

        <Stack
          direction={isMobile ? 'column' : 'row'}
          spacing={2}
          justifyContent={'space-between'}
          alignItems={isMobile ? 'stretch' : 'center'}
        >
          <Stack direction="column" spacing={1} alignItems="center">
            <Stack direction="row" spacing={1} alignItems="center">
              <Typography
                variant="body2"
                sx={{
                  color: colors.vulcan,
                  fontSize: isMobile ? '0.8rem' : '0.875rem',
                  textAlign: isMobile ? 'left' : 'inherit',
                }}
              >
                Total due:
              </Typography>
              <Typography
                variant="body2"
                sx={{
                  color: colors.vulcan,
                  fontSize: isMobile ? '0.8rem' : '0.875rem',
                  textAlign: isMobile ? 'left' : 'inherit',
                }}
              >
                ${totalDueAmount.toFixed(2)}
              </Typography>
            </Stack>

            <Stack direction="row" spacing={1} alignItems="center">
              <Typography
                variant="body2"
                sx={{
                  color: colors.vulcan,
                  fontSize: isMobile ? '0.8rem' : '0.875rem',
                  textAlign: isMobile ? 'left' : 'inherit',
                }}
              >
                Overdue on
              </Typography>
              <Typography
                variant="body2"
                sx={{
                  color: colors.red,
                  fontSize: isMobile ? '0.8rem' : '0.875rem',
                  textAlign: isMobile ? 'left' : 'inherit',
                }}
              >
                {currentBill.dueDate}
              </Typography>
            </Stack>
          </Stack>

          <Button
            variant="contained"
            fullWidth={isMobile}
            sx={{
              py: isMobile ? 1 : 'auto',
              px: isMobile ? 2 : 'auto',
              mt: isMobile ? 1 : 0,
            }}
          >
            Pay now
          </Button>
        </Stack>
      </Paper>
    </Box>
  );
};

export default OutStandingBillCard;
