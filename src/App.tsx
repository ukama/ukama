import Layout from "./layout";
import { theme } from "./theme";
import Router from "./router/Router";
import { useRecoilValue } from "recoil";
import { routes } from "./router/config";
import { CssBaseline } from "@mui/material";
import { ThemeProvider } from "@emotion/react";
import { BrowserRouter } from "react-router-dom";
import { isLoginAtom } from "./recoil";

const App = () => {
    const isLogin = useRecoilValue(isLoginAtom);
    return (
        <ThemeProvider theme={theme}>
            <CssBaseline />
            <BrowserRouter>
                {isLogin ? (
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
