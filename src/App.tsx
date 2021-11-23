import Layout from "./layout";
import { theme } from "./theme";
import Router from "./router/Router";
import { isLoginAtom } from "./recoil";
import { useRecoilValue } from "recoil";
import client from "./api/ApolloClient";
import { routes } from "./router/config";
import { CssBaseline } from "@mui/material";
import { ThemeProvider } from "@emotion/react";
import { ApolloProvider } from "@apollo/client";
import { BrowserRouter } from "react-router-dom";

const App = () => {
    const isLogin = useRecoilValue(isLoginAtom);
    return (
        <ApolloProvider client={client}>
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
        </ApolloProvider>
    );
};

export default App;
