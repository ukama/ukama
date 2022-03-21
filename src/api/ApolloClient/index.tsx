import { WebSocketLink } from "@apollo/client/link/ws";
import { getMainDefinition } from "@apollo/client/utilities";
import { ApolloClient, InMemoryCache, split, HttpLink } from "@apollo/client";

const httpLink = new HttpLink({
    uri: process.env.REACT_APP_API,
    credentials: "include",
});

const wsLink = new WebSocketLink({
    uri: process.env.REACT_APP_API_SOCKET || "ws:localhost",
    options: {
        reconnect: true,
        connectionParams: {
            credentials: "include",
        },
    },
});

const splitLink = split(
    ({ query }) => {
        const definition = getMainDefinition(query);
        return (
            definition.kind === "OperationDefinition" &&
            definition.operation === "subscription"
        );
    },
    wsLink,
    httpLink
);

const client = new ApolloClient({
    link: splitLink,
    cache: new InMemoryCache(),
    credentials: "include",
    headers: {
        Authorization: "Bearer ZCa3ktK4Q3KHBxBXmTGyqJj3QCfI2bI3",
    },
});

export default client;
