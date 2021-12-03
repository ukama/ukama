import { ApolloServer, PubSub } from "apollo-server-express";
import express from "express";
import { GraphQLSchema } from "graphql";
import { createSchema } from "../common/createSchema";

const configureApolloServer = async (): Promise<{
    server: ApolloServer;
    schema: GraphQLSchema;
}> => {
    const pubsub = new PubSub();

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
        }) => ({ req, res, pubsub }),

        playground: true,
    });
    return { server, schema };
};

export default configureApolloServer;
