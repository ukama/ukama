import { theme } from "./theme";
import Router from "./router/Router";
import client from "./api/ApolloClient";
import { routes } from "./router/config";
import { BasicDialog, UserActivationDialog } from "./components";
import { useEffect, useState } from "react";
import useWhoami from "./helpers/useWhoami";
import { CssBaseline } from "@mui/material";
import { ApolloProvider } from "@apollo/client";
import { BrowserRouter } from "react-router-dom";
import { ThemeProvider } from "@mui/material/styles";
import { useTranslation } from "react-i18next";
import { useRecoilState, useRecoilValue, useSetRecoilState } from "recoil";
import { isDarkmode, isFirstVisit, isSkeltonLoading, pageName } from "./recoil";
import "./i18n/i18n";

const App = () => {
    const { t } = useTranslation();
    const { response } = useWhoami();
    const [showSimActivationDialog, setShowSimActivationDialog] =
        useState(false);
    const setPage = useSetRecoilState(pageName);
    const _isDarkMod = useRecoilValue(isDarkmode);
    const setSkeltonLoading = useSetRecoilState(isSkeltonLoading);
    const [_isFirstVisit, _setIsFirstVisit] = useRecoilState(isFirstVisit);
    const [showValidationError, setShowValidationError] =
        useState<boolean>(false);
    const handleSimActivateClose = () => {
        setShowSimActivationDialog(false);
    };
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
                setShowSimActivationDialog(true);
                setSkeltonLoading(false);
            }
        }
    }, [response]);

    const handleGoToLogin = () => {
        window.location.replace(process.env.REACT_APP_AUTH_URL || "");
    };

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
                <UserActivationDialog
                    isOpen={showSimActivationDialog}
                    dialogTitle={t("DIALOG_MESSAGE.SimActivationDialogTitle")}
                    subTitle={t("DIALOG_MESSAGE.SimActivationDialogContent")}
                    handleClose={handleSimActivateClose}
                />
            </ThemeProvider>
        </ApolloProvider>
    );
};

export default App;
