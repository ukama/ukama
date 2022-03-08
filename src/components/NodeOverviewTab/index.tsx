import { useEffect, useState } from "react";
import { TObject } from "../../types";
import { NodeDto } from "../../generated";
import NodeStatItem from "../NodeStatItem";
import { LineChart, NodeDetailsCard, NodeStatsContainer } from "..";
import { HealtChartsConfigure, TooltipsText } from "../../constants";
import { capitalize, Grid, Paper, Stack, Typography } from "@mui/material";
import ApexLineChartIntegration from "../ApexLineChart";

interface INodeOverviewTab {
    loading: boolean;
    uptimeMetrics: any;
    graphFilters: TObject;
    nodeDetails: TObject[];
    isUpdateAvailable: boolean;
    handleUpdateNode: Function;
    handleGraphFilterChange: Function;
    selectedNode: NodeDto | undefined;
    getNodeSoftwareUpdateInfos: Function;
}

const NodeOverviewTab = ({
    loading,
    selectedNode,
    uptimeMetrics,
    handleUpdateNode,
    isUpdateAvailable,
    handleGraphFilterChange,
    getNodeSoftwareUpdateInfos,
}: INodeOverviewTab) => {
    const [selected, setSelected] = useState<number>(0);

    useEffect(() => {
        setSelected(0);
    }, [selectedNode]);

    const handleOnSelected = (value: number) => setSelected(value);
    const onfilterChange = (key: string, value: string) =>
        handleGraphFilterChange(key, value);

    return (
        <Grid container spacing={2}>
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
                        <NodeStatItem value={"11111111111"} name={"Serial #"} />
                        <NodeStatItem
                            value={"Amplifier Node 1"}
                            name={"Node Group"}
                        />
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
                                showAlertInfo={true}
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
                        nodeTitle={selectedNode?.title || "HOME"}
                        handleUpdateNode={handleUpdateNode}
                        isUpdateAvailable={isUpdateAvailable}
                    />
                )}
                {selected === 1 && (
                    <Paper sx={{ p: 3 }}>
                        <Typography variant="h6">Node Health</Typography>
                        {HealtChartsConfigure[
                            (selectedNode?.type as string) || "HOME"
                        ][0].show && (
                            <LineChart
                                hasData={true}
                                title={
                                    HealtChartsConfigure[
                                        (selectedNode?.type as string) || "HOME"
                                    ][0].name
                                }
                                filter={"DAY"}
                                handleFilterChange={(value: string) =>
                                    onfilterChange("memoryTrx", value)
                                }
                            />
                        )}
                        {HealtChartsConfigure[
                            (selectedNode?.type as string) || "HOME"
                        ][1].show && (
                            <LineChart
                                hasData={true}
                                title={
                                    HealtChartsConfigure[
                                        (selectedNode?.type as string) || "HOME"
                                    ][1].name
                                }
                            />
                        )}
                        {HealtChartsConfigure[
                            (selectedNode?.type as string) || "HOME"
                        ][2].show && (
                            <ApexLineChartIntegration
                                hasData={true}
                                data={uptimeMetrics}
                                name={
                                    HealtChartsConfigure[
                                        (selectedNode?.type as string) || "HOME"
                                    ][2].name
                                }
                                filter={"LIVE"}
                            />
                        )}
                    </Paper>
                )}
                {selected === 2 && selectedNode?.type !== "AMPLIFIER" && (
                    <Paper sx={{ p: 3 }}>
                        <Typography variant="h6">Subscribers</Typography>
                        <LineChart hasData={true} title={"Attached"} />
                        <LineChart hasData={true} title={"Active"} />
                    </Paper>
                )}
            </Grid>
        </Grid>
    );
};

export default NodeOverviewTab;
