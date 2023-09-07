import colors from "./colors";
import { createTheme, responsiveFontSizes } from "@mui/material/styles";

const theme = (isDarkmode: boolean) =>
    responsiveFontSizes(
        createTheme({
            typography: {
                fontSize: 16,
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
                    primary: isDarkmode ? colors.white : colors.black,
                    secondary: isDarkmode ? colors.white70 : colors.black70,
                    disabled: isDarkmode ? colors.white38 : colors.black38,
                },
                background: {
                    default: isDarkmode ? colors.nightGrey : colors.solitude,
                    paper: isDarkmode ? colors.nightGrey5 : colors.white,
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
                    main: colors.error,
                },
            },
            breakpoints: {
                values: {
                    xs: 0,
                    sm: 600,
                    md: 1280,
                    lg: 1440,
                    xl: 1600,
                },
                step: 8,
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
                            color: colors.white,
                            letterSpacing: "0.4px",
                            boxShadow:
                                "0px 3px 1px -2px rgba(0, 0, 0, 0.2), 0px 2px 2px rgba(0, 0, 0, 0.14), 0px 1px 5px rgba(0, 0, 0, 0.12)",
                        },
                        text: {
                            padding: "0px",
                            minWidth: "auto",
                            color: isDarkmode ? colors.white : colors.darkGray,
                            ":hover": {
                                color: colors.primaryMain,
                                backgroundColor: "transparent !important",
                            },
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
                        iconStandard: {
                            paddingLeft: "4px",
                        },
                    },
                },
                MuiDialogContent: {
                    styleOverrides: {
                        root: {
                            padding: "8px 24px 24px",
                        },
                    },
                },
                MuiDialogActions: {
                    styleOverrides: {
                        root: {
                            padding: "8px 24px 24px",
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
                        stickyHeader: {
                            backgroundColor: isDarkmode
                                ? "#292929"
                                : colors.white,
                        },
                    },
                },
                MuiPaper: {
                    styleOverrides: {
                        root: {
                            boxShadow: "2px 2px 6px rgba(0, 0, 0, 0.05)",
                        },
                    },
                },
            },
        })
    );

export { theme, colors };
