/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

'use client';
import { ThemeProvider, createTheme, useMediaQuery } from '@mui/material';
import { AppRouterCacheProvider } from '@mui/material-nextjs/v13-appRouter';
import { Rubik, Work_Sans } from 'next/font/google';
import { useEffect, useState } from 'react';
import { colors } from '.';

const rubik = Rubik({
  subsets: ['latin'],
  display: 'swap',
});
const workSans = Work_Sans({
  subsets: ['latin'],
  display: 'swap',
});

export default function AppThemeProvider({
  children,
  themeCookie,
}: Readonly<{
  children: React.ReactNode;
  themeCookie: string;
}>) {
  const prefersMode = useMediaQuery(`(prefers-color-scheme: ${themeCookie})`);

  const [mode, setMode] = useState(prefersMode);

  useEffect(() => {
    setMode(prefersMode);
  }, [prefersMode]);

  const textColor = mode ? colors.white : colors.black;

  const theme = createTheme({
    typography: {
      fontFamily: rubik.style.fontFamily,
      h1: {
        fontFamily: rubik.style.fontFamily,
        color: textColor,
      },
      h2: {
        fontFamily: rubik.style.fontFamily,
        color: textColor,
      },
      h3: {
        fontFamily: rubik.style.fontFamily,
        color: textColor,
      },
      h4: {
        fontFamily: rubik.style.fontFamily,
        color: textColor,
      },
      h5: {
        fontFamily: rubik.style.fontFamily,
        color: textColor,
      },
      h6: {
        fontFamily: rubik.style.fontFamily,
        color: textColor,
      },
      subtitle1: { fontFamily: workSans.style.fontFamily },
      subtitle2: { fontFamily: workSans.style.fontFamily },
      body1: {
        fontFamily: workSans.style.fontFamily,
        letterSpacing: '-0.02em',
      },
      body2: {
        fontFamily: workSans.style.fontFamily,
        letterSpacing: '-0.02em',
      },
      caption: {
        fontFamily: workSans.style.fontFamily,
      },
    },
    palette: {
      mode: mode ? 'dark' : 'light',
      text: {
        primary: mode ? colors.white : colors.vulcan,
        secondary: mode ? colors.white70 : colors.black70,
        disabled: mode ? colors.white38 : colors.black38,
      },
      background: {
        default: mode ? colors.nightGrey5 : colors.solitude,
        paper: mode ? colors.nightGrey5 : colors.white,
      },
      primary: {
        main: colors.primaryMain,
        light: colors.primaryLight,
        dark: colors.primaryDark,
      },
      secondary: {
        main: colors.secondaryMain,
        light: colors.secondaryLight,
        dark: colors.secondaryDark,
      },
      error: {
        main: colors.error,
      },
    },
    breakpoints: {
      values: {
        xs: 0,
        sm: 600,
        md: 900,
        lg: 1200,
        xl: 1536,
      },
      step: 8,
    },
    components: {
      MuiInputBase: {
        styleOverrides: {
          root: {
            '.MuiOutlinedInput-input': {
              fontSize: '14px',
              padding: '11px !important',
            },
          },
        },
      },
      MuiFormLabel: {
        styleOverrides: {
          root: {
            fontSize: '14px !important',
          },
        },
      },
      MuiFormControl: {
        styleOverrides: {
          root: {
            '&:hover .MuiOutlinedInput-root .MuiOutlinedInput-notchedOutline': {
              borderColor: colors.hoverColor,
            },
          },
        },
      },
      MuiDivider: {
        styleOverrides: {
          root: {
            margin: '12px 0px',
          },
        },
      },
      MuiFormHelperText: {
        styleOverrides: {
          contained: {
            marginLeft: '0px !important',
          },
        },
      },
      MuiIconButton: {
        styleOverrides: {
          root: {
            '&:hover': {
              backgroundColor: 'transparent',
            },
            '&:hover svg path': {
              fill: colors.primaryMain,
            },
            '&:hover svg circle': {
              fill: colors.primaryMain,
            },
          },
        },
      },
      MuiButton: {
        styleOverrides: {
          contained: {
            fontWeight: 500,
            color: colors.white,
            letterSpacing: '0.4px',
            boxShadow:
              '0px 3px 1px -2px rgba(0, 0, 0, 0.2), 0px 2px 2px rgba(0, 0, 0, 0.14), 0px 1px 5px rgba(0, 0, 0, 0.12)',
          },
          text: {
            padding: '0px',
            minWidth: 'auto',
            color: mode ? colors.white : colors.darkGray,
            ':hover': {
              color: colors.primaryMain,
              backgroundColor: 'transparent !important',
            },
          },
          sizeMedium: {
            padding: '7px 24px',
          },
        },
      },
      MuiSelect: {
        styleOverrides: {
          select: {
            height: '20px',
            fontSize: '16px !important',
            ':focus': {
              backgroundColor: 'transparent',
            },
          },
          iconStandard: {
            paddingLeft: '4px',
          },
        },
      },
      MuiDialogContent: {
        styleOverrides: {
          root: {
            padding: '8px 24px 24px',
          },
        },
      },
      MuiDialogActions: {
        styleOverrides: {
          root: {
            padding: '8px 24px 24px',
          },
        },
      },
      MuiDrawer: {
        styleOverrides: {
          paper: {
            borderRight: 'none',
          },
        },
      },
      MuiTableCell: {
        styleOverrides: {
          root: {
            backgroundColor: 'transparent',
          },
          stickyHeader: {
            backgroundColor: mode ? '#292929' : colors.white,
          },
        },
      },
      MuiPaper: {
        styleOverrides: {
          root: {
            boxShadow: '2px 2px 6px rgba(0, 0, 0, 0.05)',
          },
        },
      },
      MuiPopover: {
        styleOverrides: {
          paper: {
            boxShadow:
              '0px 5px 5px -3px rgba(0, 0, 0, 0.2), 0px 8px 10px 1px rgba(0, 0, 0, 0.14), 0px 3px 14px 2px rgba(0, 0, 0, 0.12)',
          },
        },
      },
      MuiListItemButton: {
        styleOverrides: {
          root: {
            marginLeft: 20,
            marginRight: 20,
            borderRadius: 4,
            '@media all and (max-width: 900px)': {
              marginLeft: 10,
              marginRight: 10,
              padding: 8,
            },
          },
        },
      },
      MuiListItemIcon: {
        styleOverrides: {
          root: {
            minWidth: 'auto',
            marginRight: '1.25rem',
            svg: {
              width: '1.5rem',
              height: '1.5rem',
              fill: mode ? colors.white : colors.vulcan,
            },
          },
        },
      },
      MuiListItemText: {
        styleOverrides: {
          root: {
            '.MuiTypography-root': {
              fontSize: '1rem',
              color: mode ? colors.white : colors.vulcan,
              fontFamily: workSans.style.fontFamily,
            },
          },
        },
      },
      MuiToolbar: {
        styleOverrides: {
          root: {
            minHeight: 'inherit !important',
            paddingLeft: '36px !important',
            paddingRight: '36px !important',
            '@media all and (max-width: 900px)': {
              paddingLeft: '10px !important',
              paddingRight: '10px !important',
            },
          },
        },
      },
      MuiChip: {
        styleOverrides: {
          root: {
            px: 1,
            fontSize: '0.875rem',
            color: colors.primaryMain,
            borderColor: colors.primaryMain,
          },
        },
      },
    },
  });

  return (
    <AppRouterCacheProvider>
      <ThemeProvider theme={theme}>{children}</ThemeProvider>
    </AppRouterCacheProvider>
  );
}
