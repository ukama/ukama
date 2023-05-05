import { useState } from "react";
import { TooltipsText } from "../../constants";
import { Paper, Grid, Stack } from "@mui/material";
import { NodeStatsContainer, NodeStatItem, ApexLineChart } from "..";

const PLACEHOLDER_VALUE = "NA";
interface INodeRadioTab {
    metrics: any;
    loading: boolean;
}
const NodeRadioTab = ({ loading, metrics }: INodeRadioTab) => {
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
                        value={PLACEHOLDER_VALUE}
                        variant={"large"}
                        name={"TX Power"}
                        nameInfo={TooltipsText.TXPOWER}
                    />
                    <NodeStatItem
                        value={PLACEHOLDER_VALUE}
                        variant={"large"}
                        name={"RX Power"}
                        nameInfo={TooltipsText.RXPOWER}
                    />
                    <NodeStatItem
                        value={PLACEHOLDER_VALUE}
                        name={"PA Power"}
                        variant={"large"}
                        nameInfo={TooltipsText.PAPOWER}
                    />
                </NodeStatsContainer>
            </Grid>
            <Grid item lg={isCollapse ? 11 : 8} md xs>
                <Paper sx={{ p: 3, width: "100%" }}>
                    <Stack spacing={4}>
                        <ApexLineChart data={metrics["txpower"]} />
                        <ApexLineChart data={metrics["rxpower"]} />
                        <ApexLineChart data={metrics["papower"]} />
                    </Stack>
                </Paper>
            </Grid>
        </Grid>
    );
};

export default NodeRadioTab;
