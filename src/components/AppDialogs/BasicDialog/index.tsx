import { makeStyles } from "@mui/styles";
import CloseIcon from "@mui/icons-material/Close";
import {
    Box,
    Typography,
    Dialog,
    DialogContent,
    DialogTitle,
    IconButton,
} from "@mui/material";

const useStyles = makeStyles(() => ({
    basicDialogHeaderStyle: {
        padding: "0px",
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
    handleClose: any;
};

const BasicDialog = ({
    title,
    isOpen,
    content,
    handleClose,
}: BasicDialogProps) => {
    const classes = useStyles();
    return (
        <Dialog open={isOpen} onClose={handleClose}>
            <Box sx={{ padding: "16px 24px 44px 24px" }}>
                <DialogTitle className={classes.basicDialogHeaderStyle}>
                    <Typography variant="h6">{title}</Typography>
                    <IconButton
                        onClick={handleClose}
                        sx={{ m: "0px 0px 0px 24px" }}
                    >
                        <CloseIcon />
                    </IconButton>
                </DialogTitle>
                <DialogContent sx={{ padding: "0px" }}>
                    <Typography variant="body1">{content}</Typography>
                </DialogContent>
            </Box>
        </Dialog>
    );
};

export { BasicDialog };
