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
            Cookie: "ukama_session=MTY2MjczODQ3MnxiNVJaTEo2YzRvNi1GdzdHVVVZSzVPOW54WHhuRVpnNlhfbEJOeG9nMmVGM05GWGh5bkRrV0hGTHIzZ29iWHBLcDNTcHVFVVFOeHo3bm9RMU8wcm5fMk5reDJMWUJKbDNJTDdJb0VKZm9oa21pdkNpQ3dOUXE0X0s2NUpRM1hQYmlvVDVJcVN3eVE9PXysbUDw8bXvLodtEHkPZwZJRi3tVEv33y37B6zWHHTx7Q==",
        },
        orgId: "d89d860d-fa7d-4523-b460-f767fd48c623",
    };
};
