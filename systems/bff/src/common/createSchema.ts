import { buildSchema } from "type-graphql";
import { GraphQLSchema } from "graphql";
import { Container } from "typedi";

export const createSchema = (): Promise<GraphQLSchema> =>
    buildSchema({
        resolvers: [__dirname + "/../modules/**/*.resolver.*"],
        validate: true,
        container: Container,
    });
