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
                <Typography variant="h6">{title}</Typography>
                <Typography
                    p="0px 8px"
                    variant="subtitle2"
                    fontWeight={600}
                    letterSpacing="6px"
                    color={colors.empress}
                >
                    ● {stats}
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
