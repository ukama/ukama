/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';

import { TEnv, TUser } from '@/types';
import React from 'react';
import { EnvContextProvider, useEnvContext } from './envContext';
import { NetworkContextProvider, useNetworkContext } from './networkContext';
import { UIContextProvider, useUIContext } from './uiContext';
import { UserContextProvider, useUserContext } from './userContext';

export { useEnvContext } from './envContext';
export { useNetworkContext } from './networkContext';
export { useUIContext } from './uiContext';
export { useUserContext } from './userContext';

interface AppContextWrapperProps {
  initEnv: TEnv;
  token: string;
  initalUserValues: TUser;
  children: React.ReactNode;
}

const AppContextWrapper: React.FC<AppContextWrapperProps> = ({
  initEnv,
  token,
  initalUserValues,
  children,
}) => {
  return (
    <EnvContextProvider initEnv={initEnv}>
      <UserContextProvider initialToken={token} initialUser={initalUserValues}>
        <NetworkContextProvider>
          <UIContextProvider>{children}</UIContextProvider>
        </NetworkContextProvider>
      </UserContextProvider>
    </EnvContextProvider>
  );
};

/** @deprecated Use the specific domain hooks: useUserContext, useNetworkContext, useUIContext, useEnvContext */
export function useAppContext() {
  const user = useUserContext();
  const network = useNetworkContext();
  const ui = useUIContext();
  const env = useEnvContext();

  return {
    ...user,
    ...network,
    ...ui,
    ...env,
  };
}

export default AppContextWrapper;
