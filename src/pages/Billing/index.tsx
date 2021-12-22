import {
    TabLayout,
    TableHeader,
    SimpleDataTable,
    CurrentBill,
} from "../../components";
import "../../i18n/i18n";
import { useState } from "react";
import { BillingTabs } from "../../constants";
import { Box, Grid, Typography } from "@mui/material";
import { CenterContainer, RoundedCard } from "../../styles";
import { CurrentBillingData } from "../../constants/stubData";
import { CurrentBillColumns } from "../../constants/tableColumns";

const Billing = () => {
    const [tab, setTab] = useState("1");
    const [selectedRows, setSelectedRows] = useState<number[]>([]);
    const handleMakePayment = () => {
        /* TODO: Handle make payment action */
    };
    const handleTabChange = (value: string) => setTab(value);
    const handleExport = () => {
        /* TODO: Handle export action */
    };

    return (
        <Box>
            <TabLayout
                tab={tab}
                tabs={BillingTabs}
                onTabChange={handleTabChange}
            />

            {tab === "1" && (
                <Grid container item spacing={2}>
                    <Grid xs={12} md={5} item>
                        <CurrentBill
                            amount={"$20.00"}
                            title={"Amount due"}
                            dueDate={"*Due 11/10/2021"}
                            periodOf={"10/10/2021 - 11/10/2021"}
                            handleMakePayment={handleMakePayment}
                        />
                    </Grid>
                    <Grid xs={12} md={7} item>
                        <RoundedCard>
                            <CenterContainer>
                                Under Developement :)
                            </CenterContainer>
                        </RoundedCard>
                    </Grid>
                    <Grid xs={12} item>
                        <RoundedCard>
                            <TableHeader
                                title={"Billing breakdown"}
                                buttonTitle={"Export"}
                                handleButtonAction={handleExport}
                            />
                            <SimpleDataTable
                                columns={CurrentBillColumns}
                                dataset={CurrentBillingData}
                            />

                            <div
                                style={{
                                    width: "100%",
                                    display: "flex",
                                    margin: "18px 0px",
                                    justifyContent: "flex-end",
                                }}
                            >
                                <Typography
                                    variant={"h6"}
                                    sx={{
                                        width: "20%",
                                    }}
                                >
                                    $20.00
                                </Typography>
                            </div>
                        </RoundedCard>
                    </Grid>
                </Grid>
            )}
            {tab === "2" && (
                <RoundedCard>
                    <TableHeader
                        title={"Billing history"}
                        buttonTitle={"Export"}
                        handleButtonAction={handleExport}
                    />
                    <SimpleDataTable
                        rowSelection={true}
                        selectedRows={selectedRows}
                        columns={CurrentBillColumns}
                        dataset={CurrentBillingData}
                        setSelectedRows={setSelectedRows}
                        totalRows={CurrentBillingData.length}
                    />
                </RoundedCard>
            )}
        </Box>
    );
};

export default Billing;
