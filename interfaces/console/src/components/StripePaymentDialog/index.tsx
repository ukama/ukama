import React, { useState } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  CircularProgress,
  Box,
  Typography,
} from '@mui/material';
import { loadStripe } from '@stripe/stripe-js';
import {
  Elements,
  PaymentElement,
  useStripe,
  useElements,
} from '@stripe/react-stripe-js';
import { useAppContext } from '@/context';

const PaymentForm: React.FC<{
  clientSecret: string;
  amount: number;
  onPaymentSuccess?: () => void;
  onPaymentError?: (error: any) => void;
}> = ({ clientSecret, amount, onPaymentSuccess, onPaymentError }) => {
  const [isProcessing, setIsProcessing] = useState(false);
  const stripe = useStripe();
  const elements = useElements();

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();

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
  };

  return (
    <form onSubmit={handleSubmit}>
      <PaymentElement
        options={{
          layout: 'tabs',
          defaultValues: {},
        }}
      />

      <Button
        type="submit"
        variant="contained"
        color="primary"
        disabled={isProcessing || !stripe}
        fullWidth
        sx={{ mt: 2, py: 1.5 }}
      >
        {isProcessing ? (
          <CircularProgress size={24} />
        ) : (
          `Pay $${amount.toFixed(2)}`
        )}
      </Button>
    </form>
  );
};

const StripePaymentDialog: React.FC<{
  open: boolean;
  onClose: () => void;
  clientSecret: string;
  amount: number;
  onPaymentSuccess?: () => void;
  onPaymentError?: (error: any) => void;
}> = ({
  open,
  onClose,
  clientSecret,
  amount,
  onPaymentSuccess,
  onPaymentError,
}) => {
  const { env } = useAppContext();
  const stripePromise = loadStripe(env.STRIPE_PK!);
  console.log('BRACKLEY :', process.env.STRIPE_PK);

  return (
    <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
      <DialogTitle>Complete Payment</DialogTitle>
      <DialogContent>
        <Elements
          stripe={stripePromise}
          options={{
            clientSecret,
            appearance: {
              theme: 'stripe',
            },
          }}
        >
          <PaymentForm
            clientSecret={clientSecret}
            amount={amount}
            onPaymentSuccess={onPaymentSuccess}
            onPaymentError={onPaymentError}
          />
        </Elements>
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose} color="secondary">
          Cancel
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default StripePaymentDialog;
