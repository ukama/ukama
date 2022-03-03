import {
    TabPanel,
    NodeStatus,
    NodeRadioTab,
    LoadingWrapper,
    NodeNetworkTab,
    NodeSoftwareTab,
    NodeOverviewTab,
    PagePlaceholder,
    NodeResourcesTab,
    NodeAppDetailsDialog,
    NodeSoftwareInfosDialog,
} from "../../components";
import {
    NodeApps,
    NodeAppLogs,
    NodePageTabs,
    NODE_ACTIONS,
} from "../../constants";
import {
    NodeDto,
    useGetNodesByOrgQuery,
    useGetNodeDetailsQuery,
    useGetMetricsCpuTrxLazyQuery,
    MetricDto,
    useGetMetricsUptimeLazyQuery,
    useGetMetricsMemoryTrxLazyQuery,
} from "../../generated";
import { useRecoilValue } from "recoil";
import { isSkeltonLoading, organizationId } from "../../recoil";
import React, { useEffect, useState } from "react";
import { parseObjectInNameValue } from "../../utils";
import { Box, Grid, Paper, Tab, Tabs } from "@mui/material";

const getDefaultList = (names: string[]) =>
    names.map(name => ({
        name: name,
        data: [],
    }));

const Nodes = () => {
    const _organizationId = useRecoilValue(organizationId);
    const [selectedTab, setSelectedTab] = useState(0);
    const skeltonLoading = useRecoilValue(isSkeltonLoading);
    const [selectedNode, setSelectedNode] = useState<NodeDto>();
    const [showNodeAppDialog, setShowNodeAppDialog] = useState(false);
    const [cpuTrxMetrics, setCpuTrxMetrics] = useState<
        {
            name: string;
            data: MetricDto[];
        }[]
    >(getDefaultList(["CPU-TRX (For demo)"]));
    const [uptimeMetrics, setUptimeMetrics] = useState<
        {
            name: string;
            data: MetricDto[];
        }[]
    >(getDefaultList(["UPTIME (For demo)"]));
    const [memoryTrxMetrics, setMemoryTrxMetrics] = useState<
        {
            name: string;
            data: MetricDto[];
        }[]
    >(getDefaultList(["MEMORY-TRX (For demo)"]));

    const [showNodeSoftwareUpdatInfos, setShowNodeSoftwareUpdatInfos] =
        useState<any>();

    const { data: nodesRes, loading: nodesLoading } = useGetNodesByOrgQuery({
        variables: {
            orgId: _organizationId || "a32485e4-d842-45da-bf3e-798889c68ad0",
        },
        onCompleted: res => {
            res?.getNodesByOrg?.nodes.length > 0 &&
                setSelectedNode(res?.getNodesByOrg?.nodes[0]);
        },
    });

    const { data: nodeDetailRes, loading: nodeDetailLoading } =
        useGetNodeDetailsQuery();

    const [
        getCpuTrxMetrics,
        { data: nodeCpuTrxRes, refetch: refetchCpuTrxMetrics },
    ] = useGetMetricsCpuTrxLazyQuery({
        fetchPolicy: "network-only",
        notifyOnNetworkStatusChange: true,
        onCompleted: res => {
            if (res?.getMetricsCpuTRX) {
                setCpuTrxMetrics(
                    cpuTrxMetrics.map(item => {
                        return {
                            name: item.name,
                            data: [...item.data, ...res.getMetricsCpuTRX],
                        };
                    })
                );
            }
        },
    });

    const [
        getUptimeMetrics,
        { data: nodeUptimeMetricsRes, refetch: refetchUptimeMetrics },
    ] = useGetMetricsUptimeLazyQuery({
        fetchPolicy: "network-only",
        notifyOnNetworkStatusChange: true,
        onCompleted: res => {
            if (res?.getMetricsUptime) {
                setUptimeMetrics(
                    uptimeMetrics.map(item => {
                        return {
                            name: item.name,
                            data: [...item.data, ...res.getMetricsUptime],
                        };
                    })
                );
            }
        },
    });

    const [
        getMemoryTrxMetrics,
        { data: nodeMemoryTrxRes, refetch: refetchMemoryTrxMetrics },
    ] = useGetMetricsMemoryTrxLazyQuery({
        fetchPolicy: "network-only",
        notifyOnNetworkStatusChange: true,
        onCompleted: res => {
            if (res?.getMetricsMemoryTRX) {
                setMemoryTrxMetrics(
                    memoryTrxMetrics.map(item => {
                        return {
                            name: item.name,
                            data: [...item.data, ...res.getMetricsMemoryTRX],
                        };
                    })
                );
            }
        },
    });

    useEffect(() => {
        if (selectedTab === 0 && selectedNode) {
            setCpuTrxMetrics(getDefaultList(["UPTIME (For demo)"]));
            setCpuTrxMetrics(getDefaultList(["CPU-TRX (For demo)"]));
            setCpuTrxMetrics(getDefaultList(["MEMORY-TRX (For demo)"]));
            getCpuTrxMetrics({
                variables: {
                    data: {
                        nodeId: selectedNode?.id || "",
                        orgId: _organizationId || "",
                        to: Math.round(Date.now() / 1000),
                        from: Math.round(Date.now() / 1000) - 240,
                        step: 1,
                    },
                },
            });
            getUptimeMetrics({
                variables: {
                    data: {
                        nodeId: selectedNode?.id || "",
                        orgId: _organizationId || "",
                        to: Math.round(Date.now() / 1000),
                        from: Math.round(Date.now() / 1000) - 240,
                        step: 1,
                    },
                },
            });
            getMemoryTrxMetrics({
                variables: {
                    data: {
                        nodeId: selectedNode?.id || "",
                        orgId: _organizationId || "",
                        to: Math.round(Date.now() / 1000),
                        from: Math.round(Date.now() / 1000) - 240,
                        step: 1,
                    },
                },
            });
        }
    }, [selectedTab, selectedNode]);

    const onTabSelected = (event: React.SyntheticEvent, value: any) =>
        setSelectedTab(value);
    const onNodeSelected = (node: NodeDto) => {
        setSelectedNode(node);
    };

    const onUpdateNodeClick = () => {
        //TODO: Handle NODE RESTART ACTION
    };
    const handleNodeActionItemSelected = () => {
        //Todo :Handle nodeAction Itemselected
    };

    const handleNodeActioOptionClicked = () => {
        //Todo :Handle nodeAction selected and clicked
    };
    const onAddNode = () => {
        //TODO: Handle NODE ADD ACTION
    };
    const handleUpdateNode = () => {
        // TODO: Handle Update node Action
    };

    const getNodeDetails = () => {
        //TODO:Handle nodeDetails
        setShowNodeAppDialog(true);
    };
    const handleNodAppDetailsDialog = () => {
        setShowNodeAppDialog(false);
    };

    const fetchCpuTrxData = () =>
        nodeCpuTrxRes &&
        nodeCpuTrxRes.getMetricsCpuTRX.length > 0 &&
        refetchCpuTrxMetrics({
            data: {
                nodeId: selectedNode?.id || "",
                orgId: _organizationId || "",
                to: Math.round(Date.now() / 1000),
                from:
                    nodeCpuTrxRes.getMetricsCpuTRX[
                        nodeCpuTrxRes.getMetricsCpuTRX.length - 1
                    ].x + 1,
                step: 1,
            },
        });

    const fetchUptimeData = () =>
        nodeUptimeMetricsRes &&
        nodeUptimeMetricsRes?.getMetricsUptime.length > 0 &&
        refetchUptimeMetrics({
            data: {
                nodeId: selectedNode?.id || "",
                orgId: _organizationId || "",
                to: Math.round(Date.now() / 1000),
                from:
                    nodeUptimeMetricsRes?.getMetricsUptime[
                        nodeUptimeMetricsRes.getMetricsUptime.length - 1
                    ].x + 1,
                step: 1,
            },
        });

    const fetchMemoryTrxData = () =>
        nodeMemoryTrxRes &&
        nodeMemoryTrxRes.getMetricsMemoryTRX.length > 0 &&
        refetchMemoryTrxMetrics({
            data: {
                nodeId: selectedNode?.id || "",
                orgId: _organizationId || "",
                to: Math.round(Date.now() / 1000),
                from:
                    nodeMemoryTrxRes?.getMetricsMemoryTRX[
                        nodeMemoryTrxRes.getMetricsMemoryTRX.length - 1
                    ].x + 1,
                step: 1,
            },
        });

    const handleCloseNodeInfos = () => {
        setShowNodeSoftwareUpdatInfos(false);
    };
    const handleSoftwareInfos = () => {
        setShowNodeSoftwareUpdatInfos(true);
    };
    const isLoading = skeltonLoading || nodesLoading;

    return (
        <Box
            component="div"
            sx={{
                p: 0,
                mt: 3,
                pb: 2,
            }}
        >
            {nodesRes || isLoading ? (
                <Grid container spacing={3}>
                    <Grid item xs={12}>
                        <NodeStatus
                            onAddNode={onAddNode}
                            loading={nodesLoading}
                            handleNodeActionClick={handleNodeActioOptionClicked}
                            selectedNode={selectedNode}
                            onNodeActionItemSelected={
                                handleNodeActionItemSelected
                            }
                            onNodeSelected={onNodeSelected}
                            nodeActionOptions={NODE_ACTIONS}
                            onUpdateNodeClick={onUpdateNodeClick}
                            nodes={nodesRes?.getNodesByOrg?.nodes || []}
                        />
                    </Grid>
                    <Grid item xs={12}>
                        <LoadingWrapper isLoading={isLoading} height={"40px"}>
                            <Tabs value={selectedTab} onChange={onTabSelected}>
                                {NodePageTabs.map(({ id, label, value }) => (
                                    <Tab
                                        key={id}
                                        label={label}
                                        id={`node-tab-${value}`}
                                    />
                                ))}
                            </Tabs>
                        </LoadingWrapper>
                    </Grid>

                    <Grid item xs={12}>
                        <TabPanel
                            id={"node-tab-0"}
                            value={selectedTab}
                            index={0}
                        >
                            <NodeOverviewTab
                                getNodeSoftwareUpdateInfos={handleSoftwareInfos}
                                isUpdateAvailable={true}
                                selectedNode={selectedNode}
                                cpuTrxMetrics={cpuTrxMetrics}
                                onRefreshTempTrx={fetchCpuTrxData}
                                memoryTrxMetrics={memoryTrxMetrics}
                                onRefreshMemoryTrx={fetchMemoryTrxData}
                                uptimeMetrics={uptimeMetrics}
                                onRefreshUptime={fetchUptimeData}
                                handleUpdateNode={handleUpdateNode}
                                loading={isLoading || nodeDetailLoading}
                                nodeDetails={parseObjectInNameValue(
                                    nodeDetailRes?.getNodeDetails
                                )}
                            />
                        </TabPanel>
                        <TabPanel
                            id={"node-tab-1"}
                            value={selectedTab}
                            index={1}
                        >
                            <NodeNetworkTab
                                loading={isLoading || nodeDetailLoading}
                            />
                        </TabPanel>
                        <TabPanel
                            id={"node-tab-2"}
                            value={selectedTab}
                            index={2}
                        >
                            <NodeResourcesTab
                                loading={isLoading || nodeDetailLoading}
                            />
                        </TabPanel>
                        <TabPanel
                            id={"node-tab-3"}
                            value={selectedTab}
                            index={3}
                        >
                            <NodeRadioTab
                                loading={isLoading || nodeDetailLoading}
                            />
                        </TabPanel>
                        <TabPanel
                            id={"node-tab-4"}
                            value={selectedTab}
                            index={4}
                        >
                            <NodeSoftwareTab
                                loading={isLoading || nodeDetailLoading}
                                nodeApps={NodeApps}
                                NodeLogs={NodeAppLogs}
                                getNodeAppDetails={getNodeDetails}
                            />
                        </TabPanel>
                        <TabPanel
                            id={"node-tab-5"}
                            value={selectedTab}
                            index={5}
                        >
                            <Paper>Schematic</Paper>
                        </TabPanel>
                    </Grid>
                </Grid>
            ) : (
                <PagePlaceholder
                    hyperlink="#"
                    linkText={"here"}
                    showActionButton={false}
                    buttonTitle="Install sims"
                    description="Your nodes have not arrived yet. View their status"
                />
            )}
            <NodeAppDetailsDialog
                closeBtnLabel="close"
                isOpen={showNodeAppDialog}
                handleClose={handleNodAppDetailsDialog}
            />
            <NodeSoftwareInfosDialog
                closeBtnLabel="close"
                isOpen={showNodeSoftwareUpdatInfos}
                handleClose={handleCloseNodeInfos}
            />
        </Box>
    );
};

export default Nodes;
