import { globalUseStyles, ContainerJustifySpaceBtw } from "../../styles";
import {
    Typography,
    Grid,
    Switch,
    Button,
    TextField,
    Stack,
} from "@mui/material";
import { Formik } from "formik";
import { useState } from "react";
import * as Yup from "yup";
import { colors } from "../../theme";
import { PHYSICAL_SIM_FORM_SCHEMA } from "../../helpers/formValidators";
interface IPhysicalSimform {
    description: string;
    handlePhysicalSimInstallation: Function;
}
const physicalSimFormSchema = Yup.object(PHYSICAL_SIM_FORM_SCHEMA);
const initialePhysicalSimFormValue = {
    iccid: "",
    securityCode: "",
};

const PhysicalSimform = ({
    description,
    handlePhysicalSimInstallation,
}: IPhysicalSimform) => {
    const [roaming, setRoaming] = useState<boolean>(false);
    const gclasses = globalUseStyles();
    return (
        <Formik
            validationSchema={physicalSimFormSchema}
            initialValues={initialePhysicalSimFormValue}
            onSubmit={async values =>
                handlePhysicalSimInstallation({ ...values, roaming })
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
                                id="iccid"
                                name="iccid"
                                label={"ICCID"}
                                onBlur={handleBlur}
                                onChange={handleChange}
                                value={values.iccid}
                                sx={{ mb: 1 }}
                                InputLabelProps={{ shrink: true }}
                                InputProps={{
                                    classes: {
                                        input: gclasses.inputFieldStyle,
                                    },
                                }}
                                helperText={touched.iccid && errors.iccid}
                                error={touched.iccid && Boolean(errors.iccid)}
                            />
                        </Grid>
                        <Grid item xs={12}>
                            <TextField
                                fullWidth
                                id="securityCode"
                                name="securityCode"
                                label={"SECURITY CODE"}
                                onBlur={handleBlur}
                                onChange={handleChange}
                                value={values.securityCode}
                                sx={{ mb: 1 }}
                                InputLabelProps={{ shrink: true }}
                                InputProps={{
                                    classes: {
                                        input: gclasses.inputFieldStyle,
                                    },
                                }}
                                helperText={
                                    touched.securityCode && errors.securityCode
                                }
                                error={
                                    touched.securityCode &&
                                    Boolean(errors.securityCode)
                                }
                            />
                        </Grid>
                        <Grid item xs={12}>
                            <ContainerJustifySpaceBtw
                                sx={{ alignItems: "end" }}
                            >
                                <Stack display="flex" alignItems="flex-start">
                                    <Typography
                                        variant="caption"
                                        color={colors.black54}
                                    >
                                        ROAMING
                                    </Typography>
                                    <Typography variant="body1">
                                        Roaming allows user to do xyz. Insert
                                        billing information.
                                    </Typography>
                                </Stack>
                                <Switch
                                    size="small"
                                    value="active"
                                    checked={roaming}
                                    onChange={e => setRoaming(e.target.checked)}
                                />
                            </ContainerJustifySpaceBtw>
                            <Stack direction="row" justifyContent="flex-end">
                                <Button sx={{ mr: 2, justifyItems: "center" }}>
                                    Cancel
                                </Button>
                                <Button variant="contained" type="submit">
                                    ADD USER
                                </Button>
                            </Stack>
                        </Grid>
                    </Grid>
                </form>
            )}
        </Formik>
    );
};

export default PhysicalSimform;
