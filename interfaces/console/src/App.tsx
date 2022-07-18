import {
    useRecoilState,
    useRecoilValue,
    useSetRecoilState,
    useResetRecoilState,
} from "recoil";
import {
    user,
    pageName,
    isDarkmode,
    snackbarMessage,
    isSkeltonLoading,
} from "./recoil";
import { theme } from "./theme";
import { useEffect } from "react";
import Router from "./router/Router";
import client from "./api/ApolloClient";
import { routes } from "./router/config";
import { ApolloProvider } from "@apollo/client";
import { BrowserRouter } from "react-router-dom";
import { ThemeProvider } from "@mui/material/styles";
import { doesHttpOnlyCookieExist, getTitleFromPath } from "./utils";
import { Alert, AlertColor, CssBaseline, Snackbar } from "@mui/material";

const SNACKBAR_TIMEOUT = 5000;

const App = () => {
    const [_user, _setUser] = useRecoilState(user);
    const setPage = useSetRecoilState(pageName);
    const _isDarkMod = useRecoilValue(isDarkmode);
    const [_snackbarMessage, setSnackbarMessage] =
        useRecoilState(snackbarMessage);
    const resetData = useResetRecoilState(user);
    const resetPageName = useResetRecoilState(pageName);
    const setSkeltonLoading = useSetRecoilState(isSkeltonLoading);

    useEffect(() => {
        const id = new URLSearchParams(window.location.search).get("id");
        const name = new URLSearchParams(window.location.search).get("name");
        const email = new URLSearchParams(window.location.search).get("email");
        if (id && name && email) {
            _setUser({ id, name, email });
            window.history.pushState(null, "", "/");
        }
        if ((id && name && email) || (_user.id && _user.name && _user.email)) {
            setPage(getTitleFromPath(window.location.pathname));

            if (
                !doesHttpOnlyCookieExist("id") &&
                doesHttpOnlyCookieExist("ukama_session")
            ) {
                resetData();
                resetPageName();
                window.location.replace(
                    `${process.env.REACT_APP_AUTH_URL}/logout`
                );
            } else if (
                doesHttpOnlyCookieExist("id") &&
                !doesHttpOnlyCookieExist("ukama_session")
            )
                handleGoToLogin();
        } else {
            if (process.env.NODE_ENV === "test") return;
            handleGoToLogin();
        }

        setSkeltonLoading(false);
    }, []);

    const handleGoToLogin = () => {
        setPage("Home");
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
