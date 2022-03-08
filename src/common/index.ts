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
    const header = { Authorization: "Bearer ZCa3ktK4Q3KHBxBXmTGyqJj3QCfI2bI3" };
    // if (ctx.token) {
    //     header = {
    //         Authorization: ctx.token,
    //     };
    // } else if (ctx.cookie) {
    //     header = {
    //         cookie: "ukama_session=MTY0Njc0MjE2NXxEdi1CQkFFQ180SUFBUkFCRUFBQVJfLUNBQUVHYzNSeWFXNW5EQThBRFhObGMzTnBiMjVmZEc5clpXNEdjM1J5YVc1bkRDSUFJRFk0TTFJMVpUUlBWa2hUZFVGR2NsSjBZVlpxUmtoeGVYWkZNRWxKTVdOTHwAqg6OTKxDME3MDaoOixFxzyb8q6fLcW15GaUAJbSUgQ==", //ctx.cookie,
    //     };
    // }

    return header;
};
