import { useEffect, useState } from "react";
import { useRecoilValue } from "recoil";
import { RoundedCard } from "../../styles";
import { Box, Grid, Stack } from "@mui/material";
import { isSkeltonLoading, organizationId } from "../../recoil";
import {
    NodeStatus,
    NodeInfoCard,
    NodeRFKpiTab,
    NodeHealthTab,
    NodeDetailsCard,
    PagePlaceholder,
    NodeMetaDataTab,
} from "../../components";
import {
    NodeDto,
    Graph_Filter,
    useGetNodesByOrgQuery,
    useGetNodeRfkpiqQuery,
    GetNodeRfkpisDocument,
    GetNodeRfkpisSubscription,
    useGetUsersAttachedMetricsQQuery,
    GetUsersAttachedMetricsSDocument,
    GetUsersAttachedMetricsSSubscription,
    useGetThroughputMetricsQQuery,
    GetThroughputMetricsSSubscription,
    GetThroughputMetricsSDocument,
    useGetTemperatureMetricsQQuery,
    GetTemperatureMetricsSDocument,
    GetTemperatureMetricsSSubscription,
    useGetCpuUsageMetricsQQuery,
    GetCpuUsageMetricsSSubscription,
    GetCpuUsageMetricsSDocument,
    useGetMemoryUsageMetricsQQuery,
    GetMemoryUsageMetricsSSubscription,
    GetMemoryUsageMetricsSDocument,
    useGetIoMetricsQQuery,
    GetIoMetricsSSubscription,
    GetIoMetricsSDocument,
    useGetNodeDetailsQuery,
} from "../../generated";
import { TObject } from "../../types";
import { parseObjectInNameValue, uniqueObjectsArray } from "../../utils";

const Nodes = () => {
    const orgId = useRecoilValue(organizationId);
    const [selectedTab, setSelectedTab] = useState(1);
    const skeltonLoading = useRecoilValue(isSkeltonLoading);
    const [selectedNode, setSelectedNode] = useState<NodeDto>();
    const [rfKpiStats, setRfKpiStats] = useState<TObject[]>([]);
    const [metaDataStats, setMetaDataStats] = useState<TObject[]>([]);
    const [healthStats, setHealthStats] = useState<TObject[]>([]);
    const { data: nodesRes, loading: nodesLoading } = useGetNodesByOrgQuery({
        variables: { orgId: orgId || "" },
        onCompleted: res => {
            res.getNodesByOrg.nodes.length > 0 &&
                setSelectedNode(res.getNodesByOrg.nodes[0]);
        },
    });

    const { data: nodeDetailRes, loading: nodeDetailLoading } =
        useGetNodeDetailsQuery();

    const {
        data: nodeRFKpiRes,
        loading: nodeRFKpiLoading,
        subscribeToMore: subscribeToNodeRFKpiMetrics,
    } = useGetNodeRfkpiqQuery({
        variables: {
            filter: Graph_Filter.Week,
        },
    });

    const {
        data: nodeUserAttachRes,
        loading: nodeUserAttachLoading,
        subscribeToMore: subscribeToUserAttachMetrics,
    } = useGetUsersAttachedMetricsQQuery({
        variables: {
            filter: Graph_Filter.Week,
        },
    });

    const {
        data: nodeThroughputRes,
        loading: nodeThroughpuLoading,
        subscribeToMore: subscribeToThroughputMetrics,
    } = useGetThroughputMetricsQQuery({
        variables: {
            filter: Graph_Filter.Week,
        },
    });

    const {
        data: nodeTempMetricsRes,
        loading: nodeTempMetricsLoading,
        subscribeToMore: subscribeToTempMetrics,
    } = useGetTemperatureMetricsQQuery({
        variables: {
            filter: Graph_Filter.Week,
        },
    });

    const {
        data: nodeCpuUsageMetricsRes,
        loading: nodeCpuUsageMetricsLoading,
        subscribeToMore: subscribeToCpuUsageMetrics,
    } = useGetCpuUsageMetricsQQuery({
        variables: {
            filter: Graph_Filter.Week,
        },
    });

    const {
        data: nodeMemoryUsageMetricsRes,
        loading: nodeMemoryUsageMetricsLoading,
        subscribeToMore: subscribeToMemoryUsageMetrics,
    } = useGetMemoryUsageMetricsQQuery({
        variables: {
            filter: Graph_Filter.Week,
        },
    });

    const {
        data: nodeIoMetricsRes,
        loading: nodeIoMetricsLoading,
        subscribeToMore: subscribeToIoMetrics,
    } = useGetIoMetricsQQuery({
        variables: {
            filter: Graph_Filter.Week,
        },
    });

    const nodeRFKpiMetricsSubscription = () =>
        subscribeToNodeRFKpiMetrics<GetNodeRfkpisSubscription>({
            document: GetNodeRfkpisDocument,
            updateQuery: (prev, { subscriptionData }) => {
                if (!subscriptionData.data) return prev;
                const metrics = subscriptionData.data.getNodeRFKPI;
                setRfKpiStats(parseObjectInNameValue(metrics));
                const spreadPrev =
                    prev && prev.getNodeRFKPI ? prev.getNodeRFKPI : [];
                return Object.assign({}, prev, {
                    getNodeRFKPI: [metrics, ...spreadPrev],
                });
            },
        });

    const nodeUserAttachMetricsSubscription = () =>
        subscribeToUserAttachMetrics<GetUsersAttachedMetricsSSubscription>({
            document: GetUsersAttachedMetricsSDocument,
            updateQuery: (prev, { subscriptionData }) => {
                if (!subscriptionData.data) return prev;
                const metrics = subscriptionData.data.getUsersAttachedMetrics;
                setMetaDataStats(_prev => [
                    ...uniqueObjectsArray("Users Attached", _prev),
                    { name: "Users Attached", value: metrics.users },
                ]);
                const spreadPrev =
                    prev && prev.getUsersAttachedMetrics
                        ? prev.getUsersAttachedMetrics
                        : [];
                return Object.assign({}, prev, {
                    getUsersAttachedMetrics: [metrics, ...spreadPrev],
                });
            },
        });

    const nodeThroughputMetricsSubscription = () =>
        subscribeToThroughputMetrics<GetThroughputMetricsSSubscription>({
            document: GetThroughputMetricsSDocument,
            updateQuery: (prev, { subscriptionData }) => {
                if (!subscriptionData.data) return prev;
                const metrics = subscriptionData.data.getThroughputMetrics;
                setMetaDataStats(_prev => [
                    ...uniqueObjectsArray("Throughput", _prev),
                    { name: "Throughput", value: metrics.amount },
                ]);
                const spreadPrev =
                    prev && prev.getThroughputMetrics
                        ? prev.getThroughputMetrics
                        : [];
                return Object.assign({}, prev, {
                    getThroughputMetrics: [metrics, ...spreadPrev],
                });
            },
        });

    const nodeTempMetricsSubscription = () =>
        subscribeToTempMetrics<GetTemperatureMetricsSSubscription>({
            document: GetTemperatureMetricsSDocument,
            updateQuery: (prev, { subscriptionData }) => {
                if (!subscriptionData.data) return prev;
                const metrics = subscriptionData.data.getTemperatureMetrics;
                setHealthStats(_prev => [
                    ...uniqueObjectsArray("Temperature", _prev),
                    {
                        name: "Temperature",
                        value: metrics.temperature,
                    },
                ]);
                const spreadPrev =
                    prev && prev.getTemperatureMetrics
                        ? prev.getTemperatureMetrics
                        : [];
                return Object.assign({}, prev, {
                    getTemperatureMetrics: [metrics, ...spreadPrev],
                });
            },
        });

    const nodeCpuUsageMetricsSubscription = () =>
        subscribeToCpuUsageMetrics<GetCpuUsageMetricsSSubscription>({
            document: GetCpuUsageMetricsSDocument,
            updateQuery: (prev, { subscriptionData }) => {
                if (!subscriptionData.data) return prev;
                const metrics = subscriptionData.data.getCpuUsageMetrics;
                setHealthStats(_prev => [
                    ...uniqueObjectsArray("CPU", _prev),
                    {
                        name: "CPU",
                        value: `${metrics.usage}%`,
                    },
                ]);
                const spreadPrev =
                    prev && prev.getCpuUsageMetrics
                        ? prev.getCpuUsageMetrics
                        : [];
                return Object.assign({}, prev, {
                    getCpuUsageMetrics: [metrics, ...spreadPrev],
                });
            },
        });

    const nodeMemoryUsageMetricsSubscription = () =>
        subscribeToMemoryUsageMetrics<GetMemoryUsageMetricsSSubscription>({
            document: GetMemoryUsageMetricsSDocument,
            updateQuery: (prev, { subscriptionData }) => {
                if (!subscriptionData.data) return prev;
                const metrics = subscriptionData.data.getMemoryUsageMetrics;
                setHealthStats(_prev => [
                    ...uniqueObjectsArray("Memory", _prev),
                    {
                        name: "Memory",
                        value: `${metrics.usage}%`,
                    },
                ]);
                const spreadPrev =
                    prev && prev.getMemoryUsageMetrics
                        ? prev.getMemoryUsageMetrics
                        : [];
                return Object.assign({}, prev, {
                    getMemoryUsageMetrics: [metrics, ...spreadPrev],
                });
            },
        });

    const nodeIoMetricsSubscription = () =>
        subscribeToIoMetrics<GetIoMetricsSSubscription>({
            document: GetIoMetricsSDocument,
            updateQuery: (prev, { subscriptionData }) => {
                if (!subscriptionData.data) return prev;
                const metrics = subscriptionData.data.getIOMetrics;
                setHealthStats(_prev => [
                    ...uniqueObjectsArray("IO", _prev),
                    {
                        name: "IO",
                        value: `${metrics.input} Input | ${metrics.output} Output`,
                    },
                ]);
                const spreadPrev =
                    prev && prev.getIOMetrics ? prev.getIOMetrics : [];
                return Object.assign({}, prev, {
                    getIOMetrics: [metrics, ...spreadPrev],
                });
            },
        });

    useEffect(() => {
        if (nodeRFKpiRes) {
            nodeRFKpiRes?.getNodeRFKPI?.length > 0 &&
                setRfKpiStats(
                    parseObjectInNameValue(nodeRFKpiRes?.getNodeRFKPI[0])
                );
        }
        let unsub = nodeRFKpiMetricsSubscription();
        return () => {
            unsub && unsub();
        };
    }, [nodeRFKpiRes]);

    useEffect(() => {
        if (nodeUserAttachRes) {
            nodeUserAttachRes?.getUsersAttachedMetrics?.length > 0 &&
                setMetaDataStats(prev => [
                    ...uniqueObjectsArray("Users Attached", prev),
                    {
                        name: "Users Attached",
                        value: nodeUserAttachRes.getUsersAttachedMetrics[0]
                            .users,
                    },
                ]);
        }
        let unsub = nodeUserAttachMetricsSubscription();
        return () => {
            unsub && unsub();
        };
    }, [nodeUserAttachRes]);

    useEffect(() => {
        if (nodeThroughputRes) {
            nodeThroughputRes?.getThroughputMetrics?.length > 0 &&
                setMetaDataStats(prev => [
                    ...uniqueObjectsArray("Throughput", prev),
                    {
                        name: "Throughput",
                        value: nodeThroughputRes?.getThroughputMetrics[0]
                            ?.amount,
                    },
                ]);
        }
        let unsub = nodeThroughputMetricsSubscription();
        return () => {
            unsub && unsub();
        };
    }, [nodeThroughputRes]);

    useEffect(() => {
        if (nodeTempMetricsRes) {
            nodeTempMetricsRes?.getTemperatureMetrics?.length > 0 &&
                setHealthStats(prev => [
                    ...uniqueObjectsArray("Temperature", prev),
                    {
                        name: "Temperature",
                        value: nodeTempMetricsRes?.getTemperatureMetrics[0]
                            ?.temperature,
                    },
                ]);
        }
        let unsub = nodeTempMetricsSubscription();
        return () => {
            unsub && unsub();
        };
    }, [nodeTempMetricsRes]);

    useEffect(() => {
        if (nodeCpuUsageMetricsRes) {
            nodeCpuUsageMetricsRes?.getCpuUsageMetrics?.length > 0 &&
                setHealthStats(prev => [
                    ...uniqueObjectsArray("CPU", prev),
                    {
                        name: "CPU",
                        value: `${nodeCpuUsageMetricsRes?.getCpuUsageMetrics[0]?.usage}%`,
                    },
                ]);
        }
        let unsub = nodeCpuUsageMetricsSubscription();
        return () => {
            unsub && unsub();
        };
    }, [nodeCpuUsageMetricsRes]);

    useEffect(() => {
        if (nodeMemoryUsageMetricsRes) {
            nodeMemoryUsageMetricsRes?.getMemoryUsageMetrics?.length > 0 &&
                setHealthStats(prev => [
                    ...uniqueObjectsArray("Memory", prev),
                    {
                        name: "Memory",
                        value: `${nodeMemoryUsageMetricsRes?.getMemoryUsageMetrics[0]?.usage}%`,
                    },
                ]);
        }
        let unsub = nodeMemoryUsageMetricsSubscription();
        return () => {
            unsub && unsub();
        };
    }, [nodeMemoryUsageMetricsRes]);

    useEffect(() => {
        if (nodeIoMetricsRes) {
            nodeIoMetricsRes?.getIOMetrics?.length > 0 &&
                setHealthStats(prev => [
                    ...uniqueObjectsArray("IO", prev),
                    {
                        name: "IO",
                        value: `${nodeIoMetricsRes.getIOMetrics[0].input} Input | ${nodeIoMetricsRes.getIOMetrics[0].output} Output`,
                    },
                ]);
        }
        let unsub = nodeIoMetricsSubscription();
        return () => {
            unsub && unsub();
        };
    }, [nodeIoMetricsRes]);

    const onTabSelected = (value: number) => setSelectedTab(value);
    const onNodeSelected = (node: NodeDto) => setSelectedNode(node);
    const onNodeRFClick = () => {
        //TODO: Handle NODE RF ACTIONS
    };
    const onNodeSwitchClick = () => {
        //TODO: Handle NODE ON/OFF ACTIONS
    };
    const onRestartNodeClick = () => {
        //TODO: Handle NODE RESTART ACTION
    };

    const isLoading = skeltonLoading || nodesLoading;

    if (nodesRes && nodesRes?.getNodesByOrg?.nodes?.length === 0)
        return (
            <RoundedCard
                sx={{
                    p: 0,
                    mt: 3,
                    mb: 2,
                    borderRadius: "4px",
                    height: "calc(100% - 15%)",
                }}
            >
                <PagePlaceholder description="Order your node now." />
            </RoundedCard>
        );

    return (
        <Box
            sx={{
                p: 0,
                mt: 3,
                pb: 2,
            }}
        >
            <Grid container spacing={2}>
                <Grid item xs={12}>
                    <NodeStatus
                        loading={isLoading}
                        selectedNode={selectedNode}
                        onNodeRFClick={onNodeRFClick}
                        onNodeSelected={onNodeSelected}
                        onNodeSwitchClick={onNodeSwitchClick}
                        onRestartNodeClick={onRestartNodeClick}
                        nodes={nodesRes?.getNodesByOrg?.nodes}
                    />
                </Grid>
                <Grid item container xs={4}>
                    <Stack spacing={2} sx={{ width: "100%" }}>
                        <NodeInfoCard
                            index={1}
                            title={"Node Details"}
                            loading={isLoading || nodeDetailLoading}
                            properties={
                                parseObjectInNameValue(
                                    nodeDetailRes?.getNodeDetails || {}
                                ) || []
                            }
                            onSelected={onTabSelected}
                            isSelected={selectedTab === 1}
                        />
                        <NodeInfoCard
                            index={2}
                            title={"Meta Data"}
                            properties={metaDataStats}
                            onSelected={onTabSelected}
                            isSelected={selectedTab === 2}
                            loading={
                                isLoading ||
                                nodeUserAttachLoading ||
                                nodeThroughpuLoading
                            }
                        />
                        <NodeInfoCard
                            index={3}
                            loading={
                                isLoading ||
                                nodeIoMetricsLoading ||
                                nodeTempMetricsLoading ||
                                nodeCpuUsageMetricsLoading ||
                                nodeMemoryUsageMetricsLoading
                            }
                            title={"Physical Health"}
                            properties={healthStats}
                            onSelected={onTabSelected}
                            isSelected={selectedTab === 3}
                        />
                        <NodeInfoCard
                            index={4}
                            title={"RF KPIs"}
                            properties={rfKpiStats}
                            onSelected={onTabSelected}
                            isSelected={selectedTab === 4}
                            loading={isLoading || nodeRFKpiLoading}
                        />
                    </Stack>
                </Grid>
                <Grid item container xs={8}>
                    {selectedTab === 1 && (
                        <NodeDetailsCard loading={isLoading} />
                    )}
                    {selectedTab === 2 && (
                        <NodeMetaDataTab
                            loading={
                                isLoading ||
                                nodeUserAttachLoading ||
                                nodeThroughpuLoading
                            }
                            usersAttachedMetrics={
                                nodeUserAttachRes?.getUsersAttachedMetrics || []
                            }
                            throughputMetrics={
                                nodeThroughputRes?.getThroughputMetrics || []
                            }
                        />
                    )}
                    {selectedTab === 3 && (
                        <NodeHealthTab
                            loading={isLoading || nodeTempMetricsLoading}
                            temperatureMetrics={
                                nodeTempMetricsRes?.getTemperatureMetrics || []
                            }
                            cpuUsageMetrics={
                                nodeCpuUsageMetricsRes?.getCpuUsageMetrics || []
                            }
                            memoryUsageMetrics={
                                nodeMemoryUsageMetricsRes?.getMemoryUsageMetrics ||
                                []
                            }
                            ioMetrics={nodeIoMetricsRes?.getIOMetrics || []}
                        />
                    )}
                    {selectedTab === 4 && (
                        <NodeRFKpiTab
                            loading={isLoading || nodeRFKpiLoading}
                            metrics={nodeRFKpiRes?.getNodeRFKPI || []}
                        />
                    )}
                </Grid>
            </Grid>
        </Box>
    );
};

export default Nodes;
