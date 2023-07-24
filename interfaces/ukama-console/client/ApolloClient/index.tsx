import { ApolloClient, HttpLink, InMemoryCache, split } from '@apollo/client';
import { GraphQLWsLink } from '@apollo/client/link/subscriptions';
import { getMainDefinition } from '@apollo/client/utilities';
import { createClient } from 'graphql-ws';

const httpLink = new HttpLink({
  uri: 'http://localhost:4042/graphql',
});

const wsLink = new GraphQLWsLink(
  createClient({
    url: 'ws://localhost:4042/graphql',
  }),
);

// The split function takes three parameters:
//
// * A function that's called for each operation to execute
// * The Link to use for an operation if the function returns a "truthy" value
// * The Link to use for an operation if the function returns a "falsy" value
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

// const httpLink = new HttpLink({
//   uri: process.env.NEXT_PUBLIC_REACT_APP_API,
//   credentials: 'include',
// });

// const wsLink = new WebSocketLink({
//   uri: process.env.NEXT_PUBLIC_REACT_APP_API_SOCKET || '',
//   options: {
//     reconnect: true,
//     connectionParams: {
//       credentials: 'include',
//     },
//   },
// });

// const splitLink = split(
//   ({ query }) => {
//     const definition = getMainDefinition(query);
//     return (
//       definition.kind === 'OperationDefinition' &&
//       definition.operation === 'subscription'
//     );
//   },
//   wsLink,
//   httpLink,
// );

const client = new ApolloClient({
  uri: process.env.NEXT_PUBLIC_REACT_APP_API,
  cache: new InMemoryCache(),
  credentials: 'include',
});

export default client;

// const metricsHttpLink = new HttpLink({
//   uri: process.env.NEXT_PUBLIC_REACT_METRICS_API,
//   credentials: 'include',
// });

// const metricsWsLink = new WebSocketLink({
//   uri: process.env.NEXT_PUBLIC_REACT_METRICS_API_SOCKET || '',
//   options: {
//     reconnect: true,
//     connectionParams: {
//       credentials: 'include',
//     },
//   },
// });

// const metricsSplitLink = split(
//   ({ query }) => {
//     const definition = getMainDefinition(query);
//     return (
//       definition.kind === 'OperationDefinition' &&
//       definition.operation === 'subscription'
//     );
//   },
//   metricsWsLink,
//   metricsHttpLink,
// );

export const metricsClient = new ApolloClient({
  link: splitLink,
  cache: new InMemoryCache(),
  credentials: 'include',
});
