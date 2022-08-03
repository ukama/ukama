import {
    Elements,
    useStripe,
    useElements,
    CardElement,
} from "@stripe/react-stripe-js";
import { loadStripe } from "@stripe/stripe-js";
import { Button, Stack, useTheme } from "@mui/material";
import { HorizontalContainerJustify } from "../../../styles";
interface ICheckoutForm {
    isPaymentOnly: boolean;
    handleBackAction: Function;
    handleCloseAction: Function;
    handleIsPaymentSuccess: Function;
}

interface IPaymentForm {
    isPaymentOnly: boolean;
    handleBackAction: Function;
    handleCloseAction: Function;
    handleIsPaymentSuccess: Function;
}

const stripePromise = loadStripe("pk_test_6pRNASCoBOKtIshFeQd4XMUh");

const CheckoutForm = ({
    isPaymentOnly,
    handleBackAction,
    handleCloseAction,
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
                    type: "card",
                    card: cardElement,
                })
                .then(res => {
                    if (!res.error) handleIsPaymentSuccess();
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

            <HorizontalContainerJustify alignItems={"center"} mt={4}>
                <Button
                    variant="text"
                    color={"primary"}
                    sx={{ visibility: isPaymentOnly ? "hidden" : "visible" }}
                    onClick={() => handleBackAction()}
                >
                    Back
                </Button>

                <Stack direction={"row"} alignItems="center" spacing={2}>
                    <Button
                        variant={"text"}
                        color={"primary"}
                        onClick={() => handleCloseAction()}
                    >
                        Close
                    </Button>

                    <button
                        type="submit"
                        disabled={!stripe || !elements}
                        style={{
                            fontSize: 16,
                            width: "100%",
                            height: "42px",
                            border: "none",
                            color: "white",
                            borderRadius: 4,
                            fontWeight: 600,
                            cursor: "pointer",
                            paddingLeft: "8px",
                            paddingRight: "8px",
                            fontFamily: "Rubik",
                            letterSpacing: "0.4px",
                            textTransform: "uppercase",
                            backgroundColor: theme.palette.primary.main,
                            boxShadow:
                                "0px 3px 1px -2px rgb(0 0 0 / 20%), 0px 2px 2px rgb(0 0 0 / 14%), 0px 1px 5px rgb(0 0 0 / 12%)",
                        }}
                    >
                        Submit payment information
                    </button>
                </Stack>
            </HorizontalContainerJustify>
        </form>
    );
};

const PaymentForm = ({
    isPaymentOnly,
    handleBackAction,
    handleCloseAction,
    handleIsPaymentSuccess,
}: IPaymentForm) => {
    return (
        <Elements stripe={stripePromise}>
            <CheckoutForm
                isPaymentOnly={isPaymentOnly}
                handleBackAction={handleBackAction}
                handleCloseAction={handleCloseAction}
                handleIsPaymentSuccess={handleIsPaymentSuccess}
            />
        </Elements>
    );
};

export default PaymentForm;
