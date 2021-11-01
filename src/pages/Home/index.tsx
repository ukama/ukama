import { useState } from "react";
import { Grid } from "@mui/material";
import { RoundedCard } from "../../styles";
import {
    NodeCard,
    StatusCard,
    NetworkStatus,
    ContainerHeader,
    DataTableWithOptions,
} from "../../components";
import {
    DashboardResidentsTable,
    DashboardStatusCard,
} from "../../constants/stubData";
import {
    NETWORKS,
    DEACTIVATE_EDIT_ACTION_MENU,
    DataTableWithOptionColumns,
} from "../../constants";

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
        <>
            <NetworkStatus
                duration={""}
                option={network}
                options={NETWORKS}
                statusType={"IN_PROGRESS"}
                status={"Your network is being configured"}
                handleStatusChange={(value: string) => setNetwork(value)}
            />
            <Grid container spacing={2}>
                <Grid xs={12} item container spacing={2}>
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
                    <RoundedCard sx={{ height: "100%" }}>
                        <ContainerHeader
                            stats="6/16"
                            title="Residents"
                            buttonTitle="ACTIVATE"
                            handleButtonAction={() => {}}
                        />
                        <DataTableWithOptions
                            columns={DataTableWithOptionColumns}
                            dataset={DashboardResidentsTable}
                            menuOptions={DEACTIVATE_EDIT_ACTION_MENU}
                            onMenuItemClick={() => {}}
                        />
                    </RoundedCard>
                </Grid>
            </Grid>
        </>
    );
};

export default Home;
