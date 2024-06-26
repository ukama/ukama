/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';

import { TNetwork, TSnackbarMessage, TUser } from '@/types';
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
  children,
  token: _token,
  initalUserValues,
}: {
  token: string;
  initalUserValues: TUser;
  children: React.ReactNode;
}) => {
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
