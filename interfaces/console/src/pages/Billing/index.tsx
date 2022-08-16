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
import {
    CurrentBillColumns,
    historyyBilling,
} from "../../constants/tableColumns";
import {
    useGetBillHistoryQuery,
    useGetCurrentBillQuery,
    useRetrivePaymentMethodsQuery,
} from "../../generated";
import { useState } from "react";
import { useRecoilValue } from "recoil";
import { RoundedCard } from "../../styles";
import { SelectItemType } from "../../types";
import { BillingTabs } from "../../constants";
import { isSkeltonLoading } from "../../recoil";
import { Box, Grid, Tabs, Typography, Tab, AlertColor } from "@mui/material";

const Billing = () => {
    const [isBilling, setIsBilling] = useState({
        isShow: false,
        isOnlypaymentFlow: false,
    });
    const [billingAlert] = useState({
        type: "info",
        btnText: "Enter now →",
        title: "Set up your payment information securely at any time.",
    });
    const [tab, setTab] = useState<number>(0);
    const _isSkeltonLoading = useRecoilValue(isSkeltonLoading);
    const [selectedRows, setSelectedRows] = useState<number[]>([]);
    const [cardsList, setCardsList] = useState<SelectItemType[]>([
        { id: "1", value: "no_payment_method_Set", label: "None set up." },
    ]);
    const { data: billingHistoryRes, loading: billingHistoryLoading } =
        useGetBillHistoryQuery();

    const { data: currentBill, loading: currenBillLoading } =
        useGetCurrentBillQuery();

    const { refetch: refetchPM } = useRetrivePaymentMethodsQuery({
        onCompleted: res => {
            if (res) {
                const list: SelectItemType[] = [];
                for (const element of res.retrivePaymentMethods) {
                    list.push({
                        id: element.id,
                        value: element.id,
                        label: `${element.brand} - ending in ${element.last4}`,
                    });
                }
                setCardsList(prev => [...list, ...prev]);
            }
        },
    });

    const handleTabChange = (_: any, value: any) => setTab(value);

    const handleAlertAction = () => {
        setIsBilling({ isShow: true, isOnlypaymentFlow: false });
    };

    const handleDialogClose = () => {
        setIsBilling({ isShow: false, isOnlypaymentFlow: false });
    };

    const handlePaymentSuccess = () => {
        refetchPM();
    };

    const handleViewPdf = () => {
        //handle-pdf-vieew
    };

    const addPaymentMethod = () => {
        setIsBilling({ isShow: true, isOnlypaymentFlow: true });
    };

    const totalCurrentBill: number | undefined =
        currentBill?.getCurrentBill?.bill.reduce(
            (totalCurrentBill, currentItem) =>
                (totalCurrentBill = totalCurrentBill + currentItem.subtotal),
            0
        );
    return (
        <Box>
            <BillingAlerts
                title={billingAlert.title}
                btnText={billingAlert.btnText}
                onActionClick={handleAlertAction}
                type={billingAlert.type as AlertColor}
            />
            <LoadingWrapper isLoading={_isSkeltonLoading} height={"300px"}>
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
                                    amount={`$ ${currentBill?.getCurrentBill?.total}`}
                                    billMonth={`${currentBill?.getCurrentBill?.billMonth}`}
                                    dueDate={`${currentBill?.getCurrentBill?.dueDate}`}
                                    loading={currenBillLoading}
                                />
                            </Grid>
                            <Grid xs={12} md={7} item>
                                <RoundedCard>
                                    <PaymentCard
                                        title={"Payment settings"}
                                        paymentMethodData={cardsList}
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
                                        dataset={
                                            currentBill?.getCurrentBill?.bill
                                        }
                                    />

                                    <div
                                        style={{
                                            width: "100%",
                                            display: "flex",
                                            margin: "18px 0px",
                                            justifyContent: "flex-end",
                                            position: "relative",
                                            right: 48,
                                        }}
                                    >
                                        <Typography
                                            variant={"h6"}
                                            sx={{
                                                width: "20%",
                                            }}
                                        >
                                            {`$ ${totalCurrentBill}`}
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
                                <LoadingWrapper
                                    isLoading={billingHistoryLoading}
                                    height={200}
                                >
                                    <SimpleDataTable
                                        isHistoryTab={true}
                                        rowSelection={true}
                                        handleViewPdf={handleViewPdf}
                                        selectedRows={selectedRows}
                                        columns={historyyBilling}
                                        dataset={
                                            billingHistoryRes?.getBillHistory
                                        }
                                        setSelectedRows={setSelectedRows}
                                        totalRows={
                                            billingHistoryRes?.getBillHistory
                                                .length
                                        }
                                    />
                                </LoadingWrapper>
                            </RoundedCard>
                        </>
                    )}
                </Box>
            </LoadingWrapper>
            {isBilling.isShow && (
                <BillingDialog
                    isOpen={isBilling.isShow}
                    handleCloseAction={handleDialogClose}
                    initPaymentFlow={isBilling.isOnlypaymentFlow}
                    handleSuccessAction={handlePaymentSuccess}
                />
            )}
        </Box>
    );
};

export default Billing;
