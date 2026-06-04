/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

import { useEffect } from 'react';
import CssBaseline from '@mui/material/CssBaseline';
import ApolloClientProvider from '@/client/ApolloClientProvider';
import { ThemeProvider } from '@mui/material/styles';
import ToastProvider from '@/components/ToastProvider';
import { AuthProvider } from '@/lib/auth/context';
import type { AuthUser } from '@/lib/auth/types';
import { useUiPrefs } from '@/lib/store';
import { theme } from '@/theme/theme';

/** Applies accent/density data-attributes to <html> after mount. */
function ThemeAttributes() {
  const accent = useUiPrefs((s) => s.accent);
  const density = useUiPrefs((s) => s.density);

  useEffect(() => {
    const el = document.documentElement;
    el.setAttribute('data-accent', accent);
    el.setAttribute('data-density', density);
  }, [accent, density]);

  return null;
}

export default function Providers({
  user,
  children,
}: {
  user: AuthUser | null;
  children: React.ReactNode;
}) {
  return (
    <ThemeProvider theme={theme} defaultMode="light">
      <CssBaseline enableColorScheme />
      <ThemeAttributes />
      <AuthProvider user={user}>
        <ApolloClientProvider>
          <ToastProvider>{children}</ToastProvider>
        </ApolloClientProvider>
      </AuthProvider>
    </ThemeProvider>
  );
}
