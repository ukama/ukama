import {
    Grid,
    Select,
    MenuItem,
    Typography,
    ToggleButton,
    ToggleButtonGroup,
} from "@mui/material";
import { colors } from "../../theme";
import { RechartsData } from "../../constants/stubData";
import { RoundedCard, SkeletonRoundedCard } from "../../styles";
import { StatsItemType, statsPeriodItemType } from "../../types";
import {
    ComposedChart,
    Bar,
    XAxis,
    YAxis,
    CartesianGrid,
    Tooltip,
    ResponsiveContainer,
    Line,
} from "recharts";
type StatsCardProps = {
    loading: boolean;
    handleSelect: any;
    selectOption: number;
    selectedButton: string;
    options: StatsItemType[];
    handleSelectedButton: any;
    periodOptions: statsPeriodItemType[];
};

const StatsCard = ({
    options,
    loading,
    selectOption,
    handleSelect,
    periodOptions,
    selectedButton,
    handleSelectedButton,
}: StatsCardProps) => {
    return (
        <>
            {loading ? (
                <SkeletonRoundedCard variant="rectangular" height={337} />
            ) : (
                <RoundedCard>
                    <Grid container spacing={1}>
                        <Grid item container xs={12}>
                            <Grid item xs={12} sm={6}>
                                <Select
                                    sx={{
                                        width: "auto",
                                    }}
                                    value={selectOption}
                                    variant="standard"
                                    disableUnderline
                                    onChange={handleSelect}
                                >
                                    {options.map(
                                        ({ id, label }: StatsItemType) => (
                                            <MenuItem key={id} value={id}>
                                                <Typography variant="h6">
                                                    {label}
                                                </Typography>
                                            </MenuItem>
                                        )
                                    )}
                                </Select>
                            </Grid>
                            <Grid
                                item
                                xs={12}
                                sm={6}
                                display="flex"
                                justifyContent={"flex-end"}
                            >
                                <ToggleButtonGroup
                                    size="small"
                                    color="primary"
                                    exclusive
                                    value={selectedButton}
                                    onChange={handleSelectedButton}
                                >
                                    {periodOptions.map(
                                        ({
                                            id,
                                            label,
                                        }: statsPeriodItemType) => (
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
                                        )
                                    )}
                                </ToggleButtonGroup>
                            </Grid>
                        </Grid>
                        <Grid item xs={12}>
                            <ResponsiveContainer width="100%" height={300}>
                                <ComposedChart data={RechartsData}>
                                    <CartesianGrid stroke="#f5f5f5" />
                                    <XAxis dataKey="name" scale="band" />
                                    <YAxis />
                                    <Tooltip />
                                    <Bar
                                        dataKey="uv"
                                        barSize={20}
                                        fill="#413ea0"
                                    />
                                    <Line
                                        type="monotone"
                                        dataKey="uv"
                                        stroke="#ff7300"
                                    />
                                </ComposedChart>
                            </ResponsiveContainer>
                        </Grid>
                    </Grid>
                </RoundedCard>
            )}
        </>
    );
};
export default StatsCard;
