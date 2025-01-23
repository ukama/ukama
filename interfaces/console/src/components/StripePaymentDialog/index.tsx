/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import React, { useState } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  Button,
  CircularProgress,
  Typography,
  Stack,
  IconButton,
} from '@mui/material';
import { loadStripe } from '@stripe/stripe-js';
import {
  Elements,
  PaymentElement,
  useStripe,
  useElements,
} from '@stripe/react-stripe-js';
import { useAppContext } from '@/context';
import { ReportDto } from '@/client/graphql/generated';
import colors from '@/theme/colors';
import { format } from 'date-fns';
import CloseIcon from '@mui/icons-material/Close';

const PaymentForm: React.FC<{
  extraKey: string;
  bill?: ReportDto;
  onPaymentSuccess?: () => void;
  onPaymentError?: (error: any) => void;
  onClose: () => void;
}> = ({ bill, onPaymentSuccess, onPaymentError, onClose }) => {
  const [isProcessing, setIsProcessing] = useState(false);
  const stripe = useStripe();
  const elements = useElements();

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();

    if (!stripe || !elements) {
      return;
    }

    setIsProcessing(true);

    try {
      const { error, paymentIntent } = await stripe.confirmPayment({
        elements,
        redirect: 'if_required',
      });

      if (error) {
        setIsProcessing(false);
        onPaymentError?.(error);
      } else if (paymentIntent?.status === 'succeeded') {
        setIsProcessing(false);
        onPaymentSuccess?.();
      }
    } catch (error) {
      setIsProcessing(false);
      onPaymentError?.(error);
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <PaymentElement
        options={{
          layout: 'tabs',
          defaultValues: {},
        }}
      />
      <Stack
        direction={'row'}
        justifyContent={'flex-end'}
        spacing={2}
        sx={{ py: 2 }}
      >
        <Button type="submit" variant="contained" disabled={isProcessing}>
          {isProcessing ? (
            <CircularProgress size={24} />
          ) : (
            `Pay $${(
              Number(bill?.rawReport?.totalAmountCents ?? 0) / 100
            ).toFixed(2)}`
          )}
        </Button>
        <Button onClick={onClose} color="secondary">
          Cancel
        </Button>
      </Stack>
    </form>
  );
};

const StripePaymentDialog: React.FC<{
  open: boolean;
  onClose: () => void;
  extraKey: string;
  onPaymentSuccess?: () => void;
  onPaymentError?: (error: any) => void;
  bill?: ReportDto;
}> = ({ open, onClose, extraKey, onPaymentSuccess, onPaymentError, bill }) => {
  const { env } = useAppContext();
  const stripePromise = loadStripe(env.STRIPE_PK!);

  return (
    <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
      <DialogTitle>
        Pay bill({format(new Date(bill?.createdAt || ''), 'MMM')})
        <IconButton
          aria-label="close"
          onClick={onClose}
          sx={{
            position: 'absolute',
            right: 8,
            top: 8,
            color: (theme) => theme.palette.grey[500],
          }}
        >
          <CloseIcon />
        </IconButton>
      </DialogTitle>

      <DialogContent>
        <Typography variant="body1" sx={{ color: colors.vulcan, mb: 2 }}>
          Please enter your payment information to pay for your current bill.
        </Typography>
        <Elements
          stripe={stripePromise}
          options={{
            clientSecret: extraKey,
            appearance: {
              theme: 'stripe',
            },
          }}
        >
          <PaymentForm
            extraKey={extraKey}
            bill={bill}
            onPaymentSuccess={onPaymentSuccess}
            onPaymentError={onPaymentError}
            onClose={onClose}
          />
        </Elements>
      </DialogContent>
    </Dialog>
  );
};

export default StripePaymentDialog;
