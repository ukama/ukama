import {
    Grid,
    Stack,
    Switch,
    TextField,
    Button,
    Typography,
} from "@mui/material";
import { colors } from "../../theme";
import { ContainerJustifySpaceBtw, globalUseStyles } from "../../styles";
import { Formik } from "formik";
import { useState } from "react";
import * as Yup from "yup";
import { ESIM_FORM_SCHEMA } from "../../helpers/formValidators";
interface IUserform {
    description: string;
    handleEsimInstallation: Function;
}
const eSimFormSchema = Yup.object(ESIM_FORM_SCHEMA);
const initialeEsimFormValue = {
    name: "",
    email: "",
};

const Userform = ({ handleEsimInstallation, description }: IUserform) => {
    const gclasses = globalUseStyles();
    const [status, setStatus] = useState<boolean>(false);

    return (
        <Formik
            validationSchema={eSimFormSchema}
            initialValues={initialeEsimFormValue}
            onSubmit={async values =>
                handleEsimInstallation({ ...values, status })
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

export default Userform;
