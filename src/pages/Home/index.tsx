import React, { useState } from "react";
import { Box, Grid } from "@mui/material";
import { RoundedCard } from "../../styles";
import {
    NodeCard,
    StatusCard,
    NetworkStatus,
    ContainerHeader,
    StatsCard,
    AlertCard,
} from "../../components";
import { DashboardStatusCard } from "../../constants/stubData";
import {
    STATS_OPTIONS,
    STATS_PERIOD,
    NETWORKS,
    ALERT_INFORMATION,
} from "../../constants";
const Home = () => {
    const [network, setNetwork] = useState("public");
    const [userStatusFilter, setUserStatusFilter] = useState("total");
    const [dataStatusFilter, setDataStatusFilter] = useState("total");
    const [billingStatusFilter, setBillingStatusFilter] = useState("july");
    const [statOptionValue, setstatOptionValue] = React.useState(3);

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

    return (
        <>
            <Box sx={{ flexGrow: 1 }}>
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
                            />
                        </Grid>

                        <Grid xs={12} item md={4} sm={12}>
                            <AlertCard alertCardItems={ALERT_INFORMATION} />
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
                                handleButtonAction={() => {}}
                            />
                            <NodeCard isConfigure={true} />
                        </RoundedCard>
                    </Grid>
                    <Grid xs={12} md={4} item>
                        <RoundedCard sx={{ height: "100%" }}></RoundedCard>
                    </Grid>
                </Grid>
            </Box>
        </>
    );
};
export default Home;
