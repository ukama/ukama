import { Meta } from "../common/types";

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

export const createPaginatedResponse = (
    pageNo: number,
    pageSize: number,
    data: any[]
): any => {
    let metrics = [];
    if (!pageNo) metrics = data;
    else {
        const index = (pageNo - 1) * pageSize;
        for (let i = index; i < index + pageSize; i++) {
            if (data[i]) metrics.push(data[i]);
        }
    }
    return metrics;
};
