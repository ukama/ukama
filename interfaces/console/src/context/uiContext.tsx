/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

import { TSnackbarMessage } from '@/types';
import React, { createContext, useContext, useMemo, useState } from 'react';

// ── Snackbar context (isolated so toast triggers don't re-render theme/page consumers) ──

const defaultSnackbarMessage: TSnackbarMessage = {
  id: 'message-id',
  message: '',
  type: 'info',
  show: false,
};

interface SnackbarContextType {
  snackbarMessage: TSnackbarMessage;
  setSnackbarMessage: (message: TSnackbarMessage) => void;
}

const SnackbarContext = createContext<SnackbarContextType | undefined>(undefined);

const SnackbarContextProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [snackbarMessage, setSnackbarMessage] = useState<TSnackbarMessage>(defaultSnackbarMessage);
  const value = useMemo(
    () => ({ snackbarMessage, setSnackbarMessage }),
    [snackbarMessage],
  );
  return <SnackbarContext.Provider value={value}>{children}</SnackbarContext.Provider>;
};

export function useSnackbarContext() {
  const ctx = useContext(SnackbarContext);
  if (!ctx) throw new Error('useSnackbarContext must be used within UIContextProvider');
  return ctx;
}

// ── UI (theme / page / skeleton) context ──────────────────────────────────────────────

interface UIContextState {
  pageName: string;
  isDarkMode: boolean;
  skeltonLoading: boolean;
}

interface UIContextActions {
  setPageName: (pageName: string) => void;
  setIsDarkMode: (isDarkMode: boolean) => void;
  setSkeltonLoading: (loading: boolean) => void;
}

export type UIContextType = UIContextState & UIContextActions;

const UIContext = createContext<UIContextType | undefined>(undefined);

const UIInnerProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [pageName, setPageName] = useState('Home');
  const [isDarkMode, setIsDarkMode] = useState(false);
  const [skeltonLoading, setSkeltonLoading] = useState(false);

  const value = useMemo(
    () => ({ pageName, isDarkMode, skeltonLoading, setPageName, setIsDarkMode, setSkeltonLoading }),
    [pageName, isDarkMode, skeltonLoading],
  );
  return <UIContext.Provider value={value}>{children}</UIContext.Provider>;
};

export const UIContextProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => (
  <SnackbarContextProvider>
    <UIInnerProvider>{children}</UIInnerProvider>
  </SnackbarContextProvider>
);

// ── Combined hook (backwards compatible — re-exports both contexts as one object) ─────

export function useUIContext() {
  const uiCtx = useContext(UIContext);
  const snackCtx = useContext(SnackbarContext);
  if (!uiCtx || !snackCtx) {
    throw new Error('useUIContext must be used within a UIContextProvider');
  }
  return { ...uiCtx, snackbarMessage: snackCtx.snackbarMessage, setSnackbarMessage: snackCtx.setSnackbarMessage };
}
