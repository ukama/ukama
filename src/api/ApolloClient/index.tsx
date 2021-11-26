import { ApolloClient, InMemoryCache } from "@apollo/client";
import config from "../../config";

const client = new ApolloClient({
    uri: config.REACT_APP_API,
    cache: new InMemoryCache(),
    headers: {
        authorization: "",
    },
});

export default client;
