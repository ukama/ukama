import { makeStyles } from "@mui/styles";
import CloseIcon from "@mui/icons-material/Close";
import colors from "../../../theme/colors";
import { StackedAreaChart } from "../..";
import {
    Box,
    Button,
    Dialog,
    IconButton,
    Typography,
    Stack,
    DialogActions,
    Divider,
    DialogContent,
} from "@mui/material";
import { NodeAppDetailsTypes } from "../../../types";
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
    nodeData?: NodeAppDetailsTypes;
};

const NodeAppDetailsDialog = ({
    isOpen,
    nodeData = {
        id: 1,
        cpu: 23,
        memory: 34,
        nodeAppName: "App 1",
    },
    handleClose,
    closeBtnLabel,
    isClosable = true,
}: BasicDialogProps) => {
    const classes = useStyles();
    const { id, nodeAppName, cpu, memory } = nodeData;

    return (
        <Dialog
            key={id}
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
                        <Typography variant="h6">{nodeAppName}</Typography>
                        <Typography
                            variant="h6"
                            sx={{ color: colors.darkBlue }}
                        >
                            - {"Version - 0.1"}
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
                        <Typography variant="body1">CPU:{cpu} %</Typography>
                        <Typography variant="body1">
                            Memory:{memory} KB
                        </Typography>
                    </Stack>

                    <Stack spacing={6} pt={2}>
                        <StackedAreaChart
                            hasData={true}
                            height={140}
                            title={"CPU"}
                        />
                        <StackedAreaChart
                            hasData={true}
                            height={140}
                            title={"Memory"}
                        />
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
