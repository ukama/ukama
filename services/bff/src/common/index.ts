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
            Cookie: "ukama_session=MTY1NjQzMjIwMXxFMGdBZjdOdHlfanA2VVFuZ28yMWZSeDREakpuWWg2OHRfckthR3ljbTRwNTE0RjF4MTNEcXh5WHl4Z2QyZWFBQlYxWWRVa2lzd0NHY2x0LVg4NkMxWldKYW5QUEhoNWtyOWJOQnp3SWx0cE8xWUFrUmlsdzVDSlFqM3JvVUltVlpyR2hOM2ZyU2c9PXzJlTfmefmDORvuugdKja8-xynaIlOLhM7I5Y_SISm7dA==",
        },
        orgId: "d0a36c51-6a66-4187-b786-72a9e09bf7a4",
    };
};
