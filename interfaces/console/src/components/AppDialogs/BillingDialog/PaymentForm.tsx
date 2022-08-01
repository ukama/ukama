import {
    CardElement,
    Elements,
    useStripe,
    useElements,
} from "@stripe/react-stripe-js";
import { useTheme } from "@mui/material";
import { loadStripe } from "@stripe/stripe-js";
interface ICheckoutForm {
    handleIsPaymentSuccess: Function;
}

interface IPaymentForm {
    handleIsPaymentSuccess: Function;
}

const stripePromise = loadStripe("pk_test_6pRNASCoBOKtIshFeQd4XMUh");

const CheckoutForm = ({ handleIsPaymentSuccess }: ICheckoutForm) => {
    const theme = useTheme();
    const stripe = useStripe();
    const elements = useElements();

    const handleSubmit = async (event: any) => {
        event.preventDefault();

        if (elements === null) {
            return;
        }

        const cardElement = elements.getElement(CardElement);

        if (cardElement)
            await stripe?.createPaymentMethod({
                type: "card",
                card: cardElement,
            });
    };

    return (
        <form onSubmit={handleSubmit} style={{ marginTop: 16 }}>
            <CardElement
                options={{
                    style: {
                        base: {
                            iconColor: theme.palette.text.primary,
                            color: theme.palette.text.primary,
                            fontFamily: "Work Sans, sans-serif",
                            fontSize: "16px",
                            fontWeight: "500",
                            fontSmoothing: "antialiased",
                            ":-webkit-autofill": {
                                color: "#fce883",
                            },
                            "::placeholder": {
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

            <button
                type="submit"
                disabled={!stripe || !elements}
                style={{
                    marginTop: 16,
                    width: "100%",
                    height: "42px",
                    border: "none",
                    color: "white",
                    borderRadius: 4,
                    fontSize: 18,
                    fontWeight: 600,
                    textTransform: "capitalize",
                    backgroundColor: theme.palette.primary.main,
                    boxShadow:
                        "0px 3px 1px -2px rgb(0 0 0 / 20%), 0px 2px 2px rgb(0 0 0 / 14%), 0px 1px 5px rgb(0 0 0 / 12%)",
                }}
                onClick={() => handleIsPaymentSuccess(true)}
            >
                Pay
            </button>
        </form>
    );
};

const PaymentForm = ({ handleIsPaymentSuccess }: IPaymentForm) => {
    return (
        <Elements stripe={stripePromise}>
            <CheckoutForm
                handleIsPaymentSuccess={(v: boolean) =>
                    handleIsPaymentSuccess(v)
                }
            />
        </Elements>
    );
};

export default PaymentForm;
