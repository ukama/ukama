import { FormikValues } from "formik";
import { useHistory } from "react-router";
import { CenterContainer } from "../../styles";
import { ResetPasswordForm } from "../../components";
const ResetPassword = () => {
    const history = useHistory();
    // eslint-disable-next-line no-unused-vars
    const handleSubmit = (values: FormikValues) => {};
    const handelCanceel = () => history.push("/login");
    return (
        <CenterContainer>
            <ResetPasswordForm
                onCancel={handelCanceel}
                onSubmit={(val: any) => handleSubmit(val)}
            />
        </CenterContainer>
    );
};

export default ResetPassword;
