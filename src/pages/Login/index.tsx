import { FormikValues } from "formik";
import { LoginForm } from "../../components";
import { CenterContainer } from "../../styles/login";
const Login = () => {
    // eslint-disable-next-line no-unused-vars
    const onSubmit = (values: FormikValues) => {};
    return (
        <CenterContainer>
            <LoginForm onSubmit={onSubmit} />
        </CenterContainer>
    );
};

export default Login;
