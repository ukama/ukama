import React, { useState } from "react";
import { Box, Typography, Grid, SelectChangeEvent } from "@mui/material";
import {
    StatusCard,
    NetworkStatus,
    StatsCard,
    AlertCard,
} from "../../components";
import { DashboardStatusCard } from "../../constants/stubData";
import { STATS_OPTIONS, STATS_PERIOD, NETWORKS } from "../../constants";
import { CloudOffIcon } from "../../assets/svg";
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
                <Box>
                    <NetworkStatus
                        duration={""}
                        option={network}
                        options={NETWORKS}
                        statusType={"IN_PROGRESS"}
                        status={"Your network is being configured"}
                        handleStatusChange={(value: string) =>
                            setNetwork(value)
                        }
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
                    <Box mt={2}>
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
                                <AlertCard
                                    title="Software update"
                                    subheader="Short description of alert."
                                    action={
                                        <Typography variant="caption">
                                            08/16/21 1PM
                                        </Typography>
                                    }
                                    Icon={<CloudOffIcon />}
                                />
                            </Grid>
                        </Grid>
                    </Box>
                </Box>
            </Box>
        </>
    );
};
export default Home;
