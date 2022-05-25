import {
    NodeGroup,
    ApexLineChart,
    NodeDetailsCard,
    NodeStatsContainer,
} from "..";
import NodeStatItem from "../NodeStatItem";
import { useEffect, useState } from "react";
import { NodeDto, NodeResponse, Node_Type } from "../../generated";
import { HealtChartsConfigure, TooltipsText } from "../../constants";
import { capitalize, Grid, Paper, Stack, Typography } from "@mui/material";

interface INodeOverviewTab {
    metrics: any;
    loading: boolean;
    metricsLoading: boolean;
    onNodeSelected: Function;
    nodeGroupLoading: boolean;
    uptime: number | undefined;
    isUpdateAvailable: boolean;
    handleUpdateNode: Function;
    selectedNode: NodeDto | undefined;
    connectedUsers: string | undefined;
    getNodeSoftwareUpdateInfos: Function;
    nodeGroupData: NodeResponse | undefined;
}

const NodeOverviewTab = ({
    metrics,
    uptime,
    loading,
    selectedNode,
    nodeGroupData,
    metricsLoading,
    connectedUsers = "0",
    onNodeSelected,
    nodeGroupLoading,
    handleUpdateNode,
    isUpdateAvailable,
    getNodeSoftwareUpdateInfos,
}: INodeOverviewTab) => {
    const [selected, setSelected] = useState<number>(0);

    useEffect(() => {
        setSelected(0);
    }, [selectedNode]);

    const handleOnSelected = (value: number) => setSelected(value);

    return (
        <Grid container spacing={3}>
            <Grid item xs={12} md={4}>
                <Stack spacing={2}>
                    <NodeStatsContainer
                        index={0}
                        loading={loading}
                        isClickable={true}
                        selected={selected}
                        title={"Node Information"}
                        handleAction={handleOnSelected}
                    >
                        <Grid container spacing={1}>
                            <Grid item xs={12}>
                                <NodeStatItem
                                    value={`${capitalize(
                                        selectedNode?.type.toLowerCase() ||
                                            "HOME"
                                    )} Node`}
                                    name={"Model type"}
                                />
                            </Grid>
                            <Grid item xs={12}>
                                <NodeStatItem
                                    value={
                                        selectedNode?.id.toLowerCase() || "-"
                                    }
                                    name={"Serial #"}
                                />
                            </Grid>
                            {selectedNode?.type === "TOWER" && (
                                <Grid item xs={12}>
                                    <NodeGroup
                                        nodes={nodeGroupData?.attached || []}
                                        loading={nodeGroupLoading}
                                        handleNodeAction={onNodeSelected}
                                    />
                                </Grid>
                            )}
                        </Grid>
                    </NodeStatsContainer>
                    <NodeStatsContainer
                        index={1}
                        loading={loading}
                        isClickable={true}
                        selected={selected}
                        title={"Node Health"}
                        handleAction={handleOnSelected}
                    >
                        {HealtChartsConfigure[
                            (selectedNode?.type as string) || "HOME"
                        ][0].show && (
                            <NodeStatItem
                                value={"24 °C"}
                                name={
                                    HealtChartsConfigure[
                                        (selectedNode?.type as string) || "HOME"
                                    ][0].name
                                }
                                showAlertInfo={false}
                                nameInfo={TooltipsText.TRX}
                            />
                        )}
                        {HealtChartsConfigure[
                            (selectedNode?.type as string) || "HOME"
                        ][1].show && (
                            <NodeStatItem
                                value={"22 °C"}
                                name={
                                    HealtChartsConfigure[
                                        (selectedNode?.type as string) || "HOME"
                                    ][1].name
                                }
                                nameInfo={TooltipsText.COM}
                            />
                        )}
                        {HealtChartsConfigure[
                            (selectedNode?.type as string) || "HOME"
                        ][2].show && (
                            <NodeStatItem
                                name={
                                    HealtChartsConfigure[
                                        (selectedNode?.type as string) || "HOME"
                                    ][2].name
                                }
                                nameInfo={TooltipsText.COM}
                                value={
                                    uptime
                                        ? `${Math.floor(
                                              uptime / 60 / 60
                                          )} hours`
                                        : "NA"
                                }
                            />
                        )}
                    </NodeStatsContainer>
                    {selectedNode?.type !== "AMPLIFIER" && (
                        <NodeStatsContainer
                            index={2}
                            loading={loading}
                            isClickable={true}
                            selected={selected}
                            title={"Subscribers"}
                            handleAction={handleOnSelected}
                        >
                            <NodeStatItem
                                name={"Attached"}
                                value={connectedUsers}
                                nameInfo={TooltipsText.ATTACHED}
                            />
                            <NodeStatItem
                                name={"Active"}
                                value={`${
                                    connectedUsers === "0"
                                        ? parseInt(connectedUsers)
                                        : parseInt(connectedUsers) - 1
                                }`}
                                nameInfo={TooltipsText.ACTIVE}
                            />
                        </NodeStatsContainer>
                    )}
                </Stack>
            </Grid>
            <Grid item xs={12} md={8}>
                {selected === 0 && (
                    <NodeDetailsCard
                        nodeType={
                            (selectedNode?.type as Node_Type) || undefined
                        }
                        getNodeUpdateInfos={getNodeSoftwareUpdateInfos}
                        loading={loading}
                        nodeTitle={selectedNode?.name || "HOME"}
                        handleUpdateNode={handleUpdateNode}
                        isUpdateAvailable={isUpdateAvailable}
                    />
                )}
                {selected === 1 && (
                    <Paper sx={{ p: 3 }}>
                        <Stack spacing={4}>
                            <Typography variant="h6">Node Health</Typography>
                            {HealtChartsConfigure[
                                (selectedNode?.type as string) || "HOME"
                            ][0].show && (
                                <ApexLineChart
                                    loading={metricsLoading}
                                    data={
                                        metrics[
                                            HealtChartsConfigure[
                                                (selectedNode?.type as string) ||
                                                    "HOME"
                                            ][0].id
                                        ]
                                    }
                                />
                            )}
                            {HealtChartsConfigure[
                                (selectedNode?.type as string) || "HOME"
                            ][1].show && (
                                <ApexLineChart
                                    loading={metricsLoading}
                                    data={
                                        metrics[
                                            HealtChartsConfigure[
                                                (selectedNode?.type as string) ||
                                                    "HOME"
                                            ][1].id
                                        ]
                                    }
                                />
                            )}
                            {HealtChartsConfigure[
                                (selectedNode?.type as string) || "HOME"
                            ][2].show && (
                                <ApexLineChart
                                    loading={metricsLoading}
                                    data={
                                        metrics[
                                            HealtChartsConfigure[
                                                (selectedNode?.type as string) ||
                                                    "HOME"
                                            ][2].id
                                        ]
                                    }
                                />
                            )}
                        </Stack>
                    </Paper>
                )}
                {selected === 2 && selectedNode?.type !== "AMPLIFIER" && (
                    <Paper sx={{ p: 3 }}>
                        <Stack spacing={4}>
                            <Typography variant="h6">Subscribers</Typography>
                            {HealtChartsConfigure[
                                (selectedNode?.type as string) || "HOME"
                            ][3].show && (
                                <ApexLineChart
                                    loading={metricsLoading}
                                    data={
                                        metrics[
                                            HealtChartsConfigure[
                                                (selectedNode?.type as string) ||
                                                    "HOME"
                                            ][3].id
                                        ]
                                    }
                                />
                            )}
                            {HealtChartsConfigure[
                                (selectedNode?.type as string) || "HOME"
                            ][4].show && (
                                <ApexLineChart
                                    loading={metricsLoading}
                                    data={
                                        metrics[
                                            HealtChartsConfigure[
                                                (selectedNode?.type as string) ||
                                                    "HOME"
                                            ][4].id
                                        ]
                                    }
                                />
                            )}
                        </Stack>
                    </Paper>
                )}
            </Grid>
        </Grid>
    );
};

export default NodeOverviewTab;
