/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

'use client';
import {
  commonData,
  isDarkmode,
  pageName,
  snackbarMessage,
  user,
} from '@/app-recoil';
import client from '@/client/ApolloClient';
import { theme } from '@/styles/theme';
import { MyAppProps, TCommonData, TSnackMessage } from '@/types';
import createEmotionCache from '@/ui/wrappers/createEmotionCache';
import ErrorBoundary from '@/ui/wrappers/errorBoundary';
import { doesHttpOnlyCookieExist } from '@/utils';
import { ApolloProvider, HttpLink } from '@apollo/client';
import { CacheProvider } from '@emotion/react';
import { Alert, AlertColor, CssBaseline, Snackbar } from '@mui/material';
import { ThemeProvider } from '@mui/material/styles';
import dynamic from 'next/dynamic';
import { useEffect, useState } from 'react';
import {
  RecoilRoot,
  useRecoilState,
  useRecoilValue,
  useResetRecoilState,
} from 'recoil';
import '../styles/global.css';
const MainApp = dynamic(() => import('@/pages/_main_app'));
const clientSideEmotionCache = createEmotionCache();
const SNACKBAR_TIMEOUT = 5000;

const ClientWrapper = (appProps: MyAppProps) => {
  const _isDarkMod = useRecoilValue<boolean>(isDarkmode);
  const _commonData = useRecoilValue<TCommonData>(commonData);
  const [_snackbarMessage, setSnackbarMessage] =
    useRecoilState<TSnackMessage>(snackbarMessage);
  const [mounted, setMounted] = useState(false);
  const resetData = useResetRecoilState(user);
  const resetPageName = useResetRecoilState(pageName);
  const httpLink = new HttpLink({
    uri: process.env.NEXT_PUBLIC_API_GW,
    credentials: 'include',
    headers: {
      'org-id': _commonData.orgId,
      'user-id': _commonData.userId,
      'org-name': _commonData.orgName,
      'x-session-token': 'abc',
    },
  });

  useEffect(() => {
    setMounted(true);
    return () => {
      setMounted(false);
    };
  });

  const getClient = (): any => {
    if (mounted && !doesHttpOnlyCookieExist('ukama_session')) handleGoToLogin();
    client.setLink(httpLink);
    return client;
  };

  const handleGoToLogin = () => {
    resetData();
    resetPageName();
    typeof window !== 'undefined' &&
      window.location.replace(process.env.NEXT_PUBLIC_AUTH_APP_URL || '');
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
