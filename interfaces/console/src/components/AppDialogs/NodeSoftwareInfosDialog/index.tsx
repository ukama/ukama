import { makeStyles } from "@mui/styles";
import CloseIcon from "@mui/icons-material/Close";
import {
    Box,
    Button,
    Dialog,
    IconButton,
    Typography,
    Stack,
    DialogActions,
    DialogContent,
} from "@mui/material";
const useStyles = makeStyles(() => ({
    basicDialogHeaderStyle: {
        display: "flex",
        flexDirection: "row",
        alignItems: "center",
        justifyContent: "space-between",
    },
}));

type BasicDialogProps = {
    isOpen: boolean;
    isClosable?: boolean;
    handleClose: Function;
    closeBtnLabel?: string;
};

const NodeSoftwareInfosDialog = ({
    isOpen = false,
    handleClose,
    closeBtnLabel,
    isClosable = true,
}: BasicDialogProps) => {
    const classes = useStyles();

    return (
        <Dialog
            open={isOpen}
            onBackdropClick={() => isClosable && handleClose()}
        >
            <Box
                component="div"
                sx={{
                    width: { xs: "100%", md: "500px" },
                    padding: "16px 24px",
                }}
            >
                <Box component="div" className={classes.basicDialogHeaderStyle}>
                    <Typography variant="h6">{"Update Notes 12.4"}</Typography>
                    {isClosable && (
                        <IconButton
                            onClick={() => handleClose()}
                            sx={{ ml: "24px", p: "8px" }}
                        >
                            <CloseIcon />
                        </IconButton>
                    )}
                </Box>
                <DialogContent sx={{ p: 0, my: 2 }}>
                    <Typography variant="body1">Short introduction.</Typography>

                    <Stack direction="column" sx={{ my: 2 }}>
                        <Typography variant="body1">TL;DR</Typography>
                        <Typography variant="body1">*** NEW ***</Typography>
                        <Typography variant="body1">
                            *** IMPROVEMENTS
                        </Typography>
                        <Typography variant="body1">*** FIXES ***</Typography>
                    </Stack>
                </DialogContent>
                <DialogActions sx={{ padding: 0 }}>
                    <Button
                        onClick={() => handleClose()}
                        sx={{
                            mr: 2,
                            left: 7,
                            position: "relative",
                            justifyItems: "center",
                        }}
                    >
                        {closeBtnLabel}
                    </Button>
                </DialogActions>
            </Box>
        </Dialog>
    );
};

export default NodeSoftwareInfosDialog;
