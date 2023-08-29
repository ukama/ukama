import { ApolloServer } from "@apollo/server";
import { ApolloServerPluginInlineTrace } from "@apollo/server/plugin/inlineTrace";
import { buildSubgraphSchema, printSubgraphSchema } from "@apollo/subgraph";
import { GraphQLScalarType } from "graphql";
import { DateTimeResolver } from "graphql-scalars";
import { gql } from "graphql-tag";
import * as tq from "type-graphql";
import { NonEmptyArray } from "type-graphql";

const SubGraphServer = async (resolvers: NonEmptyArray<any>) => {
  const ts = await tq.buildSchema({
    resolvers: resolvers,
    scalarsMap: [{ type: GraphQLScalarType, scalar: DateTimeResolver }],
    validate: { forbidUnknownValues: false },
  });

  const federatedSchema = buildSubgraphSchema({
    typeDefs: gql(printSubgraphSchema(ts)),
    resolvers: tq.createResolversMap(ts) as any,
  });

  const server = new ApolloServer({
    schema: federatedSchema,
    csrfPrevention: false,
    introspection: true,
    plugins: [ApolloServerPluginInlineTrace({})],
  });
  return server;
};

export default SubGraphServer;
