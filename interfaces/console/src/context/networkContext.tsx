/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

import { TNetwork } from '@/types';
import React, { createContext, useContext, useMemo, useState } from 'react';

const defaultNetwork: TNetwork = { id: '', name: '' };

interface NetworkContextState {
  network: TNetwork;
  selectedDefaultSite: string;
}

interface NetworkContextActions {
  setNetwork: (network: TNetwork) => void;
  setSelectedDefaultSite: (siteId: string) => void;
}

export type NetworkContextType = NetworkContextState & NetworkContextActions;

const NetworkContext = createContext<NetworkContextType | undefined>(undefined);

export const NetworkContextProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [network, setNetwork] = useState<TNetwork>(defaultNetwork);
  const [selectedDefaultSite, setSelectedDefaultSite] = useState('');

  const value = useMemo(
    () => ({ network, selectedDefaultSite, setNetwork, setSelectedDefaultSite }),
    [network, selectedDefaultSite],
  );

  return (
    <NetworkContext.Provider value={value}>{children}</NetworkContext.Provider>
  );
};

export function useNetworkContext() {
  const context = useContext(NetworkContext);
  if (context === undefined) {
    throw new Error(
      'useNetworkContext must be used within a NetworkContextProvider',
    );
  }
  return context;
}
