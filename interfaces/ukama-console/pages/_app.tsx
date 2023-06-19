'use client';
import { commonData, isDarkmode, snackbarMessage } from '@/app-recoil';
import client from '@/client/ApolloClient';
import { theme } from '@/styles/theme';
import { MyAppProps, TCommonData, TSnackMessage } from '@/types';
import createEmotionCache from '@/ui/wrappers/createEmotionCache';
import ErrorBoundary from '@/ui/wrappers/errorBoundary';
import { ApolloProvider, HttpLink } from '@apollo/client';
import { CacheProvider } from '@emotion/react';
import { Alert, AlertColor, CssBaseline, Snackbar } from '@mui/material';
import { ThemeProvider } from '@mui/material/styles';
import dynamic from 'next/dynamic';
import { RecoilRoot, useRecoilState, useRecoilValue } from 'recoil';
import '../styles/global.css';

const MainApp = dynamic(() => import('@/pages/_main_app'));
const clientSideEmotionCache = createEmotionCache();
const SNACKBAR_TIMEOUT = 5000;

const ClientWrapper = (appProps: MyAppProps) => {
  const _isDarkMod = useRecoilValue<boolean>(isDarkmode);
  const _commonData = useRecoilValue<TCommonData>(commonData);
  const [_snackbarMessage, setSnackbarMessage] =
    useRecoilState<TSnackMessage>(snackbarMessage);
  const httpLink = new HttpLink({
    uri: process.env.NEXT_PUBLIC_REACT_APP_API,
    credentials: 'include',
    headers: {
      'org-id': _commonData.orgId,
      'user-id': _commonData.userId,
      'org-name': _commonData.orgName,
    },
  });

  const getClient = (): any => {
    client.setLink(httpLink);
    return client;
  };

  const handleSnackbarClose = () =>
    setSnackbarMessage({ ..._snackbarMessage, show: false });

  return (
    <ApolloProvider client={getClient()}>
      <CacheProvider value={clientSideEmotionCache}>
        <ThemeProvider theme={theme(_isDarkMod)}>
          <CssBaseline />
          <MainApp {...appProps} />
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
    </ApolloProvider>
  );
};

const RootWrapper = (appProps: MyAppProps) => {
  return (
    <ErrorBoundary>
      <RecoilRoot>
        <ClientWrapper {...appProps} />
      </RecoilRoot>
    </ErrorBoundary>
  );
};

export default RootWrapper;
