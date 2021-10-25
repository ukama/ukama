import * as React from "react";
import Button from "@mui/material/Button";
import TextField from "@mui/material/TextField";
import Dialog from "@mui/material/Dialog";
import DialogActions from "@mui/material/DialogActions";
import DialogContent from "@mui/material/DialogContent";
import DialogContentText from "@mui/material/DialogContentText";
import DialogTitle from "@mui/material/DialogTitle";
import { globalUseStyles } from "../../styles";
import Grid from "@mui/material/Grid";
import { useTranslation } from "react-i18next";
import "../../i18n/i18n";
type FormDialogProps = {
    dialogTitle?: string;
    dialogContent?: string;
    showBackButton?: boolean;
    submitButtonLabel?: string;
    open: boolean;
    onClose: () => void;
    formField: React.ReactElement;
};
const FormDialog = ({
    dialogTitle,
    dialogContent,
    formField,
    showBackButton,
    submitButtonLabel,
    open,
    onClose,
}: FormDialogProps) => {
    const { t } = useTranslation();
    return (
        <div>
            <Dialog open={open} onClose={onClose}>
                <DialogTitle>{dialogTitle}</DialogTitle>
                <DialogContent>
                    <DialogContentText>{dialogContent}</DialogContentText>
                    {formField}
                </DialogContent>
                <DialogActions>
                    <Grid container spacing={1} style={{ margin: "10px" }}>
                        <Grid container item xs={4} justifyContent="flex-start">
                            {showBackButton ? (
                                <Button
                                    size="large"
                                    type="submit"
                                    color="primary"
                                    variant="contained"
                                    sx={{ fontWeight: 600 }}
                                >
                                    {t("CONSTANT.BackButtonLable")}
                                </Button>
                            ) : null}
                        </Grid>

                        <Grid container item xs={8} justifyContent="flex-end">
                            <Button
                                size="large"
                                type="submit"
                                color="primary"
                                variant="contained"
                                sx={{ fontWeight: 600 }}
                                style={{ marginRight: "0.5em" }}
                            >
                                {t("CONSTANT.CancelButtonLable")}
                            </Button>
                            <Button
                                size="large"
                                variant="outlined"
                                sx={{ fontWeight: 600 }}
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
