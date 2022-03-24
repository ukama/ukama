import { colors } from "../../theme";
import { styled } from "@mui/material/styles";
import InputBase from "@mui/material/InputBase";
import SearchIcon from "@mui/icons-material/Search";
import { Grid, Stack, Button, Typography } from "@mui/material";

import { useState, useEffect } from "react";
type ContainerHeaderProps = {
    title?: string;
    stats?: string;
    buttonSize?: "small" | "medium" | "large";
    buttonTitle?: string;
    showButton?: boolean;
    showSearchBox?: boolean;
    handleSearchChange?: Function;
    handleButtonAction?: Function;
};

const StyledInputBase = styled(InputBase)(() => ({
    color: "inherit",
    "& .MuiInputBase-input": {
        width: "100%",
    },
}));

const ContainerHeader = ({
    title,
    stats,
    buttonTitle,
    showButton = false,
    buttonSize = "large",
    showSearchBox = false,
    handleSearchChange = () => {
        /* Default empty function */
    },
    handleButtonAction = () => {
        /* Default function */
    },
}: ContainerHeaderProps) => {
    const [currentSearchValue, setCurrentSearchValue] = useState<string>("");

    useEffect(() => {
        handleSearchChange(currentSearchValue.toLowerCase());
    }, [currentSearchValue]);

    return (
        <Grid container spacing={2} justifyContent="space-between">
            <Grid item xs={12} md={4}>
                <Stack
                    spacing={2}
                    direction="row"
                    sx={{ alignItems: "baseline" }}
                >
                    <Typography variant="h6">{title}</Typography>
                    {stats && (
                        <Typography
                            variant="subtitle2"
                            letterSpacing="4px"
                            color={"textSecondary"}
                        >
                            &#40;{stats}&#41;
                        </Typography>
                    )}
                </Stack>
            </Grid>

            <Grid
                item
                md={8}
                xs={12}
                container
                spacing={3}
                alignItems="center"
                justifyContent="flex-end"
            >
                {showSearchBox && (
                    <Grid item>
                        <StyledInputBase
                            placeholder="Search…"
                            value={currentSearchValue}
                            onChange={(e: any) =>
                                setCurrentSearchValue(e.target.value)
                            }
                            sx={{
                                height: "42px",
                                borderRadius: 2,
                                minWidth: { xs: "100%", md: "300px" },
                                border: `1px solid ${colors.silver}`,
                                padding: "4px 8px 4px 12px !important",
                            }}
                            endAdornment={
                                <SearchIcon fontSize="small" color="primary" />
                            }
                        />
                    </Grid>
                )}
                {showButton && (
                    <Grid item justifyContent="flex-end" display="flex">
                        <Button
                            sx={{
                                px: 4,
                                width: { xs: "100%", md: "fit-content" },
                            }}
                            size={buttonSize}
                            variant="contained"
                            onClick={() => handleButtonAction()}
                        >
                            {buttonTitle}
                        </Button>
                    </Grid>
                )}
            </Grid>
        </Grid>
    );
};

export default ContainerHeader;
