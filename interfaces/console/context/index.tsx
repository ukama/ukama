'use client';

import React, { createContext, useContext, useState } from 'react';

type TNetwork = {
  id: string;
  name: string;
};

type TUser = {
  id: string;
  name: string;
  email: string;
  role: string;
  orgId: string;
  orgName: string;
};

type TSnackbarMessage = {
  id: string;
  message: string;
  type: string;
  show: boolean;
};

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

const AppContextWrapper = ({ children }: { children: React.ReactNode }) => {
  const [pageName, setPageName] = useState('Home');
  const [token, setToken] = useState('Home');
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
  const [user, setUser] = useState<TUser>({
    id: '',
    name: '',
    role: '',
    email: '',
    orgId: '',
    orgName: '',
  });

  const value = {
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
  };

  return <AppContext.Provider value={value}>{children}</AppContext.Provider>;
};

export default AppContextWrapper;

export function useAppContext() {
  return useContext(AppContext);
}
