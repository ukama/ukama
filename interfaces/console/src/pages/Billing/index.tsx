import {
    TableHeader,
    SimpleDataTable,
    CurrentBill,
    LoadingWrapper,
    BillingAlerts,
    BillingDialog,
    PaymentCard,
} from "../../components";
import "../../i18n/i18n";
import React, { useState } from "react";
import { useRecoilValue } from "recoil";
import { isSkeltonLoading } from "../../recoil";
import { RoundedCard } from "../../styles";
import { Box, Grid, Tabs, Typography, Tab, AlertColor } from "@mui/material";
import { CurrentBillColumns } from "../../constants/tableColumns";
import { BillingTabs, CurrentBillingData } from "../../constants";
import { PaymentCards } from "../../constants/stubData";

const Billing = () => {
    const [isBilling, setIsBilling] = useState(false);
    const [billingAlert, setBillingAlert] = useState({
        type: "info",
        btnText: "Enter now →",
        title: "Set up your payment information securely at any time.",
    });
    const [tab, setTab] = useState<number>(0);
    const _isSkeltonLoading = useRecoilValue(isSkeltonLoading);
    const [selectedRows, setSelectedRows] = useState<number[]>([]);

    const handleTabChange = (event: React.SyntheticEvent, value: any) =>
        setTab(value);

    const handleAlertAction = () => {
        setBillingAlert(prev => ({ ...prev, type: "error" }));
        setIsBilling(true);
    };

    const handleDialogClose = () => {
        setIsBilling(false);
    };

    const handlePaymentSuccess = () => {
        /* TODO: Handle payment success */
    };

    const handleViewPdf = () => {
        //handle-pdf-vieew
    };
    const handlePaymentMethod = () => {
        //get-payment-method
    };
    const addPaymentMethod = () => {
        //handle add pyament method
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
                        sx={{ mt: 2, mb: 4 }}
                        onChange={handleTabChange}
                    >
                        {BillingTabs.map(({ id, label, value }) => (
                            <Tab
                                key={id}
                                label={label}
                                sx={{ px: 3 }}
                                id={`billing-tab-${value}`}
                            />
                        ))}
                    </Tabs>

                    {tab === 0 && (
                        <Grid container item spacing={2}>
                            <Grid xs={12} md={5} item>
                                <CurrentBill
                                    amount={"$20.00"}
                                    title="jully bill"
                                    periodOf={"10/10/2021 - 11/10/2021"}
                                />
                            </Grid>
                            <Grid xs={12} md={7} item>
                                <RoundedCard>
                                    <PaymentCard
                                        title={"Payment settings"}
                                        handlePaymentMethod={
                                            handlePaymentMethod
                                        }
                                        paymentMethodData={PaymentCards}
                                        onAddPaymentMethod={addPaymentMethod}
                                    />
                                </RoundedCard>
                            </Grid>
                            <Grid xs={12} item>
                                <RoundedCard>
                                    <TableHeader
                                        title={"Billing breakdown"}
                                        showSecondaryButton={false}
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
                        <>
                            <RoundedCard>
                                <TableHeader
                                    title={"Billing history"}
                                    showSecondaryButton={false}
                                />

                                <SimpleDataTable
                                    isHistoryTab={true}
                                    rowSelection={true}
                                    handleViewPdf={handleViewPdf}
                                    selectedRows={selectedRows}
                                    columns={CurrentBillColumns}
                                    dataset={CurrentBillingData}
                                    setSelectedRows={setSelectedRows}
                                    totalRows={CurrentBillingData.length}
                                />
                            </RoundedCard>
                        </>
                    )}
                </Box>
            </LoadingWrapper>
            <BillingDialog
                isOpen={isBilling}
                handleCloseAction={handleDialogClose}
                handleSuccessAction={() => handlePaymentSuccess()}
            />
        </Box>
    );
};

export default Billing;
