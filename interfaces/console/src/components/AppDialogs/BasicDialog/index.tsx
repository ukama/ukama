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
import CloseIcon from "@mui/icons-material/Close";

type BasicDialogProps = {
    title: string;
    isOpen: boolean;
    description: string;
    isClosable?: boolean;
    handleCloseAction: any;
    labelSuccessBtn?: string;
    handleSuccessAction?: any;
    labelNegativeBtn?: string;
};

const BasicDialog = ({
    title,
    isOpen,
    description,
    labelSuccessBtn,
    labelNegativeBtn,
    handleCloseAction,
    isClosable = true,
    handleSuccessAction,
}: BasicDialogProps) => {
    return (
        <Dialog
            fullWidth
            open={isOpen}
            maxWidth="sm"
            onClose={handleCloseAction}
            aria-labelledby="alert-dialog-title"
            aria-describedby="alert-dialog-description"
            onBackdropClick={() => isClosable && handleCloseAction()}
        >
            <Stack
                direction="row"
                alignItems="center"
                justifyContent="space-between"
            >
                <DialogTitle>{title}</DialogTitle>
                <IconButton
                    onClick={handleCloseAction}
                    sx={{ position: "relative", right: 8 }}
                >
                    <CloseIcon />
                </IconButton>
            </Stack>

            <DialogContent>
                <Typography variant="body1">{description}</Typography>
            </DialogContent>

            <DialogActions>
                <Stack direction={"row"} alignItems="center" spacing={2}>
                    {labelNegativeBtn && (
                        <Button variant="text" onClick={handleCloseAction}>
                            {labelNegativeBtn}
                        </Button>
                    )}
                    {labelSuccessBtn && (
                        <Button
                            variant="contained"
                            onClick={() =>
                                handleSuccessAction
                                    ? handleSuccessAction()
                                    : handleCloseAction()
                            }
                        >
                            {labelSuccessBtn}
                        </Button>
                    )}
                </Stack>
            </DialogActions>
        </Dialog>
    );
};

export default BasicDialog;
