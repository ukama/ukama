import { theme } from "./theme";
import Router from "./router/Router";
import { routes } from "./router/config";
import { CssBaseline } from "@mui/material";
import { ThemeProvider } from "@emotion/react";
import { BrowserRouter } from "react-router-dom";

const App = () => {
    return (
        <ThemeProvider theme={theme}>
            <CssBaseline />
            <BrowserRouter>
                <Router routes={routes} />
            </BrowserRouter>
        </ThemeProvider>
    );
};

export default App;
