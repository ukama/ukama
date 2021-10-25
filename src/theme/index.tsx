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
        body2: {
            display: "block",
            lineHeight: "166%",
            fontSize: "0.75rem",
            color: colors.black,
            fontWeight: "normal",
            alignItems: "center",
        },
    },
    palette: themePalette as PaletteOptions,
    components: {
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
                        backgroundColor: "rgba(33, 144, 246, 0.08)",
                    },
                },
            },
        },
    },
});

export { theme, colors };
