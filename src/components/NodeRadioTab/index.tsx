import { useState } from "react";
import { NodeRadioTabConfigure, TooltipsText } from "../../constants";
import { Paper, Grid, Stack } from "@mui/material";
import ApexLineChartIntegration from "../ApexLineChart";
import { NodeStatsContainer, NodeStatItem, ApexStackAreaChart } from "..";
import { NodeDto } from "../../generated";
interface INodeRadioTab {
    loading: boolean;
    txPowerMetrics: any;
    rxPowerMetrics: any;
    paPowerMetrics: any;
    selectedNode: NodeDto | undefined;
}
const NodeRadioTab = ({
    loading,
    rxPowerMetrics,
    txPowerMetrics,
    selectedNode,
    paPowerMetrics,
}: INodeRadioTab) => {
    const [isCollapse, setIsCollapse] = useState<boolean>(false);
    const handleCollapse = () => setIsCollapse(prev => !prev);
    return (
        <Grid container spacing={3}>
            <Grid item lg={!isCollapse ? 3 : 1} md xs>
                <NodeStatsContainer
                    index={0}
                    selected={0}
                    title={"Radio"}
                    loading={loading}
                    isCollapsable={true}
                    isCollapse={isCollapse}
                    onCollapse={handleCollapse}
                >
                    <NodeStatItem
                        value={"NNN"}
                        variant={"large"}
                        name={"TX Power"}
                        nameInfo={TooltipsText.TXPOWER}
                    />
                    <NodeStatItem
                        value={"NNN"}
                        variant={"large"}
                        name={"RX Power"}
                        nameInfo={TooltipsText.RXPOWER}
                    />
                    <NodeStatItem
                        value={"NNN"}
                        name={"PA Power"}
                        variant={"large"}
                        nameInfo={TooltipsText.PAPOWER}
                    />
                </NodeStatsContainer>
            </Grid>
            <Grid item lg={isCollapse ? 11 : 9} md xs>
                <Paper sx={{ padding: "22px 18px 0px 30px", width: "100%" }}>
                    <Stack spacing={1}>
                        {/* <ApexLineChartIntegration
                            hasData={true}
                            data={txPowerMetrics}
                            name={"TX Power"}
                        />
                        <ApexLineChartIntegration
                            hasData={true}
                            data={rxPowerMetrics}
                            name={"RX Power"}
                        /> */}
                        {/* <ApexLineChartIntegration
                            hasData={true}
                            data={paPowerMetrics}
                            name={"PA Power"}
                        /> */}
                        {NodeRadioTabConfigure[
                            (selectedNode?.type as string) || ""
                        ][1].show && (
                            <ApexStackAreaChart
                                hasData={true}
                                name={
                                    NodeRadioTabConfigure[
                                        (selectedNode?.type as string) || ""
                                    ][1].name
                                }
                                data={paPowerMetrics}
                            />
                        )}
                        {/* <ApexStackAreaChart
                            hasData={true}
                            name={paPowerMetrics}
                            data={paPowerMetrics}
                        /> */}
                        {/* <StackedAreaChart hasData={true} title={"TX Power"} />
                        <StackedAreaChart hasData={true} title={"RX Power "} />
                        <StackedAreaChart hasData={true} title={"PA Power "} /> */}
                    </Stack>
                </Paper>
            </Grid>
        </Grid>
    );
};

export default NodeRadioTab;
