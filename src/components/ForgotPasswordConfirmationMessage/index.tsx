import withAuthWrapperHOC from "../withAuthWrapperHOC";
import { Box, Typography, Stack } from "@mui/material";
import { LinkStyle, MessageContainer } from "../../styles";
type ForgotPasswordConfirmationMessageProps = {
    email: string;
};
const ForgotPasswordConfirmationMessage = ({
    email,
}: ForgotPasswordConfirmationMessageProps) => {
    return (
        <Box width="100%">
            <Stack spacing={"18px"}></Stack>
            <MessageContainer>
                <Typography variant="h5" color="initial">
                    RECOVER PASSWORD
                </Typography>
            </MessageContainer>

            <MessageContainer>
                <Typography variant="body1" sx={{ letterSpacing: "0.15 px" }}>
                    If an account with the email {`"${email}"`} exists, an email
                    will be sent with further instructions.
                    <b>{` Link will expire in 30 minutes.`}</b>
                </Typography>
            </MessageContainer>
            <LinkStyle href="/login">RETURN TO LOGIN</LinkStyle>
        </Box>
    );
};

export default withAuthWrapperHOC(ForgotPasswordConfirmationMessage);
