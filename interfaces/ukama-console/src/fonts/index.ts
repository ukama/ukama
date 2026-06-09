/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import localFont from 'next/font/local';

/** Display face — page titles, KPI values, dialog titles (design: Work Sans). */
export const workSans = localFont({
  src: [
    {
      path: './WorkSans-VariableFont_wght.ttf',
      weight: '100 900',
      style: 'normal',
    },
    {
      path: './WorkSans-Italic-VariableFont_wght.ttf',
      weight: '100 900',
      style: 'italic',
    },
  ],
  variable: '--font-display',
  display: 'swap',
});

/** Body face — everything else (design: Roboto). */
export const roboto = localFont({
  src: [
    {
      path: './Roboto-VariableFont_wdth_wght.ttf',
      weight: '100 900',
      style: 'normal',
    },
    {
      path: './Roboto-Italic-VariableFont_wdth_wght.ttf',
      weight: '100 900',
      style: 'italic',
    },
  ],
  variable: '--font-body',
  display: 'swap',
});
