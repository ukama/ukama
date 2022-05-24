import { styled } from "@mui/material/styles";
import CloseIcon from "@mui/icons-material/Close";
import colors from "../../../theme/colors";
import {
    Dialog,
    IconButton,
    DialogTitle,
    DialogContent,
    DialogActions,
    Button,
    FormControlLabel,
    Typography,
    Checkbox,
} from "@mui/material";
import React from "react";
const BootstrapDialog = styled(Dialog)(({ theme }) => ({
    "& .MuiDialogContent-root": {
        padding: theme.spacing(2),
    },
    "& .MuiDialogActions-root": {
        padding: theme.spacing(1),
    },
}));

export interface DialogTitleProps {
    id: string;
    children?: React.ReactNode;
    onClose: () => void;
}
type softwareUpdateModalProps = {
    isOpen: boolean;
    handleClose: any;
    title: string;
    content: string;
    submit: any;
    btnLabel?: string;
};
const BootstrapDialogTitle = (props: DialogTitleProps) => {
    const { children, onClose, ...other } = props;

    return (
        <DialogTitle sx={{ m: 0, p: 2 }} {...other}>
            {children}
            {onClose ? (
                <IconButton
                    aria-label="close"
                    onClick={onClose}
                    sx={{
                        position: "absolute",
                        right: 8,
                        top: 8,
                        color: theme => theme.palette.grey[500],
                    }}
                >
                    <CloseIcon />
                </IconButton>
            ) : null}
        </DialogTitle>
    );
};

const SoftwareUpdateModal = ({
    isOpen,
    handleClose,
    submit,
    content,
    title,
    btnLabel = "CONTINUE WITH UPDATE ALL",
}: softwareUpdateModalProps) => {
    const [checked, setChecked] = React.useState(false);

    const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setChecked(event.target.checked);
    };

    return (
        <div>
            <BootstrapDialog
                onClose={handleClose}
                aria-labelledby="customized-dialog-title"
                open={isOpen}
            >
                <BootstrapDialogTitle
                    id="customized-dialog-title"
                    onClose={handleClose}
                >
                    {title}
                </BootstrapDialogTitle>
                <DialogContent>
                    <Typography>{content}</Typography>
                </DialogContent>
                <DialogActions sx={{ marginBottom: 2 }}>
                    <FormControlLabel
                        sx={{ position: "relative", left: 10 }}
                        control={
                            <Checkbox
                                checked={checked}
                                onChange={handleChange}
                            />
                        }
                        label="Don't show again"
                    />
                    <div
                        style={{
                            flex: "1 0 0",
                        }}
                    />
                    <Button
                        onClick={handleClose}
                        sx={{
                            marginRight: 3,
                        }}
                    >
                        CANCEL
                    </Button>
                    <Button
                        variant="contained"
                        onClick={() => submit(checked)}
                        sx={{ position: "relative", right: 10 }}
                    >
                        {btnLabel}
                    </Button>
                </DialogActions>
            </BootstrapDialog>
        </div>
    );
};
export default SoftwareUpdateModal;
