import {
    Line,
    YAxis,
    XAxis,
    Tooltip,
    LineChart,
    ResponsiveContainer,
    Area,
    AreaChart,
} from "recharts";
import { format } from "date-fns";
import { RoundedCard } from "../../styles";
import { Stack, Typography } from "@mui/material";
import { GraphTitleWrapper, LoadingWrapper } from "..";
import { UsersAttachedMetricsDto, ThroughputMetricsDto } from "../../generated";

interface INodeMetaDataTab {
    loading: boolean;
    throughputMetrics: ThroughputMetricsDto[];
    usersAttachedMetrics: UsersAttachedMetricsDto[];
}

const NodeMetaDataTab = ({
    loading,
    throughputMetrics,
    usersAttachedMetrics,
}: INodeMetaDataTab) => {
    const gradientOffset = () => {
        const dataMax = Math.max(...throughputMetrics.map(i => i.amount));
        const dataMin = Math.min(...throughputMetrics.map(i => i.amount));

        if (dataMax <= 0) {
            return 0;
        }
        if (dataMin >= 0) {
            return 1;
        }

        return dataMax / (dataMax - dataMin);
    };

    const off = gradientOffset();
    return (
        <LoadingWrapper radius={"small"} height={450} isLoading={loading}>
            <RoundedCard
                sx={{
                    width: "100%",
                    borderRadius: "4px",
                    height: "fit-content",
                }}
            >
                <Stack spacing={4}>
                    <Typography variant="h6">Meta Data</Typography>
                    <GraphTitleWrapper title="Users Attached">
                        <ResponsiveContainer width="100%" height={300}>
                            <LineChart
                                width={500}
                                height={300}
                                data={usersAttachedMetrics}
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
                                <YAxis fontSize={"14px"} />
                                <Tooltip />
                                <Line
                                    type="monotone"
                                    dataKey="users"
                                    stroke="#8884d8"
                                    activeDot={{ r: 8 }}
                                    strokeWidth={2}
                                    animationDuration={300}
                                />
                            </LineChart>
                        </ResponsiveContainer>
                    </GraphTitleWrapper>
                    <GraphTitleWrapper title="Throughput">
                        <ResponsiveContainer width={"100%"} height={300}>
                            <AreaChart
                                width={500}
                                height={300}
                                data={throughputMetrics}
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
                                <defs>
                                    <linearGradient
                                        id="splitColor"
                                        x1="0"
                                        y1="0"
                                        x2="0"
                                        y2="1"
                                    >
                                        <stop
                                            offset={off}
                                            stopColor="green"
                                            stopOpacity={1}
                                        />
                                        <stop
                                            offset={off}
                                            stopColor="red"
                                            stopOpacity={1}
                                        />
                                    </linearGradient>
                                </defs>
                                <Area
                                    type="monotone"
                                    dataKey="amount"
                                    stroke="#000"
                                    fill="url(#splitColor)"
                                />
                            </AreaChart>
                        </ResponsiveContainer>
                    </GraphTitleWrapper>
                </Stack>
            </RoundedCard>
        </LoadingWrapper>
    );
};

export default NodeMetaDataTab;
