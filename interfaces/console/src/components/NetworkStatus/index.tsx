import { LoadingWrapper } from "..";
import { colors } from "../../theme";
import { getStatusByType } from "../../utils";
import { Network_Status } from "../../generated";
import {
    Typography,
    Grid,
    Button,
    Stack,
    Box,
    useMediaQuery,
} from "@mui/material";
import SignalCellularOffIcon from "@mui/icons-material/SignalCellularOff";
import CheckCircleIcon from "@mui/icons-material/CheckCircle";
import { ReactNode } from "react";
const DOT = (icon: ReactNode) => <Box>{icon}</Box>;

const getIconByStatus = (status: string) => {
    switch (status) {
        case "DOWN":
            return DOT(<CheckCircleIcon sx={{ color: colors.green }} />);
        case "ONLINE":
            return DOT(<CheckCircleIcon sx={{ color: colors.green }} />);
        default:
            return DOT(<SignalCellularOffIcon />);
    }
};

type NetworkStatusProps = {
    loading?: boolean;
    regLoading: boolean;
    handleAddNode: Function;
    handleActivateUser: Function;
    totalNodes: number | undefined;
    liveNodes: number | undefined;
    statusType: Network_Status | undefined;
};

const NetworkStatus = ({
    loading,
    regLoading,
    handleAddNode,
    handleActivateUser,
    totalNodes = undefined,
    liveNodes = undefined,
    statusType = Network_Status.Undefined,
}: NetworkStatusProps) => {
    const isSmall = useMediaQuery("(max-width:600px)");

    return (
        <Grid container spacing={2}>
            <Grid item xs={12} md={8}>
                <LoadingWrapper height={30} isLoading={loading}>
                    <Grid container>
                        <Grid item>{getIconByStatus(statusType)}</Grid>
                        <Grid item xs={11}>
                            <Stack
                                direction={{ xs: "column", md: "row" }}
                                alignItems="flex-start"
                            >
                                <Typography
                                    variant={"h6"}
                                    sx={{ fontWeight: { xs: 400, md: 500 } }}
                                >
                                    {getStatusByType(statusType)}
                                    {isSmall &&
                                        totalNodes &&
                                        liveNodes &&
                                        ` ${liveNodes}/${totalNodes} nodes up.`}
                                </Typography>

                                {!isSmall && totalNodes && liveNodes && (
                                    <>
                                        <Typography
                                            variant={"h6"}
                                            color="secondary"
                                            sx={{
                                                fontWeight: {
                                                    xs: 400,
                                                    md: 500,
                                                },
                                                whiteSpace: "break-spaces",
                                                ml: { xs: "28px", md: "8px" },
                                            }}
                                        >
                                            {`${liveNodes}/${totalNodes} `}
                                        </Typography>
                                        <Typography
                                            variant={"h6"}
                                            sx={{
                                                fontWeight: {
                                                    xs: 400,
                                                    md: 500,
                                                },
                                            }}
                                        >
                                            nodes up
                                        </Typography>
                                    </>
                                )}
                            </Stack>
                        </Grid>
                    </Grid>
                </LoadingWrapper>
            </Grid>
            <Grid item xs={12} md={4} display="flex" justifyContent="flex-end">
                <LoadingWrapper height={30} isLoading={loading}>
                    <Grid container spacing={2} justifyContent="flex-end">
                        <Grid item xs={5} md={5} lg={4}>
                            <Button
                                fullWidth
                                variant="contained"
                                onClick={() => handleActivateUser()}
                            >
                                ADD USER
                            </Button>
                        </Grid>
                        <Grid item xs={7} md={7} lg={6} xl={5}>
                            <LoadingWrapper isLoading={regLoading} height={40}>
                                <Button
                                    fullWidth
                                    variant="contained"
                                    onClick={() => handleAddNode()}
                                >
                                    REGISTER NODE
                                </Button>
                            </LoadingWrapper>
                        </Grid>
                    </Grid>
                </LoadingWrapper>
            </Grid>
        </Grid>
    );
};

export default NetworkStatus;
