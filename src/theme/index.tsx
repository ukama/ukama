import colors from "./colors";
import themePalette from "./themePalette";
import { createTheme, PaletteOptions } from "@mui/material/styles";
const theme = createTheme({
    typography: {
        fontFamily: "Rubik, sans-serif",
        h3: {
            display: "flex",
            fontSize: "1.5rem",
            fontWeight: "normal",
            lineHeight: "133.4%",
            alignItems: "center",
            color: colors.black,
        },
        h6: {
            fontWeight: 600,
        },
        body1: {
            display: "block",
            lineHeight: "157%",
            fontSize: "0.875rem",
            fontWeight: "normal",
            alignItems: "center",
        },
        body2: {
            display: "block",
            lineHeight: "166%",
            fontSize: "0.75rem",
            fontWeight: "normal",
            alignItems: "center",
        },
    },
    palette: themePalette as PaletteOptions,
    components: {
        MuiFormControl: {
            styleOverrides: {
                root: {
                    "&:hover .MuiOutlinedInput-root .MuiOutlinedInput-notchedOutline":
                        {
                            borderColor: colors.hoverColor,
                        },
                },
            },
        },
        MuiDivider: {
            styleOverrides: {
                root: {
                    margin: "12px 0px !important",
                },
            },
        },
        MuiFormHelperText: {
            styleOverrides: {
                contained: {
                    marginLeft: "0px !important",
                },
            },
        },
        MuiListItem: {
            styleOverrides: {
                button: {
                    "&:hover": {
                        backgroundColor: colors.aliceBlue,
                    },
                },
            },
        },
        MuiIconButton: {
            styleOverrides: {
                root: {
                    width: "68px",
                    height: "68px",
                    padding: "0px",
                    "&:hover": {
                        backgroundColor: colors.solitude,
                    },
                    "&:hover svg path": {
                        fill: colors.primary,
                    },
                },
            },
        },
        MuiSelect: {
            styleOverrides: {
                select: {
                    fontSize: "0.875rem",
                    backgroundColor: "transparent",
                },
            },
        },
    },
});

export { theme, colors };
