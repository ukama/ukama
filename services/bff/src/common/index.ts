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
            Cookie: "ukama_session=MTY1MjgxMjczOXxpc2E5ZGI4RE1WUkZWTTNZb3NCUFQ0MTBuVGEyZWs5NjFJc1l2a1pjb04zdUJjZ1M2b1JqcVZibnVjMnRkeV9GRVpHY3k0SVhnRE80VmFOQzg1SzBRbTVEblNEUGtEdTJBTmJWTFpGQU9JT3l1NFJrTV9mV1lVUjM0WEltWkdZRlczeVk5a3NrZlE9PXxmdZfBLdfx3DgGNwyov6mp41KGEK4bbQmVYwgP_FHL0A==",
        };
    }
    return { header: header, orgId: "a32485e4-d842-45da-bf3e-798889c68ad0" };
};
