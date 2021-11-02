import React, { useState } from "react";
import { Box, Grid, TextField } from "@mui/material";
import { RoundedCard, globalUseStyles } from "../../styles";
import * as Yup from "yup";
import { Formik } from "formik";
import {
    NodeCard,
    StatusCard,
    NetworkStatus,
    ContainerHeader,
    StatsCard,
    AlertCard,
    FormDialog,
} from "../../components";
import { useTranslation } from "react-i18next";
import "../../i18n/i18n";
import { DashboardStatusCard } from "../../constants/stubData";
import {
    STATS_OPTIONS,
    STATS_PERIOD,
    NETWORKS,
    ALERT_INFORMATION,
} from "../../constants";
const Home = () => {
    const [network, setNetwork] = useState("public");
    const [userStatusFilter, setUserStatusFilter] = useState("total");
    const [dataStatusFilter, setDataStatusFilter] = useState("total");
    const [billingStatusFilter, setBillingStatusFilter] = useState("july");
    const [statOptionValue, setstatOptionValue] = React.useState(3);
    const [open, setOpen] = React.useState(false);
    const { t } = useTranslation();
    const classes = globalUseStyles();
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
    const addNodeSchema = Yup.object({
        nodeName: Yup.string().required("Node Name is required"),
        nodeSerialNumber: Yup.string().required("Serial number is required"),
    });

    const initialAddNodeValue = {
        nodeName: "",
        nodeSerialNumber: "",
    };

    const handleSubmit = (values: any) => {
        //console.log(JSON.stringify(values, null, 2));
    };

    const openAddNodeDialog = () => {
        setOpen(true);
    };

    const closeAddNodeDialog = () => {
        setOpen(false);
    };

    return (
        <>
            <Box sx={{ flexGrow: 1 }}>
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
                <Box mt={2} mb={2}>
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
                            <AlertCard alertCardItems={ALERT_INFORMATION} />
                        </Grid>
                    </Grid>
                </Box>

                <Grid container spacing={2}>
                    <Grid xs={12} md={8} item>
                        <RoundedCard>
                            <Formik
                                validationSchema={addNodeSchema}
                                initialValues={initialAddNodeValue}
                                onSubmit={async values => handleSubmit(values)}
                            >
                                {({
                                    values,
                                    errors,
                                    touched,
                                    handleChange,
                                    handleSubmit,
                                }) => (
                                    <form onSubmit={handleSubmit}>
                                        <FormDialog
                                            onClose={closeAddNodeDialog}
                                            open={open}
                                            dialogContent={t(
                                                "HOME.FormDialogContent"
                                            )}
                                            dialogTitle={t(
                                                "HOME.FormDialogTitle"
                                            )}
                                            submitButtonLabel={t(
                                                "CONSTANT.NextButtonLabel"
                                            )}
                                        >
                                            <>
                                                <Box pt={3}>
                                                    <Grid
                                                        container
                                                        spacing={2}
                                                        xs={12}
                                                    >
                                                        <Grid item xs={6}>
                                                            <TextField
                                                                fullWidth
                                                                id="nodeName"
                                                                name="nodeName"
                                                                label={t(
                                                                    "NODE.NodeNameLabel"
                                                                )}
                                                                value={
                                                                    values.nodeName
                                                                }
                                                                onChange={
                                                                    handleChange
                                                                }
                                                                InputLabelProps={{
                                                                    shrink: true,
                                                                }}
                                                                InputProps={{
                                                                    classes: {
                                                                        input: classes.inputFieldStyle,
                                                                    },
                                                                }}
                                                                helperText={
                                                                    touched.nodeName &&
                                                                    errors.nodeName
                                                                }
                                                                error={
                                                                    touched.nodeName &&
                                                                    Boolean(
                                                                        errors.nodeName
                                                                    )
                                                                }
                                                            />
                                                        </Grid>
                                                        <Grid item xs={6}>
                                                            <TextField
                                                                fullWidth
                                                                id="nodeSerialNumber"
                                                                name="nodeSerialNumber"
                                                                label={t(
                                                                    "NODE.NodeSerialNumberLabel"
                                                                )}
                                                                value={
                                                                    values.nodeSerialNumber
                                                                }
                                                                onChange={
                                                                    handleChange
                                                                }
                                                                InputLabelProps={{
                                                                    shrink: true,
                                                                }}
                                                                InputProps={{
                                                                    classes: {
                                                                        input: classes.inputFieldStyle,
                                                                    },
                                                                }}
                                                                helperText={
                                                                    touched.nodeSerialNumber &&
                                                                    errors.nodeSerialNumber
                                                                }
                                                                error={
                                                                    touched.nodeSerialNumber &&
                                                                    Boolean(
                                                                        errors.nodeSerialNumber
                                                                    )
                                                                }
                                                            />
                                                        </Grid>
                                                    </Grid>
                                                </Box>
                                            </>
                                        </FormDialog>
                                    </form>
                                )}
                            </Formik>
                            <ContainerHeader
                                stats="1/8"
                                title="My Nodes"
                                buttonTitle="Add Node"
                                handleButtonAction={() => {
                                    openAddNodeDialog();
                                }}
                            />
                            <NodeCard isConfigure={true} />
                        </RoundedCard>
                    </Grid>
                    <Grid xs={12} md={4} item>
                        <RoundedCard sx={{ height: "100%" }}></RoundedCard>
                    </Grid>
                </Grid>
            </Box>
        </>
    );
};
export default Home;
