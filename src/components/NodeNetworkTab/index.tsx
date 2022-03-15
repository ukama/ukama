import { TooltipsText } from "../../constants";
import { Paper, Grid, Typography } from "@mui/material";
import { NodeStatsContainer, NodeStatItem, LineChart } from "..";
import { useState } from "react";
import ApexLineChartIntegration from "../ApexLineChart";
interface INodeOverviewTab {
    loading: boolean;
    throughpuULMetric: any;
    throughpuDLMetric: any;
    erabDropRateMetrix: any;
    rrcCnxSuccessMetrix: any;
    rlsDropRateMetrics: any;
}
const NodeNetworkTab = ({
    loading,
    throughpuULMetric,
    rlsDropRateMetrics,
    erabDropRateMetrix,
    rrcCnxSuccessMetrix,
    throughpuDLMetric,
}: INodeOverviewTab) => {
    const [isCollapse, setIsCollapse] = useState<boolean>(false);
    const handleCollapse = () => setIsCollapse(prev => !prev);

    return (
        <Grid container spacing={3}>
            <Grid md xs item lg={!isCollapse ? 4 : 1}>
                <NodeStatsContainer
                    index={0}
                    selected={0}
                    loading={loading}
                    title={"Network"}
                    isCollapsable={true}
                    isCollapse={isCollapse}
                    onCollapse={handleCollapse}
                >
                    <NodeStatItem
                        variant={"large"}
                        value={"100 Mbps"}
                        name={"Throughput (D/L)"}
                        nameInfo={TooltipsText.DL}
                    />
                    <NodeStatItem
                        variant={"large"}
                        value={"400 Mbps"}
                        name={"Throughput (U/L)"}
                        nameInfo={TooltipsText.UL}
                    />
                    <NodeStatItem
                        value={"80%"}
                        variant={"large"}
                        name={"RRC CNX Success"}
                        nameInfo={TooltipsText.RRCCNX}
                    />
                    <NodeStatItem
                        value={"72%"}
                        variant={"large"}
                        name={"ERAB Drop Rate"}
                        nameInfo={TooltipsText.ERAB}
                    />
                    <NodeStatItem
                        value={"60%"}
                        variant={"large"}
                        name={"RLS  Drop Rate"}
                        nameInfo={TooltipsText.RLS}
                    />
                </NodeStatsContainer>
            </Grid>
            <Grid item lg={isCollapse ? 11 : 8} md xs>
                <Paper sx={{ padding: "22px 18px 0px 30px", width: "100%" }}>
                    <Typography variant="h6">Network</Typography>

                    <ApexLineChartIntegration
                        hasData={true}
                        data={throughpuULMetric}
                        name={"Throughput (U/L)"}
                    />

                    <ApexLineChartIntegration
                        hasData={true}
                        data={throughpuDLMetric}
                        name={"Throughput (D/L)"}
                    />
                    <ApexLineChartIntegration
                        hasData={true}
                        data={rrcCnxSuccessMetrix}
                        name={"RRC CNX Success"}
                    />
                    <ApexLineChartIntegration
                        hasData={true}
                        data={erabDropRateMetrix}
                        name={"ERAB Drop Rate"}
                    />
                    <ApexLineChartIntegration
                        hasData={true}
                        data={rlsDropRateMetrics}
                        name={"RLS  Drop Rate"}
                    />
                </Paper>
            </Grid>
        </Grid>
    );
};

export default NodeNetworkTab;
