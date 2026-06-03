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

interface UiPrefsState {
  accent: Accent;
  density: Density;
  setAccent: (accent: Accent) => void;
  setDensity: (density: Density) => void;
}

export const useUiPrefs = create<UiPrefsState>()(
  persist(
    (set) => ({
      accent: 'blue',
      density: 'comfortable',
      setAccent: (accent) => set({ accent }),
      setDensity: (density) => set({ density }),
    }),
    { name: 'uk-ui-prefs' },
  ),
);
