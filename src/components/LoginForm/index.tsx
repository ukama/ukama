import * as Yup from "yup";
import { useState } from "react";
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
import { Formik } from "formik";
import { LinkStyle } from "../../styles";
import withAuthWrapperHOC from "../withAuthWrapperHOC";
import { Visibility, VisibilityOff } from "@mui/icons-material";

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
    onSubmit: any;
};

const LoginForm = ({ onSubmit }: LoginFormProps) => {
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
                                Log in to Ukama
                            </Typography>

                            <TextField
                                fullWidth
                                id="email"
                                name="email"
                                label="Email"
                                value={values.email}
                                onChange={handleChange}
                                error={touched.email && Boolean(errors.email)}
                                helperText={touched.email && errors.email}
                            />

                            <TextField
                                fullWidth
                                id="password"
                                name="password"
                                label="Password"
                                value={values.password}
                                onChange={handleChange}
                                type={togglePassword ? "text" : "password"}
                                error={
                                    touched.password && Boolean(errors.password)
                                }
                                helperText={touched.password && errors.password}
                                InputProps={{
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

                            <LinkStyle href="/forgot-password">
                                Forgot Password?
                            </LinkStyle>

                            <Button
                                size="large"
                                type="submit"
                                color="primary"
                                variant="contained"
                                sx={{ fontWeight: 600 }}
                            >
                                LOG IN
                            </Button>

                            <Divider />

                            <Button
                                size="large"
                                variant="outlined"
                                sx={{ fontWeight: 600 }}
                            >
                                LOG IN WITH GOOGLE
                            </Button>

                            <Typography
                                variant="body2"
                                sx={{ letterSpacing: "0.4px" }}
                            >
                                {`Don't have an account? `}
                                <LinkStyle href="/signup">
                                    Sign up insted.
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
