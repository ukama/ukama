import colors from "./colors";
import themePalette from "./themePalette";
import { createTheme, PaletteOptions } from "@mui/material/styles";
const theme = createTheme({
    typography: {
        fontFamily: `Work Sans`,
        allVariants: {
            color: colors.primary,
        },
    },
    palette: themePalette as PaletteOptions,
});

export { theme, colors };
