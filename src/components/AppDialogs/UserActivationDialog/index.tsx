import { makeStyles } from "@mui/styles";
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
import { SimCardDesign } from "../..";
import { colors } from "../../../theme";
import { globalUseStyles } from "../../../styles";
import CloseIcon from "@mui/icons-material/Close";
import { ChangeEventHandler, useState } from "react";
import { SimActivateFormType } from "../../../types";
import { SimCardData } from "../../../constants/stubData";

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

const initialSimActivateFormValue = {
    number: "",
    firstName: "",
    lastName: "",
    email: "",
    phone: "",
};

type FormContainerProps = {
    values: any;
    handleChange: ChangeEventHandler<HTMLInputElement>;
};

const FormContainer = ({ values, handleChange }: FormContainerProps) => {
    const classes = globalUseStyles();
    return (
        <Box sx={{ p: "8px 0px" }}>
            <Grid item container spacing={3}>
                <Grid item xs={12}>
                    <TextField
                        fullWidth
                        id="number"
                        name="number"
                        disabled={true}
                        variant="standard"
                        value={values.number}
                        label={"ESIM NUMBER"}
                        onChange={handleChange}
                        InputLabelProps={{ shrink: true }}
                        InputProps={{
                            disableUnderline: true,
                        }}
                    />
                </Grid>
                <Grid item container xs={12} spacing={1}>
                    <Grid item xs={12} md={6}>
                        <TextField
                            fullWidth
                            id="firstName"
                            name="firstName"
                            label={"FIRST NAME"}
                            onChange={handleChange}
                            value={values.firstName}
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
                            id="lastName"
                            name="lastName"
                            label={"LAST NAME"}
                            value={values.lastName}
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
                <Grid item xs={12}>
                    <TextField
                        fullWidth
                        id="email"
                        name="email"
                        value={values.email}
                        label={"CONTACT EMAIL - Optional"}
                        onChange={handleChange}
                        InputLabelProps={{ shrink: true }}
                        InputProps={{
                            classes: {
                                input: classes.inputFieldStyle,
                            },
                        }}
                    />
                </Grid>
                <Grid item xs={12}>
                    <TextField
                        fullWidth
                        id="phone"
                        name="phone"
                        value={values.phone}
                        label={"CONTACT PHONE - Optional"}
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

type UserActivationDialogProps = {
    dialogTitle: string;
    subTitle: string;
    isOpen: boolean;
    handleClose: any;
};

const UserActivationDialog = ({
    isOpen,
    subTitle,
    dialogTitle,
    handleClose,
}: UserActivationDialogProps) => {
    const classes = useStyles();
    const [simActivateForm, setSimActivateForm] = useState<SimActivateFormType>(
        initialSimActivateFormValue
    );
    const [flowScreen, setFlowScreen] = useState(1);
    const [selectedSim, setSelectedSim] = useState<number | null>(null);

    const handleSimCardClick = (id: number) => setSelectedSim(id);
    const handleNext = () => {
        setFlowScreen(2);
        setSimActivateForm({
            ...simActivateForm,
            number: SimCardData.filter(i => i.id === selectedSim)[0].serial,
        });
    };
    const handleBack = () => setFlowScreen(1);
    const handleChange = (e: any) => {
        setSimActivateForm({
            ...simActivateForm,
            [e.target.id]: e.target.value,
        });
    };

    const handleSubmitDisableBehaviour = () =>
        selectedSim && flowScreen === 1
            ? false
            : simActivateForm.firstName && simActivateForm.lastName
            ? false
            : true;

    return (
        <Dialog open={isOpen} onClose={handleClose}>
            <Box
                sx={{
                    width: { xs: "100%", md: "600px" },
                    padding: "16px 24px",
                }}
            >
                <Box className={classes.basicDialogHeaderStyle}>
                    <Typography variant="h6">{dialogTitle}</Typography>

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
                            {SimCardData.map(
                                ({ id, title, serial, isActive }) => (
                                    <SimCardDesign
                                        id={id}
                                        key={id}
                                        title={title}
                                        serial={serial}
                                        isActivate={isActive}
                                        isSelected={id === selectedSim}
                                        handleItemClick={handleSimCardClick}
                                    />
                                )
                            )}
                        </Stack>
                    ) : (
                        <FormContainer
                            handleChange={handleChange}
                            values={simActivateForm}
                        />
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
                            {flowScreen === 2 ? "Activate" : "Next"}
                        </Button>
                    </div>
                </DialogActions>
            </Box>
        </Dialog>
    );
};

export default UserActivationDialog;
