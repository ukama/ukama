import {
    StatsCard,
    NodeDialog,
    StatusCard,
    BasicDialog,
    NodeContainer,
    NetworkStatus,
    LoadingWrapper,
    DeactivateUser,
    ContainerHeader,
    DataTableWithOptions,
    UserActivationDialog,
    SoftwareUpdateModal,
} from "../../components";
import {
    TIME_FILTER,
    MONTH_FILTER,
    UserActivation,
    DataTableWithOptionColumns,
    DEACTIVATE_EDIT_ACTION_MENU,
} from "../../constants";
import "../../i18n/i18n";
import {
    Node_Type,
    Time_Filter,
    Network_Type,
    Data_Bill_Filter,
    useGetNetworkQuery,
    useGetDataBillQuery,
    useGetDataUsageQuery,
    useGetUsersByOrgQuery,
    useGetNodesByOrgQuery,
    GetLatestNetworkDocument,
    useDeactivateUserMutation,
    useAddNodeMutation,
    useGetConnectedUsersQuery,
    GetLatestDataBillDocument,
    GetLatestDataUsageDocument,
    useGetMetricsByTabLazyQuery,
    GetLatestNetworkSubscription,
    GetLatestDataBillSubscription,
    GetLatestDataUsageSubscription,
    GetLatestConnectedUsersDocument,
    useGetMetricsByTabSSubscription,
    GetLatestConnectedUsersSubscription,
    useUpdateNodeMutation,
} from "../../generated";
import { TMetric } from "../../types";
import { Box, Grid } from "@mui/material";
import { RoundedCard } from "../../styles";
import { useEffect, useState } from "react";
import { useRecoilState, useRecoilValue, useSetRecoilState } from "recoil";
import {
    user,
    isFirstVisit,
    isSkeltonLoading,
    snackbarMessage,
} from "../../recoil";
import { DataBilling, DataUsage, UsersWithBG } from "../../assets/svg";
import { getMetricPayload, isContainNodeUpdate } from "../../utils";

const Home = () => {
    const isSkeltonLoad = useRecoilValue(isSkeltonLoading);
    const [_isFirstVisit, _setIsFirstVisit] = useRecoilState(isFirstVisit);
    const { id: orgId = "" } = useRecoilValue(user);
    const [isUserActivateOpen, setIsUserActivateOpen] = useState(false);
    const [isWelcomeDialog, setIsWelcomeDialog] = useState(false);
    const [userStatusFilter, setUserStatusFilter] = useState(Time_Filter.Total);
    const [dataStatusFilter, setDataStatusFilter] = useState(Time_Filter.Month);
    const [showNodeDialog, setShowNodeDialog] = useState({
        type: "add",
        isShow: false,
        title: "Register Node",
        subTitle:
            "Ensure node is properly set up in desired location before completing this step. Enter serial number found in your confirmation email, or on the back of your node, and we’ll take care of the rest for you.",
        nodeData: {
            type: "HOME",
            name: "",
            nodeId: "",
            orgId: "",
        },
    });
    const [deactivateUserDialog, setDeactivateUserDialog] = useState({
        isShow: false,
        userId: "",
        userName: "",
    });

    const [isSoftwaUpdate, setIsSoftwaUpdate] = useState<boolean>(false);
    const [isMetricPolling, setIsMetricPolling] = useState<boolean>(false);
    const setRegisterNodeNotification = useSetRecoilState(snackbarMessage);
    const [billingStatusFilter, setBillingStatusFilter] = useState(
        Data_Bill_Filter.July
    );
    const [uptimeMetric, setUptimeMetrics] = useState<TMetric>({
        temperaturetrx: null,
    });

    const {
        data: nodeRes,
        loading: nodeLoading,
        refetch: refetchGetNodesByOrg,
    } = useGetNodesByOrgQuery({ fetchPolicy: "network-only" });

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
        updateNode,
        {
            loading: updateNodeLoading,
            data: updateNodeRes,
            error: updateNodError,
        },
    ] = useUpdateNodeMutation({
        onCompleted: () => {
            setRegisterNodeNotification({
                id: "UpdateNodeNotification",
                message: `${updateNodeRes?.updateNode?.nodeId} has been updated successfully!`,
                type: "success",
                show: true,
            });
            refetchGetNodesByOrg();
        },
        onError: () =>
            setRegisterNodeNotification({
                id: "UpdateNodeErrorNotification",
                message: `${updateNodError?.message}`,
                type: "error",
                show: true,
            }),
    });

    const {
        data: connectedUserRes,
        loading: connectedUserloading,
        subscribeToMore: subscribeToLatestConnectedUsers,
    } = useGetConnectedUsersQuery({
        variables: {
            filter: userStatusFilter,
        },
    });
    const {
        data: dataBillingRes,
        loading: dataBillingloading,
        subscribeToMore: subscribeToLatestDataBill,
    } = useGetDataBillQuery({
        skip: true,
        variables: {
            filter: billingStatusFilter,
        },
    });
    const {
        data: dataUsageRes,
        loading: dataUsageloading,
        subscribeToMore: subscribeToLatestDataUsage,
    } = useGetDataUsageQuery({
        variables: {
            filter: dataStatusFilter,
        },
    });

    const {
        data: residentsRes,
        loading: residentsloading,
        refetch: refetchResidents,
    } = useGetUsersByOrgQuery();

    const [deactivateUser, { loading: deactivateUserLoading }] =
        useDeactivateUserMutation({
            onCompleted: res => {
                setRegisterNodeNotification({
                    id: "userDeactivated",
                    message: `${res.deactivateUser.name} has been deactivated successfully!`,
                    type: "success",
                    show: true,
                });
                refetchResidents();
            },
            onError: err =>
                setRegisterNodeNotification({
                    id: "userDeactivated",
                    message: `${err?.message}`,
                    type: "error",
                    show: true,
                }),
        });

    const {
        data: networkStatusRes,
        loading: networkStatusLoading,
        subscribeToMore: subscribeToLatestNetworkStatus,
    } = useGetNetworkQuery({
        variables: {
            filter: Network_Type.Public,
        },
    });

    const [
        getMetrics,
        {
            data: getMetricsRes,
            refetch: getMetricsRefetch,
            loading: getMetricLoading,
        },
    ] = useGetMetricsByTabLazyQuery({
        onCompleted: res => {
            if (res?.getMetricsByTab?.metrics.length > 0 && !isMetricPolling) {
                const _m: TMetric = {
                    temperaturetrx: null,
                };
                setIsMetricPolling(true);
                for (const element of res.getMetricsByTab.metrics) {
                    if (!uptimeMetric[element.type]) {
                        _m[element.type] = {
                            name: element.name,
                            data: element.data,
                        };
                    }
                }
                setUptimeMetrics({ ..._m });
            }
        },
        onError: () => {
            setUptimeMetrics(() => ({
                temperaturetrx: null,
            }));
        },
        fetchPolicy: "network-only",
    });

    useGetMetricsByTabSSubscription({
        onSubscriptionData: res => {
            if (
                isMetricPolling &&
                res?.subscriptionData?.data?.getMetricsByTab &&
                res?.subscriptionData?.data?.getMetricsByTab.length > 0
            ) {
                const _m: TMetric = {
                    temperaturetrx: null,
                };
                for (const element of res.subscriptionData.data
                    .getMetricsByTab) {
                    const metric = uptimeMetric[element.type];
                    if (
                        metric &&
                        metric.data &&
                        metric.data.length > 0 &&
                        element.data[element.data.length - 1].x >
                            metric.data[metric.data.length - 1].x
                    ) {
                        _m[element.type] = {
                            name: element.name,
                            data: [...(metric.data || []), ...element.data],
                        };
                    }
                }
                const filter = Object.fromEntries(
                    Object.entries(_m).filter(([_, v]) => v !== null)
                );

                setUptimeMetrics((_prev: TMetric) => ({
                    ..._prev,
                    ...filter,
                }));
            }
        },
    });

    useEffect(() => {
        if (_isFirstVisit && orgId) {
            setIsWelcomeDialog(true);
        }
    }, [_isFirstVisit, orgId]);

    const getFirstMetricCallPayload = () =>
        getMetricPayload({
            tab: 4,
            regPolling: false,
            nodeType: Node_Type.Home,
            nodeId: "uk-sa2209-comv1-a1-ee58",
            to: Math.floor(Date.now() / 1000) - 15,
            from: Math.floor(Date.now() / 1000) - 180,
        });

    const getMetricPollingCallPayload = (from: number) =>
        getMetricPayload({
            tab: 4,
            from: from,
            regPolling: true,
            nodeType: Node_Type.Home,
            nodeId: "uk-sa2209-comv1-a1-ee58",
        });

    useEffect(() => {
        if (nodeRes && nodeRes.getNodesByOrg.nodes.length > 0) {
            getMetrics({
                variables: {
                    ...getFirstMetricCallPayload(),
                },
            });
        }
    }, [nodeRes]);

    useEffect(() => {
        if (
            getMetricsRes &&
            getMetricsRes.getMetricsByTab.next &&
            getMetricsRes?.getMetricsByTab.metrics.length > 0
        ) {
            getMetricsRefetch({
                ...getMetricPollingCallPayload(
                    getMetricsRes?.getMetricsByTab.to
                ),
            });
        }
    }, [getMetricsRes]);

    const handleAddNodeClose = () => {
        setShowNodeDialog(prev => ({
            ...prev,
            isShow: false,
        }));
    };

    const subToConnectedUser = () =>
        subscribeToLatestConnectedUsers<GetLatestConnectedUsersSubscription>({
            document: GetLatestConnectedUsersDocument,
            updateQuery: (prev, { subscriptionData }) => {
                let data = { ...prev };
                const latestConnectedUser =
                    subscriptionData.data.getConnectedUsers;
                if (latestConnectedUser.__typename === "ConnectedUserDto")
                    data.getConnectedUsers = latestConnectedUser;
                return data;
            },
        });

    const subToDataUsage = () =>
        subscribeToLatestDataUsage<GetLatestDataUsageSubscription>({
            document: GetLatestDataUsageDocument,
            updateQuery: (prev, { subscriptionData }) => {
                let data = { ...prev };
                const latestDataUsage = subscriptionData.data.getDataUsage;
                if (latestDataUsage.__typename === "DataUsageDto")
                    data.getDataUsage = latestDataUsage;
                return data;
            },
        });

    const subToDataBill = () =>
        subscribeToLatestDataBill<GetLatestDataBillSubscription>({
            document: GetLatestDataBillDocument,
            updateQuery: (prev, { subscriptionData }) => {
                let data = { ...prev };
                const latestDataBill = subscriptionData.data.getDataBill;
                if (latestDataBill.__typename === "DataBillDto")
                    data.getDataBill = latestDataBill;
                return data;
            },
        });

    useEffect(() => {
        let unsub = subToConnectedUser();
        return () => {
            unsub && unsub();
        };
    }, [connectedUserRes]);

    useEffect(() => {
        let unsub = subToDataUsage();
        return () => {
            unsub && unsub();
        };
    }, [dataUsageRes]);

    useEffect(() => {
        let unsub = subToDataBill();
        return () => {
            unsub && unsub();
        };
    }, [dataBillingRes]);

    useEffect(() => {
        if (networkStatusRes) {
            subscribeToLatestNetworkStatus<GetLatestNetworkSubscription>({
                document: GetLatestNetworkDocument,
                updateQuery: (prev, { subscriptionData }) => {
                    let data = { ...prev };
                    const latestNewtworkStatus =
                        subscriptionData.data.getNetwork;
                    if (latestNewtworkStatus.__typename === "NetworkDto")
                        data.getNetwork = latestNewtworkStatus;
                    return data;
                },
            });
        }
    }, [networkStatusRes]);

    const handleSatusChange = (key: string, value: string) => {
        switch (key) {
            case "statusUser":
                return setUserStatusFilter(value as Time_Filter);
            case "statusUsage":
                return setDataStatusFilter(value as Time_Filter);
            case "statusBill":
                return setBillingStatusFilter(value as Data_Bill_Filter);
        }
    };
    const handleCloseSoftwareUpdate = () => {
        setIsSoftwaUpdate(false);
    };

    const getStatus = (key: string) => {
        switch (key) {
            case "statusUser":
                return userStatusFilter;
            case "statusUsage":
                return dataStatusFilter;
            case "statusBill":
                return billingStatusFilter;
            default:
                return "";
        }
    };

    const handleUserActivateClose = () => setIsUserActivateOpen(() => false);
    const handleCloseDeactivateUser = () =>
        setDeactivateUserDialog({ ...deactivateUserDialog, isShow: false });

    const handleDeactivateUser = () => {
        handleCloseDeactivateUser();
        deactivateUser({
            variables: {
                id: deactivateUserDialog.userId,
            },
        });
    };

    const onResidentsTableMenuItem = (id: string, type: string) => {
        if (type === "deactivate") {
            setDeactivateUserDialog({
                isShow: true,
                userId: id,
                userName:
                    residentsRes?.getUsersByOrg.find(item => item.id === id)
                        ?.name || "",
            });
        }
    };

    const handleNodeActions = (id: string, type: string) => {
        const node = nodeRes?.getNodesByOrg.nodes.filter(
            item => item.id === id
        );
        if (type == "edit" && node && node.length > 0) {
            setShowNodeDialog({
                ...showNodeDialog,
                type: "editNode",
                isShow: true,
                title: "Edit Node",
                nodeData: {
                    type: node[0].type,
                    name: node[0].name,
                    nodeId: node[0].id,
                    orgId: orgId,
                },
            });
        }
    };

    const handleAddNode = () => {
        setShowNodeDialog(prev => ({
            ...prev,
            type: "add",
            isShow: true,
            title: "Register Node",
        }));
    };

    const onUpdateAllNodes = () => {
        /* TODO: Handle Node Updates */
    };

    const handleNodeSubmitAction = (data: any) => {
        setShowNodeDialog(prev => ({
            ...prev,
            type: "add",
            isShow: false,
            title: "Register Node",
        }));
        if (showNodeDialog.type === "add") {
            registerNode({
                variables: {
                    data: {
                        name: data.name,
                        nodeId: data.nodeId,
                    },
                },
            });
        } else if (showNodeDialog.type === "editNode") {
            updateNode({
                variables: {
                    data: {
                        name: data.name,
                        nodeId: data.nodeId,
                    },
                },
            });
        }
    };

    const onActivateUser = () => setIsUserActivateOpen(() => true);

    // eslint-disable-next-line no-unused-vars
    const handleNodeUpdateActin = (id: string) => {
        setIsSoftwaUpdate(true);
        /* Handle node update  action */
    };

    const handleCloseWelcome = () => {
        if (_isFirstVisit) {
            _setIsFirstVisit(false);
            setIsWelcomeDialog(false);
        }
    };
    return (
        <Box component="div" sx={{ flexGrow: 1, pb: "18px" }}>
            <Grid container spacing={3}>
                <Grid xs={12} item>
                    <NetworkStatus
                        handleAddNode={handleAddNode}
                        handleActivateUser={onActivateUser}
                        loading={networkStatusLoading || isSkeltonLoad}
                        regLoading={registerNodeLoading || updateNodeLoading}
                        statusType={networkStatusRes?.getNetwork?.status || ""}
                        duration={
                            networkStatusRes?.getNetwork?.description || ""
                        }
                    />
                </Grid>
                <Grid item container xs={12} spacing={{ xs: 1.5, md: 3 }}>
                    <Grid item xs={4} md={6} lg={4}>
                        <StatusCard
                            Icon={UsersWithBG}
                            title={"Connected Users"}
                            options={TIME_FILTER}
                            subtitle1={`${
                                connectedUserRes?.getConnectedUsers
                                    ?.totalUser || 0
                            }`}
                            subtitle2={""}
                            option={getStatus("statusUser")}
                            loading={connectedUserloading || isSkeltonLoad}
                            handleSelect={(value: string) =>
                                handleSatusChange("statusUser", value)
                            }
                        />
                    </Grid>
                    <Grid item xs={4} md={6} lg={4}>
                        <StatusCard
                            title={"Data Usage"}
                            subtitle1={`${
                                dataUsageRes?.getDataUsage?.dataConsumed || 0
                            }`}
                            subtitle2={`/ ${
                                dataUsageRes?.getDataUsage?.dataPackage || "-"
                            }`}
                            Icon={DataUsage}
                            options={TIME_FILTER}
                            option={getStatus("statusUsage")}
                            loading={dataUsageloading || isSkeltonLoad}
                            handleSelect={(value: string) =>
                                handleSatusChange("statusUsage", value)
                            }
                        />
                    </Grid>
                    <Grid item xs={4} md={6} lg={4}>
                        <StatusCard
                            title={"Data Bill"}
                            subtitle1={`$ ${
                                dataBillingRes?.getDataBill?.dataBill || 0
                            }`}
                            subtitle2={
                                dataBillingRes?.getDataBill?.dataBill
                                    ? ` / due in ${dataBillingRes?.getDataBill?.billDue} days`
                                    : " due"
                            }
                            Icon={DataBilling}
                            options={MONTH_FILTER}
                            loading={dataBillingloading || isSkeltonLoad}
                            option={getStatus("statusBill")}
                            handleSelect={(value: string) =>
                                handleSatusChange("statusBill", value)
                            }
                        />
                    </Grid>
                </Grid>
                <Grid xs={12} item>
                    <StatsCard
                        metricData={uptimeMetric}
                        loading={
                            nodeLoading || isSkeltonLoad || getMetricLoading
                        }
                    />
                </Grid>
                <Grid xs={12} lg={8} item>
                    <LoadingWrapper
                        height={318}
                        isLoading={nodeLoading || isSkeltonLoad}
                    >
                        <RoundedCard>
                            <ContainerHeader
                                stats={``}
                                title="My Nodes"
                                showButton={isContainNodeUpdate(
                                    nodeRes?.getNodesByOrg.nodes
                                )}
                                buttonSize={"small"}
                                buttonTitle={"Update All"}
                            />
                            <NodeContainer
                                handleNodeUpdate={handleNodeUpdateActin}
                                items={nodeRes?.getNodesByOrg.nodes || []}
                                handleItemAction={handleNodeActions}
                            />
                        </RoundedCard>
                    </LoadingWrapper>
                </Grid>
                <Grid xs={12} lg={4} item>
                    <LoadingWrapper
                        height={318}
                        isLoading={
                            residentsloading ||
                            deactivateUserLoading ||
                            isSkeltonLoad
                        }
                    >
                        <RoundedCard sx={{ height: "100%" }}>
                            <ContainerHeader
                                title="Residents"
                                showButton={false}
                            />
                            <DataTableWithOptions
                                columns={DataTableWithOptionColumns}
                                dataset={residentsRes?.getUsersByOrg}
                                menuOptions={DEACTIVATE_EDIT_ACTION_MENU}
                                onMenuItemClick={onResidentsTableMenuItem}
                            />
                        </RoundedCard>
                    </LoadingWrapper>
                </Grid>
            </Grid>
            <SoftwareUpdateModal
                submit={onUpdateAllNodes}
                isOpen={isSoftwaUpdate}
                handleClose={handleCloseSoftwareUpdate}
                title={" Node Update all Confirmation"}
                content={` The software updates for “Tryphena’s Node,” and
                “Tryphena’s Node 2” will disrupt your network, and will
                take approximately [insert time here]. Continue updating
                all?`}
            />
            {isWelcomeDialog && (
                <BasicDialog
                    isClosable={false}
                    isOpen={isWelcomeDialog}
                    title={"Welcome to Ukama Console!"}
                    description={
                        "This is where you can manage your network, and troubleshoot things, if necessary. For now, while your nodes have not shipped, you can monitor your users’ data usage, and [insert other main use]. "
                    }
                    labelSuccessBtn={"continue to console"}
                    handleCloseAction={handleCloseWelcome}
                />
            )}
            {isUserActivateOpen && (
                <UserActivationDialog
                    isOpen={isUserActivateOpen}
                    dialogTitle={UserActivation.title}
                    subTitle={UserActivation.subTitle}
                    handleClose={handleUserActivateClose}
                />
            )}

            {showNodeDialog.isShow && (
                <NodeDialog
                    action={showNodeDialog.type}
                    isOpen={showNodeDialog.isShow}
                    handleClose={handleAddNodeClose}
                    nodeData={showNodeDialog.nodeData}
                    dialogTitle={showNodeDialog.title}
                    subTitle={showNodeDialog.subTitle}
                    handleNodeSubmitAction={handleNodeSubmitAction}
                />
            )}

            {deactivateUserDialog.isShow && (
                <DeactivateUser
                    isClosable={true}
                    isOpen={deactivateUserDialog.isShow}
                    title={"Deactivate User Confirmation"}
                    description={`${deactivateUserDialog.userName} will be deactivated permanently. Other copy depends on surrounding policy.`}
                    labelSuccessBtn={"DEACTIVATE USER"}
                    labelNegativeBtn={"cancel"}
                    handleCloseAction={handleCloseDeactivateUser}
                    handleSuccessAction={handleDeactivateUser}
                />
            )}
        </Box>
    );
};
export default Home;
