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
} from "@mui/material";

const useStyles = makeStyles(() => ({
    basicDialogHeaderStyle: {
        padding: "0px 0px 18px 0px",
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
            <Box
                sx={{
                    width: { xs: "100%", md: "500px" },
                    padding: "16px 8px 8px 24px",
                }}
            >
                <Box className={classes.basicDialogHeaderStyle}>
                    <Typography variant="h6">{title}</Typography>
                    <IconButton
                        onClick={handleClose}
                        sx={{ ml: "24px", p: "8px" }}
                    >
                        <CloseIcon />
                    </IconButton>
                </Box>
                <DialogContent sx={{ p: "18px 0px" }}>
                    <Typography variant="body1">{content}</Typography>
                </DialogContent>
                <DialogActions>
                    <Button autoFocus onClick={handleClose}>
                        Close
                    </Button>
                </DialogActions>
            </Box>
        </Dialog>
    );
};

export default BasicDialog;
