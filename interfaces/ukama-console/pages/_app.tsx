import {
  isDarkmode,
  isSkeltonLoading,
  pageName,
  snackbarMessage,
  user,
} from '@/app-recoil';
import client from '@/client/ApolloClient';
import { theme } from '@/styles/theme';
import Layout from '@/ui/layout';
import { doesHttpOnlyCookieExist, getTitleFromPath } from '@/utils';
import { ApolloProvider } from '@apollo/client';
import { Alert, AlertColor, Snackbar } from '@mui/material';
import { ThemeProvider } from '@mui/material/styles';
import type { AppProps } from 'next/app';
import { useEffect } from 'react';
import {
  RecoilRoot,
  useRecoilState,
  useRecoilValue,
  useResetRecoilState,
  useSetRecoilState,
} from 'recoil';
import '../styles/global.css';

const SNACKBAR_TIMEOUT = 5000;

const App = ({ Component, pageProps }: AppProps) => {
  const [_user, _setUser] = useRecoilState<any>(user);
  const setPage = useSetRecoilState(pageName);
  const _isDarkMod = useRecoilValue<boolean>(isDarkmode);
  const [_snackbarMessage, setSnackbarMessage] =
    useRecoilState<any>(snackbarMessage);
  const resetData = useResetRecoilState(user);
  const resetPageName = useResetRecoilState(pageName);
  const setSkeltonLoading = useSetRecoilState(isSkeltonLoading);
  useEffect(() => {
    if (typeof window !== 'undefined') {
      const id = new URLSearchParams(window.location.search).get('id');
      const name = new URLSearchParams(window.location.search).get('name');
      const email = new URLSearchParams(window.location.search).get('email');

      if (id && name && email) {
        _setUser({ id, name, email });
        window.history.pushState(null, '', '/home');
      }
      if ((id && name && email) || (_user.id && _user.name && _user.email)) {
        setPage(getTitleFromPath(window.location.pathname));

        if (
          !doesHttpOnlyCookieExist('id') &&
          doesHttpOnlyCookieExist('ukama_session')
        ) {
          resetData();
          resetPageName();
          window.location.replace(`${process.env.REACT_APP_AUTH_URL}/logout`);
        } else if (
          doesHttpOnlyCookieExist('id') &&
          !doesHttpOnlyCookieExist('ukama_session')
        )
          handleGoToLogin();
      } else {
        if (process.env.NODE_ENV === 'test') return;
        handleGoToLogin();
      }
    }

    setSkeltonLoading(false);
  }, []);
  const handleGoToLogin = () => {
    setPage('Home');
    typeof window !== 'undefined' &&
      window.location.replace(process.env.REACT_APP_AUTH_URL || '');
  };

  const handleSnackbarClose = () =>
    setSnackbarMessage({ ..._snackbarMessage, show: false });

  return (
    <ApolloProvider client={client}>
      <ThemeProvider theme={theme(_isDarkMod)}>
        <Layout>
          <Component {...pageProps} />
        </Layout>
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

const MyApp = (appProps: AppProps) => {
  return (
    <RecoilRoot>
      <App {...appProps} />
    </RecoilRoot>
  );
};

export default MyApp;
