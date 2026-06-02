/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

import { getMetricsClient } from '@/client/client';
import { TEnv } from '@/types';
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

interface EnvContextState {
  env: TEnv;
  metaInfo: MetaInfo;
  subscriptionClient: ReturnType<typeof getMetricsClient>;
}

interface EnvContextActions {
  setEnv: (env: TEnv) => void;
  setMetaInfo: (info: MetaInfo) => void;
  setSubscriptionClient: (
    client: ReturnType<typeof getMetricsClient>,
  ) => void;
}

export type EnvContextType = EnvContextState & EnvContextActions;

const EnvContext = createContext<EnvContextType | undefined>(undefined);

interface EnvContextProviderProps {
  initEnv: TEnv;
  children: React.ReactNode;
}

export const EnvContextProvider: React.FC<EnvContextProviderProps> = ({
  initEnv,
  children,
}) => {
  const initialMetaInfo = useMemo(() => {
    if (typeof window === 'undefined') return defaultMetaInfo;
    const storedInfo = localStorage.getItem('metaInfo');
    if (!storedInfo) return defaultMetaInfo;
    try {
      const parsed = JSON.parse(storedInfo);
      // Validate required shape — stale/schema-drifted values fall back to default
      if (
        parsed &&
        typeof parsed === 'object' &&
        typeof parsed.ip === 'string' &&
        typeof parsed.lat === 'number' &&
        typeof parsed.lng === 'number' &&
        typeof parsed.country_code === 'string'
      ) {
        return { ...defaultMetaInfo, ...parsed } as MetaInfo;
      }
    } catch {
      // corrupted entry — clear it so it doesn't keep failing
      localStorage.removeItem('metaInfo');
    }
    return defaultMetaInfo;
  }, []);

  const [env, setEnv] = useState<TEnv>(initEnv);
  const [metaInfo, setMetaInfo] = useState<MetaInfo>(initialMetaInfo);
  const [subscriptionClient, setSubscriptionClient] = useState(() =>
    getMetricsClient(initEnv.METRIC_URL),
  );

  const value = useMemo(
    () => ({
      env,
      metaInfo,
      subscriptionClient,
      setEnv,
      setMetaInfo,
      setSubscriptionClient,
    }),
    [env, metaInfo, subscriptionClient],
  );

  return <EnvContext.Provider value={value}>{children}</EnvContext.Provider>;
};

export function useEnvContext() {
  const context = useContext(EnvContext);
  if (context === undefined) {
    throw new Error('useEnvContext must be used within an EnvContextProvider');
  }
  return context;
}

export type { MetaInfo };
