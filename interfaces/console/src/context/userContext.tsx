/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

import { TUser } from '@/types';
import React, { createContext, useContext, useMemo, useState } from 'react';

interface UserContextState {
  user: TUser;
  token: string;
  isValidSession: boolean;
}

interface UserContextActions {
  setUser: (user: TUser) => void;
  setToken: (token: string) => void;
  setIsValidSession: (valid: boolean) => void;
}

export type UserContextType = UserContextState & UserContextActions;

const UserContext = createContext<UserContextType | undefined>(undefined);

interface UserContextProviderProps {
  initialToken: string;
  initialUser: TUser;
  children: React.ReactNode;
}

export const UserContextProvider: React.FC<UserContextProviderProps> = ({
  initialToken,
  initialUser,
  children,
}) => {
  const [user, setUser] = useState<TUser>(initialUser);
  const [token, setToken] = useState(initialToken);
  const [isValidSession, setIsValidSession] = useState(false);

  const value = useMemo(
    () => ({ user, token, isValidSession, setUser, setToken, setIsValidSession }),
    [user, token, isValidSession],
  );

  return <UserContext.Provider value={value}>{children}</UserContext.Provider>;
};

export function useUserContext() {
  const context = useContext(UserContext);
  if (context === undefined) {
    throw new Error('useUserContext must be used within a UserContextProvider');
  }
  return context;
}
