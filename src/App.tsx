import { theme } from "./theme";
import Router from "./router/Router";
import client from "./api/ApolloClient";
import { routes } from "./router/config";
import { CenterContainer } from "./styles";
import { BasicDialog } from "./components";
import { useEffect, useState } from "react";
import useWhoami from "./helpers/useWhoami";
import { ThemeProvider } from "@emotion/react";
import { ApolloProvider } from "@apollo/client";
import { BrowserRouter } from "react-router-dom";
import { isDarkmode, isFirstVisit, isSkeltonLoading, pageName } from "./recoil";
import { useRecoilState, useRecoilValue, useSetRecoilState } from "recoil";
import { CircularProgress, CssBaseline } from "@mui/material";

const App = () => {
    const { loading, response } = useWhoami();
    const setPage = useSetRecoilState(pageName);
    const _isDarkMod = useRecoilValue(isDarkmode);
    const setSkeltonLoading = useSetRecoilState(isSkeltonLoading);
    const [_isFirstVisit, _setIsFirstVisit] = useRecoilState(isFirstVisit);
    const [showValidationError, setShowValidationError] =
        useState<boolean>(false);

    useEffect(() => {
        setSkeltonLoading(true);
    }, []);

    useEffect(() => {
        if (response) {
            if (!response?.isValid) {
                setPage("Home");
                if (_isFirstVisit) {
                    handleGoToLogin();
                } else {
                    setSkeltonLoading(true);
                    setShowValidationError(true);
                }
            } else if (response?.isValid) {
                if (_isFirstVisit) {
                    _setIsFirstVisit(false);
                }
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
            <ThemeProvider theme={theme(_isDarkMod)}>
                <CssBaseline />
                <BrowserRouter>
                    <Router routes={Object.values(routes)} />
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
