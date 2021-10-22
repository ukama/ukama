import {
    Box,
    Button,
    Stack,
    Divider,
    TextField,
    IconButton,
    Typography,
    InputAdornment,
} from "@mui/material";
import * as Yup from "yup";
import { Formik } from "formik";
import { useState } from "react";
import withAuthWrapperHOC from "../withAuthWrapperHOC";
import { LinkStyle, globalUseStyles } from "../../styles";
import { Visibility, VisibilityOff } from "@mui/icons-material";
import { useTranslation } from "react-i18next";
import "../../i18n/i18n";
const loginSchema = Yup.object({
    email: Yup.string()
        .email("Enter a valid email")
        .required("Email is required"),
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
    const [togglePassword, setTogglePassword] = useState(false);
    const handleTogglePassword = () => {
        setTogglePassword(prev => !prev);
    };

    return (
        <Box width="100%">
            <Formik
                validationSchema={loginSchema}
                initialValues={initialLoginValue}
                onSubmit={async values => onSubmit(values)}
            >
                {({ values, errors, touched, handleChange, handleSubmit }) => (
                    <form onSubmit={handleSubmit}>
                        <Stack spacing={"18px"}>
                            <Typography variant="h3">
                                {t("LOGIN.FormTitle")}
                            </Typography>

                            <TextField
                                fullWidth
                                id="email"
                                name="email"
                                label={t("CONSTANT.EmailLabel")}
                                value={values.email}
                                onChange={handleChange}
                                InputLabelProps={{ shrink: true }}
                                InputProps={{
                                    classes: { input: classes.inputFieldStyle },
                                }}
                                helperText={touched.email && errors.email}
                                error={touched.email && Boolean(errors.email)}
                            />

                            <TextField
                                fullWidth
                                id="password"
                                name="password"
                                label={t("CONSTANT.PasswordLabel")}
                                value={values.password}
                                onChange={handleChange}
                                InputLabelProps={{ shrink: true }}
                                type={togglePassword ? "text" : "password"}
                                error={
                                    touched.password && Boolean(errors.password)
                                }
                                helperText={touched.password && errors.password}
                                InputProps={{
                                    classes: { input: classes.inputFieldStyle },
                                    endAdornment: (
                                        <InputAdornment position="end">
                                            <IconButton
                                                edge="end"
                                                onClick={handleTogglePassword}
                                            >
                                                {togglePassword ? (
                                                    <VisibilityOff />
                                                ) : (
                                                    <Visibility />
                                                )}
                                            </IconButton>
                                        </InputAdornment>
                                    ),
                                }}
                            />

                            <LinkStyle href="/forgot-password">
                                {t("LOGIN.ForgotPasswordLabel")}
                            </LinkStyle>

                            <Button
                                size="medium"
                                type="submit"
                                color="primary"
                                variant="contained"
                                sx={{ fontWeight: 600 }}
                            >
                                {t("LOGIN.ButtonLabel")}
                            </Button>

                            <Divider />

                            <Button
                                size="medium"
                                variant="outlined"
                                sx={{ fontWeight: 600 }}
                                onClick={() => onGoogleLogin()}
                            >
                                {t("LOGIN.ButtonWithGoogle")}
                            </Button>

                            <Typography
                                variant="body2"
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
