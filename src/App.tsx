import { theme } from "./theme";
import Router from "./router/Router";
import client from "./api/ApolloClient";
import { routes } from "./router/config";
import { CenterContainer } from "./styles";
import { BasicDialog } from "./components";
import { useSetRecoilState } from "recoil";
import { useEffect, useState } from "react";
import useWhoami from "./helpers/useWhoami";
import { isSkeltonLoading } from "./recoil";
import { ThemeProvider } from "@emotion/react";
import { ApolloProvider } from "@apollo/client";
import { BrowserRouter } from "react-router-dom";
import { CircularProgress, CssBaseline } from "@mui/material";

const App = () => {
    const { loading, response } = useWhoami();
    const setSkeltonLoading = useSetRecoilState(isSkeltonLoading);
    const [showValidationError, setShowValidationError] =
        useState<boolean>(false);

    useEffect(() => {
        if (response) {
            if (!response?.isValid) {
                setSkeltonLoading(true);
                setShowValidationError(true);
            } else if (response?.isValid) {
                setSkeltonLoading(false);
            }
        }
    }, [response]);

    const handleGoToLogin = () => {
        window.location.replace(process.env.REACT_APP_AUTH_URL || "");
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
                    <Router routes={routes} />
                </BrowserRouter>
                <BasicDialog
                    isClosable={false}
                    btnLabel={"Log In"}
                    isOpen={showValidationError}
                    handleClose={handleGoToLogin}
                    title={"Session validation failed"}
                    content={
                        "Your session is not valid or has expired. Please re-login."
                    }
                />
            </ThemeProvider>
        </ApolloProvider>
    );
};

export default App;
