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
            Cookie: "ukama_session=MTY1MTY1NDUwOXxhbmxJa1NDOHd4UGhkbWZ1N21sQmRoLTcwY0x2VGdBZDNLOUJUZTBYZ1U2bW96cElyRmN0YUw0LTd2OGg0YzlZY05FNFNFcDdCemhIdU5taXBtN09oclUwTjZId0RFN1lpNVA1WlZ3VUxYV09lTGY5OFNYaWlJZkZ1QlBEeXMyb21ZeG9ubkdTblE9PXwRhHbwoNt9JSQ5CRz7D1HSk8uv3bjT9N25AhwyiPH83A==",
        };
    }
    return { header: header, orgId: "a32485e4-d842-45da-bf3e-798889c68ad0" };
};
