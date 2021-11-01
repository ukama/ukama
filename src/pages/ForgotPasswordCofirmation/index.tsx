import { ForgotPasswordConfirmationMessage } from "../../components";
import { CenterContainer } from "../../styles";
const ForgotPasswordConfirmation = () => {
    return (
        <CenterContainer>
            <ForgotPasswordConfirmationMessage email={`john@doe.com`} />
        </CenterContainer>
    );
};

export default ForgotPasswordConfirmation;
