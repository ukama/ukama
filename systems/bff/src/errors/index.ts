import { ErrorType } from "../common/types";
import { HttpStatusCode } from "./codes";
import { HTTP400Error } from "./http400.error";
import { HTTP401Error } from "./http401.error";
import { HTTP404Error } from "./http404.error";
import { HTTP500Error } from "./http500.error";
import Messages from "./messages";

export {
    HTTP400Error,
    HTTP500Error,
    HTTP401Error,
    HttpStatusCode,
    HTTP404Error,
    Messages,
};

export const axiosErrorHandler = (error: any): ErrorType => {
    let res: ErrorType;
    if (error.response) {
        // The request was made and the server responded with a status code
        if (error.response.data) {
            return {
                code: error.response.status,
                message: error.response.data.message,
                description: error.response.statusText,
            };
        } else {
            return {
                code: error.response.status,
                message: error.response.statusText,
                description: "",
            };
        }
    } else if (error.request) {
        // The request was not made

        res = {
            code: 400,
            message: Messages.ERR_SERVER_REQUEST_FAILED,
        };
    } else {
        // Something happened in setting up the request that triggered an Error

        res = {
            code: 400,
            message: error.message,
        };
    }
    return res;
};

export const checkError = (error: any): boolean => {
    if (error.code || error.message || error.description) return true;
    return false;
};
