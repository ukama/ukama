import { axiosErrorHandler } from "../errors";
import { ApiMethodDataDto } from "./types";
import ApiMethods from "../api";

export const catchAsyncIOMethod = async <Output>(
    req: ApiMethodDataDto
): Promise<Output> => {
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
        const err = axiosErrorHandler(error);
        throw new Error(err.message);
    }
};
