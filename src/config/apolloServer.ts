import { ApolloServer } from "apollo-server-express";
import express from "express";
import { createSchema } from "../common/createSchema";

const configureApolloServer = async (): Promise<ApolloServer> => {
    const schema = await createSchema();

    const server = new ApolloServer({
        schema,
        introspection: true,
        context: ({
            req,
            res,
        }: {
            req: express.Request;
            res: express.Response;
        }) => ({ req, res }),
        playground: true,
    });
    return server;
};

export default configureApolloServer;
