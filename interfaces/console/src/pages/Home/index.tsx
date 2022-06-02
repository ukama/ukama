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
    SoftwareUpdateModal,
    UserDetailsDialog,
    AddUser,
} from "../../components";
import {
    TIME_FILTER,
    MONTH_FILTER,
    DataTableWithOptionColumns,
    DEACTIVATE_EDIT_ACTION_MENU,
} from "../../constants";
import "../../i18n/i18n";
import {
    Node_Type,
    Time_Filter,
    GetUsersDto,
    UserInputDto,
    Data_Bill_Filter,
    useGetUserLazyQuery,
    useAddNodeMutation,
    useUpdateUserMutation,
    useAddUserMutation,
    useGetDataBillQuery,
    useGetDataUsageQuery,
    useUpdateNodeMutation,
    useGetUsersByOrgQuery,
    GetUserDto,
    useGetNodesByOrgQuery,
    useDeleteNodeMutation,
    useGetEsimQrLazyQuery,
    useDeactivateUserMutation,
    useGetConnectedUsersQuery,
    GetLatestDataBillDocument,
    GetLatestDataUsageDocument,
    useGetMetricsByTabLazyQuery,
    GetLatestDataBillSubscription,
    useGetUsersDataUsageLazyQuery,
    GetLatestDataUsageSubscription,
    useUpdateUserStatusMutation,
    GetLatestConnectedUsersDocument,
    useGetMetricsByTabSSubscription,
    useGetUsersDataUsageSSubscription,
    GetLatestConnectedUsersSubscription,
    useGetNetworkStatusQuery,
    GetNetworkStatusSDocument,
    GetNetworkStatusSSubscription,
} from "../../generated";
import {
    user,
    isFirstVisit,
    isSkeltonLoading,
    snackbarMessage,
} from "../../recoil";
import { Box, Grid } from "@mui/material";
import { RoundedCard } from "../../styles";
import { useEffect, useState } from "react";
import { TMetric } from "../../types";
import { useRecoilState, useRecoilValue, useSetRecoilState } from "recoil";
import {
    getMetricPayload,
    getTowerNodeFromNodes,
    isContainNodeUpdate,
} from "../../utils";
import { DataBilling, DataUsage, UsersWithBG } from "../../assets/svg";
const userInit = {
    id: "",
    name: "",
    iccid: "",
    email: "",
    phone: "",
    dataPlan: "0",
    dataUsage: "0",
    roaming: false,
    eSimNumber: "",
    status: false,
};
const Home = () => {
    const isSkeltonLoad = useRecoilValue(isSkeltonLoading);
    const [_isFirstVisit, _setIsFirstVisit] = useRecoilState(isFirstVisit);
    const { id: orgId = "" } = useRecoilValue(user);
    const [users, setUsers] = useState<GetUsersDto[]>([]);
    const [isWelcomeDialog, setIsWelcomeDialog] = useState(false);
    const [userStatusFilter, setUserStatusFilter] = useState(Time_Filter.Total);
    const [dataStatusFilter, setDataStatusFilter] = useState(Time_Filter.Month);
    const [newAddedUserName, setNewAddedUserName] = useState<any>();
    const [isEsimAdded, setIsEsimAdded] = useState<boolean>(false);
    const [isPsimAdded, setIsPsimAdded] = useState<boolean>(false);
    const [simFlow, setSimFlow] = useState<number>(1);
    const [userStatus, setUserStatus] = useState<boolean>(true);
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
    const [deleteNodeDialog, setDeleteNodeDialog] = useState({
        isShow: false,
        nodeId: "",
    });
    const [simDialog, setSimDialog] = useState({
        isShow: false,
        type: "add",
    });
    const [selectedUser, setSelectedUser] = useState<GetUserDto>(userInit);
    const [qrCodeId, setqrCodeId] = useState<any>();
    const [isSoftwaUpdate, setIsSoftwaUpdate] = useState<boolean>(false);
    const [showInstallSim, setShowInstallSim] = useState(false);
    const [isMetricPolling, setIsMetricPolling] = useState<boolean>(false);
    const setNodeToastNotification = useSetRecoilState(snackbarMessage);
    const [billingStatusFilter, setBillingStatusFilter] = useState(
        Data_Bill_Filter.July
    );
    const [uptimeMetric, setUptimeMetrics] = useState<TMetric>({
        memorytrxused: null,
    });

    const {
        data: nodeRes,
        loading: nodeLoading,
        refetch: refetchGetNodesByOrg,
    } = useGetNodesByOrgQuery({ fetchPolicy: "network-only" });

    const [deleteNode, { loading: deleteNodeLoading }] = useDeleteNodeMutation({
        onCompleted: res => {
            setNodeToastNotification({
                id: "delete-node-success",
                message: `${res?.deleteNode?.nodeId} has been deleted successfully!`,
                type: "success",
                show: true,
            });
            refetchGetNodesByOrg();
        },
        onError: err => {
            setNodeToastNotification({
                id: "delete-node-success",
                message: `${err?.message}`,
                type: "error",
                show: true,
            });
        },
    });

    const [addUser, { loading: addUserLoading }] = useAddUserMutation({
        onCompleted: res => {
            if (res?.addUser) {
                setIsEsimAdded(true);
                setNewAddedUserName(res?.addUser?.name);
                handleGetSimQrCode(res?.addUser.id, res?.addUser?.iccid || "");
                refetchResidents();

                handleUpdateUserStatus(
                    res.addUser.id,
                    res.addUser.iccid || "",
                    userStatus
                );
            }
        },
        onError: err => {
            if (err?.message) {
                setNodeToastNotification({
                    id: "error-add-user-success",
                    message: `${err?.message}`,
                    type: "error",
                    show: true,
                });
            }
        },
    });

    const [getEsimQrdcodeId, { data: getEsimQrCodeRes }] =
        useGetEsimQrLazyQuery();

    const handleGetSimQrCode = async (userId: string, simId: string) => {
        await getEsimQrdcodeId({
            variables: {
                data: {
                    userId: userId,
                    simId: simId,
                },
            },
        });
    };

    useEffect(() => {
        setqrCodeId(getEsimQrCodeRes?.getEsimQR?.qrCode);
    }, [getEsimQrCodeRes]);
    const [registerNode, { loading: registerNodeLoading }] = useAddNodeMutation(
        {
            onCompleted: res => {
                if (res?.addNode) {
                    setNodeToastNotification({
                        id: "addNodeSuccess",
                        message: `${res?.addNode?.name} has been registered successfully!`,
                        type: "success",
                        show: true,
                    });
                    refetchGetNodesByOrg();
                }
            },

            onError: err => {
                if (err?.message) {
                    setNodeToastNotification({
                        id: "ErrorAddingNode",
                        message: `${err?.message}`,
                        type: "error",
                        show: true,
                    });
                }
            },
        }
    );

    const [updateNode, { loading: updateNodeLoading }] = useUpdateNodeMutation({
        onCompleted: res => {
            if (res?.updateNode) {
                setNodeToastNotification({
                    id: "UpdateNodeNotification",
                    message: `${res?.updateNode?.nodeId} has been updated successfully!`,
                    type: "success",
                    show: true,
                });
                refetchGetNodesByOrg();
            }
        },
        onError: err => {
            if (err?.message) {
                setNodeToastNotification({
                    id: "UpdateNodeErrorNotification",
                    message: `${err?.message}`,
                    type: "error",
                    show: true,
                });
            }
        },
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

    const [getUsersDataUsage] = useGetUsersDataUsageLazyQuery();

    useGetUsersDataUsageSSubscription({
        fetchPolicy: "network-only",
        onSubscriptionData: res => {
            if (res.subscriptionData.data?.getUsersDataUsage?.id) {
                const userRes = res.subscriptionData.data?.getUsersDataUsage;
                const index = users.findIndex(item => item.id === userRes.id);
                setUsers([
                    ...users.slice(0, index),
                    {
                        id: userRes.id,
                        name: userRes.name,
                        email: userRes.email,
                        phone: userRes.phone,
                        dataPlan: userRes.dataPlan,
                        dataUsage: userRes.dataUsage,
                    },
                    ...users.slice(index + 1),
                ]);
            }
        },
    });

    const { loading: residentsloading, refetch: refetchResidents } =
        useGetUsersByOrgQuery({
            nextFetchPolicy: "network-only",
            onCompleted: res => {
                setUsers([...res.getUsersByOrg].reverse());
                getUsersDataUsage({
                    variables: {
                        data: { ids: res.getUsersByOrg.map(u => u.id) },
                    },
                });
            },
        });

    const [deactivateUser, { loading: deactivateUserLoading }] =
        useDeactivateUserMutation({
            onCompleted: res => {
                setNodeToastNotification({
                    id: "userDeactivated",
                    message: `${res.deactivateUser.name} has been deactivated successfully!`,
                    type: "success",
                    show: true,
                });
                refetchResidents();
            },
            onError: err =>
                setNodeToastNotification({
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
    } = useGetNetworkStatusQuery();

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
                    memorytrxused: null,
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
                memorytrxused: null,
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
                    memorytrxused: null,
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
    const handleSimInstallationClose = () => {
        setSimFlow(1);
        setIsEsimAdded(false);
        setShowInstallSim(false);
    };

    const getFirstMetricCallPayload = () =>
        getMetricPayload({
            tab: 4,
            regPolling: false,
            nodeType: Node_Type.Home,
            nodeId: getTowerNodeFromNodes(nodeRes?.getNodesByOrg.nodes || []),
            to: Math.floor(Date.now() / 1000) - 10,
            from: Math.floor(Date.now() / 1000) - 180,
        });

    const getMetricPollingCallPayload = (from: number) =>
        getMetricPayload({
            tab: 4,
            from: from,
            regPolling: true,
            nodeType: Node_Type.Home,
            nodeId: getTowerNodeFromNodes(nodeRes?.getNodesByOrg.nodes || []),
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
            subscribeToLatestNetworkStatus<GetNetworkStatusSSubscription>({
                document: GetNetworkStatusSDocument,
                updateQuery: (prev, { subscriptionData }) => {
                    let data = { ...prev };
                    const latestNewtworkStatus =
                        subscriptionData.data.getNetworkStatus;
                    if (latestNewtworkStatus.__typename === "NetworkDto")
                        data.getNetworkStatus = latestNewtworkStatus;
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
    const [getUser, { loading: getUserLoading }] = useGetUserLazyQuery({
        onCompleted: res => {
            if (res.getUser) {
                setSelectedUser(res.getUser);
            }
        },
    });

    const handleDeleteNode = () => {
        deleteNode({
            variables: {
                id: deleteNodeDialog.nodeId,
            },
        });
        handleCloseDeleteNode();
    };

    const onResidentsTableMenuItem = (id: string, type: string) => {
        if (type === "deactivate") {
            setDeactivateUserDialog({
                isShow: true,
                userId: id,
                userName: users?.find(item => item.id === id)?.name || "",
            });
        } else if (type === "edit") {
            getUser({
                variables: {
                    userId: id,
                },
            });
            setSimDialog({ isShow: true, type: "edit" });
        }
    };
    const handleUserSubmitAction = () => {
        handleSimDialogClose();
        if (simDialog.type === "edit" && selectedUser.id) {
            updateUser({
                variables: {
                    userId: selectedUser.id,
                    data: {
                        email: selectedUser.email,
                        name: selectedUser.name,
                        phone: selectedUser.phone,
                        status: selectedUser.status,
                    },
                },
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
        } else {
            setDeleteNodeDialog({
                isShow: true,
                nodeId: id || "",
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
    const [updateUser, { loading: updateUserLoading }] = useUpdateUserMutation({
        onCompleted: res => {
            if (res?.updateUser) {
                setNodeToastNotification({
                    id: "updateUserNotification",
                    message: `The ${res?.updateUser?.name} has been updated successfully!`,
                    type: "success",
                    show: true,
                });
                refetchResidents();
            }
        },
        onError: err => {
            if (err?.message) {
                setNodeToastNotification({
                    id: "updateUserNotification",
                    message: `${err?.message}`,
                    type: "error",
                    show: true,
                });
            }
        },
    });
    const [updateUserStatus, { loading: updateUserStatusLoading }] =
        useUpdateUserStatusMutation({
            onCompleted: res => {
                if (res) {
                    setSelectedUser({
                        ...selectedUser,
                        status: res.updateUserStatus.carrier.services.data,
                        roaming: res.updateUserStatus.ukama.services.data,
                    });
                }
            },
        });
    const handleUpdateUserStatus = (
        id: string,
        iccid: string,
        status: boolean
    ) => {
        updateUserStatus({
            variables: {
                data: {
                    userId: id,
                    simId: iccid,
                    status: status,
                },
            },
        });
    };
    const onUpdateAllNodes = () => {
        /* TODO: Handle Node Updates */
    };
    const handleSimDialogClose = () =>
        setSimDialog({ ...simDialog, isShow: false });

    const handleCloseDeleteNode = () =>
        setDeleteNodeDialog({ ...deleteNodeDialog, isShow: false });
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

    const onActivateUser = () => setShowInstallSim(() => true);

    const handleNodeUpdateActin = () => {
        setIsSoftwaUpdate(true);
        /* Handle node update  action */
    };

    const handleCloseWelcome = () => {
        if (_isFirstVisit) {
            _setIsFirstVisit(false);
            setIsWelcomeDialog(false);
        }
    };
    const handleEsimInstallation = (eSimData: UserInputDto) => {
        setUserStatus(eSimData.status);
        if (eSimData) {
            addUser({
                variables: {
                    data: {
                        email: eSimData.email,
                        name: eSimData.name,
                        status: eSimData.status,
                        phone: "",
                    },
                },
            });
        }
    };
    const handlePhysicalSimEmailFlow = () => {
        setSimFlow(simFlow + 1);
    };

    const handlePhysicalSimSecurityFlow = () => {
        setSimFlow(simFlow + 1);
        setIsPsimAdded(true);
    };
    const handleDeactivateAction = (userId: any) => {
        setSimDialog({ ...simDialog, isShow: false });
        setDeactivateUserDialog({
            isShow: true,
            userId: userId,
            userName: users?.find(item => item.id === userId)?.name || "",
        });
    };
    return (
        <Box component="div" sx={{ flexGrow: 1, pb: "18px" }}>
            <Grid container rowSpacing={3} columnSpacing={3}>
                <Grid xs={12} item>
                    <NetworkStatus
                        handleAddNode={handleAddNode}
                        handleActivateUser={onActivateUser}
                        loading={networkStatusLoading || isSkeltonLoad}
                        regLoading={registerNodeLoading || updateNodeLoading}
                        statusType={
                            networkStatusRes?.getNetworkStatus?.status ||
                            undefined
                        }
                        duration={
                            networkStatusRes?.getNetworkStatus?.uptime ||
                            undefined
                        }
                    />
                </Grid>
                <Grid item container columnSpacing={{ xs: 1.5, md: 3 }}>
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
                            subtitle2={`${
                                dataUsageRes?.getDataUsage?.dataPackage || ""
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
                        isLoading={
                            nodeLoading || isSkeltonLoad || deleteNodeLoading
                        }
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
                            isSkeltonLoad ||
                            addUserLoading ||
                            updateUserLoading
                        }
                    >
                        <RoundedCard sx={{ height: "100%" }}>
                            <ContainerHeader
                                title="Residents"
                                showButton={false}
                            />
                            <DataTableWithOptions
                                dataset={users}
                                columns={DataTableWithOptionColumns}
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
            {deleteNodeDialog.isShow && (
                <DeactivateUser
                    isClosable={true}
                    isOpen={deleteNodeDialog.isShow}
                    title={"Delete Node Confirmation"}
                    description={`${deleteNodeDialog?.nodeId} will be deleted permanently.`}
                    labelSuccessBtn={"DELETE NODE"}
                    labelNegativeBtn={"cancel"}
                    handleCloseAction={handleCloseDeleteNode}
                    handleSuccessAction={handleDeleteNode}
                />
            )}
            {simDialog.isShow && (
                <UserDetailsDialog
                    user={selectedUser}
                    type={simDialog.type}
                    saveBtnLabel={"Save"}
                    closeBtnLabel="close"
                    loading={getUserLoading}
                    isOpen={simDialog.isShow}
                    setUserForm={setSelectedUser}
                    simDetailsTitle="SIM Details"
                    userDetailsTitle="User Details"
                    handleClose={handleSimDialogClose}
                    userStatusLoading={updateUserStatusLoading}
                    handleServiceAction={handleUpdateUserStatus}
                    handleSubmitAction={handleUserSubmitAction}
                    handleDeactivateAction={handleDeactivateAction}
                />
            )}

            {showInstallSim && (
                <AddUser
                    loading={addUserLoading}
                    handleEsimInstallation={handleEsimInstallation}
                    addedUserName={newAddedUserName}
                    qrCodeId={qrCodeId}
                    isPsimAdded={isPsimAdded}
                    iSeSimAdded={isEsimAdded}
                    handlePhysicalSimInstallationFlow1={
                        handlePhysicalSimEmailFlow
                    }
                    handlePhysicalSimInstallationFlow2={
                        handlePhysicalSimSecurityFlow
                    }
                    step={simFlow}
                    isOpen={showInstallSim}
                    handleClose={handleSimInstallationClose}
                />
            )}
        </Box>
    );
};
export default Home;
