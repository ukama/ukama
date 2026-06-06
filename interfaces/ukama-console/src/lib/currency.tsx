/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Org currency, resolved once. The currency *code* comes from the auth token
 * (org.currency); its display symbol is looked up via the BFF
 * getCurrencySymbol exactly once per session and shared through context, so
 * no component fetches it directly. Apollo caches the result for 24h on top
 * of that (static reference data), making remounts free too.
 *
 * Use `useCurrency().symbol` for the symbol, or `useCurrency().format(n)` to
 * render an amount.
 */
'use client';

import { createContext, useContext, useMemo } from 'react';

import { useGetCurrencySymbolQuery } from '@/client/graphql/currency.generated';
import { useAuth } from '@/lib/auth/context';

interface CurrencyContextValue {
  /** ISO currency code (e.g. USD), from the org. */
  code: string;
  /** Display symbol (e.g. $). Falls back to the code, then '$'. */
  symbol: string;
  /** Formats an amount with the resolved symbol, e.g. "$1,200". */
  format: (amount: number) => string;
}

const CurrencyContext = createContext<CurrencyContextValue | null>(null);

export function CurrencyProvider({ children }: { children: React.ReactNode }) {
  const user = useAuth();
  const code = user?.currency ?? '';

  const { data } = useGetCurrencySymbolQuery({
    variables: { code },
    skip: !code,
    // Static reference data — keep it cached, never refetch on mount.
    fetchPolicy: 'cache-first',
  });

  const value = useMemo<CurrencyContextValue>(() => {
    const symbol = data?.getCurrencySymbol.symbol || code || '$';
    return {
      code,
      symbol,
      format: (amount: number) =>
        `${symbol}${amount.toLocaleString()}`,
    };
  }, [data, code]);

  return (
    <CurrencyContext.Provider value={value}>
      {children}
    </CurrencyContext.Provider>
  );
}

export function useCurrency(): CurrencyContextValue {
  return (
    useContext(CurrencyContext) ?? {
      code: '',
      symbol: '$',
      format: (amount: number) => `$${amount.toLocaleString()}`,
    }
  );
}
