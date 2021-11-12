import { ExecutionResult, graphql, GraphQLSchema } from "graphql";
import { Maybe } from "graphql/jsutils/Maybe";
import { createSchema } from "../common/createSchema";
import nock from "nock";
import env from "../config/env";

export const BASE_URL = env.BASE_URL;

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
    nock(BASE_URL).get(path).reply(200, response);
};
