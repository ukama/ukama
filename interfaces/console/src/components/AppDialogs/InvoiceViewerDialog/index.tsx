import {
    Button,
    Dialog,
    IconButton,
    Typography,
    DialogActions,
    DialogContent,
    Stack,
    DialogTitle,
    Grid,
} from "@mui/material";
import colors from "../../../theme/colors";
import CloseIcon from "@mui/icons-material/Close";
import React from "react";

type BasicDialogProps = {
    isOpen: boolean;
    data: any;
};
const Logo = React.lazy(() =>
    import("../../../assets/svg").then(module => ({
        default: module.Logo,
    }))
);

const InvoiceViewerDialog = ({ isOpen, data }: BasicDialogProps) => {
    return (
        <Dialog
            fullWidth
            open={isOpen}
            maxWidth="sm"
            // onClose={handleCloseAction}
            aria-labelledby="alert-dialog-title"
            aria-describedby="alert-dialog-description"
            // onBackdropClick={() => isClosable && handleCloseAction()}
        >
            <Stack
                direction="row"
                alignItems="center"
                justifyContent="space-between"
            >
                {/* <DialogTitle>{title}</DialogTitle> */}
                {/* <IconButton
                    onClick={handleCloseAction}
                    sx={{ position: "relative", right: 8 }}
                >
                    <CloseIcon />
                </IconButton> */}
            </Stack>

            <DialogContent>
                {/* <Typography variant="body1">{description}</Typography> */}
                <Grid container spacing={1}>
                    <Grid item xs={6}>
                        <Logo
                            width={"100%"}
                            height={"36px"}
                            color={colors.primaryMain}
                        />
                    </Grid>
                    <Grid item xs={6}>
                        <Stack direction="column" spacing={1}>
                            <Stack direction="row" spacing={3}>
                                <Typography variant="body2">
                                    {"title"}
                                </Typography>
                                <Typography variant="body2">
                                    {"data"}
                                </Typography>
                            </Stack>
                        </Stack>
                    </Grid>
                </Grid>
            </DialogContent>

            {/* <DialogActions>
                <Stack direction={"row"} alignItems="center" spacing={2}>
                    {labelNegativeBtn && (
                        <Button
                            variant="text"
                            color={"primary"}
                            onClick={handleCloseAction}
                        >
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
            </DialogActions> */}
        </Dialog>
    );
};

export default InvoiceViewerDialog;
