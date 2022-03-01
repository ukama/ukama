import { useState } from "react";
import { TObject } from "../../types";
import { NodeDto } from "../../generated";
import NodeStatItem from "../NodeStatItem";
import { TooltipsText } from "../../constants";
import ApexLineChartIntegration from "../ApexLineChart";
import { NodeDetailsCard, NodeStatsContainer } from "..";
import { Grid, Paper, Stack, Typography } from "@mui/material";

const REFRESH_INTERVAL = 10000; // 20 sec

interface INodeOverviewTab {
    loading: boolean;
    cpuTrxMetrics: any;
    uptimeMetrics: any;
    memoryTrxMetrics: any;
    nodeDetails: TObject[];
    isUpdateAvailable: boolean;
    handleUpdateNode: Function;
    onRefreshTempTrx: Function;
    onRefreshMemoryTrx: Function;
    onRefreshUptime: Function;
    selectedNode: NodeDto | undefined;
}

const NodeOverviewTab = ({
    loading,
    selectedNode,
    handleUpdateNode,
    onRefreshTempTrx,
    onRefreshMemoryTrx,
    onRefreshUptime,
    isUpdateAvailable,
    cpuTrxMetrics = [],
    uptimeMetrics = [],
    memoryTrxMetrics = [],
}: INodeOverviewTab) => {
    const [selected, setSelected] = useState<number>(0);

    const handleOnSelected = (value: number) => setSelected(value);

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
                        <NodeStatItem value={"Home Node"} name={"Model type"} />
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
                        <NodeStatItem
                            value={"50 °C"}
                            name={"Temp. (TRX)"}
                            showAlertInfo={true}
                            nameInfo={TooltipsText.TRX}
                            valueInfo={TooltipsText.TRX_ALERT}
                        />
                        <NodeStatItem
                            value={"50 °C"}
                            name={"Temp. (COM)"}
                            nameInfo={TooltipsText.COM}
                            valueInfo={TooltipsText.COM_ALERT}
                        />
                    </NodeStatsContainer>
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
                </Stack>
            </Grid>
            <Grid item xs={12} md={8}>
                {selected === 0 && (
                    <NodeDetailsCard
                        loading={loading}
                        nodeTitle={selectedNode?.title || ""}
                        handleUpdateNode={handleUpdateNode}
                        isUpdateAvailable={isUpdateAvailable}
                    />
                )}
                {selected === 1 && (
                    <Paper sx={{ padding: "28px 18px" }}>
                        <Typography variant="h6">Node Health</Typography>
                        <Stack spacing={6} pt={2}>
                            <ApexLineChartIntegration
                                hasData={true}
                                data={cpuTrxMetrics}
                                name={"CPU-TRX (For demo)"} //"Temperature-TRX"
                                onRefreshData={onRefreshTempTrx}
                                refreshInterval={REFRESH_INTERVAL}
                            />
                            <ApexLineChartIntegration
                                hasData={true}
                                data={uptimeMetrics}
                                name={"UPTIME (For demo)"} //"Temperature-TRX"
                                onRefreshData={onRefreshUptime}
                                refreshInterval={REFRESH_INTERVAL}
                            />
                            <ApexLineChartIntegration
                                hasData={true}
                                data={memoryTrxMetrics}
                                name={"MEMORY-TRX (For demo)"} //"Temperature-TRX"
                                onRefreshData={onRefreshMemoryTrx}
                                refreshInterval={REFRESH_INTERVAL}
                            />
                            {/* <ApexLineChartIntegration
                                hasData={true}
                                initData={[]}
                                name={"Temperature-COM"}
                            /> */}
                        </Stack>
                    </Paper>
                )}
                {selected === 2 && (
                    <Paper sx={{ padding: "28px 18px" }}>
                        <Typography variant="h6">Subscribers</Typography>
                        <Stack spacing={6} pt={2}>
                            {/* <ApexLineChartIntegration
                                hasData={true}
                                initData={[]}
                                name={"Attached"}
                            />
                            <ApexLineChartIntegration
                                hasData={true}
                                initData={[]}
                                name={"Active"}
                            /> */}
                        </Stack>
                    </Paper>
                )}
            </Grid>
        </Grid>
    );
};

export default NodeOverviewTab;
