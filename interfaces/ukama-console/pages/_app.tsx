import {
  isDarkmode,
  isSkeltonLoading,
  pageName,
  snackbarMessage,
  user,
} from '@/app-recoil';
import client from '@/client/ApolloClient';
import { useWhoamiLazyQuery } from '@/generated';
import { theme } from '@/styles/theme';
import { TSnackMessage, TUser } from '@/types';
import Layout from '@/ui/layout';
import ErrorBoundary from '@/ui/wrappers/errorBoundary';
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
  const [_user, _setUser] = useRecoilState<TUser>(user);
  const setPage = useSetRecoilState(pageName);
  const _isDarkMod = useRecoilValue<boolean>(isDarkmode);
  const [_snackbarMessage, setSnackbarMessage] =
    useRecoilState<TSnackMessage>(snackbarMessage);
  const resetData = useResetRecoilState(user);
  const resetPageName = useResetRecoilState(pageName);
  const setSkeltonLoading = useSetRecoilState(isSkeltonLoading);
  const [getWhoami, { data, loading, error }] = useWhoamiLazyQuery();

  useEffect(() => {
    if (!_user?.id) getWhoami();
  }, []);

  useEffect(() => {
    const { id, name, email } = _user;
    if (id && name && email) {
      if (
        !doesHttpOnlyCookieExist('id') &&
        doesHttpOnlyCookieExist('ukama_session')
      ) {
        resetData();
        resetPageName();
        window.location.replace(
          `${process.env.NEXT_PUBLIC_REACT_AUTH_APP_URL}/logout`,
        );
      } else if (
        doesHttpOnlyCookieExist('id') &&
        !doesHttpOnlyCookieExist('ukama_session')
      )
        handleGoToLogin();
    } else {
      if (process.env.NEXT_PUBLIC_NODE_ENV === 'test') return;
      handleGoToLogin();
    }
  }, []);

  useEffect(() => {
    if (loading) setSkeltonLoading(true);
  }, [loading]);

  useEffect(() => {
    if (data?.whoami) {
      setPage(getTitleFromPath(window.location.pathname));
      _setUser({
        id: data.whoami.id,
        name: data.whoami.name,
        email: data.whoami.email,
        role: data.whoami.role,
        isFirstVisit: data.whoami.isFirstVisit,
      });
      setSkeltonLoading(false);
    }
  }, [data]);

  useEffect(() => {
    if (error) {
      setSnackbarMessage({
        id: 'whoami-msg',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
      resetData();
      window.location.replace(`${process.env.NEXT_PUBLIC_REACT_AUTH_APP_URL}`);
    }
  }, [error]);

  const handleGoToLogin = () => {
    setPage('Home');
    typeof window !== 'undefined' &&
      window.location.replace(process.env.NEXT_PUBLIC_REACT_AUTH_APP_URL || '');
  };

  const handleSnackbarClose = () =>
    setSnackbarMessage({ ..._snackbarMessage, show: false });

  return (
    <ThemeProvider theme={theme(_isDarkMod)}>
      <ErrorBoundary>
        <Layout>
          <Component {...pageProps} />
        </Layout>
      </ErrorBoundary>
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
  );
};

const MyApp = (appProps: AppProps) => {
  return (
    <RecoilRoot>
      <ApolloProvider client={client}>
        <App {...appProps} />
      </ApolloProvider>
    </RecoilRoot>
  );
};

export default MyApp;
