import { globalUseStyles, ContainerJustifySpaceBtw } from "../../styles";
import {
    Typography,
    Grid,
    Switch,
    Button,
    TextField,
    Stack,
} from "@mui/material";
import * as Yup from "yup";
import { Formik } from "formik";
import { useState } from "react";
import { ESIM_FORM_SCHEMA } from "../../helpers/formValidators";
interface IPhysicalSimform {
    description: string;
    handleClose: Function;
    handlePhysicalSimInstallation: Function;
}
const physicalSimFormSchema = Yup.object(ESIM_FORM_SCHEMA);
const initialePhysicalSimFormValue = {
    name: "",
    email: "",
};

const PhysicalSimform = ({
    description,
    handleClose,
    handlePhysicalSimInstallation,
}: IPhysicalSimform) => {
    const [status, setStatus] = useState<boolean>(true);
    const gclasses = globalUseStyles();
    return (
        <Formik
            validationSchema={physicalSimFormSchema}
            initialValues={initialePhysicalSimFormValue}
            onSubmit={async values =>
                handlePhysicalSimInstallation({ ...values, status })
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
                                sx={{ mb: 1 }}
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
                                    <Typography variant="caption">
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
                                    checked={status}
                                    onChange={e => setStatus(e.target.checked)}
                                />
                            </ContainerJustifySpaceBtw>
                            <Stack
                                mt={1}
                                direction="row"
                                justifyContent="flex-end"
                            >
                                <Button
                                    sx={{ mr: 2, justifyItems: "center" }}
                                    onClick={() => handleClose()}
                                >
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
