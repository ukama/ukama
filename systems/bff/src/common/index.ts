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
            Cookie: "ukama_session=MTY4NTk1NDYxOXxkRWpsRnZtbVFBZ1lHV0Jnc1ppeDFJYnhuMWtsbTJGeERRNVZFRWJrakxsdk5MN2ptYjA3UEdEWXkzSklzZHlRVmxEdTBjcElMQmswUFNIMzNHTTNYYkJCZV81R2tSRG5UTUFxem9qQzRzUXpwdkNaQUJTSmlpSkRrbGFRMGNKSHIxd3VrXzdFTlFkWEhITXpueFFaekctdW5paDRXZDJ0aGtuLWpobU1LTmFSQlRvbEk3WE5YWFRSc1k0OV9JaUd2TFBCVkFySHZSVDlrR2lRWkJFQThPUURtdjlCYnBYeU9XQkNHLTFrcEdwWEZQamdfWGFpcEd1MnM5dGxFdk1xT1FPd0ZsWXFJNW1taVptYW44Q3JUQT09fGG25BGdPznBtGqyNf65yGWgB_tfIyOHMHhbnP_ITXEo",
            Authorization: "",
        },
        orgId: "aac1ed88-2546-4f9c-a808-fb9c4d0ef24b",
        userId: "851e0abc-7e91-4206-8c85-498e16f91e67",
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
