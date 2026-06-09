/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { roboto, workSans } from '@/fonts';
import { getCurrentUser } from '@/lib/auth/server';
import { runtimeEnvScript } from '@/lib/runtime-env';
import { AppRouterCacheProvider } from '@mui/material-nextjs/v15-appRouter';
import InitColorSchemeScript from '@mui/material/InitColorSchemeScript';
import type { Metadata } from 'next';
import './components.css';
import './globals.css';
import Providers from './providers';

export const metadata: Metadata = {
  title: 'Ukama Console',
  description:
    'Ukama network operator console — Business, Network and Customer lenses.',
};

export default async function RootLayout({
  children,
}: Readonly<{ children: React.ReactNode }>) {
  const user = await getCurrentUser();

  return (
    <html
      lang="en"
      className={`${workSans.variable} ${roboto.variable}`}
      suppressHydrationWarning
    >
      <head>
        <script dangerouslySetInnerHTML={{ __html: runtimeEnvScript() }} />
      </head>
      <body>
        <InitColorSchemeScript attribute="class" defaultMode="light" />
        <AppRouterCacheProvider>
          <Providers user={user}>{children}</Providers>
        </AppRouterCacheProvider>
      </body>
    </html>
  );
}
