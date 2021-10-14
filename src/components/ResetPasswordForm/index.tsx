import * as Yup from "yup";
import { Formik } from "formik";
import withAuthWrapperHOC from "../withAuthWrapperHOC";
import { Visibility, VisibilityOff } from "@mui/icons-material";
import { ContainerJustifySpaceBtw, globalUseStyles } from "../../styles";
import {
    Box,
    Stack,
    Button,
    TextField,
    Typography,
    InputAdornment,
    IconButton,
} from "@mui/material";
import { useState } from "react";

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
                                Recover Password
                            </Typography>

                            <TextField
                                fullWidth
                                id="newPassword"
                                name="newPassword"
                                label="New Password"
                                value={values.newPassword}
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
                                                    <Visibility />
                                                ) : (
                                                    <VisibilityOff />
                                                )}
                                            </IconButton>
                                        </InputAdornment>
                                    ),
                                }}
                            />

                            <ContainerJustifySpaceBtw>
                                <Button
                                    size="large"
                                    variant="text"
                                    sx={{ fontWeight: 600 }}
                                    onClick={() => onCancel()}
                                >
                                    CANCEL
                                </Button>

                                <Button
                                    size="large"
                                    type="submit"
                                    color="primary"
                                    variant="contained"
                                    sx={{ fontWeight: 600 }}
                                >
                                    CHANGE PASSWORD{` & `}LOGIN
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
