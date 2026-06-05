/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Apollo Client factory (Phase 2 API layer).
 *
 * Decisions (plan §15.5):
 *  - client-only + skeletons: the client is created once in the browser
 *    (ApolloClientProvider), no RSC/streaming SSR integration needed
 *  - cache-first with data expiry: InvalidationPolicyCache (TTL) from
 *    @nerdwallet/apollo-cache-policies on @apollo/client ^3.14
 *  - polling-only v1 (no websockets); the separate metrics-endpoint
 *    client + SSE subscriptions come with the metrics phase
 *  - auth: cookies (`credentials: 'include'`) against the API gateway,
 *    same as the legacy console; session/role headers via proxy.ts later
 */
import { env } from '@/env';
import { ApolloClient, HttpLink, from } from '@apollo/client';
import { onError } from '@apollo/client/link/error';
import { InvalidationPolicyCache } from '@nerdwallet/apollo-cache-policies';

const MINUTE = 60 * 1000;

/** Default freshness window — reads older than this refetch (plan §5). */
export const DEFAULT_TTL = 5 * MINUTE;

/** Per-type freshness overrides: volatile data expires faster. */
const TYPE_TTLS: Record<string, number> = {
  // ops status changes quickly
  NodeDto: 1 * MINUTE,
  // activation state must react quickly to setup progress (alert bar/guards)
  OnboardingStatusDto: 0.5 * MINUTE,
  SiteDto: 2 * MINUTE,
  NotificationsResDto: 1 * MINUTE,
  // money/people move slower
  PackageDto: 10 * MINUTE,
  SubscriberDto: 5 * MINUTE,
  MemberDto: 10 * MINUTE,
  NetworkDto: 10 * MINUTE,
  // reference data barely changes
  CountryDto: 24 * 60 * MINUTE,
  TimezoneDto: 24 * 60 * MINUTE,
};

function makeCache() {
  return new InvalidationPolicyCache({
    typePolicies: {
      Query: {
        fields: {
          // paginated/filterable lists: replace on refetch (legacy parity)
          getSites: { merge: (_, incoming: unknown) => incoming },
          getNodes: { merge: (_, incoming: unknown) => incoming },
        },
      },
    },
    invalidationPolicies: {
      timeToLive: DEFAULT_TTL,
      types: Object.fromEntries(
        Object.entries(TYPE_TTLS).map(([t, ttl]) => [t, { timeToLive: ttl }]),
      ),
    },
  });
}

/** Guards against repeated redirects if a refresh still returns 401. */
let isRecovering = false;

export function makeApolloClient() {
  const errorLink = onError(({ graphQLErrors, networkError, operation }) => {
    if (graphQLErrors) {
      for (const err of graphQLErrors) {
        console.error(`[gql] ${operation.operationName}: ${err.message}`);
        // Gateway rejected the token (expired/rotated): drop the cached token
        // cookie and let proxy.ts mint a fresh one from the still-valid session.
        if (
          err.extensions?.code === 'UNAUTHENTICATED' &&
          typeof window !== 'undefined' &&
          !isRecovering
        ) {
          isRecovering = true;
          window.location.assign('/api/auth/refresh');
        }
      }
    }
    if (networkError) {
      console.error(`[gql/network] ${operation.operationName}:`, networkError);
    }
  });

  const httpLink = new HttpLink({
    uri: `${env.NEXT_PUBLIC_API_GW}/graphql`,
    credentials: 'include',
  });

  return new ApolloClient({
    link: from([errorLink, httpLink]),
    cache: makeCache(),
    defaultOptions: {
      // cache-first + TTL: instant from cache while fresh, refetch when stale
      watchQuery: {
        fetchPolicy: 'cache-first',
        errorPolicy: 'all',
        refetchWritePolicy: 'merge',
      },
      query: { fetchPolicy: 'cache-first', errorPolicy: 'all' },
    },
  });
}
