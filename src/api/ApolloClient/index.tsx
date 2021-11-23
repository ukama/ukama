import { ApolloClient, InMemoryCache } from "@apollo/client";

const client = new ApolloClient({
    uri: process.env.REACT_APP_API,
    cache: new InMemoryCache(),
    headers: {
        CSRF_TOKEN: "token",
        KRATOS_SESSION: "session",
    },
});

export default client;
