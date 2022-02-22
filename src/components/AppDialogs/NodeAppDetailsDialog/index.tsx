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
        padding: "0px 0px 18px 0px",
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
            hideBackdrop
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
                    <Typography variant="h6">
                        {nodeAppName} - {"Version - 0.1"}
                    </Typography>
                    {isClosable && (
                        <IconButton
                            onClick={() => handleClose()}
                            sx={{ ml: "24px", p: "8px" }}
                        >
                            <CloseIcon />
                        </IconButton>
                    )}
                </Box>
                <DialogContent sx={{ padding: 0, mb: 4 }}>
                    <Typography variant="body1">App Information</Typography>
                    <Divider />
                    <Stack direction="column" sx={{ mb: 4 }}>
                        <Typography
                            variant="body1"
                            sx={{ color: colors.black70 }}
                        >
                            CPU:{cpu} %
                        </Typography>
                        <Typography
                            variant="body1"
                            sx={{ color: colors.black70 }}
                        >
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
                            justifyItems: "center",
                            color: colors.primaryMain,
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
