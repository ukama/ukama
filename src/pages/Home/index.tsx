import {
    NodeCard,
    StatusCard,
    NetworkStatus,
    ContainerHeader,
    StatsCard,
    AlertCard,
    MultiSlideCarousel,
    DataTableWithOptions,
    UserActivationDialog,
    LoadingWrapper,
} from "../../components";
import {
    ALERT_INFORMATION,
    DashboardSliderData,
    DashboardStatusCard,
    DashboardResidentsTable,
} from "../../constants/stubData";
import {
    NETWORKS,
    STATS_PERIOD,
    STATS_OPTIONS,
    UserActivation,
    DEACTIVATE_EDIT_ACTION_MENU,
    DataTableWithOptionColumns,
} from "../../constants";
import "../../i18n/i18n";
import {
    Box,
    Grid,
    List,
    ListItem,
    Typography,
    useMediaQuery,
} from "@mui/material";
import React, { useState } from "react";
import { useRecoilValue } from "recoil";
import { RoundedCard } from "../../styles";
import { AlertItemType } from "../../types";
import { useTranslation } from "react-i18next";
import { isSkeltonLoading } from "../../recoil";

let slides = [
    {
        id: 1,
        title: "",
        subTitle: "",
        users: "",
        isConfigure: true,
    },
];

const Home = () => {
    const { t } = useTranslation();
    const isSliderLarge = useMediaQuery("(min-width:1410px)");
    const isSliderMedium = useMediaQuery("(min-width:1160px)") ? 2 : 1;
    const slidesToShow = isSliderLarge ? 3 : isSliderMedium;
    const [network, setNetwork] = useState("public");
    const [isAddNode, setIsAddNode] = useState(false);
    const [selectedBtn, setSelectedBtn] = useState("DAY");
    const isSkeltonLoad = useRecoilValue(isSkeltonLoading);
    const [statOptionValue, setstatOptionValue] = React.useState(3);
    const [userStatusFilter, setUserStatusFilter] = useState("current");
    const [dataStatusFilter, setDataStatusFilter] = useState("month");
    const [isUserActivateOpen, setIsUserActivateOpen] = useState(false);
    const [billingStatusFilter, setBillingStatusFilter] = useState("july");

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
                return setUserStatusFilter(value);
            case "statusUsage":
                return setDataStatusFilter(value);
            case "statusBill":
                return setBillingStatusFilter(value);
        }
    };
    const onResidentsTableMenuItem = () => {};
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

    const getNodesContainerData = (items: any[], slidesToShow: number) =>
        items.length > 3 ? (
            <MultiSlideCarousel numberOfSlides={slidesToShow}>
                {items.map(({ id, title, users, subTitle, isConfigure }) => (
                    <NodeCard
                        key={id}
                        title={title}
                        users={users}
                        loading={false}
                        subTitle={subTitle}
                        isConfigure={isConfigure}
                    />
                ))}
            </MultiSlideCarousel>
        ) : (
            <Grid
                item
                xs={12}
                container
                spacing={6}
                sx={{
                    display: "flex",
                    justifyContent: { xs: "center", md: "flex-start" },
                }}
            >
                {items.map(i => (
                    <Grid key={i} item>
                        <NodeCard isConfigure={true} />
                    </Grid>
                ))}
            </Grid>
        );

    return (
        <>
            <Box sx={{ flexGrow: 1, pb: "18px" }}>
                <NetworkStatus
                    duration={""}
                    option={network}
                    options={NETWORKS}
                    loading={isSkeltonLoad}
                    statusType={"IN_PROGRESS"}
                    status={"Your network is being configured"}
                    handleStatusChange={(value: string) => setNetwork(value)}
                />
                <Grid container spacing={2}>
                    {DashboardStatusCard.map(
                        ({
                            id,
                            Icon,
                            title,
                            options,
                            subtitle1,
                            subtitle2,
                        }: any) => (
                            <Grid key={id} item xs={12} md={6} lg={4}>
                                <StatusCard
                                    Icon={Icon}
                                    title={title}
                                    options={options}
                                    subtitle1={subtitle1}
                                    subtitle2={subtitle2}
                                    option={getStatus(id)}
                                    loading={isSkeltonLoad}
                                    handleSelect={(value: string) =>
                                        handleSatusChange(id, value)
                                    }
                                />
                            </Grid>
                        )
                    )}
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
                                height={337}
                                isLoading={isSkeltonLoad}
                            >
                                <RoundedCard>
                                    <Typography
                                        variant="h6"
                                        sx={{ mb: "14px" }}
                                    >
                                        {t("ALERT.Title")}
                                    </Typography>
                                    <List
                                        sx={{
                                            pr: "4px",
                                            maxHeight: 305,
                                            overflow: "auto",
                                            position: "relative",
                                        }}
                                    >
                                        {ALERT_INFORMATION.map(
                                            ({
                                                id,
                                                date,
                                                description,
                                                title,
                                                Icon,
                                            }: AlertItemType) => (
                                                <ListItem
                                                    key={id}
                                                    style={{
                                                        padding: 1,
                                                        marginBottom: "4px",
                                                    }}
                                                >
                                                    <AlertCard
                                                        id={id}
                                                        date={date}
                                                        Icon={Icon}
                                                        title={title}
                                                        description={
                                                            description
                                                        }
                                                    />
                                                </ListItem>
                                            )
                                        )}
                                    </List>
                                </RoundedCard>
                            </LoadingWrapper>
                        </Grid>
                    </Grid>
                </Box>

                <Grid container spacing={2}>
                    <Grid xs={12} lg={8} item>
                        <LoadingWrapper height={312} isLoading={isSkeltonLoad}>
                            <RoundedCard>
                                <ContainerHeader
                                    stats="1/8"
                                    title="My Nodes"
                                    buttonTitle="Add Node"
                                    handleButtonAction={() =>
                                        setIsAddNode(prev => !prev)
                                    }
                                />
                                {getNodesContainerData(
                                    isAddNode ? DashboardSliderData : slides,
                                    slidesToShow
                                )}
                            </RoundedCard>
                        </LoadingWrapper>
                    </Grid>
                    <Grid xs={12} lg={4} item>
                        <LoadingWrapper height={312} isLoading={isSkeltonLoad}>
                            <RoundedCard sx={{ height: "100%" }}>
                                <ContainerHeader
                                    stats="6/16"
                                    title="Residents"
                                    buttonTitle="ACTIVATE"
                                    handleButtonAction={onActivateButton}
                                />
                                <DataTableWithOptions
                                    columns={DataTableWithOptionColumns}
                                    dataset={DashboardResidentsTable}
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
