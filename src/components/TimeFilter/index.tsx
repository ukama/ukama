import {
    Box,
    ToggleButton,
    ToggleButtonGroup,
    Typography,
} from "@mui/material";
import { colors } from "../../theme";
import { STATS_PERIOD } from "../../constants";
import { statsPeriodItemType } from "../../types";

interface ITimeFilter {
    filter: string;
    handleFilterSelect: Function;
}

const TimeFilter = ({ filter, handleFilterSelect }: ITimeFilter) => {
    return (
        <Box>
            <ToggleButtonGroup
                exclusive
                size="small"
                color="primary"
                value={filter}
                onChange={(_, value: string) => handleFilterSelect(value)}
            >
                {STATS_PERIOD.map(({ id, label }: statsPeriodItemType) => (
                    <ToggleButton
                        fullWidth
                        key={id}
                        value={label}
                        style={{
                            height: "32px",
                            color: colors.lightBlue,
                            border: `1px solid ${colors.lightBlue}`,
                        }}
                    >
                        <Typography
                            variant="body2"
                            sx={{
                                p: "0px 2px",
                                fontWeight: 600,
                            }}
                        >
                            {label}
                        </Typography>
                    </ToggleButton>
                ))}
            </ToggleButtonGroup>
        </Box>
    );
};
export default TimeFilter;
