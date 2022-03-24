import { makeStyles } from "@mui/styles";
import CloseIcon from "@mui/icons-material/Close";
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

const useStyles = makeStyles(() => ({
    basicDialogHeaderStyle: {
        display: "flex",
        flexDirection: "row",
        alignItems: "center",
        justifyContent: "space-between",
    },
}));

type BasicDialogProps = {
    title: string;
    isOpen: boolean;
    content: string;
    btnLabel: string;
    handleClose: any;
    isClosable?: boolean;
    btnVariant?: "text" | "outlined" | "contained";
};

const BasicDialog = ({
    title,
    isOpen,
    content,
    btnLabel,
    handleClose,
    btnVariant = "text",
    isClosable = true,
}: BasicDialogProps) => {
    const classes = useStyles();
    return (
        <Dialog
            open={isOpen}
            onClose={handleClose}
            aria-labelledby="alert-dialog-title"
            aria-describedby="alert-dialog-description"
            onBackdropClick={() => isClosable && handleClose()}
        >
            <Stack
                spacing={3}
                sx={{
                    width: { xs: "100%", md: "500px" },
                    padding: "16px 8px 8px 24px",
                }}
            >
                <Box component="div" className={classes.basicDialogHeaderStyle}>
                    <Typography variant="h6">{title}</Typography>
                    {isClosable && (
                        <IconButton
                            onClick={handleClose}
                            sx={{ ml: "24px", p: "8px" }}
                        >
                            <CloseIcon />
                        </IconButton>
                    )}
                </Box>
                <DialogContent sx={{ p: 0 }}>
                    <Typography variant="body1">{content}</Typography>
                </DialogContent>
                <DialogActions>
                    <Button variant={btnVariant} onClick={handleClose}>
                        {btnLabel}
                    </Button>
                </DialogActions>
            </Stack>
        </Dialog>
    );
};

export default BasicDialog;
