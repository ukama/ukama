import {
    Button,
    Dialog,
    IconButton,
    Typography,
    DialogActions,
    DialogContent,
    Stack,
    DialogTitle,
} from "@mui/material";
import { SimCardDesign } from "../..";
import { colors } from "../../../theme";
import { SimCardData } from "../../../constants";
import CloseIcon from "@mui/icons-material/Close";
import { useState } from "react";

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
    const [selectedSim, setSelectedSim] = useState<number | null>(null);

    const handleSimCardClick = (id: number) => setSelectedSim(id);

    return (
        <Dialog open={isOpen} onClose={handleClose} maxWidth="md" fullWidth>
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
                    sx={{ color: colors.primaryMain }}
                    onClick={handleClose}
                >
                    Close
                </Button>
            </DialogActions>
        </Dialog>
    );
};

export default UserActivationDialog;
