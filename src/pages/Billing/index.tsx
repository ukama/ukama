import {
    BillingTabs,
    BILLING_TABLE_HEADER_OPTIONS,
    BILLING_HOSTORY_TABLE_HEADER_OPTIONS,
} from "../../constants";
import { CREDIT_CARD, CurrentBilling } from "../../constants/stubData";
import MenuItem from "@mui/material/MenuItem";
import { useTranslation } from "react-i18next";
import "../../i18n/i18n";
import { BillingDataTable } from "../../components";
import TabLayoutHeader from "../../components/TabLayout";

import {
    Box,
    TableCell,
    Paper,
    Grid,
    Typography,
    Button,
    TextField,
    FormGroup,
    TableBody,
    Stack,
    TableRow,
    FormControlLabel,
    Switch,
} from "@mui/material";
import { useState } from "react";
import AddIcon from "@mui/icons-material/Add";
import { RoundedCard } from "../../styles";
import { AmountBalanceImg } from "../../assets/svg";
import colors from "../../theme/colors";
import { PaymentMethodType, currentBillType } from "../../types";
const TabCellStyle = {
    borderBottom: "none",
};
const Billing = () => {
    const { t } = useTranslation();
    const [tab, setTab] = useState("1");
    const handleTabChange = (value: string) => setTab(value);
    const [paymentMethod, setPaymentMethod] = useState(1);
    const [autoCharge, setAutoCharge] = useState(false);

    const handlePaymentMethod = (event: any) =>
        setPaymentMethod(event.target.value);
    const handleAutoCharge = (event: any) =>
        setAutoCharge(event.target.checked);
    return (
        <Box sx={{ mt: "24px" }}>
            <TabLayoutHeader
                tab={tab}
                tabs={BillingTabs}
                onTabChange={handleTabChange}
            />
            {tab === "1" ? (
                <Grid container spacing={3} sx={{ pt: 2 }}>
                    <Grid item xs={12} md={5}>
                        <RoundedCard>
                            <Grid container spacing={1}>
                                <Grid item xs={6}>
                                    <Stack direction="column">
                                        <Typography variant="h5">
                                            {t("BILLING.AmountDue")}
                                        </Typography>

                                        <Typography variant="body2">
                                            10/10/2021 - 11/10/2021
                                        </Typography>
                                    </Stack>
                                </Grid>
                                <Grid
                                    item
                                    xs={6}
                                    container
                                    justifyContent="flex-end"
                                >
                                    <Paper
                                        elevation={0}
                                        style={{
                                            borderRadius: "30px",
                                            width: "80%",
                                            backgroundColor: colors.lightRed,
                                            textAlign: "center",
                                        }}
                                    >
                                        <Box pt={2}>
                                            <Typography variant="body2">
                                                Due 10/10/2021
                                            </Typography>
                                        </Box>
                                    </Paper>
                                </Grid>
                                <Grid item xs={3}>
                                    <Box pt={5}>
                                        <Typography variant="h3">
                                            $20.00
                                        </Typography>
                                    </Box>
                                </Grid>
                                <Grid
                                    item
                                    xs={9}
                                    container
                                    justifyContent="flex-end"
                                >
                                    <AmountBalanceImg />
                                </Grid>
                            </Grid>
                        </RoundedCard>
                    </Grid>
                    <Grid item xs={12} md={7}>
                        <RoundedCard>
                            <Grid container spacing={1}>
                                <Grid item container xs={6}>
                                    <Typography variant="h6">
                                        {t("BILLING.PaymentInformation")}
                                    </Typography>
                                </Grid>
                                <Grid
                                    item
                                    container
                                    justifyContent="flex-end"
                                    xs={6}
                                >
                                    <FormGroup>
                                        <FormControlLabel
                                            control={
                                                <Switch
                                                    checked={autoCharge}
                                                    onChange={handleAutoCharge}
                                                />
                                            }
                                            label={t(
                                                "BILLING.AutoChargeSwitchButtonLabel"
                                            )}
                                        />
                                    </FormGroup>
                                </Grid>
                                <Grid item container xs={6}>
                                    <Button
                                        variant="contained"
                                        size="large"
                                        sx={{ width: "80%" }}
                                    >
                                        {t("BILLING.MakeAPayment")}
                                    </Button>
                                </Grid>
                                <Grid item container xs={6}>
                                    <TextField
                                        InputLabelProps={{
                                            style: {
                                                color: colors.lightGrey,
                                            },
                                        }}
                                        id="payment-method"
                                        select
                                        label="PAYMENT-METHOD"
                                        value={paymentMethod}
                                        onChange={handlePaymentMethod}
                                        sx={{ width: "100%" }}
                                    >
                                        {CREDIT_CARD.map(
                                            ({
                                                id,
                                                card_experintionDetails,
                                            }: PaymentMethodType) => (
                                                <MenuItem key={id} value={id}>
                                                    <Typography variant="body1">
                                                        {
                                                            card_experintionDetails
                                                        }
                                                    </Typography>
                                                </MenuItem>
                                            )
                                        )}
                                        <MenuItem>
                                            <Button
                                                variant="text"
                                                startIcon={<AddIcon />}
                                            >
                                                {t("BILLING.FormTitle")}
                                            </Button>
                                        </MenuItem>
                                    </TextField>
                                </Grid>
                                <Grid item mt={4}>
                                    <Typography variant="caption">
                                        *Bill due on November 10, 2021
                                    </Typography>
                                </Grid>
                            </Grid>
                        </RoundedCard>
                    </Grid>
                    <Grid container spacing={1}>
                        <Grid item xs={12} ml={3} mt={2}>
                            <BillingDataTable
                                tableTitle="Bill breakDown"
                                headerOptions={BILLING_TABLE_HEADER_OPTIONS}
                                isBillingHistory={false}
                            >
                                {CurrentBilling.map(
                                    ({
                                        id,
                                        name,
                                        rate,
                                        dataUsage,
                                        subTotal,
                                    }: currentBillType) => (
                                        <>
                                            <TableRow key={id}>
                                                <TableCell align="left">
                                                    {name}
                                                </TableCell>
                                                <TableCell align="left">
                                                    {rate}
                                                </TableCell>
                                                <TableCell align="left">
                                                    {dataUsage}
                                                </TableCell>
                                                <TableCell align="left">
                                                    {subTotal}
                                                </TableCell>
                                            </TableRow>
                                        </>
                                    )
                                )}
                                <TableBody>
                                    <TableRow>
                                        <TableCell
                                            style={{ ...TabCellStyle }}
                                        />
                                        <TableCell
                                            style={{ ...TabCellStyle }}
                                        />
                                        <TableCell
                                            style={{ ...TabCellStyle }}
                                        />
                                        <TableCell
                                            align="left"
                                            style={{ ...TabCellStyle }}
                                        >
                                            <Typography variant="h6">
                                                {" "}
                                                $20.0
                                            </Typography>
                                        </TableCell>
                                    </TableRow>
                                </TableBody>
                            </BillingDataTable>
                        </Grid>
                    </Grid>
                </Grid>
            ) : (
                <Grid container spacing={1}>
                    <Grid item xs={12} ml={3} mt={2}>
                        <BillingDataTable
                            tableTitle="Billing History"
                            headerOptions={BILLING_HOSTORY_TABLE_HEADER_OPTIONS}
                            isBillingHistory={true}
                        ></BillingDataTable>
                    </Grid>
                </Grid>
            )}
        </Box>
    );
};

export default Billing;
