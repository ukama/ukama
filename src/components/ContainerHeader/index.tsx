import { Typography, Button, IconButton, Paper, Stack } from "@mui/material";
import { HorizontalContainerJustify, HorizontalContainer } from "../../styles";
import { styled } from "@mui/material/styles";
import { colors } from "../../theme";
import SearchIcon from "@mui/icons-material/Search";
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
        paddingLeft: `calc(1em + ${theme.spacing(1)})`,
        width: "100%",
        [theme.breakpoints.up("md")]: {
            width: "20ch",
        },
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
                <Paper
                    sx={{
                        mr: 1,
                        border: `1px solid ${colors.darkGray}`,
                    }}
                    elevation={0}
                >
                    <Stack direction="row">
                        <StyledInputBase placeholder="Search…" />
                        <IconButton
                            color="primary"
                            aria-label="simSearch"
                            component="span"
                        >
                            <SearchIcon sx={{ color: colors.darkGray }} />
                        </IconButton>
                    </Stack>
                </Paper>
            )}

            <Button
                variant="contained"
                sx={{ width: "25%", py: "8px" }}
                onClick={() => handleButtonAction()}
            >
                {buttonTitle}
            </Button>
        </HorizontalContainerJustify>
    );
};

export default ContainerHeader;
