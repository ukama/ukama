import Layout from "./layout";
import { theme } from "./theme";
import { useEffect } from "react";
import Router from "./router/Router";
import client from "./api/ApolloClient";
import { routes } from "./router/config";
import { CenterContainer } from "./styles";
import useWhoami from "./helpers/useWhoami";
import { ThemeProvider } from "@emotion/react";
import { ApolloProvider } from "@apollo/client";
import { BrowserRouter } from "react-router-dom";
import { CircularProgress, CssBaseline } from "@mui/material";

const App = () => {
    const { loading, response } = useWhoami();

    useEffect(() => {
        if (response && process.env.NODE_ENV === "production") {
            if (!response?.isValid) {
                window.close();
                window.location.replace(`${process.env.REACT_APP_AUTH_URL}`);
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
