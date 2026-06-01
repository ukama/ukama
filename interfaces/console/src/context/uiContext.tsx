/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';

import { TSnackbarMessage } from '@/types';
import React, { createContext, useContext, useMemo, useState } from 'react';

const defaultSnackbarMessage: TSnackbarMessage = {
  id: 'message-id',
  message: '',
  type: 'info',
  show: false,
};

interface UIContextState {
  pageName: string;
  isDarkMode: boolean;
  skeltonLoading: boolean;
  snackbarMessage: TSnackbarMessage;
}

interface UIContextActions {
  setPageName: (pageName: string) => void;
  setIsDarkMode: (isDarkMode: boolean) => void;
  setSkeltonLoading: (loading: boolean) => void;
  setSnackbarMessage: (message: TSnackbarMessage) => void;
}

export type UIContextType = UIContextState & UIContextActions;

const UIContext = createContext<UIContextType | undefined>(undefined);

export const UIContextProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [pageName, setPageName] = useState('Home');
  const [isDarkMode, setIsDarkMode] = useState(false);
  const [skeltonLoading, setSkeltonLoading] = useState(false);
  const [snackbarMessage, setSnackbarMessage] = useState<TSnackbarMessage>(
    defaultSnackbarMessage,
  );

  const value = useMemo(
    () => ({
      pageName,
      isDarkMode,
      skeltonLoading,
      snackbarMessage,
      setPageName,
      setIsDarkMode,
      setSkeltonLoading,
      setSnackbarMessage,
    }),
    [pageName, isDarkMode, skeltonLoading, snackbarMessage],
  );

  return <UIContext.Provider value={value}>{children}</UIContext.Provider>;
};

export function useUIContext() {
  const context = useContext(UIContext);
  if (context === undefined) {
    throw new Error('useUIContext must be used within a UIContextProvider');
  }
  return context;
}
