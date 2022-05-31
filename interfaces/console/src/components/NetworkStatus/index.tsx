import { LoadingWrapper } from "..";
import { colors } from "../../theme";
import { getStatusByType, secToHoursNMints } from "../../utils";
import { Typography, Grid, Button, Stack, useMediaQuery } from "@mui/material";

const DOT = (color: string) => (
    <span style={{ color: `${color}`, fontSize: "24px", marginRight: 14 }}>
        ●
    </span>
);

const getIconByStatus = (status: string) => {
    switch (status) {
        case "BEING_CONFIGURED":
            return DOT(colors.yellow);
        case "ONLINE":
            return DOT(colors.green);
        default:
            return DOT(colors.red);
    }
};

type NetworkStatusProps = {
    loading?: boolean;
    duration?: string;
    statusType: string;
    regLoading: boolean;
    handleAddNode: Function;
    handleActivateUser: Function;
};

const NetworkStatus = ({
    loading,
    duration,
    statusType,
    regLoading,
    handleAddNode,
    handleActivateUser,
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
                                        duration &&
                                        duration !== "0" &&
                                        duration}
                                </Typography>

                                {!isSmall && duration && duration !== "0" && (
                                    <Typography
                                        ml={{ xs: "28px", md: "8px" }}
                                        variant={"h6"}
                                        color="secondary"
                                        sx={{
                                            fontWeight: { xs: 400, md: 500 },
                                        }}
                                    >
                                        {secToHoursNMints(
                                            new Date().getTime() / 1000 -
                                                parseInt(duration),
                                            " hours and "
                                        )}
                                    </Typography>
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
