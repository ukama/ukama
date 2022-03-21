import {
    Button,
    Dialog,
    TextField,
    Typography,
    DialogActions,
    DialogContentText,
    DialogTitle,
    DialogContent,
    Stack,
} from "@mui/material";
import { makeStyles } from "@mui/styles";
import { globalUseStyles } from "../../../styles";
import { ChangeEventHandler, useState } from "react";
import { UserActivateFormType } from "../../../types";
import { colors } from "../../../theme";

const useStyles = makeStyles(() => ({
    basicDialogHeaderStyle: {
        padding: "0px 0px 18px 0px",
        display: "flex",
        flexDirection: "row",
        alignItems: "center",
        justifyContent: "space-between",
    },
    actionContainer: {
        padding: "0px",
        marginTop: "16px",
        justifyContent: "space-between",
    },
    stepButtonStyle: {
        "&:disabled": {
            color: colors.white,
            backgroundColor: colors.nightGrey,
        },
    },
}));

const initialActivationFormValue = {
    nodeName: "",
    serialNumber: "",
    securityCode: "",
};

type FormContainerProps = {
    values: any;
    handleChange: ChangeEventHandler<HTMLInputElement>;
};

const FormFlowOne = ({ values, handleChange }: FormContainerProps) => {
    const classes = globalUseStyles();
    return (
        <Stack direction="row" spacing={1} sx={{ mt: 3 }}>
            <TextField
                fullWidth
                id="nodeName"
                name="nodeName"
                label={"NODE NAME"}
                onChange={handleChange}
                value={values.nodeName}
                InputLabelProps={{ shrink: true }}
                InputProps={{
                    classes: {
                        input: classes.inputFieldStyle,
                    },
                }}
            />

            <TextField
                fullWidth
                id="serialNumber"
                name="serialNumber"
                label={"SERIAL NUMBER"}
                onChange={handleChange}
                value={values.serialNumber}
                InputLabelProps={{ shrink: true }}
                InputProps={{
                    classes: {
                        input: classes.inputFieldStyle,
                    },
                }}
            />
        </Stack>
    );
};

type ActivationDialogProps = {
    isOpen: boolean;
    subTitle: string;
    handleClose: any;
    subTitle2?: string;
    dialogTitle: string;
    handleActivationSubmit: Function;
};

const ActivationDialog = ({
    isOpen,
    subTitle,
    dialogTitle,
    handleClose,
    handleActivationSubmit,
}: ActivationDialogProps) => {
    const classes = useStyles();
    const [userActivateForm, setUserActivateForm] =
        useState<UserActivateFormType>(initialActivationFormValue);

    const handleRegisterNode = () => {
        handleActivationSubmit(userActivateForm);
    };

    const handleChange = (e: any) => {
        setUserActivateForm({
            ...userActivateForm,
            [e.target.id]: e.target.value,
        });
    };

    return (
        <Dialog open={isOpen} onClose={handleClose}>
            <DialogTitle>{dialogTitle}</DialogTitle>
            <DialogContent>
                <DialogContentText>
                    <Typography variant="body1" sx={{ color: colors.black }}>
                        {subTitle}
                    </Typography>
                </DialogContentText>
                <FormFlowOne
                    handleChange={handleChange}
                    values={userActivateForm}
                />
            </DialogContent>
            <DialogActions sx={{ mr: 2, paddingBottom: 3 }}>
                <Button
                    sx={{ color: colors.primaryMain, mr: 2 }}
                    onClick={handleClose}
                >
                    Cancel
                </Button>
                <Button
                    variant="contained"
                    type="submit"
                    onClick={handleRegisterNode}
                    className={classes.stepButtonStyle}
                >
                    REGISTER NODE
                </Button>
            </DialogActions>
        </Dialog>
    );
};

export default ActivationDialog;
