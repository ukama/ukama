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

    // return { header: header, orgId: cookieObj["id"] };
    return {
        header: {
            Cookie: "ukama_session=MTY1Njk1MjA3M3xvMGt4R1dPcjEwamdTdDdlVkttNTFCak5QZmpackprb0FsalNLMXZtaFhJcmFJLXVpTEZmbXpyekl1LTQycTJJcWUwN3hFc21BazNZVmgtdDZGZWdub013SXVkUEhSYmdUeEJhS0ZaSXZyZDN0VGk1aGNabXQ3UDU4ZUJHNjNNQVVXaHVEenE4NkE9PXzXLvIyD-BEuWF5mEQGbTSZVD_eDcfc4kseLI5PoxOT_A==",
        },
        orgId: "a32485e4-d842-45da-bf3e-798889c68ad0",
    };
};
