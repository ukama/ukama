import {
    Grid,
    Stack,
    Radio,
    Switch,
    TextField,
    Paper,
    Button,
    Typography,
} from "@mui/material";
import * as Yup from "yup";
import { Formik } from "formik";
import React, { useState } from "react";
import { ESIM_FORM_SCHEMA } from "../../helpers/formValidators";
import { ContainerJustifySpaceBtw, globalUseStyles } from "../../styles";
import colors from "../../theme/colors";
interface IUserform {
    description: string;
    handleClose: Function;
    handleSimInstallation: Function;
    eSimLeft?: number;
    physicalSimLeft?: number;
    getSimType: any;
}
const eSimFormSchema = Yup.object(ESIM_FORM_SCHEMA);
const initialeEsimFormValue = {
    name: "",
    email: "",
};

const Userform = ({
    handleClose,
    description,
    eSimLeft,
    physicalSimLeft,
    handleSimInstallation,
    getSimType,
}: IUserform) => {
    const gclasses = globalUseStyles();
    const [status, setStatus] = useState<boolean>(true);
    const [selectedSimType, setSelectedSimType] = useState("eSim");

    const handleSimTypeChange = (
        event: React.ChangeEvent<HTMLInputElement>
    ) => {
        setSelectedSimType(event.target.value);
        getSimType(event.target.value);
    };

    return (
        <Formik
            validationSchema={eSimFormSchema}
            initialValues={initialeEsimFormValue}
            onSubmit={async values =>
                handleSimInstallation({ ...values, status, selectedSimType })
            }
        >
            {({
                values,
                errors,
                touched,
                handleChange,
                handleSubmit,
                handleBlur,
            }) => (
                <form onSubmit={handleSubmit}>
                    <Grid container spacing={2}>
                        <Grid item xs={12} mb={1}>
                            <Typography variant="body1">
                                {description}
                            </Typography>
                        </Grid>

                        <Grid item xs={12}>
                            <TextField
                                fullWidth
                                id="name"
                                name="name"
                                label="NAME"
                                onBlur={handleBlur}
                                onChange={handleChange}
                                value={values.name}
                                sx={{ mb: 1 / 2 }}
                                InputLabelProps={{ shrink: true }}
                                InputProps={{
                                    classes: {
                                        input: gclasses.inputFieldStyle,
                                    },
                                }}
                                helperText={touched.name && errors.name}
                                error={touched.name && Boolean(errors.name)}
                            />
                        </Grid>
                        <Grid
                            item
                            xs={12}
                            container
                            spacing={2}
                            sx={{ mb: 1 / 2 }}
                        >
                            <Grid item xs={6}>
                                <Paper variant="outlined" sx={{}}>
                                    <Stack
                                        direction="row"
                                        spacing={1}
                                        alignItems="center"
                                    >
                                        <Radio
                                            checked={selectedSimType === "eSim"}
                                            onChange={handleSimTypeChange}
                                            value="eSim"
                                            name="eSim"
                                            inputProps={{
                                                "aria-label": "eSim",
                                            }}
                                        />
                                        <Typography variant="body1">
                                            {`eSIM (${eSimLeft || 0} left) `}
                                        </Typography>
                                    </Stack>
                                </Paper>
                            </Grid>
                            <Grid item xs={6}>
                                <Paper variant="outlined">
                                    <Stack
                                        direction="row"
                                        spacing={1}
                                        alignItems="center"
                                    >
                                        <Radio
                                            checked={
                                                selectedSimType ===
                                                "physicalSim"
                                            }
                                            onChange={handleSimTypeChange}
                                            value="physicalSim"
                                            name="physicalSim"
                                            inputProps={{
                                                "aria-label": "PhysicalSim",
                                            }}
                                        />
                                        <Typography variant="body1">
                                            {`   Physical SIM (${
                                                physicalSimLeft || 0
                                            } left) `}
                                        </Typography>
                                    </Stack>
                                </Paper>
                            </Grid>
                        </Grid>
                        <Grid item xs={12}>
                            <TextField
                                fullWidth
                                id="email"
                                name="email"
                                label="EMAIL"
                                onBlur={handleBlur}
                                onChange={handleChange}
                                value={values.email}
                                sx={{ mb: 1 }}
                                InputLabelProps={{ shrink: true }}
                                InputProps={{
                                    classes: {
                                        input: gclasses.inputFieldStyle,
                                    },
                                }}
                                helperText={touched.email && errors.email}
                                error={touched.email && Boolean(errors.email)}
                            />
                        </Grid>

                        <Grid item xs={12}>
                            <ContainerJustifySpaceBtw
                                sx={{ alignItems: "end" }}
                            >
                                <Stack display="flex" alignItems="flex-start">
                                    <Typography
                                        variant="caption"
                                        sx={{ color: colors.black38 }}
                                    >
                                        ALLOW ROAMING
                                    </Typography>
                                    <Typography variant="body1">
                                        Roaming allows user to do xyz. Insert
                                        billing information.
                                    </Typography>
                                </Stack>
                                <Switch
                                    size="small"
                                    value="active"
                                    checked={status}
                                    onChange={e => setStatus(e.target.checked)}
                                />
                            </ContainerJustifySpaceBtw>
                            <Stack
                                direction="row"
                                justifyContent="flex-end"
                                mt={1}
                            >
                                <Button
                                    sx={{ mr: 2, justifyItems: "center" }}
                                    onClick={() => handleClose()}
                                >
                                    Cancel
                                </Button>
                                <Button variant="contained" type="submit">
                                    NEXT
                                </Button>
                            </Stack>
                        </Grid>
                    </Grid>
                </form>
            )}
        </Formik>
    );
};

export default Userform;
