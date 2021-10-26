import {
    Box,
    Stack,
    Button,
    TextField,
    Typography,
    InputAdornment,
    IconButton,
} from "@mui/material";
import * as Yup from "yup";
import { Formik } from "formik";
import { useState } from "react";
import withAuthWrapperHOC from "../withAuthWrapperHOC";
import { Visibility, VisibilityOff } from "@mui/icons-material";
import { ContainerJustifySpaceBtw, globalUseStyles } from "../../styles";
import PasswordRequirementIndicator from "../PasswordRequirementIndicator";
import { useTranslation } from "react-i18next";
import "../../i18n/i18n";
const ResetPasswordSchema = Yup.object({
    newPassword: Yup.string().required("Password is required"),
});

const ResetPasswordValue = {
    newPassword: "",
};

type ResetPasswordFormProps = {
    onSubmit: Function;
    onCancel: Function;
};

const ResetPasswordForm = ({ onSubmit, onCancel }: ResetPasswordFormProps) => {
    const classes = globalUseStyles();
    const { t } = useTranslation();
    const [togglePassword, setTogglePassword] = useState(false);
    const handleTogglePassword = () => {
        setTogglePassword(prev => !prev);
    };
    return (
        <Box width="100%">
            <Formik
                initialValues={ResetPasswordValue}
                validationSchema={ResetPasswordSchema}
                onSubmit={async values => onSubmit(values)}
            >
                {({ values, errors, touched, handleChange, handleSubmit }) => (
                    <form onSubmit={handleSubmit}>
                        <Stack spacing={"18px"}>
                            <Typography variant="h3">
                                {t("RESET_PASSWORD.FormTitle")}
                            </Typography>

                            <TextField
                                fullWidth
                                id="newPassword"
                                name="newPassword"
                                label={t("CONSTANT.NewPasswordLabel")}
                                value={values.newPassword}
                                className={classes.inputFieldBorder}
                                onChange={handleChange}
                                InputLabelProps={{ shrink: true }}
                                type={togglePassword ? "text" : "password"}
                                error={
                                    touched.newPassword &&
                                    Boolean(errors.newPassword)
                                }
                                helperText={
                                    touched.newPassword && errors.newPassword
                                }
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

                            <PasswordRequirementIndicator
                                password={values.newPassword}
                            />

                            <ContainerJustifySpaceBtw>
                                <Button
                                    size="large"
                                    variant="text"
                                    sx={{ fontWeight: 600 }}
                                    onClick={() => onCancel()}
                                >
                                    {t("CONSTANT.CancelButtonLable")}
                                </Button>

                                <Button
                                    size="large"
                                    type="submit"
                                    color="primary"
                                    variant="contained"
                                    sx={{ fontWeight: 600 }}
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
