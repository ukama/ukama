import { useState } from "react";
import { NodeDto } from "../../generated";
import ApexLineChart from "../ApexLineChart";
import { Paper, Grid, Stack } from "@mui/material";
import { NodeResourcesTabConfigure, TooltipsText } from "../../constants";
import { NodeStatsContainer, NodeStatItem } from "..";
interface INodeResourcesTab {
    metrics: any;
    loading: boolean;
    selectedNode: NodeDto | undefined;
}
const NodeResourcesTab = ({
    metrics,
    loading,
    selectedNode,
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
                                data={
                                    metrics[
                                        NodeResourcesTabConfigure[
                                            (selectedNode?.type as string) || ""
                                        ][0].id
                                    ]
                                }
                            />
                        )}
                        {NodeResourcesTabConfigure[
                            (selectedNode?.type as string) || ""
                        ][1].show && (
                            <ApexLineChart
                                data={
                                    metrics[
                                        NodeResourcesTabConfigure[
                                            (selectedNode?.type as string) || ""
                                        ][1].id
                                    ]
                                }
                            />
                        )}
                        {NodeResourcesTabConfigure[
                            (selectedNode?.type as string) || ""
                        ][2].show && (
                            <ApexLineChart
                                data={
                                    metrics[
                                        NodeResourcesTabConfigure[
                                            (selectedNode?.type as string) || ""
                                        ][2].id
                                    ]
                                }
                            />
                        )}
                        {NodeResourcesTabConfigure[
                            (selectedNode?.type as string) || ""
                        ][3].show && (
                            <ApexLineChart
                                data={
                                    metrics[
                                        NodeResourcesTabConfigure[
                                            (selectedNode?.type as string) || ""
                                        ][3].id
                                    ]
                                }
                            />
                        )}
                        {NodeResourcesTabConfigure[
                            (selectedNode?.type as string) || ""
                        ][4].show && (
                            <ApexLineChart
                                data={
                                    metrics[
                                        NodeResourcesTabConfigure[
                                            (selectedNode?.type as string) || ""
                                        ][4].id
                                    ]
                                }
                            />
                        )}
                        {NodeResourcesTabConfigure[
                            (selectedNode?.type as string) || ""
                        ][5].show && (
                            <ApexLineChart
                                data={
                                    metrics[
                                        NodeResourcesTabConfigure[
                                            (selectedNode?.type as string) || ""
                                        ][5].id
                                    ]
                                }
                            />
                        )}
                        {NodeResourcesTabConfigure[
                            (selectedNode?.type as string) || ""
                        ][6].show && (
                            <ApexLineChart
                                data={
                                    metrics[
                                        NodeResourcesTabConfigure[
                                            (selectedNode?.type as string) || ""
                                        ][6].id
                                    ]
                                }
                            />
                        )}
                    </Stack>
                </Paper>
            </Grid>
        </Grid>
    );
};

export default NodeResourcesTab;
