/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/**
 * App toast — dark snackbar with check icon and optional action (ds.jsx
 * useToast; design finding #11: routine changes get an Undo snackbar).
 */
import { createContext, useCallback, useContext, useMemo, useState } from 'react';
import Snackbar from '@mui/material/Snackbar';
import CheckCircleRounded from '@mui/icons-material/CheckCircleRounded';

export interface ToastOptions {
  action?: { label: string; fn: () => void };
  ms?: number;
}

type ShowToast = (msg: string, opts?: ToastOptions) => void;

const ToastContext = createContext<ShowToast>(() => undefined);

export function useToast(): ShowToast {
  return useContext(ToastContext);
}

interface ToastState {
  msg: string;
  action?: { label: string; fn: () => void };
  ms: number;
  key: number;
}

export default function ToastProvider({ children }: { children: React.ReactNode }) {
  const [toast, setToast] = useState<ToastState | null>(null);

  const show = useCallback<ShowToast>((msg, opts) => {
    setToast({
      msg,
      ms: opts?.ms ?? 3200,
      key: Date.now(),
      ...(opts?.action ? { action: opts.action } : {}),
    });
  }, []);

  const value = useMemo(() => show, [show]);

  return (
    <ToastContext.Provider value={value}>
      {children}
      <Snackbar
        key={toast?.key}
        open={!!toast}
        autoHideDuration={toast?.ms}
        onClose={(_, reason) => {
          if (reason !== 'clickaway') setToast(null);
        }}
        anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
        slotProps={{
          content: {
            sx: {
              bgcolor: '#15181F',
              color: '#fff',
              borderRadius: 2.5,
              fontSize: 13.5,
              boxShadow: 'var(--uk-shadow-lg)',
            },
          },
        }}
        message={
          toast ? (
            <span style={{ display: 'flex', alignItems: 'center', gap: 14 }}>
              <CheckCircleRounded
                sx={{ fontSize: 19, color: 'var(--uk-success-bright)' }}
              />
              {toast.msg}
            </span>
          ) : undefined
        }
        action={
          toast?.action ? (
            <button
              type="button"
              onClick={() => {
                toast.action?.fn();
                setToast(null);
              }}
              style={{
                background: 'none',
                border: 'none',
                color: 'var(--uk-ac-light)',
                fontWeight: 600,
                cursor: 'pointer',
                fontSize: 13.5,
                fontFamily: 'inherit',
              }}
            >
              {toast.action.label}
            </button>
          ) : undefined
        }
      />
    </ToastContext.Provider>
  );
}
