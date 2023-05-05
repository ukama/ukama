import { useState } from "react";
import { TooltipsText } from "../../constants";
import { Paper, Grid, Stack } from "@mui/material";
import { NodeStatsContainer, NodeStatItem, ApexLineChart } from "..";

const PLACEHOLDER_VALUE = "NA";
interface INodeOverviewTab {
    metrics: any;
    loading: boolean;
}
const NodeNetworkTab = ({ loading, metrics }: INodeOverviewTab) => {
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
                        value={PLACEHOLDER_VALUE}
                        name={"Throughput (D/L)"}
                        nameInfo={TooltipsText.DL}
                    />
                    <NodeStatItem
                        variant={"large"}
                        value={PLACEHOLDER_VALUE}
                        name={"Throughput (U/L)"}
                        nameInfo={TooltipsText.UL}
                    />
                    <NodeStatItem
                        value={PLACEHOLDER_VALUE}
                        variant={"large"}
                        name={"RRC CNX Success"}
                        nameInfo={TooltipsText.RRCCNX}
                    />
                    <NodeStatItem
                        value={PLACEHOLDER_VALUE}
                        variant={"large"}
                        name={"ERAB Drop Rate"}
                        nameInfo={TooltipsText.ERAB}
                    />
                    <NodeStatItem
                        value={PLACEHOLDER_VALUE}
                        variant={"large"}
                        name={"RLS  Drop Rate"}
                        nameInfo={TooltipsText.RLS}
                    />
                </NodeStatsContainer>
            </Grid>
            <Grid item lg={isCollapse ? 11 : 8} md xs>
                <Paper sx={{ p: 3, width: "100%" }}>
                    <Stack spacing={4}>
                        <ApexLineChart data={metrics["throughputuplink"]} />
                        <ApexLineChart data={metrics["throughputdownlink"]} />
                        <ApexLineChart data={metrics["rrc"]} />
                        <ApexLineChart data={metrics["erab"]} />
                        <ApexLineChart data={metrics["rlc"]} />
                    </Stack>
                </Paper>
            </Grid>
        </Grid>
    );
};

export default NodeNetworkTab;
