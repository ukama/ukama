import React, { useState } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  CircularProgress,
} from '@mui/material';
import { loadStripe } from '@stripe/stripe-js';
import {
  Elements,
  CardElement,
  useStripe,
  useElements,
} from '@stripe/react-stripe-js';

interface StripePaymentDialogProps {
  open: boolean;
  onClose: () => void;
  clientSecret: string;
  amount: number;
  currency?: string;
  loading?: boolean;
  onPaymentSuccess?: () => void;
  onPaymentError?: (error: any) => void;
}

const stripePromise = loadStripe(
  process.env.NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY!,
);

const PaymentForm: React.FC<{
  clientSecret: string;
  amount: number;
  currency: string;
  onPaymentSuccess?: () => void;
  onPaymentError?: (error: any) => void;
}> = ({ clientSecret, amount, onPaymentSuccess, onPaymentError }) => {
  const [isProcessing, setIsProcessing] = useState(false);
  const stripe = useStripe();
  const elements = useElements();

  const isPaymentReady = clientSecret && stripe && elements;

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();

    if (!isPaymentReady) return;

    setIsProcessing(true);

    const cardElement = elements.getElement(CardElement);

    if (!cardElement) return;

    const { error, paymentIntent } = await stripe.confirmCardPayment(
      clientSecret,
      {
        payment_method: {
          card: cardElement,
          billing_details: {
            // Optional: Add minimal billing details if needed
          },
        },
      },
    );

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
      <CardElement
        options={{
          style: {
            base: {
              fontSize: '16px',
              color: '#424770',
              '::placeholder': {
                color: '#aab7c4',
              },
            },
            invalid: {
              color: '#9e2146',
            },
          },
          hidePostalCode: true,
        }}
      />
      <Button
        type="submit"
        variant="contained"
        color="primary"
        disabled={isProcessing || !isPaymentReady}
        fullWidth
        sx={{ mt: 2 }}
      >
        {isProcessing ? <CircularProgress size={24} /> : `Pay $${amount}`}
      </Button>
    </form>
  );
};
const StripePaymentDialog: React.FC<StripePaymentDialogProps> = ({
  open,
  onClose,
  clientSecret,
  amount,
  currency = 'usd',
  loading = false,
  onPaymentSuccess,
  onPaymentError,
}) => {
  return (
    <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
      <DialogTitle>Complete Payment</DialogTitle>
      <DialogContent>
        {loading ? (
          <CircularProgress />
        ) : (
          <Elements stripe={stripePromise}>
            <PaymentForm
              clientSecret={clientSecret}
              amount={amount}
              currency={currency}
              onPaymentSuccess={onPaymentSuccess}
              onPaymentError={onPaymentError}
            />
          </Elements>
        )}
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
