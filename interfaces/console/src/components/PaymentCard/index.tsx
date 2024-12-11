'use client';
import React, { useState, useEffect } from 'react';
import {
  Box,
  Typography,
  Paper,
  Select,
  MenuItem,
  Stack,
  Divider,
  Collapse,
  Grid,
  IconButton,
  SelectChangeEvent,
  Skeleton,
} from '@mui/material';
import { CreditCard, PaymentOutlined } from '@mui/icons-material';
import colors from '@/theme/colors';

interface PaymentCardProps {
  amount: string;
  startDate: string;
  endDate: string;
  onChangePaymentMethod: (method: string) => void;
  paymentMethods: string[];
  onPaymentMethodSelect: (method: string) => void;
  children?: React.ReactNode;
  isLoading?: boolean;
}

const PaymentCard: React.FC<PaymentCardProps> = ({
  amount,
  startDate,
  endDate,
  onChangePaymentMethod,
  paymentMethods,
  onPaymentMethodSelect,
  children,
  isLoading = false,
}) => {
  const [isPaymentFormOpen, setIsPaymentFormOpen] = useState(false);
  const [paymentMethod, setPaymentMethod] = useState('');

  useEffect(() => {
    if (paymentMethods.includes('Stripe')) {
      setPaymentMethod('Stripe');
      onChangePaymentMethod('Stripe');
      onPaymentMethodSelect('Stripe');
      setIsPaymentFormOpen(true);
    }
  }, [paymentMethods]);

  const handlePaymentMethodChange = (event: SelectChangeEvent) => {
    const method = event.target.value;
    setPaymentMethod(method);
    onChangePaymentMethod(method);
    onPaymentMethodSelect(method);
    setIsPaymentFormOpen(method === 'Stripe');
  };

  const StripeFormSkeleton = () => (
    <Box>
      <Skeleton variant="rectangular" width="100%" height={50} sx={{ mb: 2 }} />
      <Skeleton variant="rectangular" width="100%" height={50} sx={{ mb: 2 }} />
    </Box>
  );

  return (
    <Grid container spacing={3}>
      <Grid item xs={12} md={6}>
        <Paper
          elevation={3}
          sx={{
            padding: 3,
            borderRadius: 2,
            backgroundColor: colors.white,
            boxShadow: '0 4px 6px rgba(0,0,0,0.1)',
          }}
        >
          <Stack spacing={2}>
            <Box
              display="flex"
              justifyContent="space-between"
              alignItems="center"
            >
              <Typography variant="h6" color="text.primary">
                Next Payment
              </Typography>
              <Typography variant="body2" color="text.secondary">
                {startDate} - {endDate}
              </Typography>
            </Box>

            <Divider />

            <Box
              display="flex"
              alignItems="center"
              justifyContent="space-between"
            >
              <Typography variant="h4" fontWeight="bold" color="primary">
                ${amount}
              </Typography>
              <IconButton color="primary">
                <PaymentOutlined />
              </IconButton>
            </Box>

            <Typography
              variant="body2"
              color="text.secondary"
              sx={{ fontStyle: 'italic' }}
            >
              Detailed breakdown available below
            </Typography>
          </Stack>
        </Paper>
      </Grid>

      <Grid item xs={12} md={6}>
        <Paper
          elevation={3}
          sx={{
            padding: 3,
            borderRadius: 2,
            backgroundColor: colors.white,
            boxShadow: '0 4px 6px rgba(0,0,0,0.1)',
          }}
        >
          <Stack spacing={2}>
            <Typography
              variant="h6"
              color="text.primary"
              sx={{ display: 'flex', alignItems: 'center', gap: 1 }}
            >
              <CreditCard /> Payment Method
            </Typography>
            {isLoading ? (
              <Skeleton
                variant="rectangular"
                width="100%"
                height={50}
                sx={{ mb: 2 }}
              />
            ) : (
              <Select
                value={paymentMethod}
                onChange={handlePaymentMethodChange}
                fullWidth
                variant="outlined"
                displayEmpty
                sx={{
                  '& .MuiSelect-select': {
                    display: 'flex',
                    alignItems: 'center',
                    gap: 1,
                  },
                }}
              >
                {paymentMethods.map((method, index) => (
                  <MenuItem key={index} value={method}>
                    {method === 'Stripe' && <CreditCard sx={{ mr: 1 }} />}
                    {method}
                  </MenuItem>
                ))}
              </Select>
            )}

            <Collapse in={isPaymentFormOpen}>
              <Box mt={2}>{isLoading ? <StripeFormSkeleton /> : children}</Box>
            </Collapse>
          </Stack>
        </Paper>
      </Grid>
    </Grid>
  );
};

export default PaymentCard;
