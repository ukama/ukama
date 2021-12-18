import {
    Typography,
    Button,
    IconButton,
    Paper,
    Stack,
    Grid,
} from "@mui/material";
import { styled } from "@mui/material/styles";
import { colors } from "../../theme";
import SearchIcon from "@mui/icons-material/Search";
import InputBase from "@mui/material/InputBase";
import { useState, useEffect } from "react";
type ContainerHeaderProps = {
    title?: string;
    stats?: string;
    buttonTitle?: string;
    handleButtonAction: Function;
    withSearchBox?: boolean;
    values?: any;
    handleSearchChange?: Function;
};

const StyledInputBase = styled(InputBase)(({ theme }) => ({
    color: "inherit",
    "& .MuiInputBase-input": {
        paddingLeft: `calc(1em + ${theme.spacing(1)})`,
        width: "100%",
    },
}));

const ContainerHeader = ({
    title,
    stats,
    withSearchBox,
    buttonTitle,
    handleButtonAction,
    handleSearchChange = () => {
        /* Default empty function */
    },
}: ContainerHeaderProps) => {
    const [currentSearchValue, setCurrentSearchValue] = useState<any>("");

    useEffect(() => {
        handleSearchChange(currentSearchValue);
    }, [currentSearchValue]);

    return (
        <Grid container spacing={2} sx={{ mb: 2 }}>
            <Grid item xs={12} md={6}>
                <Stack
                    direction="row"
                    spacing={2}
                    sx={{ alignItems: "center" }}
                >
                    <Typography variant="h6">{title}</Typography>
                    <Typography
                        variant="subtitle2"
                        letterSpacing="4px"
                        color={colors.empress}
                    >
                        &#40;{stats}&#41;
                    </Typography>
                </Stack>
            </Grid>

            <Grid item xs={12} md={6} container>
                <Grid container spacing={1} sx={{ alignItems: "center" }}>
                    {withSearchBox && (
                        <Grid item xs={6}>
                            <Paper
                                sx={{
                                    border: `1px solid ${colors.darkGray}`,
                                    width: "100%",
                                }}
                                elevation={0}
                            >
                                <Stack
                                    direction="row"
                                    justifyContent="space-between"
                                >
                                    <StyledInputBase
                                        placeholder="Search…"
                                        value={currentSearchValue}
                                        onChange={(e: any) =>
                                            setCurrentSearchValue(
                                                e.target.value
                                            )
                                        }
                                    />
                                    <IconButton
                                        color="primary"
                                        aria-label="simSearch"
                                        component="span"
                                    >
                                        <SearchIcon
                                            sx={{
                                                fill: colors.darkGray,
                                                color: colors.darkGray,
                                            }}
                                            fontSize="small"
                                        />
                                    </IconButton>
                                </Stack>
                            </Paper>
                        </Grid>
                    )}

                    <Grid
                        item
                        container
                        xs={withSearchBox ? 6 : 12}
                        justifyContent="flex-end"
                    >
                        <Button
                            variant="contained"
                            sx={{ width: withSearchBox ? "100%" : null }}
                            onClick={() => handleButtonAction()}
                        >
                            {buttonTitle}
                        </Button>
                    </Grid>
                </Grid>
            </Grid>
        </Grid>
    );
};

export default ContainerHeader;
