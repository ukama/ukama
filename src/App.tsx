import Layout from "./layout";
import config from "./config";
import { theme } from "./theme";
import { useEffect } from "react";
import Router from "./router/Router";
import { routes } from "./router/config";
import client from "./api/ApolloClient";
import { CenterContainer } from "./styles";
import useWhoami from "./helpers/useWhoami";
import { ThemeProvider } from "@emotion/react";
import { ApolloProvider } from "@apollo/client";
import { BrowserRouter } from "react-router-dom";
import { CircularProgress, CssBaseline } from "@mui/material";

const App = () => {
    const { loading, response } = useWhoami();

    useEffect(() => {
        if (response && config.ENVIROMENT === "production") {
            if (!response?.isValid) {
                window.close();
                window.location.replace(`${config.REACT_APP_AUTH_URL}`);
            }
        }
    }, [response]);

    if (loading)
        return (
            <CenterContainer>
                <CircularProgress />
            </CenterContainer>
        );

    return (
        <ApolloProvider client={client}>
            <ThemeProvider theme={theme}>
                <CssBaseline />
                <BrowserRouter>
                    <Layout>
                        <Router routes={routes} />
                    </Layout>
                </BrowserRouter>
            </ThemeProvider>
        </ApolloProvider>
    );
};

export default App;
