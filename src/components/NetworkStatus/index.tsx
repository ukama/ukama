import { colors } from "../../theme";
import { SelectItemType } from "../../types";
import { Box, Typography, Grid, Select, MenuItem } from "@mui/material";

const DOT = (color: string) => (
    <span style={{ color: `${color}`, fontSize: "24px", marginRight: 14 }}>
        ●
    </span>
);

const getIconByStatus = (status: string) => {
    switch (status) {
        case "DONE":
            return DOT(colors.green);
        case "IN_PROGRESS":
            return DOT(colors.yellow);
        case "FAILED":
            return DOT(colors.red);
        default:
            return DOT(colors.red);
    }
};

type NetworkStatusProps = {
    statusType: string;
    status: string;
    option: string;
    duration?: string;
    options: SelectItemType[];
    handleStatusChange: Function;
};

const NetworkStatus = ({
    status,
    option,
    options,
    duration,
    statusType,
    handleStatusChange,
}: NetworkStatusProps) => {
    return (
        <Grid width="100%" container p="18px 8px">
            <Grid item xs={12} md={10}>
                <Box display="flex" flexDirection="row" alignItems="center">
                    {getIconByStatus(statusType)}
                    <Typography variant={"h6"}>{status}</Typography>
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
            </Grid>
            <Grid item xs={12} md={2} display="flex" justifyContent="flex-end">
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
            </Grid>
        </Grid>
    );
};

export default NetworkStatus;
