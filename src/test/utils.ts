import { ExecutionResult, graphql, GraphQLSchema } from "graphql";
import { Maybe } from "graphql/jsutils/Maybe";
import { createSchema } from "../common/createSchema";
import nock from "nock";

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

export const nockCall = (path: string, response: Object): void => {
    nock("http://localhost:8081").get(path).reply(200, response);
};
