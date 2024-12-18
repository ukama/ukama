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
} from '@mui/material';
import { loadStripe } from '@stripe/stripe-js';
import {
  Elements,
  PaymentElement,
  useStripe,
  useElements,
} from '@stripe/react-stripe-js';
import { useAppContext } from '@/context';
import { GetReportResDto } from '@/client/graphql/generated';
import colors from '@/theme/colors';

const StripePaymentDialog: React.FC<{
  open: boolean;
  onClose: () => void;
  extraKey: string;
  amount: number;
  onPaymentSuccess?: () => void;
  onPaymentError?: (error: any) => void;
  bill?: GetReportResDto;
}> = ({
  open,
  onClose,
  extraKey,
  amount,
  onPaymentSuccess,
  onPaymentError,
  bill,
}) => {
  const { env } = useAppContext();
  const stripePromise = loadStripe(env.STRIPE_PK!);
  const [isProcessing, setIsProcessing] = useState(false);

  return (
    <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
      <DialogTitle>Pay bill {bill?.createdAt}</DialogTitle>

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
          <form
            onSubmit={async (event) => {
              event.preventDefault();
              const stripe = useStripe();
              const elements = useElements();

              if (!stripe || !elements) {
                return;
              }

              setIsProcessing(true);

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
            }}
          >
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
                  `Pay $${amount.toFixed(2)}`
                )}
              </Button>
              <Button onClick={onClose} color="secondary">
                Cancel
              </Button>
            </Stack>
          </form>
        </Elements>
      </DialogContent>
    </Dialog>
  );
};

export default StripePaymentDialog;
