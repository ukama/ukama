import config from "../../config";
import { ApolloClient, InMemoryCache } from "@apollo/client";

const client = new ApolloClient({
    uri: config.REACT_APP_API,
    cache: new InMemoryCache(),
    headers: {
        "csrf-token": "test",
        "kratos-session": "test",
    },
});

export default client;
