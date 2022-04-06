import { axiosErrorHandler } from "../errors";
import { ApiMethodDataDto, Context, HeaderType } from "./types";
import ApiMethods from "../api";

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

export const getHeaders = (ctx: Context): HeaderType => {
    let header = {};
    if (ctx.token) {
        header = {
            Authorization: ctx.token,
        };
    } else if (ctx.cookie) {
        header = {
            cookie: "ukama_session=MTY0OTIzMjU4MHxEdi1CQkFFQ180SUFBUkFCRUFBQVJfLUNBQUVHYzNSeWFXNW5EQThBRFhObGMzTnBiMjVmZEc5clpXNEdjM1J5YVc1bkRDSUFJR0ZyT1doNFQyTmxhVTlFU1RoS1NrVmFibWgzU0ROWFMwbFhOSGx4VkhGWXwl4p8XeNJc5_jXjJUXuur4OG7MVmggIG2xNKkvmyF4kw==",
        };
    }
    return header;
};
