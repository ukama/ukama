import { Typography, Button } from "@mui/material";
import { HorizontalContainerJustify, HorizontalContainer } from "../../styles";

type TableHeaderProps = {
    title: string;
    buttonTitle?: string;
    showSecondaryButton: boolean;
    handleButtonAction?: any;
};

const TableHeader = ({
    title,
    buttonTitle,
    handleButtonAction,
    showSecondaryButton,
}: TableHeaderProps) => {
    return (
        <HorizontalContainerJustify sx={{ marginBottom: "18px" }}>
            <HorizontalContainer>
                <Typography variant="h6" marginRight="2px">
                    {title}
                </Typography>
            </HorizontalContainer>
            {showSecondaryButton && (
                <Button
                    variant="outlined"
                    sx={{ width: "144px" }}
                    onClick={() => handleButtonAction()}
                >
                    {buttonTitle}
                </Button>
            )}
        </HorizontalContainerJustify>
    );
};

export default TableHeader;
