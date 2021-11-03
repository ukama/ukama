import * as Yup from "yup";
import { Formik } from "formik";
import withAuthWrapperHOC from "../withAuthWrapperHOC";
import { ContainerJustifySpaceBtw, globalUseStyles } from "../../styles";
import { Box, Stack, Button, TextField, Typography } from "@mui/material";
import { useTranslation } from "react-i18next";
import "../../i18n/i18n";
const forgotPasswordSchema = Yup.object({
    email: Yup.string()
        .email("Please enter a valid email")
        .required("Please enter a valid email"),
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
    const { t } = useTranslation();
    return (
        <Box width="100%">
            <Formik
                initialValues={forgotPasswordValue}
                validationSchema={forgotPasswordSchema}
                onSubmit={async values => onSubmit(values)}
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
                        <Stack spacing={"18px"}>
                            <Typography variant="h3">
                                {t("FORGOT_PASSWORD.FormTitle")}
                            </Typography>

                            <TextField
                                fullWidth
                                id="email"
                                name="email"
                                onBlur={handleBlur}
                                value={values.email}
                                onChange={handleChange}
                                label={t("CONSTANT.EmailLabel")}
                                InputLabelProps={{ shrink: true }}
                                InputProps={{
                                    classes: { input: classes.inputFieldStyle },
                                }}
                                helperText={touched.email && errors.email}
                                error={touched.email && Boolean(errors.email)}
                            />
                            <ContainerJustifySpaceBtw>
                                <Button
                                    size="medium"
                                    variant="text"
                                    onClick={() => onBack()}
                                >
                                    {t("CONSTANT.BackButtonLable")}
                                </Button>

                                <Button
                                    size="medium"
                                    type="submit"
                                    color="primary"
                                    variant="contained"
                                >
                                    {t("CONSTANT.RecoveryEmail")}
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
