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
