import { HorizontalContainerJustify } from '@/styles/global';
import { Button, CircularProgress, Stack, useTheme } from '@mui/material';
import {
  CardElement,
  Elements,
  useElements,
  useStripe,
} from '@stripe/react-stripe-js';
import { loadStripe } from '@stripe/stripe-js';
interface ICheckoutForm {
  loading: boolean;
  isPaymentOnly: boolean;
  handleBackAction: Function;
  handleIsPaymentSuccess: Function;
}

interface IPaymentForm {
  loading: boolean;
  isPaymentOnly: boolean;
  handleBackAction: Function;
  handleIsPaymentSuccess: Function;
}

const stripePromise = loadStripe(
  'pk_test_51LN9vGHBOiFTwZOsILdYKGAyT3JpOJt55PLXT7RgcwMrezgETce1GDYP3iEFIQCy6OsgS51Z0B1lVorApjBwqkMu001gz6uBbS',
);

const CheckoutForm = ({
  loading,
  isPaymentOnly,
  handleBackAction,
  handleIsPaymentSuccess,
}: ICheckoutForm) => {
  const theme = useTheme();
  const stripe = useStripe();
  const elements = useElements();

  const handleSubmit = async (event: any) => {
    event.preventDefault();

    if (elements === null) {
      return;
    }

    const cardElement = elements.getElement(CardElement);

    if (cardElement) {
      await stripe
        ?.createPaymentMethod({
          type: 'card',
          card: cardElement,
        })
        .then((res: any) => {
          if (!res.error) {
            handleIsPaymentSuccess(res.paymentMethod.id);
          }
        });
    }
  };

  return (
    <form onSubmit={handleSubmit} style={{ marginTop: 16 }}>
      <CardElement
        options={{
          style: {
            base: {
              iconColor: theme.palette.text.primary,
              color: theme.palette.text.primary,
              fontFamily: 'Work Sans, sans-serif',
              fontSize: '16px',
              fontWeight: '500',
              fontSmoothing: 'antialiased',
              ':-webkit-autofill': {
                color: '#fce883',
              },
              '::placeholder': {
                color: theme.palette.text.disabled,
              },
            },

            invalid: {
              iconColor: theme.palette.error.main,
              color: theme.palette.error.main,
            },
          },
        }}
      />

      <HorizontalContainerJustify alignItems={'center'} mt={4}>
        <Button
          variant="text"
          color={'primary'}
          sx={{ visibility: isPaymentOnly ? 'hidden' : 'visible' }}
          onClick={() => handleBackAction()}
        >
          Back
        </Button>
        {loading ? (
          <CircularProgress />
        ) : (
          <Stack direction={'row'} alignItems="center" spacing={2}>
            <button
              type="submit"
              disabled={!stripe || !elements}
              style={{
                fontSize: 16,
                width: '100%',
                height: '42px',
                border: 'none',
                color: 'white',
                borderRadius: 4,
                fontWeight: 600,
                cursor: 'pointer',
                paddingLeft: '8px',
                paddingRight: '8px',
                fontFamily: 'Rubik',
                letterSpacing: '0.4px',
                textTransform: 'uppercase',
                backgroundColor: theme.palette.primary.main,
                boxShadow:
                  '0px 3px 1px -2px rgb(0 0 0 / 20%), 0px 2px 2px rgb(0 0 0 / 14%), 0px 1px 5px rgb(0 0 0 / 12%)',
              }}
            >
              Submit payment information
            </button>
          </Stack>
        )}
      </HorizontalContainerJustify>
    </form>
  );
};

const PaymentForm = ({
  loading,
  isPaymentOnly,
  handleBackAction,
  handleIsPaymentSuccess,
}: IPaymentForm) => {
  return (
    <Elements stripe={stripePromise}>
      <CheckoutForm
        loading={loading}
        isPaymentOnly={isPaymentOnly}
        handleBackAction={handleBackAction}
        handleIsPaymentSuccess={handleIsPaymentSuccess}
      />
    </Elements>
  );
};

export default PaymentForm;
