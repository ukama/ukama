import {
    StatsCard,
    StatusCard,
    NodeContainer,
    NetworkStatus,
    LoadingWrapper,
    ContainerHeader,
    ActivationDialog,
    DataTableWithOptions,
    UserActivationDialog,
    SoftwareUpdateModal,
    BasicDialog,
} from "../../components";
import { useTranslation } from "react-i18next";
import {
    TIME_FILTER,
    MONTH_FILTER,
    UserActivation,
    DataTableWithOptionColumns,
    DEACTIVATE_EDIT_ACTION_MENU,
} from "../../constants";
import "../../i18n/i18n";
import {
    MetricDto,
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
    useGetConnectedUsersQuery,
    GetLatestDataBillDocument,
    GetLatestDataUsageDocument,
    GetLatestNetworkSubscription,
    useGetMetricsUptimeLazyQuery,
    GetLatestDataBillSubscription,
    GetLatestDataUsageSubscription,
    GetLatestConnectedUsersDocument,
    GetLatestConnectedUsersSubscription,
    useGetMetricsUptimeSSubscription,
} from "../../generated";
import { useRecoilValue } from "recoil";
import { RoundedCard } from "../../styles";
import { useEffect, useState } from "react";
import { isSkeltonLoading, user } from "../../recoil";
import { Box, Grid } from "@mui/material";
import { DataBilling, DataUsage, UsersWithBG } from "../../assets/svg";
import {
    getDefaultMetricList,
    getMetricPayload,
    isContainNodeUpdate,
} from "../../utils";

const Home = () => {
    const { t } = useTranslation();
    const isSkeltonLoad = useRecoilValue(isSkeltonLoading);
    const { id: orgId = "" } = useRecoilValue(user);
    const [showSimActivationDialog, setShowSimActivationDialog] =
        useState<boolean>(false);
    const [isUserActivateOpen, setIsUserActivateOpen] = useState(false);
    const [isWelcomeDialog, setIsWelcomeDialog] = useState(false);
    const [userStatusFilter, setUserStatusFilter] = useState(Time_Filter.Total);
    const [dataStatusFilter, setDataStatusFilter] = useState(Time_Filter.Month);
    const [isAddNode, setIsAddNode] = useState(false);
    const [isSoftwaUpdate, setIsSoftwaUpdate] = useState<boolean>(false);
    const [billingStatusFilter, setBillingStatusFilter] = useState(
        Data_Bill_Filter.July
    );
    const [uptimeMetric, setUptimeMetrics] = useState<{
        name: string;
        data: MetricDto[];
    }>(getDefaultMetricList("UPTIME"));

    const [deactivateUser, { loading: deactivateUserLoading }] =
        useDeactivateUserMutation();
    const handleAddNodeClose = () => setIsAddNode(() => false);
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
        refetch: refetchUser,
    } = useGetUsersByOrgQuery({
        variables: {
            orgId: orgId,
        },
    });

    const { data: nodeRes, loading: nodeLoading } = useGetNodesByOrgQuery({
        variables: {
            orgId: orgId,
        },
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
        getMetricUptime,
        {
            data: metricUptimeRes,
            loading: metricUptimeLoading,
            refetch: metricUptimeRefetch,
        },
    ] = useGetMetricsUptimeLazyQuery({
        fetchPolicy: "network-only",
        onCompleted: res => {
            if (
                uptimeMetric.data.length === 0 &&
                res.getMetricsUptime.length > 0
            ) {
                setUptimeMetrics({
                    name: uptimeMetric.name,
                    data: [
                        ...uptimeMetric.data,
                        ...(res.getMetricsUptime || []),
                    ],
                });
            }
        },
    });

    useEffect(() => {
        if (orgId && !isSkeltonLoad) {
            setShowSimActivationDialog(true);
            setIsWelcomeDialog(true);
        }
    }, [orgId, isSkeltonLoad]);

    useGetMetricsUptimeSSubscription({
        onSubscriptionData: res => {
            if (res.subscriptionData.data && uptimeMetric.data.length > 0) {
                res.subscriptionData?.data?.getMetricsUptime[0].x >
                    uptimeMetric.data[uptimeMetric.data.length - 1].x &&
                    setUptimeMetrics(_prev => ({
                        name: _prev.name,
                        data: [
                            ..._prev.data,
                            ...(res.subscriptionData.data?.getMetricsUptime ||
                                []),
                        ],
                    }));
            }
        },
    });

    const getFirstMetricCallPayload = () =>
        getMetricPayload({
            nodeId: "uk-sa2209-comv1-a1-ee58",
            orgId: orgId,
            regPolling: false,
            to: Math.floor(Date.now() / 1000) - 10,
            from: Math.floor(Date.now() / 1000) - 180,
        });

    const getMetricPollingCallPayload = (from: number) =>
        getMetricPayload({
            nodeId: "uk-sa2209-comv1-a1-ee58",
            orgId: orgId,
            from: from,
        });

    useEffect(() => {
        if (nodeRes && nodeRes.getNodesByOrg.nodes.length > 0) {
            getMetricUptime({
                variables: {
                    ...getFirstMetricCallPayload(),
                },
            });
        }
    }, [nodeRes]);

    useEffect(() => {
        if (
            metricUptimeRes &&
            metricUptimeRes.getMetricsUptime &&
            metricUptimeRes.getMetricsUptime.length > 0
        ) {
            metricUptimeRefetch({
                ...getMetricPollingCallPayload(
                    metricUptimeRes.getMetricsUptime[
                        metricUptimeRes.getMetricsUptime.length - 1
                    ].x + 1
                ),
            });
        }
    }, [metricUptimeRes]);

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
    const handleSimActivateClose = () => {
        setShowSimActivationDialog(false);
    };
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
    const onResidentsTableMenuItem = (id: string, type: string) => {
        if (type === "deactivate") {
            deactivateUser({
                variables: {
                    id,
                },
            });

            refetchUser();
        }
    };

    // eslint-disable-next-line no-unused-vars
    const handleNodeActions = (id: string, type: string) => {
        /* TODO: Handle Node Action */
    };

    // eslint-disable-next-line no-unused-vars
    const handleAddNode = (value: any) => {
        setIsAddNode(true);
    };
    const onUpdateAllNodes = () => {
        /* TODO: Handle Node Updates */
    };
    const handleActivationSubmit = () => {
        /* Handle submit activation action */
    };
    const onActivateUser = () => setIsUserActivateOpen(() => true);

    // eslint-disable-next-line no-unused-vars
    const handleNodeUpdateActin = (id: string) => {
        setIsSoftwaUpdate(true);
        /* Handle node update  action */
    };

    return (
        <Box component="div" sx={{ flexGrow: 1, pb: "18px" }}>
            <Grid container spacing={3}>
                <Grid xs={12} item>
                    <NetworkStatus
                        handleAddNode={handleAddNode}
                        handleActivateUser={onActivateUser}
                        loading={networkStatusLoading || isSkeltonLoad}
                        statusType={networkStatusRes?.getNetwork?.status || ""}
                        duration={
                            networkStatusRes?.getNetwork?.description || ""
                        }
                    />
                </Grid>
                <Grid item xs={12} md={6} lg={4}>
                    <StatusCard
                        Icon={UsersWithBG}
                        title={"Connected Users"}
                        options={TIME_FILTER}
                        subtitle1={`${
                            connectedUserRes?.getConnectedUsers?.totalUser || 0
                        }`}
                        subtitle2={""}
                        option={getStatus("statusUser")}
                        loading={connectedUserloading || isSkeltonLoad}
                        handleSelect={(value: string) =>
                            handleSatusChange("statusUser", value)
                        }
                    />
                </Grid>
                <Grid item xs={12} md={6} lg={4}>
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
                <Grid item xs={12} md={6} lg={4}>
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
                <Grid xs={12} item>
                    <StatsCard
                        hasMetricData={uptimeMetric.data.length > 0}
                        metricData={uptimeMetric}
                        loading={
                            nodeLoading || isSkeltonLoad || metricUptimeLoading
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
                                title="My Nodes"
                                showButton={isContainNodeUpdate(
                                    nodeRes?.getNodesByOrg.nodes
                                )}
                                buttonSize={"small"}
                                buttonTitle={"Update All"}
                                stats={`${
                                    nodeRes?.getNodesByOrg.activeNodes || "0"
                                }/${nodeRes?.getNodesByOrg.totalNodes || "-"}`}
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
                                dataset={residentsRes?.getUsersByOrg.users}
                                menuOptions={DEACTIVATE_EDIT_ACTION_MENU}
                                onMenuItemClick={onResidentsTableMenuItem}
                            />
                        </RoundedCard>
                    </LoadingWrapper>
                </Grid>
            </Grid>
            <UserActivationDialog
                isOpen={showSimActivationDialog}
                dialogTitle={t("DIALOG_MESSAGE.SimActivationDialogTitle")}
                subTitle={t("DIALOG_MESSAGE.SimActivationDialogContent")}
                handleClose={handleSimActivateClose}
            />
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
            <BasicDialog
                isClosable={false}
                btnVariant="contained"
                isOpen={isWelcomeDialog}
                title={"Welcome to Ukama Console!"}
                content={
                    "This is where you can manage your network, and troubleshoot things, if necessary. For now, while your nodes have not shipped, you can monitor your users’ data usage, and [insert other main use]. "
                }
                btnLabel={"continue to console"}
                handleClose={() => setIsWelcomeDialog(false)}
            />
            {isUserActivateOpen && (
                <UserActivationDialog
                    isOpen={isUserActivateOpen}
                    dialogTitle={UserActivation.title}
                    subTitle={UserActivation.subTitle}
                    handleClose={handleUserActivateClose}
                />
            )}
            {isAddNode && (
                <ActivationDialog
                    isOpen={isAddNode}
                    dialogTitle={"Register Node"}
                    handleClose={handleAddNodeClose}
                    handleActivationSubmit={handleActivationSubmit}
                    subTitle={
                        "Ensure node is properly set up in desired location before completing this step. Enter serial number found in your confirmation email, or on the back of your node, and we’ll take care of the rest for you."
                    }
                />
            )}
        </Box>
    );
};
export default Home;
