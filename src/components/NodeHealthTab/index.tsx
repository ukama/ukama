import {
    Bar,
    Cell,
    YAxis,
    XAxis,
    Tooltip,
    BarChart,
    ReferenceLine,
    ResponsiveContainer,
} from "recharts";
import { format } from "date-fns";
import { LoadingWrapper } from "..";
import { RoundedCard } from "../../styles";
import { Typography } from "@mui/material";
import { TemperatureMetricsDto } from "../../generated";

interface INodeHealthTab {
    loading: boolean;
    temperatureMetrics: TemperatureMetricsDto[];
}

const NodeHealthTab = ({ loading, temperatureMetrics }: INodeHealthTab) => {
    return (
        <LoadingWrapper radius={"small"} height={450} isLoading={loading}>
            <RoundedCard
                sx={{
                    borderRadius: "4px",
                    height: "fit-content",
                }}
            >
                <Typography variant="h6" mb={4}>
                    Physical Health
                </Typography>
                <ResponsiveContainer width="100%" height={300}>
                    <BarChart
                        width={500}
                        height={300}
                        data={temperatureMetrics}
                        margin={{
                            top: 5,
                            right: 30,
                            left: 20,
                            bottom: 5,
                        }}
                    >
                        <XAxis
                            dataKey="timestamp"
                            fontSize={"14px"}
                            tickFormatter={(value: any) =>
                                format(value, "MMM dd HH:mm:ss")
                            }
                        />
                        <YAxis />
                        <Tooltip />
                        <ReferenceLine y={0} stroke="#000" />
                        <Bar dataKey="temperature">
                            {temperatureMetrics.map((entry, index) => (
                                <Cell
                                    key={`cell-${index}`}
                                    fill={
                                        entry.temperature < 0
                                            ? "#E30000"
                                            : "#82ca9d"
                                    }
                                    strokeWidth={index === 2 ? 4 : 1}
                                />
                            ))}
                        </Bar>
                    </BarChart>
                </ResponsiveContainer>
            </RoundedCard>
        </LoadingWrapper>
    );
};

export default NodeHealthTab;
