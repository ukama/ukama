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
            Cookie: `ukama_session=${cookieObj["ukama_session"]}`,
        };
    }
    // return { header: header, orgId: cookieObj["orgId"] };
    return {
        header: {
            Cookie: "ukama_session=MTY1MDA1MTg4MHxXZm91Yi1KN3U1aHhDRkNKT1d6RFFKT2lLWjBKX3BFVHpGUDJJVHF6Z2Q0OWRvZnRqekZUWnlZYWdiSm5SVlRMeGp6MFlHYXM3MkdXc1JHYmJpQjN1NXlweDZVd191Rko1TTdTUHRjcVdmQXRhOVlYX3puWVZQUW5NS19HTkl5MXFrSVhlZDVSNWc9PXxRj49ElI307naAjVwazS5OZSrdbB7W7yWWf00SrCzAdw==",
        },
        orgId: "a32485e4-d842-45da-bf3e-798889c68ad0",
    };
};
