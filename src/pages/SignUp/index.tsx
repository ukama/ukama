import { Typography } from "@mui/material";
import { CenterContainer } from "../../styles/welcome";
import SignUpForm from "../../components/SignUpForm";
import withAuthWrapperHOC from "../../components/withAuthWrapperHOC/index";
const SignUp = () => {
    return (
        <CenterContainer>
            <Typography variant="h2">signUp page</Typography>
            <SignUpForm />
        </CenterContainer>
    );
};

export default SignUp;
