import {
    Box,
    ToggleButton,
    ToggleButtonGroup,
    Typography,
} from "@mui/material";
import { colors } from "../../theme";
import { statsPeriodItemType } from "../../types";

interface ITimeFilter {
    filter?: string;
    options?: statsPeriodItemType[];
    handleFilterSelect: Function;
}

const TimeFilter = ({
    filter = "LIVE",
    handleFilterSelect,
    options = [{ id: "1", label: "LIVE" }],
}: ITimeFilter) => {
    return (
        <Box component="div">
            <ToggleButtonGroup
                exclusive
                size="small"
                color="primary"
                value={filter}
                onChange={(_, value: string) => handleFilterSelect(value)}
            >
                {options.map(({ id, label }: statsPeriodItemType) => (
                    <ToggleButton
                        fullWidth
                        key={id}
                        value={label}
                        style={{
                            height: "32px",
                            color: colors.hoverColor,
                            border: `1px solid ${colors.hoverColor}`,
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
