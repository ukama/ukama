'use client';
import React, { use } from 'react';
import { Box, Typography, Paper, Select, MenuItem } from '@mui/material';

interface PaymentCardProps {
  amount: string;
  startDate: string;
  endDate: string;
  paymentMethod: string;
  onChangePaymentMethod: (method: string) => void;
  paymentMethods: string[];
}

const PaymentCard: React.FC<PaymentCardProps> = ({
  amount,
  startDate,
  endDate,
  paymentMethod,
  onChangePaymentMethod,
  paymentMethods,
}) => {
  return (
    <Box display="flex" gap={2}>
      <Paper elevation={2} sx={{ padding: 2, flex: 1 }}>
        <Typography variant="subtitle1" fontWeight="bold">
          Next payment
        </Typography>
        <Typography
          variant="body2"
          color="text.secondary"
          sx={{ marginBottom: 1 }}
        >
          Detailed breakdown available below.
        </Typography>
        <Typography variant="h4" fontWeight="bold">
          {amount}
        </Typography>
        <Typography variant="body2" color="text.secondary">
          {startDate} - {endDate}
        </Typography>
      </Paper>

      <Paper elevation={2} sx={{ padding: 2, flex: 1 }}>
        <Typography variant="subtitle1" fontWeight="bold">
          Payment information
        </Typography>
        <Typography variant="caption" color="text.secondary">
          PAYMENT METHOD
        </Typography>
        <Select
          value={paymentMethod}
          onChange={(e) => onChangePaymentMethod(e.target.value)}
          fullWidth
          displayEmpty
          sx={{ marginTop: 1, marginBottom: 1 }}
        >
          {paymentMethods.map((method, index) => (
            <MenuItem key={index} value={method}>
              {method}
            </MenuItem>
          ))}
        </Select>
        <Typography variant="caption" color="text.secondary">
          *Automatically charged EOD on the last day of the billing cycle
        </Typography>
      </Paper>
    </Box>
  );
};

export default PaymentCard;
