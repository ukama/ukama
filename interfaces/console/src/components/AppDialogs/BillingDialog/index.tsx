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
import { HorizontalContainerJustify } from "../../../styles";

interface IBillingDialog {
    isOpen: boolean;
    handleCloseAction: Function;
    handleSuccessAction: Function;
}

const DialogList = [
    {
        id: 0,
        title: "Choose roaming plan",
        description:
            "Choose a roaming plan below, and it will apply to all residents that have roaming enabled. Your selection can always be changed later.",
    },
    {
        id: 1,
        title: "Enter payment information",
        description: "Enter your payment information",
    },
    {
        id: 2,
        title: "Customize preferences",
        description: "Monitor and budget data usage with these settings.",
    },
    {
        id: 3,
        title: "Payment set up successfully ",
        description:
            "Your payment and preferences have been set up successfully! You can change your settings at any time.",
    },
];

const BillingDialog = ({
    isOpen,
    handleCloseAction,
    handleSuccessAction,
}: IBillingDialog) => {
    const [flow, setFlow] = useState(0);

    const handleFlowChange = (i: number) => {
        if (flow === 2) handleSuccessAction();
        setFlow(i);
    };

    const handleClose = () => {
        setFlow(0);
        handleCloseAction();
    };

    return (
        <Dialog
            fullWidth
            open={isOpen}
            maxWidth="sm"
            onClose={() => handleClose()}
            aria-labelledby="alert-dialog-title"
            aria-describedby="alert-dialog-description"
            onBackdropClick={() => handleClose()}
        >
            <Stack
                direction="row"
                alignItems="center"
                justifyContent="space-between"
            >
                <DialogTitle>{DialogList[flow].title}</DialogTitle>
                <IconButton
                    onClick={() => handleClose()}
                    sx={{ position: "relative", right: 8 }}
                >
                    <CloseIcon />
                </IconButton>
            </Stack>

            <DialogContent>
                <Typography variant="body1">
                    {DialogList[flow].description}
                </Typography>
                {flow === 0 && <ChoosePlan />}
                {flow === 1 && <PaymentForm />}
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
