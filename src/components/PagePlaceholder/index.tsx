import { colors } from "../../theme";
import { LinkStyle } from "../../styles";
import { PagePlaceholderSvg } from "../../assets/svg";
import { Button, Stack, styled, Typography, Box } from "@mui/material";
const Container = styled(Box)((props: any) => ({
    background:
        props.theme.palette.mode === "dark" ? colors.darkGreen12 : colors.white,
}));
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
    return (
        <Container>
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
                <PagePlaceholderSvg />
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
        </Container>
    );
};

export default PagePlaceholder;
