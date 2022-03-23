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
    ActivationDialog,
} from "../../components";
import { NodePageTabs, NODE_ACTIONS } from "../../constants";
import {
    NodeDto,
    MetricDto,
    useGetNodeDetailsQuery,
    useGetNodeAppsQuery,
    useGetMetricsThroughputUlLazyQuery,
    useGetMetricsSubAttachedLazyQuery,
    useGetMetricsSubAttachedsSubscription,
    useGetMetricsThroughputDlLazyQuery,
    useGetMetricsThroughputUlsSubscription,
    useGetMetricsThroughputDlsSubscription,
    useGetMetricsCpuTrxLazyQuery,
    useGetMetricsEraBsSubscription,
    useGetMetricsRlCsSubscription,
    useGetMetricsCpuCoMsSubscription,
    useGetMetricsCpuComLazyQuery,
    useGetMetricsErabLazyQuery,
    useGetMetricsRrCsSubscription,
    useGetMetricsDiskCoMsSubscription,
    useGetMetricsDiskComLazyQuery,
    useGetMetricsRlcLazyQuery,
    useGetMetricsMemoryComsSubscription,
    useGetMetricsMemoryComLazyQuery,
    useGetMetricsMemoryTrxLazyQuery,
    useGetMetricsUptimeLazyQuery,
    useGetMetricsUptimeSSubscription,
    useGetMetricsMemoryTrxsSubscription,
    useGetMetricsCpuTrxsSubscription,
    useGetNodeAppsVersionLogsQuery,
    useGetMetricsPowerLazyQuery,
    useGetMetricsPowerSSubscription,
    useGetMetricsTempTrxLazyQuery,
    useGetMetricsTempTrxsSubscription,
    useGetMetricsTempComLazyQuery,
    useGetMetricsTempComsSubscription,
    useGetMetricsRrcLazyQuery,
    useGetMetricsDiskTrXsSubscription,
    useGetMetricsDiskTrxLazyQuery,
    useGetMetricsSubActiveLazyQuery,
    useGetMetricsSubActivesSubscription,
    useGetNodesByOrgLazyQuery,
    useGetMetricsTxPowerLazyQuery,
    useGetMetricsTxPowersSubscription,
    useGetMetricsPaPowerLazyQuery,
    useGetMetricsRxPowersSubscription,
    useGetMetricsPaPowersSubscription,
    Org_Node_State,
    useGetMetricsRxPowerLazyQuery,
} from "../../generated";
import { useRecoilValue } from "recoil";
import { isSkeltonLoading, user } from "../../recoil";
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
    const { id: orgId = "" } = useRecoilValue(user);
    const [selectedTab, setSelectedTab] = useState(0);
    const [isAddNode, setIsAddNode] = useState(false);
    const skeltonLoading = useRecoilValue(isSkeltonLoading);
    const [nodeAppDetails, setNodeAppDetails] = useState<any>();
    const [selectedNode, setSelectedNode] = useState<NodeDto | undefined>({
        id: "",
        type: "",
        title: "",
        totalUser: 0,
        description: "",
        updateVersion: "",
        updateShortNote: "",
        updateDescription: "",
        isUpdateAvailable: false,
        status: Org_Node_State.Undefined,
    });
    const [showNodeAppDialog, setShowNodeAppDialog] = useState(false);
    const [rrcCnxSuccessMetrix, setRrcCnxSuccessMetrix] = useState<
        {
            name: string;
            data: MetricDto[];
        }[]
    >(getDefaultList(["RRC CNX Success"]));
    const [diskComMetrics, setDiskComMetrics] = useState<
        {
            name: string;
            data: MetricDto[];
        }[]
    >(getDefaultList(["DISK-COM"]));
    const [paPowerMetrics, setPaPowerMetrics] = useState<
        {
            name: string;
            data: MetricDto[];
        }[]
    >(getDefaultList(["PA POWER"]));
    const [diskTrxMatrics, setDiskTrxMatrics] = useState<
        {
            name: string;
            data: MetricDto[];
        }[]
    >(getDefaultList(["DISK-COM"]));
    const [erabDropRateMetrix, setErabDropRateMetrix] = useState<
        {
            name: string;
            data: MetricDto[];
        }[]
    >(getDefaultList(["ERAB Drop Rate"]));

    const [cpuComMetrics, setCpuComMetrics] = useState<
        {
            name: string;
            data: MetricDto[];
        }[]
    >(getDefaultList(["CPU-COM"]));
    const [rxPowerMetrics, setRxPowerMetrics] = useState<
        {
            name: string;
            data: MetricDto[];
        }[]
    >(getDefaultList(["RX POWER"]));
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
    >(getDefaultList(["UPTIME"]));
    const [cpuTrxMetric, setCpuTrxMetric] = useState<
        {
            name: string;
            data: MetricDto[];
        }[]
    >(getDefaultList(["CPU (TRX)"]));
    const [attachedSubcriberMetrics, setAttachedSubcriberMetrics] = useState<
        {
            name: string;
            data: MetricDto[];
        }[]
    >(getDefaultList(["Attached"]));
    const [txPowerMetrics, setTxPowerMetrics] = useState<
        {
            name: string;
            data: MetricDto[];
        }[]
    >(getDefaultList(["TX Power"]));
    const [activeSubcriberMetrics, setActiveSubcriberMetrics] = useState<
        {
            name: string;
            data: MetricDto[];
        }[]
    >(getDefaultList(["Active"]));
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
    const [memoryComMetrics, setMemoryComMetrics] = useState<
        {
            name: string;
            data: MetricDto[];
        }[]
    >(getDefaultList(["MEMORY (COM)"]));
    const [showNodeSoftwareUpdatInfos, setShowNodeSoftwareUpdatInfos] =
        useState<boolean>(false);

    const getFirstMetricCallPayload = () =>
        getMetricPayload({
            nodeId: selectedNode?.id,
            orgId: orgId,
            regPolling: false,
            to: Math.floor(Date.now() / 1000) - 10,
            from: Math.floor(Date.now() / 1000) - 180,
        });

    const getMetricPollingCallPayload = (from: number) =>
        getMetricPayload({
            nodeId: selectedNode?.id,
            orgId: orgId,
            from: from + 1,
        });

    const [getNodesByOrg, { data: nodesRes, loading: nodesLoading }] =
        useGetNodesByOrgLazyQuery({
            onCompleted: res => {
                res?.getNodesByOrg?.nodes.length > 0 &&
                    setSelectedNode(res?.getNodesByOrg?.nodes[0]);
            },
        });

    const { data: nodeDetailRes, loading: nodeDetailLoading } =
        useGetNodeDetailsQuery();
    const [
        getMetricsCpuCOM,
        { data: cpuComMetricsRes, refetch: cpuComMetricsRefetch },
    ] = useGetMetricsCpuComLazyQuery();

    useGetMetricsCpuCoMsSubscription({
        skip: selectedTab !== 1,
        onSubscriptionData: res => {
            setCpuComMetrics(
                cpuComMetrics.map(item => {
                    return {
                        name: item.name,
                        data: [
                            ...item.data,
                            ...(res.subscriptionData.data?.getMetricsCpuCOM ||
                                []),
                        ],
                    };
                })
            );
        },
    });
    const [
        getMetricsRxPower,
        { data: rxPowerMetricsRes, refetch: rxPowerMetricsRefetch },
    ] = useGetMetricsRxPowerLazyQuery();

    useGetMetricsRxPowersSubscription({
        skip: selectedTab !== 3,
        onSubscriptionData: res => {
            setRxPowerMetrics(
                rxPowerMetrics.map(item => {
                    return {
                        name: item.name,
                        data: [
                            ...item.data,
                            ...(res.subscriptionData.data?.getMetricsRxPower ||
                                []),
                        ],
                    };
                })
            );
        },
    });
    const [
        getMetricsPaPower,
        { data: paPowerMetricsRes, refetch: paPowerMetricsRefetch },
    ] = useGetMetricsPaPowerLazyQuery();

    useGetMetricsPaPowersSubscription({
        skip: selectedTab !== 3,
        onSubscriptionData: res => {
            setPaPowerMetrics(
                paPowerMetrics.map(item => {
                    return {
                        name: item.name,
                        data: [
                            ...item.data,
                            ...(res.subscriptionData.data?.getMetricsPaPower ||
                                []),
                        ],
                    };
                })
            );
        },
    });
    const [
        getMetricsTxPower,
        { data: txPowerMetricsRes, refetch: txPowerMetricsRefetch },
    ] = useGetMetricsTxPowerLazyQuery();

    useGetMetricsTxPowersSubscription({
        skip: selectedTab !== 3,
        onSubscriptionData: res => {
            setTxPowerMetrics(
                txPowerMetrics.map(item => {
                    return {
                        name: item.name,
                        data: [
                            ...item.data,
                            ...(res.subscriptionData.data?.getMetricsTxPower ||
                                []),
                        ],
                    };
                })
            );
        },
    });
    const [
        getMetricsDiskTRX,
        { data: diskTrxMetricsRes, refetch: diskTrxMetricsResRefetch },
    ] = useGetMetricsDiskTrxLazyQuery();

    useGetMetricsDiskTrXsSubscription({
        skip: selectedTab !== 2,
        onSubscriptionData: res => {
            setDiskTrxMatrics(
                diskTrxMatrics.map(item => {
                    return {
                        name: item.name,
                        data: [
                            ...item.data,
                            ...(res.subscriptionData.data?.getMetricsDiskTRX ||
                                []),
                        ],
                    };
                })
            );
        },
    });
    const [
        getMetricsMemoryCOM,
        { data: metricsMemoryComres, refetch: metricsMemoryComRefetch },
    ] = useGetMetricsMemoryComLazyQuery();

    useGetMetricsMemoryComsSubscription({
        skip: selectedTab !== 1,
        onSubscriptionData: res => {
            setMemoryComMetrics(
                memoryComMetrics.map(item => {
                    return {
                        name: item.name,
                        data: [
                            ...item.data,
                            ...(res.subscriptionData.data
                                ?.getMetricsMemoryCOM || []),
                        ],
                    };
                })
            );
        },
    });

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
    const { data: nodeAppsRes, loading: nodeAppsLoading } =
        useGetNodeAppsQuery();
    const { data: nodeAppsLogsRes, loading: nodeAppsLogsLoading } =
        useGetNodeAppsVersionLogsQuery();
    const [
        getMetricsSubAttached,
        {
            data: subscriberAttachedmetricRes,
            refetch: subscriberAttachedmetricRefetch,
        },
    ] = useGetMetricsSubAttachedLazyQuery();

    useGetMetricsSubAttachedsSubscription({
        skip: selectedTab !== 0,
        onSubscriptionData: res => {
            setAttachedSubcriberMetrics(
                attachedSubcriberMetrics.map(item => {
                    return {
                        name: item.name,
                        data: [
                            ...item.data,
                            ...(res.subscriptionData.data
                                ?.getMetricsSubAttached || []),
                        ],
                    };
                })
            );
        },
    });
    const [
        getMetricsSubActive,
        {
            data: subscriberActiveMetricRes,
            refetch: subscriberActiveMetricRefetch,
        },
    ] = useGetMetricsSubActiveLazyQuery();

    useGetMetricsSubActivesSubscription({
        skip: selectedTab !== 0,
        onSubscriptionData: res => {
            setActiveSubcriberMetrics(
                activeSubcriberMetrics.map(item => {
                    return {
                        name: item.name,
                        data: [
                            ...item.data,
                            ...(res.subscriptionData.data
                                ?.getMetricsSubActive || []),
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
        getMetricsDiskCOM,
        { data: metricsDiskComRes, refetch: metricsDiskComRefetch },
    ] = useGetMetricsDiskComLazyQuery();

    useGetMetricsDiskCoMsSubscription({
        skip: selectedTab !== 2,
        onSubscriptionData: res => {
            setDiskComMetrics(
                diskComMetrics.map(item => {
                    return {
                        name: item.name,
                        data: [
                            ...item.data,
                            ...(res.subscriptionData.data?.getMetricsDiskCOM ||
                                []),
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
        getNodesByOrg({ variables: { orgId: orgId } });
    }, []);

    useEffect(() => {
        if (selectedTab === 0) {
            getMetricUptime({
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
            getMetricsSubAttached({
                variables: {
                    ...getFirstMetricCallPayload(),
                },
            });
            getMetricsSubActive({
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
            getMetricsERAB({
                variables: {
                    ...getFirstMetricCallPayload(),
                },
            });
            getMetricsRLC({
                variables: {
                    ...getFirstMetricCallPayload(),
                },
            });
            getMetricThroughtpuDl({
                variables: {
                    ...getFirstMetricCallPayload(),
                },
            });
            getMetricsRRC({
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
            getMetricsDiskCOM({
                variables: {
                    ...getFirstMetricCallPayload(),
                },
            });
            getMetricsDiskTRX({
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
            getMetricsCpuCOM({
                variables: {
                    ...getFirstMetricCallPayload(),
                },
            });

            getMetricsMemoryCOM({
                variables: {
                    ...getFirstMetricCallPayload(),
                },
            });
        } else if (selectedTab === 3) {
            getMetricsPaPower({
                variables: {
                    ...getFirstMetricCallPayload(),
                },
            });
            getMetricsTxPower({
                variables: {
                    ...getFirstMetricCallPayload(),
                },
            });
            getMetricsRxPower({
                variables: {
                    ...getFirstMetricCallPayload(),
                },
            });
        }
    }, [selectedTab, selectedNode]);

    useEffect(() => {
        if (
            selectedTab == 3 &&
            paPowerMetricsRes &&
            paPowerMetricsRes.getMetricsPaPower.length > 0
        ) {
            if (!isMetricData(paPowerMetrics)) {
                setPaPowerMetrics(
                    paPowerMetrics.map(item => {
                        return {
                            name: item.name,
                            data: [
                                ...item.data,
                                ...(paPowerMetricsRes.getMetricsPaPower || []),
                            ],
                        };
                    })
                );
            }
            paPowerMetricsRefetch({
                ...getMetricPollingCallPayload(
                    paPowerMetricsRes.getMetricsPaPower[
                        paPowerMetricsRes.getMetricsPaPower.length - 1
                    ].x
                ),
            });
        }
    }, [paPowerMetricsRes]);
    useEffect(() => {
        if (
            selectedTab == 2 &&
            txPowerMetricsRes &&
            txPowerMetricsRes.getMetricsTxPower.length > 0
        ) {
            if (!isMetricData(txPowerMetrics)) {
                setTxPowerMetrics(
                    txPowerMetrics.map(item => {
                        return {
                            name: item.name,
                            data: [
                                ...item.data,
                                ...(txPowerMetricsRes.getMetricsTxPower || []),
                            ],
                        };
                    })
                );
            }
            txPowerMetricsRefetch({
                ...getMetricPollingCallPayload(
                    txPowerMetricsRes.getMetricsTxPower[
                        txPowerMetricsRes.getMetricsTxPower.length - 1
                    ].x
                ),
            });
        }
    }, [txPowerMetricsRes]);
    useEffect(() => {
        if (
            selectedTab == 3 &&
            rxPowerMetricsRes &&
            rxPowerMetricsRes.getMetricsRxPower.length > 0
        ) {
            if (!isMetricData(rxPowerMetrics)) {
                setRxPowerMetrics(
                    rxPowerMetrics.map(item => {
                        return {
                            name: item.name,
                            data: [
                                ...item.data,
                                ...(rxPowerMetricsRes.getMetricsRxPower || []),
                            ],
                        };
                    })
                );
            }
            rxPowerMetricsRefetch({
                ...getMetricPollingCallPayload(
                    rxPowerMetricsRes.getMetricsRxPower[
                        rxPowerMetricsRes.getMetricsRxPower.length - 1
                    ].x
                ),
            });
        }
    }, [rxPowerMetricsRes]);

    useEffect(() => {
        if (
            selectedTab == 2 &&
            diskTrxMetricsRes &&
            diskTrxMetricsRes.getMetricsDiskTRX.length > 0
        ) {
            if (!isMetricData(diskTrxMatrics)) {
                setDiskTrxMatrics(
                    diskTrxMatrics.map(item => {
                        return {
                            name: item.name,
                            data: [
                                ...item.data,
                                ...(diskTrxMetricsRes.getMetricsDiskTRX || []),
                            ],
                        };
                    })
                );
            }
            diskTrxMetricsResRefetch({
                ...getMetricPollingCallPayload(
                    diskTrxMetricsRes.getMetricsDiskTRX[
                        diskTrxMetricsRes.getMetricsDiskTRX.length - 1
                    ].x
                ),
            });
        }
    }, [diskTrxMetricsRes]);
    useEffect(() => {
        if (
            selectedTab == 1 &&
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
            selectedTab == 2 &&
            metricsDiskComRes &&
            metricsDiskComRes.getMetricsDiskCOM.length > 0
        ) {
            if (!isMetricData(cpuComMetrics)) {
                setDiskComMetrics(
                    diskComMetrics.map(item => {
                        return {
                            name: item.name,
                            data: [
                                ...item.data,
                                ...(metricsDiskComRes.getMetricsDiskCOM || []),
                            ],
                        };
                    })
                );
            }
            metricsDiskComRefetch({
                ...getMetricPollingCallPayload(
                    metricsDiskComRes.getMetricsDiskCOM[
                        metricsDiskComRes.getMetricsDiskCOM.length - 1
                    ].x
                ),
            });
        }
    }, [metricsDiskComRes]);
    useEffect(() => {
        if (
            selectedTab !== 2 &&
            cpuComMetricsRes &&
            cpuComMetricsRes.getMetricsCpuCOM.length > 0
        ) {
            if (!isMetricData(cpuComMetrics)) {
                setCpuComMetrics(
                    cpuComMetrics.map(item => {
                        return {
                            name: item.name,
                            data: [
                                ...item.data,
                                ...(cpuComMetricsRes.getMetricsCpuCOM || []),
                            ],
                        };
                    })
                );
            }
            cpuComMetricsRefetch({
                ...getMetricPollingCallPayload(
                    cpuComMetricsRes.getMetricsCpuCOM[
                        cpuComMetricsRes.getMetricsCpuCOM.length - 1
                    ].x
                ),
            });
        }
    }, [cpuComMetricsRes]);
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
            selectedTab !== 2 &&
            metricsMemoryComres &&
            metricsMemoryComres.getMetricsMemoryCOM.length > 0
        ) {
            if (!isMetricData(memoryComMetrics)) {
                setMemoryComMetrics(
                    memoryComMetrics.map(item => {
                        return {
                            name: item.name,
                            data: [
                                ...item.data,
                                ...(metricsMemoryComres.getMetricsMemoryCOM ||
                                    []),
                            ],
                        };
                    })
                );
            }
            metricsMemoryComRefetch({
                ...getMetricPollingCallPayload(
                    metricsMemoryComres.getMetricsMemoryCOM[
                        metricsMemoryComres.getMetricsMemoryCOM.length - 1
                    ].x
                ),
            });
        }
    }, [metricsMemoryComres]);
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
            selectedTab == 1 &&
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
            selectedTab === 0 &&
            subscriberAttachedmetricRes &&
            subscriberAttachedmetricRes.getMetricsSubAttached.length > 0
        ) {
            if (!isMetricData(attachedSubcriberMetrics)) {
                setAttachedSubcriberMetrics(
                    attachedSubcriberMetrics.map(item => {
                        return {
                            name: item.name,
                            data: [
                                ...item.data,
                                ...(subscriberAttachedmetricRes.getMetricsSubAttached ||
                                    []),
                            ],
                        };
                    })
                );
            }
            subscriberAttachedmetricRefetch({
                ...getMetricPollingCallPayload(
                    subscriberAttachedmetricRes.getMetricsSubAttached[
                        subscriberAttachedmetricRes.getMetricsSubAttached
                            .length - 1
                    ].x
                ),
            });
        }
    }, [subscriberAttachedmetricRes]);
    useEffect(() => {
        if (
            selectedTab === 0 &&
            subscriberActiveMetricRes &&
            subscriberActiveMetricRes.getMetricsSubActive.length > 0
        ) {
            if (!isMetricData(activeSubcriberMetrics)) {
                setActiveSubcriberMetrics(
                    attachedSubcriberMetrics.map(item => {
                        return {
                            name: item.name,
                            data: [
                                ...item.data,
                                ...(subscriberActiveMetricRes.getMetricsSubActive ||
                                    []),
                            ],
                        };
                    })
                );
            }
            subscriberActiveMetricRefetch({
                ...getMetricPollingCallPayload(
                    subscriberActiveMetricRes.getMetricsSubActive[
                        subscriberActiveMetricRes.getMetricsSubActive.length - 1
                    ].x
                ),
            });
        }
    }, [subscriberActiveMetricRes]);
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
            selectedTab == 1 &&
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
        //TODO: Handle NODE action button click
    };
    const onAddNode = () => {
        setIsAddNode(true);
    };
    const handleUpdateNode = () => {
        // TODO: Handle Update node Action
    };

    const handleAddNodeClose = () => setIsAddNode(() => false);

    const handleActivationSubmit = () => {
        /* Handle submit activation action */
    };

    const getNodeAppDetails = (id: any) => {
        setShowNodeAppDialog(true);
        nodeAppsRes?.getNodeApps
            .filter(nodeApp => nodeApp.id == id)
            .map(filteredNodeApp => setNodeAppDetails(filteredNodeApp));
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
            {(nodesRes && nodesRes?.getNodesByOrg?.nodes.length > 0) ||
            isLoading ? (
                <Grid container spacing={3}>
                    <Grid item xs={12}>
                        <NodeStatus
                            onAddNode={onAddNode}
                            loading={isLoading || nodesLoading}
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
                                attachedSubcriberMetrics={
                                    attachedSubcriberMetrics
                                }
                                activeSubcriberMetrics={activeSubcriberMetrics}
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
                                memoryComMetrics={memoryComMetrics}
                                cpuComMetrics={cpuComMetrics}
                                diskTrxMatrics={diskTrxMatrics}
                                diskComMetrics={diskComMetrics}
                                powerMetrics={powerMetric}
                            />
                        </TabPanel>
                        <TabPanel
                            id={"node-tab-3"}
                            value={selectedTab}
                            index={3}
                        >
                            <NodeRadioTab
                                loading={isLoading || nodeDetailLoading}
                                txPowerMetrics={txPowerMetrics}
                                rxPowerMetrics={rxPowerMetrics}
                                paPowerMetrics={paPowerMetrics}
                            />
                        </TabPanel>
                        <TabPanel
                            id={"node-tab-4"}
                            value={selectedTab}
                            index={4}
                        >
                            <NodeSoftwareTab
                                loading={
                                    isLoading ||
                                    nodeDetailLoading ||
                                    nodeAppsLogsLoading ||
                                    nodeAppsLoading
                                }
                                nodeApps={nodeAppsRes?.getNodeApps}
                                NodeLogs={
                                    nodeAppsLogsRes?.getNodeAppsVersionLogs
                                }
                                getNodeAppDetails={getNodeAppDetails}
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
                nodeData={nodeAppDetails}
                handleClose={handleNodAppDetailsDialog}
            />
            <NodeSoftwareInfosDialog
                closeBtnLabel="close"
                isOpen={showNodeSoftwareUpdatInfos}
                handleClose={handleCloseNodeInfos}
            />
            <ActivationDialog
                isOpen={isAddNode}
                dialogTitle={"Register Node"}
                handleClose={handleAddNodeClose}
                handleActivationSubmit={handleActivationSubmit}
                subTitle={
                    "Ensure node is properly set up in desired location before completing this step. Enter serial number found in your confirmation email, or on the back of your node, and we’ll take care of the rest for you."
                }
            />
        </Box>
    );
};

export default Nodes;
