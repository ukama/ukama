import * as Yup from "yup";
import { Formik } from "formik";
import withAuthWrapperHOC from "../withAuthWrapperHOC";
import { ContainerJustifySpaceBtw, globalUseStyles } from "../../styles";
import { Box, Stack, Button, TextField, Typography } from "@mui/material";

const forgotPasswordSchema = Yup.object({
    email: Yup.string()
        .email("Enter a valid email")
        .required("Email is required"),
});

const forgotPasswordValue = {
    email: "",
};

type ForgotPasswordFormProps = {
    onSubmit: Function;
    onBack: Function;
};

const ForgotPasswordForm = ({ onSubmit, onBack }: ForgotPasswordFormProps) => {
    const classes = globalUseStyles();

    return (
        <Box width="100%">
            <Formik
                initialValues={forgotPasswordValue}
                validationSchema={forgotPasswordSchema}
                onSubmit={async values => onSubmit(values)}
            >
                {({ values, errors, touched, handleChange, handleSubmit }) => (
                    <form onSubmit={handleSubmit}>
                        <Stack spacing={"18px"}>
                            <Typography variant="h3">
                                Recover Password
                            </Typography>

                            <TextField
                                fullWidth
                                id="email"
                                name="email"
                                label="Email"
                                value={values.email}
                                onChange={handleChange}
                                InputLabelProps={{ shrink: true }}
                                InputProps={{
                                    classes: { input: classes.inputFieldStyle },
                                }}
                                helperText={touched.email && errors.email}
                                error={touched.email && Boolean(errors.email)}
                            />
                            <ContainerJustifySpaceBtw>
                                <Button
                                    size="large"
                                    variant="text"
                                    sx={{ fontWeight: 600 }}
                                    onClick={() => onBack()}
                                >
                                    BACK
                                </Button>

                                <Button
                                    size="large"
                                    type="submit"
                                    color="primary"
                                    variant="contained"
                                    sx={{ fontWeight: 600 }}
                                >
                                    EMAIL RECOVERY LINK
                                </Button>
                            </ContainerJustifySpaceBtw>
                        </Stack>
                    </form>
                )}
            </Formik>
        </Box>
    );
};

export default withAuthWrapperHOC(ForgotPasswordForm);
