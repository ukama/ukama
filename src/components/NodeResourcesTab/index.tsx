import { TooltipsText } from "../../constants";
import { Stack, Paper, Typography } from "@mui/material";
import { NodeStatsContainer, NodeStatItem, StackedAreaChart } from "..";

interface INodeResourcesTab {
    loading: boolean;
}
const NodeResourcesTab = ({ loading }: INodeResourcesTab) => {
    return (
        <Stack spacing={2} direction="row">
            <NodeStatsContainer
                index={0}
                selected={0}
                loading={loading}
                title={"Network"}
                isCollapsable={true}
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
            <Paper sx={{ p: 3, width: "100%" }}>
                <Typography variant="h6">Resources</Typography>
                <Stack spacing={6} pt={2}>
                    <StackedAreaChart hasData={true} title={"Memory-TRX"} />
                    <StackedAreaChart hasData={true} title={"Memory-COM"} />
                    <StackedAreaChart hasData={true} title={"CPU-TRX"} />
                    <StackedAreaChart hasData={true} title={"CPU-COM"} />
                    <StackedAreaChart hasData={true} title={"DISK-TRX"} />
                    <StackedAreaChart hasData={true} title={"DISK-COM"} />
                    <StackedAreaChart hasData={true} title={"POWER"} />
                </Stack>
            </Paper>
        </Stack>
    );
};

export default NodeResourcesTab;
