import * as React from "react";
import { useTranslation } from "react-i18next";
import CloseIcon from "@mui/icons-material/Close";
import {
    Box,
    Button,
    IconButton,
    DialogTitle,
    DialogContentText,
    DialogContent,
    DialogActions,
    Dialog,
    Grid,
    Typography,
} from "@mui/material";
import "../../../i18n/i18n";

type FormDialogProps = {
    buttonProps?: {
        size: "medium";
        type: "submit";
    };
    dialogTitle?: string;
    dialogContent?: string;
    showBackButton?: boolean;
    submitButtonLabel?: string;
    open: boolean;
    onClose?: () => void;
    children?: React.ReactElement;
};
const FormDialog = ({
    dialogTitle,
    buttonProps,
    dialogContent,
    children,
    showBackButton,
    submitButtonLabel,
    open,
    onClose,
}: FormDialogProps) => {
    const { t } = useTranslation();
    return (
        <div>
            <Dialog open={open} onClose={onClose}>
                <DialogTitle>
                    <Box display="flex" alignItems="center">
                        <Box flexGrow={1}> {dialogTitle}</Box>

                        <IconButton onClick={onClose}>
                            <CloseIcon />
                        </IconButton>
                    </Box>
                </DialogTitle>
                <DialogContent>
                    <DialogContentText>
                        <Typography variant="subtitle2" color="initial">
                            {dialogContent}
                        </Typography>
                    </DialogContentText>
                    {children}
                </DialogContent>
                <DialogActions>
                    <Grid container spacing={1} style={{ margin: "10px" }}>
                        <Grid container item xs={4} justifyContent="flex-start">
                            {showBackButton && (
                                <Button
                                    sx={{ fontWeight: 600 }}
                                    {...buttonProps}
                                >
                                    {t("CONSTANT.BackButtonLable")}
                                </Button>
                            )}
                        </Grid>

                        <Grid container item xs={8} justifyContent="flex-end">
                            <Button
                                {...buttonProps}
                                sx={{ fontWeight: 600 }}
                                style={{ marginRight: "0.5em" }}
                            >
                                {t("CONSTANT.CancelButtonLable")}
                            </Button>
                            <Button
                                size="medium"
                                variant="contained"
                                sx={{ fontWeight: 600 }}
                                type="submit"
                            >
                                {submitButtonLabel}
                            </Button>
                        </Grid>
                    </Grid>
                </DialogActions>
            </Dialog>
        </div>
    );
};
export default FormDialog;
