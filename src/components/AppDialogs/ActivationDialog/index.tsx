import {
    Box,
    Grid,
    Button,
    Dialog,
    TextField,
    IconButton,
    Typography,
    DialogActions,
    DialogContent,
    Stack,
} from "@mui/material";
import { colors } from "../../../theme";
import { makeStyles } from "@mui/styles";
import { globalUseStyles } from "../../../styles";
import CloseIcon from "@mui/icons-material/Close";
import { ChangeEventHandler, useState } from "react";
import { UserActivateFormType } from "../../../types";

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
            backgroundColor: `${colors.darkGrey}`,
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
        <Box sx={{ p: "8px 0px" }}>
            <Grid item container spacing={3}>
                <Grid item container xs={12} spacing={1}>
                    <Grid item xs={12} md={6}>
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
                    </Grid>
                    <Grid item xs={12} md={6}>
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
                    </Grid>
                </Grid>
            </Grid>
        </Box>
    );
};

const FormFlowTwo = ({ values, handleChange }: FormContainerProps) => {
    const classes = globalUseStyles();
    return (
        <Box sx={{ p: "8px 0px" }}>
            <Grid item container spacing={3}>
                <Grid item xs={12}>
                    <TextField
                        fullWidth
                        id="securityCode"
                        name="securityCode"
                        label={"SECURITY CODE"}
                        value={values.securityCode}
                        onChange={handleChange}
                        InputLabelProps={{ shrink: true }}
                        InputProps={{
                            classes: {
                                input: classes.inputFieldStyle,
                            },
                        }}
                    />
                </Grid>
            </Grid>
        </Box>
    );
};

type ActivationDialogProps = {
    isOpen: boolean;
    subTitle: string;
    handleClose: any;
    subTitle2: string;
    dialogTitle: string;
};

const ActivationDialog = ({
    isOpen,
    subTitle,
    subTitle2,
    dialogTitle,
    handleClose,
}: ActivationDialogProps) => {
    const classes = useStyles();
    const [flowScreen, setFlowScreen] = useState(1);
    const [userActivateForm, setUserActivateForm] =
        useState<UserActivateFormType>(initialActivationFormValue);

    const handleNext = () => setFlowScreen(2);
    const handleBack = () => setFlowScreen(1);
    const handleChange = (e: any) => {
        setUserActivateForm({
            ...userActivateForm,
            [e.target.id]: e.target.value,
        });
    };

    const handleSubmitDisableBehaviour = () =>
        flowScreen === 1
            ? false
            : userActivateForm.nodeName && userActivateForm.serialNumber
            ? false
            : true;

    return (
        <Dialog open={isOpen} onClose={handleClose}>
            <Box
                sx={{
                    width: { xs: "100%", md: "560px" },
                    padding: "16px 24px",
                }}
            >
                <Box className={classes.basicDialogHeaderStyle}>
                    <Typography variant="h6">
                        {flowScreen === 1
                            ? dialogTitle
                            : `${dialogTitle} [Serial Number]`}
                    </Typography>

                    <IconButton
                        onClick={handleClose}
                        sx={{ ml: "24px", p: "0px" }}
                    >
                        <CloseIcon />
                    </IconButton>
                </Box>

                <DialogContent sx={{ p: "1px" }}>
                    {flowScreen === 1 ? (
                        <Stack spacing={2}>
                            <Typography variant="body1">{subTitle}</Typography>
                            <FormFlowOne
                                handleChange={handleChange}
                                values={userActivateForm}
                            />
                        </Stack>
                    ) : (
                        <Stack spacing={2}>
                            <Typography variant="body1">{subTitle2}</Typography>
                            <FormFlowTwo
                                handleChange={handleChange}
                                values={userActivateForm}
                            />
                        </Stack>
                    )}
                </DialogContent>
                <DialogActions className={classes.actionContainer}>
                    <Button
                        onClick={handleBack}
                        sx={{
                            visibility: flowScreen === 2 ? "visible" : "hidden",
                        }}
                    >
                        Back
                    </Button>
                    <div>
                        <Button
                            variant="text"
                            sx={{ mr: "20px" }}
                            onClick={handleClose}
                        >
                            Cancel
                        </Button>
                        <Button
                            variant="contained"
                            onClick={handleNext}
                            className={classes.stepButtonStyle}
                            disabled={handleSubmitDisableBehaviour()}
                        >
                            {flowScreen === 2 ? "Add node" : "Next"}
                        </Button>
                    </div>
                </DialogActions>
            </Box>
        </Dialog>
    );
};

export default ActivationDialog;
