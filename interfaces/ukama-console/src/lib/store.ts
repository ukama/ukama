/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/**
 * UI preference store (client-only): accent + density. Color mode itself is
 * handled by MUI's useColorScheme. Persisted to localStorage; the attributes
 * are applied to <html> by <ThemeAttributes/> after mount (no SSR mismatch).
 */
import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import type { Accent, Density } from '@/theme/tokens';

export type Rail = 'full' | 'icon';

interface UiPrefsState {
  accent: Accent;
  density: Density;
  rail: Rail;
  networkId: string;
  /** Last /configure URL (path+query) — onboarding resume point. */
  lastConfigureUrl: string | null;
  setAccent: (accent: Accent) => void;
  setDensity: (density: Density) => void;
  toggleRail: () => void;
  setNetworkId: (networkId: string) => void;
  setLastConfigureUrl: (url: string | null) => void;
}

export const useUiPrefs = create<UiPrefsState>()(
  persist(
    (set) => ({
      accent: 'blue',
      density: 'comfortable',
      rail: 'full',
      networkId: 'kwacha',
      lastConfigureUrl: null,
      setAccent: (accent) => set({ accent }),
      setDensity: (density) => set({ density }),
      toggleRail: () =>
        set((s) => ({ rail: s.rail === 'full' ? 'icon' : 'full' })),
      setNetworkId: (networkId) => set({ networkId }),
      setLastConfigureUrl: (lastConfigureUrl) => set({ lastConfigureUrl }),
    }),
    { name: 'uk-ui-prefs' },
  ),
);
