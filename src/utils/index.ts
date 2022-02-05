import { Meta } from "../common/types";
import { GRAPH_FILTER } from "../constants";

export const getPaginatedOutput = (
    page: number,
    pageSize: number,
    count: number
): Meta => {
    return {
        count,
        page: page ? page : 1,
        size: pageSize ? pageSize : count,
        pages: pageSize ? Math.ceil(count / pageSize) : 1,
    };
};

export const getUniqueTimeStamp = (index?: number, length?: number): number =>
    Math.floor(new Date().valueOf() / 1000) -
    (length ? length - (index || 1) : 0);

export const getRecordsLengthByFilter = (
    filter: string | undefined
): number => {
    switch (filter) {
        case GRAPH_FILTER.WEEK:
            return 50;
        case GRAPH_FILTER.MONTH:
            return 100;
        default:
            return 10;
    }
};
