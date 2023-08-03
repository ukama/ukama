import { ApolloClient, HttpLink, InMemoryCache, split } from '@apollo/client';
import { GraphQLWsLink } from '@apollo/client/link/subscriptions';
import { getMainDefinition } from '@apollo/client/utilities';
import { createClient } from 'graphql-ws';

const client = new ApolloClient({
  uri: process.env.NEXT_PUBLIC_REACT_APP_API,
  cache: new InMemoryCache(),
  credentials: 'include',
});

export default client;

const httpLink = new HttpLink({
  uri: process.env.NEXT_PUBLIC_METRIC_URL,
});

const wsLink = new GraphQLWsLink(
  createClient({
    url: process.env.NEXT_PUBLIC_METRIC_WEBSOCKET_URL || '',
  }),
);

const splitLink = split(
  ({ query }) => {
    const definition = getMainDefinition(query);
    return (
      definition.kind === 'OperationDefinition' &&
      definition.operation === 'subscription'
    );
  },
  wsLink,
  httpLink,
);

export const metricsClient = new ApolloClient({
  link: splitLink,
  cache: new InMemoryCache(),
  credentials: 'include',
});
