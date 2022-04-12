import {
    Button,
    Dialog,
    TextField,
    Typography,
    DialogActions,
    DialogContentText,
    DialogTitle,
    DialogContent,
    IconButton,
    Stack,
} from "@mui/material";
import CloseIcon from "@mui/icons-material/Close";
import { useState } from "react";
import { colors } from "../../../theme";
import { makeStyles } from "@mui/styles";
import { globalUseStyles } from "../../../styles";

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
    const gclasses = globalUseStyles();
    const [nodeName, setNodeName] = useState("");
    const [nodeSerial, setNodeSerial] = useState("");

    const handleRegisterNode = () =>
        handleActivationSubmit({ name: nodeName, serial: nodeSerial });

    return (
        <Dialog open={isOpen} onClose={handleClose}>
            <Stack
                direction="row"
                alignItems="center"
                justifyContent="space-between"
            >
                <DialogTitle>{dialogTitle}</DialogTitle>
                <IconButton
                    onClick={handleClose}
                    sx={{ position: "relative", right: 8 }}
                >
                    <CloseIcon />
                </IconButton>
            </Stack>

            <DialogContent>
                <DialogContentText>
                    <Typography variant="body1" sx={{ color: colors.black }}>
                        {subTitle}
                    </Typography>
                </DialogContentText>
                <Stack direction="row" spacing={1} sx={{ mt: 3 }}>
                    <TextField
                        fullWidth
                        value={nodeName}
                        label={"NODE NAME"}
                        InputLabelProps={{ shrink: true }}
                        InputProps={{
                            classes: {
                                input: gclasses.inputFieldStyle,
                            },
                        }}
                        onChange={(e: any) => setNodeName(e.target.value)}
                    />

                    <TextField
                        fullWidth
                        value={nodeSerial}
                        label={"SERIAL NUMBER"}
                        InputLabelProps={{ shrink: true }}
                        InputProps={{
                            classes: {
                                input: gclasses.inputFieldStyle,
                            },
                        }}
                        onChange={(e: any) => setNodeSerial(e.target.value)}
                    />
                </Stack>
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
