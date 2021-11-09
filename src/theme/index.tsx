import colors from "./colors";
import themePalette from "./themePalette";
import { createTheme, PaletteOptions } from "@mui/material/styles";
const theme = createTheme({
    typography: {
        fontFamily: "Rubik",
        subtitle1: { fontFamily: "Work Sans" },
        subtitle2: { fontFamily: "Work Sans" },
        body1: {
            fontFamily: "Work Sans",
            letterSpacing: "-0.02em",
        },
        body2: {
            fontFamily: "Work Sans",
            letterSpacing: "-0.02em",
        },
        caption: {
            fontFamily: "Work Sans",
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
                    margin: "12px 0px",
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
        MuiButton: {
            styleOverrides: {
                contained: {
                    fontWeight: 500,
                    color: colors.white,
                    letterSpacing: "0.4px",
                    boxShadow:
                        "0px 3px 1px -2px rgba(0, 0, 0, 0.2), 0px 2px 2px rgba(0, 0, 0, 0.14), 0px 1px 5px rgba(0, 0, 0, 0.12)",
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
        MuiDrawer: {
            styleOverrides: {
                paper: {
                    borderRight: "none",
                },
            },
        },
    },
});

export { theme, colors };
