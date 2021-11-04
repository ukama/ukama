import {
    Grid,
    Select,
    MenuItem,
    ToggleButton,
    ToggleButtonGroup,
    Typography,
} from "@mui/material";
import { RoundedCard } from "../../styles";
import { colors } from "../../theme";
import { RechartsData } from "../../constants/stubData";
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
    options: StatsItemType[];
    periodOptions: statsPeriodItemType[];
    handleSelect: any;
    selectOption: number;
    handleSelectedButton: any;
    selectedButton: string;
};

const StatsCard = ({
    handleSelect,
    options,
    handleSelectedButton,
    selectedButton,
    periodOptions,
    selectOption,
}: StatsCardProps) => {
    return (
        <>
            <RoundedCard>
                <Grid container spacing={2}>
                    <Grid item xs={12} md={6} container>
                        <Select
                            value={selectOption}
                            variant="standard"
                            disableUnderline
                            onChange={handleSelect}
                            style={{
                                minWidth: "30%",

                                color: colors.empress,
                            }}
                        >
                            {options.map(({ id, label }: StatsItemType) => (
                                <MenuItem key={id} value={id}>
                                    <Typography variant="h6" color="initial">
                                        {label}
                                    </Typography>
                                </MenuItem>
                            ))}
                        </Select>
                    </Grid>
                    <Grid
                        item
                        xs={12}
                        md={6}
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
                                ({ id, label }: statsPeriodItemType) => (
                                    <ToggleButton
                                        key={id}
                                        value={label}
                                        style={{
                                            border: `1px solid ${colors.lightBlue}`,
                                            color: colors.lightBlue,
                                        }}
                                    >
                                        <Typography variant="h6">
                                            {label}
                                        </Typography>
                                    </ToggleButton>
                                )
                            )}
                        </ToggleButtonGroup>
                    </Grid>
                    <Grid item xs={12}>
                        <ResponsiveContainer width="100%" height={300}>
                            <ComposedChart data={RechartsData}>
                                <CartesianGrid stroke="#f5f5f5" />
                                <XAxis dataKey="name" scale="band" />
                                <YAxis />
                                <Tooltip />
                                <Bar dataKey="uv" barSize={20} fill="#413ea0" />
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
        </>
    );
};
export default StatsCard;
