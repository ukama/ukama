import { axiosErrorHandler } from "../errors";
import { ApiMethodDataDto, Context, HeaderType } from "./types";
import ApiMethods from "../api";

export const catchAsyncIOMethod = async (
    req: ApiMethodDataDto
): Promise<any> => {
    try {
        const res = await ApiMethods.fetch({
            type: req.type,
            path: req.path,
            params: req.params,
            body: req.body,
            headers: req.headers,
        });

        return res.data;
    } catch (error) {
        return axiosErrorHandler(error);
    }
};

export const getHeaders = (ctx: Context): HeaderType => {
    let header = {};
    if (ctx.token) {
        header = {
            Authorization: ctx.token,
        };
    } else if (ctx.cookie) {
        header = {
            Cookie: ctx.cookie,
        };
    }

    return header;
};
