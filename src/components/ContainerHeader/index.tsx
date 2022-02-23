import { colors } from "../../theme";
import { styled } from "@mui/material/styles";
import InputBase from "@mui/material/InputBase";
import SearchIcon from "@mui/icons-material/Search";
import { Grid, Stack, Button, Typography, IconButton } from "@mui/material";

import { useState, useEffect } from "react";
type ContainerHeaderProps = {
    title?: string;
    stats?: string;
    buttonTitle?: string;
    handleButtonAction?: Function;
    showSearchBox?: boolean;
    handleSearchChange?: Function;
    showButton?: boolean;
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
    showButton = false,
    showSearchBox = false,
    buttonTitle,
    handleSearchChange = () => {
        /* Default empty function */
    },
    handleButtonAction = () => {
        /* Default function */
    },
}: ContainerHeaderProps) => {
    const [currentSearchValue, setCurrentSearchValue] = useState<string>("");

    useEffect(() => {
        handleSearchChange(currentSearchValue);
    }, [currentSearchValue]);

    return (
        <Grid container spacing={2} sx={{ mb: 2 }}>
            <Grid item xs={12} md={8}>
                <Stack
                    spacing={2}
                    direction="row"
                    sx={{ alignItems: "baseline" }}
                >
                    <Typography variant="h6">{title}</Typography>
                    <Typography
                        variant="subtitle2"
                        letterSpacing="4px"
                        color={"textSecondary"}
                    >
                        &#40;{stats}&#41;
                    </Typography>
                </Stack>
            </Grid>

            <Grid container item xs={12} md={4} spacing={2}>
                {showSearchBox && (
                    <Grid item xs={8}>
                        <StyledInputBase
                            placeholder="Search…"
                            value={currentSearchValue}
                            onChange={(e: any) =>
                                setCurrentSearchValue(e.target.value)
                            }
                            sx={{
                                width: "100%",
                                height: "48px",
                                borderRadius: 2,
                                border: `1px solid ${colors.silver}`,
                                padding: "4px 8px 4px 12px !important",
                            }}
                            endAdornment={
                                <IconButton
                                    color="primary"
                                    aria-label="simSearch"
                                    component="span"
                                >
                                    <SearchIcon fontSize="small" />
                                </IconButton>
                            }
                        />
                    </Grid>
                )}
                {showButton && (
                    <Grid item xs={4}>
                        <Button
                            variant="contained"
                            sx={{ height: "48px" }}
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
