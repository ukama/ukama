import { ApolloServer } from "apollo-server-express";
import express from "express";
import { createSchema } from "../utils/createSchema";

const configureApolloServer = async (): Promise<ApolloServer> => {
    const schema = await createSchema();

    const server = new ApolloServer({
        schema,
        introspection: true,
        context: ({ req }: { req: express.Request }) => ({ req }),
        playground: true,
    });
    return server;
};

export default configureApolloServer;
