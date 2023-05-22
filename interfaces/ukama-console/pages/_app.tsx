import {
  isDarkmode,
  isSkeltonLoading,
  pageName,
  snackbarMessage,
  user,
} from '@/app-recoil';
import { networkId } from '@/app-recoil/atom';
import client from '@/client/ApolloClient';
import { useWhoamiLazyQuery } from '@/generated';
import { theme } from '@/styles/theme';
import { TSnackMessage, TUser } from '@/types';
import createEmotionCache from '@/ui/wrappers/createEmotionCache';
import ErrorBoundary from '@/ui/wrappers/errorBoundary';
import { ApolloProvider } from '@apollo/client';
import { CacheProvider, EmotionCache } from '@emotion/react';
import { Alert, AlertColor, CssBaseline, Snackbar } from '@mui/material';
import { ThemeProvider } from '@mui/material/styles';
import type { AppProps } from 'next/app';
import dynamic from 'next/dynamic';
import { useEffect } from 'react';
import {
  RecoilRoot,
  useRecoilState,
  useRecoilValue,
  useResetRecoilState,
} from 'recoil';
import '../styles/global.css';
const Layout = dynamic(() => import('@/ui/layout'));
const SNACKBAR_TIMEOUT = 5000;

const clientSideEmotionCache = createEmotionCache();

export interface MyAppProps extends AppProps {
  emotionCache?: EmotionCache;
}

const NETWORKS = [
  {
    id: '1',
    value: "Joe's Testnet",
    label: "Joe's Testnet",
  },
  {
    id: '2',
    value: "Sam's Testnet",
    label: "Sam's Testnet",
  },
];

const App = ({
  Component,
  pageProps,
  emotionCache = clientSideEmotionCache,
}: MyAppProps) => {
  const [_user, _setUser] = useRecoilState<TUser>(user);
  const [page, setPage] = useRecoilState(pageName);
  const _isDarkMod = useRecoilValue<boolean>(isDarkmode);
  const [_snackbarMessage, setSnackbarMessage] =
    useRecoilState<TSnackMessage>(snackbarMessage);
  const [skeltonLoading, setSkeltonLoading] =
    useRecoilState<boolean>(isSkeltonLoading);
  const [network, setNetwork] = useRecoilState<string>(networkId);
  const resetData = useResetRecoilState(user);
  const resetPageName = useResetRecoilState(pageName);
  const [getWhoami, { data, loading, error }] = useWhoamiLazyQuery();

  useEffect(() => {
    // if (!_user?.id) getWhoami();
  }, []);

  useEffect(() => {
    // const { id, name, email } = _user;
    // const pathname =
    //   typeof window !== 'undefined' && window.location.pathname
    //     ? window.location.pathname
    //     : '';
    // setPage(getTitleFromPath(pathname));
    // if (id && name && email) {
    //   if (
    //     !doesHttpOnlyCookieExist('id') &&
    //     doesHttpOnlyCookieExist('ukama_session')
    //   ) {
    //     resetData();
    //     resetPageName();
    //     window.location.replace(
    //       `${process.env.NEXT_PUBLIC_REACT_AUTH_APP_URL}/logout`,
    //     );
    //   } else if (
    //     doesHttpOnlyCookieExist('id') &&
    //     !doesHttpOnlyCookieExist('ukama_session')
    //   )
    //     handleGoToLogin();
    // } else {
    //   if (process.env.NEXT_PUBLIC_NODE_ENV === 'test') return;
    //   handleGoToLogin();
    // }
  }, []);

  useEffect(() => {
    if (loading) setSkeltonLoading(true);
  }, [loading]);

  useEffect(() => {
    // if (data?.whoami) {
    //   _setUser({
    //     id: data.whoami.id,
    //     name: data.whoami.name,
    //     email: data.whoami.email,
    //     role: data.whoami.role,
    //     isFirstVisit: data.whoami.isFirstVisit,
    //   });
    //   setSkeltonLoading(false);
    // }
  }, [data]);

  useEffect(() => {
    if (error) {
      setSnackbarMessage({
        id: 'whoami-msg',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
      // resetData();
      // window.location.replace(`${process.env.NEXT_PUBLIC_REACT_AUTH_APP_URL}`);
    }
  }, [error]);

  const handleGoToLogin = () => {
    setPage('Home');
    typeof window !== 'undefined' &&
      window.location.replace(process.env.NEXT_PUBLIC_REACT_AUTH_APP_URL || '');
  };

  const handleSnackbarClose = () =>
    setSnackbarMessage({ ..._snackbarMessage, show: false });

  const handlePageChange = (page: string) => setPage(page);
  const handleNetworkChange = (id: string) => setNetwork(id);

  return (
    <CacheProvider value={emotionCache}>
      <ThemeProvider theme={theme(_isDarkMod)}>
        <CssBaseline />
        <ErrorBoundary>
          <Layout
            page={page}
            networkId={network}
            networks={NETWORKS}
            isDarkMode={_isDarkMod}
            isLoading={skeltonLoading}
            handlePageChange={handlePageChange}
            handleNetworkChange={handleNetworkChange}
          >
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
    </CacheProvider>
  );
};

const MyApp = (appProps: MyAppProps) => {
  return (
    <RecoilRoot>
      <ApolloProvider client={client}>
        <App {...appProps} />
      </ApolloProvider>
    </RecoilRoot>
  );
};

export default MyApp;
