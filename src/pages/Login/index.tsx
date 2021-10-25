import { Button, TextField } from "@mui/material";
import { FormikValues } from "formik";
import { LoginForm, FormDialog } from "../../components";
import { CenterContainer } from "../../styles";
import { useState } from "react";
const Login = () => {
    // eslint-disable-next-line no-unused-vars

    const handleSubmit = (values: FormikValues) => {};
    return (
        <>
            <CenterContainer>
                <LoginForm
                    onGoogleLogin={() => {}}
                    onSubmit={(val: any) => handleSubmit(val)}
                />
            </CenterContainer>
        </>
    );
};

export default Login;
