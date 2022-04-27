import { makeStyles } from "@mui/styles";
import CloseIcon from "@mui/icons-material/Close";
import colors from "../../../theme/colors";
import {
    Box,
    Button,
    Dialog,
    IconButton,
    Typography,
    Stack,
    DialogActions,
    DialogContent,
} from "@mui/material";
import ApexLineChart from "../../ApexLineChart";
const useStyles = makeStyles(() => ({
    basicDialogHeaderStyle: {
        padding: "0px 0px 10px 0px",
        display: "flex",
        flexDirection: "row",
        alignItems: "center",
        justifyContent: "space-between",
    },
}));

type BasicDialogProps = {
    isOpen: boolean;
    isClosable?: boolean;
    handleClose: Function;
    closeBtnLabel?: string;
    nodeData?: any;
};

const NodeAppDetailsDialog = ({
    isOpen,
    nodeData,
    handleClose,
    closeBtnLabel,
    isClosable = true,
}: BasicDialogProps) => {
    const classes = useStyles();
    return (
        <Dialog
            open={isOpen}
            maxWidth="xl"
            onBackdropClick={() => isClosable && handleClose()}
        >
            <Box
                component="div"
                sx={{
                    width: { xs: "100%", md: "500px" },
                    padding: "16px 24px",
                }}
            >
                <Box component="div" className={classes.basicDialogHeaderStyle}>
                    <Stack direction="row" spacing={1}>
                        <Typography variant="h6">{nodeData?.title}</Typography>
                        <Typography
                            variant="h6"
                            sx={{ color: colors.darkBlue }}
                        >
                            - {nodeData?.version}
                        </Typography>
                    </Stack>

                    {isClosable && (
                        <IconButton
                            onClick={() => handleClose()}
                            sx={{
                                ml: "25px",
                                p: "8px",
                                position: "relative",
                                left: 10,
                            }}
                        >
                            <CloseIcon />
                        </IconButton>
                    )}
                </Box>
                <DialogContent sx={{ padding: 0, mb: 4 }}>
                    <Stack direction="column" sx={{ mb: 4 }}>
                        <Typography variant="body1">
                            CPU:{nodeData?.cpu} %
                        </Typography>
                        <Typography variant="body1">
                            MEMORY:{nodeData?.memory} KB
                        </Typography>
                    </Stack>

                    <Stack spacing={6} pt={2}>
                        <ApexLineChart data={{ name: "CPU", data: [] }} />
                        <ApexLineChart data={{ name: "MEMORY", data: [] }} />
                    </Stack>
                </DialogContent>
                <DialogActions sx={{ padding: 0 }}>
                    <Button
                        onClick={() => handleClose()}
                        sx={{
                            mr: 2,
                            color: colors.primaryMain,
                            position: "relative",
                            left: 10,
                        }}
                    >
                        {closeBtnLabel}
                    </Button>
                </DialogActions>
            </Box>
        </Dialog>
    );
};

export default NodeAppDetailsDialog;
