/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/**
 * Client-side auth context. The user is resolved on the server (proxy +
 * getCurrentUser) and handed to AuthProvider, so client components can read
 * it synchronously via useAuth() without an extra round-trip.
 */
import { createContext, useContext } from 'react';
import type { AuthUser } from './types';

const AuthContext = createContext<AuthUser | null>(null);

export function AuthProvider({
  user,
  children,
}: {
  user: AuthUser | null;
  children: React.ReactNode;
}) {
  return <AuthContext.Provider value={user}>{children}</AuthContext.Provider>;
}

/** Returns the authenticated user, or null when unauthenticated. */
export function useAuth(): AuthUser | null {
  return useContext(AuthContext);
}

/** Like useAuth, but throws if no user — for screens that require auth. */
export function useRequiredAuth(): AuthUser {
  const user = useContext(AuthContext);
  if (!user) {
    throw new Error('useRequiredAuth called without an authenticated user');
  }
  return user;
}
