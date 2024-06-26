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

const httpLink = () =>
  new HttpLink({
    uri: `${process.env.NEXT_PUBLIC_METRIC_URL}/graphql`,
    credentials: 'include',
  });

const wsLink = new GraphQLWsLink(
  createClient({
    url: `${process.env.NEXT_PUBLIC_METRIC_WEBSOCKET_URL}/graphql`,
  }),
);

export const MetricLink = () => {
  return split(
    ({ query }) => {
      const definition = getMainDefinition(query);
      return (
        definition.kind === 'OperationDefinition' &&
        definition.operation === 'subscription'
      );
    },

    wsLink,
    httpLink(),
  );
};

export const metricsClient = new ApolloClient({
  link: MetricLink(),
  cache: new InMemoryCache(),
  credentials: 'include',
});
