import { Grid, Select, MenuItem, ButtonGroup, Button } from "@mui/material";
import { RoundedCard } from "../../styles";
import { colors } from "../../theme";
import { statsItemType, statsPeriodItemType } from "../../types";
import { RechartsData } from "../../constants/rechartsData";
type StatsCardProps = {
    options: statsItemType[];
    periodOptions: statsPeriodItemType[];
    handleSelect: any;
    selectOption: string;
};
import {
    ComposedChart,
    Line,
    Bar,
    XAxis,
    YAxis,
    CartesianGrid,
    Tooltip,
    Legend,
} from "recharts";
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
                            value={selectOption}
                            variant="standard"
                            disableUnderline
                            sx={{ width: "64px", color: colors.empress }}
                            onChange={handleSelect}
                        >
                            {options.map(
                                ({ id, label, value }: statsItemType) => (
                                    <MenuItem key={id} value={value}>
                                        {label}
                                    </MenuItem>
                                )
                            )}
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
                    <ComposedChart
                        width={630}
                        height={200}
                        data={RechartsData}
                        margin={{
                            top: 20,
                            right: 20,
                            bottom: 20,
                        }}
                    >
                        <CartesianGrid stroke="#f5f5f5" />
                        <XAxis dataKey="name" scale="band" />
                        <YAxis />
                        <Tooltip />
                        <Legend />
                        <Bar dataKey="uv" barSize={20} fill="#413ea0" />
                        <Line type="monotone" dataKey="uv" stroke="#ff7300" />
                    </ComposedChart>
                </Grid>
            </RoundedCard>
        </>
    );
};
export default StatsCard;
