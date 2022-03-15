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
    useGetMetricsCpuTrxLazyQuery,
    useGetMetricsEraBsSubscription,
    useGetMetricsRlCsSubscription,
    useGetMetricsErabLazyQuery,
    useGetMetricsRrCsSubscription,
    useGetMetricsRlcLazyQuery,
    useGetMetricsMemoryTrxLazyQuery,
    useGetMetricsUptimeLazyQuery,
    useGetMetricsUptimeSSubscription,
    useGetMetricsMemoryTrxsSubscription,
    useGetMetricsCpuTrxsSubscription,
    useGetMetricsPowerLazyQuery,
    useGetMetricsPowerSSubscription,
    useGetMetricsTempTrxLazyQuery,
    useGetMetricsTempTrxsSubscription,
    useGetMetricsTempComLazyQuery,
    useGetMetricsTempComsSubscription,
    useGetMetricsRrcLazyQuery,
} from "../../generated";
import { useRecoilValue } from "recoil";
import { isSkeltonLoading, organizationId } from "../../recoil";
import React, { useEffect, useState } from "react";
import {
    getMetricPayload,
    isMetricData,
    parseObjectInNameValue,
} from "../../utils";
import { Box, Grid, Paper, Tab, Tabs } from "@mui/material";

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
    const [rrcCnxSuccessMetrix, setRrcCnxSuccessMetrix] = useState<
        {
            name: string;
            data: MetricDto[];
        }[]
    >(getDefaultList(["RRC CNX Success"]));
    const [erabDropRateMetrix, setErabDropRateMetrix] = useState<
        {
            name: string;
            data: MetricDto[];
        }[]
    >(getDefaultList(["ERAB Drop Rate"]));
    const [rlsDropRateMetrics, setRlsDropRateMetrics] = useState<
        {
            name: string;
            data: MetricDto[];
        }[]
    >(getDefaultList(["RLS Drop Rate"]));
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
    const [powerMetric, setPowerMetric] = useState<
        {
            name: string;
            data: MetricDto[];
        }[]
    >(getDefaultList(["Power"]));
    const [tempTrxMetric, setTempTrxMetric] = useState<
        {
            name: string;
            data: MetricDto[];
        }[]
    >(getDefaultList(["Temp. (TRX)"]));
    const [tempComMetric, setTempComMetric] = useState<
        {
            name: string;
            data: MetricDto[];
        }[]
    >(getDefaultList(["Temp. (COM)"]));

    const [showNodeSoftwareUpdatInfos, setShowNodeSoftwareUpdatInfos] =
        useState<boolean>(false);

    const getFirstMetricCallPayload = () =>
        getMetricPayload({
            nodeId: selectedNode?.id,
            orgId: _organizationId,
            regPolling: false,
            to: Math.floor(Date.now() / 1000) - 20,
            from: Math.floor(Date.now() / 1000) - 180,
        });

    const getMetricPollingCallPayload = (from: number) =>
        getMetricPayload({
            nodeId: selectedNode?.id,
            orgId: _organizationId,
            from: from + 1,
        });

    const { data: nodesRes, loading: nodesLoading } = useGetNodesByOrgQuery({
        variables: {
            orgId: _organizationId,
        },
        onCompleted: res => {
            res?.getNodesByOrg?.nodes.length > 0 &&
                setSelectedNode(res?.getNodesByOrg?.nodes[0]);
        },
    });

    const { data: nodeDetailRes, loading: nodeDetailLoading } =
        useGetNodeDetailsQuery();
    // const [getMetricsRLC, { data: metricsRLCres, refetch: metricsRLCRefetch }] =
    //     useGetMetricsRlcLazyQuery();
    // useGetMetricsRlcsSubscription({
    //     skip: selectedTab !== 1,
    //     onSubscriptionData: res => {
    //         setRlsDropRateMetrics(
    //             rlsDropRateMetrics.map(item => {
    //                 return {
    //                     name: item.name,
    //                     data: [
    //                         ...item.data,
    //                         ...(res.subscriptionData.data?.getMetricsRLC || []),
    //                     ],
    //                 };
    //             })
    //         );
    //     },
    // });
    const [getMetricsRLC, { data: metricsRLCres, refetch: metricsRlcRefetch }] =
        useGetMetricsRlcLazyQuery();

    useGetMetricsRlCsSubscription({
        skip: selectedTab !== 1,
        onSubscriptionData: res => {
            setRlsDropRateMetrics(
                rlsDropRateMetrics.map(item => {
                    return {
                        name: item.name,
                        data: [
                            ...item.data,
                            ...(res.subscriptionData.data?.getMetricsRLC || []),
                        ],
                    };
                })
            );
        },
    });
    const [
        getMetricsERAB,
        { data: metricsERABres, refetch: metricsERABresRefetch },
    ] = useGetMetricsErabLazyQuery();

    useGetMetricsEraBsSubscription({
        skip: selectedTab !== 1,
        onSubscriptionData: res => {
            setErabDropRateMetrix(
                erabDropRateMetrix.map(item => {
                    return {
                        name: item.name,
                        data: [
                            ...item.data,
                            ...(res.subscriptionData.data?.getMetricsERAB ||
                                []),
                        ],
                    };
                })
            );
        },
    });
    const [
        getMetricsRRC,
        { data: metricsRRCRes, refetch: metricsRRCResRefetch },
    ] = useGetMetricsRrcLazyQuery();

    useGetMetricsRrCsSubscription({
        skip: selectedTab !== 1,
        onSubscriptionData: res => {
            setRrcCnxSuccessMetrix(
                rrcCnxSuccessMetrix.map(item => {
                    return {
                        name: item.name,
                        data: [
                            ...item.data,
                            ...(res.subscriptionData.data?.getMetricsRRC || []),
                        ],
                    };
                })
            );
        },
    });
    const [
        getMetricThroughtpuUl,
        { data: metricThroughtputUlRes, refetch: metricThroughtputUlRefetch },
    ] = useGetMetricsThroughputUlLazyQuery();

    useGetMetricsThroughputUlsSubscription({
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

    useGetMetricsThroughputDlsSubscription({
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

    useGetMetricsCpuTrxsSubscription({
        skip: selectedTab !== 2,
        onSubscriptionData: res => {
            setCpuTrxMetric(
                cpuTrxMetric.map(item => {
                    return {
                        name: item.name,
                        data: [
                            ...item.data,
                            ...(res.subscriptionData.data?.getMetricsCpuTRX ||
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

    useGetMetricsMemoryTrxsSubscription({
        skip: selectedTab !== 2,
        onSubscriptionData: res => {
            setMemoryTrxMetric(
                memoryTrxMetric.map(item => {
                    return {
                        name: item.name,
                        data: [
                            ...item.data,
                            ...(res.subscriptionData.data
                                ?.getMetricsMemoryTRX || []),
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

    useGetMetricsUptimeSSubscription({
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

    const [
        getMetricPower,
        { data: metricPowerRes, refetch: metricPowerRefetch },
    ] = useGetMetricsPowerLazyQuery();

    useGetMetricsPowerSSubscription({
        skip: selectedTab !== 2,
        onSubscriptionData: res => {
            setPowerMetric(
                powerMetric.map(item => {
                    return {
                        name: item.name,
                        data: [
                            ...item.data,
                            ...(res.subscriptionData.data?.getMetricsPower ||
                                []),
                        ],
                    };
                })
            );
        },
    });

    const [
        getMetricTempTrx,
        { data: metricTempTrxRes, refetch: metricTempTrxRefetch },
    ] = useGetMetricsTempTrxLazyQuery();

    useGetMetricsTempTrxsSubscription({
        skip: selectedTab !== 0,
        onSubscriptionData: res => {
            setTempTrxMetric(
                tempTrxMetric.map(item => {
                    return {
                        name: item.name,
                        data: [
                            ...item.data,
                            ...(res.subscriptionData.data?.getMetricsTempTRX ||
                                []),
                        ],
                    };
                })
            );
        },
    });

    const [
        getMetricTempCom,
        { data: metricTempComRes, refetch: metricTempComRefetch },
    ] = useGetMetricsTempComLazyQuery();

    useGetMetricsTempComsSubscription({
        skip: selectedTab !== 0,
        onSubscriptionData: res => {
            setTempComMetric(
                tempComMetric.map(item => {
                    return {
                        name: item.name,
                        data: [
                            ...item.data,
                            ...(res.subscriptionData.data?.getMetricsTempCOM ||
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
                    ...getFirstMetricCallPayload(),
                },
            });
            getMetricsRLC({
                variables: {
                    ...getFirstMetricCallPayload(),
                },
            });
            getMetricsRRC({
                variables: {
                    ...getFirstMetricCallPayload(),
                },
            });
            getMetricTempTrx({
                variables: {
                    ...getFirstMetricCallPayload(),
                },
            });
            getMetricTempCom({
                variables: {
                    ...getFirstMetricCallPayload(),
                },
            });
            getMetricsERAB({
                variables: {
                    ...getFirstMetricCallPayload(),
                },
            });
        } else if (selectedTab === 1) {
            getMetricThroughtpuUl({
                variables: {
                    ...getFirstMetricCallPayload(),
                },
            });
            getMetricThroughtpuDl({
                variables: {
                    ...getFirstMetricCallPayload(),
                },
            });
        } else if (selectedTab === 2) {
            getMetricCpuTrx({
                variables: {
                    ...getFirstMetricCallPayload(),
                },
            });
            getMetricMemoryTrx({
                variables: {
                    ...getFirstMetricCallPayload(),
                },
            });
            getMetricPower({
                variables: {
                    ...getFirstMetricCallPayload(),
                },
            });
        }
    }, [selectedTab, selectedNode]);

    useEffect(() => {
        if (
            selectedTab === 0 &&
            metricsERABres &&
            metricsERABres.getMetricsERAB.length > 0
        ) {
            if (!isMetricData(uptimeMetric)) {
                setErabDropRateMetrix(
                    erabDropRateMetrix.map(item => {
                        return {
                            name: item.name,
                            data: [
                                ...item.data,
                                ...(metricsERABres.getMetricsERAB || []),
                            ],
                        };
                    })
                );
            }
            metricsERABresRefetch({
                ...getMetricPollingCallPayload(
                    metricsERABres.getMetricsERAB[
                        metricsERABres.getMetricsERAB.length - 1
                    ].x
                ),
            });
        }
    }, [metricsERABres]);
    useEffect(() => {
        if (
            selectedTab === 0 &&
            metricUptimeTrxRes &&
            metricUptimeTrxRes.getMetricsUptime.length > 0
        ) {
            if (!isMetricData(uptimeMetric)) {
                setUptimeMetrics(
                    uptimeMetric.map(item => {
                        return {
                            name: item.name,
                            data: [
                                ...item.data,
                                ...(metricUptimeTrxRes.getMetricsUptime || []),
                            ],
                        };
                    })
                );
            }
            metricUptimeRefetch({
                ...getMetricPollingCallPayload(
                    metricUptimeTrxRes.getMetricsUptime[
                        metricUptimeTrxRes.getMetricsUptime.length - 1
                    ].x
                ),
            });
        }
    }, [metricUptimeTrxRes]);

    useEffect(() => {
        if (
            selectedTab === 0 &&
            metricTempTrxRes &&
            metricTempTrxRes.getMetricsTempTRX.length > 0
        ) {
            if (!isMetricData(tempTrxMetric)) {
                setTempTrxMetric(
                    tempTrxMetric.map(item => {
                        return {
                            name: item.name,
                            data: [
                                ...item.data,
                                ...(metricTempTrxRes.getMetricsTempTRX || []),
                            ],
                        };
                    })
                );
            }
            metricTempTrxRefetch({
                ...getMetricPollingCallPayload(
                    metricTempTrxRes.getMetricsTempTRX[
                        metricTempTrxRes.getMetricsTempTRX.length - 1
                    ].x
                ),
            });
        }
    }, [metricTempTrxRes]);

    useEffect(() => {
        if (
            selectedTab === 0 &&
            metricsRLCres &&
            metricsRLCres.getMetricsRLC.length > 0
        ) {
            if (!isMetricData(rlsDropRateMetrics)) {
                setRlsDropRateMetrics(
                    rlsDropRateMetrics.map(item => {
                        return {
                            name: item.name,
                            data: [
                                ...item.data,
                                ...(metricsRLCres.getMetricsRLC || []),
                            ],
                        };
                    })
                );
            }
            metricsRlcRefetch({
                ...getMetricPollingCallPayload(
                    metricsRLCres.getMetricsRLC[
                        metricsRLCres.getMetricsRLC.length - 1
                    ].x
                ),
            });
        }
    }, [metricsRLCres]);
    useEffect(() => {
        if (
            selectedTab === 0 &&
            metricTempComRes &&
            metricTempComRes.getMetricsTempCOM.length > 0
        ) {
            if (!isMetricData(tempComMetric)) {
                setTempComMetric(
                    tempComMetric.map(item => {
                        return {
                            name: item.name,
                            data: [
                                ...item.data,
                                ...(metricTempComRes.getMetricsTempCOM || []),
                            ],
                        };
                    })
                );
            }
            metricTempComRefetch({
                ...getMetricPollingCallPayload(
                    metricTempComRes.getMetricsTempCOM[
                        metricTempComRes.getMetricsTempCOM.length - 1
                    ].x
                ),
            });
        }
    }, [metricTempComRes]);

    useEffect(() => {
        if (
            selectedTab === 1 &&
            metricThroughtputUlRes &&
            metricThroughtputUlRes.getMetricsThroughputUL.length > 0
        ) {
            if (!isMetricData(throughputULMetric)) {
                setThroughputULMetric(
                    throughputULMetric.map(item => {
                        return {
                            name: item.name,
                            data: [
                                ...item.data,
                                ...(metricThroughtputUlRes.getMetricsThroughputUL ||
                                    []),
                            ],
                        };
                    })
                );
            }
            metricThroughtputUlRefetch({
                ...getMetricPollingCallPayload(
                    metricThroughtputUlRes.getMetricsThroughputUL[
                        metricThroughtputUlRes.getMetricsThroughputUL.length - 1
                    ].x
                ),
            });
        }
    }, [metricThroughtputUlRes]);

    useEffect(() => {
        if (
            selectedTab === 1 &&
            metricThroughtputDlRes &&
            metricThroughtputDlRes.getMetricsThroughputDL.length > 0
        ) {
            if (!isMetricData(throughputDLMetric)) {
                setThroughputDLMetric(
                    throughputDLMetric.map(item => {
                        return {
                            name: item.name,
                            data: [
                                ...item.data,
                                ...(metricThroughtputDlRes.getMetricsThroughputDL ||
                                    []),
                            ],
                        };
                    })
                );
            }
            metricThroughtputDlRefetch({
                ...getMetricPollingCallPayload(
                    metricThroughtputDlRes.getMetricsThroughputDL[
                        metricThroughtputDlRes.getMetricsThroughputDL.length - 1
                    ].x
                ),
            });
        }
    }, [metricThroughtputDlRes]);

    useEffect(() => {
        if (
            selectedTab === 2 &&
            metricCpuTrxRes &&
            metricCpuTrxRes.getMetricsCpuTRX.length > 0
        ) {
            if (!isMetricData(cpuTrxMetric)) {
                setCpuTrxMetric(
                    cpuTrxMetric.map(item => {
                        return {
                            name: item.name,
                            data: [
                                ...item.data,
                                ...(metricCpuTrxRes.getMetricsCpuTRX || []),
                            ],
                        };
                    })
                );
            }
            metricCpuTrxRefetch({
                ...getMetricPollingCallPayload(
                    metricCpuTrxRes.getMetricsCpuTRX[
                        metricCpuTrxRes.getMetricsCpuTRX.length - 1
                    ].x
                ),
            });
        }
    }, [metricCpuTrxRes]);

    useEffect(() => {
        if (
            selectedTab === 2 &&
            metricMemoryTrxRes &&
            metricMemoryTrxRes.getMetricsMemoryTRX.length > 0
        ) {
            if (!isMetricData(memoryTrxMetric)) {
                setMemoryTrxMetric(
                    memoryTrxMetric.map(item => {
                        return {
                            name: item.name,
                            data: [
                                ...item.data,
                                ...(metricMemoryTrxRes.getMetricsMemoryTRX ||
                                    []),
                            ],
                        };
                    })
                );
            }
            metricMemoryTrxRefetch({
                ...getMetricPollingCallPayload(
                    metricMemoryTrxRes.getMetricsMemoryTRX[
                        metricMemoryTrxRes.getMetricsMemoryTRX.length - 1
                    ].x
                ),
            });
        }
    }, [metricMemoryTrxRes]);

    useEffect(() => {
        if (
            selectedTab === 2 &&
            metricPowerRes &&
            metricPowerRes.getMetricsPower.length > 0
        ) {
            if (!isMetricData(powerMetric)) {
                setPowerMetric(
                    powerMetric.map(item => {
                        return {
                            name: item.name,
                            data: [
                                ...item.data,
                                ...(metricPowerRes.getMetricsPower || []),
                            ],
                        };
                    })
                );
            }
            metricPowerRefetch({
                ...getMetricPollingCallPayload(
                    metricPowerRes.getMetricsPower[
                        metricPowerRes.getMetricsPower.length - 1
                    ].x
                ),
            });
        }
    }, [metricPowerRes]);

    useEffect(() => {
        if (
            selectedTab === 0 &&
            metricsRRCRes &&
            metricsRRCRes.getMetricsRRC.length > 0
        ) {
            if (!isMetricData(rrcCnxSuccessMetrix)) {
                setRrcCnxSuccessMetrix(
                    rrcCnxSuccessMetrix.map(item => {
                        return {
                            name: item.name,
                            data: [
                                ...item.data,
                                ...(metricsRRCRes.getMetricsRRC || []),
                            ],
                        };
                    })
                );
            }
            metricsRRCResRefetch({
                ...getMetricPollingCallPayload(
                    metricsRRCRes.getMetricsRRC[
                        metricsRRCRes.getMetricsRRC.length - 1
                    ].x
                ),
            });
        }
    }, [metricsRRCRes]);
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
                                isUpdateAvailable={true}
                                selectedNode={selectedNode}
                                uptimeMetrics={uptimeMetric}
                                tempTrxMetric={tempTrxMetric}
                                tempComMetric={tempComMetric}
                                handleUpdateNode={handleUpdateNode}
                                getNodeSoftwareUpdateInfos={handleSoftwareInfos}
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
                                rrcCnxSuccessMetrix={rrcCnxSuccessMetrix}
                                erabDropRateMetrix={erabDropRateMetrix}
                                rlsDropRateMetrics={rlsDropRateMetrics}
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
