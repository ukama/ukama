import { Typography, Button } from "@mui/material";
import { HorizontalContainerJustify, HorizontalContainer } from "../../styles";
import { colors } from "../../theme";

type ContainerHeaderProps = {
    title: string;
    stats: string;
    buttonTitle: string;
    handleButtonAction: Function;
};

const ContainerHeader = ({
    title,
    stats,
    buttonTitle,
    handleButtonAction,
}: ContainerHeaderProps) => {
    return (
        <HorizontalContainerJustify sx={{ marginBottom: "18px" }}>
            <HorizontalContainer>
                <Typography variant="h6" marginRight="2px">
                    {title}
                </Typography>
                <Typography
                    p="0px 8px"
                    variant="subtitle2"
                    letterSpacing="4px"
                    color={colors.empress}
                >
                    &#40;{stats}&#41;
                </Typography>
            </HorizontalContainer>
            <Button
                variant="contained"
                sx={{ width: "144px" }}
                onClick={() => handleButtonAction()}
            >
                {buttonTitle}
            </Button>
        </HorizontalContainerJustify>
    );
};

export default ContainerHeader;
