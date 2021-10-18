import withAuthWrapperHOC from "../withAuthWrapperHOC";
import { Box, Typography, Stack } from "@mui/material";
import { LinkStyle, MessageContainer } from "../../styles";
import { useTranslation } from "react-i18next";
import "../../i18n/i18n";
type ForgotPasswordConfirmationMessageProps = {
    email: string;
};
const ForgotPasswordConfirmationMessage = ({
    email,
}: ForgotPasswordConfirmationMessageProps) => {
    const { t } = useTranslation();

    return (
        <Box width="100%">
            <Stack spacing={"18px"}></Stack>
            <MessageContainer>
                <Typography variant="h5" color="initial">
                    {t("RECOVER_PASSWORD.RecoverPasswordTitle")}
                </Typography>
            </MessageContainer>

            <MessageContainer>
                <Typography variant="body1" sx={{ letterSpacing: "0.15 px" }}>
                    {t("RECOVER_PASSWORD.RecoveryNote")}
                    {email}

                    <b>{t("RECOVER_PASSWORD.RecoveryImportantNote")}</b>
                </Typography>
            </MessageContainer>
            <LinkStyle href="/login">
                {t("CONSTANT.ReturnToLoginLink")}
            </LinkStyle>
        </Box>
    );
};

export default withAuthWrapperHOC(ForgotPasswordConfirmationMessage);
