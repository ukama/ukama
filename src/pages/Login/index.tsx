import { FormikValues } from "formik";
import { LoginForm } from "../../components";
import { CenterContainer } from "../../styles";
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
