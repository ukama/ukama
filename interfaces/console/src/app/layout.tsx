/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Sim_Types } from '@/client/graphql/generated';
import AppContextWrapper from '@/context';
import '@/styles/global.css';
import AppThemeProvider from '@/theme/AppThemeProvider';
import { ApolloWrapper } from '@/wrappers/apolloWrapper';
import type { Metadata } from 'next';
import { Inter } from 'next/font/google';
import { cookies, headers } from 'next/headers';

const inter = Inter({ subsets: ['latin'] });

export const metadata: Metadata = {
  title: 'Ukama Console',
  description: 'Ukama Conosle app to manage your network',
  icons: {
    icon: [
      {
        url: '/ulogo.svg',
      },
    ],
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const headersList = headers();
  const cookieStore = cookies();
  const cookieTheme = cookieStore.get('theme') ?? {
    name: 'theme',
    value: 'light',
  };
  const role = headersList.get('role');
  const name = headersList.get('name');
  const email = headersList.get('email');
  const orgId = headersList.get('org-id');
  const userId = headersList.get('user-id');
  const country = headersList.get('country');
  const orgName = headersList.get('org-name');
  const currency = headersList.get('currency');
  const tokenStr = cookieStore.get('token') ?? {
    name: 'token',
    value: '',
  };
  return (
    <html lang="en">
      <head>
        <link rel="icon" type="image/svg+xml" href="/ulogo.svg" />
      </head>
      <body className={inter.className}>
        <ApolloWrapper baseUrl={process.env.NEXT_PUBLIC_API_GW ?? ''}>
          <AppContextWrapper
            token={tokenStr.value}
            initEnv={{
              APP_URL: process.env.NEXT_PUBLIC_APP_URL ?? '',
              SIM_TYPE:
                process.env.NEXT_PUBLIC_SIM_TYPE ?? Sim_Types.OperatorData,
              METRIC_URL: process.env.NEXT_PUBLIC_METRIC_URL ?? '',
              API_GW_URL: process.env.NEXT_PUBLIC_API_GW ?? '',
              AUTH_APP_URL: process.env.NEXT_PUBLIC_AUTH_APP_URL ?? '',
              MAP_BOX_TOKEN: process.env.NEXT_PUBLIC_MAP_BOX_TOKEN ?? '',
              METRIC_WEBSOCKET_URL:
                process.env.NEXT_PUBLIC_METRIC_WEBSOCKET_URL ?? '',
              STRIPE_PK: process.env.NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY ?? '',
            }}
            initalUserValues={{
              id: userId ?? '',
              name: name ?? '',
              role: role ?? '',
              email: email ?? '',
              orgId: orgId ?? '',
              orgName: orgName ?? '',
              country: country ?? '',
              currency: currency ?? '',
            }}
          >
            <AppThemeProvider themeCookie={cookieTheme?.value}>
              {children}
            </AppThemeProvider>
          </AppContextWrapper>
        </ApolloWrapper>
      </body>
    </html>
  );
}
