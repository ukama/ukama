import {
    TabPanel,
    NodeDialog,
    NodeStatus,
    NodeRadioTab,
    LoadingWrapper,
    NodeNetworkTab,
    NodeSoftwareTab,
    NodeOverviewTab,
    PagePlaceholder,
    NodeResourcesTab,
    NodeSchematicTab,
    SoftwareUpdateModal,
    NodeAppDetailsDialog,
    NodeSoftwareInfosDialog,
} from "../../components";
import {
    NodeDto,
    Node_Type,
    Org_Node_State,
    useAddNodeMutation,
    useGetNodeAppsQuery,
    useGetNodeLazyQuery,
    useGetNodeStatusQuery,
    useGetNodesByOrgLazyQuery,
    useGetMetricsByTabLazyQuery,
    useGetNodeAppsVersionLogsQuery,
    useGetMetricsByTabSSubscription,
    Time_Filter,
    useGetConnectedUsersQuery,
} from "../../generated";
import {
    getMetricPayload,
    getMetricsInitObj,
    getMetricObjectByKey,
} from "../../utils";
import Fab from "@mui/material/Fab";
import { TMetric } from "../../types";
import { globalUseStyles } from "../../styles";
import { Box, Grid, Tab, Tabs } from "@mui/material";
import { SpecsDocsData } from "../../constants/stubData";
import ArrowBackIcon from "@mui/icons-material/ArrowBack";
import { useRecoilValue, useSetRecoilState } from "recoil";
import { NodePageTabs, NODE_ACTIONS } from "../../constants";
import React, { useCallback, useEffect, useState } from "react";
import { isSkeltonLoading, snackbarMessage } from "../../recoil";

let abortController = new AbortController();

const NODE_INIT = {
    type: "HOME",
    name: "",
    nodeId: "",
    orgId: "",
};

const Nodes = () => {
    const classes = globalUseStyles();
    const [selectedTab, setSelectedTab] = useState(0);
    const getFirstMetricCallPayload = (nodeId: string) =>
        getMetricPayload({
            tab: selectedTab,
            regPolling: false,
            nodeId: nodeId,
            to: Math.floor(Date.now() / 1000) - 10,
            from: Math.floor(Date.now() / 1000) - 180,
            nodeType: selectedNode?.type || Node_Type.Home,
        });

    const getMetricPollingCallPayload = (from: number) =>
        getMetricPayload({
            nodeId: selectedNode?.id,
            from: from + 1,
            tab: selectedTab,
            nodeType: selectedNode?.type || Node_Type.Home,
        });
    const [isAddNode, setIsAddNode] = useState(false);
    const skeltonLoading = useRecoilValue(isSkeltonLoading);
    const [nodeAppDetails, setNodeAppDetails] = useState<any>();
    const [isNodeUpdate, setIsNodeUpdate] = useState<boolean>(false);
    const [isSwitchOffRF, setIsSwitchOffRF] = useState<boolean>(false);
    const setRegisterNodeNotification = useSetRecoilState(snackbarMessage);
    const [isNodeRestart, setIsNodeRestart] = useState<boolean>(false);
    const [isSwitchOffNode, setIsSwitchOffNode] = useState<boolean>(false);
    const [towerNodeGroup, setTowerNodeGroup] = useState<NodeDto>();
    const [isTowerNode, setIsTowerNode] = useState<any>();
    const [selectedNode, setSelectedNode] = useState<NodeDto | undefined>({
        id: "",
        name: "",
        totalUser: 0,
        description: "",
        updateVersion: "",
        updateShortNote: "",
        type: Node_Type.Home,
        updateDescription: "",
        isUpdateAvailable: false,
        status: Org_Node_State.Undefined,
    });

    const [showNodeAppDialog, setShowNodeAppDialog] = useState(false);
    const [isMetricPolling, setIsMetricPolling] = useState(false);
    const [metrics, setMetrics] = useState<TMetric>(getMetricsInitObj());
    const [showNodeSoftwareUpdatInfos, setShowNodeSoftwareUpdatInfos] =
        useState<boolean>(false);
    const [backToPreviousNode, setBackToPreviousNode] =
        useState<boolean>(false);
    const { data: nodeAppsRes, loading: nodeAppsLoading } =
        useGetNodeAppsQuery();

    const { data: nodeAppsLogsRes, loading: nodeAppsLogsLoading } =
        useGetNodeAppsVersionLogsQuery();

    const [
        getNodesByOrg,
        {
            data: nodesRes,
            loading: nodesLoading,
            refetch: refetchGetNodesByOrg,
        },
    ] = useGetNodesByOrgLazyQuery({
        fetchPolicy: "cache-and-network",
    });

    const [
        registerNode,
        {
            loading: registerNodeLoading,
            data: registerNodeRes,
            error: registerNodeError,
        },
    ] = useAddNodeMutation({
        onCompleted: () => {
            setRegisterNodeNotification({
                id: "addNodeSuccess",
                message: `${registerNodeRes?.addNode?.name} has been registered successfully!`,
                type: "success",
                show: true,
            });
            refetchGetNodesByOrg();
        },

        onError: () =>
            setRegisterNodeNotification({
                id: "ErrorAddingNode",
                message: `${registerNodeError?.message}`,
                type: "error",
                show: true,
            }),
    });

    const [
        getNode,
        {
            data: getNodeData,
            loading: getNodeLoading,
            variables: getNodeVariables,
        },
    ] = useGetNodeLazyQuery();

    const { data: getNodeStatusData, loading: nodeStatusLoading } =
        useGetNodeStatusQuery({
            skip: !selectedNode?.id,
            variables: {
                data: {
                    nodeId: selectedNode?.id || "",
                    nodeType:
                        (selectedNode?.type as Node_Type) || Node_Type.Home,
                },
            },
        });

    const [
        getMetrics,
        {
            data: getMetricsRes,
            loading: metricsLoading,
            refetch: getMetricsRefetch,
            variables: lastMetricsFetchVariables,
        },
    ] = useGetMetricsByTabLazyQuery({
        context: {
            fetchOptions: {
                signal: abortController.signal,
            },
        },
        onCompleted: res => {
            if (res?.getMetricsByTab?.metrics.length > 0 && !isMetricPolling) {
                const _m: TMetric = getMetricsInitObj();
                setIsMetricPolling(true);

                for (const element of res.getMetricsByTab.metrics) {
                    _m[element.type] = {
                        name: element.name,
                        data: element.data,
                    };
                }

                setMetrics((_prev: TMetric) => ({ ..._m }));
            }
        },
        onError: () => {
            const obj: TMetric = getMetricsInitObj();
            Object.keys(obj).forEach(
                (k: string) =>
                    (obj[k as keyof TMetric] = getMetricObjectByKey(k))
            );
            setMetrics(() => ({ ...obj }));
        },
        fetchPolicy: "network-only",
    });

    const refetchMetrics = useCallback(() => {
        if (getMetricsRes && getMetricsRes.getMetricsByTab.to) {
            getMetricsRefetch({
                ...getMetricPollingCallPayload(
                    getMetricsRes?.getMetricsByTab.to
                ),
            });
        }
    }, [getMetricsRes]);

    useGetMetricsByTabSSubscription({
        fetchPolicy: "network-only",
        onSubscriptionData: res => {
            if (
                isMetricPolling &&
                res?.subscriptionData?.data?.getMetricsByTab &&
                res?.subscriptionData?.data?.getMetricsByTab.length > 0
            ) {
                const _m: TMetric = {};

                for (const element of res.subscriptionData.data
                    .getMetricsByTab) {
                    const metric = metrics[element.type];

                    if (
                        metric &&
                        metric.data &&
                        metric.data.length > 0 &&
                        element.data.length > 0 &&
                        element.data[element.data.length - 1].x >
                            metric.data[metric.data.length - 1].x
                    ) {
                        _m[element.type] = {
                            name: element.name,
                            data: [...(metric.data || []), ...element.data],
                        };
                    }
                }

                setMetrics((_prev: TMetric) => ({
                    ..._prev,
                    ..._m,
                }));

                let next = false;
                for (const element of res.subscriptionData.data
                    .getMetricsByTab) {
                    if (!next && element.next) next = true;
                }
                if (next) {
                    refetchMetrics();
                }
            }
        },
    });

    const { data: connectedUserRes } = useGetConnectedUsersQuery({
        variables: {
            filter: Time_Filter.Total,
        },
    });

    useEffect(() => {
        getNodesByOrg();
    }, []);

    useEffect(() => {
        if (
            !!selectedNode &&
            nodesRes?.getNodesByOrg &&
            nodesRes.getNodesByOrg.nodes.length > 0 &&
            !metricsLoading
        ) {
            if (nodesRes.getNodesByOrg.nodes[0].type === Node_Type.Tower) {
                getNode({
                    variables: { nodeId: nodesRes.getNodesByOrg.nodes[0].id },
                });
            }
            setSelectedNode(nodesRes.getNodesByOrg.nodes[0]);
            setMetrics(getMetricsInitObj());
            getMetrics({
                variables: {
                    ...getFirstMetricCallPayload(
                        nodesRes.getNodesByOrg.nodes[0].id || ""
                    ),
                },
            });
        }
    }, [nodesRes]);

    useEffect(() => {
        if (selectedNode && selectedNode.id && !metricsLoading) {
            if (
                selectedNode.type === Node_Type.Tower &&
                selectedNode.id !== getNodeVariables?.nodeId
            ) {
                getNode({
                    variables: { nodeId: selectedNode.id },
                });
            }
            abortController.abort();
            setTimeout(() => {
                setIsMetricPolling(false);
                abortController = new AbortController();
                setMetrics(getMetricsInitObj());
                getMetrics({
                    variables: {
                        ...getFirstMetricCallPayload(selectedNode?.id || ""),
                    },
                });
            }, 500);
        }
    }, [selectedNode, selectedTab]);

    useEffect(() => {
        if (
            isMetricPolling &&
            getMetricsRes &&
            getMetricsRes.getMetricsByTab.next &&
            !lastMetricsFetchVariables?.data.regPolling
        ) {
            getMetricsRefetch({
                ...getMetricPollingCallPayload(
                    getMetricsRes?.getMetricsByTab.to
                ),
            });
        }
    }, [isMetricPolling, getMetricsRes]);

    const onTabSelected = (event: React.SyntheticEvent, value: any) =>
        setSelectedTab(value);

    const onNodeSelected = (node: NodeDto) => {
        setSelectedNode(node);
        setIsTowerNode(selectedNode?.type);
    };

    const onNodeSelectedFromGroup = (id: string) => {
        if (id) {
            setSelectedNode(
                nodesRes?.getNodesByOrg?.nodes.find(ele => ele.id === id)
            );
            if (selectedNode?.type === "TOWER") {
                setTowerNodeGroup(selectedNode);

                setBackToPreviousNode(true);
            } else {
                setBackToPreviousNode(false);
            }
        }
    };
    useEffect(() => {
        if (isTowerNode === "AMPLIFIER" || isTowerNode === "HOME") {
            setBackToPreviousNode(false);
        }
    }, [isTowerNode]);
    const onUpdateNodeClick = () => {
        //TODO: Handle NODE RESTART ACTION
        setIsNodeUpdate(true);
    };
    const onRestartNode = () => {
        //TODO: Handle NODE RESTART ACTION
    };
    const handleNodeActionItemSelected = () => {
        //Todo :Handle nodeAction Itemselected
    };
    const onUpdateNode = () => {
        //Todo :Handle nodeAction update
        setIsNodeUpdate(true);
    };
    const handleCloseUpdateAllNode = () => {
        setIsNodeUpdate(false);
    };
    const onSwitchOffNode = () => {
        //Todo :Handle nodeAction Itemselected
    };
    const handleNodeActioOptionClicked = (nodeAction: any) => {
        if (nodeAction == "Turn node off") {
            setIsSwitchOffNode(true);
        } else if (nodeAction == "Turn RF off") {
            setIsSwitchOffRF(true);
        } else if (nodeAction == "Restart") {
            setIsNodeRestart(true);
        }
    };
    const onAddNode = () => {
        setIsAddNode(true);
    };
    const getSpecsSchematicSearch = () => {
        //GetSpecs search
    };
    const onSwitchOffRF = () => {
        // TODO: Handle Update node Action
    };
    const handleUpdateNode = () => {
        // TODO: Handle Update node Action
        setIsNodeUpdate(true);
    };
    const handleCloseTurnOffRF = () => {
        setIsSwitchOffRF(false);
    };
    const handleAddNodeClose = () => setIsAddNode(() => false);

    const handleActivationSubmit = (data: any) => {
        registerNode({
            variables: {
                data: {
                    name: data.name,
                    nodeId: data.nodeId,
                },
            },
        });
        setIsAddNode(() => registerNodeLoading);
    };
    const handleCloseNodeRestart = () => {
        setIsNodeRestart(false);
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

    const handleBackToSingleTowerNode = () => {
        setSelectedNode(towerNodeGroup);
        setBackToPreviousNode(false);
    };

    const handleSoftwareInfos = () => {
        setShowNodeSoftwareUpdatInfos(true);
    };
    const handleCloseTurnOffNode = () => {
        setIsSwitchOffNode(false);
    };
    const isLoading = skeltonLoading || nodesLoading;

    return (
        <Box
            component="div"
            sx={{
                p: 0,
                mt: 3,
                pb: 2,
                height: "100%",
            }}
        >
            {(nodesRes && nodesRes?.getNodesByOrg?.nodes.length > 0) ||
            isLoading ? (
                <Grid container spacing={3}>
                    <Grid item xs={12}>
                        <NodeStatus
                            onAddNode={onAddNode}
                            loading={
                                isLoading ||
                                nodesLoading ||
                                registerNodeLoading ||
                                skeltonLoading
                            }
                            handleNodeActionClick={handleNodeActioOptionClicked}
                            selectedNode={selectedNode}
                            onNodeActionItemSelected={
                                handleNodeActionItemSelected
                            }
                            onNodeSelected={onNodeSelected}
                            nodeActionOptions={NODE_ACTIONS}
                            onUpdateNodeClick={onUpdateNodeClick}
                            nodeStatus={getNodeStatusData?.getNodeStatus}
                            nodeStatusLoading={nodeStatusLoading}
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
                                metrics={metrics}
                                isUpdateAvailable={true}
                                selectedNode={selectedNode}
                                metricsLoading={metricsLoading}
                                nodeGroupLoading={getNodeLoading}
                                handleUpdateNode={handleUpdateNode}
                                nodeGroupData={getNodeData?.getNode}
                                connectedUsers={
                                    connectedUserRes?.getConnectedUsers
                                        .totalUser
                                }
                                onNodeSelected={onNodeSelectedFromGroup}
                                uptime={getNodeStatusData?.getNodeStatus.uptime}
                                getNodeSoftwareUpdateInfos={handleSoftwareInfos}
                                loading={
                                    isLoading || nodesLoading || !selectedNode
                                }
                            />
                        </TabPanel>
                        <TabPanel
                            id={"node-tab-1"}
                            value={selectedTab}
                            index={1}
                        >
                            <NodeNetworkTab
                                metrics={metrics}
                                loading={isLoading}
                            />
                        </TabPanel>
                        <TabPanel
                            id={"node-tab-2"}
                            value={selectedTab}
                            index={2}
                        >
                            <NodeResourcesTab
                                metrics={metrics}
                                selectedNode={selectedNode}
                                loading={isLoading}
                            />
                        </TabPanel>
                        <TabPanel
                            id={"node-tab-3"}
                            value={selectedTab}
                            index={3}
                        >
                            <NodeRadioTab
                                metrics={metrics}
                                loading={isLoading}
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
                            <NodeSchematicTab
                                getSearchValue={getSpecsSchematicSearch}
                                schematicsSpecsData={SpecsDocsData}
                                nodeTitle={selectedNode?.name}
                                loading={nodesLoading}
                            />
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
            <SoftwareUpdateModal
                submit={onSwitchOffNode}
                isOpen={isSwitchOffNode}
                handleClose={handleCloseTurnOffNode}
                btnLabel={"TURN NODE OFF"}
                title={"Continue Turning Node Off?"}
                content={`Continue turning node off? You will lose Ukama coverage where this node is located, but will still be able to connect to the network through roaming. `}
            />
            <SoftwareUpdateModal
                submit={onSwitchOffRF}
                isOpen={isSwitchOffRF}
                handleClose={handleCloseTurnOffRF}
                btnLabel={"TURN RF OFF"}
                title={"Continue Turning RF Off?"}
                content={`Continue turning RF off? You will lose Ukama coverage for a few minutes while it restarts, but will still be able to connect to the network through roaming.`}
            />
            <SoftwareUpdateModal
                submit={onRestartNode}
                isOpen={isNodeRestart}
                handleClose={handleCloseNodeRestart}
                btnLabel={"RESTART NODE"}
                title={"Continue Restarting Node?"}
                content={`Continue restarting node? You will lose Ukama coverage for a few minutes while it restarts, but will still be able to connect to the network through roaming. `}
            />
            <SoftwareUpdateModal
                submit={onUpdateNode}
                isOpen={isNodeUpdate}
                handleClose={handleCloseUpdateAllNode}
                title={"Node Update all Confirmation"}
                btnLabel="continue with update all"
                content={`The software updates for “Tryphena’s Node,” and “Tryphena’s Node 2” will disrupt your network, and will take approximately [insert time here]. Continue updating all?`}
            />

            <NodeSoftwareInfosDialog
                closeBtnLabel="close"
                isOpen={showNodeSoftwareUpdatInfos}
                handleClose={handleCloseNodeInfos}
            />
            <NodeDialog
                isOpen={isAddNode}
                nodeData={NODE_INIT}
                dialogTitle={"Register Node"}
                handleClose={handleAddNodeClose}
                handleNodeSubmitAction={handleActivationSubmit}
                subTitle={
                    "Ensure node is properly set up in desired location before completing this step. Enter serial number found in your confirmation email, or on the back of your node, and we’ll take care of the rest for you."
                }
            />
            {backToPreviousNode && (
                <Fab
                    variant="extended"
                    color="primary"
                    aria-label="back"
                    className={classes.backToNodeGroupButtonStyle}
                    onClick={() => handleBackToSingleTowerNode()}
                >
                    <ArrowBackIcon sx={{ mr: 1 }} />
                    RETURN TO ORIGINAL NODE
                </Fab>
            )}
        </Box>
    );
};

export default Nodes;
