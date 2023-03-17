import { ExecutionResult, graphql, GraphQLSchema } from "graphql";
import { Maybe } from "graphql/jsutils/Maybe";
import { createSchema } from "./createSchema";
import nock from "nock";
import { BASE_URL } from "../constants";

interface Options {
    source: string;
    variableValues?: Maybe<{
        [key: string]: any;
    }>;
    contextValue?: any;
}

let schema: GraphQLSchema;

export const gCall = async ({
    source,
    variableValues,
    contextValue,
}: Options): Promise<ExecutionResult> => {
    if (!schema) {
        schema = await createSchema();
    }
    return graphql({
        schema,
        source,
        variableValues,
        contextValue,
    });
};

export const beforeEachGetCall = (
    path: string,
    response: Object,
    responseHttpCode: number,
): void => {
    beforeEach(() => {
        nock(BASE_URL).get(path).reply(responseHttpCode, response);
    });
};

export const beforeEachPostCall = (
    path: string,
    body: any,
    response: Object,
    responseHttpCode: number,
): void => {
    beforeEach(() => {
        nock(BASE_URL).post(path, body).reply(responseHttpCode, response);
    });
};
