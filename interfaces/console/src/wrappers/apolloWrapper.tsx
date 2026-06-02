/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';

import { ApolloLink, HttpLink } from '@apollo/client';
import { onError } from '@apollo/client/link/error';
import { RetryLink } from '@apollo/client/link/retry';
import {
  ApolloClient,
  ApolloNextAppProvider,
  InMemoryCache,
  SSRMultipartLink,
} from '@apollo/experimental-nextjs-app-support';

const REQUEST_TIMEOUT_MS = 15_000;

function fetchWithTimeout(input: RequestInfo | URL, init?: RequestInit): Promise<Response> {
  const controller = new AbortController();
  const timer = setTimeout(() => controller.abort(), REQUEST_TIMEOUT_MS);
  return fetch(input, { ...init, signal: controller.signal }).finally(() =>
    clearTimeout(timer),
  );
}

function makeClient(baseUrl: string) {
  const httpLink = new HttpLink({
    uri: `${baseUrl}/graphql`,
    credentials: 'include',
    fetch: fetchWithTimeout,
  });

  const errorLink = onError(({ graphQLErrors, networkError, operation }) => {
    if (graphQLErrors) {
      graphQLErrors.forEach(({ message, locations, path }) => {
        console.error(
          `[GraphQL error] op=${operation.operationName} message=${message} path=${path} locations=${JSON.stringify(locations)}`,
        );
      });
    }
    if (networkError) {
      console.error(`[Network error] op=${operation.operationName}:`, networkError);
      if ('statusCode' in networkError && networkError.statusCode === 401) {
        if (typeof window !== 'undefined') {
          window.location.href = `${process.env.NEXT_PUBLIC_AUTH_APP_URL}/auth/login`;
        }
      }
    }
  });

  const retryLink = new RetryLink({
    delay: { initial: 300, max: 3000, jitter: true },
    attempts: { max: 3, retryIf: (error) => !!error && error?.statusCode !== 401 },
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
          getSubscribers: {
            keyArgs: ['data', ['orgName', 'networkId']],
            merge(_existing, incoming) {
              return incoming;
            },
          },
          getMembers: {
            keyArgs: ['data', ['orgName']],
            merge(_existing, incoming) {
              return incoming;
            },
          },
          getInvitations: {
            keyArgs: ['data', ['orgName']],
            merge(_existing, incoming) {
              return incoming;
            },
          },
          getPackages: {
            keyArgs: ['data', ['orgName']],
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
      SubscriberDto: { keyFields: ['uuid'] },
      PackageDto: { keyFields: ['uuid'] },
      MemberDto: { keyFields: ['memberId'] },
      InvitationDto: { keyFields: ['id'] },
      SimDto: { keyFields: ['iccid'] },
      NetworkDto: { keyFields: ['id'] },
      UserResDto: { keyFields: ['uuid'] },
    },
  });

  const clientLink = ApolloLink.from([errorLink, retryLink, httpLink]);

  return new ApolloClient({
    cache,
    connectToDevTools: process.env.NODE_ENV === 'development',
    link:
      typeof window === 'undefined'
        ? ApolloLink.from([
            new SSRMultipartLink({ stripDefer: true }),
            errorLink,
            httpLink,
          ])
        : clientLink,
    defaultOptions: {
      watchQuery: {
        fetchPolicy: 'cache-and-network',
        nextFetchPolicy: 'cache-first',
        refetchWritePolicy: 'merge',
      },
      query: {
        fetchPolicy: 'cache-first',
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
