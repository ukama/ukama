import { TooltipsText } from "../../constants";
import { Stack, Paper, Grid, Typography } from "@mui/material";
import { NodeStatsContainer, NodeStatItem, StackedAreaChart } from "..";
import { useState } from "react";
interface INodeResourcesTab {
    loading: boolean;
}
const NodeResourcesTab = ({ loading }: INodeResourcesTab) => {
    const [isCollapse, setIsCollapse] = useState<boolean>(false);
    const handleCollapse = () => setIsCollapse(prev => !prev);
    return (
        <Grid container spacing={3}>
            <Grid item lg={!isCollapse ? 3 : 1} md xs>
                <NodeStatsContainer
                    index={0}
                    selected={0}
                    loading={loading}
                    title={"Resources"}
                    isCollapsable={true}
                    isCollapse={isCollapse}
                    onCollapse={handleCollapse}
                >
                    <NodeStatItem
                        value={"NNN"}
                        variant={"large"}
                        name={"Memory-TRX"}
                        nameInfo={TooltipsText.MTRX}
                    />
                    <NodeStatItem
                        value={"NNN"}
                        variant={"large"}
                        name={"Memory-COM"}
                        nameInfo={TooltipsText.MCOM}
                    />
                    <NodeStatItem
                        value={"NNN"}
                        name={"CPU-TRX"}
                        variant={"large"}
                        nameInfo={TooltipsText.CPUTRX}
                    />
                    <NodeStatItem
                        value={"NNN"}
                        name={"CPU-COM"}
                        variant={"large"}
                        nameInfo={TooltipsText.CPUCOM}
                    />
                    <NodeStatItem
                        value={"NNN"}
                        variant={"large"}
                        name={"DISK - TRX"}
                        nameInfo={TooltipsText.DISKTRX}
                    />
                    <NodeStatItem
                        value={"NNN"}
                        variant={"large"}
                        name={"DISK - COM"}
                        nameInfo={TooltipsText.DISKCOM}
                    />
                    <NodeStatItem
                        value={"NNN"}
                        name={"Power"}
                        variant={"large"}
                        nameInfo={TooltipsText.POWER}
                    />
                </NodeStatsContainer>
            </Grid>
            <Grid item lg={isCollapse ? 11 : 9} md xs>
                <Paper sx={{ padding: "22px 18px 0px 30px", width: "100%" }}>
                    <Typography variant="h6">Resources</Typography>
                    <Stack spacing={6}>
                        <StackedAreaChart hasData={true} title={"Memory-TRX"} />
                        <StackedAreaChart hasData={true} title={"Memory-COM"} />
                        <StackedAreaChart hasData={true} title={"CPU-TRX"} />
                        <StackedAreaChart hasData={true} title={"CPU-COM"} />
                        <StackedAreaChart hasData={true} title={"DISK-TRX"} />
                        <StackedAreaChart hasData={true} title={"DISK-COM"} />
                        <StackedAreaChart hasData={true} title={"POWER"} />
                    </Stack>
                </Paper>
            </Grid>
        </Grid>
    );
};

export default NodeResourcesTab;
