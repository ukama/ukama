'use client';

import React, { useState, useMemo, useCallback } from 'react';
import {
  Tabs,
  Tab,
  Paper,
  Box,
  Typography,
  Stack,
  TextField,
  AlertColor,
  Button,
  Skeleton,
} from '@mui/material';
import LoadingWrapper from '@/components/LoadingWrapper';
import colors from '@/theme/colors';
import PaymentCard from '@/components/PaymentCard';
import { globalUseStyles } from '@/styles/global';
import {
  useGetReportsQuery,
  useGetPaymentsQuery,
} from '@/client/graphql/generated';
import { useAppContext } from '@/context';
import {
  Elements,
  CardElement,
  useStripe,
  useElements,
} from '@stripe/react-stripe-js';
import { loadStripe } from '@stripe/stripe-js';
import { Stripe, StripeElements } from '@stripe/stripe-js';
import { PAYMENT_METHODS } from '@/constants';

const stripeKey = process.env.NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY ?? '';

if (!stripeKey) {
  console.error(
    'Stripe publishable key is missing. Ensure NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY is set.',
  );
}
const stripePromise = stripeKey ? loadStripe(stripeKey) : null;

interface CurrentBillingProps {
  dataUsagePaid: number;
  notificationEmail: string;
  nextPaymentAmount: number;
  nextPaymentDate: string;
  clientSecret: string;
}

const StripePaymentForm: React.FC<{
  onStripeSubmit: (stripe: Stripe, elements: StripeElements) => Promise<void>;
}> = ({ onStripeSubmit }) => {
  const [isProcessing, setIsProcessing] = useState(false);
  const { setSnackbarMessage } = useAppContext();
  const stripe = useStripe();
  const elements = useElements();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!stripe || !elements) {
      setSnackbarMessage({
        id: 'payment-error',
        message: "Stripe.js hasn't loaded yet.",
        type: 'error' as AlertColor,
        show: true,
      });
      return;
    }

    setIsProcessing(true);

    try {
      await onStripeSubmit(stripe, elements);
    } catch (error) {
      console.error(error);
      setSnackbarMessage({
        id: 'payment-error',
        message: 'An unexpected error occurred',
        type: 'error' as AlertColor,
        show: true,
      });
    } finally {
      setIsProcessing(false);
    }
  };

  const cardElementOptions = {
    hidePostalCode: true,
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
  };

  return (
    <Box sx={{ mt: 2, p: 2 }}>
      <form onSubmit={handleSubmit}>
        <Box
          sx={{
            border: '1px solid #ccc',
            borderRadius: '4px',
            padding: '10px',
            marginBottom: '15px',
          }}
        >
          <CardElement options={cardElementOptions} />
        </Box>
        <Button
          type="submit"
          disabled={isProcessing || !stripe}
          variant="contained"
          sx={{
            width: '100%',
            background: isProcessing ? colors.black38 : colors.primaryMain,
          }}
        >
          {isProcessing ? 'Processing...' : 'Pay Now'}
        </Button>
      </form>
    </Box>
  );
};

const CurrentBilling: React.FC<CurrentBillingProps> = ({
  dataUsagePaid,
  notificationEmail,
  nextPaymentAmount,
  nextPaymentDate,
  clientSecret,
}) => {
  const gclasses = globalUseStyles();
  const { setSnackbarMessage } = useAppContext();
  const [selectedPaymentMethod, setSelectedPaymentMethod] = useState('Stripe');

  const handlePaymentMethodChange = useCallback((method: string) => {
    setSelectedPaymentMethod(method);
  }, []);

  const handleStripeSubmit = async (
    stripe: Stripe,
    elements: StripeElements,
  ) => {
    try {
      const cardElement = elements.getElement(CardElement);

      if (!cardElement) {
        throw new Error('Card Element not found');
      }

      const { error, paymentIntent } = await stripe.confirmCardPayment(
        clientSecret,
        {
          payment_method: {
            card: cardElement,
            billing_details: {
              email: notificationEmail,
            },
          },
        },
      );

      if (error) {
        console.error(error);
        setSnackbarMessage({
          id: 'payment-failed',
          message: `Payment failed. ${error?.message}`,
          type: 'error' as AlertColor,
          show: true,
        });
      } else if (paymentIntent?.status === 'succeeded') {
        setSnackbarMessage({
          id: 'payment-success',
          message: 'Payment completed successfully.',
          type: 'success' as AlertColor,
          show: true,
        });
      }
    } catch (error) {
      console.error(error);
      setSnackbarMessage({
        id: 'payment-error',
        message: 'An unexpected error occurred',
        type: 'error' as AlertColor,
        show: true,
      });
    }
  };

  return (
    <Box sx={{ py: 2 }}>
      {!clientSecret && (
        <Skeleton
          variant="rectangular"
          width="100%"
          height={50}
          sx={{ mb: 2 }}
        />
      )}
      <Elements stripe={stripePromise} options={{ clientSecret }}>
        <PaymentCard
          amount={nextPaymentAmount.toString()}
          startDate={nextPaymentDate}
          endDate={nextPaymentDate}
          onChangePaymentMethod={handlePaymentMethodChange}
          paymentMethods={PAYMENT_METHODS}
          onPaymentMethodSelect={handlePaymentMethodChange}
          isLoading={!stripePromise}
        >
          {selectedPaymentMethod === 'Stripe' && (
            <StripePaymentForm onStripeSubmit={handleStripeSubmit} />
          )}
        </PaymentCard>
      </Elements>

      <Paper
        elevation={2}
        sx={{
          p: 4,
          mt: 2,
          borderRadius: '10px',
          bgcolor: colors.white,
        }}
      >
        <Typography variant="h6">Data Usage</Typography>
        <Stack
          direction="row"
          justifyContent="space-between"
          alignItems="center"
        >
          <Typography variant="body2" sx={{ color: colors.black54 }}>
            Data usage paid for by subscribers
          </Typography>
          <Typography variant="h6">$ {dataUsagePaid}</Typography>
        </Stack>
      </Paper>

      <Paper
        elevation={2}
        sx={{
          p: 4,
          mt: 2,
          borderRadius: '10px',
          bgcolor: colors.white,
        }}
      >
        <Typography variant="h6">Notification Settings</Typography>
        <Stack direction="column" spacing={2}>
          <Typography variant="body2" sx={{ color: colors.black54 }}>
            All entered emails will receive receipts for the monthly bill
          </Typography>
          <TextField
            fullWidth
            label="PRIMARY EMAIL"
            value={notificationEmail}
            variant="outlined"
            disabled
            InputLabelProps={{ shrink: true }}
            InputProps={{
              classes: { input: gclasses.inputFieldStyle },
            }}
          />
        </Stack>
      </Paper>
    </Box>
  );
};

const BillingHistory: React.FC = () => (
  <Box sx={{ p: 3 }}>
    <Typography variant="h6">Billing History</Typography>
  </Box>
);

const BillingSettingsPage: React.FC = () => {
  const [currentTab, setCurrentTab] = useState(0);
  const [clientSecret, setClientSecret] = useState<string>();
  const { setSnackbarMessage, user } = useAppContext();

  const { data: paymentsData, loading: paymentsLoading } = useGetPaymentsQuery({
    variables: {
      data: {
        paymentMethod: 'stripe',
        status: 'processing',
        type: 'package',
      },
    },
    fetchPolicy: 'network-only',
    onError: (error) => {
      setSnackbarMessage({
        id: 'payments-error',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });

  const { data: reportsData, loading: reportsLoading } = useGetReportsQuery({
    variables: {
      data: {
        isPaid: false,
        report_type: 'invoice',
        count: 0,
        networkId: '',
        ownerId: '',
        ownerType: '',
        sort: false,
      },
    },
    fetchPolicy: 'network-only',
    onError: (error) => {
      setSnackbarMessage({
        id: 'reports-error',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });

  const totalAmountUSD = useMemo(() => {
    return (
      paymentsData?.getPayments.payments.reduce((total, payment) => {
        return total + parseFloat(payment.amount) * 100;
      }, 0) || 0
    );
  }, [paymentsData]);

  const handleTabChange = (_: React.SyntheticEvent, newValue: number) => {
    setCurrentTab(newValue);
  };

  return (
    <Elements stripe={stripePromise} options={{ clientSecret }}>
      <LoadingWrapper
        width="100%"
        radius="medium"
        isLoading={paymentsLoading || reportsLoading || !clientSecret}
        height="calc(100vh - 244px)"
      >
        <Tabs value={currentTab} onChange={handleTabChange}>
          <Tab label="Current Billing" />
          <Tab label="Billing History" />
        </Tabs>
        {currentTab === 0 && (
          <CurrentBilling
            dataUsagePaid={
              totalAmountUSD ? parseFloat(totalAmountUSD.toFixed(2)) : 0
            }
            notificationEmail={user.email}
            nextPaymentAmount={
              reportsData?.getReports?.reports[0]?.rawReport?.amountCents || 0
            }
            nextPaymentDate={
              reportsData?.getReports?.reports[0]?.rawReport?.issuingDate || ''
            }
            clientSecret={clientSecret}
          />
        )}

        {currentTab === 1 && <BillingHistory />}
      </LoadingWrapper>
    </Elements>
  );
};

export default BillingSettingsPage;
