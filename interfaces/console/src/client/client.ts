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

const MetricLink = (baseUrl: string, websocketBaseUrl: string) => {
  return split(
    ({ query }) => {
      const definition = getMainDefinition(query);
      return (
        definition.kind === 'OperationDefinition' &&
        definition.operation === 'subscription'
      );
    },

    new GraphQLWsLink(
      createClient({
        url: `${websocketBaseUrl}/graphql`,
      }),
    ),
    new HttpLink({
      uri: `${baseUrl}/graphql`,
      credentials: 'include',
    }),
  );
};

export const getMetricsClient = (baseUrl: string, websocketBaseUrl: string) => {
  return new ApolloClient({
    link: MetricLink(baseUrl, websocketBaseUrl),
    cache: new InMemoryCache(),
    credentials: 'include',
  });
};
