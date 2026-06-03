/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Design tokens — single source of truth, ported from the design package
 * (`app.css` + `assets/colors_and_type.css`). See BUILD-PLAN §7 / §7.2.
 */

export type Accent = 'blue' | 'indigo' | 'teal';
export type Density = 'comfortable' | 'compact';

export interface AccentSet {
  main: string;
  dark: string;
  light: string;
  soft: string;
  /** space-separated RGB channel triplet (MUI CSS-vars channel format) */
  channel: string;
}

export const ACCENTS: Record<Accent, AccentSet> = {
  blue: {
    main: '#2190F6',
    dark: '#006AC3',
    light: '#8FCDF9',
    soft: '#EBF5FF',
    channel: '33 144 246',
  },
  indigo: {
    main: '#6974F8',
    dark: '#284AC4',
    light: '#A0A4FF',
    soft: '#EEF0FF',
    channel: '105 116 248',
  },
  teal: {
    main: '#0FB5C9',
    dark: '#087E8F',
    light: '#7FE3F0',
    soft: '#E5FAFD',
    channel: '15 181 201',
  },
};

/** Brand / status colors (mode-independent hues). */
export const STATUS = {
  successBright: '#1DCD9F',
  successDeep: '#03744B',
  warning: '#FBC34D',
  error: '#F50533',
  errorMatt: '#D83C3E',
  orange: '#E27429',
  cyan: '#00D3EB',
  secondary: '#6974F8',
  secondaryDark: '#284AC4',
  secondaryLight: '#A0A4FF',
} as const;

/** Light scheme surfaces & ink (prototype `:root`). */
export const LIGHT = {
  ink: '#1C1E22',
  ink2: '#5A5E66',
  ink3: '#898D95',
  line: '#E7E9EE',
  lineSoft: '#EEF0F4',
  page: '#F4F5F8',
  panel: '#FFFFFF',
  hover: '#F2F4F7',
  shadow: '0 1px 2px rgba(20,22,30,.04), 0 4px 16px rgba(20,22,30,.05)',
  shadowLg: '0 8px 30px rgba(20,22,30,.12)',
} as const;

/** Dark scheme surfaces & ink (prototype `[data-theme="dark"]`). */
export const DARK = {
  ink: '#E9EBEF',
  ink2: '#A6ABB4',
  ink3: '#6C727C',
  line: '#2A2E37',
  lineSoft: '#21242B',
  page: '#0E1015',
  panel: '#171A21',
  hover: '#21252E',
  shadow: '0 1px 2px rgba(0,0,0,.4), 0 6px 20px rgba(0,0,0,.4)',
  shadowLg: '0 16px 44px rgba(0,0,0,.55)',
} as const;

/** Radii (px). */
export const RADII = { sm: 6, md: 10, lg: 14 } as const;

/** Density-dependent sizing (CSS vars; compact values in globals.css). */
export const DENSITY = {
  comfortable: { gap: 16, cardPad: 22, rowH: 52, secGap: 18 },
  compact: { gap: 12, cardPad: 16, rowH: 44, secGap: 14 },
} as const;
