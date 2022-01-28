import colors from "./colors";
import { createTheme } from "@mui/material/styles";

const theme = (isDarkmode: boolean) =>
    createTheme({
        typography: {
            fontFamily: "Rubik, sans-serif;",
            subtitle1: { fontFamily: "Work Sans, sans-serif" },
            subtitle2: { fontFamily: "Work Sans, sans-serif" },
            body1: {
                fontFamily: "Work Sans, sans-serif",
                letterSpacing: "-0.02em",
            },
            body2: {
                fontFamily: "Work Sans, sans-serif",
                letterSpacing: "-0.02em",
            },
            caption: {
                fontFamily: "Work Sans, sans-serif",
            },
        },
        palette: {
            mode: isDarkmode ? "dark" : "light",
            text: {
                primary: isDarkmode ? colors._white : colors.black,
                secondary: isDarkmode ? colors._white38 : colors._vulcan60,
            },
            background: {
                default: isDarkmode ? colors._nightGrey : colors._solitude,
                paper: isDarkmode ? colors._nightGrey5 : colors._white,
            },
            primary: {
                main: colors.primaryMain,
                light: colors.primaryLight,
                dark: colors.primaryDark,
            },
            secondary: {
                main: colors.secondaryMain,
                light: colors.secondaryLight,
                dark: colors.secondaryDark,
            },
            error: {
                main: colors._error,
            },
        },
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
                        // "&:hover": {
                        //     backgroundColor: colors.aliceBlue,
                        // },
                    },
                },
            },
            MuiIconButton: {
                styleOverrides: {
                    root: {
                        "&:hover": {
                            backgroundColor: "transparent",
                        },
                        "&:hover svg path": {
                            fill: colors.primaryMain,
                        },
                    },
                },
            },
            MuiButton: {
                styleOverrides: {
                    contained: {
                        fontWeight: 500,
                        color: colors._white,
                        letterSpacing: "0.4px",
                        boxShadow:
                            "0px 3px 1px -2px rgba(0, 0, 0, 0.2), 0px 2px 2px rgba(0, 0, 0, 0.14), 0px 1px 5px rgba(0, 0, 0, 0.12)",
                    },
                    text: {
                        padding: "0px",
                        minWidth: "auto",
                    },
                },
            },
            MuiSelect: {
                styleOverrides: {
                    select: {
                        fontSize: "0.875rem",
                        ":focus": {
                            backgroundColor: "transparent",
                        },
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
            MuiTableCell: {
                styleOverrides: {
                    root: {
                        backgroundColor: "transparent",
                    },
                },
            },
        },
    });

export { theme, colors };
