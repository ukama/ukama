import {
    NodeCard,
    StatusCard,
    NetworkStatus,
    ContainerHeader,
    StatsCard,
    AlertCard,
    MultiSlideCarousel,
    DataTableWithOptions,
} from "../../components";
import {
    DashboardSliderData,
    DashboardStatusCard,
    DashboardResidentsTable,
    ALERT_INFORMATION,
} from "../../constants/stubData";
import {
    NETWORKS,
    STATS_OPTIONS,
    STATS_PERIOD,
    DEACTIVATE_EDIT_ACTION_MENU,
    DataTableWithOptionColumns,
} from "../../constants";
import "../../i18n/i18n";
import React, { useState } from "react";
import { RoundedCard } from "../../styles";
import { useTranslation } from "react-i18next";
import {
    Box,
    Grid,
    List,
    ListItem,
    Typography,
    useMediaQuery,
} from "@mui/material";
import { AlertItemType } from "../../types";

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
    const isSliderLarge = useMediaQuery("(min-width:1430px)");
    const isSliderMedium = useMediaQuery("(min-width:1160px)") ? 2 : 1;
    const slidesToShow = isSliderLarge ? 3 : isSliderMedium;
    const [isAddNode, setIsAddNode] = useState(false);
    const [network, setNetwork] = useState("public");
    const [userStatusFilter, setUserStatusFilter] = useState("total");
    const [dataStatusFilter, setDataStatusFilter] = useState("total");
    const [billingStatusFilter, setBillingStatusFilter] = useState("july");
    const [statOptionValue, setstatOptionValue] = React.useState(3);
    const [selectedBtn, setSelectedBtn] = useState("DAY");
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
    const onActivateButton = () => {};
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
    const getNodesContainerData = (items: any[], slidesToShow: number) =>
        items.length > 3 ? (
            <MultiSlideCarousel numberOfSlides={slidesToShow}>
                {items.map(({ id, title, users, subTitle, isConfigure }) => (
                    <NodeCard
                        key={id}
                        title={title}
                        users={users}
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
                        <Grid xs={12} item sm={12} md={8}>
                            <StatsCard
                                selectOption={statOptionValue}
                                options={STATS_OPTIONS}
                                periodOptions={STATS_PERIOD}
                                handleSelect={handleStatsChange}
                                handleSelectedButton={
                                    handleSelectedButtonChange
                                }
                                selectedButton={selectedBtn}
                            />
                        </Grid>

                        <Grid xs={12} item md={4}>
                            <RoundedCard>
                                <Typography variant="h6" sx={{ mb: "14px" }}>
                                    {t("ALERT.Title")}
                                </Typography>
                                <List
                                    sx={{
                                        p: "0px",
                                        maxHeight: 300,
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
                                                    description={description}
                                                />
                                            </ListItem>
                                        )
                                    )}
                                </List>
                            </RoundedCard>
                        </Grid>
                    </Grid>
                </Box>

                <Grid container spacing={2}>
                    <Grid xs={12} md={8} item>
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
                    </Grid>
                    <Grid xs={12} md={4} item>
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
                    </Grid>
                </Grid>
            </Box>
        </>
    );
};
export default Home;
