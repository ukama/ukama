import {
    Bar,
    Cell,
    YAxis,
    XAxis,
    Tooltip,
    BarChart,
    ReferenceLine,
    ResponsiveContainer,
    Area,
    AreaChart,
    Line,
    LineChart,
} from "recharts";
import { format } from "date-fns";
import { LoadingWrapper } from "..";
import { RoundedCard } from "../../styles";
import { Stack, Typography } from "@mui/material";
import {
    IoMetricsDto,
    CpuUsageMetricsDto,
    TemperatureMetricsDto,
    MemoryUsageMetricsDto,
} from "../../generated";

interface INodeHealthTab {
    loading: boolean;
    ioMetrics: IoMetricsDto[];
    cpuUsageMetrics: CpuUsageMetricsDto[];
    memoryUsageMetrics: MemoryUsageMetricsDto[];
    temperatureMetrics: TemperatureMetricsDto[];
}

const NodeHealthTab = ({
    loading,
    ioMetrics,
    cpuUsageMetrics,
    temperatureMetrics,
    memoryUsageMetrics,
}: INodeHealthTab) => {
    return (
        <LoadingWrapper radius={"small"} height={450} isLoading={loading}>
            <RoundedCard
                sx={{
                    borderRadius: "4px",
                    height: "fit-content",
                }}
            >
                <Stack spacing={6}>
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
                                {temperatureMetrics.map(
                                    (
                                        {
                                            temperature,
                                            id,
                                        }: TemperatureMetricsDto,
                                        index: number
                                    ) => (
                                        <Cell
                                            key={`cell-${id}`}
                                            fill={
                                                temperature < 0
                                                    ? "#E30000"
                                                    : "#82ca9d"
                                            }
                                            strokeWidth={index === 2 ? 4 : 1}
                                        />
                                    )
                                )}
                            </Bar>
                        </BarChart>
                    </ResponsiveContainer>
                    <ResponsiveContainer width="100%" height={300}>
                        <BarChart
                            width={500}
                            height={300}
                            data={cpuUsageMetrics}
                            stackOffset="expand"
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
                            <YAxis
                                tickFormatter={(value: any) => `${value}%`}
                            />

                            <Tooltip />
                            <Bar dataKey="usage" fill="#ffc658" />
                        </BarChart>
                    </ResponsiveContainer>
                    <ResponsiveContainer width="100%" height={300}>
                        <LineChart
                            width={500}
                            height={300}
                            data={memoryUsageMetrics}
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
                            <YAxis
                                tickFormatter={(value: any) => `${value}%`}
                            />
                            <Tooltip />
                            <Line
                                type="monotone"
                                dataKey="usage"
                                stroke="#009d5f"
                                strokeWidth={2}
                            />
                        </LineChart>
                    </ResponsiveContainer>
                    <ResponsiveContainer width="100%" height={300}>
                        <AreaChart
                            width={500}
                            height={300}
                            data={ioMetrics}
                            margin={{
                                top: 5,
                                right: 30,
                                left: 20,
                                bottom: 5,
                            }}
                        >
                            <YAxis />
                            <Tooltip />
                            <Area
                                type="monotone"
                                dataKey="input"
                                stackId="1"
                                stroke="#FFBB28"
                                fill="#FFBB28"
                            />
                            <Area
                                type="monotone"
                                dataKey="output"
                                stackId="1"
                                stroke="#FF8042"
                                fill="#FF8042"
                            />
                        </AreaChart>
                    </ResponsiveContainer>
                </Stack>
            </RoundedCard>
        </LoadingWrapper>
    );
};

export default NodeHealthTab;
