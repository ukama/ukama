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
            Cookie: "ukama_session=MTY1MzI0NjAyNXxxcW5zbTE2eXdfZGxLS29mRW5MNTdBYkVnOFVEM1Z0NEdQSlp2UWVvM2RKNVJhMzdvY1Z2eGFrcXRVY2hEV1cxdk9wWkg5VkhGdUkxeF9zTDZYWld1ODRhZGFVbDNlSDgxWjlXTEVhdGlqTS1Yek1pdFR1Rm1aQlRRSnVILUZPR0JHeDVtUnNxSVE9PXzzIoGbNSQKLzGvG2UAyq_QHHfElLouB7RxY3vsbBQGWA==",
        };
    }
    return { header: header, orgId: "d89d860d-fa7d-4523-b460-f767fd48c623" };
};
