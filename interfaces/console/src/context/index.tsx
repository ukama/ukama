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
  selectedDefaultSite: '',
  setSelectedDefaultSite: (siteId: string) => {},
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
  subscriptionClient: undefined,
  setSubscriptionClient: (client: any) => {},
  metaInfo: {
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
  },
  setMetaInfo: (info: any) => {},
});

const INIT_META_INFO = {
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
  const JInfo = localStorage.getItem('metaInfo');
  const info = JSON.parse(JInfo || '{}');
  const [subscriptionClient, setSubscriptionClient] = useState<any>(
    getMetricsClient(initEnv.METRIC_URL),
  );
  const [env, setEnv] = useState<TEnv>(initEnv);
  const [token, setToken] = useState(_token);
  const [pageName, setPageName] = useState('Home');
  const [isDarkMode, setIsDarkMode] = useState(false);
  const [metaInfo, setMetaInfo] = useState(
    info.ip
      ? {
          ip: info.ip,
          city: info.city,
          lat: info.lat,
          lng: info.lng,
          languages: info.languages,
          currency: info.currency,
          timezone: info.timezone,
          region_code: info.region_code,
          country_code: info.country_code,
          country_name: info.country_name,
          country_calling_code: info.country_calling_code,
        }
      : INIT_META_INFO,
  );
  const [skeltonLoading, setSkeltonLoading] = useState(false);
  const [isValidSession, setIsValidSession] = useState(false);
  const [selectedDefaultSite, setSelectedDefaultSite] = useState('');
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
      selectedDefaultSite,
      setSelectedDefaultSite,
      subscriptionClient,
      setSubscriptionClient,
      metaInfo,
      setMetaInfo,
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
      selectedDefaultSite,
      setSelectedDefaultSite,
      subscriptionClient,
      setSubscriptionClient,
      metaInfo,
      setMetaInfo,
    ],
  );

  return <AppContext.Provider value={value}>{children}</AppContext.Provider>;
};

export default AppContextWrapper;

export function useAppContext() {
  return useContext(AppContext);
}
