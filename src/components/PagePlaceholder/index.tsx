import { colors } from "../../theme";
import { LinkStyle } from "../../styles";
import { PagePlaceholderSvg } from "../../assets/svg";
import { Button, Stack, Typography } from "@mui/material";
import { isDarkmode } from "../../recoil";
import { useRecoilValue } from "recoil";

type PagePlaceholderProps = {
    linkText?: string;
    hyperlink?: string;
    description: string;
    buttonTitle?: string;
    handleAction?: Function;
    showActionButton?: boolean;
};

const PagePlaceholder = ({
    hyperlink = "",
    linkText = "",
    description = "",
    buttonTitle = "",
    showActionButton = false,
    handleAction = () => {
        //Default behaviour
    },
}: PagePlaceholderProps) => {
    const _isDarkmode = useRecoilValue(isDarkmode);
    return (
        <Stack
            spacing={4}
            sx={{
                height: "100%",
                borderRadius: "5px",
                alignItems: "center",
                justifyContent: "center",
                p: 10,
            }}
        >
            <PagePlaceholderSvg
                color={_isDarkmode ? colors.greyish : colors.whiteGrey}
                color2={_isDarkmode ? colors.nightGrey12 : colors.white}
            />
            <Typography variant="body1">
                {`${description} `}
                {hyperlink && (
                    <LinkStyle
                        href={hyperlink}
                        sx={{
                            typography: "body1",
                        }}
                    >
                        {linkText}
                    </LinkStyle>
                )}
            </Typography>

            {showActionButton && (
                <Button
                    variant="contained"
                    sx={{ width: 190 }}
                    onClick={() => handleAction()}
                >
                    {buttonTitle}
                </Button>
            )}
        </Stack>
    );
};

export default PagePlaceholder;
