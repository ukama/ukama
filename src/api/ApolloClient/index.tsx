import { ApolloClient, InMemoryCache } from "@apollo/client";

const client = new ApolloClient({
    uri: process.env.REACT_APP_API,
    cache: new InMemoryCache(),
    headers: {
        "csrf-token": "test",
        "ukama-session": "test",
    },
});

export default client;
