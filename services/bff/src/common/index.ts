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
            Cookie: "ukama_session=MTY1NjY3NTU4N3xyRjkzLUZGTEYzaGVkamJVVmJXanNrbjQ5YzlTZllad292OEt6SDVNLUsxNUF3Nks4c1VISFpYMDRTWExFazRXRHFEU1U4ZnNkdF9YbEVSd21SS2Q2b0kyUEUwajY1elRnMjd5VTRfYnZwQVRBN1o2SWw4dWoyTEtpbWFpQ1JiN1JKenBuTE9yMEE9PXwwKhXp-1c_swFCbwvGTJ4uptxaY1o16v23j6_z29LFDA==",
        },
        orgId: "a32485e4-d842-45da-bf3e-798889c68ad0",
    };
};
