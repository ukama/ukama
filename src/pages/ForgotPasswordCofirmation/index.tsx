import { ForgotPasswordConfirmationMessage } from "../../components";
import { CenterContainer } from "../../styles";
const ForgotPasswordConfirmation = () => {
    return (
        <CenterContainer>
            <ForgotPasswordConfirmationMessage email={`Brackley@ukama.com`} />
        </CenterContainer>
    );
};

export default ForgotPasswordConfirmation;
