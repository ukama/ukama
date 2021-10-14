import { FormikValues } from "formik";
import { useHistory } from "react-router";
import { ForgotPasswordForm } from "../../components";
import { CenterContainer } from "../../styles/welcome";
const ForgotPassword = () => {
    const history = useHistory();
    // eslint-disable-next-line no-unused-vars
    const handleSubmit = (values: FormikValues) => {};
    const handleBack = () => history.goBack();

    return (
        <CenterContainer>
            <ForgotPasswordForm
                onBack={handleBack}
                onSubmit={(val: any) => handleSubmit(val)}
            />
        </CenterContainer>
    );
};

export default ForgotPassword;
