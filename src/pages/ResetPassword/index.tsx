import { FormikValues } from "formik";
import { useHistory } from "react-router";
import { CenterContainer } from "../../styles";
import { BasicDialog, ResetPasswordForm } from "../../components";
import { useState } from "react";
const ResetPassword = () => {
    const history = useHistory();
    const [successDialog, setSuccessDialog] = useState(false);

    // eslint-disable-next-line no-unused-vars
    const handleSubmit = (values: FormikValues) => {
        setSuccessDialog(true);
    };

    const handelCanceel = () => history.push("/login");

    return (
        <CenterContainer>
            <ResetPasswordForm
                onCancel={handelCanceel}
                onSubmit={(val: any) => handleSubmit(val)}
            />
            <BasicDialog
                isOpen={successDialog}
                title={"Password Changed Successfully"}
                handleClose={() => setSuccessDialog(false)}
                content={"Your password has been changed successfully!"}
            />
        </CenterContainer>
    );
};

export default ResetPassword;
