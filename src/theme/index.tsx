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
        subtitle1: {
            fontWeight: 500,
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
        MuiButton: {
            styleOverrides: {
                contained: {
                    fontWeight: 600,
                    color: colors.white,
                    letterSpacing: "0.4px",
                    boxShadow:
                        "0px 3px 1px -2px rgba(0, 0, 0, 0.2), 0px 2px 2px rgba(0, 0, 0, 0.14), 0px 1px 5px rgba(0, 0, 0, 0.12)",
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
