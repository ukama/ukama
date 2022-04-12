import { Typography, Button } from "@mui/material";
import { HorizontalContainerJustify, HorizontalContainer } from "../../styles";

type TableHeaderProps = {
    title: string;
    buttonTitle: string;
    handleButtonAction: Function;
};

const TableHeader = ({
    title,
    buttonTitle,
    handleButtonAction,
}: TableHeaderProps) => {
    return (
        <HorizontalContainerJustify sx={{ marginBottom: "18px" }}>
            <HorizontalContainer>
                <Typography variant="h6" marginRight="2px">
                    {title}
                </Typography>
            </HorizontalContainer>
            <Button
                variant="outlined"
                sx={{ width: "144px" }}
                onClick={() => handleButtonAction()}
            >
                {buttonTitle}
            </Button>
        </HorizontalContainerJustify>
    );
};

export default TableHeader;
