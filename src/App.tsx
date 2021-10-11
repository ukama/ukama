import { theme } from "./theme";
import Router from "./router/Router";
import { routes } from "./router/config";
import { BrowserRouter } from "react-router-dom";
import { ThemeProvider } from "@emotion/react";
import { CssBaseline } from "@mui/material";

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
