import "../../i18n/i18n";
import * as Yup from "yup";
import { Formik } from "formik";
import { useTranslation } from "react-i18next";
import withAuthWrapperHOC from "../withAuthWrapperHOC";
import { Box, Stack, Button, Typography } from "@mui/material";
import { ContainerJustifySpaceBtw, globalUseStyles } from "../../styles";
import PasswordRequirementIndicator from "../PasswordFieldWithIndicator";

const ResetPasswordSchema = Yup.object({
    password: Yup.string().required("Password is required"),
});

const ResetPasswordValue = {
    password: "",
};

type ResetPasswordFormProps = {
    onSubmit: Function;
    onCancel: Function;
};

const ResetPasswordForm = ({ onSubmit, onCancel }: ResetPasswordFormProps) => {
    const classes = globalUseStyles();
    const { t } = useTranslation();

    return (
        <Box width="100%">
            <Formik
                initialValues={ResetPasswordValue}
                validationSchema={ResetPasswordSchema}
                onSubmit={async values => onSubmit(values)}
            >
                {({
                    values,
                    errors,
                    touched,
                    handleBlur,
                    handleChange,
                    handleSubmit,
                }) => (
                    <form onSubmit={handleSubmit}>
                        <Stack spacing={"18px"}>
                            <Typography variant="h5">
                                {t("RESET_PASSWORD.FormTitle")}
                            </Typography>

                            <PasswordRequirementIndicator
                                errors={errors}
                                touched={touched}
                                onBlur={handleBlur}
                                withIndicator={true}
                                value={values.password}
                                handleChange={handleChange}
                                label={t("CONSTANT.NewPasswordLabel")}
                                fieldStyle={classes.inputFieldStyle}
                            />

                            <ContainerJustifySpaceBtw>
                                <Button
                                    size="medium"
                                    variant="text"
                                    onClick={() => onCancel()}
                                >
                                    {t("CONSTANT.CancelButtonLable")}
                                </Button>

                                <Button
                                    size="medium"
                                    type="submit"
                                    color="primary"
                                    variant="contained"
                                >
                                    {t("RESET_PASSWORD.ButtonLabel")}
                                </Button>
                            </ContainerJustifySpaceBtw>
                        </Stack>
                    </form>
                )}
            </Formik>
        </Box>
    );
};

export default withAuthWrapperHOC(ResetPasswordForm);
