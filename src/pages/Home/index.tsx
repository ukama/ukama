/* eslint-disable no-unused-vars */
import {
    StatsCard,
    AlertCard,
    StatusCard,
    NetworkStatus,
    NodeContainer,
    LoadingWrapper,
    ContainerHeader,
    DataTableWithOptions,
    UserActivationDialog,
} from "../../components";
import {
    NETWORKS,
    TIME_FILTER,
    MONTH_FILTER,
    STATS_PERIOD,
    STATS_OPTIONS,
    UserActivation,
    DEACTIVATE_EDIT_ACTION_MENU,
    DataTableWithOptionColumns,
} from "../../constants";
import "../../i18n/i18n";
import { Box, Grid, Typography, useMediaQuery } from "@mui/material";
import {
    Time_Filter,
    Network_Type,
    Data_Bill_Filter,
    useGetDataBillQuery,
    useGetDataUsageQuery,
    useGetConnectedUsersQuery,
    useGetNodesQuery,
    useGetNetworkQuery,
    useGetAlertsQuery,
    useGetResidentsQuery,
    useDeleteNodeMutation,
    AddNodeDto,
    useAddNodeMutation,
    useUpdateNodeMutation,
    UpdateNodeDto,
} from "../../generated";
import React, { useState } from "react";
import { useRecoilValue } from "recoil";
import { RoundedCard } from "../../styles";
import { useTranslation } from "react-i18next";
import { isSkeltonLoading } from "../../recoil";
import { DataBilling, DataUsage, UsersWithBG } from "../../assets/svg";

const Home = () => {
    const { t } = useTranslation();
    const isSliderLarge = useMediaQuery("(min-width:1430px)");
    const isSliderMedium = useMediaQuery("(min-width:1160px)") ? 2 : 1;
    const slidesToShow = isSliderLarge ? 3 : isSliderMedium;
    const [selectedBtn, setSelectedBtn] = useState("DAY");
    const isSkeltonLoad = useRecoilValue(isSkeltonLoading);
    const [statOptionValue, setstatOptionValue] = useState(3);
    const [isUserActivateOpen, setIsUserActivateOpen] = useState(false);
    const [userStatusFilter, setUserStatusFilter] = useState(Time_Filter.Total);
    const [dataStatusFilter, setDataStatusFilter] = useState(Time_Filter.Total);
    const [networkType, setNetworkType] = useState<Network_Type>(
        Network_Type.Public
    );
    const [billingStatusFilter, setBillingStatusFilter] = useState(
        Data_Bill_Filter.July
    );
    // eslint-disable-next-line no-unused-vars
    const [deleteNode, { data: deleNodeRes }] = useDeleteNodeMutation();
    const [addNode, { data: addNodeRes }] = useAddNodeMutation();
    const [editNode, { data: editNodeRes }] = useUpdateNodeMutation();
    const { data: connectedUserRes, loading: connectedUserloading } =
        useGetConnectedUsersQuery({
            variables: {
                filter: userStatusFilter,
            },
        });
    const { data: alertsInfoRes, loading: alertsloading } = useGetAlertsQuery({
        variables: {
            data: {
                pageNo: 1,
                pageSize: 50,
            },
        },
    });

    const { data: residentsRes, loading: residentsloading } =
        useGetResidentsQuery({
            variables: {
                data: {
                    pageNo: 1,
                    pageSize: 50,
                },
            },
        });
    const { data: dataUsageRes, loading: dataUsageloading } =
        useGetDataUsageQuery({
            variables: {
                filter: dataStatusFilter,
            },
        });

    const { data: dataBillingRes, loading: dataBillingloading } =
        useGetDataBillQuery({
            variables: {
                filter: billingStatusFilter,
            },
        });

    const {
        data: nodeRes,
        loading: nodeLoading,
        refetch: refetchNodes,
    } = useGetNodesQuery({
        variables: {
            data: {
                pageNo: 1,
                pageSize: 50,
            },
        },
    });

    const { data: networkStatusRes, loading: networkStatusLoading } =
        useGetNetworkQuery({
            variables: {
                filter: networkType,
            },
        });

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

    const onActivateButton = () => setIsUserActivateOpen(() => true);
    const handleUserActivateClose = () => setIsUserActivateOpen(() => false);
    const onResidentsTableMenuItem = (id: string, type: string) => {};
    const handleNodeActions = (id: string, type: string) => {
        if (type === "delete") {
            deleteNode({
                variables: { id },
            });
            refetchNodes();
        }
    };
    const handleAddNode = (data: AddNodeDto) => {
        addNode({
            variables: {
                data,
            },
        });
    };
    const handleEditNode = (data: UpdateNodeDto) => {
        editNode({
            variables: {
                data,
            },
        });
    };
    return (
        <>
            <Box sx={{ flexGrow: 1, pb: "18px" }}>
                <NetworkStatus
                    options={NETWORKS}
                    option={networkType}
                    loading={networkStatusLoading}
                    statusType={networkStatusRes?.getNetwork?.status || ""}
                    duration={networkStatusRes?.getNetwork?.description || ""}
                    handleStatusChange={(value: Network_Type) =>
                        setNetworkType(value)
                    }
                />
                <Grid container spacing={2}>
                    <Grid item xs={12} md={6} lg={4}>
                        <StatusCard
                            Icon={UsersWithBG}
                            title={"Connected Users"}
                            options={TIME_FILTER}
                            subtitle1={`${
                                connectedUserRes?.getConnectedUsers
                                    ?.totalUser || 0
                            }`}
                            subtitle2={
                                connectedUserRes?.getConnectedUsers?.totalUser
                                    ? `| ${connectedUserRes?.getConnectedUsers?.residentUsers} residents; ${connectedUserRes?.getConnectedUsers?.guestUsers} guests`
                                    : ""
                            }
                            option={getStatus("statusUser")}
                            loading={connectedUserloading}
                            handleSelect={(value: string) =>
                                handleSatusChange("statusUser", value)
                            }
                        />
                    </Grid>
                    <Grid item xs={12} md={6} lg={4}>
                        <StatusCard
                            title={"Data usage"}
                            subtitle1={`${
                                dataUsageRes?.getDataUsage?.dataConsumed || 0
                            }`}
                            subtitle2={` GBs / ${
                                dataUsageRes?.getDataUsage?.dataPackage ||
                                "unlimited"
                            }`}
                            Icon={DataUsage}
                            options={TIME_FILTER}
                            option={getStatus("statusUsage")}
                            loading={dataUsageloading}
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
                                    ? ` / due in ${dataBillingRes?.getDataBill?.billDue}`
                                    : " due"
                            }
                            Icon={DataBilling}
                            options={MONTH_FILTER}
                            loading={dataBillingloading}
                            option={getStatus("statusBill")}
                            handleSelect={(value: string) =>
                                handleSatusChange("statusBill", value)
                            }
                        />
                    </Grid>
                </Grid>
                <Box mt={2} mb={2}>
                    <Grid container spacing={2}>
                        <Grid xs={12} item lg={8}>
                            <StatsCard
                                loading={isSkeltonLoad}
                                options={STATS_OPTIONS}
                                selectedButton={selectedBtn}
                                periodOptions={STATS_PERIOD}
                                selectOption={statOptionValue}
                                handleSelect={handleStatsChange}
                                handleSelectedButton={
                                    handleSelectedButtonChange
                                }
                            />
                        </Grid>
                        <Grid xs={12} item lg={4}>
                            <LoadingWrapper
                                height={387}
                                isLoading={alertsloading}
                            >
                                <RoundedCard>
                                    <Typography
                                        variant="h6"
                                        sx={{ mb: "14px" }}
                                    >
                                        {t("ALERT.Title")}
                                    </Typography>
                                    <AlertCard
                                        alertOptions={
                                            alertsInfoRes?.getAlerts?.alerts
                                        }
                                    />
                                </RoundedCard>
                            </LoadingWrapper>
                        </Grid>
                    </Grid>
                </Box>

                <Grid container spacing={2}>
                    <Grid xs={12} lg={8} item>
                        <LoadingWrapper height={312} isLoading={nodeLoading}>
                            <RoundedCard>
                                <ContainerHeader
                                    title="My Nodes"
                                    buttonTitle="Add Node"
                                    handleButtonAction={() => {}}
                                    stats={`${
                                        nodeRes?.getNodes?.nodes?.activeNodes ||
                                        "0"
                                    }/${
                                        nodeRes?.getNodes?.nodes?.totalNodes ||
                                        "-"
                                    }`}
                                />
                                <NodeContainer
                                    slidesToShow={slidesToShow}
                                    items={nodeRes?.getNodes?.nodes.nodes}
                                    count={nodeRes?.getNodes.meta.size}
                                    handleItemAction={handleNodeActions}
                                />
                            </RoundedCard>
                        </LoadingWrapper>
                    </Grid>
                    <Grid xs={12} lg={4} item>
                        <LoadingWrapper
                            height={337}
                            isLoading={residentsloading}
                        >
                            <RoundedCard sx={{ height: "100%" }}>
                                <ContainerHeader
                                    title="Residents"
                                    buttonTitle="ACTIVATE"
                                    handleButtonAction={onActivateButton}
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
            </Box>
        </>
    );
};
export default Home;
