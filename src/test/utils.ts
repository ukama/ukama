import { ExecutionResult, graphql, GraphQLSchema } from "graphql";
import { Maybe } from "graphql/jsutils/Maybe";
import { createSchema } from "../common/createSchema";

interface Options {
    source: string;
    variableValues?: Maybe<{
        [key: string]: any;
    }>;
}

let schema: GraphQLSchema;

export const gCall = async ({
    source,
    variableValues,
}: Options): Promise<ExecutionResult> => {
    if (!schema) {
        schema = await createSchema();
    }
    return graphql({
        schema,
        source,
        variableValues,
    });
};
