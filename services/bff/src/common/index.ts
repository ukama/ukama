import { axiosErrorHandler } from "../errors";
import { ApiMethodDataDto, Context, ParsedCookie } from "./types";
import ApiMethods from "../api";
import { converCookieToObj } from "../utils";

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

export const parseCookie = (ctx: Context): ParsedCookie => {
    let header = {};
    const cookieObj: any = converCookieToObj(ctx.cookie);
    if (ctx.token) {
        header = {
            Authorization: ctx.token,
        };
    } else if (ctx.cookie) {
        header = {
            Cookie: "ukama_session=MTY1Mjc5Nzg2M3w3TVRyTEUyUEFXcHJuOWNEWm5iN0tmeno2LXJWdnZWSGRPdlF6WVRQZE5BWDcycktWUnc0cFIwNjBmTG9zWkZCUE5TOWtDQzFVZWRkV2hQcXBJWG9RYkRYNHFheVp4alNnM19RbUdDV0lYNnVfekdQSVQyaWd5LUl6WTFlZnl4LU9PTXRUTnktZ1E9PXy32OcVwiQGkXkVfrypzqk-EsX8oNG8JQ10gqAu0su_gA==",
        };
    }
    return { header: header, orgId: "a32485e4-d842-45da-bf3e-798889c68ad0" };
};
