import {
    TableHeader,
    SimpleDataTable,
    CurrentBill,
    LoadingWrapper,
    BillingAlerts,
} from "../../components";
import "../../i18n/i18n";
import React, { useState } from "react";
import { useRecoilValue } from "recoil";
import { isSkeltonLoading } from "../../recoil";
import { CenterContainer, RoundedCard } from "../../styles";
import { Box, Grid, Tabs, Typography, Tab, AlertColor } from "@mui/material";
import { CurrentBillColumns } from "../../constants/tableColumns";
import { BillingTabs, CurrentBillingData } from "../../constants";

const Billing = () => {
    const [billingAlert, setBillingAlert] = useState({
        type: "info",
        btnText: "Enter now →",
        title: "Set up your payment information securely at any time.",
    });
    const [tab, setTab] = useState<number>(0);
    const _isSkeltonLoading = useRecoilValue(isSkeltonLoading);
    const [selectedRows, setSelectedRows] = useState<number[]>([]);
    const handleMakePayment = () => {
        /* TODO: Handle make payment action */
    };
    const handleTabChange = (event: React.SyntheticEvent, value: any) =>
        setTab(value);
    const handleExport = () => {
        /* TODO: Handle export action */
    };

    const handleAlertAction = () => {
        /* TODO: Handle Alert notification action */
        setBillingAlert(prev => ({ ...prev, type: "error" }));
    };

    return (
        <Box>
            <BillingAlerts
                title={billingAlert.title}
                btnText={billingAlert.btnText}
                onActionClick={handleAlertAction}
                type={billingAlert.type as AlertColor}
            />
            <LoadingWrapper isLoading={_isSkeltonLoading} height={"90%"}>
                <Box component="div">
                    <Tabs
                        value={tab}
                        onChange={handleTabChange}
                        sx={{ mt: 2, mb: 4 }}
                    >
                        {BillingTabs.map(({ id, label, value }) => (
                            <Tab
                                key={id}
                                label={label}
                                id={`billing-tab-${value}`}
                            />
                        ))}
                    </Tabs>

                    {tab === 0 && (
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
                    {tab === 1 && (
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
            </LoadingWrapper>
        </Box>
    );
};

export default Billing;
