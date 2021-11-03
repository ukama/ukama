import { Grid, Select, MenuItem, ButtonGroup, Button } from "@mui/material";
import { RoundedCard } from "../../styles";
import { colors } from "../../theme";
import { StatsItemType, statsPeriodItemType } from "../../types";
import { RechartsData } from "../../constants/stubData";
import {
    ComposedChart,
    Bar,
    XAxis,
    YAxis,
    CartesianGrid,
    Tooltip,
    ResponsiveContainer,
    Legend,
} from "recharts";
type StatsCardProps = {
    options: StatsItemType[];
    periodOptions: statsPeriodItemType[];
    handleSelect: any;
    selectOption: number;
};

const StatsCard = ({
    handleSelect,
    options,
    periodOptions,
    selectOption,
}: StatsCardProps) => {
    return (
        <>
            <RoundedCard>
                <Grid container spacing={1} justifyContent="space-between">
                    <Grid item xs={6} container>
                        <Select
                            style={{
                                minWidth: 120,
                                color: colors.black,
                            }}
                            value={selectOption}
                            variant="standard"
                            disableUnderline
                            sx={{ width: "64px", color: colors.empress }}
                            onChange={handleSelect}
                        >
                            {options.map(({ id, label }: StatsItemType) => (
                                <MenuItem key={id} value={id}>
                                    {label}
                                </MenuItem>
                            ))}
                        </Select>
                    </Grid>
                    <Grid item>
                        <ButtonGroup size="small" variant="outlined">
                            {periodOptions.map(
                                ({ id, label }: statsPeriodItemType) => (
                                    <Button key={id}>{label}</Button>
                                )
                            )}
                        </ButtonGroup>
                    </Grid>
                </Grid>
                <Grid container spacing={1}>
                    <ResponsiveContainer width="100%" height={300}>
                        <ComposedChart
                            data={RechartsData}
                            margin={{
                                top: 20,
                                right: 5,
                                bottom: 20,
                            }}
                        >
                            <CartesianGrid stroke="#f5f5f5" />
                            <XAxis dataKey="name" scale="band" />
                            <YAxis />
                            <Tooltip />
                            <Legend />
                            <Bar dataKey="uv" barSize={20} fill="#413ea0" />
                        </ComposedChart>
                    </ResponsiveContainer>
                </Grid>
            </RoundedCard>
        </>
    );
};
export default StatsCard;
