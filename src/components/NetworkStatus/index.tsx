import { colors } from "../../theme";
import { LoadingWrapper } from "..";
import { HorizontalContainer } from "../../styles";
import { Box, Typography, Grid, Button } from "@mui/material";
import { getStatusByType } from "../../utils";

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
    handleAddNode: Function;
    handleActivateUser: Function;
};

const NetworkStatus = ({
    loading,
    duration,
    statusType,
    handleAddNode,
    handleActivateUser,
}: NetworkStatusProps) => {
    return (
        <Grid width="100%" container>
            <Grid item xs={12} md={8}>
                <LoadingWrapper height={30} width={280} isLoading={loading}>
                    <Box
                        component="div"
                        display="flex"
                        flexDirection="row"
                        alignItems="center"
                    >
                        {getIconByStatus(statusType)}
                        <Typography variant={"h6"}>
                            {getStatusByType(statusType)}
                        </Typography>
                        {duration && (
                            <Typography
                                ml="8px"
                                variant={"h6"}
                                color="secondary"
                            >
                                {duration}
                            </Typography>
                        )}
                    </Box>
                </LoadingWrapper>
            </Grid>
            <Grid item xs={12} md={4} display="flex" justifyContent="flex-end">
                <LoadingWrapper height={30} isLoading={loading}>
                    <HorizontalContainer>
                        <Button
                            variant="contained"
                            sx={{ width: "144px", mr: "18px" }}
                            onClick={() => handleActivateUser()}
                        >
                            INVITE USER
                        </Button>
                        <Button
                            variant="contained"
                            sx={{ width: "164px" }}
                            onClick={() => handleAddNode()}
                        >
                            REGISTER NODE
                        </Button>
                    </HorizontalContainer>
                </LoadingWrapper>
            </Grid>
        </Grid>
    );
};

export default NetworkStatus;
