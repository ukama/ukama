import { ApolloClient, InMemoryCache } from "@apollo/client";

const client = new ApolloClient({
    uri: process.env.REACT_APP_API,
    cache: new InMemoryCache(),
    credentials: "include",
    headers: {
        "Content-Type": "application/json",
    },
});

export default client;
