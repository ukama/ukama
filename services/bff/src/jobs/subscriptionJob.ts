import "reflect-metadata";
import schedule from "node-schedule";
import {
    GET_ALERTS_QUERY,
    GET_CONNECTED_USERS_QUERY,
    GET_NETWORK_QUERY,
} from "../common/graphql";

import { HEADER, NETWORK_TYPE, TIME_FILTER } from "../constants";
import { graphql, GraphQLSchema } from "graphql";
import { PaginationDto } from "../common/types";
import setupLogger from "../config/setupLogger";

const logger = setupLogger("Job");
const rule = new schedule.RecurrenceRule();
rule.second = [0, 10, 20];

export const job = (schema: GraphQLSchema): void => {
    schedule.scheduleJob(rule, async function () {
        const meta: PaginationDto = {
            pageNo: 1,
            pageSize: 3,
        };

        await graphql({
            schema,
            source: GET_CONNECTED_USERS_QUERY,
            variableValues: {
                data: TIME_FILTER.WEEK,
            },
            contextValue: {
                req: HEADER,
            },
        });
        await graphql({
            schema,
            source: GET_ALERTS_QUERY,
            variableValues: {
                input: meta,
            },
            contextValue: {
                req: HEADER,
            },
        });
        logger.info(`Job Completed`);
    });
};
