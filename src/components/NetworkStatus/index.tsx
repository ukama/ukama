import { colors } from "../../theme";
import { LoadingWrapper } from "..";
import { SelectItemType } from "../../types";
import { Box, Typography, Grid, Select, MenuItem } from "@mui/material";

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
    option: string;
    loading?: boolean;
    duration?: string;
    statusType: string;
    options: SelectItemType[];
    handleStatusChange: Function;
};

const NetworkStatus = ({
    option,
    options,
    loading,
    duration,
    statusType,
    handleStatusChange,
}: NetworkStatusProps) => {
    const getStatusByType = (status: string) => {
        if (status === "BEING_CONFIGURED")
            return "Your network is being configured.";
        else if (status === "ONLINE")
            return "Your network is online and well for ";
        else return "Something went wrong.";
    };

    return (
        <Grid width="100%" container p="18px 2px">
            <Grid item xs={12} md={10}>
                <LoadingWrapper height={30} width={280} isLoading={loading}>
                    <Box display="flex" flexDirection="row" alignItems="center">
                        {getIconByStatus(statusType)}
                        <Typography variant={"h6"}>
                            {getStatusByType(statusType)}
                        </Typography>
                        {duration && (
                            <Typography
                                ml="8px"
                                variant={"h6"}
                                color={colors.empress}
                            >
                                {duration}
                            </Typography>
                        )}
                    </Box>
                </LoadingWrapper>
            </Grid>
            <Grid item xs={12} md={2} display="flex" justifyContent="flex-end">
                <LoadingWrapper height={30} isLoading={loading}>
                    <Select
                        value={option}
                        disableUnderline
                        variant="standard"
                        sx={{
                            color: colors.black,
                        }}
                        onChange={e => handleStatusChange(e.target.value)}
                    >
                        {options.map(({ id, label, value }: SelectItemType) => (
                            <MenuItem key={id} value={value}>
                                <Typography variant="body1">{label}</Typography>
                            </MenuItem>
                        ))}
                    </Select>
                </LoadingWrapper>
            </Grid>
        </Grid>
    );
};

export default NetworkStatus;
