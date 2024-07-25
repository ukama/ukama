/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';

import { TEnv, TNetwork, TSnackbarMessage, TUser } from '@/types';
import React, { createContext, useContext, useMemo, useState } from 'react';

const AppContext = createContext({
  pageName: 'Home',
  setPageName: (pageName: string) => {},
  token: '',
  setToken: (token: string) => {},
  isDarkMode: false,
  setIsDarkMode: (isDarkMode: boolean) => {},
  skeltonLoading: false,
  setSkeltonLoading: (loading: boolean) => {},
  isValidSession: false,
  setIsValidSession: (valid: boolean) => {},
  snackbarMessage: {
    id: 'message-id',
    message: '',
    type: 'info',
    show: false,
  },
  setSnackbarMessage: (s: TSnackbarMessage) => {},
  network: {
    id: '',
    name: '',
  },
  env: {
    APP_URL: '',
    SIM_TYPE: '',
    METRIC_URL: '',
    API_GW_URL: '',
    AUTH_APP_URL: '',
    MAP_BOX_TOKEN: '',
    METRIC_WEBSOCKET_URL: '',
  },
  setEnv: (e: any) => {},
  setNetwork: (n: TNetwork) => {},
  user: {
    id: '',
    name: '',
    email: '',
    role: '',
    orgId: '',
    orgName: '',
  },
  setUser: (u: TUser) => {},
});

const AppContextWrapper = ({
  initEnv,
  children,
  token: _token,
  initalUserValues,
}: {
  initEnv: TEnv;
  token: string;
  initalUserValues: TUser;
  children: React.ReactNode;
}) => {
  const [env, setEnv] = useState<TEnv>(initEnv);
  const [token, setToken] = useState(_token);
  const [pageName, setPageName] = useState('Home');
  const [isDarkMode, setIsDarkMode] = useState(false);
  const [skeltonLoading, setSkeltonLoading] = useState(false);
  const [isValidSession, setIsValidSession] = useState(false);
  const [snackbarMessage, setSnackbarMessage] = useState<TSnackbarMessage>({
    id: 'message-id',
    message: '',
    type: 'info',
    show: false,
  });
  const [network, setNetwork] = useState<TNetwork>({
    id: '',
    name: '',
  });
  const [user, setUser] = useState<TUser>(initalUserValues);

  const value = useMemo(
    () => ({
      env,
      setEnv,
      isDarkMode,
      setIsDarkMode,
      user,
      setUser,
      token,
      setToken,
      network,
      setNetwork,
      pageName,
      setPageName,
      skeltonLoading,
      setSkeltonLoading,
      isValidSession,
      setIsValidSession,
      snackbarMessage,
      setSnackbarMessage,
    }),
    [
      env,
      setEnv,
      isDarkMode,
      setIsDarkMode,
      user,
      setUser,
      token,
      setToken,
      network,
      setNetwork,
      pageName,
      setPageName,
      skeltonLoading,
      setSkeltonLoading,
      isValidSession,
      setIsValidSession,
      snackbarMessage,
      setSnackbarMessage,
    ],
  );

  return <AppContext.Provider value={value}>{children}</AppContext.Provider>;
};

export default AppContextWrapper;

export function useAppContext() {
  return useContext(AppContext);
}
