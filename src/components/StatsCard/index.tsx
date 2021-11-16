import {
    Grid,
    Select,
    MenuItem,
    ToggleButton,
    ToggleButtonGroup,
    Typography,
    Button,
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
const buttons = [
    <Button key="one" fullWidth>
        One
    </Button>,
    <Button key="two" fullWidth>
        Two
    </Button>,
    <Button key="three" fullWidth>
        Three
    </Button>,
];
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
                    <Grid container spacing={2}>
                        <Grid item container spacing={2}>
                            <Grid item xs={12} sm={6}>
                                <Select
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
                                container
                                justifyContent="flex-end"
                            >
                                <ToggleButtonGroup
                                    size="small"
                                    color="primary"
                                    value={selectedButton}
                                    exclusive
                                    onChange={handleSelectedButton}
                                >
                                    {periodOptions.map(
                                        ({
                                            id,
                                            label,
                                        }: statsPeriodItemType) => (
                                            <ToggleButton
                                                key={id}
                                                value={label}
                                                style={{
                                                    border: `2.2px solid ${colors.lightBlue}`,
                                                    color: colors.lightBlue,
                                                }}
                                            >
                                                <Typography
                                                    variant="body2"
                                                    sx={{
                                                        p: "2px",
                                                        fontWeight: 900,
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
