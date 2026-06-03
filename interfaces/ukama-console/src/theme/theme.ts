/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * MUI v7 theme — CSS variables + light/dark color schemes (BUILD-PLAN §7.2).
 * Built ONCE at module scope; mode switches flip a class on <html> (no theme
 * rebuild). Accent/density swap CSS custom properties via data-attributes
 * (see globals.css). Mode-specific overrides use theme.applyStyles('dark').
 */
import { createTheme } from '@mui/material/styles';
import { ACCENTS, DARK, LIGHT, RADII, STATUS } from './tokens';

const displayFont = 'var(--font-display), "Work Sans", sans-serif';
const bodyFont = 'var(--font-body), Roboto, sans-serif';

const heading = (size: number, weight = 500) => ({
  fontFamily: displayFont,
  fontSize: size,
  fontWeight: weight,
  lineHeight: 1.2,
});

export const theme = createTheme({
  cssVariables: { colorSchemeSelector: 'class' },
  colorSchemes: {
    light: {
      palette: {
        primary: {
          main: ACCENTS.blue.main,
          dark: ACCENTS.blue.dark,
          light: ACCENTS.blue.light,
          contrastText: '#FFFFFF',
        },
        secondary: {
          main: STATUS.secondary,
          dark: STATUS.secondaryDark,
          light: STATUS.secondaryLight,
          contrastText: '#FFFFFF',
        },
        success: { main: STATUS.successBright, dark: STATUS.successDeep },
        warning: { main: STATUS.warning },
        error: { main: STATUS.error },
        background: { default: LIGHT.page, paper: LIGHT.panel },
        text: {
          primary: LIGHT.ink,
          secondary: LIGHT.ink2,
          disabled: LIGHT.ink3,
        },
        divider: LIGHT.line,
        action: { hover: LIGHT.hover },
      },
    },
    dark: {
      palette: {
        primary: {
          main: ACCENTS.blue.main,
          dark: '#5FB0FF',
          light: ACCENTS.blue.light,
          contrastText: '#FFFFFF',
        },
        secondary: {
          main: STATUS.secondary,
          dark: STATUS.secondaryDark,
          light: STATUS.secondaryLight,
          contrastText: '#FFFFFF',
        },
        success: { main: STATUS.successBright, dark: STATUS.successDeep },
        warning: { main: STATUS.warning },
        error: { main: '#FF7A7A' },
        background: { default: DARK.page, paper: DARK.panel },
        text: {
          primary: DARK.ink,
          secondary: DARK.ink2,
          disabled: DARK.ink3,
        },
        divider: DARK.line,
        action: { hover: DARK.hover },
      },
    },
  },
  typography: {
    fontFamily: bodyFont,
    h1: heading(34),
    h2: heading(30),
    h3: heading(27),
    h4: heading(24),
    h5: heading(20),
    h6: heading(17, 600),
    subtitle1: { fontWeight: 600 },
    button: { textTransform: 'none', fontWeight: 600 },
  },
  shape: { borderRadius: RADII.md },
  components: {
    MuiButton: {
      defaultProps: { disableElevation: true },
      styleOverrides: {
        root: { borderRadius: 8 },
      },
    },
    MuiPaper: {
      styleOverrides: {
        root: { backgroundImage: 'none' },
      },
    },
    MuiCard: {
      styleOverrides: {
        root: ({ theme: t }) => ({
          borderRadius: RADII.lg,
          border: `1px solid ${(t.vars ?? t).palette.divider}`,
          boxShadow: 'var(--uk-shadow)',
        }),
      },
    },
    MuiChip: {
      styleOverrides: {
        root: { fontWeight: 600 },
      },
    },
    MuiTooltip: {
      styleOverrides: {
        tooltip: { fontSize: 12.5 },
      },
    },
  },
});
