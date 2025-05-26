/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';

import { getMetricsClient } from '@/client/client';
import { TEnv, TNetwork, TSnackbarMessage, TUser } from '@/types';
import React, { createContext, useContext, useMemo, useState } from 'react';

interface MetaInfo {
  ip: string;
  city: string;
  lat: number;
  lng: number;
  languages: string;
  currency: string;
  timezone: string;
  region_code: string;
  country_code: string;
  country_name: string;
  country_calling_code: string;
}

interface AppContextState {
  pageName: string;
  token: string;
  isDarkMode: boolean;
  skeltonLoading: boolean;
  isValidSession: boolean;
  selectedDefaultSite: string;
  snackbarMessage: TSnackbarMessage;
  network: TNetwork;
  user: TUser;
  env: TEnv;
  metaInfo: MetaInfo;
  subscriptionClient: ReturnType<typeof getMetricsClient>;
}

interface AppContextActions {
  setPageName: (pageName: string) => void;
  setToken: (token: string) => void;
  setIsDarkMode: (isDarkMode: boolean) => void;
  setSkeltonLoading: (loading: boolean) => void;
  setIsValidSession: (valid: boolean) => void;
  setSelectedDefaultSite: (siteId: string) => void;
  setSnackbarMessage: (message: TSnackbarMessage) => void;
  setNetwork: (network: TNetwork) => void;
  setUser: (user: TUser) => void;
  setEnv: (env: TEnv) => void;
  setMetaInfo: (info: MetaInfo) => void;
  setSubscriptionClient: (client: ReturnType<typeof getMetricsClient>) => void;
}

type AppContextType = AppContextState & AppContextActions;

const defaultMetaInfo: MetaInfo = {
  ip: '',
  city: '',
  lat: 0,
  lng: 0,
  languages: '',
  currency: '',
  timezone: '',
  region_code: '',
  country_code: '',
  country_name: '',
  country_calling_code: '',
};

const defaultSnackbarMessage: TSnackbarMessage = {
  id: 'message-id',
  message: '',
  type: 'info',
  show: false,
};

const defaultNetwork: TNetwork = {
  id: '',
  name: '',
};

const AppContext = createContext<AppContextType | undefined>(undefined);

interface AppContextWrapperProps {
  initEnv: TEnv;
  token: string;
  initalUserValues: TUser;
  children: React.ReactNode;
}

const AppContextWrapper: React.FC<AppContextWrapperProps> = ({
  initEnv,
  token: initialToken,
  initalUserValues,
  children,
}) => {
  const initialMetaInfo = useMemo(() => {
    if (typeof window === 'undefined') return defaultMetaInfo;
    const storedInfo = localStorage.getItem('metaInfo');
    return storedInfo ? JSON.parse(storedInfo) : defaultMetaInfo;
  }, []);

  const [subscriptionClient, setSubscriptionClient] = useState(() =>
    getMetricsClient(initEnv.METRIC_URL),
  );
  const [env, setEnv] = useState<TEnv>(initEnv);
  const [token, setToken] = useState(initialToken);
  const [pageName, setPageName] = useState('Home');
  const [isDarkMode, setIsDarkMode] = useState(false);
  const [metaInfo, setMetaInfo] = useState<MetaInfo>(initialMetaInfo);
  const [skeltonLoading, setSkeltonLoading] = useState(false);
  const [isValidSession, setIsValidSession] = useState(false);
  const [selectedDefaultSite, setSelectedDefaultSite] = useState('');
  const [snackbarMessage, setSnackbarMessage] = useState<TSnackbarMessage>(
    defaultSnackbarMessage,
  );
  const [network, setNetwork] = useState<TNetwork>(defaultNetwork);
  const [user, setUser] = useState<TUser>(initalUserValues);

  const value = useMemo(
    () => ({
      env,
      isDarkMode,
      user,
      token,
      network,
      pageName,
      skeltonLoading,
      isValidSession,
      snackbarMessage,
      selectedDefaultSite,
      subscriptionClient,
      metaInfo,
      setEnv,
      setIsDarkMode,
      setUser,
      setToken,
      setNetwork,
      setPageName,
      setSkeltonLoading,
      setIsValidSession,
      setSnackbarMessage,
      setSelectedDefaultSite,
      setSubscriptionClient,
      setMetaInfo,
    }),
    [
      env,
      isDarkMode,
      user,
      token,
      network,
      pageName,
      skeltonLoading,
      isValidSession,
      snackbarMessage,
      selectedDefaultSite,
      subscriptionClient,
      metaInfo,
    ],
  );

  return <AppContext.Provider value={value}>{children}</AppContext.Provider>;
};

export function useAppContext() {
  const context = useContext(AppContext);
  if (context === undefined) {
    throw new Error('useAppContext must be used within an AppContextProvider');
  }
  return context;
}

export default AppContextWrapper;
