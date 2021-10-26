import {
    Stack,
    Button,
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
import PasswordRequirementIndicator from "../PasswordRequirementIndicator";
import { useTranslation } from "react-i18next";
import "../../i18n/i18n";
const signUpSchema = Yup.object({
    email: Yup.string()
        .email("Please enter a valid email")
        .required("Email is required"),
    password: Yup.string().required("Password is required"),
});

const initialSignUpValue = {
    email: "",
    password: "",
};

type SignUpFormProps = {
    onSubmit: Function;
    onGoogleSignUp: Function;
};
const SignUpForm = ({ onSubmit, onGoogleSignUp }: SignUpFormProps) => {
    const classes = globalUseStyles();
    const { t } = useTranslation();
    const [togglePassword, setTogglePassword] = useState(false);
    const handleTogglePassword = () => {
        setTogglePassword(prev => !prev);
    };

    return (
        <Formik
            validationSchema={signUpSchema}
            initialValues={initialSignUpValue}
            onSubmit={async values => onSubmit(values)}
        >
            {({ errors, touched, values, handleChange, handleSubmit }) => (
                <>
                    <form onSubmit={handleSubmit}>
                        <Stack spacing={"18px"}>
                            <Typography variant="h3">
                                {t("SIGNUP.FormTitle")}
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
                                onChange={event => {
                                    handleChange(event);
                                }}
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
                                                    <Visibility />
                                                ) : (
                                                    <VisibilityOff />
                                                )}
                                            </IconButton>
                                        </InputAdornment>
                                    ),
                                }}
                            />

                            <PasswordRequirementIndicator
                                password={values.password}
                            />

                            <Button
                                size="medium"
                                type="submit"
                                color="primary"
                                variant="contained"
                                sx={{ fontWeight: 600 }}
                            >
                                {t("SIGNUP.ButtonLabel")}
                            </Button>

                            <Divider />

                            <Button
                                size="medium"
                                variant="outlined"
                                sx={{
                                    fontWeight: 600,
                                    marginTop: "0px !important",
                                }}
                                onClick={() => onGoogleSignUp()}
                            >
                                {t("SIGNUP.ButtonWithGoogleLabel")}
                            </Button>

                            <Typography
                                variant="body2"
                                sx={{
                                    letterSpacing: "0.4px",
                                }}
                            >
                                {t("SIGNUP.FooterNoteText")}
                                <LinkStyle href="/login">
                                    {t("CONSTANT.LinkLabel")}
                                </LinkStyle>
                            </Typography>
                        </Stack>
                    </form>
                </>
            )}
        </Formik>
    );
};

export default withAuthWrapperHOC(SignUpForm);
