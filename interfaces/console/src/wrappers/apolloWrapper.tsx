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

function makeClient(baseUrl: string) {
  const httpLink = new HttpLink({
    uri: `${baseUrl}/graphql`,
    credentials: 'include',
  });

  const cache = new InMemoryCache({
    typePolicies: {
      Query: {
        fields: {
          getSites: {
            merge(_existing, incoming) {
              return incoming;
            },
          },
          getNodes: {
            merge(_existing, incoming) {
              return incoming;
            },
            read(existing) {
              return existing;
            },
          },
          getCurrencySymbol: {
            merge(_existing, incoming) {
              return incoming;
            },
          },
        },
      },
      Site: {
        keyFields: ['id'],
        fields: {
          nodes: {
            merge: true,
          },
        },
      },
      Node: {
        keyFields: ['id'],
        fields: {
          status: {
            merge: true,
          },
          site: {
            merge: true,
          },
        },
      },
    },
  });

  return new ApolloClient({
    cache,
    link:
      typeof window === 'undefined'
        ? ApolloLink.from([
            new SSRMultipartLink({
              stripDefer: true,
            }),
            httpLink,
          ])
        : httpLink,
    defaultOptions: {
      watchQuery: {
        fetchPolicy: 'cache-and-network',
        nextFetchPolicy: 'cache-first',
        refetchWritePolicy: 'merge',
      },
      query: {
        fetchPolicy: 'network-only',
        errorPolicy: 'all',
      },
    },
  });
}

export function ApolloWrapper({
  baseUrl,
  children,
}: {
  baseUrl: string;
  children: React.ReactNode;
}) {
  return (
    <ApolloNextAppProvider makeClient={() => makeClient(baseUrl)}>
      {children}
    </ApolloNextAppProvider>
  );
}
