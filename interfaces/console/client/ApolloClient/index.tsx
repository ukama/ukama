/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { ApolloClient, HttpLink, InMemoryCache, split } from '@apollo/client';
import { GraphQLWsLink } from '@apollo/client/link/subscriptions';
import { getMainDefinition } from '@apollo/client/utilities';
import { createClient } from 'graphql-ws';

const client = new ApolloClient({
  uri: process.env.NEXT_PUBLIC_API_GW,
  cache: new InMemoryCache(),
  credentials: 'include',
});

export default client;

const httpLink = (headers: any) =>
  new HttpLink({
    uri: "http://localhost:5042/graphql",
    credentials: 'include',
    headers: {
      ...headers,
    },
  });

const wsLink = new GraphQLWsLink(
  createClient({
    url: "ws://localhost:5042/graphql" || '',
  }),
);

export const MetricLink = () => {
  const _commonData = {
    orgId: '',
    userId: '',
    orgName: '',
  };

  if (typeof window !== 'undefined' && window.localStorage) {
    let data = localStorage.getItem('recoil-persist');
    if (data) {
      let parsedData = JSON.parse(data);
      if (parsedData['commonData']) {
        _commonData.orgId = parsedData['commonData']['orgId'];
        _commonData.userId = parsedData['commonData']['userId'];
        _commonData.orgName = parsedData['commonData']['orgName'];
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
      'org-id': _commonData.orgId,
      'user-id': _commonData.userId,
      'org-name': _commonData.orgName,
    }),
  );
};

export const metricsClient = new ApolloClient({
  link: MetricLink(),
  cache: new InMemoryCache(),
  credentials: 'include',
});
