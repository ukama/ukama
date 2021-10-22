import * as React from "react";
import Button from "@mui/material/Button";
import TextField from "@mui/material/TextField";
import Dialog from "@mui/material/Dialog";
import DialogActions from "@mui/material/DialogActions";
import DialogContent from "@mui/material/DialogContent";
import DialogContentText from "@mui/material/DialogContentText";
import DialogTitle from "@mui/material/DialogTitle";
import { globalUseStyles } from "../../styles";
type FieldsProps = {
    id: string;
    name: string;
    label: string;
    value: string;
    onChange: any;
    helperText: string;
    error: boolean;
};
type FormDialogProps = {
    dialogTitle?: string;
    dialogContent?: string;
    formFields?: [];
    cancelButton?: any;
    backButton?: any;
    submitButton?: any;
};
const FormDialog = ({
    dialogTitle,
    dialogContent,
    formFields,
    cancelButton,
    backButton,
    submitButton,
}: FormDialogProps) => {
    const [open, setOpen] = React.useState(false);
    const classes = globalUseStyles();
    const handleClose = () => {
        setOpen(false);
    };

    return (
        <div>
            <Dialog open={open} onClose={handleClose}>
                <DialogTitle>{dialogTitle}</DialogTitle>
                <DialogContent>
                    <DialogContentText>{dialogContent}</DialogContentText>
                    {formFields &&
                        formFields.map((fields: FieldsProps) => {
                            return (
                                <>
                                    <TextField
                                        fullWidth
                                        id={fields.id}
                                        name={fields.name}
                                        label={fields.label}
                                        value={fields.value}
                                        onChange={fields.onChange}
                                        InputLabelProps={{ shrink: true }}
                                        InputProps={{
                                            classes: {
                                                input: classes.inputFieldStyle,
                                            },
                                        }}
                                        helperText={fields.helperText}
                                        error={fields.error}
                                    />
                                </>
                            );
                        })}
                </DialogContent>
                <DialogActions>
                    <Button
                        size="large"
                        type="submit"
                        color="primary"
                        variant="contained"
                        sx={{ fontWeight: 600 }}
                    >
                        {cancelButton}
                    </Button>
                    <Button
                        size="large"
                        type="submit"
                        color="primary"
                        variant="contained"
                        sx={{ fontWeight: 600 }}
                    >
                        {submitButton}
                    </Button>

                    <Button
                        size="large"
                        type="submit"
                        color="primary"
                        variant="contained"
                        sx={{ fontWeight: 600 }}
                    >
                        {backButton}
                    </Button>
                </DialogActions>
            </Dialog>
        </div>
    );
};
export default FormDialog;
