import ApiMethods from "../api";
import { axiosErrorHandler } from "../errors";
import { converCookieToObj } from "../utils";
import { ApiMethodDataDto, Context, ParsedCookie } from "./types";

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
            Cookie: `ukama_session=${cookieObj["ukama_session"]}`,
        };
    }

    return {
        header: header,
        orgId: cookieObj["org_id"],
        orgName: cookieObj["org_name"],
        userId: cookieObj["user_id"],
    };
};
