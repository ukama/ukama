import {
    Line,
    YAxis,
    XAxis,
    Tooltip,
    LineChart,
    ResponsiveContainer,
} from "recharts";
import { format } from "date-fns";
import { LoadingWrapper } from "..";
import { RoundedCard } from "../../styles";
import { NodeRfDto } from "../../generated";
import { Typography } from "@mui/material";

interface INodeRFKpiTab {
    loading: boolean;
    metrics: NodeRfDto[];
}

const NodeRFKpiTab = ({ loading, metrics }: INodeRFKpiTab) => {
    return (
        <LoadingWrapper radius={"small"} height={450} isLoading={loading}>
            <RoundedCard
                sx={{
                    borderRadius: "4px",
                    height: "fit-content",
                }}
            >
                <Typography variant="h6">RF KPIs</Typography>
                <ResponsiveContainer width="100%" height={300}>
                    <LineChart
                        width={500}
                        height={300}
                        data={metrics}
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
                            dataKey="qam"
                            stroke="#8884d8"
                            activeDot={{ r: 8 }}
                            strokeWidth={2}
                            animationDuration={300}
                        />
                        <Line
                            type="monotone"
                            dataKey="rfOutput"
                            stroke="#82ca9d"
                            strokeWidth={2}
                            animationDuration={300}
                        />
                        <Line
                            type="monotone"
                            dataKey="rssi"
                            stroke="#E6534E"
                            strokeWidth={2}
                            animationDuration={300}
                        />
                    </LineChart>
                </ResponsiveContainer>
            </RoundedCard>
        </LoadingWrapper>
    );
};

export default NodeRFKpiTab;
