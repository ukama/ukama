import { BillingTabs, BILLING_TYPE } from "../../constants";
import { CREDIT_CARD } from "../../constants/stubData";
import TabLayoutHeader from "../../components/TabLayout";
import MenuItem from "@mui/material/MenuItem";
import {
    Box,
    Card,
    CardContent,
    Grid,
    Typography,
    Button,
    TextField,
} from "@mui/material";
import { useState } from "react";
import { RoundedCard } from "../../styles";
import { AmountBalanceImg } from "../../assets/svg";
import colors from "../../theme/colors";
import { BillingType, PaymentMethodType } from "../../types";
const Billing = () => {
    const [tab, setTab] = useState("1");
    const handleTabChange = (value: string) => setTab(value);
    const [billingTime, setBillingTime] = useState("AUTO");
    const [paymentMethod, setPaymentMethod] = useState(1);
    const handlePaymentMethod = (event: any) =>
        setPaymentMethod(event.target.value);
    const handleBillingTime = (event: any) =>
        setBillingTime(event.target.value);

    return (
        <Box sx={{ mt: "24px" }}>
            <TabLayoutHeader
                tab={tab}
                tabs={BillingTabs}
                onTabChange={handleTabChange}
            />
            {tab === "1" ? (
                <Grid container spacing={3} sx={{ pt: 2 }}>
                    <Grid item xs={5}>
                        <RoundedCard>
                            <Card sx={{ display: "flex" }} elevation={0}>
                                <Box
                                    sx={{
                                        display: "flex",
                                        flexDirection: "column",
                                        width: "100%",
                                    }}
                                >
                                    <CardContent sx={{ flex: "1 0" }}>
                                        <Typography
                                            component="div"
                                            variant="h6"
                                        >
                                            Amount due
                                        </Typography>
                                        <Typography
                                            variant="body2"
                                            color="text.secondary"
                                            component="div"
                                        >
                                            10/10/2021 - 11/10/2021
                                        </Typography>
                                    </CardContent>
                                    <Box
                                        sx={{
                                            alignItems: "center",

                                            pl: 1,
                                        }}
                                    >
                                        <Typography variant="h3">
                                            $20.00
                                        </Typography>
                                    </Box>
                                </Box>
                                <Box sx={{ height: "100%" }}>
                                    <AmountBalanceImg />
                                </Box>
                            </Card>
                            <Box sx={{ pl: 1 }}>
                                <Typography variant="caption">
                                    *Automatically charged to card on November
                                    10, 2021
                                </Typography>
                            </Box>
                        </RoundedCard>
                    </Grid>
                    <Grid item xs={7}>
                        <RoundedCard>
                            <Grid container spacing={1}>
                                <Grid item container xs={6}>
                                    <Typography variant="h6">
                                        Payment method
                                    </Typography>
                                </Grid>
                                <Grid
                                    item
                                    container
                                    justifyContent="flex-end"
                                    xs={6}
                                >
                                    <Button variant="text" color="primary">
                                        Edit
                                    </Button>
                                </Grid>
                                <Grid item container xs={6}>
                                    <TextField
                                        InputLabelProps={{
                                            style: {
                                                color: colors.lightGrey,
                                            },
                                        }}
                                        id="billing-time"
                                        select
                                        label="BILLING TIME"
                                        value={billingTime}
                                        onChange={handleBillingTime}
                                        sx={{ width: "100%" }}
                                    >
                                        {BILLING_TYPE.map(
                                            ({ value, label }: BillingType) => (
                                                <MenuItem
                                                    key={value}
                                                    value={value}
                                                >
                                                    <Typography variant="body1">
                                                        {label}
                                                    </Typography>
                                                </MenuItem>
                                            )
                                        )}
                                    </TextField>
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
                                    </TextField>
                                </Grid>
                            </Grid>
                        </RoundedCard>
                    </Grid>
                </Grid>
            ) : (
                <h4>Billing History</h4>
            )}
        </Box>
    );
};

export default Billing;
