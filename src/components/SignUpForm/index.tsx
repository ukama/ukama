import "../../i18n/i18n";
import * as Yup from "yup";
import { Formik } from "formik";
import { useTranslation } from "react-i18next";
import { PasswordFieldWithIndicator } from "..";
import withAuthWrapperHOC from "../withAuthWrapperHOC";
import { LinkStyle, globalUseStyles } from "../../styles";
import { Stack, Button, Divider, TextField, Typography } from "@mui/material";

const signUpSchema = Yup.object({
    email: Yup.string()
        .email("Please enter a valid email")
        .required("Please enter a valid email"),
    name: Yup.string().required("Name is required"),
    password: Yup.string().required("Password is required"),
});

const initialSignUpValue = {
    email: "",
    name: "",
    password: "",
};

type SignUpFormProps = {
    onSubmit: Function;
    onGoogleSignUp: Function;
};
const SignUpForm = ({ onSubmit, onGoogleSignUp }: SignUpFormProps) => {
    const classes = globalUseStyles();
    const { t } = useTranslation();

    return (
        <Formik
            validationSchema={signUpSchema}
            initialValues={initialSignUpValue}
            onSubmit={async values => onSubmit(values)}
        >
            {({
                errors,
                touched,
                values,
                handleChange,
                handleSubmit,
                handleBlur,
            }) => (
                <>
                    <form onSubmit={handleSubmit}>
                        <Typography variant="h5" sx={{ mb: "24px" }}>
                            {t("SIGNUP.FormTitle")}
                        </Typography>
                        <Stack spacing={"12px"}>
                            <TextField
                                fullWidth
                                id="name"
                                name="name"
                                value={values.name}
                                onBlur={handleBlur}
                                onChange={handleChange}
                                sx={{ mb: "12px  !important" }}
                                label={t("CONSTANT.NameLabel")}
                                InputLabelProps={{ shrink: true }}
                                InputProps={{
                                    classes: {
                                        input: classes.inputFieldStyle,
                                    },
                                }}
                                helperText={touched.name && errors.name}
                                error={touched.name && Boolean(errors.name)}
                            />

                            <TextField
                                fullWidth
                                id="email"
                                name="email"
                                value={values.email}
                                onBlur={handleBlur}
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
                                withIndicator={true}
                                value={values.password}
                                handleChange={handleChange}
                                label={t("CONSTANT.PasswordLabel")}
                                fieldStyle={classes.inputFieldStyle}
                            />

                            <Button
                                id="signUpButton"
                                size="medium"
                                type="submit"
                                color="primary"
                                variant="contained"
                            >
                                {t("SIGNUP.ButtonLabel")}
                            </Button>

                            <Divider />

                            <Button
                                size="medium"
                                variant="outlined"
                                onClick={() => onGoogleSignUp()}
                            >
                                {t("SIGNUP.ButtonWithGoogleLabel")}
                            </Button>

                            <Typography
                                variant="caption"
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
