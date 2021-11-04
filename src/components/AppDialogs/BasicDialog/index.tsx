import { makeStyles } from "@mui/styles";
import CloseIcon from "@mui/icons-material/Close";
import {
    Typography,
    Dialog,
    DialogContent,
    DialogTitle,
    IconButton,
} from "@mui/material";

const useStyles = makeStyles(() => ({
    basicDialogHeaderStyle: {
        display: "flex",
        flexDirection: "row",
        alignItems: "center",
        justifyContent: "space-between",
    },
    headerTitleStyle: {
        lineHeight: "160%",
        fontSize: "1.25rem",
        marginRight: "24px",
        letterSpacing: "0.15px",
    },
    contentStyle: {
        fontSize: "16px",
        lineHeight: "19px",
        letterSpacing: "-0.02em",
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
        <Dialog onClose={handleClose} open={isOpen}>
            <DialogTitle className={classes.basicDialogHeaderStyle}>
                <Typography className={classes.headerTitleStyle}>
                    {title}
                </Typography>
                <IconButton onClick={handleClose} sx={{ padding: "0px" }}>
                    <CloseIcon />
                </IconButton>
            </DialogTitle>
            <DialogContent className={classes.contentStyle}>
                {content}
            </DialogContent>
        </Dialog>
    );
};

export { BasicDialog };
