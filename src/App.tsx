import {
    pageName,
    isDarkmode,
    isFirstVisit,
    snackbarMessage,
    isSkeltonLoading,
} from "./recoil";
import { theme } from "./theme";
import Router from "./router/Router";
import client from "./api/ApolloClient";
import { routes } from "./router/config";
import { getTitleFromPath } from "./utils";
import { useEffect } from "react";
import useWhoami from "./helpers/useWhoami";
import { ApolloProvider } from "@apollo/client";
import { BrowserRouter } from "react-router-dom";
import { ThemeProvider } from "@mui/material/styles";
import { Alert, AlertColor, CssBaseline, Snackbar } from "@mui/material";
import { useRecoilState, useRecoilValue, useSetRecoilState } from "recoil";

const SNACKBAR_TIMEOUT = 5000;

const App = () => {
    const { response } = useWhoami();
    const setPage = useSetRecoilState(pageName);
    const _isDarkMod = useRecoilValue(isDarkmode);
    const [_snackbarMessage, setSnackbarMessage] =
        useRecoilState(snackbarMessage);
    const setSkeltonLoading = useSetRecoilState(isSkeltonLoading);
    const _isFirstVisit = useRecoilValue(isFirstVisit);
    useEffect(() => {
        if (response) {
            if (!response?.isValid) {
                setPage("Home");
                if (_isFirstVisit) {
                    handleGoToLogin();
                } else {
                    handleGoToLogin();
                }
            } else if (response?.isValid) {
                setPage(getTitleFromPath(window.location.pathname));
                setSkeltonLoading(false);
            }
        }
    }, [response]);

    const handleGoToLogin = () => {
        window.location.replace(process.env.REACT_APP_AUTH_URL || "");
    };

    const handleSnackbarClose = () =>
        setSnackbarMessage({ ..._snackbarMessage, show: false });

    return (
        <ApolloProvider client={client}>
            <ThemeProvider theme={theme(_isDarkMod)}>
                <CssBaseline />
                <BrowserRouter>
                    <Router routes={Object.values(routes)} />
                </BrowserRouter>
                <Snackbar
                    open={_snackbarMessage.show}
                    autoHideDuration={SNACKBAR_TIMEOUT}
                    onClose={handleSnackbarClose}
                >
                    <Alert
                        id={_snackbarMessage.id}
                        severity={_snackbarMessage.type as AlertColor}
                        onClose={handleSnackbarClose}
                    >
                        {_snackbarMessage.message}
                    </Alert>
                </Snackbar>
            </ThemeProvider>
        </ApolloProvider>
    );
};

export default App;
