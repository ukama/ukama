import { theme } from "./theme";
import Router from "./router/Router";
import { routes } from "./router/config";
import { CssBaseline } from "@mui/material";
import { ThemeProvider } from "@emotion/react";
import { BrowserRouter } from "react-router-dom";
import Layout from "./layout";
import { useState } from "react";

const App = () => {
    const [login] = useState(true);
    return (
        <ThemeProvider theme={theme}>
            <CssBaseline />
            <BrowserRouter>
                {login ? (
                    <Layout>
                        <Router routes={routes} />
                    </Layout>
                ) : (
                    <Router routes={routes} />
                )}
            </BrowserRouter>
        </ThemeProvider>
    );
};

export default App;
