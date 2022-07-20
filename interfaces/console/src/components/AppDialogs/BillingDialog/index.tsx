import {
    Stack,
    Button,
    Dialog,
    IconButton,
    Typography,
    DialogTitle,
    DialogContent,
    DialogActions,
} from "@mui/material";
import { useState } from "react";
import ChoosePlan from "./ChoosePlan";
import PaymentForm from "./PaymentForm";
import CustomizePref from "./CustomizePref";
import CloseIcon from "@mui/icons-material/Close";
import { BillingDialogList } from "../../../constants";
import { HorizontalContainerJustify } from "../../../styles";

interface IBillingDialog {
    isOpen: boolean;
    handleCloseAction: Function;
    handleSuccessAction: Function;
}

const BillingDialog = ({
    isOpen,
    handleCloseAction,
    handleSuccessAction,
}: IBillingDialog) => {
    const [flow, setFlow] = useState(0);
    const [isPaymentSuccess, setIsPaymentSuccess] = useState(false);

    const handleFlowChange = (i: number) => {
        if (flow === 2) handleSuccessAction();
        setFlow(i);
    };

    const handleClose = () => {
        setFlow(0);
        handleCloseAction();
    };

    const handleIsPaymentSuccess = (isSuccess: boolean) =>
        setIsPaymentSuccess(isSuccess);

    const isNextDiable = () => (flow === 1 && !isPaymentSuccess ? true : false);

    return (
        <Dialog
            fullWidth
            open={isOpen}
            maxWidth="sm"
            aria-labelledby="alert-dialog-title"
            aria-describedby="alert-dialog-description"
        >
            <Stack
                direction="row"
                alignItems="center"
                justifyContent="space-between"
            >
                <DialogTitle>{BillingDialogList[flow].title}</DialogTitle>
                <IconButton
                    onClick={() => handleClose()}
                    sx={{ position: "relative", right: 8 }}
                >
                    <CloseIcon />
                </IconButton>
            </Stack>

            <DialogContent>
                <Typography variant="body1">
                    {BillingDialogList[flow].description}
                </Typography>
                {flow === 0 && <ChoosePlan />}
                {flow === 1 && (
                    <PaymentForm
                        handleIsPaymentSuccess={handleIsPaymentSuccess}
                    />
                )}
                {flow === 2 && <CustomizePref />}
                {flow === 3 && <></>}
            </DialogContent>

            <DialogActions>
                <HorizontalContainerJustify>
                    <Button
                        variant="text"
                        color={"primary"}
                        sx={{
                            visibility:
                                flow !== 0 && flow !== 3 ? "visible" : "hidden",
                        }}
                        onClick={() => handleFlowChange(flow - 1)}
                    >
                        Back
                    </Button>

                    <Stack direction={"row"} alignItems="center" spacing={2}>
                        <Button
                            variant={flow === 3 ? "contained" : "text"}
                            color={"primary"}
                            onClick={() => handleClose()}
                        >
                            Close
                        </Button>

                        {flow !== 3 && (
                            <Button
                                variant="contained"
                                disabled={isNextDiable()}
                                onClick={() => handleFlowChange(flow + 1)}
                            >
                                {flow === 2 ? "Save" : "Next"}
                            </Button>
                        )}
                    </Stack>
                </HorizontalContainerJustify>
            </DialogActions>
        </Dialog>
    );
};

export default BillingDialog;
