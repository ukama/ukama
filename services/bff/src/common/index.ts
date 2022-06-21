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
            Cookie: "ukama_session=MTY1NTc0MDQyOHxBRldBWkpRTm93N21JVkNWNjBoR0lqZEg2N3NHSkpLcnNPN3NWaV9SU0I3WFNZZlRHLUhLRFEtTTB2M2hRUHZDLXlOc2pGMWZmQkhCU0s3QjZzOFZpc0szUkZCOFVYMjRYY29Bc2xTQ0RWWmN3X2g5QklOMEo1cU0yYzczZW0zbXhhSFdIVDdsTWc9PXzyMqp_jQPabM8ke3kVqZ8yUgphK8ir21t3lsGOyibzOw==",
        },
        orgId: "a32485e4-d842-45da-bf3e-798889c68ad0",
    };
};
