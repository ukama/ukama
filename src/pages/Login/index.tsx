import { FormikValues } from "formik";
import { useSetRecoilState } from "recoil";
import { LoginForm } from "../../components";
import { isLoginAtom } from "../../recoil";
import { CenterContainer } from "../../styles";
const Login = () => {
    const setIsLogin = useSetRecoilState(isLoginAtom);
    // eslint-disable-next-line no-unused-vars
    const handleSubmit = (values: FormikValues) => {
        setIsLogin(true);
    };
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
