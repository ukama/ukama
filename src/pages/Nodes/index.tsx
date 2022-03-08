/* eslint-disable no-empty-pattern */
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
    MetricDto,
    useGetNodesByOrgQuery,
    useGetNodeDetailsQuery,
    useGetMetricsThroughputUlLazyQuery,
    useGetMetricsThroughputDlLazyQuery,
    useGetMetricsThroughputUlsSubscription,
    useGetMetricsThroughputDlsSubscription,
    useGetMetricsCpuTrxSSubscription,
    useGetMetricsCpuTrxLazyQuery,
    useGetMetricsMemoryTrxLazyQuery,
    useGetMetricsMemoryTrxSSubscription,
    useGetMetricsUptimeLazyQuery,
    useGetMetricsUptimeSSubscription,
} from "../../generated";
import { useRecoilValue } from "recoil";
import { isSkeltonLoading, organizationId } from "../../recoil";
import React, { useEffect, useState } from "react";
import { getMetricPayload, parseObjectInNameValue } from "../../utils";
import { Box, Grid, Paper, Tab, Tabs } from "@mui/material";
import { TObject } from "../../types";

const getDefaultList = (names: string[]) =>
    names.map(name => ({
        name: name,
        data: [],
    }));

const Nodes = () => {
    const _organizationId = useRecoilValue(organizationId) || "";
    const [selectedTab, setSelectedTab] = useState(0);
    const skeltonLoading = useRecoilValue(isSkeltonLoading);
    const [selectedNode, setSelectedNode] = useState<NodeDto>();
    const [showNodeAppDialog, setShowNodeAppDialog] = useState(false);
    const [uptimeMetric, setUptimeMetrics] = useState<
        {
            name: string;
            data: MetricDto[];
        }[]
    >(getDefaultList(["UPTIME (For demo)"]));
    const [cpuTrxMetric, setCpuTrxMetric] = useState<
        {
            name: string;
            data: MetricDto[];
        }[]
    >(getDefaultList(["CPU (TRX)"]));
    const [memoryTrxMetric, setMemoryTrxMetric] = useState<
        {
            name: string;
            data: MetricDto[];
        }[]
    >(getDefaultList(["MEMORY (TRX)"]));
    const [throughputULMetric, setThroughputULMetric] = useState<
        {
            name: string;
            data: MetricDto[];
        }[]
    >(getDefaultList(["Throughput (UL)"]));
    const [throughputDLMetric, setThroughputDLMetric] = useState<
        {
            name: string;
            data: MetricDto[];
        }[]
    >(getDefaultList(["Throughput (DL)"]));

    const [showNodeSoftwareUpdatInfos, setShowNodeSoftwareUpdatInfos] =
        useState<boolean>(false);

    const [graphFilters, setGraphFilters] = useState<TObject>({
        cpuTrx: "DAY",
        uptime: "DAY",
        tempTrx: "DAY",
        tempCom: "DAY",
        memoryTrx: "DAY",
        subActive: "DAY",
        subAttached: "DAY",
    });

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
        getMetricThroughtpuUl,
        { data: metricThroughtputUlRes, refetch: metricThroughtputUlRefetch },
    ] = useGetMetricsThroughputUlLazyQuery();

    const {} = useGetMetricsThroughputUlsSubscription({
        skip: selectedTab !== 1,
        onSubscriptionData: res => {
            setThroughputULMetric(
                throughputULMetric.map(item => {
                    return {
                        name: item.name,
                        data: [
                            ...item.data,
                            ...(res.subscriptionData.data
                                ?.getMetricsThroughputUL || []),
                        ],
                    };
                })
            );
        },
    });

    const [
        getMetricThroughtpuDl,
        { data: metricThroughtputDlRes, refetch: metricThroughtputDlRefetch },
    ] = useGetMetricsThroughputDlLazyQuery();

    const {} = useGetMetricsThroughputDlsSubscription({
        skip: selectedTab !== 1,
        onSubscriptionData: res => {
            setThroughputDLMetric(
                throughputDLMetric.map(item => {
                    return {
                        name: item.name,
                        data: [
                            ...item.data,
                            ...(res.subscriptionData.data
                                ?.getMetricsThroughputDL || []),
                        ],
                    };
                })
            );
        },
    });

    const [
        getMetricCpuTrx,
        { data: metricCpuTrxRes, refetch: metricCpuTrxRefetch },
    ] = useGetMetricsCpuTrxLazyQuery();

    const {} = useGetMetricsCpuTrxSSubscription({
        skip: selectedTab !== 2,
        onSubscriptionData: res => {
            setCpuTrxMetric(
                cpuTrxMetric.map(item => {
                    return {
                        name: item.name,
                        data: [
                            ...item.data,
                            ...(res.subscriptionData.data?.getMetricsCpuTrx ||
                                []),
                        ],
                    };
                })
            );
        },
    });
    const [
        getMetricMemoryTrx,
        { data: metricMemoryTrxRes, refetch: metricMemoryTrxRefetch },
    ] = useGetMetricsMemoryTrxLazyQuery();

    const {} = useGetMetricsMemoryTrxSSubscription({
        skip: selectedTab !== 2,
        onSubscriptionData: res => {
            setMemoryTrxMetric(
                memoryTrxMetric.map(item => {
                    return {
                        name: item.name,
                        data: [
                            ...item.data,
                            ...(res.subscriptionData.data
                                ?.getMetricsMemoryTrx || []),
                        ],
                    };
                })
            );
        },
    });

    const [
        getMetricUptime,
        { data: metricUptimeTrxRes, refetch: metricUptimeRefetch },
    ] = useGetMetricsUptimeLazyQuery();

    const {} = useGetMetricsUptimeSSubscription({
        skip: selectedTab !== 0,
        onSubscriptionData: res => {
            setUptimeMetrics(
                uptimeMetric.map(item => {
                    return {
                        name: item.name,
                        data: [
                            ...item.data,
                            ...(res.subscriptionData.data?.getMetricsUptime ||
                                []),
                        ],
                    };
                })
            );
        },
    });

    useEffect(() => {
        if (selectedTab === 0) {
            getMetricUptime({
                variables: {
                    ...getMetricPayload(
                        "uk-test36-hnode-a1-30df",
                        _organizationId
                    ),
                },
            });
        } else if (selectedTab === 1) {
            getMetricThroughtpuUl({
                variables: {
                    ...getMetricPayload(
                        "uk-test36-hnode-a1-30df",
                        _organizationId
                    ),
                },
            });
            getMetricThroughtpuDl({
                variables: {
                    ...getMetricPayload(
                        "uk-test36-hnode-a1-30df",
                        _organizationId
                    ),
                },
            });
        } else if (selectedTab === 2) {
            getMetricCpuTrx({
                variables: {
                    ...getMetricPayload(
                        "uk-test36-hnode-a1-30df",
                        _organizationId
                    ),
                },
            });
            getMetricMemoryTrx({
                variables: {
                    ...getMetricPayload(
                        "uk-test36-hnode-a1-30df",
                        _organizationId
                    ),
                },
            });
        }
    }, [selectedTab, selectedNode]);

    useEffect(() => {
        if (selectedTab === 0 && metricUptimeTrxRes) {
            metricUptimeRefetch({
                ...getMetricPayload("uk-test36-hnode-a1-30df", _organizationId),
            });
        }
    }, [metricUptimeTrxRes]);

    useEffect(() => {
        if (selectedTab === 1 && metricThroughtputUlRes) {
            metricThroughtputUlRefetch({
                ...getMetricPayload("uk-test36-hnode-a1-30df", _organizationId),
            });
        }
    }, [metricThroughtputUlRes]);

    useEffect(() => {
        if (selectedTab === 1 && metricThroughtputDlRes) {
            metricThroughtputDlRefetch({
                ...getMetricPayload("uk-test36-hnode-a1-30df", _organizationId),
            });
        }
    }, [metricThroughtputDlRes]);

    useEffect(() => {
        if (selectedTab === 2 && metricCpuTrxRes) {
            metricCpuTrxRefetch({
                ...getMetricPayload("uk-test36-hnode-a1-30df", _organizationId),
            });
        }
    }, [metricCpuTrxRes]);

    useEffect(() => {
        if (selectedTab === 2 && metricMemoryTrxRes) {
            metricMemoryTrxRefetch({
                ...getMetricPayload("uk-test36-hnode-a1-30df", _organizationId),
            });
        }
    }, [metricMemoryTrxRes]);

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

    const handleCloseNodeInfos = () => {
        setShowNodeSoftwareUpdatInfos(false);
    };

    const handleSoftwareInfos = () => {
        setShowNodeSoftwareUpdatInfos(true);
    };

    const handleGraphFilterChange = (key: string, value: string) =>
        setGraphFilters(prev => ({ ...prev, [key]: value }));

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
                                        sx={{
                                            display:
                                                (selectedNode?.type ===
                                                    "HOME" &&
                                                    label === "Radio") ||
                                                (selectedNode?.type ===
                                                    "AMPLIFIER" &&
                                                    label === "Network")
                                                    ? "none"
                                                    : "block",
                                        }}
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
                                graphFilters={graphFilters}
                                handleGraphFilterChange={
                                    handleGraphFilterChange
                                }
                                getNodeSoftwareUpdateInfos={handleSoftwareInfos}
                                isUpdateAvailable={true}
                                selectedNode={selectedNode}
                                uptimeMetrics={uptimeMetric}
                                handleUpdateNode={handleUpdateNode}
                                loading={
                                    isLoading ||
                                    nodeDetailLoading ||
                                    nodesLoading ||
                                    !selectedNode
                                }
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
                                throughpuULMetric={throughputULMetric}
                                throughpuDLMetric={throughputDLMetric}
                            />
                        </TabPanel>
                        <TabPanel
                            id={"node-tab-2"}
                            value={selectedTab}
                            index={2}
                        >
                            <NodeResourcesTab
                                selectedNode={selectedNode}
                                cpuTrxMetric={cpuTrxMetric}
                                memoryTrxMetric={memoryTrxMetric}
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
