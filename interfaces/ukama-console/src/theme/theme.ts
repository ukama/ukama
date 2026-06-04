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
        root: ({ ownerState, theme: t }) => {
          const tone = (
            bg: string,
            color: string,
            dark: { bg: string; color: string },
          ) => ({
            background: bg,
            color,
            ...t.applyStyles('dark', { background: dark.bg, color: dark.color }),
          });
          const badge = {
            fontSize: 12,
            fontWeight: 600,
            height: 'auto',
            borderRadius: 20,
            lineHeight: 1.4,
            padding: '3px 9px',
            '& .MuiChip-label': { padding: 0 },
            '& .MuiChip-icon': { margin: '0 6px 0 -1px' },
          };
          return {
            fontWeight: 600,
            ...(ownerState.variant === 'ok' && {
              ...badge,
              ...tone('var(--uk-success-fill)', 'var(--uk-success)', {
                bg: 'rgba(29,205,159,.16)',
                color: 'var(--uk-success-bright)',
              }),
            }),
            ...(ownerState.variant === 'warn' && {
              ...badge,
              ...tone('rgba(251,195,77,.2)', '#946600', {
                bg: 'rgba(251,195,77,.16)',
                color: '#fbc34d',
              }),
            }),
            ...(ownerState.variant === 'err' && {
              ...badge,
              ...tone('var(--uk-error-fill)', 'var(--uk-error-deep, #cf121b)', {
                bg: 'rgba(245,5,51,.18)',
                color: '#ff7a7a',
              }),
            }),
            ...(ownerState.variant === 'neut' && {
              ...badge,
              ...tone('var(--uk-page)', 'var(--uk-ink-2)', {
                bg: 'var(--uk-hover)',
                color: 'var(--uk-ink-2)',
              }),
            }),
            ...(ownerState.variant === 'info' && {
              ...badge,
              ...tone('var(--uk-ac-soft)', 'var(--uk-ac-dark)', {
                bg: 'rgba(33,144,246,.2)',
                color: 'var(--uk-ac-dark)',
              }),
            }),
            ...(ownerState.variant === 'chipFilter' && {
              height: 32,
              borderRadius: 8,
              border: '1px solid var(--uk-line)',
              background: 'var(--uk-panel)',
              fontSize: 13,
              fontWeight: 400,
              color: 'var(--uk-ink-2)',
              cursor: 'pointer',
              transition: '.12s',
              '&:hover': { background: 'var(--uk-hover)' },
              '& .MuiChip-label': {
                padding: '0 12px',
                display: 'flex',
                gap: 7,
                alignItems: 'center',
              },
              '& .MuiChip-icon': { marginLeft: 10, marginRight: -6 },
            }),
          };
        },
      },
    },
    MuiTabs: {
      styleOverrides: {
        root: {
          minHeight: 40,
          borderBottom: '1px solid var(--uk-line)',
          marginBottom: 20,
        },
        indicator: { backgroundColor: 'var(--uk-ac)', height: 2 },
      },
    },
    MuiTab: {
      styleOverrides: {
        root: {
          textTransform: 'none',
          fontSize: 14,
          fontWeight: 500,
          color: 'var(--uk-ink-2)',
          minHeight: 40,
          minWidth: 0,
          padding: '10px 14px',
          '&.Mui-selected': { color: 'var(--uk-ac-dark)' },
        },
      },
    },
    MuiTextField: {
      defaultProps: { size: 'small' },
    },
    MuiOutlinedInput: {
      styleOverrides: {
        root: {
          borderRadius: 8,
          background: 'var(--uk-panel)',
          fontSize: 14,
          '& .MuiOutlinedInput-notchedOutline': { borderColor: 'var(--uk-line)' },
          '&:hover .MuiOutlinedInput-notchedOutline': { borderColor: 'var(--uk-line)' },
          '&.Mui-focused': {
            boxShadow: '0 0 0 3px color-mix(in srgb, var(--uk-ac) 12%, transparent)',
          },
          '&.Mui-focused .MuiOutlinedInput-notchedOutline': {
            borderColor: 'var(--uk-ac)',
            borderWidth: 1,
          },
          '&.Mui-error .MuiOutlinedInput-notchedOutline': {
            borderColor: 'var(--uk-error)',
          },
        },
        input: { padding: '8px 11px' },
      },
    },
    MuiSwitch: {
      styleOverrides: {
        root: { width: 40, height: 22, padding: 0 },
        switchBase: {
          padding: 2,
          '&.Mui-checked': {
            transform: 'translateX(18px)',
            color: '#fff',
            '& + .MuiSwitch-track': { backgroundColor: 'var(--uk-ac)', opacity: 1 },
          },
        },
        thumb: { width: 18, height: 18, boxShadow: '0 1px 2px rgba(0,0,0,.3)' },
        track: { borderRadius: 20, backgroundColor: 'var(--uk-line)', opacity: 1 },
      },
    },
    MuiLinearProgress: {
      styleOverrides: {
        root: { height: 6, borderRadius: 6, backgroundColor: 'var(--uk-line)' },
        bar: { borderRadius: 6 },
      },
    },
    MuiTable: {
      styleOverrides: {
        root: { tableLayout: 'fixed', minWidth: 760, width: '100%' },
      },
    },
    MuiTableCell: {
      styleOverrides: {
        root: { fontFamily: 'inherit' },
        head: {
          fontSize: 11.5,
          fontWeight: 600,
          letterSpacing: '.04em',
          textTransform: 'uppercase',
          color: 'var(--uk-ink-3)',
          padding: '0 14px 11px',
          borderBottom: '1px solid var(--uk-line)',
          whiteSpace: 'nowrap',
        },
        body: {
          padding: '0 14px',
          height: 'var(--uk-row-h)',
          fontSize: 13.5,
          color: 'var(--uk-ink)',
          borderBottom: '1px solid var(--uk-line-soft)',
          whiteSpace: 'nowrap',
          overflow: 'hidden',
          textOverflow: 'ellipsis',
        },
      },
    },
    MuiTableRow: {
      styleOverrides: {
        root: {
          '&:last-child .MuiTableCell-body': { borderBottom: 0 },
          '&.MuiTableRow-hover:hover': { backgroundColor: 'var(--uk-hover)' },
        },
      },
    },
    MuiTooltip: {
      styleOverrides: {
        tooltip: { fontSize: 12.5 },
      },
    },
  },
});
