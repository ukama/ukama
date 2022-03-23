import { useState } from "react";
import { NodeDto } from "../../generated";
import ApexLineChart from "../ApexLineChart";
import { Paper, Grid, Stack } from "@mui/material";
import { NodeResourcesTabConfigure, TooltipsText } from "../../constants";
import { NodeStatsContainer, NodeStatItem } from "..";
interface INodeResourcesTab {
    loading: boolean;
    cpuTrxMetric: any;
    memoryTrxMetric: any;
    diskTrxMatrics: any;
    powerMetrics: any;
    diskComMetrics: any;
    memoryComMetrics: any;
    cpuComMetrics: any;
    selectedNode: NodeDto | undefined;
}
const NodeResourcesTab = ({
    loading,
    selectedNode,
    diskComMetrics = [],
    diskTrxMatrics = [],
    cpuComMetrics = [],
    powerMetrics = [],
    memoryComMetrics = [],
    cpuTrxMetric = [],
    memoryTrxMetric = [],
}: INodeResourcesTab) => {
    const [isCollapse, setIsCollapse] = useState<boolean>(false);
    const handleCollapse = () => setIsCollapse(prev => !prev);
    return (
        <Grid container spacing={3}>
            <Grid item lg={!isCollapse ? 4 : 1} md xs>
                <NodeStatsContainer
                    index={0}
                    selected={0}
                    loading={loading}
                    title={"Resources"}
                    isCollapsable={true}
                    isCollapse={isCollapse}
                    onCollapse={handleCollapse}
                >
                    {NodeResourcesTabConfigure[
                        (selectedNode?.type as string) || ""
                    ][0].show && (
                        <NodeStatItem
                            value={"NNN"}
                            variant={"large"}
                            name={
                                NodeResourcesTabConfigure[
                                    (selectedNode?.type as string) || ""
                                ][0].name
                            }
                            nameInfo={TooltipsText.MTRX}
                        />
                    )}
                    {NodeResourcesTabConfigure[
                        (selectedNode?.type as string) || ""
                    ][1].show && (
                        <NodeStatItem
                            value={"NNN"}
                            variant={"large"}
                            name={
                                NodeResourcesTabConfigure[
                                    (selectedNode?.type as string) || ""
                                ][1].name
                            }
                            nameInfo={TooltipsText.MCOM}
                        />
                    )}
                    {NodeResourcesTabConfigure[
                        (selectedNode?.type as string) || ""
                    ][2].show && (
                        <NodeStatItem
                            value={"NNN"}
                            name={
                                NodeResourcesTabConfigure[
                                    (selectedNode?.type as string) || ""
                                ][2].name
                            }
                            variant={"large"}
                            nameInfo={TooltipsText.CPUTRX}
                        />
                    )}
                    {NodeResourcesTabConfigure[
                        (selectedNode?.type as string) || ""
                    ][3].show && (
                        <NodeStatItem
                            value={"NNN"}
                            name={
                                NodeResourcesTabConfigure[
                                    (selectedNode?.type as string) || ""
                                ][3].name
                            }
                            variant={"large"}
                            nameInfo={TooltipsText.CPUCOM}
                        />
                    )}
                    {NodeResourcesTabConfigure[
                        (selectedNode?.type as string) || ""
                    ][4].show && (
                        <NodeStatItem
                            value={"NNN"}
                            variant={"large"}
                            name={
                                NodeResourcesTabConfigure[
                                    (selectedNode?.type as string) || ""
                                ][4].name
                            }
                            nameInfo={TooltipsText.DISKTRX}
                        />
                    )}
                    {NodeResourcesTabConfigure[
                        (selectedNode?.type as string) || ""
                    ][5].show && (
                        <NodeStatItem
                            value={"NNN"}
                            variant={"large"}
                            name={
                                NodeResourcesTabConfigure[
                                    (selectedNode?.type as string) || ""
                                ][5].name
                            }
                            nameInfo={TooltipsText.DISKCOM}
                        />
                    )}
                    {NodeResourcesTabConfigure[
                        (selectedNode?.type as string) || ""
                    ][6].show && (
                        <NodeStatItem
                            value={"NNN"}
                            name={
                                NodeResourcesTabConfigure[
                                    (selectedNode?.type as string) || ""
                                ][6].name
                            }
                            variant={"large"}
                            nameInfo={TooltipsText.POWER}
                        />
                    )}
                </NodeStatsContainer>
            </Grid>
            <Grid item lg={isCollapse ? 11 : 8} md xs>
                <Paper sx={{ p: 3, width: "100%" }}>
                    <Stack spacing={4}>
                        {NodeResourcesTabConfigure[
                            (selectedNode?.type as string) || ""
                        ][0].show && (
                            <ApexLineChart
                                hasData={true}
                                name={
                                    NodeResourcesTabConfigure[
                                        (selectedNode?.type as string) || ""
                                    ][0].name
                                }
                                data={memoryTrxMetric}
                            />
                        )}
                        {NodeResourcesTabConfigure[
                            (selectedNode?.type as string) || ""
                        ][1].show && (
                            <ApexLineChart
                                hasData={true}
                                name={
                                    NodeResourcesTabConfigure[
                                        (selectedNode?.type as string) || ""
                                    ][1].name
                                }
                                data={memoryComMetrics}
                            />
                        )}
                        {NodeResourcesTabConfigure[
                            (selectedNode?.type as string) || ""
                        ][2].show && (
                            <ApexLineChart
                                hasData={true}
                                name={
                                    NodeResourcesTabConfigure[
                                        (selectedNode?.type as string) || ""
                                    ][2].name
                                }
                                data={cpuTrxMetric}
                            />
                        )}
                        {NodeResourcesTabConfigure[
                            (selectedNode?.type as string) || ""
                        ][3].show && (
                            <ApexLineChart
                                hasData={true}
                                name={
                                    NodeResourcesTabConfigure[
                                        (selectedNode?.type as string) || ""
                                    ][3].name
                                }
                                data={cpuComMetrics}
                            />
                        )}
                        {NodeResourcesTabConfigure[
                            (selectedNode?.type as string) || ""
                        ][4].show && (
                            <ApexLineChart
                                hasData={true}
                                name={
                                    NodeResourcesTabConfigure[
                                        (selectedNode?.type as string) || ""
                                    ][4].name
                                }
                                data={diskTrxMatrics}
                            />
                        )}
                        {NodeResourcesTabConfigure[
                            (selectedNode?.type as string) || ""
                        ][5].show && (
                            <ApexLineChart
                                hasData={true}
                                name={
                                    NodeResourcesTabConfigure[
                                        (selectedNode?.type as string) || ""
                                    ][5].name
                                }
                                data={diskComMetrics}
                            />
                        )}
                        {NodeResourcesTabConfigure[
                            (selectedNode?.type as string) || ""
                        ][6].show && (
                            <ApexLineChart
                                hasData={true}
                                name={
                                    NodeResourcesTabConfigure[
                                        (selectedNode?.type as string) || ""
                                    ][6].name
                                }
                                data={powerMetrics}
                            />
                        )}
                    </Stack>
                </Paper>
            </Grid>
        </Grid>
    );
};

export default NodeResourcesTab;
