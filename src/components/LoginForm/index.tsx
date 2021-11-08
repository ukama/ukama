import {
    Box,
    Button,
    Stack,
    Divider,
    TextField,
    Typography,
} from "@mui/material";
import "../../i18n/i18n";
import * as Yup from "yup";
import { Formik } from "formik";
import { useTranslation } from "react-i18next";
import { PasswordFieldWithIndicator } from "..";
import withAuthWrapperHOC from "../withAuthWrapperHOC";
import { LinkStyle, globalUseStyles } from "../../styles";

const loginSchema = Yup.object({
    email: Yup.string()
        .email("Please enter a valid email")
        .required("Please enter a valid email"),
    password: Yup.string().required("Password is required"),
});

const initialLoginValue = {
    email: "",
    password: "",
};

type LoginFormProps = {
    onSubmit: Function;
    onGoogleLogin: Function;
};

const LoginForm = ({ onSubmit, onGoogleLogin }: LoginFormProps) => {
    const classes = globalUseStyles();
    const { t } = useTranslation();

    return (
        <Box width="100%">
            <Formik
                validationSchema={loginSchema}
                initialValues={initialLoginValue}
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
                            <Typography variant="h5">
                                {t("LOGIN.FormTitle")}
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

                            <PasswordFieldWithIndicator
                                errors={errors}
                                touched={touched}
                                onBlur={handleBlur}
                                withIndicator={false}
                                value={values.password}
                                handleChange={handleChange}
                                label={t("CONSTANT.PasswordLabel")}
                                fieldStyle={classes.inputFieldStyle}
                            />

                            <LinkStyle
                                href="/forgot-password"
                                sx={{
                                    marginTop: "8px !important",
                                }}
                            >
                                {t("LOGIN.ForgotPasswordLabel")}
                            </LinkStyle>

                            <Button
                                size="medium"
                                type="submit"
                                color="primary"
                                variant="contained"
                            >
                                {t("LOGIN.ButtonLabel")}
                            </Button>

                            <Divider />

                            <Button
                                size="medium"
                                variant="outlined"
                                sx={{
                                    marginTop: "0px !important",
                                }}
                                onClick={() => onGoogleLogin()}
                            >
                                {t("LOGIN.ButtonWithGoogle")}
                            </Button>
                            <Typography
                                variant="caption"
                                sx={{ letterSpacing: "0.4px" }}
                            >
                                {t("LOGIN.FooterNoteText")}
                                <LinkStyle href="/signup">
                                    {t("CONSTANT.SignUpLinkLabel")}
                                </LinkStyle>
                            </Typography>
                        </Stack>
                    </form>
                )}
            </Formik>
        </Box>
    );
};

export default withAuthWrapperHOC(LoginForm);
