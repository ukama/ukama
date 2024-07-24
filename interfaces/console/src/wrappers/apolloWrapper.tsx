/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';

import { ApolloLink, HttpLink } from '@apollo/client';
import {
  ApolloClient,
  ApolloNextAppProvider,
  InMemoryCache,
  SSRMultipartLink,
} from '@apollo/experimental-nextjs-app-support';

function makeClient() {
  const httpLink = new HttpLink({
    uri: `${process.env.NEXT_PUBLIC_API_GW}/graphql`,
    credentials: 'include',
  });

  return new ApolloClient({
    cache: new InMemoryCache(),
    link:
      typeof window === 'undefined'
        ? ApolloLink.from([
            new SSRMultipartLink({
              stripDefer: true,
            }),
            httpLink,
          ])
        : httpLink,
  });
}

export function ApolloWrapper({ children }: Readonly<React.PropsWithChildren>) {
  console.log(process.env.NEXT_PUBLIC_API_GW);
  console.log(process.env.NEXT_PUBLIC_APP_URL);
  console.log(process.env.NEXT_PUBLIC_METRIC_URL);
  console.log(process.env.NEXT_PUBLIC_AUTH_APP_URL);
  console.log(process.env.NEXT_PUBLIC_METRIC_WEBSOCKET_URL);
  console.log(process.env.NEXT_PUBLIC_MAP_BOX_TOKEN);
  console.log(process.env.NEXT_PUBLIC_API_GW_4SS);
  console.log(process.env.NEXT_PUBLIC_SIM_TYPE);

  return (
    <ApolloNextAppProvider makeClient={makeClient}>
      {children}
    </ApolloNextAppProvider>
  );
}
