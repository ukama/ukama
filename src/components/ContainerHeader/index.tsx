import { Typography, Button } from "@mui/material";
import { HorizontalContainerJustify, HorizontalContainer } from "../../styles";
import { styled } from "@mui/material/styles";
import { colors } from "../../theme";
// import SearchIcon from "@mui/icons-material/Search";
import InputBase from "@mui/material/InputBase";
type ContainerHeaderProps = {
    title?: string;
    stats?: string;
    buttonTitle?: string;
    handleButtonAction: Function;
    withSearchBox?: boolean;
};

const StyledInputBase = styled(InputBase)(({ theme }) => ({
    color: "inherit",
    "& .MuiInputBase-input": {
        padding: theme.spacing(1, 1, 1, 0),
        border: "1px solid #E0E0E0",
        borderRadius: "4px",
        paddingLeft: `calc(1em + ${theme.spacing(1)})`,
        transition: theme.transitions.create("width"),
        width: "100%",
        [theme.breakpoints.up("md")]: {
            width: "20ch",
        },
    },
}));
const Search = styled("div")(({ theme }) => ({
    position: "relative",
    borderRadius: theme.shape.borderRadius,
    backgroundColor: "white",
    "&:hover": {
        backgroundColor: "white",
    },
    marginRight: theme.spacing(2),
    marginLeft: 0,
    width: "100%",
    [theme.breakpoints.up("sm")]: {
        marginLeft: theme.spacing(3),
        width: "auto",
    },
}));
const ContainerHeader = ({
    title,
    stats,
    withSearchBox,
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
            {withSearchBox && (
                <Search>
                    <StyledInputBase
                        placeholder="Search…"
                        inputProps={{ "aria-label": "search" }}
                    />
                </Search>
            )}

            <Button
                variant="contained"
                sx={{ width: "25%" }}
                onClick={() => handleButtonAction()}
            >
                {buttonTitle}
            </Button>
        </HorizontalContainerJustify>
    );
};

export default ContainerHeader;
