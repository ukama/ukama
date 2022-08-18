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
import {
    AlertColor,
    Tabs,
    Tab,
    Grid,
    Typography,
    Stack,
    Box,
} from "@mui/material";
import { useState } from "react";
import colors from "../../theme/colors";
import { useRecoilValue } from "recoil";
import { RoundedCard } from "../../styles";
import { NoBillYet } from "../../assets/svg";
import { SelectItemType } from "../../types";
import { BillingTabs } from "../../constants";
import { isSkeltonLoading, isDarkmode } from "../../recoil";

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
    const _isDarkmode = useRecoilValue(isDarkmode);

    const _isSkeltonLoading = useRecoilValue(isSkeltonLoading);
    const [selectedRows, setSelectedRows] = useState<number[]>([]);
    const [cardsList, setCardsList] = useState<SelectItemType[]>([
        { id: "1", value: "no_payment_method_Set", label: "None set up." },
    ]);
    const { data: billingHistoryRes, loading: billingHistoryLoading } =
        useGetBillHistoryQuery();
    const isSkeltonLoad = useRecoilValue(isSkeltonLoading);

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
                                    amount={`$ ${
                                        currentBill?.getCurrentBill?.total ||
                                        "0.00"
                                    }`}
                                    billMonth={`${
                                        currentBill
                                            ? currentBill?.getCurrentBill
                                                  ?.billMonth
                                            : ""
                                    }`}
                                    dueDate={`${
                                        currentBill
                                            ? currentBill?.getCurrentBill
                                                  ?.dueDate
                                            : ""
                                    }`}
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
                                    {totalCurrentBill !== undefined || null ? (
                                        <>
                                            <SimpleDataTable
                                                columns={CurrentBillColumns}
                                                dataset={
                                                    currentBill?.getCurrentBill
                                                        ?.bill
                                                }
                                            />

                                            <div
                                                style={{
                                                    width: "100%",
                                                    display: "flex",
                                                    justifyContent: "flex-end",
                                                }}
                                            >
                                                <Typography
                                                    variant={"h6"}
                                                    sx={{
                                                        width: "20%",
                                                        position: "relative",
                                                        right: 50,
                                                    }}
                                                >
                                                    {`$ ${totalCurrentBill}`}
                                                </Typography>
                                            </div>
                                        </>
                                    ) : (
                                        <Stack
                                            direction="column"
                                            spacing={2}
                                            justifyItems={"center"}
                                            alignItems={"center"}
                                        >
                                            <NoBillYet />
                                            <Typography variant="body1">
                                                Nothing in your bill yet!
                                            </Typography>
                                        </Stack>
                                    )}
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
                                    isLoading={
                                        isSkeltonLoad || billingHistoryLoading
                                    }
                                    height={200}
                                >
                                    {billingHistoryRes !== undefined || null ? (
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
                                                billingHistoryRes
                                                    ?.getBillHistory.length
                                            }
                                        />
                                    ) : (
                                        <Box
                                            display="flex"
                                            justifyContent="center"
                                            alignItems="center"
                                            minHeight="60vh"
                                        >
                                            <Stack
                                                direction="column"
                                                spacing={2}
                                            >
                                                <NoBillYet
                                                    color={
                                                        _isDarkmode
                                                            ? colors.white38
                                                            : colors.silver
                                                    }
                                                    color2={
                                                        _isDarkmode
                                                            ? colors.nightGrey12
                                                            : colors.white
                                                    }
                                                />
                                                <Typography variant="body1">
                                                    No bill History yet!
                                                </Typography>
                                            </Stack>
                                        </Box>
                                    )}
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
