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
            Cookie: "ukama_session=MTY1NjY4ODE4MnxCV3J6RC0tUnBkNGVLV2k0ejE0T1ROTkZEbmNzdWc3aTk0NHBrcUE4M05xN1JEU25hcFRaZXZxbHk4eFowNkVnWEhxR1phc0RPdHJXRUhHN1JGX2ZSOFpqMS16YV9yVzA1MHZ0b1VUcWRYbXlSaUtrVHJRVEdWcEpwU3VsT2FSc0w2VkRJOXp5NXc9PXzWzPtKEVwlPm9dG0IFv6R4Bm_5Q4Vd5wNtrpUUlY5f4A==",
        },
        orgId: "a32485e4-d842-45da-bf3e-798889c68ad0",
    };
};
