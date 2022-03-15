import { makeStyles } from "@mui/styles";
import {
    Box,
    Button,
    Dialog,
    IconButton,
    Typography,
    DialogActions,
    DialogContent,
    Stack,
} from "@mui/material";
import { SimCardDesign } from "../..";
import { colors } from "../../../theme";
import { SimCardData } from "../../../constants";
import CloseIcon from "@mui/icons-material/Close";
import { useState } from "react";

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
            backgroundColor: colors.primaryMain,
        },
    },
}));

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
    const [selectedSim, setSelectedSim] = useState<number | null>(null);

    const handleSimCardClick = (id: number) => setSelectedSim(id);

    return (
        <Dialog open={isOpen} onClose={handleClose}>
            <Box
                component="div"
                sx={{
                    width: { xs: "100%", md: "600px" },
                    padding: "16px 24px",
                }}
            >
                <Box component="div" className={classes.basicDialogHeaderStyle}>
                    <Typography variant="h6">{dialogTitle}</Typography>

                    <IconButton
                        onClick={handleClose}
                        sx={{ ml: "24px", p: "0px" }}
                    >
                        <CloseIcon />
                    </IconButton>
                </Box>

                <DialogContent sx={{ p: "1px" }}>
                    <Stack spacing={2}>
                        <Typography variant="body1">{subTitle}</Typography>
                        {SimCardData.map(({ id, title, serial, isActive }) => (
                            <SimCardDesign
                                id={id}
                                key={id}
                                title={title}
                                serial={serial}
                                isActivate={isActive}
                                isSelected={id === selectedSim}
                                handleItemClick={handleSimCardClick}
                            />
                        ))}
                    </Stack>
                </DialogContent>
                <DialogActions>
                    <Button
                        sx={{ color: colors.primaryMain, mt: 2 }}
                        onClick={handleClose}
                    >
                        Close
                    </Button>
                </DialogActions>
            </Box>
        </Dialog>
    );
};

export default UserActivationDialog;
