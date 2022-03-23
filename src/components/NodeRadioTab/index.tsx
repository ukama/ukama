import { useState } from "react";
import { TooltipsText } from "../../constants";
import { Paper, Grid, Stack } from "@mui/material";
import { NodeStatsContainer, NodeStatItem, ApexStackAreaChart } from "..";
interface INodeRadioTab {
    loading: boolean;
    txPowerMetrics: any;
    rxPowerMetrics: any;
    paPowerMetrics: any;
}
const NodeRadioTab = ({
    loading,
    rxPowerMetrics,
    txPowerMetrics,
    paPowerMetrics,
}: INodeRadioTab) => {
    const [isCollapse, setIsCollapse] = useState<boolean>(false);
    const handleCollapse = () => setIsCollapse(prev => !prev);
    return (
        <Grid container spacing={3}>
            <Grid item lg={!isCollapse ? 4 : 1} md xs>
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
            <Grid item lg={isCollapse ? 11 : 8} md xs>
                <Paper sx={{ p: 3, width: "100%" }}>
                    <Stack spacing={4}>
                        <ApexStackAreaChart
                            hasData={true}
                            name={"PX Power"}
                            data={txPowerMetrics}
                        />
                        <ApexStackAreaChart
                            hasData={true}
                            name={"PA POWER"}
                            data={paPowerMetrics}
                        />
                        <ApexStackAreaChart
                            hasData={true}
                            name={"RX POWER"}
                            data={rxPowerMetrics}
                        />
                    </Stack>
                </Paper>
            </Grid>
        </Grid>
    );
};

export default NodeRadioTab;
