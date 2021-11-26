import Layout from "./layout";
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
import config from "./config";

const App = () => {
    const { loading, response } = useWhoami();

    useEffect(() => {
        if (response && config.ENVIROMENT === "PROD") {
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
            {window.alert(process.env.REACT_APP_KRATOS_URL)}
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
