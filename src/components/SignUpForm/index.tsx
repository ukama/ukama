import {
    Box,
    Button,
    Stack,
    Divider,
    TextField,
    IconButton,
    Typography,
    InputAdornment,
    Grid,
} from "@mui/material";
import * as Yup from "yup";
import { Formik } from "formik";
import { useState } from "react";
import withAuthWrapperHOC from "../withAuthWrapperHOC";
import { LinkStyle, globalUseStyles } from "../../styles";
import { Visibility, VisibilityOff } from "@mui/icons-material";

const signUpSchema = Yup.object({
    email: Yup.string()
        .email("Enter a valid email")
        .required("Email is required"),
    firstName: Yup.string().required("FirstName is required"),
    lastName: Yup.string().required("LastName is required"),

    password: Yup.string().required("Password is required"),
});

const initialSignUpValue = {
    email: "",
    password: "",
    firstName: "",
    lastName: "",
};

type SignUpFormProps = {
    onSubmit: Function;
    onGoogleSignUp: Function;
};

const SignUpForm = ({ onSubmit, onGoogleSignUp }: SignUpFormProps) => {
    const classes = globalUseStyles();
    const [togglePassword, setTogglePassword] = useState(false);
    const handleTogglePassword = () => {
        setTogglePassword(prev => !prev);
    };

    return (
        <Box width="100%">
            <Formik
                validationSchema={signUpSchema}
                initialValues={initialSignUpValue}
                onSubmit={async values => onSubmit(values)}
            >
                {({ values, errors, touched, handleChange, handleSubmit }) => (
                    <form onSubmit={handleSubmit}>
                        <Stack spacing={"18px"}>
                            <Typography variant="h3">
                                Sign up for Ukama
                            </Typography>
                            <Box sx={{ flexGrow: 1 }}>
                                <Grid container spacing={2}>
                                    <Grid item xs={6}>
                                        <TextField
                                            fullWidth
                                            id="firstName"
                                            name="firstName"
                                            label="FirstName"
                                            value={values.firstName}
                                            onChange={handleChange}
                                            InputLabelProps={{ shrink: true }}
                                            InputProps={{
                                                classes: {
                                                    input: classes.inputFieldStyle,
                                                },
                                            }}
                                            helperText={
                                                touched.firstName &&
                                                errors.firstName
                                            }
                                            error={
                                                touched.firstName &&
                                                Boolean(errors.firstName)
                                            }
                                        />
                                    </Grid>
                                    <Grid item xs={6}>
                                        <TextField
                                            fullWidth
                                            id="lastName"
                                            name="lastName"
                                            label="LastName"
                                            value={values.lastName}
                                            onChange={handleChange}
                                            InputLabelProps={{ shrink: true }}
                                            InputProps={{
                                                classes: {
                                                    input: classes.inputFieldStyle,
                                                },
                                            }}
                                            helperText={
                                                touched.lastName &&
                                                errors.lastName
                                            }
                                            error={
                                                touched.lastName &&
                                                Boolean(errors.lastName)
                                            }
                                        />
                                    </Grid>
                                </Grid>
                            </Box>

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
                            <TextField
                                fullWidth
                                id="password"
                                name="password"
                                label="Password"
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
                                                    <Visibility />
                                                ) : (
                                                    <VisibilityOff />
                                                )}
                                            </IconButton>
                                        </InputAdornment>
                                    ),
                                }}
                            />

                            <Button
                                size="large"
                                type="submit"
                                color="primary"
                                variant="contained"
                                sx={{ fontWeight: 600 }}
                            >
                                SIGN UP
                            </Button>

                            <Divider />

                            <Button
                                size="large"
                                variant="outlined"
                                sx={{ fontWeight: 600 }}
                                onClick={() => onGoogleSignUp()}
                            >
                                SIGN UP WITH GOOGLE
                            </Button>

                            <Typography
                                variant="body2"
                                sx={{
                                    letterSpacing: "0.4px",
                                }}
                            >
                                {`By signing up , I am agreeing to the Terms and Conditions and Privacy Policy. 
                                Already have an account? `}
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

export default withAuthWrapperHOC(SignUpForm);
