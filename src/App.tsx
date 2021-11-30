import Layout from "./layout";
import { theme } from "./theme";
import Router from "./router/Router";
import client from "./api/ApolloClient";
import { routes } from "./router/config";
import { CenterContainer } from "./styles";
import { useEffect, useState } from "react";
import useWhoami from "./helpers/useWhoami";
import { ThemeProvider } from "@emotion/react";
import { ApolloProvider } from "@apollo/client";
import { BrowserRouter } from "react-router-dom";
import { CircularProgress, CssBaseline } from "@mui/material";
import { BasicDialog } from "./components";

const App = () => {
    const { loading, response } = useWhoami();
    const [showValidationError, setShowValidationError] =
        useState<boolean>(false);

    useEffect(() => {
        if (response && !response?.isValid) {
            setShowValidationError(true);
        }
    }, [response]);

    const handleGoToLogin = () => {
        setShowValidationError(false);
        window.close();
        window.location.replace(`${process.env.REACT_APP_AUTH_URL}`);
    };

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
                    <BasicDialog
                        isClosable={false}
                        btnLabel={"Try Login"}
                        isOpen={showValidationError}
                        handleClose={handleGoToLogin}
                        title={"Session validation failed"}
                        content={
                            "Your session is not valid or expires. Please re-login."
                        }
                    />
                </BrowserRouter>
            </ThemeProvider>
        </ApolloProvider>
    );
};

export default App;
