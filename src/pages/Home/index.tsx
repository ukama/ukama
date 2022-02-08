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
} from "../../components";
import {
    TIME_FILTER,
    MONTH_FILTER,
    STATS_OPTIONS,
    UserActivation,
    DataTableWithOptionColumns,
    DEACTIVATE_EDIT_ACTION_MENU,
} from "../../constants";
import "../../i18n/i18n";
import {
    Time_Filter,
    Network_Type,
    Data_Bill_Filter,
    useGetNetworkQuery,
    useGetDataBillQuery,
    useGetDataUsageQuery,
    useGetResidentsQuery,
    useGetNodesByOrgQuery,
    GetLatestNetworkDocument,
    useDeactivateUserMutation,
    useGetConnectedUsersQuery,
    GetLatestDataBillDocument,
    GetLatestDataUsageDocument,
    GetLatestNetworkSubscription,
    GetLatestDataBillSubscription,
    GetLatestDataUsageSubscription,
    GetLatestConnectedUsersDocument,
    GetLatestConnectedUsersSubscription,
} from "../../generated";
import { useRecoilValue } from "recoil";
import { RoundedCard } from "../../styles";
import React, { useEffect, useState } from "react";
import { Box, Grid, useMediaQuery } from "@mui/material";
import { isSkeltonLoading, organizationId } from "../../recoil";
import { DataBilling, DataUsage, UsersWithBG } from "../../assets/svg";

const Home = () => {
    const isSliderLarge = useMediaQuery("(min-width:1410px)");
    const isSliderMedium = useMediaQuery("(min-width:1160px)") ? 2 : 1;
    const slidesToShow = isSliderLarge ? 3 : isSliderMedium;
    const [selectedBtn, setSelectedBtn] = useState("DAY");
    const isSkeltonLoad = useRecoilValue(isSkeltonLoading);
    const orgId = useRecoilValue(organizationId);
    const [statOptionValue, setstatOptionValue] = useState(3);
    const [isUserActivateOpen, setIsUserActivateOpen] = useState(false);
    const [userStatusFilter, setUserStatusFilter] = useState(Time_Filter.Total);
    const [dataStatusFilter, setDataStatusFilter] = useState(Time_Filter.Month);
    const [isAddNode, setIsAddNode] = useState(false);
    const [billingStatusFilter, setBillingStatusFilter] = useState(
        Data_Bill_Filter.July
    );

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
    useEffect(() => {
        if (connectedUserRes) {
            subscribeToLatestConnectedUsers<GetLatestConnectedUsersSubscription>(
                {
                    document: GetLatestConnectedUsersDocument,
                    updateQuery: (prev, { subscriptionData }) => {
                        let data = { ...prev };
                        const latestConnectedUser =
                            subscriptionData.data.getConnectedUsers;
                        if (
                            latestConnectedUser.__typename ===
                            "ConnectedUserDto"
                        )
                            data.getConnectedUsers = latestConnectedUser;
                        return data;
                    },
                }
            );
        }
    }, [connectedUserRes]);

    const {
        data: residentsRes,
        loading: residentsloading,
        refetch: refetchUser,
    } = useGetResidentsQuery({
        variables: {
            data: {
                pageNo: 1,
                pageSize: 50,
            },
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
    useEffect(() => {
        if (dataUsageRes) {
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
        }
    }, [dataUsageRes]);

    const {
        data: dataBillingRes,
        loading: dataBillingloading,
        subscribeToMore: subscribeToLatestDataBill,
    } = useGetDataBillQuery({
        variables: {
            filter: billingStatusFilter,
        },
    });
    useEffect(() => {
        if (dataBillingRes) {
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
        }
    }, [dataBillingRes]);

    const { data: nodeRes, loading: nodeLoading } = useGetNodesByOrgQuery({
        variables: { orgId: orgId || "" },
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

    const handleSelectedButtonChange = (
        event: React.MouseEvent<HTMLElement>,
        newSelected: string
    ) => {
        setSelectedBtn(newSelected);
    };

    const handleStatsChange = (event: {
        target: { value: React.SetStateAction<number> };
    }) => {
        setstatOptionValue(event.target.value);
    };

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

    const handleActivationSubmit = () => {
        /* Handle submit activation action */
    };
    const onActivateUser = () => setIsUserActivateOpen(() => true);

    return (
        <Box sx={{ flexGrow: 1, pb: "18px" }}>
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
                        loading={isSkeltonLoad}
                        options={STATS_OPTIONS}
                        selectedButton={selectedBtn}
                        selectOption={statOptionValue}
                        handleSelect={handleStatsChange}
                        handleSelectedButton={handleSelectedButtonChange}
                    />
                </Grid>
                <Grid xs={12} lg={8} item>
                    <LoadingWrapper
                        height={312}
                        isLoading={nodeLoading || isSkeltonLoad}
                    >
                        <RoundedCard>
                            <ContainerHeader
                                title="My Nodes"
                                showButton={false}
                                stats={`${
                                    nodeRes?.getNodesByOrg.activeNodes || "0"
                                }/${nodeRes?.getNodesByOrg.totalNodes || "-"}`}
                            />
                            <NodeContainer
                                slidesToShow={slidesToShow}
                                items={nodeRes?.getNodesByOrg.nodes}
                                count={nodeRes?.getNodesByOrg.nodes.length}
                                handleItemAction={handleNodeActions}
                            />
                        </RoundedCard>
                    </LoadingWrapper>
                </Grid>
                <Grid xs={12} lg={4} item>
                    <LoadingWrapper
                        height={312}
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
                                stats={`${
                                    residentsRes?.getResidents?.residents
                                        ?.activeResidents || "0"
                                }/${
                                    residentsRes?.getResidents?.residents
                                        ?.totalResidents || "-"
                                }`}
                            />
                            <DataTableWithOptions
                                columns={DataTableWithOptionColumns}
                                dataset={
                                    residentsRes?.getResidents.residents
                                        .residents
                                }
                                menuOptions={DEACTIVATE_EDIT_ACTION_MENU}
                                onMenuItemClick={onResidentsTableMenuItem}
                            />
                        </RoundedCard>
                    </LoadingWrapper>
                </Grid>
            </Grid>

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
                    dialogTitle={"Add Node"}
                    subTitle2={
                        "To confirm this node is yours, we have emailed you a security code that  expires in 15 minutes."
                    }
                    handleClose={handleAddNodeClose}
                    handleActivationSubmit={handleActivationSubmit}
                    subTitle={
                        "Add more nodes to expand your network coverage. Enter the serial number found in your purchase confirmation email, and it will be automatically configured."
                    }
                />
            )}
        </Box>
    );
};
export default Home;
