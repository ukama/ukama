import { ApolloClient, InMemoryCache } from "@apollo/client";

const client = new ApolloClient({
    uri: process.env.REACT_APP_API,
    cache: new InMemoryCache(),
    headers: {
        authorization: "",
    },
});

export default client;
