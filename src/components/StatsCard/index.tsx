import { colors } from "../../theme";
import { RoundedCard } from "../../styles";
import { RechartsData } from "../../constants/stubData";
import { StatsItemType, statsPeriodItemType } from "../../types";
import { Grid, Select, MenuItem, ButtonGroup, Button } from "@mui/material";
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
                <Grid container spacing={2}>
                    <Grid item xs={12} md={6} container>
                        <Select
                            value={selectOption}
                            variant="standard"
                            disableUnderline
                            size="medium"
                            sx={{
                                minWidth: 120,
                                color: colors.balck,
                            }}
                            onChange={handleSelect}
                        >
                            {options.map(({ id, label }: StatsItemType) => (
                                <MenuItem key={id} value={id}>
                                    {label}
                                </MenuItem>
                            ))}
                        </Select>
                    </Grid>
                    <Grid
                        item
                        xs={12}
                        md={6}
                        display="flex"
                        justifyContent="flex-end"
                    >
                        <ButtonGroup size="small" variant="outlined">
                            {periodOptions.map(
                                ({ id, label }: statsPeriodItemType) => (
                                    <Button key={id}>{label}</Button>
                                )
                            )}
                        </ButtonGroup>
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
