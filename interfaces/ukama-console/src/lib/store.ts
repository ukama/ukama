/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/**
 * UI preference store (client-only): accent + density + map view. Color mode
 * itself is handled by MUI's useColorScheme. Persisted to localStorage; the
 * attributes are applied to <html> by <ThemeAttributes/> after mount (no SSR
 * mismatch).
 */
import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { DEFAULT_MAP_VIEW, type MapView } from '@/components/Map/basemaps';
import type { Accent, Density } from '@/theme/tokens';

export type Rail = 'full' | 'icon';

interface UiPrefsState {
  accent: Accent;
  density: Density;
  rail: Rail;
  /** Base layer used by every map (street / satellite / terrain). */
  mapView: MapView;
  networkId: string;
  /** Last /configure URL (path+query) — onboarding resume point. */
  lastConfigureUrl: string | null;
  setAccent: (accent: Accent) => void;
  setDensity: (density: Density) => void;
  toggleRail: () => void;
  setMapView: (mapView: MapView) => void;
  setNetworkId: (networkId: string) => void;
  setLastConfigureUrl: (url: string | null) => void;
}

export const useUiPrefs = create<UiPrefsState>()(
  persist(
    (set) => ({
      accent: 'blue',
      density: 'compact',
      rail: 'full',
      mapView: DEFAULT_MAP_VIEW,
      networkId: 'kwacha',
      lastConfigureUrl: null,
      setAccent: (accent) => set({ accent }),
      setDensity: (density) => set({ density }),
      toggleRail: () =>
        set((s) => ({ rail: s.rail === 'full' ? 'icon' : 'full' })),
      setMapView: (mapView) => set({ mapView }),
      setNetworkId: (networkId) => set({ networkId }),
      setLastConfigureUrl: (lastConfigureUrl) => set({ lastConfigureUrl }),
    }),
    { name: 'uk-ui-prefs' },
  ),
);
