import { useEffect, useState } from "react";
import { NodeDto } from "../../generated";
import NodeStatItem from "../NodeStatItem";
import { HealtChartsConfigure, TooltipsText } from "../../constants";
import { NodeDetailsCard, NodeStatsContainer, ApexLineChart } from "..";
import { capitalize, Grid, Paper, Stack, Typography } from "@mui/material";

interface INodeOverviewTab {
    metrics: any;
    loading: boolean;
    metricsLoading: boolean;
    isUpdateAvailable: boolean;
    handleUpdateNode: Function;
    selectedNode: NodeDto | undefined;
    getNodeSoftwareUpdateInfos: Function;
}

const NodeOverviewTab = ({
    metrics,
    loading,
    selectedNode,
    metricsLoading,
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
                        <NodeStatItem
                            value={`${capitalize(
                                selectedNode?.type.toLowerCase() || "HOME"
                            )} Node`}
                            name={"Model type"}
                        />
                        <NodeStatItem
                            value={selectedNode?.id || "-"}
                            name={"Serial #"}
                        />
                        {selectedNode?.type === "TOWER" && (
                            <NodeStatItem
                                value={"Amplifier Node 1"}
                                name={"Node Group"}
                            />
                        )}
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
                                value={"50 °C"}
                                name={
                                    HealtChartsConfigure[
                                        (selectedNode?.type as string) || "HOME"
                                    ][0].name
                                }
                                showAlertInfo={false}
                                nameInfo={TooltipsText.TRX}
                                valueInfo={TooltipsText.TRX_ALERT}
                            />
                        )}
                        {HealtChartsConfigure[
                            (selectedNode?.type as string) || "HOME"
                        ][1].show && (
                            <NodeStatItem
                                value={"50 °C"}
                                name={
                                    HealtChartsConfigure[
                                        (selectedNode?.type as string) || "HOME"
                                    ][1].name
                                }
                                nameInfo={TooltipsText.COM}
                                valueInfo={TooltipsText.COM_ALERT}
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
                                value={"200 hours"}
                                nameInfo={TooltipsText.COM}
                                valueInfo={TooltipsText.COM_ALERT}
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
                                value={"100"}
                                nameInfo={TooltipsText.ATTACHED}
                            />
                            <NodeStatItem
                                name={"Active"}
                                value={"100000"}
                                nameInfo={TooltipsText.ACTIVE}
                            />
                        </NodeStatsContainer>
                    )}
                </Stack>
            </Grid>
            <Grid item xs={12} md={8}>
                {selected === 0 && (
                    <NodeDetailsCard
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
