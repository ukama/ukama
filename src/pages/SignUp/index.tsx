import { FormikValues } from "formik";
import { SignUpForm } from "../../components";
import { CenterContainer } from "../../styles/login";
const SignUp = () => {
    // eslint-disable-next-line no-unused-vars
    const handleSubmit = (values: FormikValues) => {};
    return (
        <CenterContainer>
            <SignUpForm
                onGoogleSignUp={() => {}}
                onSubmit={(val: any) => handleSubmit(val)}
            />
        </CenterContainer>
    );
};

export default SignUp;
