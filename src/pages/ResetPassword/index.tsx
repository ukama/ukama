import { FormikValues } from "formik";
import { useHistory } from "react-router";
import { CenterContainer } from "../../styles";
import { BasicDialog, ResetPasswordForm } from "../../components";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import "../../i18n/i18n";
const ResetPassword = () => {
    const history = useHistory();
    const [successDialog, setSuccessDialog] = useState(false);
    const { t } = useTranslation();
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
                btnLabel={"Close"}
                isOpen={successDialog}
                handleClose={() => setSuccessDialog(false)}
                title={t("DIALOG_MESSAGE.SuccessTitle")}
                content={t("DIALOG_MESSAGE.SuccessContent")}
            />
        </CenterContainer>
    );
};

export default ResetPassword;
