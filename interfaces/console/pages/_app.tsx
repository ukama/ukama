/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

'use client';
import client from '@/client/ApolloClient';
import { IPFY_URL, IP_API_BASE_URL } from '@/constants';
import AppContextWrapper, { useAppContext } from '@/context';
import { theme } from '@/styles/theme';
import { MyAppProps } from '@/types';
import AuthWrapper from '@/ui/wrappers/authWrapper';
import createEmotionCache from '@/ui/wrappers/createEmotionCache';
import ErrorBoundary from '@/ui/wrappers/errorBoundary';
import { ApolloProvider, HttpLink } from '@apollo/client';
import { CacheProvider, ThemeProvider } from '@emotion/react';
import { Alert, AlertColor, CssBaseline, Snackbar } from '@mui/material';
import dynamic from 'next/dynamic';
import { useEffect } from 'react';
import '../styles/global.css';
const MainApp = dynamic(() => import('@/pages/_main_app'));
const clientSideEmotionCache = createEmotionCache();
const SNACKBAR_TIMEOUT = 5000;

const getMetaInfo = async () => {
  return await fetch(IPFY_URL, {
    method: 'GET',
  })
    .then((response) => response.text())
    .then((data) => JSON.parse(data))
    .then((data) =>
      fetch(`${IP_API_BASE_URL}/${data.ip}/json/`, {
        method: 'GET',
      }),
    )
    .then((response) => response.text())
    .then((data) => JSON.parse(data))
    .catch((err) => {
      return {};
    });
};

const ClientWrapper = (appProps: MyAppProps) => {
  const {
    token,
    isDarkMode,
    snackbarMessage,
    setSnackbarMessage,
    skeltonLoading,
    setSkeltonLoading,
    isValidSession,
  } = useAppContext();

  useEffect(() => {
    if (isValidSession) {
      const httpLink = new HttpLink({
        uri: `${process.env.NEXT_PUBLIC_API_GW}/graphql`,
        credentials: 'include',
        headers: {
          token: token,
        },
      });
      client.setLink(httpLink);
    }
  }, [isValidSession]);

  const handleSnackbarClose = () =>
    setSnackbarMessage({ ...snackbarMessage, show: false });

  return (
    <ApolloProvider client={client}>
      <AuthWrapper>
        <CacheProvider value={clientSideEmotionCache}>
          <ThemeProvider theme={theme(isDarkMode)}>
            <CssBaseline />
            <MainApp {...appProps} />
            <Snackbar
              open={snackbarMessage.show}
              autoHideDuration={SNACKBAR_TIMEOUT}
              onClose={handleSnackbarClose}
            >
              <Alert
                id={snackbarMessage.id}
                severity={snackbarMessage.type as AlertColor}
                onClose={handleSnackbarClose}
              >
                {snackbarMessage.message}
              </Alert>
            </Snackbar>
          </ThemeProvider>
        </CacheProvider>
      </AuthWrapper>
    </ApolloProvider>
  );
};

const RootWrapper = (appProps: MyAppProps) => {
  return (
    <ErrorBoundary>
      <AppContextWrapper>
        <ClientWrapper {...appProps} />
      </AppContextWrapper>
    </ErrorBoundary>
  );
};

export default RootWrapper;
