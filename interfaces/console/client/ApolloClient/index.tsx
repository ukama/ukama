/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { useAppContext } from '@/context';
import { ApolloClient, HttpLink, InMemoryCache, split } from '@apollo/client';
import { GraphQLWsLink } from '@apollo/client/link/subscriptions';
import { getMainDefinition } from '@apollo/client/utilities';
import { createClient } from 'graphql-ws';

const client = new ApolloClient({
  uri: `${process.env.NEXT_PUBLIC_API_GW}/graphql`,
  cache: new InMemoryCache(),
  credentials: 'include',
});

export default client;

const httpLink = (headers: any) =>
  new HttpLink({
    uri: process.env.NEXT_PUBLIC_METRICS_URL,
    credentials: 'include',
    headers: {
      ...headers,
    },
  });

const wsLink = new GraphQLWsLink(
  createClient({
    url: process.env.NEXT_PUBLIC_METRICS_WEBSOCKET_URL || '',
  }),
);

export const MetricLink = () => {
  // const { user } = useAppContext();

  if (typeof window !== 'undefined' && window.localStorage) {
    let data = localStorage.getItem('recoil-persist');
    if (data) {
      let parsedData = JSON.parse(data);
      if (parsedData['commonData']) {
        user.orgId = parsedData['commonData']['orgId'];
        user.id = parsedData['commonData']['userId'];
        user.orgName = parsedData['commonData']['orgName'];
      }
    }
  }
  return split(
    ({ query }) => {
      const definition = getMainDefinition(query);
      return (
        definition.kind === 'OperationDefinition' &&
        definition.operation === 'subscription'
      );
    },

    wsLink,
    httpLink({
      'user-id': 'user.id',
      'org-id': 'user.orgId',
      'org-name': 'user.orgName',
    }),
  );
};

export const metricsClient = new ApolloClient({
  link: MetricLink(),
  cache: new InMemoryCache(),
  credentials: 'include',
});
