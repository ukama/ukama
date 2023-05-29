import ApiMethods from "../api";
import { HTTP401Error, Messages, axiosErrorHandler } from "../errors";
import { convertCookieToObj } from "../utils";
import { ApiMethodDataDto, Context, THeaders } from "./types";

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

export const parseHeaders = (ctx: Context): THeaders => {
    const headers: THeaders = {
        auth: {
            Cookie: "ukama_session=MTY4NDkwNzgyN3xCR0lyREZ2WWNPekNTVXhOd0g0OXh2QWZFelhUcVJPNEw4SGRVdnNWYTNaYXE1Z2tNWEVGUlRNRGhtaXFuSUxodk5WWDBrOG1IdTQ2OURGOTMzR3U3ejNORTBJZ2J3Tk5OSWdRTGMwcGR5Nm5lUHc5N3VLcVFXdGxZbUpneHhtTzE0V3FmU3JQTzJfdGFsU1pOMnFDTTU1ZGozSkNFcmdYTy0tVzNXZFlkWGdkeWh2T2dhbUZIVnNTamtCdzNtUGJ4VmVBOUN3YzZtNlk2bnZJdUVvS2N4eU03Q0NzX1Z2azJFcjhVUmJuZ2hvdk5NSk5HTDBEbS00djlJdVFBSHZDdFJoQmhoXzNCLW90emlOU1hpNXNiQT09fJLjUeEPUYkcsU22rg3TOGP_gcNMacpKQ-zlAvb3U-PC",
            Authorization: "",
        },
        orgId: "bf184df7-0ce6-4100-a9c6-497c181b87cf",
        userId: "a9a3dc45-fe06-43d6-b148-7508c9674627",
        orgName: "ukama",
    };
    const orgId = ctx.req.headers["org-id"];
    const userId = ctx.req.headers["user-id"];
    const orgName = ctx.req.headers["org-name"];
    if (!orgId || !userId || !orgName)
        throw new HTTP401Error(Messages.ERR_REQUIRED_HEADER_NOT_FOUND);
    else {
        headers.orgId = orgId as string;
        headers.userId = userId as string;
        headers.orgName = orgName as string;
    }
    if (ctx.authType === "token") {
        headers.auth.Authorization = ctx.req.headers[
            "x-session-token"
        ] as string;
    } else {
        const cookieObj: any = convertCookieToObj(
            ctx.req.headers["cookie"] || ""
        );
        headers.auth.Cookie = `ukama_session=${cookieObj["ukama_session"]}`;
    }
    return headers;
};
