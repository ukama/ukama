import { useState } from "react";
import { Box, Grid } from "@mui/material";
import { NETWORKS } from "../../constants";
import { StatusCard, NetworkStatus } from "../../components";
import { DashboardStatusCard } from "../../constants/stubData";

const Home = () => {
    const [network, setNetwork] = useState("public");
    const [userStatusFilter, setUserStatusFilter] = useState("total");
    const [dataStatusFilter, setDataStatusFilter] = useState("total");
    const [billingStatusFilter, setBillingStatusFilter] = useState("july");

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
        <Box>
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
        </Box>
    );
};

export default Home;
