/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/**
 * Stripe payment dialog for paying a billing invoice. Ported from the legacy
 * console (interfaces/console) into the ukama-console design system. The
 * `clientSecret` is the PaymentIntent secret returned by the payment service
 * in `payment.extra` after updatePayment(paymentMethod: 'stripe').
 */
import { useMemo, useState } from 'react';
import Button from '@mui/material/Button';
import CircularProgress from '@mui/material/CircularProgress';
import Dialog from '@mui/material/Dialog';
import DialogContent from '@mui/material/DialogContent';
import DialogTitle from '@mui/material/DialogTitle';
import IconButton from '@mui/material/IconButton';
import Stack from '@mui/material/Stack';
import Typography from '@mui/material/Typography';
import CloseRounded from '@mui/icons-material/CloseRounded';
import {
  Elements,
  PaymentElement,
  useElements,
  useStripe,
} from '@stripe/react-stripe-js';
import { loadStripe, type Stripe } from '@stripe/stripe-js';

import { useCurrency } from '@/lib/currency';
import { publicEnv } from '@/lib/runtime-env';

export interface StripePaymentDialogProps {
  open: boolean;
  onClose: () => void;
  /** Stripe PaymentIntent client secret (payment.extra). */
  clientSecret: string;
  /** Invoice total in minor units (cents) as a string, for the Pay label. */
  amountCents?: string | null;
  /** Short label for the dialog title, e.g. the invoice period. */
  periodLabel?: string;
  onPaymentSuccess?: () => void;
  onPaymentError?: (message: string) => void;
}

function PaymentForm({
  amountCents,
  onPaymentSuccess,
  onPaymentError,
  onClose,
}: {
  amountCents?: string | null;
  onPaymentSuccess?: () => void;
  onPaymentError?: (message: string) => void;
  onClose: () => void;
}) {
  const { money } = useCurrency();
  const stripe = useStripe();
  const elements = useElements();
  const [isProcessing, setIsProcessing] = useState(false);

  const amount = Number(amountCents);
  const payLabel = Number.isFinite(amount) ? `Pay ${money(amount / 100)}` : 'Pay';

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();
    if (!stripe || !elements) return;

    setIsProcessing(true);
    try {
      const { error, paymentIntent } = await stripe.confirmPayment({
        elements,
        redirect: 'if_required',
      });
      if (error) {
        onPaymentError?.(error.message ?? 'Payment failed');
      } else if (paymentIntent?.status === 'succeeded') {
        onPaymentSuccess?.();
      }
    } catch (err) {
      onPaymentError?.(err instanceof Error ? err.message : 'Payment failed');
    } finally {
      setIsProcessing(false);
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <PaymentElement options={{ layout: 'tabs' }} />
      <Stack direction="row" justifyContent="flex-end" spacing={1.5} sx={{ pt: 2 }}>
        <Button onClick={onClose} color="inherit">
          Cancel
        </Button>
        <Button type="submit" variant="contained" disabled={!stripe || isProcessing}>
          {isProcessing ? <CircularProgress size={22} color="inherit" /> : payLabel}
        </Button>
      </Stack>
    </form>
  );
}

export default function StripePaymentDialog({
  open,
  onClose,
  clientSecret,
  amountCents,
  periodLabel,
  onPaymentSuccess,
  onPaymentError,
}: StripePaymentDialogProps) {
  // Stripe instance is keyed off the publishable key; memoize so it isn't
  // re-created on every render.
  const stripePromise = useMemo<Promise<Stripe | null>>(() => {
    const pk = publicEnv().stripePk;
    return pk ? loadStripe(pk) : Promise.resolve(null);
  }, []);

  return (
    <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
      <DialogTitle>
        Pay invoice{periodLabel ? ` · ${periodLabel}` : ''}
        <IconButton
          aria-label="close"
          onClick={onClose}
          sx={{ position: 'absolute', right: 8, top: 8, color: 'var(--uk-ink-3)' }}
        >
          <CloseRounded />
        </IconButton>
      </DialogTitle>
      <DialogContent>
        <Typography variant="body2" sx={{ color: 'var(--uk-ink-2)', mb: 2 }}>
          Enter your payment details to settle this invoice.
        </Typography>
        {clientSecret ? (
          <Elements
            stripe={stripePromise}
            options={{ clientSecret, appearance: { theme: 'stripe' } }}
          >
            <PaymentForm
              amountCents={amountCents}
              onPaymentSuccess={onPaymentSuccess}
              onPaymentError={onPaymentError}
              onClose={onClose}
            />
          </Elements>
        ) : null}
      </DialogContent>
    </Dialog>
  );
}
